package basic

import (
	"github.com/temporalio/maru/target/modules"
)

// WorkflowRequest is used for starting workflow
// the parameters control the pattern of the workflow
// components that need no enable: queries, updates
type WorkflowRequest struct {
	SequenceCount                int    `json:"sequenceCount"`
	ParallelCount                int    `json:"parallelCount"`
	ActivityDurationMilliseconds int    `json:"activityDurationMilliseconds"`
	Payload                      string `json:"payload"`
	ResultPayload                string `json:"resultPayload"`

	// injected modules
	EnabledAwaitSignal   bool `json:"enabledAwaitSignal"`
	EnabledLocalActivity bool `json:"enabledLocalActivity"`

	modules.ChildWorkflowRequest
	modules.SagaTransferRequest
}
