package controller

import (
	"fmt"
	"strconv"

	"github.com/hoisie/mustache"
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var controllerConfig = `
	controller.helix.cluster.name={{clusterName}}
	controller.port={{controllerPort}}
	{{#controllerVIPHost}}
	controller.vip.host={{controllerVIPHost}}
	{{/controllerVIPEnabled}}
	{{#controllerVIPPort}}
	controller.vip.port={{controllerVIPPort}}
	{{#controllerVIPort}}
	controller.data.dir={{controllerDataDir}}
	controller.zk.str={{zookeeperURL}}
	pinot.set.instance.id.to.hostname=true
`

func (r *Reconciler) configmap() runtime.Object {
	// TODO: fix name of zookeeper server
	zookeeperURL := fmt.Sprintf("%s:%s", "pinot-zookeeper", "2181")

	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMetaWithAnnotations(configmapName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		Data: map[string]string{
			"pinot-controller.conf": mustache.Render(controllerConfig, map[string]string{
				"clusterName":       r.Config.Spec.ClusterName,
				"controllerPort":    strconv.Itoa(r.Config.Spec.Controller.Service.Port),
				"controllerVIPHost": r.Config.Spec.Controller.VIPHost,
				"controllerVIPPort": r.Config.Spec.Controller.VIPPort,
				"controllerDataDir": "/var/pinot/controller/data",
				"zookeeperURL":      zookeeperURL,
			}),
		},
	}
}
