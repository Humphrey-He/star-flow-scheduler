package workflow

import (
	"fmt"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
)

type DependencyResolver interface {
	CanTrigger(triggerCondition string, upstreamStatuses []types.WorkflowNodeStatus) (bool, error)
}

type defaultResolver struct{}

func NewResolver() DependencyResolver {
	return &defaultResolver{}
}

func (r *defaultResolver) CanTrigger(triggerCondition string, upstreamStatuses []types.WorkflowNodeStatus) (bool, error) {
	switch triggerCondition {
	case "all_success":
		for _, status := range upstreamStatuses {
			if status != types.WorkflowNodeStatusSuccess {
				return false, nil
			}
		}
		return true, nil
	case "any_success":
		for _, status := range upstreamStatuses {
			if status == types.WorkflowNodeStatusSuccess {
				return true, nil
			}
		}
		return false, nil
	case "all_finished":
		for _, status := range upstreamStatuses {
			if !isFinished(status) {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("unsupported trigger condition: %s", triggerCondition)
	}
}

func isFinished(status types.WorkflowNodeStatus) bool {
	switch status {
	case types.WorkflowNodeStatusSuccess,
		types.WorkflowNodeStatusFailed,
		types.WorkflowNodeStatusSkipped,
		types.WorkflowNodeStatusCanceled:
		return true
	default:
		return false
	}
}
