/*
Copyright 2021 The Pinot Operator authors.

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

	"github.com/spaghettifunk/pinot-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// OperatorVersion current operator version
const OperatorVersion = "v0.0.1"

// ServicePort extends the corev1 ServicePort object
type ServicePort struct {
	corev1.ServicePort `json:",inline"`
	TargetPort         *int32 `json:"targetPort,omitempty"`
}

// ServicePorts is an array of ServicePort
type ServicePorts []ServicePort

// Convert wraps the corev1.ServicePort object into a ServicePort
func (ps ServicePorts) Convert() []corev1.ServicePort {
	ports := make([]corev1.ServicePort, 0)
	for _, po := range ps {
		port := corev1.ServicePort{
			Name:     po.Name,
			Protocol: po.Protocol,
			Port:     po.Port,
			NodePort: po.NodePort,
		}
		if po.TargetPort != nil {
			port.TargetPort = intstr.FromInt(int(util.PointerToInt32(po.TargetPort)))
		}
		ports = append(ports, port)
	}

	return ports
}

// CommonResourceConfiguration defines basic K8s resource spec configurations
type CommonResourceConfiguration struct {
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Optional: node selector to be used by Pinot statefulsets
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Optional: affinity to be used to for enabling node, pod affinity and anti-affinity
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Optional: toleration to be used in order to run Pinot on nodes tainted
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// Optional: custom annotations to be populated in Pinot pods
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`
	// Optional: By default it is set to "parallel"
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`
	// Optional: custom labels to be populated in Pinot pods
	PodLabels       map[string]string       `json:"podLabels,omitempty"`
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
	ReplicaCount    *int32                  `json:"replicaCount,omitempty"`
	Env             []corev1.EnvVar         `json:"env,omitempty"`
	// Optional: If set to true then operator checks the rollout status of previous version StateSets before updating next.
	// Used only for updates.
	RollingDeploy bool `json:"rollingDeploy,omitempty"`
	// Optional
	UpdateStrategy *appsv1.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`
	// Optional, port is set to pinot.port if not specified with httpGet handler
	LivenessProbe *v1.Probe `json:"livenessProbe,omitempty"`
	// Optional, port is set to pinot.port if not specified with httpGet handler
	ReadinessProbe *v1.Probe `json:"readinessProbe,omitempty"`
	// Optional: StartupProbe for nodeSpec
	StartUpProbes *v1.Probe `json:"startUpProbes,omitempty"`
	// Optional: volumes etc for the Pinot pods
	VolumeClaimTemplates []v1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	VolumeMounts         []v1.VolumeMount           `json:"volumeMounts,omitempty"`
	Volumes              []v1.Volume                `json:"volumes,omitempty"`
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
	ClusterName string       `json:"clusterName"`
	Version     PinotVersion `json:"version"`
	Image       *string      `json:"image,omitempty"`
	// Optional: imagePullSecrets for private registries
	ImagePullSecrets []v1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Optional: image pull policy for the docker image
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Optional: log4j config file directory
	Log4jConfigPath string `json:"log4j.path,omitempty"`
	// Components for the Pinot cluster
	Controller  *ControllerConfiguration  `json:"controller"`
	Broker      *BrokerConfiguration      `json:"broker"`
	Server      *ServerConfiguration      `json:"server"`
	Zookeeper   *ZookeeperConfiguration   `json:"zookeeper"`
	DeepStorage *DeepStorageConfiguration `json:"deepStorage,omitempty"`
}

// ControllerConfiguration defines the k8s spec configuration for the Pinot controller
type ControllerConfiguration struct {
	CommonResourceConfiguration `json:",inline"`
	DiskSize                    string                               `json:"diskSize,omitempty"`
	JvmOptions                  string                               `json:"jvmOptions,omitempty"`
	VIPHost                     string                               `json:"vip.host,omitempty"`
	VIPPort                     string                               `json:"vip.port,omitempty"`
	Service                     ServiceResourceConfiguration         `json:"service,omitempty"`
	ExternalService             ExternalServiceResourceConfiguration `json:"externalService,omitempty"`
}

// BrokerConfiguration defines the k8s spec configuration for the Pinot broker
type BrokerConfiguration struct {
	CommonResourceConfiguration `json:",inline"`
	JvmOptions                  string                               `json:"jvmOptions,omitempty"`
	Service                     ServiceResourceConfiguration         `json:"service,omitempty"`
	ExternalService             ExternalServiceResourceConfiguration `json:"externalService,omitempty"`
}

// ServerConfiguration defines the k8s spec configuration for the Pinot server
type ServerConfiguration struct {
	CommonResourceConfiguration `json:",inline"`
	DiskSize                    string                       `json:"diskSize,omitempty"`
	JvmOptions                  string                       `json:"jvmOptions,omitempty"`
	Service                     ServiceResourceConfiguration `json:"service,omitempty"`
	AdminPort                   int                          `json:"adminPort,omitport"`
}

// ZookeeperConfiguration defines the desired state of Zookeeper
type ZookeeperConfiguration struct {
	Image      *string
	Replicas   int                          `json:"replicas,omitempty"`
	Resources  *corev1.ResourceRequirements `json:"resources,omitempty"`
	Storage    *zookeeperStorage            `json:"storage,omitempty"`
	JvmOptions string                       `json:"jvmOptions,omitempty"`
}

// zookeeperStorage defines the inner parameters for setting up the storage
type zookeeperStorage struct {
	Size string `json:"size,omitempty"`
}

// DeepStorageConfiguration defines the desired state of the DeepStorege
type DeepStorageConfiguration struct {
	Spec json.RawMessage `json:"spec"`
}

// PinotStatus defines the observed state of Pinot
type PinotStatus struct {
	Status       ConfigState `json:"Status,omitempty"`
	ErrorMessage string      `json:"ErrorMessage,omitempty"`
}

// +kubebuilder:object:root=true

// Pinot is the Schema for the pinots API
type Pinot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PinotSpec   `json:"spec,omitempty"`
	Status PinotStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PinotList contains a list of Pinot
type PinotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pinot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pinot{}, &PinotList{})
}
