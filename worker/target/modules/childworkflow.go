package modules

import (
	cw "github.com/temporalio/samples-go/child-workflow"
	"go.temporal.io/sdk/worker"
)

type ChildWorkflowRequest struct {
	EnabledChildWorkflow bool `json:"enabledChildWorkflow"`
}

func RegisterChildWorkflow(w worker.Worker) {
	w.RegisterWorkflow(cw.SampleChildWorkflow)
}
