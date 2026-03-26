package handler

import "context"

type JobHandler interface {
	Name() string
	Execute(ctx context.Context, payload []byte) error
}
