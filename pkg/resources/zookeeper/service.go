package zookeeper

import (
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: templates.ObjectMetaWithAnnotations(serviceName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeClusterIP,
			Selector: r.selector(componentName),
			Ports: []apiv1.ServicePort{
				{
					Name:       "client",
					Port:       int32(zookeeperClientPort),
					TargetPort: intstr.FromString("client"),
					Protocol:   apiv1.ProtocolTCP,
				},
			},
		},
	}
}
