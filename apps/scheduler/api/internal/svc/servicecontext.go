package svc

import (
	"database/sql"
	"fmt"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/config"
	"github.com/Humphrey-He/star-flow-scheduler/internal/db"
	"github.com/Humphrey-He/star-flow-scheduler/internal/repo"
)

type ServiceContext struct {
	Config    config.Config
	DB        *sql.DB
	Jobs      *repo.JobRepository
	Instances *repo.JobInstanceRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.Open(c.MySQLDSN)
	if err != nil {
		panic(fmt.Sprintf("open mysql failed: %v", err))
	}

	return &ServiceContext{
		Config:    c,
		DB:        database,
		Jobs:      repo.NewJobRepository(database),
		Instances: repo.NewJobInstanceRepository(database),
	}
}
