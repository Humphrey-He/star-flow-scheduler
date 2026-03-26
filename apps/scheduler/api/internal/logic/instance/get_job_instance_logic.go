package instance

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

type GetJobInstanceLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	instSvc *service.InstanceService
}

func NewGetJobInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJobInstanceLogic {
	return &GetJobInstanceLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		instSvc: service.NewInstanceService(svcCtx.Instances),
	}
}

func (l *GetJobInstanceLogic) GetJobInstance(req *types.GetJobInstanceRequest) (resp *types.GetJobInstanceResponse, err error) {
	item, err := l.instSvc.GetByInstanceNo(l.ctx, req.InstanceNo)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errx.NotFound("instance not found")
		}
		return nil, err
	}

	return &types.GetJobInstanceResponse{
		Instance: assembler.MapJobInstance(item),
	}, nil
}
