package controller

import (
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/runtime"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) serviceHeadless() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMeta(serviceHeadlessName, r.labels(), r.Config),
		Spec: apiv1.ServiceSpec{
			Type:      apiv1.ServiceTypeClusterIP,
			ClusterIP: "None",
			Selector:  r.labels(),
			Ports: []apiv1.ServicePort{
				{
					Port: int32(r.Config.Spec.Controller.Service.Port),
				},
			},
		},
	}
}
