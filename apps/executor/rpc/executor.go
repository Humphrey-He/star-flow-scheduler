package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/config"
	dispatchserviceServer "github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/server/dispatchservice"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/svc"
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
	if err := svc.ValidateConfig(c); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		schedulerv1_schedulev1.RegisterDispatchServiceServer(grpcServer, dispatchserviceServer.NewDispatchServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	bgCtx, cancel := context.WithCancel(context.Background())
	ctx.Start(bgCtx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("executor shutdown signal received")
		s.Stop()
		ctx.Drain()
		timeout := time.Duration(c.Runtime.ShutdownTimeoutSec) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
		defer shutdownCancel()
		_ = ctx.Shutdown(shutdownCtx)
		cancel()
	}()
	s.Start()
}
