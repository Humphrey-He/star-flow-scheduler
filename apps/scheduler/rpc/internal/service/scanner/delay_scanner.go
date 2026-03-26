package scanner

import (
	"context"
	"errors"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	TickInterval time.Duration
	BatchSize    int64
	LockTTL      time.Duration
	RequeueDelay time.Duration
}

type DelayScanner struct {
	cfg         Config
	delayQueue  redisx.DelayQueue
	readyQueue  redisx.ReadyQueue
	locker      redisx.Locker
	instances   *repo.JobInstanceRepository
	lockKey     string
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
	ok, err := s.locker.TryLock(ctx, s.lockKey, s.cfg.LockTTL)
	if err != nil || !ok {
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
		}
		return
	}
	if len(items) == 0 {
		return
	}

	for _, instanceNo := range items {
		instance, err := s.instances.GetByInstanceNo(ctx, instanceNo)
		if err != nil {
			logx.WithContext(ctx).Errorf("delay scanner load instance=%s failed: %v", instanceNo, err)
			continue
		}
		if instance.Status != string(state.StatusPending) {
			continue
		}
		if err := s.readyQueue.Push(ctx, instanceNo); err != nil {
			logx.WithContext(ctx).Errorf("delay scanner push ready instance=%s failed: %v", instanceNo, err)
			_ = s.delayQueue.Add(ctx, instanceNo, now.Add(s.cfg.RequeueDelay))
		}
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
	return cfg
}
