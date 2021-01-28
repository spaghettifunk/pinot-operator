/*
Copyright 2021 the Apache Pinot Kubernetes Operator authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	_ "github.com/shurcooL/vfsgen"
	"k8s.io/apimachinery/pkg/runtime"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/spaghettifunk/pinot-operator/controllers/pinot"
	"github.com/spaghettifunk/pinot-operator/controllers/schema"
	"github.com/spaghettifunk/pinot-operator/controllers/table"
	"github.com/spaghettifunk/pinot-operator/controllers/tenant"
	api "github.com/spaghettifunk/pinot-operator/pkg/apis"
	operatorsv1alpha1 "github.com/spaghettifunk/pinot-operator/pkg/apis/pinot/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	// +kubebuilder:scaffold:imports
)

const watchNamespaceEnvVar = "WATCH_NAMESPACE"
const podNamespaceEnvVar = "POD_NAMESPACE"

var (
	scheme                 = runtime.NewScheme()
	setupLog               = ctrl.Log.WithName("setup")
	shutdownWaitDuration   = time.Duration(30) * time.Second
	waitBeforeExitDuration = time.Duration(3) * time.Second
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = operatorsv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	namespace, err := getWatchNamespace()
	if err != nil {
		setupLog.Error(err, "")
		os.Exit(1)
	}
	if namespace != "" {
		setupLog.Info("watch namespace", "namespace", namespace)
	} else {
		setupLog.Info("watch all namespaces")
	}

	mgr, err := manager.New(ctrl.GetConfigOrDie(), manager.Options{
		Scheme:             scheme,
		Namespace:          namespace,
		MetricsBindAddress: metricsAddr,
		MapperProvider:     k8sutil.NewCachedRESTMapper,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "366fcac3.apache.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info("registering components")

	// Setup Scheme for all resources
	setupLog.Info("setting up scheme")
	if err := api.AddToScheme(mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	stop := setupSignalHandler(mgr, setupLog, shutdownWaitDuration)

	// Setup all Controllers
	setupLog.Info("setting up controller")

	// pinot cluster controller
	if err := pinot.Add(mgr); err != nil {
		setupLog.Error(err, "problem adding manager")
		os.Exit(1)
	}

	// tenants controller
	if err := tenant.Add(mgr); err != nil {
		setupLog.Error(err, "problem adding manager")
		os.Exit(1)
	}

	// schema controller
	if err := schema.Add(mgr); err != nil {
		setupLog.Error(err, "problem adding manager")
		os.Exit(1)
	}

	// table controller
	if err := table.Add(mgr); err != nil {
		setupLog.Error(err, "problem adding manager")
		os.Exit(1)
	}

	// readines and liveness
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(stop); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

	// Wait a bit for the workers to stop
	time.Sleep(waitBeforeExitDuration)

	// Cleanup
	setupLog.Info("removing finalizer from Pinot resources")
	err = pinot.RemoveFinalizers(mgr.GetClient())
	if err != nil {
		setupLog.Error(err, "could not remove finalizers from Pinot resources")
	}

	setupLog.Info("removing finalizer from Tenant resources")
	err = tenant.RemoveFinalizers(mgr.GetClient())
	if err != nil {
		setupLog.Error(err, "could not remove finalizers from Tenant resources")
	}
}

func getWatchNamespace() (string, error) {
	podNamespace, found := os.LookupEnv(podNamespaceEnvVar)
	if !found {
		return "", errors.Errorf("%s env variable must be specified and cannot be empty", podNamespaceEnvVar)
	}
	watchNamespace, found := os.LookupEnv(watchNamespaceEnvVar)
	if found {
		if watchNamespace != "" && watchNamespace != podNamespace {
			return "", errors.New("watch namespace must be either empty or equal to pod namespace")
		}
	}
	return watchNamespace, nil
}

func setupSignalHandler(mgr manager.Manager, log logr.Logger, shutdownWaitDuration time.Duration) (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("termination signal arrived, shutdown gracefully")
		// wait a bit for deletion requests to arrive
		log.Info("wait a bit for CR deletion events to arrive", "waitSeconds", shutdownWaitDuration)
		time.Sleep(shutdownWaitDuration)
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
