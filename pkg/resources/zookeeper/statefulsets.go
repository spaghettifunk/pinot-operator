package zookeeper

import (
	"fmt"
	"strconv"

	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"github.com/spaghettifunk/pinot-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	jvmDefaultOptions = "-XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime"
)

func (r *Reconciler) statefulsets() runtime.Object {
	return &appsv1.StatefulSet{
		ObjectMeta: templates.ObjectMetaWithAnnotations(statefulsetName, r.labels(), templates.DefaultAnnotations(string(r.Config.Spec.Version)), r.Config),
		Spec: appsv1.StatefulSetSpec{
			Replicas:            util.IntPointer(int32(r.Config.Spec.Zookeeper.Replicas)),
			ServiceName:         serviceHeadlessName,
			PodManagementPolicy: appsv1.ParallelPodManagement,
			Selector: &v1.LabelSelector{
				MatchLabels: r.labels(),
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: r.labels(),
				},
				Spec: apiv1.PodSpec{
					TerminationGracePeriodSeconds: util.Int64Pointer(1800),
					SecurityContext: &apiv1.PodSecurityContext{
						FSGroup:   util.Int64Pointer(1000),
						RunAsUser: util.Int64Pointer(1000),
					},
					RestartPolicy: apiv1.RestartPolicyAlways,
					Containers:    r.containers(),
					Volumes:       r.volumes(),
				},
			},
			VolumeClaimTemplates: r.volumeClaimTemplates(),
		},
	}
}

func (r *Reconciler) containers() []apiv1.Container {
	containers := []apiv1.Container{
		{
			Name:            componentName,
			Image:           *r.Config.Spec.Zookeeper.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Command: []string{
				"/bin/bash",
				"-xec",
				"/config-scripts/run",
			},
			Env:       r.envs(),
			Resources: *r.Config.Spec.Zookeeper.Resources,
			LivenessProbe: &apiv1.Probe{
				Handler: apiv1.Handler{
					Exec: &apiv1.ExecAction{
						Command: []string{
							"sh",
							"/config-scripts/ok",
						},
					},
				},
				InitialDelaySeconds: 20,
				PeriodSeconds:       30,
				TimeoutSeconds:      5,
				FailureThreshold:    2,
				SuccessThreshold:    1,
			},
			ReadinessProbe: &apiv1.Probe{
				Handler: apiv1.Handler{
					Exec: &apiv1.ExecAction{
						Command: []string{
							"sh",
							"/config-scripts/ready",
						},
					},
				},
				InitialDelaySeconds: 20,
				PeriodSeconds:       30,
				TimeoutSeconds:      5,
				FailureThreshold:    2,
				SuccessThreshold:    1,
			},
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("client", zookeeperClientPort, apiv1.ProtocolTCP),
				templates.DefaultContainerPort("election", zookeeperElectionPort, apiv1.ProtocolTCP),
				templates.DefaultContainerPort("server", zookeeperServerPort, apiv1.ProtocolTCP),
			},
			VolumeMounts: r.volumeMounts(),
		},
	}
	return containers
}

func (r *Reconciler) envs() []apiv1.EnvVar {
	jvm := fmt.Sprintf("%s %s", r.Config.Spec.Zookeeper.JvmOptions, jvmDefaultOptions)
	envs := []apiv1.EnvVar{
		{
			Name:  "ZK_REPLICAS",
			Value: strconv.Itoa(r.Config.Spec.Zookeeper.Replicas),
		},
		{
			Name:  "JMXAUTH",
			Value: "false",
		},
		{
			Name:  "JMXDISABLE",
			Value: "false",
		},
		{
			Name:  "JMXPORT",
			Value: "1099",
		},
		{
			Name:  "JMXSSL",
			Value: "false",
		},
		{
			Name:  "ZK_HEAP_SIZE",
			Value: "256M",
		},
		{
			Name:  "ZOO_STANDALONE_ENABLED",
			Value: "false",
		},
		{
			Name:  "JAVA_OPTS",
			Value: jvm,
		},
	}
	return envs
}

func (r *Reconciler) volumes() []apiv1.Volume {
	volumes := []apiv1.Volume{
		{
			Name: zookeeperConfigVolumeName,
			VolumeSource: apiv1.VolumeSource{
				ConfigMap: &apiv1.ConfigMapVolumeSource{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: configmapName,
					},
					DefaultMode: util.IntPointer(0555),
				},
			},
		},
	}
	return volumes
}

func (r *Reconciler) volumeMounts() []apiv1.VolumeMount {
	volumeMounts := []apiv1.VolumeMount{
		{
			Name:      zookeeperConfigVolumeName,
			MountPath: "/config-scripts",
		},
		{
			Name:      zookeeperDataVolumeName,
			MountPath: "/data",
		},
	}
	return volumeMounts
}

func (r *Reconciler) volumeClaimTemplates() []apiv1.PersistentVolumeClaim {
	pvc := []apiv1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: zookeeperDataVolumeName,
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
				Resources: apiv1.ResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceStorage: resource.MustParse(r.Config.Spec.Zookeeper.Storage.Size),
					},
				},
			},
		},
	}
	return append(pvc, r.Config.Spec.Controller.VolumeClaimTemplates...)
}
