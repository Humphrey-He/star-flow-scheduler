package types

type ShardStatus string

const (
	ShardStatusPending    ShardStatus = "pending"
	ShardStatusDispatched ShardStatus = "dispatched"
	ShardStatusRunning    ShardStatus = "running"
	ShardStatusSuccess    ShardStatus = "success"
	ShardStatusFailed     ShardStatus = "failed"
	ShardStatusCanceled   ShardStatus = "canceled"
)
