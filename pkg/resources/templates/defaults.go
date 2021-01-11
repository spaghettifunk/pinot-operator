package templates

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/spaghettifunk/pinot-operator/pkg/util"
)

// DefaultAnnotations are the default annotations for deployments
func DefaultAnnotations(version string) map[string]string {
	return map[string]string{
		"pinot.cluster/created-by": version,
	}
}

// GetResourcesRequirementsOrDefault sets the new resources constraints or use the defaults
func GetResourcesRequirementsOrDefault(requirements *apiv1.ResourceRequirements, defaults *apiv1.ResourceRequirements) apiv1.ResourceRequirements {
	if requirements != nil {
		return *requirements
	}
	return *defaults
}

// DefaultRollingUpdateStrategy defines the default rolling update strategy
func DefaultRollingUpdateStrategy() appsv1.DeploymentStrategy {
	return appsv1.DeploymentStrategy{
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       util.IntstrPointer(1),
			MaxUnavailable: util.IntstrPointer(0),
		},
	}
}

// DefaultReadinessProbe returns the default readiness probe values
func DefaultReadinessProbe(path string, port, failureThreshold, timeoutSeconds int) *apiv1.Probe {
	return &apiv1.Probe{
		TimeoutSeconds:   int32(timeoutSeconds),
		FailureThreshold: int32(failureThreshold),
		Handler: apiv1.Handler{
			HTTPGet: &apiv1.HTTPGetAction{
				Path: path,
				Port: intstr.FromInt(port),
			},
		},
	}
}

// DefaultLivenessProbe returns the default liveness probe values
func DefaultLivenessProbe(path string, port, initialDelaySeconds, timeoutSeconds int) *apiv1.Probe {
	return &apiv1.Probe{
		TimeoutSeconds:      int32(timeoutSeconds),
		InitialDelaySeconds: int32(initialDelaySeconds),
		Handler: apiv1.Handler{
			HTTPGet: &apiv1.HTTPGetAction{
				Path: path,
				Port: intstr.FromInt(port),
			},
		},
	}
}

// DefaultServicePort returns the default values for the servicePort
func DefaultServicePort(name string, port, targetPort int) apiv1.ServicePort {
	return apiv1.ServicePort{
		Name:       name,
		Port:       int32(port),
		TargetPort: intstr.FromInt(targetPort),
	}
}

// DefaultContainerPort returns the default values for the containerPort
func DefaultContainerPort(name string, containerPort int) apiv1.ContainerPort {
	return apiv1.ContainerPort{
		Name:          name,
		ContainerPort: int32(containerPort),
	}
}
