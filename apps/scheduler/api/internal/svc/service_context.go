package svc

import (
	"fmt"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/config"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/db"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config    config.Config
	DB        *db.DB
	Ent       *ent.Client
	Jobs      *repo.JobRepository
	Instances *repo.JobInstanceRepository
	Redis     *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.Open(c.PostgresDSN)
	if err != nil {
		panic(fmt.Sprintf("open postgres failed: %v", err))
	}
	redisClient, err := redisx.NewRedis(c.Redis)
	if err != nil {
		panic(fmt.Sprintf("open redis failed: %v", err))
	}

	return &ServiceContext{
		Config:    c,
		DB:        database,
		Ent:       database.Client,
		Jobs:      repo.NewJobRepository(database.Client),
		Instances: repo.NewJobInstanceRepository(database.Client),
		Redis:     redisClient,
	}
}
