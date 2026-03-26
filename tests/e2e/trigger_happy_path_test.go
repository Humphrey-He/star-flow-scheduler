package e2e

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/db"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobdefinition"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobinstance"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestE2ETriggerHappyPath(t *testing.T) {
	if os.Getenv("RUN_E2E") != "1" {
		t.Skip("set RUN_E2E=1 to run")
	}

	rpcAddr := os.Getenv("SCHEDULER_RPC_ADDR")
	if rpcAddr == "" {
		t.Fatal("SCHEDULER_RPC_ADDR is required")
	}
	postgres := os.Getenv("POSTGRES_DSN")
	if postgres == "" {
		t.Fatal("POSTGRES_DSN is required")
	}

	jobCode := os.Getenv("JOB_CODE")
	if jobCode == "" {
		jobCode = "demo_job"
	}

	dbConn, err := db.Open(postgres)
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	defer dbConn.Client.Close()

	if err := ensureJobDefinition(context.Background(), dbConn.Client, jobCode); err != nil {
		t.Fatalf("ensure job definition failed: %v", err)
	}

	conn, err := grpc.Dial(rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial scheduler rpc failed: %v", err)
	}
	defer conn.Close()

	client := schedulerv1_schedulev1.NewSchedulerInternalServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := map[string]string{"k": "v"}
	raw, _ := json.Marshal(payload)

	createResp, err := client.CreateInstance(ctx, &schedulerv1_schedulev1.CreateInstanceRequest{
		JobCode:     jobCode,
		TriggerType: "manual",
		Payload: &schedulerv1_schedulev1.JobPayload{
			Raw:         raw,
			ContentType: "application/json",
		},
	})
	if err != nil {
		t.Fatalf("create instance failed: %v", err)
	}

	_, err = client.DispatchInstance(ctx, &schedulerv1_schedulev1.DispatchInstanceRequest{InstanceNo: createResp.InstanceNo})
	if err != nil {
		t.Fatalf("dispatch instance failed: %v", err)
	}

	if err := waitInstanceSuccess(context.Background(), dbConn.Client, createResp.InstanceNo, 20*time.Second); err != nil {
		t.Fatalf("instance not success: %v", err)
	}
}

func ensureJobDefinition(ctx context.Context, client *ent.Client, jobCode string) error {
	exists, err := client.JobDefinition.Query().Where(jobdefinition.JobCodeEQ(jobCode)).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	jobRepo := repo.NewJobRepository(client)
	_, err = jobRepo.Create(ctx, repo.JobDefinitionCreate{
		JobCode:       jobCode,
		JobName:       "Demo Job",
		JobType:       "once",
		ExecuteMode:   "standalone",
		HandlerName:   "demo_print",
		TimeoutMs:     10000,
		RetryLimit:    0,
		RetryBackoff:  "1s",
		Priority:      5,
		ShardTotal:    1,
		RouteStrategy: "least_load",
		Status:        "enabled",
	})
	return err
}

func waitInstanceSuccess(ctx context.Context, client *ent.Client, instanceNo string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		instance, err := client.JobInstance.Query().Where(jobinstance.InstanceNoEQ(instanceNo)).Only(ctx)
		if err == nil {
			if instance.Status == "success" && instance.StartTime != nil && instance.FinishTime != nil {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return context.DeadlineExceeded
		}
		time.Sleep(500 * time.Millisecond)
	}
}
