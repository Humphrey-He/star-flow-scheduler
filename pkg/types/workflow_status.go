package types

type WorkflowStatus string

const (
	WorkflowStatusPending  WorkflowStatus = "pending"
	WorkflowStatusRunning  WorkflowStatus = "running"
	WorkflowStatusSuccess  WorkflowStatus = "success"
	WorkflowStatusFailed   WorkflowStatus = "failed"
	WorkflowStatusCanceled WorkflowStatus = "canceled"
)
