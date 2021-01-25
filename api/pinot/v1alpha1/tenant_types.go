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
	// The tenant role to be used
	// +kubebuilder:validation:Enum=broker;server
	Role string `json:"role" protobuf:"byte,1,opt,name=role"`
	// Name of the tenant
	Name string `json:"name" protobuf:"byte,2,opt,name=name"`
	// Number of instances to be associated with the tenant. It is used only
	// when creating a tenant with Role Broker
	// +optional
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:default:=0
	NumberOfInstances *int32 `json:"numberOfInstances,omitempty" protobuf:"varint,3,opt,name=numberOfInstances"`
	// Number of Offline instances to be associted with the tenant. It is used only
	// when creating a tenant with Role Server
	// +optional
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:default:=0
	OfflineInstances *int32 `json:"offlineInstances,omitempty" protobuf:"varint,4,opt,name=offlineInstances"`
	// Number of Realtime instances to be associted with the tenant. It is used only
	// when creating a tenant with Role Server
	// +optional
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:default:=0
	RealtimeInstances *int32 `json:"realtimeInstances,omitempty" protobuf:"varint,5,opt,name=realtimeInstances"`
	// +optional
	PinotServer *NamespacedName `json:"pinotServer,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Tenant is the Schema for the Tenants API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.spec.role`
// +kubebuilder:printcolumn:name="Tenant Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.ErrorMessage",description="Error message"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Pinot Cluster",type="string",JSONPath=".spec.pinotServer"
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=tenants,shortName=tn
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
