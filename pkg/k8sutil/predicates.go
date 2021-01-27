/*
Copyright 2021 the Apache Pinot Kubernetes Operator authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8sutil

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	pinotv1alpha1 "github.com/spaghettifunk/pinot-operator/api/pinot/v1alpha1"
)

func GetWatchPredicateForPinot() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			old := e.ObjectOld.(*pinotv1alpha1.Pinot)
			new := e.ObjectNew.(*pinotv1alpha1.Pinot)
			if !reflect.DeepEqual(old.Spec, new.Spec) ||
				old.GetDeletionTimestamp() != new.GetDeletionTimestamp() ||
				old.GetGeneration() != new.GetGeneration() {
				return true
			}
			return false
		},
	}
}

func GetWatchPredicateForPinotServicePods() predicate.Funcs {
	return predicate.Funcs{
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			if _, ok := e.Meta.GetLabels()["pinot"]; ok {
				return true
			}
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if _, ok := e.MetaNew.GetLabels()["pinot"]; ok {
				return true
			}
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}

func GetWatchPredicateForPinotService(name string) predicate.Funcs {
	return predicate.Funcs{
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			if value, ok := e.Meta.GetLabels()["pinot"]; ok && value == name {
				return true
			}
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if value, ok := e.MetaNew.GetLabels()["pinot"]; ok && value == name {
				return true
			}
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}

func GetWatchPredicateForOwnedResources(owner runtime.Object, isController bool, scheme *runtime.Scheme, logger logr.Logger) predicate.Funcs {
	ownerMatcher := NewOwnerReferenceMatcher(owner, isController, scheme)
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// If a new namespace is created, we need to reconcile to mutate the injection labels
			if _, ok := e.Object.(*corev1.Namespace); ok {
				return true
			}
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// We don't want to run reconcile if a namespace is deleted
			if _, ok := e.Object.(*corev1.Namespace); ok {
				return false
			}
			related, object, err := ownerMatcher.Match(e.Object)
			if err != nil {
				logger.Error(err, "could not determine relation", "kind", e.Object.GetObjectKind())
			}
			if related {
				logger.Info("related object deleted", "trigger", object.GetName())
			}
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// If a namespace is updated, we need to reconcile to mutate the injection labels
			if _, ok := e.MetaNew.(*corev1.Namespace); ok {
				return true
			}
			related, object, err := ownerMatcher.Match(e.ObjectNew)
			if err != nil {
				logger.Error(err, "could not determine relation", "kind", e.ObjectNew.GetObjectKind())
			}
			if related {
				changed, err := IsObjectChanged(e.ObjectOld, e.ObjectNew, true)
				if err != nil {
					logger.Error(err, "could not check whether object is changed", "kind", e.ObjectNew.GetObjectKind())
				}
				if !changed {
					return false
				}

				logger.Info("related object changed", "trigger", object.GetName())
			}
			return true
		},
	}
}

func GetWatchPredicateForTenant() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if _, ok := e.ObjectOld.(*pinotv1alpha1.Tenant); !ok {
				return false
			}
			old := e.ObjectOld.(*pinotv1alpha1.Tenant)
			new := e.ObjectNew.(*pinotv1alpha1.Tenant)
			if !reflect.DeepEqual(old.Spec, new.Spec) ||
				old.GetDeletionTimestamp() != new.GetDeletionTimestamp() ||
				old.GetGeneration() != new.GetGeneration() {
				return true
			}
			return false
		},
	}
}

func GetWatchPredicateForSchema() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if _, ok := e.ObjectOld.(*pinotv1alpha1.Tenant); !ok {
				return false
			}
			old := e.ObjectOld.(*pinotv1alpha1.Tenant)
			new := e.ObjectNew.(*pinotv1alpha1.Tenant)
			if !reflect.DeepEqual(old.Spec, new.Spec) ||
				old.GetDeletionTimestamp() != new.GetDeletionTimestamp() ||
				old.GetGeneration() != new.GetGeneration() {
				return true
			}
			return false
		},
	}
}

func GetWatchPredicateForTable() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if _, ok := e.ObjectOld.(*pinotv1alpha1.Tenant); !ok {
				return false
			}
			old := e.ObjectOld.(*pinotv1alpha1.Tenant)
			new := e.ObjectNew.(*pinotv1alpha1.Tenant)
			if !reflect.DeepEqual(old.Spec, new.Spec) ||
				old.GetDeletionTimestamp() != new.GetDeletionTimestamp() ||
				old.GetGeneration() != new.GetGeneration() {
				return true
			}
			return false
		},
	}
}
