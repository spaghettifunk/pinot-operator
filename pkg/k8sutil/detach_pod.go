package k8sutil

import (
	"context"

	"github.com/banzaicloud/istio-operator/pkg/util"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DetachPodsFromStatefulset .
func DetachPodsFromStatefulset(c client.Client, statefulset *appsv1.StatefulSet, log logr.Logger, additionalLabels ...map[string]string) error {
	ls, err := metav1.LabelSelectorAsSelector(statefulset.Spec.Selector)
	if err != nil {
		return err
	}

	pods := &corev1.PodList{}
	err = c.List(context.Background(), pods, client.MatchingLabelsSelector{
		Selector: ls,
	})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if len(pod.OwnerReferences) != 1 {
			log.V(1).Info("evaluting pod for detaching", "action", "skip", "statefulsetName", statefulset.Name, "name", pod.Name, "reason", "notExactlyOneOwnerReferenceForPod")
			continue
		}
		ownerRef := pod.OwnerReferences[0]
		if ownerRef.Kind != "ReplicaSet" {
			log.V(1).Info("evaluting pod for detaching", "action", "skip", "statefulsetName", statefulset.Name, "name", pod.Name, "reason", "ownerIsNotReplicaSet")
			continue
		}
		rs := &appsv1.ReplicaSet{}
		err = c.Get(context.Background(), client.ObjectKey{
			Name:      ownerRef.Name,
			Namespace: pod.Namespace,
		}, rs)
		if err != nil {
			return err
		}

		if len(rs.OwnerReferences) != 1 {
			log.V(1).Info("evaluting pod for detaching", "action", "skip", "statefulsetName", statefulset.Name, "name", pod.Name, "reason", "notExactlyOneOwnerReferenceForReplicaSet")
			continue
		}

		if rs.OwnerReferences[0].UID != statefulset.UID {
			log.V(1).Info("evaluting pod for detaching", "action", "skip", "statefulsetName", statefulset.Name, "name", pod.Name, "reason", "replicaSetOwnerMismatch")
			continue
		}

		log.V(1).Info("evaluting pod for detaching", "action", "detach", "statefulsetName", statefulset.Name, "name", pod.Name)

		p := pod.DeepCopy()
		p.OwnerReferences = nil
		if p.Labels == nil {
			p.Labels = make(map[string]string)
		} else {
			delete(p.Labels, "pod-template-hash")
		}
		p.Labels = util.MergeStringMaps(p.Labels, util.MergeMultipleStringMaps(additionalLabels...))

		err = c.Update(context.Background(), p)
		if err != nil {
			return err
		}
	}

	return nil
}
