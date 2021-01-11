package v1alpha1

import (
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"

	apiv1 "k8s.io/api/core/v1"
)

const (
	pinotImageHub          = "docker.io/davideberdin"
	pinotImageVersion      = ""
	defaultImageHub        = "gcr.io/apache"
	defaultImageVersion    = ""
	defaultImagePullPolicy = "IfNotPresent"
	defaultNetworkName     = "cluster.local"
	// replicas
	defaultReplicaCount = 1
	defaultMinReplicas  = 1
	defaultMaxReplicas  = 5
	// images
	defaultControllerImage = defaultImageHub + "/" + "controller" + ":" + defaultImageVersion
	defaultBrokerImage     = defaultImageHub + "/" + "broker" + ":" + defaultImageVersion
	defaultServerImage     = defaultImageHub + "/" + "server" + ":" + defaultImageVersion
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
	// controller
	if config.Spec.Image == nil {
		config.Spec.Image = util.StrPointer(defaultControllerImage)
	}
	if config.Spec.Controller.Resources == nil {
		config.Spec.Controller.Resources = defaultResources
	}
	// broker
	if config.Spec.Broker.Resources == nil {
		config.Spec.Broker.Resources = defaultResources
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
