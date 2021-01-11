package controller

import (
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Reconciler) statefulsets() runtime.Object {
	labels := util.MergeStringMaps(r.labels(), r.deploymentLabels())
	return &appsv1.Deployment{
		ObjectMeta: templates.ObjectMetaWithAnnotations(
			statefulsetName,
			util.MergeMultipleStringMaps(r.deploymentLabels(), r.labels()),
			templates.DefaultAnnotations(string(r.Config.Spec.Version)),
			r.Config,
		),
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				// TODO: enable only when podAntiAffinity is true
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{IntVal: 1},
				},
			},
			Replicas: r.Config.Spec.Controller.ReplicaCount,
			Selector: &v1.LabelSelector{
				MatchLabels: r.labels(),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: templates.DefaultAnnotations(string(r.Config.Spec.Version)),
				},
				Spec: apiv1.PodSpec{
					Containers: r.containers(),
					Volumes: []apiv1.Volume{
						{
							Name: "config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "pinot-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *Reconciler) containers() []apiv1.Container {

	args := []string{
		"public-api",
		"-controller-namespace=" + r.Config.Namespace,
		"-log-level=info",
	}

	controllerConfig := r.Config.Spec.Controller
	containers := []apiv1.Container{
		{
			Name:            "public-api",
			Image:           *r.Config.Spec.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args:            args,
			LivenessProbe:   templates.DefaultLivenessProbe("/ping", 9995, 10, 30),
			ReadinessProbe:  templates.DefaultReadinessProbe("/ready", 9995, 7, 30),
			Resources:       *controllerConfig.Resources,
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("http", 8085),
				templates.DefaultContainerPort("admin-http", 9995),
			},
			SecurityContext: &apiv1.SecurityContext{
				RunAsUser: util.Int64Pointer(2103),
			},
			TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
			TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		},
	}

	return containers
}
