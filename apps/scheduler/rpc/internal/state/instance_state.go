package state

import "fmt"

type InstanceStatus string

const (
	StatusPending    InstanceStatus = "pending"
	StatusDispatched InstanceStatus = "dispatched"
	StatusRunning    InstanceStatus = "running"
	StatusSuccess    InstanceStatus = "success"
	StatusFailed     InstanceStatus = "failed"
	StatusRetryWait  InstanceStatus = "retry_wait"
	StatusDead       InstanceStatus = "dead"
	StatusCanceled   InstanceStatus = "canceled"
)

var allowedTransitions = map[InstanceStatus]map[InstanceStatus]bool{
	StatusPending: {
		StatusDispatched: true,
		StatusCanceled:   true,
	},
	StatusDispatched: {
		StatusRunning:  true,
		StatusFailed:   true,
		StatusCanceled: true,
	},
	StatusRunning: {
		StatusSuccess: true,
		StatusFailed:  true,
	},
	StatusFailed: {
		StatusRetryWait: true,
		StatusDead:      true,
	},
	StatusRetryWait: {
		StatusDispatched: true,
		StatusDead:       true,
	},
}

func CanTransition(from, to InstanceStatus) bool {
	if nexts, ok := allowedTransitions[from]; ok {
		return nexts[to]
	}
	return false
}

func ValidateTransition(from, to InstanceStatus) error {
	if CanTransition(from, to) {
		return nil
	}
	return fmt.Errorf("invalid transition: %s -> %s", from, to)
}
