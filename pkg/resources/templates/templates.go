package templates

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/spaghettifunk/pinot-operator/api/v1alpha1"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
)

// ObjectMeta .
func ObjectMeta(name string, labels map[string]string, config runtime.Object) metav1.ObjectMeta {
	obj := config.DeepCopyObject()
	objMeta, _ := meta.Accessor(obj)
	ovk := config.GetObjectKind().GroupVersionKind()

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: objMeta.GetNamespace(),
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion:         ovk.GroupVersion().String(),
				Kind:               ovk.Kind,
				Name:               objMeta.GetName(),
				UID:                objMeta.GetUID(),
				Controller:         util.BoolPointer(true),
				BlockOwnerDeletion: util.BoolPointer(true),
			},
		},
	}
}

// ObjectMetaNamespace .
func ObjectMetaNamespace(name, namespace string, labels map[string]string, config runtime.Object) metav1.ObjectMeta {
	obj := config.DeepCopyObject()
	objMeta, _ := meta.Accessor(obj)
	ovk := config.GetObjectKind().GroupVersionKind()

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion:         ovk.GroupVersion().String(),
				Kind:               ovk.Kind,
				Name:               objMeta.GetName(),
				UID:                objMeta.GetUID(),
				Controller:         util.BoolPointer(true),
				BlockOwnerDeletion: util.BoolPointer(true),
			},
		},
	}
}

// ObjectMetaWithAnnotations .
func ObjectMetaWithAnnotations(name string, labels map[string]string, annotations map[string]string, config runtime.Object) metav1.ObjectMeta {
	o := ObjectMeta(name, labels, config)
	o.Annotations = annotations
	return o
}

func ObjectMetaWithRevision(name string, labels map[string]string, config *v1alpha1.Pinot) metav1.ObjectMeta {
	return ObjectMeta(config.WithRevision(name), util.MergeStringMaps(labels, config.RevisionLabels()), config)
}

// ObjectMetaClusterScope .
func ObjectMetaClusterScope(name string, labels map[string]string, config runtime.Object) metav1.ObjectMeta {
	obj := config.DeepCopyObject()
	objMeta, _ := meta.Accessor(obj)
	ovk := config.GetObjectKind().GroupVersionKind()

	return metav1.ObjectMeta{
		Name:   name,
		Labels: labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion:         ovk.GroupVersion().String(),
				Kind:               ovk.Kind,
				Name:               objMeta.GetName(),
				UID:                objMeta.GetUID(),
				Controller:         util.BoolPointer(true),
				BlockOwnerDeletion: util.BoolPointer(true),
			},
		},
	}
}
