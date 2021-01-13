package zookeeper

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	apiv1 "k8s.io/api/core/v1"
)

func (r *Reconciler) serviceHeadless() runtime.Object {
	return &apiv1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   serviceHeadlessName,
			Labels: r.labels(),
		},
		Spec: apiv1.ServiceSpec{
			Type:      apiv1.ServiceTypeClusterIP,
			ClusterIP: "None",
			Selector:  r.selector(componentName),
			Ports: []apiv1.ServicePort{
				{
					Name:       "client",
					Port:       int32(zookeeperClientPort),
					TargetPort: intstr.FromString("client"),
					Protocol:   apiv1.ProtocolTCP,
				},
				{
					Name:       "election",
					Port:       int32(zookeeperElectionPort),
					TargetPort: intstr.FromString("election"),
					Protocol:   apiv1.ProtocolTCP,
				},
				{
					Name:       "server",
					Port:       int32(zookeeperServerPort),
					TargetPort: intstr.FromString("server"),
					Protocol:   apiv1.ProtocolTCP,
				},
			},
		},
	}
}
