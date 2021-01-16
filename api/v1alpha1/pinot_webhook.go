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
	"fmt"

	"github.com/spaghettifunk/pinot-operator/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	pinotImageHub     = "apachepinot/pinot"
	pinotImageVersion = "latest"
	// replicas
	defaultReplicaCount = 1
	defaultMinReplicas  = 1
	defaultMaxReplicas  = 5
)

var defaultResources = &apiv1.ResourceRequirements{
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("512m"),
		apiv1.ResourceMemory: resource.MustParse("2Gi"),
	},
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("256m"),
		apiv1.ResourceMemory: resource.MustParse("1Gi"),
	},
}

var zookeeperDefaultResources = &apiv1.ResourceRequirements{
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("512m"),
		apiv1.ResourceMemory: resource.MustParse("2Gi"),
	},
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("256m"),
		apiv1.ResourceMemory: resource.MustParse("1Gi"),
	},
}

// log is for logging in this package.
var pinotlog = logf.Log.WithName("pinot-resource")

// SetupWebhookWithManager injects the webhook into the manager client
func (r *Pinot) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-operators-apache-io-v1alpha1-pinot,mutating=true,failurePolicy=fail,groups=operators.apache.io,resources=pinots,verbs=create;update,versions=v1alpha1,name=mpinot.kb.io

var _ webhook.Defaulter = &Pinot{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Pinot) Default() {
	pinotlog.Info("default", "name", r.Name)

	// common
	if r.Spec.Image == nil {
		r.Spec.Image = util.StrPointer(fmt.Sprintf("%s:%s", pinotImageHub, pinotImageVersion))
	}
	// controller
	if r.Spec.Controller == nil {
		r.Spec.Controller = &ControllerConfiguration{}
	}
	if r.Spec.Controller.Resources == nil {
		r.Spec.Controller.Resources = defaultResources
	}
	if r.Spec.Controller.Service.Port == 0 {
		r.Spec.Controller.Service.Port = 9000
	}
	if r.Spec.Controller.DiskSize == "" {
		r.Spec.Controller.DiskSize = "1G"
	}
	if r.Spec.Controller.JvmOptions == "" {
		r.Spec.Controller.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	// broker
	if r.Spec.Broker == nil {
		r.Spec.Broker = &BrokerConfiguration{}
	}
	if r.Spec.Broker.Resources == nil {
		r.Spec.Broker.Resources = defaultResources
	}
	if r.Spec.Broker.JvmOptions == "" {
		r.Spec.Broker.JvmOptions = "-Xms256M -Xmx1G -XX:+UseG1GC -XX:MaxGCPauseMillis=200"
	}
	if r.Spec.Broker.Service.Port == 0 {
		r.Spec.Broker.Service.Port = 8099
	}
	if r.Spec.Broker.ExternalService.Port == 0 {
		r.Spec.Broker.ExternalService.Port = 8099
	}
	// server
	if r.Spec.Server == nil {
		r.Spec.Server = &ServerConfiguration{}
	}
	if r.Spec.Server.Service.Port == 0 {
		r.Spec.Server.Service.Port = 8098
	}
	if r.Spec.Server.AdminPort == 0 {
		r.Spec.Server.AdminPort = 8097
	}
	if r.Spec.Server.Resources == nil {
		r.Spec.Server.Resources = defaultResources
	}
	if r.Spec.Server.DiskSize == "" {
		r.Spec.Server.DiskSize = "4G"
	}
	// zookeeper
	if r.Spec.Zookeeper == nil {
		r.Spec.Zookeeper = &ZookeeperConfiguration{
			Storage: &zookeeperStorage{},
		}
	}
	if r.Spec.Zookeeper.Replicas == 0 {
		r.Spec.Zookeeper.Replicas = 1
	}
	if r.Spec.Zookeeper.Resources == nil {
		r.Spec.Zookeeper.Resources = zookeeperDefaultResources
	}
	if r.Spec.Zookeeper.JvmOptions == "" {
		r.Spec.Zookeeper.JvmOptions = "-Xmx2G -Xms2G"
	}
	if r.Spec.Zookeeper.Storage.Size == "" {
		r.Spec.Zookeeper.Storage.Size = "5Gi"
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-operators-apache-io-v1alpha1-pinot,mutating=false,failurePolicy=fail,groups=operators.apache.io,resources=pinots,versions=v1alpha1,name=vpinot.kb.io

var _ webhook.Validator = &Pinot{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Pinot) ValidateCreate() error {
	pinotlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Pinot) ValidateUpdate(old runtime.Object) error {
	pinotlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Pinot) ValidateDelete() error {
	pinotlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
