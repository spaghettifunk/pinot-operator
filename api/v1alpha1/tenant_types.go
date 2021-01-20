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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantSpec defines the desired state of Tenant
type TenantSpec struct {
	// +optional
	PinotServer *NamespacedName `json:"pinotServer"`
	// +optional
	Labels map[string]string `json:"labels"`
}

// TenantStatus defines the observed state of Tenant
type TenantStatus struct {
	Status       ConfigState `json:"Status,omitempty"`
	ErrorMessage string      `json:"ErrorMessage,omitempty"`
}

// NamespacedName contains reference to a resource
type NamespacedName struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// +kubebuilder:object:root=true

// Tenant is the Schema for the Tenants API
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
