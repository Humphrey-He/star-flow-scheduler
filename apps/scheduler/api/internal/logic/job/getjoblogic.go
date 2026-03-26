package job

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/assembler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/service"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	jobSvc *service.JobService
}

func NewGetJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobLogic {
	return &GetJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		jobSvc: service.NewJobService(svcCtx.Jobs),
	}
}

func (l *GetJobLogic) GetJob(req *types.GetJobRequest) (resp *types.GetJobResponse, err error) {
	job, err := l.jobSvc.GetByCode(l.ctx, req.JobCode)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errx.NotFound("job not found")
		}
		return nil, err
	}

	return &types.GetJobResponse{
		Job: assembler.MapJobDefinition(job),
	}, nil
}
