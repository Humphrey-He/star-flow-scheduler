package executorregistryservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

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
	// todo: add your logic here and delete this line

	return &schedulerv1_schedulev1.RegisterExecutorResponse{}, nil
}
