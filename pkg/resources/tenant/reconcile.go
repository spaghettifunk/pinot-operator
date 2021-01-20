package tenant

import (
	pinotsdk "github.com/spaghettifunk/pinot-go-client/client"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(pinotClient *pinotsdk.PinotSdk, client runtimeClient.Client, desired runtime.Object, desiredState k8sutil.DesiredState) error {
	// contact Apache Pinot APIs to create/update/delete the tenant
	return nil
}
