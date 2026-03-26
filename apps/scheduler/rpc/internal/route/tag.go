package route

import (
	"context"
	"errors"
	"strings"
)

type TagStrategy struct{}

func (s *TagStrategy) Name() string {
	return "tag"
}

func (s *TagStrategy) Select(ctx context.Context, job JobSnapshot, executors []ExecutorNode) (*ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, errors.New("no executors")
	}
	if job.ExecutorTag == "" {
		return nil, errors.New("executor_tag required")
	}

	var candidates []ExecutorNode
	for _, exec := range executors {
		for _, tag := range exec.Tags {
			if strings.EqualFold(tag, job.ExecutorTag) {
				candidates = append(candidates, exec)
				break
			}
		}
	}
	if len(candidates) == 0 {
		return nil, errors.New("no executor matches tag")
	}

	best := candidates[0]
	for _, exec := range candidates[1:] {
		if exec.CurrentLoad < best.CurrentLoad {
			best = exec
		}
	}
	return &best, nil
}
