package controller

import (
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) serviceExternal() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMetaWithAnnotations(serviceExternalName, r.labels(), r.Config.Spec.Controller.ExternalService.Annotations, r.Config),
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeLoadBalancer,
			Selector: r.selector(componentName),
			Ports: []apiv1.ServicePort{
				templates.DefaultServicePort(
					"external-controller",
					r.Config.Spec.Controller.ExternalService.Port,
					r.Config.Spec.Controller.ExternalService.Port,
				),
			},
		},
	}
}
