package target

import (
	"github.com/temporalio/maru/internal"
	"github.com/temporalio/samples-go/greetingslocal"
	"go.temporal.io/sdk/worker"
)

func RegisterModules(w worker.Worker) {
	internal.RegisterChildWorkflow(w)
	internal.RegisterSagaTransfer(w)
	// todo: move to module
	activities := &greetingslocal.Activities{Name: "MaruX-Temporal", Greeting: "MaruX-Hello"}
	w.RegisterActivity(activities)
}
