package broker

import (
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) serviceExternal() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMeta(serviceExternalName, r.labels(), r.Config),
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeLoadBalancer,
			Selector: r.labels(),
			Ports: []apiv1.ServicePort{
				templates.DefaultServicePort(
					"external-broker",
					r.Config.Spec.Broker.ExternalService.Port,
					r.Config.Spec.Broker.ExternalService.Port,
				),
			},
		},
	}
}
