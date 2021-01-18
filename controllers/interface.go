package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PinotReconciler interface {
	reconcile.Reconciler
	initWatches(watchCreatedResourcesEvents bool) error
	setController(ctrl controller.Controller)
}
