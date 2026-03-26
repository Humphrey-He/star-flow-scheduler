package executorregistryservicelogic

import (
	"context"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
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

	upsert := repo.ExecutorUpsert{
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
		Metadata:      nil,
	}

	_, err := l.svcCtx.Executors.Upsert(l.ctx, upsert)
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.RegisterExecutorResponse{
		ExecutorId:           0,
		HeartbeatIntervalSec: heartbeatIntervalSec,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
