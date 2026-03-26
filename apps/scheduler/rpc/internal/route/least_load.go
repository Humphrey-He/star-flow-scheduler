package route

import (
	"context"
	"errors"
)

type LeastLoadStrategy struct{}

func (s *LeastLoadStrategy) Name() string {
	return "least_load"
}

func (s *LeastLoadStrategy) Select(ctx context.Context, job JobSnapshot, executors []ExecutorNode) (*ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, errors.New("no executors")
	}
	best := executors[0]
	for _, exec := range executors[1:] {
		if exec.CurrentLoad < best.CurrentLoad {
			best = exec
		}
	}
	return &best, nil
}
