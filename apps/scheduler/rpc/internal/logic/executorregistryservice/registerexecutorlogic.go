package executorregistryservicelogic

import (
	"context"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/internal/models"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

const heartbeatIntervalSec = 10

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
	tags := ""
	if len(in.Tags) > 0 {
		tags = strings.Join(in.Tags, ",")
	}

	exec := &models.Executor{
		ExecutorCode:  in.ExecutorCode,
		Host:          in.Host,
		IP:            in.Ip,
		GrpcAddr:      in.GrpcAddr,
		HttpAddr:      strPtr(in.HttpAddr),
		Tags:          strPtr(tags),
		Capacity:      int(in.Capacity),
		CurrentLoad:   0,
		Version:       strPtr(in.Version),
		Status:        "online",
		LastHeartbeat: time.Now(),
		Metadata:      strPtr(in.MetadataJson),
	}

	id, err := l.svcCtx.Executors.Upsert(l.ctx, exec)
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.RegisterExecutorResponse{
		ExecutorId:           id,
		HeartbeatIntervalSec: heartbeatIntervalSec,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
