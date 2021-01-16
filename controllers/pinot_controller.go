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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/pkg/errors"

	pinotbroker "github.com/spaghettifunk/pinot-operator/pkg/resources/broker"
	pinotcontroller "github.com/spaghettifunk/pinot-operator/pkg/resources/controller"
	pinotserver "github.com/spaghettifunk/pinot-operator/pkg/resources/server"
	pinotzookeeper "github.com/spaghettifunk/pinot-operator/pkg/resources/zookeeper"
	corev1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/resources"
	"github.com/spaghettifunk/pinot-operator/pkg/util"

	clusterv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"
)

const finalizerID = "pinot-operator.finalizer.apache.io"

var log = logf.Log.WithName("controller")
var watchCreatedResourcesEvents bool

// Add creates a new Pinot Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &PinotReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("pinot-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Pinot
	err = c.Watch(&source.Kind{Type: &pinotv1alpha1.Pinot{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Pinot
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &pinotv1alpha1.Pinot{},
	})
	if err != nil {
		return err
	}
	return nil
}

// PinotReconciler reconciles a Pinot object
type PinotReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=operators.apache.io,resources=pinots,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operators.apache.io,resources=pinots/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operators,resources=configmaps;statefulsets;services;secrets;poddisruptionbudgets,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups=policy;apps,resources=poddisruptionbudgets;statefulsets,verbs=*
// +kubebuilder:rbac:groups="",resources=events;statefulsets;configmaps;services;poddisruptionbudgets,verbs=get;list;watch;create;update;delete

func (r *PinotReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	// Fetch the Pinot instance
	config := &pinotv1alpha1.Pinot{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, config)
	if err != nil {
		if k8errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Pinot resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get Pinot")
		return reconcile.Result{}, err
	}

	logger.Info("Reconciling Pinot")

	// start reconciling loop
	result, err := r.reconcile(logger, config)
	if err != nil {
		updateErr := updateStatus(r.Client, config, pinotv1alpha1.ReconcileFailed, err.Error(), logger)
		if updateErr != nil {
			logger.Error(updateErr, "failed to update state")
			return result, errors.WithStack(err)
		}
		return result, emperror.Wrap(err, "could not reconcile Pinot")
	}
	return result, nil
}

func (r *PinotReconciler) reconcile(logger logr.Logger, config *pinotv1alpha1.Pinot) (reconcile.Result, error) {
	if config.Status.Status == "" {
		err := updateStatus(r.Client, config, pinotv1alpha1.Created, "", logger)
		if err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	// for each component do a reconciliation
	reconcilers := []resources.ComponentReconciler{
		pinotzookeeper.New(r.Client, config),
		pinotserver.New(r.Client, config),
		pinotcontroller.New(r.Client, config),
		pinotbroker.New(r.Client, config),
	}
	for _, rec := range reconcilers {
		err := rec.Reconcile(logger)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	err := updateStatus(r.Client, config, pinotv1alpha1.Available, "", logger)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	logger.Info("reconcile finished")

	return reconcile.Result{}, nil
}

func updateStatus(c client.Client, config *pinotv1alpha1.Pinot, status pinotv1alpha1.ConfigState, errorMessage string, logger logr.Logger) error {

	typeMeta := config.TypeMeta
	config.Status.Status = status
	config.Status.ErrorMessage = errorMessage

	err := c.Status().Update(context.Background(), config)
	if k8errors.IsNotFound(err) {
		err = c.Update(context.Background(), config)
	}

	if err != nil {
		if !k8errors.IsConflict(err) {
			return emperror.Wrapf(err, "could not update Pinot state to '%s'", status)
		}

		var actualConfig pinotv1alpha1.Pinot
		err := c.Get(context.TODO(), types.NamespacedName{
			Namespace: config.Namespace,
			Name:      config.Name,
		}, &actualConfig)

		if err != nil {
			return emperror.Wrap(err, "could not get config for updating status")
		}

		actualConfig.Status.Status = status
		actualConfig.Status.ErrorMessage = errorMessage

		err = c.Status().Update(context.Background(), &actualConfig)
		if k8errors.IsNotFound(err) {
			err = c.Update(context.Background(), &actualConfig)
		}
		if err != nil {
			return emperror.Wrapf(err, "could not update Pinot state to '%s'", status)
		}
	}

	// update loses the typeMeta of the config that's used later when setting ownerrefs
	config.TypeMeta = typeMeta
	logger.Info("Pinot state updated", "status", status)

	return nil
}

// RemoveFinalizers removes the finalizers from the context
func RemoveFinalizers(c client.Client) error {
	var pinots pinotv1alpha1.PinotList

	// fix this!
	// err := c.List(context.TODO(), &client.ListOptions{}, &pinots)
	// if err != nil {
	// 	return emperror.Wrap(err, "could not list Pinot resources")
	// }

	for _, pinot := range pinots.Items {
		pinot.ObjectMeta.Finalizers = util.RemoveString(pinot.ObjectMeta.Finalizers, finalizerID)
		if err := c.Update(context.Background(), &pinot); err != nil {
			return emperror.WrapWith(err, "could not remove finalizer from Pinot resource", "name", pinot.GetName())
		}
		if err := updateStatus(c, &pinot, pinotv1alpha1.Unmanaged, "", log); err != nil {
			return emperror.Wrap(err, "could not update status of Pinot resource")
		}
	}
	return nil
}

func (r *PinotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1alpha1.Pinot{}).
		Complete(r)
}
