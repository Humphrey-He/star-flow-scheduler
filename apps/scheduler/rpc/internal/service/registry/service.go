package registry

import (
	"context"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type Service struct {
	executors *repo.ExecutorRepository
}

func NewService(executors *repo.ExecutorRepository) *Service {
	return &Service{executors: executors}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (int64, error) {
	tags := ""
	if len(req.Tags) > 0 {
		tags = strings.Join(req.Tags, ",")
	}

	upsert := repo.ExecutorUpsert{
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

	return exec.ID, nil
}

func (s *Service) Heartbeat(ctx context.Context, executorCode string, currentLoad int) error {
	return s.executors.UpdateHeartbeat(ctx, executorCode, currentLoad)
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
