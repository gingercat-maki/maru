package basic

import (
	"github.com/temporalio/maru/target/modules"
	"github.com/temporalio/samples-go/greetingslocal"
	"go.temporal.io/sdk/worker"
)

func RegisterModules(w worker.Worker) {
	modules.RegisterChildWorkflow(w)
	modules.RegisterSagaTransfer(w)
	// todo: move to module
	activities := &greetingslocal.Activities{Name: "MaruX-Temporal", Greeting: "MaruX-Hello"}
	w.RegisterActivity(activities)
}
