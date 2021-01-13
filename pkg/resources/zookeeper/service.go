package zookeeper

import (
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   serviceName,
			Labels: r.labels(),
		},
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
