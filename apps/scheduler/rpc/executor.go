package main

import (
	"flag"
	"fmt"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/config"
	dispatchserviceServer "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/server/dispatchservice"
	executorregistryserviceServer "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/server/executorregistryservice"
	schedulerinternalserviceServer "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/server/schedulerinternalservice"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/executor.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		schedulev1.RegisterExecutorRegistryServiceServer(grpcServer, executorregistryserviceServer.NewExecutorRegistryServiceServer(ctx))
		schedulev1.RegisterDispatchServiceServer(grpcServer, dispatchserviceServer.NewDispatchServiceServer(ctx))
		schedulev1.RegisterSchedulerInternalServiceServer(grpcServer, schedulerinternalserviceServer.NewSchedulerInternalServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
