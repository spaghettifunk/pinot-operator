package zookeeper

import (
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	"k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) podDisruptionBuget() runtime.Object {
	return &v1beta1.PodDisruptionBudget{
		ObjectMeta: v1.ObjectMeta{
			Name:   podDisruptionBudgetName,
			Labels: r.labels(),
		},
		Spec: v1beta1.PodDisruptionBudgetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: r.selector(componentName),
			},
			MaxUnavailable: util.IntstrPointer(podDisruptionBudgetMaxUnavailable),
		},
	}
}
