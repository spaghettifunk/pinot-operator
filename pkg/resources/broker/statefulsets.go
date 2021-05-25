package broker

import (
	"fmt"

	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	jvmDefaultOptions = "-XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime -Xloggc:/opt/pinot/gc-pinot-broker.log -Dlog4j2.configurationFile=/opt/pinot/conf/pinot-broker-log4j2.xml -Dplugins.dir=/opt/pinot/plugins"
)

func (r *Reconciler) statefulsets() runtime.Object {
	return &appsv1.StatefulSet{
		ObjectMeta: templates.ObjectMeta(statefulsetName, r.labels(), r.Config),
		Spec: appsv1.StatefulSetSpec{
			Replicas:            r.Config.Spec.Broker.ReplicaCount,
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
					Annotations: r.Config.Spec.Broker.PodAnnotations,
				},
				Spec: apiv1.PodSpec{
					NodeSelector:  r.Config.Spec.Broker.NodeSelector,
					Affinity:      r.Config.Spec.Broker.Affinity,
					Tolerations:   r.Config.Spec.Broker.Tolerations,
					RestartPolicy: apiv1.RestartPolicyAlways,
					Containers:    r.containers(),
					Volumes:       r.volumes(),
				},
			},
			VolumeClaimTemplates: r.Config.Spec.Broker.VolumeClaimTemplates,
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
				"StartBroker",
				"-clusterName",
				r.Config.Spec.ClusterName,
				"-zkAddress",
				zookeeperURL,
				"-configFileName",
				"/var/pinot/broker/config/pinot-broker.conf",
			},
			Env:            r.envs(),
			Resources:      r.Config.Spec.Broker.Resources,
			LivenessProbe:  templates.DefaultLivenessProbe("/health", r.Config.Spec.Broker.Service.Port, 60, 30),
			ReadinessProbe: templates.DefaultReadinessProbe("/health", r.Config.Spec.Broker.Service.Port, 60, 30),
			Ports: []apiv1.ContainerPort{
				templates.DefaultContainerPort("", r.Config.Spec.Broker.Service.Port, apiv1.ProtocolTCP),
			},
			VolumeMounts: r.volumeMounts(),
		},
	}
	return containers
}

func (r *Reconciler) envs() []apiv1.EnvVar {
	jvm := fmt.Sprintf("%s %s", r.Config.Spec.Broker.JvmOptions, jvmDefaultOptions)
	envs := []apiv1.EnvVar{
		{
			Name:  "JAVA_OPTS",
			Value: jvm,
		},
	}
	return append(envs, r.Config.Spec.Broker.Env...)
}

func (r *Reconciler) volumes() []apiv1.Volume {
	volumes := []apiv1.Volume{
		{
			Name: brokerConfigVolumeName,
			VolumeSource: apiv1.VolumeSource{
				ConfigMap: &apiv1.ConfigMapVolumeSource{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: configmapName,
					},
				},
			},
		},
	}
	return append(volumes, r.Config.Spec.Broker.Volumes...)
}

func (r *Reconciler) volumeMounts() []apiv1.VolumeMount {
	volumeMounts := []apiv1.VolumeMount{
		{
			Name:      brokerConfigVolumeName,
			MountPath: "/var/pinot/broker/config",
		},
	}
	return append(volumeMounts, r.Config.Spec.Broker.VolumeMounts...)
}
