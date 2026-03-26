package scheduler

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobLogic {
	return &GetJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJobLogic) GetJob(req *types.GetJobRequest) (resp *types.GetJobResponse, err error) {
	job, err := l.svcCtx.Jobs.GetByCode(l.ctx, req.JobCode)
	if err != nil {
		return nil, err
	}

	return &types.GetJobResponse{
		Job: mapJobDefinition(job),
	}, nil
}
