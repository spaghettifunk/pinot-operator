package server

import (
	"fmt"

	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	jvmDefaultOptions = "-XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime -Xloggc:/opt/pinot/gc-pinot-server.log -Dlog4j2.configurationFile=/opt/pinot/conf/pinot-server-log4j2.xml -Dplugins.dir=/opt/pinot/plugins"
)

func (r *Reconciler) statefulsets() runtime.Object {
	return &appsv1.StatefulSet{
		ObjectMeta: templates.ObjectMeta(statefulsetName, r.labels(), r.Config),
		Spec: appsv1.StatefulSetSpec{
			Replicas:            r.Config.Spec.Server.ReplicaCount,
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
					Labels:      r.labels(),
					Annotations: r.Config.Spec.Server.PodAnnotations,
				},
				Spec: apiv1.PodSpec{
					NodeSelector:  r.Config.Spec.Server.NodeSelector,
					Affinity:      r.Config.Spec.Server.Affinity,
					Tolerations:   r.Config.Spec.Server.Tolerations,
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
	// TODO: fix name of zookeeper server
	zookeeperURL := fmt.Sprintf("%s:%s", "pinot-zookeeper", "2181")
	containers := []apiv1.Container{
		{
			Name:            componentName,
			Image:           *r.Config.Spec.Image,
			ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
			Args: []string{
				"StartServer",
				"-clusterName",
				r.Config.Spec.ClusterName,
				"-zkAddress",
				zookeeperURL,
				"-configFileName",
				"/var/pinot/server/config/pinot-server.conf",
			},
			Env:       r.envs(),
			Resources: r.Config.Spec.Server.Resources,
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("", r.Config.Spec.Server.Service.Port, apiv1.ProtocolTCP),
				templates.DefaultContainerPort("", r.Config.Spec.Server.AdminPort, apiv1.ProtocolTCP),
			},
			VolumeMounts: r.volumeMounts(),
		},
	}
	return containers
}

func (r *Reconciler) envs() []apiv1.EnvVar {
	jvm := fmt.Sprintf("%s %s", r.Config.Spec.Server.JvmOptions, jvmDefaultOptions)
	envs := []apiv1.EnvVar{
		{
			Name:  "JAVA_OPTS",
			Value: jvm,
		},
	}
	return append(envs, r.Config.Spec.Server.Env...)
}

func (r *Reconciler) volumes() []apiv1.Volume {
	volumes := []apiv1.Volume{
		{
			Name: serverConfigVolumeName,
			VolumeSource: apiv1.VolumeSource{
				ConfigMap: &apiv1.ConfigMapVolumeSource{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: configmapName,
					},
				},
			},
		},
	}
	return append(volumes, r.Config.Spec.Server.Volumes...)
}

func (r *Reconciler) volumeMounts() []apiv1.VolumeMount {
	volumeMounts := []apiv1.VolumeMount{
		{
			Name:      serverConfigVolumeName,
			MountPath: "/var/pinot/server/config",
		},
		{
			Name:      serverDataVolumeName,
			MountPath: "/var/pinot/server/data",
		},
	}
	return append(volumeMounts, r.Config.Spec.Server.VolumeMounts...)
}

func (r *Reconciler) volumeClaimTemplates() []apiv1.PersistentVolumeClaim {
	pvc := []apiv1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: serverDataVolumeName,
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
				Resources: apiv1.ResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceStorage: resource.MustParse(r.Config.Spec.Server.DiskSize),
					},
				},
			},
		},
	}
	return append(pvc, r.Config.Spec.Server.VolumeClaimTemplates...)
}
