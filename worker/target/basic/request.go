package basic

import (
	"github.com/temporalio/maru/target/modules"
)

// WorkflowRequest is used for starting workflow
// the parameters control the pattern of the workflow
type WorkflowRequest struct {
	SequenceCount                int    `json:"sequenceCount"`
	ParallelCount                int    `json:"parallelCount"`
	ActivityDurationMilliseconds int    `json:"activityDurationMilliseconds"`
	Payload                      string `json:"payload"`
	ResultPayload                string `json:"resultPayload"`

	// injected modules
	LocalActivitiesModule
	SignalAwaitsModule
	modules.ChildWorkflowRequest
	modules.SagaTransferRequest
}

// SimpleModules: with no payload/workflow/register change
type LocalActivitiesModule struct {
	EnabledLocalActivity bool `json:"enabledLocalActivity"`
}

type SignalAwaitsModule struct {
	EnabledSignal bool `json:"enabledSignal"`
}
