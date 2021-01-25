package pinot

import (
	"errors"

	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/pinot/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/crds"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func (r *ReconcilePinot) initWatches(watchCreatedResourcesEvents bool) error {
	if r.Ctrl == nil {
		return errors.New("controller is not set")
	}

	var err error
	for _, f := range []func() error{
		r.watchPinotConfig,
		r.watchTenant,
		r.watchNamespace,
	} {
		err = f()
		if err != nil {
			return err
		}
	}

	if !watchCreatedResourcesEvents {
		return nil
	}

	createdResourceTypes := []runtime.Object{
		&corev1.ServiceAccount{TypeMeta: metav1.TypeMeta{Kind: "ServiceAccount", APIVersion: corev1.SchemeGroupVersion.String()}},
		&rbacv1.Role{TypeMeta: metav1.TypeMeta{Kind: "ClusterRole", APIVersion: rbacv1.SchemeGroupVersion.String()}},
		&rbacv1.RoleBinding{TypeMeta: metav1.TypeMeta{Kind: "ClusterRoleBinding", APIVersion: rbacv1.SchemeGroupVersion.String()}},
		&rbacv1.ClusterRole{TypeMeta: metav1.TypeMeta{Kind: "ClusterRole", APIVersion: rbacv1.SchemeGroupVersion.String()}},
		&rbacv1.ClusterRoleBinding{TypeMeta: metav1.TypeMeta{Kind: "ClusterRoleBinding", APIVersion: rbacv1.SchemeGroupVersion.String()}},
		&corev1.ConfigMap{TypeMeta: metav1.TypeMeta{Kind: "ConfigMap", APIVersion: corev1.SchemeGroupVersion.String()}},
		&corev1.Service{TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: corev1.SchemeGroupVersion.String()}},
		&appsv1.StatefulSet{TypeMeta: metav1.TypeMeta{Kind: "Statefulset", APIVersion: appsv1.SchemeGroupVersion.String()}},
		&autoscalingv2beta2.HorizontalPodAutoscaler{TypeMeta: metav1.TypeMeta{Kind: "HorizontalPodAutoscaler", APIVersion: autoscalingv2beta2.SchemeGroupVersion.String()}},
	}

	// Watch for changes to resources managed by the operator
	for _, resource := range createdResourceTypes {
		err = r.watchResource(resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcilePinot) watchNamespace() error {
	return r.Ctrl.Watch(
		&source.Kind{
			Type: &corev1.Namespace{TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: corev1.SchemeGroupVersion.String()}},
		},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(object handler.MapObject) []reconcile.Request {
				if revision, ok := object.Meta.GetLabels()[pinotv1alpha1.RevisionedAutoInjectionLabelKey]; ok {
					nn := pinotv1alpha1.NamespacedNameFromRevision(revision)
					if nn.Name == "" {
						return nil
					}
					return []reconcile.Request{
						{
							NamespacedName: nn,
						},
					}
				}

				return nil
			}),
		},
	)
}

func (r *ReconcilePinot) watchPinotConfig() error {
	return r.Ctrl.Watch(
		&source.Kind{
			Type: &pinotv1alpha1.Pinot{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pinot",
					APIVersion: pinotv1alpha1.GroupVersion.String(),
				},
			},
		},
		&handler.EnqueueRequestForObject{},
		k8sutil.GetWatchPredicateForPinot(),
	)
}

func (r *ReconcilePinot) watchTenant() error {
	return r.Ctrl.Watch(
		&source.Kind{
			Type: &pinotv1alpha1.Tenant{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Tenant",
					APIVersion: pinotv1alpha1.GroupVersion.String(),
				},
			},
		},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(object handler.MapObject) []reconcile.Request {
				if t, ok := object.Object.(*pinotv1alpha1.Tenant); ok && t.Spec.PinotServer != nil {
					return []reconcile.Request{
						{
							NamespacedName: types.NamespacedName(*t.Spec.PinotServer),
						},
					}
				}
				return nil
			}),
		},
		k8sutil.GetWatchPredicateForTenant(),
	)
}

func (r *ReconcilePinot) watchResource(resource runtime.Object) error {
	return r.Ctrl.Watch(
		&source.Kind{
			Type: resource,
		},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &pinotv1alpha1.Pinot{},
		},
		k8sutil.GetWatchPredicateForOwnedResources(&pinotv1alpha1.Pinot{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pinot",
				APIVersion: pinotv1alpha1.GroupVersion.String(),
			},
		}, true, r.Manager.GetScheme(), log),
	)
}

func (r *ReconcilePinot) watchCRDs(nn types.NamespacedName) error {
	return r.Ctrl.Watch(
		&source.Kind{
			Type: &apiextensionsv1.CustomResourceDefinition{
				TypeMeta: metav1.TypeMeta{
					Kind:       "CustomResourceDefinition",
					APIVersion: apiextensionsv1.SchemeGroupVersion.String(),
				},
			},
		},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(object handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{
						NamespacedName: nn,
					},
				}
			}),
		},
		crds.GetWatchPredicateForCRDs(),
	)
}
