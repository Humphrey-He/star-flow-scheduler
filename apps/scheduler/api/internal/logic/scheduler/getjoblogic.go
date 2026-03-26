package scheduler

import (
	"context"
	"database/sql"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"

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
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &types.GetJobResponse{
		Job: mapJobDefinition(job),
	}, nil
}
