package v1alpha1

// SetDefaults sets some of the specs of the various components
func (r *Pinot) SetDefaults() {
	// controller
	if r.Spec.Controller == nil {
		r.Spec.Controller = &ControllerConfiguration{}
	}
	if r.Spec.Controller.Service.Port == 0 {
		r.Spec.Controller.Service.Port = 9000
	}
	// broker
	if r.Spec.Broker == nil {
		r.Spec.Broker = &BrokerConfiguration{}
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
	// zookeeper
	if r.Spec.Zookeeper == nil {
		r.Spec.Zookeeper = &ZookeeperConfiguration{
			Storage: &zookeeperStorage{},
		}
	}
}
