package reporter

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/client/dispatchservice"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type Reporter struct {
	client        dispatchservice.DispatchService
	queue         chan *model.TaskResult
	retryTimes    int
	retryInterval time.Duration
}

func NewReporter(client dispatchservice.DispatchService, queueSize int, retryTimes int, retryInterval time.Duration) *Reporter {
	return &Reporter{
		client:        client,
		queue:         make(chan *model.TaskResult, queueSize),
		retryTimes:    retryTimes,
		retryInterval: retryInterval,
	}
}

func (r *Reporter) Start(ctx context.Context) {
	go r.loop(ctx)
}

func (r *Reporter) Report(result *model.TaskResult) {
	select {
	case r.queue <- result:
	default:
		logx.WithContext(context.Background()).Errorf("reporter queue full instance=%s shard=%s", result.InstanceNo, result.ShardNo)
	}
}

func (r *Reporter) Pending() int {
	return len(r.queue)
}

func (r *Reporter) loop(ctx context.Context) {
	logger := logx.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case result := <-r.queue:
			if result == nil {
				continue
			}
			if err := r.sendResult(ctx, result); err != nil {
				logger.Errorf("report result failed instance=%s err=%v", result.InstanceNo, err)
			}
		}
	}
}

func (r *Reporter) sendResult(ctx context.Context, result *model.TaskResult) error {
	metricsx.Inc("executor_report_result_total")
	var lastErr error
	for i := 0; i <= r.retryTimes; i++ {
		_, err := r.client.ReportResult(ctx, &schedulev1.ReportResultRequest{
			InstanceNo:    result.InstanceNo,
			ShardNo:       result.ShardNo,
			Status:        result.Status,
			StartTime:     result.StartTime.UnixMilli(),
			FinishTime:    result.FinishTime.UnixMilli(),
			ErrorCode:     result.ErrorCode,
			ErrorMessage:  result.ErrorMessage,
			ResultSummary: result.ResultSummary,
		})
		if err == nil {
			return nil
		}
		lastErr = err
		if i < r.retryTimes {
			time.Sleep(r.retryInterval)
		}
	}
	metricsx.Inc("executor_report_result_fail_total")
	return lastErr
}

func (r *Reporter) Stop(ctx context.Context) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if r.Pending() == 0 {
				return nil
			}
		}
	}
}
