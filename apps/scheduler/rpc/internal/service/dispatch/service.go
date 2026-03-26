package dispatch

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	rpcrepo "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/route"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"
	"github.com/robfig/cron/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	jobs      *rpcrepo.JobRepository
	instances *rpcrepo.JobInstanceRepository
	executors *rpcrepo.ExecutorRepository
	strategy  route.Strategy
	cache     redisx.HeartbeatCache
	delayQueue redisx.DelayQueue
}

func NewService(jobs *rpcrepo.JobRepository, instances *rpcrepo.JobInstanceRepository, executors *rpcrepo.ExecutorRepository, strategy route.Strategy, cache redisx.HeartbeatCache, delayQueue redisx.DelayQueue) *Service {
	if strategy == nil {
		strategy = route.NewStrategy("least_load")
	}
	return &Service{
		jobs:      jobs,
		instances: instances,
		executors: executors,
		strategy:  strategy,
		cache:     cache,
		delayQueue: delayQueue,
	}
}

func (s *Service) CreateInstance(ctx context.Context, jobCode string, triggerType string, payload *string) (*ent.JobInstance, error) {
	job, err := s.jobs.GetByCode(ctx, jobCode)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	scheduledAt, shouldDelay, err := nextScheduleTime(job, now)
	if err != nil {
		return nil, err
	}
	instanceNo := newInstanceNo()
	req := pkgrepo.JobInstanceCreate{
		InstanceNo:    instanceNo,
		JobID:         int64(job.ID),
		TriggerType:   triggerType,
		TriggerTime:   now,
		ScheduledTime: scheduledAt,
		Status:        string(state.StatusPending),
		Payload:       payload,
		ShardTotal:    job.ShardTotal,
	}

	instance, err := s.instances.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	if shouldDelay {
		if s.delayQueue == nil {
			return nil, fmt.Errorf("delay queue not configured")
		}
		if err := s.delayQueue.Add(ctx, instanceNo, scheduledAt); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (s *Service) DispatchInstance(ctx context.Context, instanceNo string) (*ent.Executor, error) {
	instance, err := s.instances.GetByInstanceNo(ctx, instanceNo)
	if err != nil {
		return nil, err
	}

	job, err := s.jobs.GetByID(ctx, instance.JobID)
	if err != nil {
		return nil, err
	}

	execs, err := s.executors.ListOnline(ctx, 100)
	if err != nil {
		return nil, err
	}
	if len(execs) == 0 {
		return nil, fmt.Errorf("no online executor")
	}

	nodes := toExecutorNodes(ctx, execs, s.cache)
	jobSnap := route.JobSnapshot{
		JobCode:     job.JobCode,
		RouteKey:    ptrString(instance.Payload),
		ExecutorTag: ptrString(job.ExecutorTag),
	}
	strategy := s.strategy
	if job.RouteStrategy != "" {
		strategy = route.NewStrategy(job.RouteStrategy)
	}
	selected, err := strategy.Select(ctx, jobSnap, nodes)
	if err != nil {
		return nil, err
	}
	exec := findExecutor(execs, selected)
	if exec == nil {
		return nil, fmt.Errorf("selected executor not found")
	}

	if err := s.dispatchToExecutor(ctx, exec, instance, job); err != nil {
		return nil, err
	}

	_, err = s.instances.UpdateDispatchInfoIfStatus(ctx, instanceNo, string(state.StatusPending), string(state.StatusDispatched), int64(exec.ID), time.Now())
	if err != nil {
		return nil, err
	}

	return exec, nil
}

func (s *Service) dispatchToExecutor(ctx context.Context, exec *ent.Executor, instance *ent.JobInstance, job *ent.JobDefinition) error {
	conn, err := grpc.DialContext(ctx, exec.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := schedulev1.NewDispatchServiceClient(conn)

	payload := &schedulev1.JobPayload{Raw: []byte(ptrString(instance.Payload))}
	_, err = client.DispatchJob(ctx, &schedulev1.DispatchJobRequest{
		InstanceNo:  instance.InstanceNo,
		JobCode:     job.JobCode,
		HandlerName: job.HandlerName,
		ShardNo:     "",
		TimeoutMs:   int32(job.TimeoutMs),
		Payload:     payload,
		TraceId:     "",
	})
	return err
}

func ptrString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func toExecutorNodes(ctx context.Context, execs []*ent.Executor, cache redisx.HeartbeatCache) []route.ExecutorNode {
	nodes := make([]route.ExecutorNode, 0, len(execs))
	for _, exec := range execs {
		currentLoad := exec.CurrentLoad
		if cache != nil {
			if hb, err := cache.Get(ctx, exec.ExecutorCode); err == nil {
				currentLoad = int32(hb.CurrentLoad)
				metricsx.Inc("route_cache_hit_total")
			} else if !errors.Is(err, redisx.ErrNotFound) {
				// ignore cache errors to keep dispatch safe
				metricsx.Inc("route_cache_error_total")
			} else {
				metricsx.Inc("route_cache_miss_total")
			}
		} else {
			metricsx.Inc("route_cache_disabled_total")
		}
		nodes = append(nodes, route.ExecutorNode{
			ID:           int64(exec.ID),
			ExecutorCode: exec.ExecutorCode,
			Tags:         splitTags(ptrString(exec.Tags)),
			CurrentLoad:  currentLoad,
		})
	}
	return nodes
}

func splitTags(tags string) []string {
	if tags == "" {
		return []string{}
	}
	parts := strings.Split(tags, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func findExecutor(execs []*ent.Executor, selected *route.ExecutorNode) *ent.Executor {
	if selected == nil {
		return nil
	}
	for _, exec := range execs {
		if int64(exec.ID) == selected.ID {
			return exec
		}
	}
	return nil
}

func newInstanceNo() string {
	return fmt.Sprintf("JI-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func nextScheduleTime(job *ent.JobDefinition, now time.Time) (time.Time, bool, error) {
	jobType := strings.ToLower(job.JobType)
	switch jobType {
	case "delay":
		if job.DelayMs <= 0 {
			return now, false, nil
		}
		return now.Add(time.Duration(job.DelayMs) * time.Millisecond), true, nil
	case "cron":
		if job.ScheduleExpr == nil || strings.TrimSpace(*job.ScheduleExpr) == "" {
			return now, false, nil
		}
		schedule, err := cron.ParseStandard(*job.ScheduleExpr)
		if err != nil {
			return now, false, err
		}
		next := schedule.Next(now)
		if next.IsZero() {
			return now, false, nil
		}
		return next, true, nil
	default:
		return now, false, nil
	}
}
