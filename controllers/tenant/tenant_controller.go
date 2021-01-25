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
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"github.com/spaghettifunk/pinot-go-client/client/tenant"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"github.com/spaghettifunk/pinot-operator/pkg/sdk"
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

	pinotsdk "github.com/spaghettifunk/pinot-go-client/client"
	"github.com/spaghettifunk/pinot-go-client/models"
	operatorsv1alpha1 "github.com/spaghettifunk/pinot-operator/api/pinot/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller")

const finalizerID = "pinot-tenant.finalizer.apache.io"

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
	Log         logr.Logger
	Scheme      *runtime.Scheme
	PinotClient *pinotsdk.PinotSdk
}

// +kubebuilder:rbac:groups=operators.apache.io,resources=tenants;tenants/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operators.apache.io,resources=tenants/status,verbs=get;update;patch

func (r *ReconcilerTenant) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	logger := log.WithValues("trigger", request.Namespace+"/"+request.Name, "correlationID", uuid.Must(uuid.NewV4()).String())

	// Fetch the Tenant instance
	config := &operatorsv1alpha1.Tenant{}
	err := r.Get(context.TODO(), request.NamespacedName, config)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logger.Info("Reconciling Tenant")

	pinot, err := r.getRelatedPinotCR(config)
	if err != nil {
		updateErr := updateStatus(r.Client, config, operatorsv1alpha1.ReconcileFailed, err.Error(), logger)
		if updateErr != nil {
			logger.Error(updateErr, "failed to update state")
			return reconcile.Result{}, errors.WithStack(err)
		}
		return reconcile.Result{
			Requeue: false,
		}, errors.WithStack(err)
	}
	operatorsv1alpha1.SetPinotDefaults(pinot)

	// if the pinot sdk is not initialized, do it now
	if r.PinotClient == nil {
		r.PinotClient = pinotsdk.NewHTTPClientWithConfig(strfmt.Default, &pinotsdk.TransportConfig{
			Host:     sdk.GeneratePinotControllerAddress(pinot),
			BasePath: pinotsdk.DefaultBasePath,
			Schemes:  []string{"http"},
		})
	}

	result, err := r.reconcile(logger, config)
	if err != nil {
		updateErr := updateStatus(r.Client, config, operatorsv1alpha1.ReconcileFailed, err.Error(), logger)
		if updateErr != nil {
			logger.Error(updateErr, "failed to update state")
			return result, errors.WithStack(err)
		}
		return result, emperror.Wrap(err, "could not reconcile Tenant")
	}
	return result, nil
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

func (r *ReconcilerTenant) reconcile(logger logr.Logger, config *operatorsv1alpha1.Tenant) (reconcile.Result, error) {
	if config.Status.Status == "" {
		err := updateStatus(r.Client, config, operatorsv1alpha1.Created, "", logger)
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
			if config.Status.Status == operatorsv1alpha1.Reconciling && config.Status.ErrorMessage == "" {
				logger.Info("cannot remove Tenant while reconciling")
				return reconcile.Result{}, nil
			}
			if err := r.deleteTenantResources(config); err != nil {
				return ctrl.Result{}, emperror.Wrap(err, "could not remove tenant")
			}

			config.ObjectMeta.Finalizers = util.RemoveString(config.ObjectMeta.Finalizers, finalizerID)
			if err := r.Update(context.Background(), config); err != nil {
				return reconcile.Result{}, emperror.Wrap(err, "could not remove finalizer from config")
			}
		}
		logger.Info("Tenant removed")
		return reconcile.Result{}, nil
	}

	err := updateStatus(r.Client, config, operatorsv1alpha1.Reconciling, "", logger)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	// upsert tenant
	if err := r.upsertTenantResource(config); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	err = updateStatus(r.Client, config, operatorsv1alpha1.Available, "", logger)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}
	return reconcile.Result{}, nil
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

func (r *ReconcilerTenant) upsertTenantResource(config *operatorsv1alpha1.Tenant) error {
	role := strings.ToUpper(config.Spec.Role)

	// get tenant metadata
	res, err := r.PinotClient.Tenant.GetTenantMetadata(&tenant.GetTenantMetadataParams{
		Type:    util.StrPointer(role),
		Context: context.Background(),
	})
	if _, ok := err.(*tenant.GetTenantMetadataNotFound); !ok {
		return err
	}

	// if tenant exists, update it
	if res != nil && (len(res.Payload.BrokerInstances) > 0 || len(res.Payload.ServerInstances) > 0) {
		_, err = r.PinotClient.Tenant.UpdateTenant(&tenant.UpdateTenantParams{
			Body: &models.Tenant{
				TenantRole:        role,
				TenantName:        config.Spec.Name,
				NumberOfInstances: util.PointerToInt32(config.Spec.NumberOfInstances),
				OfflineInstances:  util.PointerToInt32(config.Spec.OfflineInstances),
				RealtimeInstances: util.PointerToInt32(config.Spec.RealtimeInstances),
			},
			Context: context.Background(),
		})
	} else {
		// create the new tenant
		_, err = r.PinotClient.Tenant.CreateTenant(&tenant.CreateTenantParams{
			Body: &models.Tenant{
				TenantRole:        role,
				TenantName:        config.Spec.Name,
				NumberOfInstances: util.PointerToInt32(config.Spec.NumberOfInstances),
				OfflineInstances:  util.PointerToInt32(config.Spec.OfflineInstances),
				RealtimeInstances: util.PointerToInt32(config.Spec.RealtimeInstances),
			},
			Context: context.Background(),
		})
	}
	return err
}

func (r *ReconcilerTenant) deleteTenantResources(config *operatorsv1alpha1.Tenant) error {
	_, err := r.PinotClient.Tenant.DeleteTenant(&tenant.DeleteTenantParams{
		TenantName: config.Spec.Name,
		Type:       strings.ToUpper(config.Spec.Role),
		Context:    context.Background(),
	})
	return err
}
