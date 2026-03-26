package job

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/service"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	jobSvc *service.JobService
}

func NewCreateJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobLogic {
	return &CreateJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		jobSvc: service.NewJobService(svcCtx.Jobs),
	}
}

func (l *CreateJobLogic) CreateJob(req *types.CreateJobRequest) (resp *types.CreateJobResponse, err error) {
	created, err := l.jobSvc.Create(l.ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.CreateJobResponse{
		JobId:   int64(created.ID),
		JobCode: created.JobCode,
	}, nil
}
