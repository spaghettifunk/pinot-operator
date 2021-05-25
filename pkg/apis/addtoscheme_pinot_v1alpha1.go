package api

import (
	"github.com/spaghettifunk/pinot-operator/pkg/apis/pinot/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
