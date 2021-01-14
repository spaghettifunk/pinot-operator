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
	controller.access.protocols.http.port={{controllerPort}}
	{{#isVIPHostSet}}
	controller.vip.host={{controllerVIPHost}}
	{{/isVIPHostSet}}
	{{#isVIPPortSet}}
	controller.vip.port={{controllerVIPPort}}
	{{/isVIPPortSet}}
	controller.data.dir={{controllerDataDir}}
	controller.zk.str={{zookeeperURL}}
	pinot.set.instance.id.to.hostname=true
`

func (r *Reconciler) configmap() runtime.Object {
	// TODO: fix name of zookeeper server
	zookeeperURL := fmt.Sprintf("%s:%s", "pinot-zookeeper", "2181")

	isVIPHostSet := r.Config.Spec.Controller.VIPHost != ""
	isVIPPortSet := r.Config.Spec.Controller.VIPPort != 0

	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMeta(configmapName, r.labels(), r.Config),
		Data: map[string]string{
			"pinot-controller.conf": mustache.Render(controllerConfig, map[string]interface{}{
				"clusterName":       r.Config.Spec.ClusterName,
				"controllerPort":    strconv.Itoa(r.Config.Spec.Controller.Service.Port),
				"isVIPHostSet":      isVIPHostSet,
				"isVIPPortSet":      isVIPPortSet,
				"controllerVIPHost": r.Config.Spec.Controller.VIPHost,
				"controllerVIPPort": strconv.Itoa(r.Config.Spec.Controller.VIPPort),
				"controllerDataDir": "/var/pinot/controller/data",
				"zookeeperURL":      zookeeperURL,
			}),
		},
	}
}
