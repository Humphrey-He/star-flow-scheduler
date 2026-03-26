package executorregistryservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/registry"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

const heartbeatIntervalSec = 10

type RegisterExecutorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterExecutorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterExecutorLogic {
	return &RegisterExecutorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterExecutorLogic) RegisterExecutor(in *schedulerv1_schedulev1.RegisterExecutorRequest) (*schedulerv1_schedulev1.RegisterExecutorResponse, error) {
	id, err := l.svcCtx.RegistrySvc.Register(l.ctx, registry.RegisterRequest{
		ExecutorCode: in.ExecutorCode,
		Host:         in.Host,
		IP:           in.Ip,
		GrpcAddr:     in.GrpcAddr,
		HttpAddr:     in.HttpAddr,
		Tags:         in.Tags,
		Capacity:     int(in.Capacity),
		CurrentLoad:  0,
		Version:      in.Version,
		Status:       "online",
	})
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.RegisterExecutorResponse{
		ExecutorId:           id,
		HeartbeatIntervalSec: heartbeatIntervalSec,
	}, nil
}
