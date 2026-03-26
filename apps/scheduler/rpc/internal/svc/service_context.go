package svc

import (
	"fmt"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/config"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/dispatch"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/instance"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/registry"
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

	return &ServiceContext{
		Config:       c,
		DB:           database,
		Ent:          database.Client,
		ExecutorRepo: executorRepo,
		InstanceRepo: instanceRepo,
		JobRepo:      jobRepo,
		RegistrySvc:  registry.NewService(executorRepo),
		DispatchSvc:  dispatch.NewService(jobRepo, instanceRepo, executorRepo, nil),
		InstanceSvc:  instance.NewService(instanceRepo),
		Redis:        redisClient,
	}
}
