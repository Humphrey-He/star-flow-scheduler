package types

type WorkflowNodeStatus string

const (
	WorkflowNodeStatusPending  WorkflowNodeStatus = "pending"
	WorkflowNodeStatusReady    WorkflowNodeStatus = "ready"
	WorkflowNodeStatusRunning  WorkflowNodeStatus = "running"
	WorkflowNodeStatusSuccess  WorkflowNodeStatus = "success"
	WorkflowNodeStatusFailed   WorkflowNodeStatus = "failed"
	WorkflowNodeStatusSkipped  WorkflowNodeStatus = "skipped"
	WorkflowNodeStatusCanceled WorkflowNodeStatus = "canceled"
)
