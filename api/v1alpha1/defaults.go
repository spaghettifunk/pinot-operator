package v1alpha1

import "github.com/spaghettifunk/pinot-operator/pkg/util"

// SetDefaults sets some of the specs of the various components
func (r *Pinot) SetDefaults() {
	// controller
	if r.Spec.Controller.Service.Port == 0 {
		r.Spec.Controller.Service.Port = 9000
	}
	if r.Spec.Controller.DiskSize == "" {
		r.Spec.Controller.DiskSize = "1Gi"
	}
	// broker
	if r.Spec.Broker.Service.Port == 0 {
		r.Spec.Broker.Service.Port = 8099
	}
	if r.Spec.Broker.ExternalService.Port == 0 {
		r.Spec.Broker.ExternalService.Port = 8099
	}
	// server
	if r.Spec.Server.Service.Port == 0 {
		r.Spec.Server.Service.Port = 8098
	}
	if r.Spec.Server.AdminPort == 0 {
		r.Spec.Server.AdminPort = 8097
	}
	if r.Spec.Server.DiskSize == "" {
		r.Spec.Server.DiskSize = "4Gi"
	}
	// zookeeper
	if r.Spec.Zookeeper.Image == nil {
		r.Spec.Zookeeper.Image = util.StrPointer("zookeeper:3.5.5")
	}
	if r.Spec.Zookeeper.Storage.Size == "" {
		r.Spec.Zookeeper.Storage.Size = "5Gi"
	}
}
