package internal

import (
	"github.com/temporalio/samples-go/saga"
)

// Workflow Request and Patterns
type WorkflowRequest struct {
	// workflow common parameters
	Payload       string `json:"payload"`       // optional
	ResultPayload string `json:"resultPayload"` // optional

	// modules of workflow
	EnabledAwaitSignal   bool `json:"enabledAwaitSignal"`
	EnabledLocalActivity bool `json:"enabledLocalActivity"`
	EnabledChildWorkflow bool `json:"enabledChildWorkflow"`

	EnabledSagaWithTransfer bool                 `json:"enabledSagaWithTransfer"`
	TransferDetailsRequest  saga.TransferDetails `json:"transferDetails"`

	EnabledNormalActivity bool `json:"enabledNormalActivities"`
	NormalActivityRequest NormalActivityRequest
}

type NormalActivityRequest struct {
	ActivityPayload              string `json:"activityPayload"`
	ActivityResultPayload        string `json:"activityResultPayload"`
	SequenceCount                int    `json:"sequenceCount"`
	ParallelCount                int    `json:"parallelCount"`
	ActivityDurationMilliseconds int    `json:"activityDurationMilliseconds"`
}

type WorkflowRunningContext struct {
	RawRequestPayload WorkflowRequest `json:"rawRequestPayload"`
	Counter           int64           `json:"counter"`
}

// Plan of Test
type PlanRequest struct {
	// default to use the defaultPlan of comprehensive workflow types
	PlayType PlanType `json:"PlanType" yaml:"PlanType"`

	// default test plan: test all features and patterns
	DefaultPlan TrafficParameters `json:"defaultPlan" yaml:"defaultPlan"`

	// one plan for one workflow type and they will start in concurrency
	CustomizedPlan []PlanForWorkflow `json:"loadPlan" yaml:"loadPlan"` //percentages should sum to 100
}

type PlanType int

const (
	DefaultPlayType PlanType = iota
	CustomizedPlanType
)

type TrafficParameters struct {
	DurationInMinutes        int `json:"durationInMinutes" yaml:"durationInMinutes"`               // duration of test plan
	StartWorkflowConcurrency int `json:"startWorkflowConcurrency" yaml:"startWorkflowConcurrency"` // qps
	WaitTimeInMilliSeconds   int `json:"waitTimeInMilliSeconds" yaml:"waitTimeInMilliSeconds"`
	JitterInMilliSeconds     int `json:"jitterInMilliSeconds" yaml:"jitterInMilliSeconds"`
}

type PlanForWorkflow struct {
	WorkflowType
	LoadPlan
	SpikePlan
}

type WorkflowType struct {
	ID              string          `json:"ID" yaml:"ID"`
	WorkflowType    string          `json:"workflowType" yaml:"workflowType"`
	WorkflowRequest WorkflowRequest `json:"workflowRequest" yaml:"workflowRequest"`
}

type LoadPlan struct {
	TrafficParameters
	WeightPercent int `json:"weightPercent" yaml:"weightPercent"`
}

type SpikePlan struct {
	DurationInMinutes        int `json:"durationInMinutes" yaml:"durationInMinutes"`               // duration of test plan
	SpikeConcurrentyIncrease int `json:"startWorkflowConcurrency" yaml:"startWorkflowConcurrency"` // qps
	TotalCountLimit          int `json:"totalCountLimit" yaml:"totalCountLimit"`
	InternalInMilliSeconds   int `json:"intervalInMilliSeconds" yaml:"intervalMilliSeconds"`
	JitterInMilliSeconds     int `json:"jitterInMilliSeconds" yaml:"jitterInMilliSeconds"`
}
