package types

type ExecutorStatus string

const (
	ExecutorStatusOnline   ExecutorStatus = "online"
	ExecutorStatusOffline  ExecutorStatus = "offline"
	ExecutorStatusDraining ExecutorStatus = "draining"
)
