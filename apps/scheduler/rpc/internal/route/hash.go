package route

import (
	"context"
	"errors"
	"hash/fnv"
)

type HashStrategy struct{}

func (s *HashStrategy) Name() string {
	return "hash"
}

func (s *HashStrategy) Select(ctx context.Context, job JobSnapshot, executors []ExecutorNode) (*ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, errors.New("no executors")
	}
	if job.RouteKey == "" {
		return nil, errors.New("route_key required")
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(job.RouteKey))
	idx := int(hash.Sum32()) % len(executors)
	exec := executors[idx]
	return &exec, nil
}
