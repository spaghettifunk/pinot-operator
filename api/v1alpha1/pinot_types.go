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
	"encoding/json"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OperatorVersion current operator version
const OperatorVersion = "v0.0.1"

// CommonResourceConfiguration defines basic K8s resource spec configurations
type CommonResourceConfiguration struct {
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Node selector to be used by Pinot statefulsets
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Affinity scheduling rules to be applied on created Pods.
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Tolerations is the list of Toleration resources attached to each Pod in the Pinot cluster.
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// custom annotations to be populated in Pinot pods
	// +optional
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`
	// Optional: By default it is set to "parallel"
	// +optional
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`
	// custom labels to be populated in Pinot pods
	// +optional
	PodLabels map[string]string `json:"podLabels,omitempty"`
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:default:=1
	ReplicaCount *int32 `json:"replicaCount,omitempty"`
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Optional: If set to true then operator checks the rollout status of previous version StateSets before updating next.
	// Used only for updates.
	// +optional
	RollingDeploy bool `json:"rollingDeploy,omitempty"`
	// +optional
	UpdateStrategy *appsv1.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`
	// Port is set to pinot.port if not specified with httpGet handler
	// +optional
	LivenessProbe *v1.Probe `json:"livenessProbe,omitempty"`
	// Port is set to pinot.port if not specified with httpGet handler
	// +optional
	ReadinessProbe *v1.Probe `json:"readinessProbe,omitempty"`
	// StartupProbe for nodeSpec
	// +optional
	StartUpProbes *v1.Probe `json:"startUpProbes,omitempty"`
	// Volumes etc for the Pinot pods
	// +optional
	VolumeClaimTemplates []v1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	// +optional
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// +optional
	Volumes []v1.Volume `json:"volumes,omitempty"`
}

// ServiceResourceConfiguration defines some definition for a service resource
type ServiceResourceConfiguration struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Type        string            `json:"type,omitempty"`
	Port        int               `json:"port,omitempty"`
	NodePort    int               `json:"nodePort,omitempty"`
}

// ExternalServiceResourceConfiguration defines some definition for a service resource
type ExternalServiceResourceConfiguration struct {
	Enabled     bool              `json:"enabled,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Type        string            `json:"type,omitempty"`
	Port        int               `json:"port,omitempty"`
}

// PinotVersion stores the intended Pinot version
type PinotVersion string

// PinotSpec defines the desired state of Pinot
type PinotSpec struct {
	// Required: cluster name for the pinot deployment
	ClusterName string `json:"clusterName"`
	// +kubebuilder:default:="0.6.0"
	Version PinotVersion `json:"version"`
	// Image is the name of the Apache Pinot docker image to use for Brokers/Coordinator/Server nodes in the Pinot cluster.
	// Must be provided together with ImagePullSecrets in order to use an image in a private registry.
	// +kubebuilder:default:="apachepinot/pinot:latest"
	Image *string `json:"image,omitempty"`
	// List of Secret resource containing access credentials to the registry for the Apache Pinot image. Required if the docker registry is private.
	// +optional
	ImagePullSecrets []v1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Optional: image pull policy for the docker image
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Optional: log4j config file directory
	// +optional
	Log4jConfigPath string `json:"log4j.path,omitempty"`
	// The desired state of the Controller service to create for the cluster.
	// +optional
	Controller *ControllerConfiguration `json:"controller"`
	// The desired state of the Broker service to create for the cluster.
	// +optional
	Broker *BrokerConfiguration `json:"broker"`
	// The desired state of the Server service to create for the cluster.
	// +optional
	Server *ServerConfiguration `json:"server"`
	// The desired state of the Zookeeper service to create for the cluster.
	// +optional
	Zookeeper *ZookeeperConfiguration `json:"zookeeper"`
	// The desired state of the DeepStorage service to create for the cluster.
	// +optional
	DeepStorage *DeepStorageConfiguration `json:"deepStorage,omitempty"`
}

// ControllerConfiguration defines the k8s spec configuration for the Pinot controller
type ControllerConfiguration struct {
	// +optional
	CommonResourceConfiguration `json:",inline"`
	// +kubebuilder:default:="1Gi"
	DiskSize string `json:"diskSize,omitempty"`
	// +optional
	JvmOptions string `json:"jvmOptions,omitempty"`
	// +optional
	VIPHost string `json:"vip.host,omitempty"`
	// +optional
	VIPPort int `json:"vip.port,omitempty"`
	// +optional
	Service ServiceResourceConfiguration `json:"service,omitempty"`
	// +optional
	ExternalService ExternalServiceResourceConfiguration `json:"externalService,omitempty"`
}

// BrokerConfiguration defines the k8s spec configuration for the Pinot broker
type BrokerConfiguration struct {
	// +optional
	CommonResourceConfiguration `json:",inline"`
	// +optional
	JvmOptions string `json:"jvmOptions,omitempty"`
	// +optional
	Service ServiceResourceConfiguration `json:"service,omitempty"`
	// +optional
	ExternalService ExternalServiceResourceConfiguration `json:"externalService,omitempty"`
}

// ServerConfiguration defines the k8s spec configuration for the Pinot server
type ServerConfiguration struct {
	// +optional
	CommonResourceConfiguration `json:",inline"`
	// +kubebuilder:default:="4Gi"
	DiskSize string `json:"diskSize,omitempty"`
	// +optional
	JvmOptions string `json:"jvmOptions,omitempty"`
	// +optional
	Service ServiceResourceConfiguration `json:"service,omitempty"`
	// +optional
	AdminPort int `json:"adminPort,omitport"`
}

// ZookeeperConfiguration defines the desired state of Zookeeper
type ZookeeperConfiguration struct {
	// +optional
	Image *string `json:"-"`
	// +optional
	Replicas int `json:"replicas,omitempty"`
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// +kubebuilder:default:="5Gi"
	Storage *zookeeperStorage `json:"storage,omitempty"`
	// +optional
	JvmOptions string `json:"jvmOptions,omitempty"`
}

// zookeeperStorage defines the inner parameters for setting up the storage
type zookeeperStorage struct {
	// +optional
	Size string `json:"size,omitempty"`
}

// DeepStorageConfiguration defines the desired state of the DeepStorege
type DeepStorageConfiguration struct {
	// +optional
	Spec json.RawMessage `json:"spec"`
}

// PinotStatus defines the observed state of Pinot
type PinotStatus struct {
	Status       ConfigState `json:"Status,omitempty"`
	ErrorMessage string      `json:"ErrorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Pinot is the Schema for the pinots API
// +genclient
// +kubebuilder:resource:shortName=pn
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Pinot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PinotSpec   `json:"spec,omitempty"`
	Status PinotStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PinotList contains a list of Pinot
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PinotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pinot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pinot{}, &PinotList{})
}
