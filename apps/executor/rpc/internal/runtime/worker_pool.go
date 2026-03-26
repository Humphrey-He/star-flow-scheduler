package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/handler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type Reporter interface {
	Report(result *model.TaskResult)
}

type WorkerPool struct {
	queue          *TaskQueue
	registry       *handler.Registry
	reporter       Reporter
	workerCount    int
	defaultTimeout time.Duration
	running        int64
}

func NewWorkerPool(queue *TaskQueue, registry *handler.Registry, reporter Reporter, workerCount int, defaultTimeout time.Duration) *WorkerPool {
	return &WorkerPool{
		queue:          queue,
		registry:       registry,
		reporter:       reporter,
		workerCount:    workerCount,
		defaultTimeout: defaultTimeout,
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		go p.worker(ctx, i)
	}
}

func (p *WorkerPool) worker(ctx context.Context, idx int) {
	logger := logx.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-p.queue.Channel():
			if task == nil {
				continue
			}
			logger.Infow("executor task received",
				logx.Field("worker", idx),
				logx.Field("instance_no", task.InstanceNo),
				logx.Field("job_code", task.JobCode),
				logx.Field("handler_name", task.HandlerName),
				logx.Field("shard_no", task.ShardNo),
				logx.Field("trace_id", task.TraceID),
			)
			atomic.AddInt64(&p.running, 1)
			p.executeTask(ctx, task)
			atomic.AddInt64(&p.running, -1)
		}
	}
}

func (p *WorkerPool) executeTask(ctx context.Context, task *model.Task) {
	start := time.Now()
	metricsx.Inc("executor_task_started_total")
	result := &model.TaskResult{
		InstanceNo: task.InstanceNo,
		ShardNo:    task.ShardNo,
		Status:     schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS,
		StartTime:  start,
	}

	h, ok := p.registry.Get(task.HandlerName)
	if !ok {
		result.Status = schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED
		result.ErrorCode = "handler_not_found"
		result.ErrorMessage = fmt.Sprintf("handler not registered: %s", task.HandlerName)
		result.FinishTime = time.Now()
		logx.WithContext(ctx).Errorw("executor handler not found",
			logx.Field("instance_no", task.InstanceNo),
			logx.Field("job_code", task.JobCode),
			logx.Field("handler_name", task.HandlerName),
			logx.Field("trace_id", task.TraceID),
		)
		p.reporter.Report(result)
		return
	}

	timeout := p.defaultTimeout
	if task.TimeoutMs > 0 {
		timeout = time.Duration(task.TimeoutMs) * time.Millisecond
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := p.runHandler(execCtx, h, task.Payload)
	if err != nil {
		result.Status = schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED
		result.ErrorCode = "execute_failed"
		result.ErrorMessage = err.Error()
		logx.WithContext(ctx).Errorw("executor task execute failed",
			logx.Field("instance_no", task.InstanceNo),
			logx.Field("job_code", task.JobCode),
			logx.Field("handler_name", task.HandlerName),
			logx.Field("error_message", err.Error()),
			logx.Field("trace_id", task.TraceID),
		)
	}

	if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
		result.Status = schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED
		result.ErrorCode = "timeout"
		result.ErrorMessage = "task timeout"
		logx.WithContext(ctx).Errorw("executor task timeout",
			logx.Field("instance_no", task.InstanceNo),
			logx.Field("job_code", task.JobCode),
			logx.Field("handler_name", task.HandlerName),
			logx.Field("trace_id", task.TraceID),
		)
	}

	result.FinishTime = time.Now()
	metricsx.ObserveDurationMs("executor_task_duration_ms", result.FinishTime.Sub(start))
	if result.Status == schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS {
		metricsx.Inc("executor_task_success_total")
	} else {
		metricsx.Inc("executor_task_fail_total")
	}
	p.reporter.Report(result)
}

func (p *WorkerPool) runHandler(ctx context.Context, h handler.JobHandler, payload []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return h.Execute(ctx, payload)
}

func (p *WorkerPool) Running() int64 {
	return atomic.LoadInt64(&p.running)
}
