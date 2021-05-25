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

package pinot

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/gofrs/uuid"
	"github.com/goph/emperror"
	"github.com/pkg/errors"

	"github.com/spaghettifunk/pinot-operator/pkg/crds"
	pinotbroker "github.com/spaghettifunk/pinot-operator/pkg/resources/broker"
	pinotcontroller "github.com/spaghettifunk/pinot-operator/pkg/resources/controller"
	pinotserver "github.com/spaghettifunk/pinot-operator/pkg/resources/server"
	pinotzookeeper "github.com/spaghettifunk/pinot-operator/pkg/resources/zookeeper"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/pkg/apis/pinot/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/resources"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
)

const finalizerID = "pinot-operator.finalizer.apache.io"

var log = logf.Log.WithName("controller")
var watchCreatedResourcesEvents bool

func init() {
	flag.BoolVar(&watchCreatedResourcesEvents, "watch-created-resources-events", true, "Whether to watch created resources events")
}

// Add creates a new Pinot Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	crd, err := crds.New(mgr, pinotv1alpha1.OperatorVersion)
	if err != nil {
		return emperror.Wrap(err, "unable to set up CRD reconciler")
	}
	err = crd.LoadCRDs()
	if err != nil {
		return emperror.Wrap(err, "unable to load CRDs from manifests")
	}
	r := newReconciler(mgr, crd)
	err = newController(mgr, r)
	if err != nil {
		return emperror.Wrap(err, "failed to create controller")
	}
	return nil
}

type PinotReconciler interface {
	reconcile.Reconciler
	initWatches(watchCreatedResourcesEvents bool) error
	setController(ctrl controller.Controller)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, crd *crds.CRDReconciler) reconcile.Reconciler {
	return &ReconcilePinot{
		Client:        mgr.GetClient(),
		CRDReconciler: crd,
		Manager:       mgr,
		Scheme:        mgr.GetScheme(),
		Recorder:      mgr.GetEventRecorderFor("pinot-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func newController(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	ctrl, err := controller.New("pinot-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	if r, ok := r.(PinotReconciler); ok {
		r.setController(ctrl)
		err = r.initWatches(watchCreatedResourcesEvents)
		if err != nil {
			return emperror.Wrapf(err, "could not init watches")
		}
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcilePinot{}

// ReconcilePinot reconciles a Pinot object
type ReconcilePinot struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Log              logr.Logger
	CRDReconciler    *crds.CRDReconciler
	Manager          manager.Manager
	Scheme           *runtime.Scheme
	Recorder         record.EventRecorder
	Ctrl             controller.Controller
	WatchersInitOnce sync.Once
}

type ReconcileComponent func(log logr.Logger, pinot *pinotv1alpha1.Pinot) error

// the rbac rule requires an empty row at the end to render
// +kubebuilder:rbac:groups=pinot.apache.io,resources=pinots,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups=pinot.apache.io,resources=pinots/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=pinot.apache.io,resources=pinots/finalizers,verbs=update
// +kubebuilder:rbac:groups=operators,resources=configmaps;statefulsets;services;secrets;poddisruptionbudgets,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups=policy;apps,resources=poddisruptionbudgets;statefulsets,verbs=*
// +kubebuilder:rbac:groups="",resources=events;statefulsets;configmaps;services;poddisruptionbudgets,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=roles,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=rolebindings,verbs=get;list;watch;create;update

// Reconcile reads that state of the cluster for a Config object and makes changes based on the state read
// and what is in the Config.Spec
func (r *ReconcilePinot) Reconcile(request ctrl.Request) (ctrl.Result, error) {

	logger := log.WithValues("trigger", request.Namespace+"/"+request.Name, "correlationID", uuid.Must(uuid.NewV4()).String())

	// Fetch the Config instance
	config := &pinotv1alpha1.Pinot{}
	err := r.Get(context.TODO(), request.NamespacedName, config)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logger.Info("Reconciling Pinot")

	// Set default values where not set
	pinotv1alpha1.SetPinotDefaults(config)

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

func (r *ReconcilePinot) setController(ctrl controller.Controller) {
	r.Ctrl = ctrl
}

func (r *ReconcilePinot) reconcile(logger logr.Logger, config *pinotv1alpha1.Pinot) (reconcile.Result, error) {
	if config.Status.Status == "" {
		err := updateStatus(r.Client, config, pinotv1alpha1.Created, "", logger)
		if err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	// add finalizer strings and update
	if config.ObjectMeta.DeletionTimestamp.IsZero() {
		if !util.ContainsString(config.ObjectMeta.Finalizers, finalizerID) {
			config.ObjectMeta.Finalizers = append(config.ObjectMeta.Finalizers, finalizerID)
			if err := r.Update(context.Background(), config); err != nil {
				return reconcile.Result{}, emperror.Wrap(err, "could not add finalizer to config")
			}
			return reconcile.Result{
				RequeueAfter: time.Second * 1,
			}, nil
		}
	} else {
		// Deletion timestamp set, config is marked for deletion
		if util.ContainsString(config.ObjectMeta.Finalizers, finalizerID) {
			if config.Status.Status == pinotv1alpha1.Reconciling && config.Status.ErrorMessage == "" {
				logger.Info("cannot remove Pinot while reconciling")
				return reconcile.Result{}, nil
			}
			config.ObjectMeta.Finalizers = util.RemoveString(config.ObjectMeta.Finalizers, finalizerID)
			if err := r.Update(context.Background(), config); err != nil {
				return reconcile.Result{}, emperror.Wrap(err, "could not remove finalizer from config")
			}
		}
		logger.Info("Pinot removed")
		return reconcile.Result{}, nil
	}

	err := updateStatus(r.Client, config, pinotv1alpha1.Reconciling, "", logger)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	// reconcile here
	logger.Info("reconciling CRDs")
	err = r.CRDReconciler.Reconcile(config, logger)
	if err != nil {
		logger.Error(err, "unable to reconcile CRDs")
		return reconcile.Result{}, err
	}

	r.WatchersInitOnce.Do(func() {
		nn := types.NamespacedName{
			Namespace: config.Namespace,
			Name:      config.Name,
		}
		err = r.watchCRDs(nn)
		if err != nil {
			logger.Error(err, "unable to watch CRDs")
		}
	})

	// reconcile Zookeeper first
	zkReconciler := pinotzookeeper.New(r.Client, config)
	err = zkReconciler.Reconcile(logger)
	if err != nil {
		return reconcile.Result{}, err
	}

	// wait for zookeeper
	if err := zkReconciler.WaitForCreation(); err != nil {
		logger.Info("zookeeper is not ready")
		return reconcile.Result{
			RequeueAfter: time.Second * 10,
		}, nil
	}

	// reconcile all the rest of the resources
	reconcilers := []resources.ComponentReconciler{
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

	err = updateStatus(r.Client, config, pinotv1alpha1.Available, "", logger)
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
	if k8serrors.IsNotFound(err) {
		err = c.Update(context.Background(), config)
	}

	if err != nil {
		if !k8serrors.IsConflict(err) {
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
		if k8serrors.IsNotFound(err) {
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
