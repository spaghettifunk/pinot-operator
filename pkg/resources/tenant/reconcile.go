package tenant

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	pinotsdk "github.com/spaghettifunk/pinot-go-client/client"
	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/pinot/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(pinotClient *pinotsdk.PinotSdk, client runtimeClient.Client, desired *pinotv1alpha1.Tenant, desiredState k8sutil.DesiredState) error {
	// contact Apache Pinot APIs to create/update/delete the tenant
	// tenants, err := pinotClient.Tenant.GetAllTenants(&tenant.GetAllTenantsParams{
	// 	Type:    util.StrPointer(string(desired.Spec.Role)),
	// 	Context: context.TODO(),
	// })
	resp, err := http.Get("http://pinot-controller.pinot-system:9000/tenants")
	if err != nil {
		return err
	}

	spew.Dump(resp)

	return nil
}
