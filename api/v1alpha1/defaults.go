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
	// zookeeper image
	zookeeperImageHub     = "zookeeper"
	zookeeperImageVersion = "3.5.5"
	// replicas
	defaultReplicaCount = 1
	defaultMinReplicas  = 1
	defaultMaxReplicas  = 5
	// resources
)

var defaultResources = &apiv1.ResourceRequirements{
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("512m"),
		apiv1.ResourceMemory: resource.MustParse("2Gi"),
	},
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("256m"),
		apiv1.ResourceMemory: resource.MustParse("1Gi"),
	},
}

var zookeeperDefaultResources = &apiv1.ResourceRequirements{
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("512m"),
		apiv1.ResourceMemory: resource.MustParse("2Gi"),
	},
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("256m"),
		apiv1.ResourceMemory: resource.MustParse("1Gi"),
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
	if config.Spec.Controller == nil {
		config.Spec.Controller = &ControllerConfiguration{}
	}
	if config.Spec.Controller.Resources == nil {
		config.Spec.Controller.Resources = defaultResources
	}
	if config.Spec.Controller.Service.Port == 0 {
		config.Spec.Controller.Service.Port = 9000
	}
	if config.Spec.Controller.DiskSize == "" {
		config.Spec.Controller.DiskSize = "1G"
	}
	if config.Spec.Controller.JvmOptions == "" {
		config.Spec.Controller.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	// broker
	if config.Spec.Broker == nil {
		config.Spec.Broker = &BrokerConfiguration{}
	}
	if config.Spec.Broker.Resources == nil {
		config.Spec.Broker.Resources = defaultResources
	}
	if config.Spec.Broker.JvmOptions == "" {
		config.Spec.Broker.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	if config.Spec.Broker.Service.Port == 0 {
		config.Spec.Broker.Service.Port = 8099
	}
	if config.Spec.Broker.ExternalService.Port == 0 {
		config.Spec.Broker.ExternalService.Port = 8099
	}
	// server
	if config.Spec.Server == nil {
		config.Spec.Server = &ServerConfiguration{}
	}
	if config.Spec.Server.Service.Port == 0 {
		config.Spec.Server.Service.Port = 8098
	}
	if config.Spec.Server.AdminPort == 0 {
		config.Spec.Server.AdminPort = 8097
	}
	if config.Spec.Server.Resources == nil {
		config.Spec.Server.Resources = defaultResources
	}
	if config.Spec.Server.DiskSize == "" {
		config.Spec.Server.DiskSize = "4G"
	}
	// zookeeper
	if config.Spec.Zookeeper == nil {
		config.Spec.Zookeeper = &ZookeeperConfiguration{
			Storage: &zookeeperStorage{},
		}
	}
	if config.Spec.Zookeeper.Image == nil {
		config.Spec.Zookeeper.Image = util.StrPointer(fmt.Sprintf("%s:%s", zookeeperImageHub, zookeeperImageVersion))
	}
	if config.Spec.Zookeeper.Replicas == 0 {
		config.Spec.Zookeeper.Replicas = 1
	}
	if config.Spec.Zookeeper.Resources == nil {
		config.Spec.Zookeeper.Resources = zookeeperDefaultResources
	}
	if config.Spec.Zookeeper.JvmOptions == "" {
		config.Spec.Zookeeper.JvmOptions = "-Xmx2G -Xms2G"
	}
	if config.Spec.Zookeeper.Storage.Size == "" {
		config.Spec.Zookeeper.Storage.Size = "5Gi"
	}

	// TODO: DeepStorage
	// if config.Spec.DeepStorage.Image == nil {
	// 	config.Spec.DeepStorage.Image = util.StrPointer(defaultControllerImage)
	// }
	// if config.Spec.DeepStorage.Resources == nil {
	// 	config.Spec.DeepStorage.Resources = defaultResources
	// }
}
