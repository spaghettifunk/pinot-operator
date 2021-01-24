package tenant

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-openapi/strfmt"
	"github.com/goph/emperror"
	pinotsdk "github.com/spaghettifunk/pinot-go-client/client"
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
	PinotClient *pinotsdk.PinotSdk
}

// New .
func New(client client.Client, config *pinotv1alpha1.Pinot) *Reconciler {
	pc := pinotsdk.NewHTTPClientWithConfig(strfmt.Default, &pinotsdk.TransportConfig{
		Host:     fmt.Sprintf("%s:%d", "localhost", config.Spec.Controller.Service.Port),
		BasePath: pinotsdk.DefaultBasePath,
		Schemes:  []string{"http"},
	})
	return &Reconciler{
		Reconciler: resources.Reconciler{
			Client: client,
			Config: config,
		},
		PinotClient: pc,
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

	if err := Reconcile(r.PinotClient, r.Client, object, desiredState); err != nil {
		return emperror.WrapWith(err, "failed to reconcile resource", "resource", object.GetObjectKind().GroupVersionKind())
	}

	log.Info("Reconciled")

	return nil
}

func (r *Reconciler) labels() map[string]string {
	return nil
}
