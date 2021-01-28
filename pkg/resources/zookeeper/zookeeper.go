package zookeeper

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/pkg/apis/pinot/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"github.com/spaghettifunk/pinot-operator/pkg/resources"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName                     = "zookeeper"
	statefulsetName                   = "pinot-zookeeper"
	configmapName                     = "pinot-zookeeper-config"
	serviceName                       = "pinot-zookeeper"
	serviceHeadlessName               = "pinot-zookeeper-headless"
	podDisruptionBudgetName           = "pinot-zookeeper"
	podDisruptionBudgetMaxUnavailable = 1
	zookeeperDataVolumeName           = "pinot-zookeeper-data"
	zookeeperConfigVolumeName         = "pinot-zookeeper-config"
	// ports
	zookeeperClientPort   = 2181
	zookeeperServerPort   = 2888
	zookeeperElectionPort = 3888
)

// Reconciler .
type Reconciler struct {
	resources.Reconciler
}

// New .
func New(client client.Client, config *pinotv1alpha1.Pinot) *Reconciler {
	return &Reconciler{
		Reconciler: resources.Reconciler{
			Client: client,
			Config: config,
		},
	}
}

// Reconcile .
func (r *Reconciler) Reconcile(log logr.Logger) error {
	log = log.WithValues("component", componentName)

	desiredState := k8sutil.DesiredStatePresent

	log.Info("Reconciling")

	for _, res := range []resources.ResourceWithDesiredState{
		{Name: configmapName, Resource: r.configmap, DesiredState: desiredState},
		{Name: serviceName, Resource: r.service, DesiredState: desiredState},
		{Name: serviceHeadlessName, Resource: r.serviceHeadless, DesiredState: desiredState},
		{Name: podDisruptionBudgetName, Resource: r.podDisruptionBuget, DesiredState: desiredState},
		{Name: statefulsetName, Resource: r.statefulsets, DesiredState: desiredState},
	} {
		// reconcile resource
		o := res.Resource()
		err := k8sutil.Reconcile(log, r.Client, o, res.DesiredState)
		if err != nil {
			return emperror.WrapWith(err, "failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
	}

	log.Info("Reconciled")

	return nil
}

func (r *Reconciler) labels() map[string]string {
	return map[string]string{
		"app":             r.Config.Spec.ClusterName,
		"component":       componentName,
		"release":         "pinot",
		"release-version": pinotv1alpha1.OperatorVersion,
	}
}

func (r *Reconciler) selector(name string) map[string]string {
	return map[string]string{
		"app":       "pinot",
		"release":   "pinot",
		"component": name,
	}
}

func (r *Reconciler) WaitForCreation() error {
	var pods v1.PodList
	ls, err := labels.Parse(fmt.Sprintf("component=%s", componentName))
	if err != nil {
		return err
	}

	err = r.Client.List(context.Background(), &pods, client.InNamespace(r.Config.Namespace), client.MatchingLabelsSelector{
		Selector: ls,
	})
	if err != nil {
		return emperror.Wrap(err, "could not list pods")
	}

	readyContainers := 0
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.Ready {
					readyContainers++
				}
			}
			// wait for all the containers to be ready before continuing
			if readyContainers == (len(pod.Status.ContainerStatuses) * len(pods.Items)) {
				return nil
			}
		}
	}
	return errors.Errorf("Zookeeper is not ready yet")
}
