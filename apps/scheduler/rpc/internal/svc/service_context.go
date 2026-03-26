package svc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/config"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/dispatch"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/instance"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/registry"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/scanner"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/db"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config       config.Config
	DB           *db.DB
	Ent          *ent.Client
	ExecutorRepo *repo.ExecutorRepository
	InstanceRepo *repo.JobInstanceRepository
	JobRepo      *repo.JobRepository
	RegistrySvc  *registry.Service
	DispatchSvc  *dispatch.Service
	InstanceSvc  *instance.Service
	Redis        *redis.Client
	DelayScanner *scanner.DelayScanner
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.Open(c.PostgresDSN)
	if err != nil {
		panic(fmt.Sprintf("open postgres failed: %v", err))
	}

	executorRepo := repo.NewExecutorRepository(pkgrepo.NewExecutorRepository(database.Client))
	instanceRepo := repo.NewJobInstanceRepository(database.Client)
	jobRepo := repo.NewJobRepository(pkgrepo.NewJobRepository(database.Client))
	redisClient, err := redisx.NewRedis(c.Redis)
	if err != nil {
		panic(fmt.Sprintf("open redis failed: %v", err))
	}

	var delayQueue redisx.DelayQueue
	var readyQueue redisx.ReadyQueue
	var locker redisx.Locker
	if redisClient != nil {
		delayQueue = redisx.NewDelayQueue(redisClient)
		readyQueue = redisx.NewReadyQueue(redisClient)
		locker = redisx.NewLocker(redisClient)
	}

	dispatchSvc := dispatch.NewService(jobRepo, instanceRepo, executorRepo, nil)

	return &ServiceContext{
		Config:       c,
		DB:           database,
		Ent:          database.Client,
		ExecutorRepo: executorRepo,
		InstanceRepo: instanceRepo,
		JobRepo:      jobRepo,
		RegistrySvc:  registry.NewService(executorRepo),
		DispatchSvc:  dispatchSvc,
		InstanceSvc:  instance.NewService(instanceRepo),
		Redis:        redisClient,
		DelayScanner: scanner.NewDelayScanner(scanner.Config{
			TickInterval: time.Duration(c.Scanner.TickIntervalMs) * time.Millisecond,
			BatchSize:    c.Scanner.BatchSize,
			LockTTL:      time.Duration(c.Scanner.LockTTLms) * time.Millisecond,
			RequeueDelay: time.Duration(c.Scanner.RequeueDelayMs) * time.Millisecond,
		}, delayQueue, readyQueue, locker, instanceRepo),
	}
}

func (s *ServiceContext) Start(ctx context.Context) {
	bg, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	if s.DelayScanner != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.DelayScanner.Start(bg)
		}()
	}
}

func (s *ServiceContext) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
