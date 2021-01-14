package server

import (
	"strconv"

	"github.com/hoisie/mustache"
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var controllerConfig = `
	pinot.server.netty.port={{servicePort}}
	pinot.server.adminapi.port={{adminPort}}
	pinot.server.instance.dataDir={{dataDir}}
	pinot.server.instance.segmentTarDir={{segmentDir}}
	pinot.set.instance.id.to.hostname=true
	pinot.server.instance.realtime.alloc.offheap=true
`

func (r *Reconciler) configmap() runtime.Object {
	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMeta(configmapName, r.labels(), r.Config),
		Data: map[string]string{
			"pinot-server.conf": mustache.Render(controllerConfig, map[string]string{
				"servicePort": strconv.Itoa(r.Config.Spec.Server.Service.Port),
				"adminPort":   strconv.Itoa(r.Config.Spec.Server.AdminPort),
				"dataDir":     "/var/pinot/server/data/index",
				"segmentDir":  "/var/pinot/server/data/segment",
			}),
		},
	}
}
