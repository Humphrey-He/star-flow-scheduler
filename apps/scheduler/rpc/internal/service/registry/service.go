package registry

import (
	"context"
	"strings"
	"time"

	rpcrepo "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
)

type Service struct {
	executors      *rpcrepo.ExecutorRepository
	heartbeatCache redisx.HeartbeatCache
	heartbeatTTL   time.Duration
}

func NewService(executors *rpcrepo.ExecutorRepository, heartbeatCache redisx.HeartbeatCache, heartbeatTTL time.Duration) *Service {
	return &Service{
		executors:      executors,
		heartbeatCache: heartbeatCache,
		heartbeatTTL:   heartbeatTTL,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (int64, error) {
	tags := ""
	if len(req.Tags) > 0 {
		tags = strings.Join(req.Tags, ",")
	}

	upsert := pkgrepo.ExecutorUpsert{
		ExecutorCode:  req.ExecutorCode,
		Host:          req.Host,
		IP:            req.IP,
		GrpcAddr:      req.GrpcAddr,
		HttpAddr:      strPtr(req.HttpAddr),
		Tags:          strPtr(tags),
		Capacity:      req.Capacity,
		CurrentLoad:   req.CurrentLoad,
		Version:       strPtr(req.Version),
		Status:        req.Status,
		LastHeartbeat: time.Now(),
	}

	exec, err := s.executors.Upsert(ctx, upsert)
	if err != nil {
		return 0, err
	}
	s.setHeartbeatCache(ctx, req.ExecutorCode, req.CurrentLoad, 0)
	s.updateOnlineMetrics(ctx)

	return int64(exec.ID), nil
}

func (s *Service) Heartbeat(ctx context.Context, executorCode string, currentLoad int, runningJobs int) error {
	if err := s.executors.UpdateHeartbeat(ctx, executorCode, currentLoad); err != nil {
		return err
	}
	s.setHeartbeatCache(ctx, executorCode, currentLoad, runningJobs)
	s.updateOnlineMetrics(ctx)
	return nil
}

type RegisterRequest struct {
	ExecutorCode string
	Host         string
	IP           string
	GrpcAddr     string
	HttpAddr     string
	Tags         []string
	Capacity     int
	CurrentLoad  int
	Version      string
	Status       string
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *Service) setHeartbeatCache(ctx context.Context, executorCode string, currentLoad int, runningJobs int) {
	if s.heartbeatCache == nil || s.heartbeatTTL <= 0 {
		return
	}
	_ = s.heartbeatCache.Set(ctx, executorCode, redisx.ExecutorHeartbeat{
		ExecutorCode: executorCode,
		CurrentLoad:  int64(currentLoad),
		RunningJobs:  int64(runningJobs),
		UpdatedAtMs:  time.Now().UnixMilli(),
	}, s.heartbeatTTL)
}

func (s *Service) updateOnlineMetrics(ctx context.Context) {
	if s.executors == nil {
		return
	}
	count, err := s.executors.CountOnline(ctx)
	if err != nil {
		return
	}
	metricsx.Set("scheduler_executor_online_total", int64(count))
}
