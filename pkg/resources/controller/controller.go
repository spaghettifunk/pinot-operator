package controller

import (
	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil"
	"github.com/spaghettifunk/pinot-operator/pkg/resources"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	componentName   = "controller"
	statefulsetName = "pinot-controller"
	configmapName   = "pinot-controller"
	serviceName     = "pinot-controller-api"
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
		{Resource: r.configmap, DesiredState: desiredState},
		{Resource: r.statefulsets, DesiredState: desiredState},
		{Resource: r.service, DesiredState: desiredState},
		{Resource: r.serviceHeadless, DesiredState: desiredState},
	} {
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
		"pinot.operator/app":             r.Config.ClusterName,
		"pinot.operator/component":       componentName,
		"pinot.operator/release-version": pinotv1alpha1.OperatorVersion,
	}
}

// deploymentLabels returns the labels used for the deployment of the web component
func (r *Reconciler) deploymentLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name": componentName,
	}
}

func (r *Reconciler) annotations() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name": componentName,
	}
}
