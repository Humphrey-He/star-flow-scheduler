package runtime

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/handler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
)

type Runtime struct {
	queue  *TaskQueue
	pool   *WorkerPool
	config Config
}

type Config struct {
	WorkerCount      int
	QueueSize        int
	DefaultTimeoutMs int64
}

func NewRuntime(cfg Config, registry *handler.Registry, reporter Reporter) *Runtime {
	queue := NewTaskQueue(cfg.QueueSize)
	pool := NewWorkerPool(queue, registry, reporter, cfg.WorkerCount, time.Duration(cfg.DefaultTimeoutMs)*time.Millisecond)
	return &Runtime{
		queue:  queue,
		pool:   pool,
		config: cfg,
	}
}

func (r *Runtime) Start(ctx context.Context) {
	r.pool.Start(ctx)
}

func (r *Runtime) Enqueue(task *model.Task) error {
	return r.queue.Enqueue(task)
}
