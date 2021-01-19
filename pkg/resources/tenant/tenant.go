package tenant

import (
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"github.com/spaghettifunk/pinot-operator/pkg/resources"
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"
)

const (
	ResourceName  = "pinot-tenant"
	componentName = "tenant"
)

var (
	resourceLabels = map[string]string{
		"app":   "pinot-tenant",
		"pinot": "tenant",
	}
)

type Reconciler struct {
	resources.Reconciler
}

// New .
func New(client client.Client, config *pinotv1alpha1.Pinot) *Reconciler {
	return &Reconciler{
		Reconciler: resources.Reconciler{
			Client: client,
			Config: config,
		},
	}
}

func (r *Reconciler) Reconcile(log logr.Logger) error {
	log = log.WithValues("component", componentName)

	desiredState := k8sutil.DesiredStatePresent

	log.Info("Reconciling")

	spec := pinotv1alpha1.TenantSpec{}
	spec.PinotServer = &pinotv1alpha1.NamespacedName{
		Name:      r.Config.Name,
		Namespace: r.Config.Namespace,
	}
	spec.Labels = r.labels()
	objectMeta := templates.ObjectMetaWithRevision(ResourceName, spec.Labels, r.Config)

	object := &pinotv1alpha1.Tenant{
		ObjectMeta: objectMeta,
		Spec:       spec,
	}
	pinotv1alpha1.SetTenantDefaults(object)

	err := k8sutil.Reconcile(log, r.Client, object, desiredState)
	if err != nil {
		return emperror.WrapWith(err, "failed to reconcile resource", "resource", object.GetObjectKind().GroupVersionKind())
	}

	// selector := resourceLabels

	// var drs = []resources.DynamicResourceWithDesiredState{
	// 	{DynamicResource: func() *k8sutil.DynamicObject { return r.meshExpansionGateway(selector) }, DesiredState: meshExpansionDesiredState},
	// }
	// for _, dr := range drs {
	// 	o := dr.DynamicResource()
	// 	err := o.Reconcile(log, r.dynamic, dr.DesiredState)
	// 	if err != nil {
	// 		return emperror.WrapWith(err, "failed to reconcile dynamic resource", "resource", o.Gvr)
	// 	}
	// }

	log.Info("Reconciled")

	return nil
}

func (r *Reconciler) labels() map[string]string {
	return nil
}
