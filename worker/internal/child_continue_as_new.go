package internal

import (
	cw "github.com/temporalio/samples-go/child-workflow"
	"go.temporal.io/sdk/worker"
)

func RegisterChildWorkflow(w worker.Worker) {
	w.RegisterWorkflow(cw.SampleChildWorkflow)
}
