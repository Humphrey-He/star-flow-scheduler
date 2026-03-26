package types

type InstanceStatus string

const (
	InstanceStatusPending    InstanceStatus = "pending"
	InstanceStatusDispatched InstanceStatus = "dispatched"
	InstanceStatusRunning    InstanceStatus = "running"
	InstanceStatusSuccess    InstanceStatus = "success"
	InstanceStatusFailed     InstanceStatus = "failed"
	InstanceStatusRetryWait  InstanceStatus = "retry_wait"
	InstanceStatusDead       InstanceStatus = "dead"
	InstanceStatusCanceled   InstanceStatus = "canceled"
)
