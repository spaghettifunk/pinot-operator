package v1alpha1

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	pinotImageHub     = "apachepinot/pinot"
	pinotImageVersion = "latest"
	// replicas
	defaultReplicaCount = 1
	defaultMinReplicas  = 1
	defaultMaxReplicas  = 5
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

// SetDefaults sets some of the specs of the various components
func (r *Pinot) SetDefaults() {
	// controller
	if r.Spec.Controller == nil {
		r.Spec.Controller = &ControllerConfiguration{}
	}
	if r.Spec.Controller.Resources == nil {
		r.Spec.Controller.Resources = defaultResources
	}
	if r.Spec.Controller.Service.Port == 0 {
		r.Spec.Controller.Service.Port = 9000
	}
	if r.Spec.Controller.JvmOptions == "" {
		r.Spec.Controller.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	// broker
	if r.Spec.Broker == nil {
		r.Spec.Broker = &BrokerConfiguration{}
	}
	if r.Spec.Broker.Resources == nil {
		r.Spec.Broker.Resources = defaultResources
	}
	if r.Spec.Broker.JvmOptions == "" {
		r.Spec.Broker.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	if r.Spec.Broker.Service.Port == 0 {
		r.Spec.Broker.Service.Port = 8099
	}
	if r.Spec.Broker.ExternalService.Port == 0 {
		r.Spec.Broker.ExternalService.Port = 8099
	}
	// server
	if r.Spec.Server == nil {
		r.Spec.Server = &ServerConfiguration{}
	}
	if r.Spec.Server.Service.Port == 0 {
		r.Spec.Server.Service.Port = 8098
	}
	if r.Spec.Server.AdminPort == 0 {
		r.Spec.Server.AdminPort = 8097
	}
	if r.Spec.Server.Resources == nil {
		r.Spec.Server.Resources = defaultResources
	}
	// zookeeper
	if r.Spec.Zookeeper == nil {
		r.Spec.Zookeeper = &ZookeeperConfiguration{
			Storage: &zookeeperStorage{},
		}
	}
	if r.Spec.Zookeeper.Replicas == 0 {
		r.Spec.Zookeeper.Replicas = 1
	}
	if r.Spec.Zookeeper.Resources == nil {
		r.Spec.Zookeeper.Resources = zookeeperDefaultResources
	}
	if r.Spec.Zookeeper.JvmOptions == "" {
		r.Spec.Zookeeper.JvmOptions = "-Xmx2G -Xms2G"
	}
}
