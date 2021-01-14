package broker

import (
	"strconv"

	"github.com/hoisie/mustache"
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var controllerConfig = `
	pinot.broker.client.queryPort={{queryPort}}
	pinot.broker.routing.table.builder.class=random
	pinot.set.instance.id.to.hostname=true
`

func (r *Reconciler) configmap() runtime.Object {
	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMeta(configmapName, r.labels(), r.Config),
		Data: map[string]string{
			"pinot-broker.conf": mustache.Render(controllerConfig, map[string]string{
				"queryPort": strconv.Itoa(r.Config.Spec.Broker.Service.Port),
			}),
		},
	}
}
