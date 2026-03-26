package executorregistryservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartbeatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatLogic {
	return &HeartbeatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HeartbeatLogic) Heartbeat(in *schedulerv1_schedulev1.HeartbeatRequest) (*schedulerv1_schedulev1.HeartbeatResponse, error) {
	// todo: add your logic here and delete this line

	return &schedulerv1_schedulev1.HeartbeatResponse{}, nil
}
