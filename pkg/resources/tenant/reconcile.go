package tenant

import (
	"context"

	pinotsdk "github.com/spaghettifunk/pinot-go-client/client"
	"github.com/spaghettifunk/pinot-go-client/client/tenant"
	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(pinotClient *pinotsdk.PinotSdk, client runtimeClient.Client, desired *pinotv1alpha1.Tenant, desiredState k8sutil.DesiredState) error {
	// contact Apache Pinot APIs to create/update/delete the tenant
	pinotClient.Tenant.GetAllTenants(&tenant.GetAllTenantsParams{
		Type:    util.StrPointer(""),
		Context: context.TODO(),
	})

	return nil
}
