package scanner

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

type Config struct {
	TickInterval time.Duration
	BatchSize    int64
	LockTTL      time.Duration
	RequeueDelay time.Duration
	FallbackInterval time.Duration
}

type DelayScanner struct {
	cfg         Config
	delayQueue  redisx.DelayQueue
	readyQueue  redisx.ReadyQueue
	locker      redisx.Locker
	instances   *repo.JobInstanceRepository
	lockKey     string
	lastFallback time.Time
}

func NewDelayScanner(cfg Config, delayQueue redisx.DelayQueue, readyQueue redisx.ReadyQueue, locker redisx.Locker, instances *repo.JobInstanceRepository) *DelayScanner {
	return &DelayScanner{
		cfg:         withDefaults(cfg),
		delayQueue:  delayQueue,
		readyQueue:  readyQueue,
		locker:      locker,
		instances:   instances,
		lockKey:     redisx.LockKey("scanner:delay"),
	}
}

func (s *DelayScanner) Start(ctx context.Context) {
	if s.delayQueue == nil || s.readyQueue == nil || s.locker == nil || s.instances == nil {
		return
	}
	ticker := time.NewTicker(s.cfg.TickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *DelayScanner) tick(ctx context.Context) {
	metricsx.Inc("scheduler_scanner_tick_total")
	ok, err := s.locker.TryLock(ctx, s.lockKey, s.cfg.LockTTL)
	if err != nil || !ok {
		if err == nil {
			metricsx.Inc("scheduler_scanner_lock_miss_total")
		} else {
			metricsx.Inc("scheduler_scanner_lock_error_total")
		}
		return
	}
	defer func() {
		_ = s.locker.Unlock(ctx, s.lockKey)
	}()

	now := time.Now()
	items, err := s.delayQueue.PopDue(ctx, now, s.cfg.BatchSize)
	if err != nil {
		if !errors.Is(err, redisx.ErrNotFound) {
			logx.WithContext(ctx).Errorf("delay scanner pop due failed: %v", err)
			metricsx.Inc("scheduler_scanner_pop_error_total")
			s.maybeFallback(ctx, now)
		}
		return
	}
	if len(items) == 0 {
		metricsx.Inc("scheduler_scanner_pop_empty_total")
		return
	}
	metricsx.Add("scheduler_scanner_pop_due_total", int64(len(items)))

	for _, instanceNo := range items {
		instance, err := s.instances.GetByInstanceNo(ctx, instanceNo)
		if err != nil {
			logx.WithContext(ctx).Errorf("delay scanner load instance=%s failed: %v", instanceNo, err)
			metricsx.Inc("scheduler_scanner_load_error_total")
			continue
		}
		if instance.Status != string(state.StatusPending) {
			continue
		}
		if err := s.readyQueue.Push(ctx, instanceNo); err != nil {
			logx.WithContext(ctx).Errorf("delay scanner push ready instance=%s failed: %v", instanceNo, err)
			metricsx.Inc("scheduler_scanner_push_ready_error_total")
			_ = s.delayQueue.Add(ctx, instanceNo, now.Add(s.cfg.RequeueDelay))
			metricsx.Inc("scheduler_scanner_requeue_total")
			continue
		}
		metricsx.Inc("scheduler_scanner_push_ready_total")
	}
}

func withDefaults(cfg Config) Config {
	if cfg.TickInterval <= 0 {
		cfg.TickInterval = time.Second
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100
	}
	if cfg.LockTTL <= 0 {
		cfg.LockTTL = 5 * time.Second
	}
	if cfg.RequeueDelay <= 0 {
		cfg.RequeueDelay = time.Second
	}
	if cfg.FallbackInterval <= 0 {
		cfg.FallbackInterval = 5 * time.Second
	}
	return cfg
}

func (s *DelayScanner) maybeFallback(ctx context.Context, now time.Time) {
	if s.readyQueue == nil || s.instances == nil {
		return
	}
	if s.cfg.FallbackInterval <= 0 {
		return
	}
	if !s.shouldFallback(now) {
		return
	}

	items, err := s.instances.ListDueInstances(ctx, now, int(s.cfg.BatchSize))
	if err != nil {
		logx.WithContext(ctx).Errorf("delay scanner fallback list failed: %v", err)
		metricsx.Inc("scheduler_scanner_fallback_error_total")
		return
	}
	if len(items) == 0 {
		metricsx.Inc("scheduler_scanner_fallback_empty_total")
		return
	}
	metricsx.Add("scheduler_scanner_fallback_total", int64(len(items)))
	for _, instanceNo := range items {
		if err := s.readyQueue.Push(ctx, instanceNo); err != nil {
			logx.WithContext(ctx).Errorf("delay scanner fallback push instance=%s failed: %v", instanceNo, err)
			metricsx.Inc("scheduler_scanner_fallback_push_error_total")
			continue
		}
		metricsx.Inc("scheduler_scanner_fallback_push_ready_total")
	}
}

func (s *DelayScanner) shouldFallback(now time.Time) bool {
	if s.lastFallback.IsZero() || now.Sub(s.lastFallback) >= s.cfg.FallbackInterval {
		s.lastFallback = now
		return true
	}
	return false
}
