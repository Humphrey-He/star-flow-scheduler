package workflow

import (
	"testing"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
)

func TestResolverConditions(t *testing.T) {
	tests := []struct {
		name      string
		condition string
		statuses  []types.WorkflowNodeStatus
		want      bool
	}{
		{"all_success_true", "all_success", []types.WorkflowNodeStatus{types.WorkflowNodeStatusSuccess, types.WorkflowNodeStatusSuccess}, true},
		{"all_success_false", "all_success", []types.WorkflowNodeStatus{types.WorkflowNodeStatusSuccess, types.WorkflowNodeStatusFailed}, false},
		{"any_success_true", "any_success", []types.WorkflowNodeStatus{types.WorkflowNodeStatusFailed, types.WorkflowNodeStatusSuccess}, true},
		{"any_success_false", "any_success", []types.WorkflowNodeStatus{types.WorkflowNodeStatusFailed, types.WorkflowNodeStatusCanceled}, false},
		{"all_finished_true", "all_finished", []types.WorkflowNodeStatus{types.WorkflowNodeStatusFailed, types.WorkflowNodeStatusSuccess}, true},
		{"all_finished_false", "all_finished", []types.WorkflowNodeStatus{types.WorkflowNodeStatusRunning, types.WorkflowNodeStatusSuccess}, false},
	}

	resolver := NewResolver()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.CanTrigger(tt.condition, tt.statuses)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %v got %v", tt.want, got)
			}
		})
	}
}
