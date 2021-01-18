package v1alpha1

import "github.com/spaghettifunk/pinot-operator/pkg/util"

// SetDefaults sets some of the specs of the various components
func SetDefaults(config *Pinot) {
	// controller
	if config.Spec.Controller.Service.Port == 0 {
		config.Spec.Controller.Service.Port = 9000
	}
	if config.Spec.Controller.DiskSize == "" {
		config.Spec.Controller.DiskSize = "1Gi"
	}
	// broker
	if config.Spec.Broker.Service.Port == 0 {
		config.Spec.Broker.Service.Port = 8099
	}
	if config.Spec.Broker.ExternalService.Port == 0 {
		config.Spec.Broker.ExternalService.Port = 8099
	}
	// server
	if config.Spec.Server.Service.Port == 0 {
		config.Spec.Server.Service.Port = 8098
	}
	if config.Spec.Server.AdminPort == 0 {
		config.Spec.Server.AdminPort = 8097
	}
	if config.Spec.Server.DiskSize == "" {
		config.Spec.Server.DiskSize = "4Gi"
	}
	// zookeeper
	if config.Spec.Zookeeper.Image == nil {
		config.Spec.Zookeeper.Image = util.StrPointer("zookeeper:3.5.5")
	}
	if config.Spec.Zookeeper.Storage.Size == "" {
		config.Spec.Zookeeper.Storage.Size = "5Gi"
	}
}
