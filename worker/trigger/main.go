package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"

	"github.com/temporalio/maru/internal"
	"github.com/temporalio/maru/target"
	"github.com/temporalio/maru/target/basic"
)

var c client.Client
var logger *zap.Logger

func init() {
	// The client is a heavyweight object that should be created once per process.
	var err error
	c, err = client.Dial(client.Options{
		HostPort: "127.0.0.1:7233", //client.DefaultHostPort,
	})
	if err != nil {
		panic(err)
	}
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer c.Close()
	// request
	request := basic.WorkflowRequest{
		SequenceCount:                1,
		ParallelCount:                3,
		ActivityDurationMilliseconds: 100,
		Payload:                      "test-1",
		ResultPayload:                "test-2",
		EnabledAwaitSignal:           true,
	}
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("trigger-%v", uuid.New()),
		TaskQueue: target.WorkflowTaskQueue,
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, basic.Workflow, request)
	if err != nil {
		logger.Fatal("Unable to execute workflow", zap.Error(err))
	}
	logger.Info("start workflow", zap.String("workflowID", we.GetID()), zap.String("runID", we.GetRunID()))

	if request.EnabledAwaitSignal {
		internal.TriggerAwaitSignal(c, we)
	}
}
