package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/config"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/handler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/handler/builtin"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/receiver"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/registry"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/reporter"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/runtime"
	dispatchservice "github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/service/dispatch"
	registryservice "github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/service/registry"
	dispatchclient "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/client/dispatchservice"
	registryclient "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/client/executorregistryservice"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	SchedulerRpc zrpc.Client

	HandlerRegistry *handler.Registry
	Runtime         *runtime.Runtime
	Receiver        *receiver.Receiver
	Reporter        *reporter.Reporter
	RegistrySvc     *registryservice.Service
	DispatchSvc     *dispatchservice.Service
	Redis           *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	schedulerClient := zrpc.MustNewClient(c.SchedulerRpc)
	dispatchSvcClient := dispatchclient.NewDispatchService(schedulerClient)
	registrySvcClient := registryclient.NewExecutorRegistryService(schedulerClient)

	redisClient, err := redisx.NewRedis(c.Redis)
	if err != nil {
		panic(fmt.Sprintf("open redis failed: %v", err))
	}

	handlerRegistry := handler.NewRegistry()
	_ = handlerRegistry.Register(&builtin.NoopHandler{})
	_ = handlerRegistry.Register(&builtin.DemoHandler{})

	reporterSvc := reporter.NewReporter(
		dispatchSvcClient,
		c.Reporter.QueueSize,
		c.Reporter.RetryTimes,
		time.Duration(c.Reporter.RetryIntervalMs)*time.Millisecond,
	)

	rt := runtime.NewRuntime(runtime.Config{
		WorkerCount:        c.Runtime.WorkerCount,
		QueueSize:          c.Runtime.QueueSize,
		DefaultTimeoutMs:   c.Runtime.DefaultTimeoutMs,
		ShutdownTimeoutSec: c.Runtime.ShutdownTimeoutSec,
	}, handlerRegistry, reporterSvc)

	recv := receiver.NewReceiver(rt)
	dispatchSvc := dispatchservice.NewService(recv)

	registrar := registry.NewRegistrar(registrySvcClient)
	heartbeat := registry.NewHeartbeat(registrySvcClient, time.Duration(c.Registry.HeartbeatIntervalSec)*time.Second)
	regSvc := registryservice.NewService(registrar, heartbeat)

	return &ServiceContext{
		Config:          c,
		SchedulerRpc:    schedulerClient,
		HandlerRegistry: handlerRegistry,
		Runtime:         rt,
		Receiver:        recv,
		Reporter:        reporterSvc,
		RegistrySvc:     regSvc,
		DispatchSvc:     dispatchSvc,
		Redis:           redisClient,
	}
}

func (s *ServiceContext) Start(ctx context.Context) {
	logx.AddGlobalFields(logx.Field("executor_code", s.Config.Executor.ExecutorCode))
	s.Runtime.Start(ctx)
	s.Reporter.Start(ctx)

	req := s.buildRegisterRequest()
	resp, err := s.RegistrySvc.Register(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("register executor failed: %v", err)
	} else {
		logx.WithContext(ctx).Infof("register executor ok executor_id=%d", resp.ExecutorId)
	}

	s.RegistrySvc.StartHeartbeat(ctx, time.Duration(s.Config.Registry.HeartbeatIntervalSec)*time.Second, s.buildHeartbeatRequest)
}

func (s *ServiceContext) Drain() {
	s.Runtime.Drain()
}

func (s *ServiceContext) Shutdown(ctx context.Context) error {
	if err := s.Runtime.Wait(ctx); err != nil {
		return err
	}
	return s.Reporter.Stop(ctx)
}

func (s *ServiceContext) buildRegisterRequest() *registryclient.RegisterExecutorRequest {
	cfg := s.Config.Executor
	return &registryclient.RegisterExecutorRequest{
		ExecutorCode: cfg.ExecutorCode,
		Host:         cfg.Host,
		Ip:           cfg.IP,
		GrpcAddr:     cfg.GrpcAddr,
		HttpAddr:     cfg.HttpAddr,
		Tags:         cfg.Tags,
		Capacity:     cfg.Capacity,
		Version:      cfg.Version,
	}
}

func (s *ServiceContext) buildHeartbeatRequest() *registryclient.HeartbeatRequest {
	cfg := s.Config.Executor
	running := s.Runtime.RunningJobs()
	queueSize := s.Runtime.QueueSize()
	currentLoad := int32(running + int64(queueSize))
	return &registryclient.HeartbeatRequest{
		ExecutorCode: cfg.ExecutorCode,
		CurrentLoad:  currentLoad,
		RunningJobs:  int32(running),
		Timestamp:    time.Now().UnixMilli(),
	}
}

func ValidateConfig(c config.Config) error {
	if c.Executor.ExecutorCode == "" {
		return fmt.Errorf("executor executor_code is empty")
	}
	if c.Runtime.WorkerCount <= 0 {
		return fmt.Errorf("runtime worker_count must be positive")
	}
	if c.Runtime.QueueSize <= 0 {
		return fmt.Errorf("runtime queue_size must be positive")
	}
	if c.Runtime.ShutdownTimeoutSec <= 0 {
		return fmt.Errorf("runtime shutdown_timeout_sec must be positive")
	}
	if c.Registry.HeartbeatIntervalSec <= 0 {
		return fmt.Errorf("registry heartbeat_interval_sec must be positive")
	}
	return nil
}
