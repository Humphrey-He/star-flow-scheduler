package instance

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/assembler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/service"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListJobInstancesLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	instSvc *service.InstanceService
}

func NewListJobInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJobInstancesLogic {
	return &ListJobInstancesLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		instSvc: service.NewInstanceService(svcCtx.Instances),
	}
}

func (l *ListJobInstancesLogic) ListJobInstances(req *types.ListJobInstancesRequest) (resp *types.ListJobInstancesResponse, err error) {
	startTime, err := assembler.ParseTime(req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := assembler.ParseTime(req.EndTime)
	if err != nil {
		return nil, err
	}

	filter := repo.JobInstanceFilter{
		JobCode:  req.JobCode,
		Status:   req.Status,
		StartAt:  startTime,
		EndAt:    endTime,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	items, total, err := l.instSvc.List(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	respItems := make([]types.JobInstance, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, assembler.MapJobInstance(item))
	}

	return &types.ListJobInstancesResponse{
		Items:    respItems,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}
