package registry

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/client/executorregistryservice"
)

type Registrar struct {
	client executorregistryservice.ExecutorRegistryService
}

func NewRegistrar(client executorregistryservice.ExecutorRegistryService) *Registrar {
	return &Registrar{client: client}
}

func (r *Registrar) Register(ctx context.Context, req *executorregistryservice.RegisterExecutorRequest) (*executorregistryservice.RegisterExecutorResponse, error) {
	return r.client.RegisterExecutor(ctx, req)
}
