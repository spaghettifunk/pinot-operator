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
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// OperatorVersion current operator version
const OperatorVersion = "v0.0.1"

const (
	RevisionedAutoInjectionLabelKey = "pinot.io/rev"
)

func NamespacedNameFromRevision(revision string) types.NamespacedName {
	nn := types.NamespacedName{}
	p := strings.SplitN(revision, ".", 2)
	if len(p) == 2 {
		nn.Name = p[0]
		nn.Namespace = p[1]
	}

	return nn
}

// CommonResourceConfiguration defines basic K8s resource spec configurations
type CommonResourceConfiguration struct {
	// +kubebuilder:default:={limits: {cpu: "512m", memory: "2Gi"}, requests: {cpu: "256m", memory: "1Gi"}}
	Resources corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,1,opt,name=resources"`
	// Node selector to be used by Pinot statefulsets
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,2,opt,name=nodeSelector"`
	// Affinity scheduling rules to be applied on created Pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty" protobuf:"bytes,3,opt,name=affinity"`
	// Tolerations is the list of Toleration resources attached to each Pod in the Pinot cluster.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,4,opt,name=tolerations"`
	// PoAnnotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	PodAnnotations map[string]string `json:"podAnnotations,omitempty" protobuf:"bytes,5,opt,name=podAnnotations"`
	// PodManagementPolicy controls how pods are created during initial scale up,
	// when replacing pods on nodes, or when scaling down. The default policy is
	// `OrderedReady`, where pods are created in increasing order (pod-0, then
	// pod-1, etc) and the controller will wait until each pod is ready before
	// continuing. When scaling down, the pods are removed in the opposite order.
	// The alternative policy is `Parallel` which will create pods in parallel
	// to match the desired scale without waiting, and on scale down will delete
	// all pods at once.
	// +optional
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty" protobuf:"bytes,6,opt,name=podManagementPolicy,casttype=PodManagementPolicyType"`
	// Custom labels to be populated in Pinot pods
	// +optional
	PodLabels map[string]string `json:"podLabels,omitempty" protobuf:"bytes,7,opt,name=podLabels"`
	// Defines privilege and access control settings for a Pod or Container
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty" protobuf:"bytes,8,opt,name=securityContext"`
	// Replicas is the number of nodes in the service. Each node is deployed as a Replica in a StatefulSet. Only 1, 3, 5 replicas clusters are tested.
	// This value should be an odd number to ensure the resultant cluster can establish exactly one quorum of nodes
	// in the event of a fragmenting network partition.
	// +optional
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:default:=1
	ReplicaCount *int32 `json:"replicaCount,omitempty" protobuf:"varint,9,opt,name=replicaCount"`
	// Extra environment variables to pass to the service
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty" protobuf:"bytes,10,opt,name=env"`
	// If set to true then operator checks the rollout status of previous version StateSets before updating next.
	// Used only for updates.
	// +optional
	RollingDeploy bool `json:"rollingDeploy,omitempty" protobuf:"varint,11,opt,name=rollingDeploy"`
	// UpdateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	// +optional
	UpdateStrategy *appsv1.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty" protobuf:"bytes,12,opt,name=updateStrategy"`
	// Describes a health check to be performed against a container to determine whether it is alive or not
	// +optional
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty" protobuf:"bytes,13,opt,name=livenessProbe"`
	// Describes a health check to be performed against a container to determine whether it is ready to receive traffic or not
	// +optional
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty" protobuf:"bytes,14,opt,name=readinessProbe"`
	// VolumeClaimTemplates is a list of claims that pods are allowed to reference.
	// The StatefulSet controller is responsible for mapping network identities to
	// claims in a way that maintains the identity of a pod. Every claim in
	// this list must have at least one matching (by name) volumeMount in one
	// container in the template. A claim in this list takes precedence over
	// any volumes in the template, with the same name.
	// +optional
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty" protobuf:"bytes,15,opt,name=volumeClaimTemplates"`
	// Describes a mounting of a Volume within a container
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty" protobuf:"bytes,16,opt,name=volumeMounts"`
	// Represents a named volume in a pod that may be accessed by any container in the pod
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty" protobuf:"bytes,17,opt,name=volumes"`
}

// ServiceResourceConfiguration defines some definition for a service resource
type ServiceResourceConfiguration struct {
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Type of Service to create for the cluster. Must be one of: ClusterIP, LoadBalancer, NodePort.
	// For more info see https://pkg.go.dev/k8s.io/api/core/v1#ServiceType
	// +kubebuilder:validation:Enum=ClusterIP;LoadBalancer;NodePort
	// +kubebuilder:default:="ClusterIP"
	Type corev1.ServiceType `json:"type,omitempty"`
	// +optional
	Port int `json:"port,omitempty"`
	// +optional
	NodePort int `json:"nodePort,omitempty"`
}

// ExternalServiceResourceConfiguration defines some definition for a service resource
type ExternalServiceResourceConfiguration struct {
	// Whether enabling the external service or not
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Type of Service to create for the cluster. Must be one of: ClusterIP, LoadBalancer, NodePort.
	// For more info see https://pkg.go.dev/k8s.io/api/core/v1#ServiceType
	// +kubebuilder:validation:Enum=ClusterIP;LoadBalancer;NodePort
	// +kubebuilder:default:="LoadBalancer"
	Type corev1.ServiceType `json:"type,omitempty"`
	// +optional
	Port int `json:"port,omitempty"`
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
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Image pull policy for the docker image
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Log4j config file directory
	// +optional
	Log4jConfigPath string `json:"log4j.path,omitempty"`
	// The desired state of the Controller service to create for the cluster.
	// +optional
	Controller ControllerConfiguration `json:"controller"`
	// The desired state of the Broker service to create for the cluster.
	// +optional
	Broker BrokerConfiguration `json:"broker,omitempty"`
	// The desired state of the Server service to create for the cluster.
	// +optional
	Server ServerConfiguration `json:"server,omitempty"`
	// The desired state of the Zookeeper service to create for the cluster.
	// +optional
	Zookeeper ZookeeperConfiguration `json:"zookeeper,omitempty"`
}

// ControllerConfiguration defines the k8s spec configuration for the Pinot controller
type ControllerConfiguration struct {
	// +optional
	CommonResourceConfiguration `json:",inline"`
	// Size of the persisten disk for the controller service
	// +kubebuilder:default:="1Gi"
	DiskSize string `json:"diskSize,omitempty"`
	// Extra JVM parameters to be passed to the controller service
	// +kubebuilder:default:="-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	JvmOptions string `json:"jvmOptions,omitempty"`
	//
	// +optional
	VIPHost string `json:"vip.host,omitempty"`
	//
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
	// Extra JVM parameters to be passed to the controller service
	// +kubebuilder:default:="-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
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
	// Size of the persisten disk for the server service
	// +kubebuilder:default:="4Gi"
	DiskSize string `json:"diskSize,omitempty"`
	// Extra JVM parameters to be passed to the controller service
	// +kubebuilder:default:="-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	JvmOptions string `json:"jvmOptions,omitempty"`
	// +optional
	Service ServiceResourceConfiguration `json:"service,omitempty"`
	// Service port for the service controller
	// +optional
	AdminPort int `json:"adminPort,omitport"`
}

// ZookeeperConfiguration defines the desired state of Zookeeper
type ZookeeperConfiguration struct {
	// Image is the name of the Apache Zookeeper docker image
	// +kubebuilder:default:="zookeeper:3.5.5"
	// +optional
	Image *string `json:"image"`
	// ReplicaCount is the number of nodes in the zookeeper service. Each node is deployed as a Replica in a StatefulSet. Only 1, 3, 5 replicas clusters are tested.
	// This value should be an odd number to ensure the resultant cluster can establish exactly one quorum of nodes
	// in the event of a fragmenting network partition.
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:default:=1
	// +optional
	ReplicaCount *int32 `json:"replicaCount,omitempty"`
	// The desired compute resource requirements of Pods in the cluster.
	// +kubebuilder:default:={limits: {cpu: "512m", memory: "2Gi"}, requests: {cpu: "256m", memory: "1Gi"}}
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Defines the inner parameters for setting up the storage
	// +optional
	Storage zookeeperStorage `json:"storage,omitempty"`
	// Extra JVM parameters to be passed to the zookeeper service
	// +kubebuilder:default:="-Xmx2G -Xms2G"
	// +optional
	JvmOptions string `json:"jvmOptions,omitempty"`
}

// zookeeperStorage defines the inner parameters for setting up the storage
type zookeeperStorage struct {
	// The requested size of the persistent volume attached to each Pod in the RabbitmqCluster.
	// The format of this field matches that defined by kubernetes/apimachinery.
	// See https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Quantity for more info on the format of this field.
	// +kubebuilder:default:="5Gi"
	Size string `json:"storage,omitempty"`
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
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.clusterStatus"
// +kubebuilder:resource:shortName={"apc"}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Pinot struct {
	// Embedded metadata identifying a Kind and API Verison of an object.
	// For more info, see: https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#TypeMeta
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the Pinot Custom Resource.
	Spec PinotSpec `json:"spec,omitempty"`
	// Status presents the observed state of Pinot
	Status PinotStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PinotList contains a list of Pinot
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PinotList struct {
	// Embedded metadata identifying a Kind and API Verison of an object.
	// For more info, see: https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#TypeMeta
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Array of Pinot resources.
	Items []Pinot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pinot{}, &PinotList{})
}
