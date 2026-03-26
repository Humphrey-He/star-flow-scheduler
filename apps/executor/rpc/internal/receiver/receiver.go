package receiver

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/runtime"
)

type Receiver struct {
	runtime *runtime.Runtime
}

func NewReceiver(rt *runtime.Runtime) *Receiver {
	return &Receiver{runtime: rt}
}

func (r *Receiver) Accept(ctx context.Context, task *model.Task) error {
	return r.runtime.Enqueue(task)
}
