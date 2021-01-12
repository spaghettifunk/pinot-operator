package resources

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"

	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceWithDesiredState defines the desidered state based on the resources
type ResourceWithDesiredState struct {
	Name         string
	Resource     Resource
	DesiredState k8sutil.DesiredState
}

// Reconciler is the object holding the client and the configuration of the operator
type Reconciler struct {
	client.Client
	Config *pinotv1alpha1.Pinot
}

// ComponentReconciler is the interface that is used for each sub-component to reconcile with the config
type ComponentReconciler interface {
	Reconcile(log logr.Logger) error
}

// Resource defines a runtime.Object type
type Resource func() runtime.Object
