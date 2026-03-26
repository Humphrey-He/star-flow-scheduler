package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "trigger":
		runTrigger(os.Args[2:])
	default:
		printUsage()
	}
}

func runTrigger(args []string) {
	fs := flag.NewFlagSet("trigger", flag.ExitOnError)
	jobCode := fs.String("job_code", "", "job code")
	payload := fs.String("payload", "", "json payload")
	rpcAddr := fs.String("rpc_addr", "127.0.0.1:8080", "scheduler rpc address")
	fs.Parse(args)

	if *jobCode == "" {
		fmt.Println("job_code is required")
		os.Exit(2)
	}

	if *payload != "" && !json.Valid([]byte(*payload)) {
		fmt.Println("payload must be valid json")
		os.Exit(2)
	}

	conn, err := grpc.Dial(*rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("dial rpc failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := schedulerv1_schedulev1.NewSchedulerInternalServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &schedulerv1_schedulev1.CreateInstanceRequest{
		JobCode:     *jobCode,
		TriggerType: "manual",
	}
	if *payload != "" {
		req.Payload = &schedulerv1_schedulev1.JobPayload{
			Raw:         []byte(*payload),
			ContentType: "application/json",
		}
	}

	createResp, err := client.CreateInstance(ctx, req)
	if err != nil {
		fmt.Printf("create instance failed: %v\n", err)
		os.Exit(1)
	}

	dispatchResp, err := client.DispatchInstance(ctx, &schedulerv1_schedulev1.DispatchInstanceRequest{InstanceNo: createResp.InstanceNo})
	if err != nil {
		fmt.Printf("dispatch instance failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("instance_no=%s status=%s dispatched=%v executor=%s\n",
		createResp.InstanceNo,
		createResp.Status.String(),
		dispatchResp.Dispatched,
		dispatchResp.ExecutorCode,
	)
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  starflow trigger --job_code=demo_job --payload='{\"k\":\"v\"}' --rpc_addr=127.0.0.1:8080")
}
