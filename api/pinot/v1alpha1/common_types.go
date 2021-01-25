package v1alpha1

// ConfigState describes the state of the operator
type ConfigState string

const (
	// Created status after the initialization
	Created ConfigState = "Created"
	// ReconcileFailed status when an error happened
	ReconcileFailed ConfigState = "ReconcileFailed"
	// Reconciling status when the reconciliation loop is started
	Reconciling ConfigState = "Reconciling"
	// Available status when the resources are created successfully
	Available ConfigState = "Available"
	// Unmanaged status when the system is not sure what to do
	Unmanaged ConfigState = "Unmanaged"
)
