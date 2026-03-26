package scheduler

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJobInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobInstanceLogic {
	return &GetJobInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobInstanceLogic) GetJobInstance(req *types.GetJobInstanceRequest) (resp *types.GetJobInstanceResponse, err error) {
	item, err := l.svcCtx.Instances.GetByInstanceNo(l.ctx, req.InstanceNo)
	if err != nil {
		return nil, err
	}

	return &types.GetJobInstanceResponse{
		Instance: mapJobInstance(item),
	}, nil
}
