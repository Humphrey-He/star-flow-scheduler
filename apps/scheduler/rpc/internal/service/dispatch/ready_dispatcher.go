package dispatch

import (
	"context"
	"errors"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReadyDispatcherConfig struct {
	PopTimeout time.Duration
	IdleSleep  time.Duration
	Requeue    time.Duration
}

type ReadyDispatcher struct {
	cfg        ReadyDispatcherConfig
	queue      redisx.ReadyQueue
	instances  *repo.JobInstanceRepository
	dispatcher *Service
}

func NewReadyDispatcher(cfg ReadyDispatcherConfig, queue redisx.ReadyQueue, instances *repo.JobInstanceRepository, dispatcher *Service) *ReadyDispatcher {
	return &ReadyDispatcher{
		cfg:        readyDispatcherDefaults(cfg),
		queue:      queue,
		instances:  instances,
		dispatcher: dispatcher,
	}
}

func (d *ReadyDispatcher) Start(ctx context.Context) {
	if d.queue == nil || d.instances == nil || d.dispatcher == nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		instanceNo, err := d.queue.Pop(ctx, d.cfg.PopTimeout)
		if err != nil {
			if errors.Is(err, redisx.ErrNotFound) {
				metricsx.Inc("scheduler_dispatch_pop_empty_total")
				time.Sleep(d.cfg.IdleSleep)
				continue
			}
			logx.WithContext(ctx).Errorf("ready dispatcher pop failed: %v", err)
			metricsx.Inc("scheduler_dispatch_pop_error_total")
			time.Sleep(d.cfg.IdleSleep)
			continue
		}
		metricsx.Inc("scheduler_dispatch_pop_total")

		instance, err := d.instances.GetByInstanceNo(ctx, instanceNo)
		if err != nil {
			logx.WithContext(ctx).Errorf("ready dispatcher load instance=%s failed: %v", instanceNo, err)
			metricsx.Inc("scheduler_dispatch_load_error_total")
			continue
		}
		if instance.Status != string(state.StatusPending) {
			continue
		}

		metricsx.Inc("scheduler_dispatch_total")
		if _, err := d.dispatcher.DispatchInstance(ctx, instanceNo); err != nil {
			logx.WithContext(ctx).Errorf("ready dispatcher dispatch instance=%s failed: %v", instanceNo, err)
			metricsx.Inc("scheduler_dispatch_fail_total")
			if err := d.queue.Push(ctx, instanceNo); err != nil {
				logx.WithContext(ctx).Errorf("ready dispatcher requeue instance=%s failed: %v", instanceNo, err)
				metricsx.Inc("scheduler_dispatch_requeue_error_total")
			}
			time.Sleep(d.cfg.Requeue)
			continue
		}
		metricsx.Inc("scheduler_dispatch_success_total")
	}
}

func readyDispatcherDefaults(cfg ReadyDispatcherConfig) ReadyDispatcherConfig {
	if cfg.PopTimeout <= 0 {
		cfg.PopTimeout = time.Second
	}
	if cfg.IdleSleep <= 0 {
		cfg.IdleSleep = 300 * time.Millisecond
	}
	if cfg.Requeue <= 0 {
		cfg.Requeue = time.Second
	}
	return cfg
}
