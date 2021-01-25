package sdk

import (
	"fmt"

	operatorsv1alpha1 "github.com/spaghettifunk/pinot-operator/api/pinot/v1alpha1"
)

const (
	pinotControllerHeadless = "pinot-controller-headless"
	kubernetesDomain        = "cluster.local"
)

// GeneratePinotControlleAddressWithoutPort returns the host of the controller service
func GeneratePinotControlleAddressWithoutPort(config *operatorsv1alpha1.Pinot) string {
	// return "localhost"
	return fmt.Sprintf("%s.%s.svc.%s",
		pinotControllerHeadless,
		config.Namespace,
		kubernetesDomain,
	)
}

// GeneratePinotControllerAddress returns the fully qualified DNS of the controller server
func GeneratePinotControllerAddress(config *operatorsv1alpha1.Pinot) string {
	return fmt.Sprintf("%s:%d", GeneratePinotControlleAddressWithoutPort(config), config.Spec.Controller.Service.Port)
}
