package runtime

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/handler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
)

type Runtime struct {
	queue    *TaskQueue
	pool     *WorkerPool
	config   Config
	draining int32
}

type Config struct {
	WorkerCount        int
	QueueSize          int
	DefaultTimeoutMs   int64
	ShutdownTimeoutSec int64
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
	if r.IsDraining() {
		return errors.New("runtime draining")
	}
	return r.queue.Enqueue(task)
}

func (r *Runtime) Drain() {
	atomic.StoreInt32(&r.draining, 1)
}

func (r *Runtime) IsDraining() bool {
	return atomic.LoadInt32(&r.draining) == 1
}

func (r *Runtime) QueueSize() int {
	return r.queue.Len()
}

func (r *Runtime) RunningJobs() int64 {
	return r.pool.Running()
}

func (r *Runtime) Wait(ctx context.Context) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if r.QueueSize() == 0 && r.RunningJobs() == 0 {
				return nil
			}
		}
	}
}
