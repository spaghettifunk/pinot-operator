package k8sutil

type DesiredState string

const (
	DesiredStatePresent DesiredState = "present"
	DesiredStateAbsent  DesiredState = "absent"
	DesiredStateExists  DesiredState = "exists"
)
