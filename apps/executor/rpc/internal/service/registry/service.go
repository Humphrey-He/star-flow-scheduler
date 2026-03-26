package registry

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/registry"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/client/executorregistryservice"
)

type Service struct {
	registrar *registry.Registrar
	heartbeat *registry.Heartbeat
}

func NewService(registrar *registry.Registrar, heartbeat *registry.Heartbeat) *Service {
	return &Service{
		registrar: registrar,
		heartbeat: heartbeat,
	}
}

func (s *Service) Register(ctx context.Context, req *executorregistryservice.RegisterExecutorRequest) (*executorregistryservice.RegisterExecutorResponse, error) {
	return s.registrar.Register(ctx, req)
}

func (s *Service) StartHeartbeat(ctx context.Context, interval time.Duration, reqBuilder func() *executorregistryservice.HeartbeatRequest) {
	s.heartbeat.Start(ctx, reqBuilder)
}
