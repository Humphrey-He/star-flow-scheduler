package scheduler

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListJobInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListJobInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJobInstancesLogic {
	return &ListJobInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListJobInstancesLogic) ListJobInstances(req *types.ListJobInstancesRequest) (resp *types.ListJobInstancesResponse, err error) {
	startTime, err := parseTime(req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(req.EndTime)
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

	items, total, err := l.svcCtx.Instances.List(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	respItems := make([]types.JobInstance, 0, len(items))
	for i := range items {
		item := items[i]
		respItems = append(respItems, mapJobInstance(&item))
	}

	return &types.ListJobInstancesResponse{
		Items:    respItems,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}
