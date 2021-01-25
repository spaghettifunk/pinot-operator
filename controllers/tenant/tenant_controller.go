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

package tenant

import (
	"context"
	"sort"

	"github.com/go-logr/logr"
	"github.com/gofrs/uuid"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"github.com/spaghettifunk/pinot-operator/pkg/resources/tenant"
	"github.com/spaghettifunk/pinot-operator/pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	operatorsv1alpha1 "github.com/spaghettifunk/pinot-operator/api/pinot/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilerTenant{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("tenants-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Tenant
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.Tenant{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Tenant",
			APIVersion: "operators.apache.io/v1alpha1",
		},
	},
	}, &handler.EnqueueRequestForObject{}, k8sutil.GetWatchPredicateForTenant())
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcilerTenant{}

// ReconcilerTenant reconciles a Tenant object
type ReconcilerTenant struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=operators.apache.io,resources=Tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operators.apache.io,resources=Tenants/status,verbs=get;update;patch

func (r *ReconcilerTenant) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	logger := log.WithValues("trigger", request.Namespace+"/"+request.Name, "correlationID", uuid.Must(uuid.NewV4()).String())

	// Fetch the Tenant instance
	t := &operatorsv1alpha1.Tenant{}
	err := r.Get(context.TODO(), request.NamespacedName, t)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	operatorsv1alpha1.SetTenantDefaults(t)

	pinot, err := r.getRelatedPinotCR(t)
	if err != nil {
		updateErr := updateStatus(r.Client, t, operatorsv1alpha1.ReconcileFailed, err.Error(), logger)
		if updateErr != nil {
			logger.Error(updateErr, "failed to update state")
			return reconcile.Result{}, errors.WithStack(err)
		}
		return reconcile.Result{
			Requeue: false,
		}, errors.WithStack(err)
	}
	operatorsv1alpha1.SetPinotDefaults(pinot)

	t.Spec.Labels = util.MergeStringMaps(t.Spec.Labels, pinot.RevisionLabels())

	if err := updateStatus(r.Client, t, operatorsv1alpha1.Reconciling, "", logger); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	reconciler := tenant.New(r.Client, pinot)
	if err := reconciler.Reconcile(log); err != nil {
		logger.Error(err, "failed to reconcile tenant")
		return reconcile.Result{}, errors.WithStack(err)
	}

	if err = updateStatus(r.Client, t, operatorsv1alpha1.Available, "", logger); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	return ctrl.Result{}, nil
}

func (r *ReconcilerTenant) getRelatedPinotCR(instance *operatorsv1alpha1.Tenant) (*operatorsv1alpha1.Pinot, error) {
	pinot := &operatorsv1alpha1.Pinot{}

	// try to get specified Pinot CR
	if instance.Spec.PinotServer != nil {
		err := r.Client.Get(context.Background(), client.ObjectKey{
			Name:      instance.Spec.PinotServer.Name,
			Namespace: instance.Spec.PinotServer.Namespace,
		}, pinot)
		if err != nil {
			return nil, emperror.Wrap(err, "could not get related Pinot CR")
		}
		return pinot, nil
	}

	// get the oldest otherwise for backward compatibility
	var configs operatorsv1alpha1.PinotList
	err := r.Client.List(context.TODO(), &configs)
	if err != nil {
		return nil, emperror.Wrap(err, "could not list pinot resources")
	}
	if len(configs.Items) == 0 {
		return nil, errors.New("no Pinot CRs were found")
	}

	sort.Sort(operatorsv1alpha1.SortablePinotItems(configs.Items))

	config := configs.Items[0]
	gvk := config.GroupVersionKind()
	gvk.Version = operatorsv1alpha1.GroupVersion.Version
	gvk.Group = operatorsv1alpha1.GroupVersion.Group
	gvk.Kind = "Pinot"
	config.SetGroupVersionKind(gvk)

	return &config, nil
}

func updateStatus(c client.Client, instance *operatorsv1alpha1.Tenant, status operatorsv1alpha1.ConfigState, errorMessage string, logger logr.Logger) error {
	typeMeta := instance.TypeMeta
	instance.Status.Status = status
	instance.Status.ErrorMessage = errorMessage

	err := c.Status().Update(context.Background(), instance)
	if k8serrors.IsNotFound(err) {
		err = c.Update(context.Background(), instance)
	}
	if err != nil {
		if !k8serrors.IsConflict(err) {
			return emperror.Wrapf(err, "could not update tenant state to '%s'", status)
		}
		var actualInstance operatorsv1alpha1.Tenant
		err := c.Get(context.TODO(), types.NamespacedName{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		}, &actualInstance)
		if err != nil {
			return emperror.Wrap(err, "could not get resource for updating status")
		}
		actualInstance.Status.Status = status
		actualInstance.Status.ErrorMessage = errorMessage
		err = c.Status().Update(context.Background(), &actualInstance)
		if k8serrors.IsNotFound(err) {
			err = c.Update(context.Background(), &actualInstance)
		}
		if err != nil {
			return emperror.Wrapf(err, "could not update tenant state to '%s'", status)
		}
	}

	// update loses the typeMeta of the instace that's used later when setting ownerrefs
	instance.TypeMeta = typeMeta
	logger.Info("tenant state updated", "status", status)

	return nil
}
