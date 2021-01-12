package v1alpha1

import (
	"fmt"

	"github.com/spaghettifunk/pinot-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"

	apiv1 "k8s.io/api/core/v1"
)

const (
	pinotImageHub     = "apachepinot/pinot"
	pinotImageVersion = "latest"
	// replicas
	defaultReplicaCount = 1
	defaultMinReplicas  = 1
	defaultMaxReplicas  = 5
	// resources
)

var defaultResources = &apiv1.ResourceRequirements{
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("100m"),
		apiv1.ResourceMemory: resource.MustParse("50Mi"),
	},
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("100m"),
		apiv1.ResourceMemory: resource.MustParse("50Mi"),
	},
}

// var defaultControllerServicePorts = []ServicePort{
// 	{ServicePort: corev1.ServicePort{Name: "http", Port: int32(8085), TargetPort: intstr.FromString("8085")}},
// }

// SetDefaults sets the defaults values for all the components
func SetDefaults(config *Pinot) {
	// common
	if config.Spec.Image == nil {
		config.Spec.Image = util.StrPointer(fmt.Sprintf("%s:%s", pinotImageHub, pinotImageVersion))
	}
	// controller
	if config.Spec.Controller.Resources == nil {
		config.Spec.Controller.Resources = defaultResources
	}
	if config.Spec.Controller.DiskSize == "" {
		config.Spec.Controller.DiskSize = "1G"
	}
	if config.Spec.Controller.JvmOptions == "" {
		config.Spec.Controller.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	// broker
	if config.Spec.Broker.Resources == nil {
		config.Spec.Broker.Resources = defaultResources
	}
	if config.Spec.Broker.JvmOptions == "" {
		config.Spec.Broker.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	// server
	if config.Spec.Server.Resources == nil {
		config.Spec.Server.Resources = defaultResources
	}

	// TODO: Zookeeper
	// if config.Spec.Zookeeper.Image == nil {
	// 	config.Spec.Zookeeper.Image = util.StrPointer(defaultControllerImage)
	// }
	// if config.Spec.Zookeeper.Resources == nil {
	// 	config.Spec.Zookeeper.Resources = defaultResources
	// }

	// TODO: DeepStorage
	// if config.Spec.DeepStorage.Image == nil {
	// 	config.Spec.DeepStorage.Image = util.StrPointer(defaultControllerImage)
	// }
	// if config.Spec.DeepStorage.Resources == nil {
	// 	config.Spec.DeepStorage.Resources = defaultResources
	// }
}
