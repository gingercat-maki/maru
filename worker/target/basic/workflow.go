// The MIT License
//
// Copyright (c) 2021 Temporal Technologies Inc.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.  //
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package basic

import (
	"encoding/json"
	"time"

	samples "github.com/temporalio/samples-go/await-signals"
	sampleschild "github.com/temporalio/samples-go/child-workflow"
	sampleslocal "github.com/temporalio/samples-go/greetingslocal"
	samplessaga "github.com/temporalio/samples-go/saga"

	"go.temporal.io/sdk/workflow"
)

const taskQueue = "temporal-basic-act"

type WorkflowRunningContext struct {
	RawRequestPayload WorkflowRequest `json:"rawRequestPayload"`
}

// TODO: how to query from the bench? traffic-generator/trigger; vary-patterned-workflows; reporters;
// Workflow implements a basic bench scenario to schedule activities in sequence.
func Workflow(ctx workflow.Context, request WorkflowRequest) (string, error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("basic workflow started", "activity task queue", taskQueue)
	wrc := WorkflowRunningContext{}

	err := workflow.SetQueryHandler(
		ctx, "queryContext", func() (string, error) {
			wrcInBytes, err := json.Marshal(wrc)
			return string(wrcInBytes[:]), err
		})
	if err != nil {
		logger.Error("setQueryHandler fails,", "err", err)
		return "", err
	}

	// Currently: 1) we just put the modules here in sequential to prove they can be run
	// 2) combinations of modules

	// poc:implement a signal module
	if request.EnabledSignal {
		// TODO: this doesn't need add activity workers, but others may need this
		err := samples.AwaitSignalsWorkflow(ctx)
		if err != nil {
			logger.Error("awaitSignals fails,", "err", err)
			return "", err
		}
	}

	// module: localacitivty
	if request.EnabledLocalActivity {
		// TODO: this has activities,
		// need to register activities
		// but local activities doesn't need a seperate queue
		// nor a seprate worker
		_, err := sampleslocal.GreetingSample(ctx)
		if err != nil {
			logger.Error("localActivity fails,", "err", err)
			return "", err
		}
	}

	// module:normal part (saga with transfers)
	// TODO: this has activities, should register
	// TODO: but we can stay with one queues? this kinds of becomes a handler of the bottleneck? all traffic of all types in one queue?
	// queue: is handled in starter & worker, so we manipulate queue in the starter and worker part
	if request.EnabledSagaWithTransfer {
		err = samplessaga.TransferMoney(ctx, request.TransferDetailsRequest)
		if err != nil {
			// TODO: this fails the workflow after rollback,
			// so cannot return err to stop the wf
			logger.Error("sagaWithTransfer fails with rollback,", "err", err)
		}
	}

	if request.EnabledChildWorkflow {
		_, err = sampleschild.SampleParentWorkflow(ctx)
		if err != nil {
			logger.Error("sampleParentWorkflow fails,", "err", err)
			return "", err
		}
	}

	// normal Activity Part
	// TODO: plans to disable this part of workflow
	ao := workflow.ActivityOptions{
		TaskQueue: taskQueue,
		StartToCloseTimeout: time.Duration(
			request.ActivityDurationMilliseconds)*time.Millisecond + 10*time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	parallelCount := 1
	if request.ParallelCount > 1 {
		parallelCount = request.ParallelCount
	}

	for i := 0; i < request.SequenceCount; i++ {
		req := basicActivityRequest{
			ActivityDelayMilliseconds: request.ActivityDurationMilliseconds,
			Payload:                   request.Payload,
			ResultPayload:             request.ResultPayload,
		}

		futures := make([]workflow.Future, parallelCount)
		for i := 0; i < parallelCount; i++ {
			futures[i] = workflow.ExecuteActivity(ctx, "basic-activity", req)
		}

		allResults := make([]string, parallelCount)
		for i := 0; i < parallelCount; i++ {
			var result string
			err := futures[i].Get(ctx, &result)
			if err != nil {
				return "", err
			}
			allResults[i] = result
		}

		logger.Info("activity returned result to the workflow", "value", allResults)
	}

	logger.Info("basic workflow completed")
	return request.ResultPayload, nil
}
