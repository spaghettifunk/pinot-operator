package k8sutil

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/spaghettifunk/pinot-operator/pkg/k8sutil/wait"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	detachedPodLabel = "operators.apache.io/detached-pod"
)

// DesiredState .
type DesiredState interface {
	AfterRecreate(current, desired runtime.Object) error
	BeforeRecreate(current, desired runtime.Object) error
	ShouldRecreate(current, desired runtime.Object) (bool, error)
	AfterCreate(desired runtime.Object) error
	BeforeCreate(desired runtime.Object) error
	ShouldCreate(desired runtime.Object) (bool, error)
	AfterUpdate(current, desired runtime.Object, inSync bool) error
	BeforeUpdate(current, desired runtime.Object) error
	ShouldUpdate(current, desired runtime.Object) (bool, error)
	AfterDelete(current runtime.Object) error
	BeforeDelete(current runtime.Object) error
	ShouldDelete(current runtime.Object) (bool, error)
}

// StaticDesiredState .
type StaticDesiredState string

// AfterRecreate .
func (s StaticDesiredState) AfterRecreate(_, _ runtime.Object) error {
	return nil
}

// BeforeRecreate .
func (s StaticDesiredState) BeforeRecreate(_, _ runtime.Object) error {
	return nil
}

// ShouldRecreate .
func (s StaticDesiredState) ShouldRecreate(_, _ runtime.Object) (bool, error) {
	return true, nil
}

// AfterCreate .
func (s StaticDesiredState) AfterCreate(_ runtime.Object) error {
	return nil
}

// BeforeCreate .
func (s StaticDesiredState) BeforeCreate(_ runtime.Object) error {
	return nil
}

// ShouldCreate .
func (s StaticDesiredState) ShouldCreate(_ runtime.Object) (bool, error) {
	return true, nil
}

// AfterUpdate .
func (s StaticDesiredState) AfterUpdate(_, _ runtime.Object, _ bool) error {
	return nil
}

// BeforeUpdate .
func (s StaticDesiredState) BeforeUpdate(_, _ runtime.Object) error {
	return nil
}

// ShouldUpdate .
func (s StaticDesiredState) ShouldUpdate(_, _ runtime.Object) (bool, error) {
	return true, nil
}

// AfterDelete .
func (s StaticDesiredState) AfterDelete(_ runtime.Object) error {
	return nil
}

// BeforeDelete .
func (s StaticDesiredState) BeforeDelete(_ runtime.Object) error {
	return nil
}

// ShouldDelete .
func (s StaticDesiredState) ShouldDelete(_ runtime.Object) (bool, error) {
	return true, nil
}

const (
	// DesiredStatePresent .
	DesiredStatePresent StaticDesiredState = "present"
	// DesiredStateAbsent .
	DesiredStateAbsent StaticDesiredState = "absent"
	// DesiredStateExists .
	DesiredStateExists StaticDesiredState = "exists"
)

// RecreateAwareStatefulsetDesiredState .
type RecreateAwareStatefulsetDesiredState struct {
	StaticDesiredState

	client    client.Client
	scheme    *runtime.Scheme
	log       logr.Logger
	podLabels map[string]string
}

// NewRecreateAwareDeploymentDesiredState .
func NewRecreateAwareDeploymentDesiredState(c client.Client, scheme *runtime.Scheme, log logr.Logger, podLabels map[string]string) RecreateAwareStatefulsetDesiredState {
	podLabels = util.MergeMultipleStringMaps(map[string]string{
		detachedPodLabel: "true",
	}, podLabels)

	return RecreateAwareStatefulsetDesiredState{
		client:    c,
		scheme:    scheme,
		log:       log,
		podLabels: podLabels,
	}
}

// AfterRecreate .
func (r RecreateAwareStatefulsetDesiredState) AfterRecreate(current, desired runtime.Object) error {
	var statefulset *appsv1.StatefulSet
	var ok bool
	if statefulset, ok = desired.(*appsv1.StatefulSet); !ok {
		return nil
	}

	return r.waitForStatefulsetAndRemoveDetached(statefulset)
}

// AfterUpdate .
func (r RecreateAwareStatefulsetDesiredState) AfterUpdate(current, desired runtime.Object, inSync bool) error {
	var statefulset *appsv1.StatefulSet
	var ok bool
	if statefulset, ok = desired.(*appsv1.StatefulSet); !ok {
		return nil
	}

	return r.waitForStatefulsetAndRemoveDetached(statefulset)
}

// BeforeRecreate .
func (r RecreateAwareStatefulsetDesiredState) BeforeRecreate(current, desired runtime.Object) error {
	var statefulset *appsv1.StatefulSet
	var ok bool
	if statefulset, ok = current.(*appsv1.StatefulSet); !ok {
		return nil
	}

	err := DetachPodsFromStatefulset(r.client, statefulset, r.log, r.podLabels)
	if err != nil {
		return err
	}

	return nil
}

func (r RecreateAwareStatefulsetDesiredState) waitForStatefulsetAndRemoveDetached(statefulset *appsv1.StatefulSet) error {
	rcc := wait.NewResourceConditionChecks(r.client, wait.Backoff{
		Duration: time.Second * 5,
		Factor:   1,
		Jitter:   0,
		Steps:    12,
	}, r.log.WithName("wait"), r.scheme)

	err := rcc.WaitForResources("readiness", []runtime.Object{statefulset}, wait.ExistsConditionCheck, wait.ReadyReplicasConditionCheck)
	if err != nil {
		return err
	}

	pods := &corev1.PodList{}
	err = r.client.List(context.Background(), pods, client.InNamespace(statefulset.GetNamespace()), client.MatchingLabelsSelector{
		Selector: labels.Set(r.podLabels).AsSelector(),
	})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		r.log.Info("removing detached pods")
		err = r.client.Delete(context.Background(), &pod)
		if err != nil {
			return err
		}
	}

	return nil
}
