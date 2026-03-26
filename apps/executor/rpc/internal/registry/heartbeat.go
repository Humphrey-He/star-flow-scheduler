package registry

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/client/executorregistryservice"
)

type Heartbeat struct {
	client   executorregistryservice.ExecutorRegistryService
	interval time.Duration
}

func NewHeartbeat(client executorregistryservice.ExecutorRegistryService, interval time.Duration) *Heartbeat {
	return &Heartbeat{client: client, interval: interval}
}

func (h *Heartbeat) Start(ctx context.Context, reqBuilder func() *executorregistryservice.HeartbeatRequest) {
	ticker := time.NewTicker(h.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				req := reqBuilder()
				if req == nil {
					continue
				}
				_, _ = h.client.Heartbeat(ctx, req)
			}
		}
	}()
}
