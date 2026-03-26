package route

import (
	"context"
)

type ExecutorNode struct {
	ID           int64
	ExecutorCode string
	Tags         []string
	CurrentLoad  int
}

type JobSnapshot struct {
	JobCode     string
	RouteKey    string
	ExecutorTag string
}

type Strategy interface {
	Name() string
	Select(ctx context.Context, job JobSnapshot, executors []ExecutorNode) (*ExecutorNode, error)
}

func NewStrategy(name string) Strategy {
	switch name {
	case "tag":
		return &TagStrategy{}
	case "hash":
		return &HashStrategy{}
	default:
		return &LeastLoadStrategy{}
	}
}
