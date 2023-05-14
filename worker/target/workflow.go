// The MIT License
//
// Copyright (c) 2021 Temporal Technologies Inc.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy // of this software and associated documentation files (the "Software"), to deal
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

package target

import (
	"encoding/json"
	"fmt"

	"github.com/temporalio/maru/internal"
	"github.com/temporalio/maru/target/basic"
	sampleschild "github.com/temporalio/samples-go/child-workflow"
	sampleslocal "github.com/temporalio/samples-go/greetingslocal"
	samplessaga "github.com/temporalio/samples-go/saga"

	"github.com/temporalio/maru"
	"go.temporal.io/sdk/workflow"
)

// TODO: how to query from the bench? traffic-generator/trigger; vary-patterned-workflows; reporters;
// Workflow implements a basic bench scenario to schedule activities in sequence.
func Workflow(ctx workflow.Context, request internal.WorkflowRequest) (string, error) {

	logger := workflow.GetLogger(ctx)
	wrc := internal.WorkflowRunningContext{}

	// innate feature: Query
	err := workflow.SetQueryHandler(
		ctx, maru.QueryHandlerName, func() (string, error) {
			wrcInBytes, err := json.Marshal(wrc)
			return string(wrcInBytes[:]), err
		})
	if err != nil {
		logger.Error("setQueryHandler-queryContext fails,", "err", err)
		return "", err
	}

	// innate feature: Update
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		maru.UpdateHandlerName,
		func(ctx workflow.Context, i int64) (int64, error) {
			tmp := wrc.Counter
			wrc.Counter += i
			logger.Info("counter updated", "addend", i, "new-value", wrc.Counter)
			return tmp, nil
		},
		workflow.UpdateHandlerOptions{Validator: nonNegative},
	); err != nil {
		return "", err
	}

	// module: await on signal, this module requires signals to proceed
	if request.EnabledAwaitSignal {
		if err != nil {
			logger.Error("awaitSignals fails,", "err", err)
			return "", err
		}
	}

	// module: a group of childworfklows and continues as new on end
	if request.EnabledChildWorkflow {
		_, err = sampleschild.SampleParentWorkflow(ctx)
		if err != nil {
			logger.Error("sampleParentWorkflow fails,", "err", err)
			return "", err
		}
	}

	// module: localacitivty
	if request.EnabledLocalActivity {
		_, err := sampleslocal.GreetingSample(ctx)
		if err != nil {
			logger.Error("localActivity fails,", "err", err)
			return "", err
		}
	}

	// module: transfer in saga pattern, and has a designed error and will rollback
	if request.EnabledSagaWithTransfer {
		err = samplessaga.TransferMoney(ctx, request.TransferDetailsRequest)
		if err != nil {
			// the err is expected
			// so cannot return err to stop the wf
			logger.Error("sagaWithTransfer fails with rollback,", "err", err)
		}
	}

	// normal Activity Part
	if request.EnabledNormalActivity {
		_, err = basic.WorkflowWithNormalActivity(ctx, request.NormalActivityRequest)
		if err != nil {
			logger.Error("wf with normal activity fails,", "err", err)
			return "", err
		}
	}

	return request.ResultPayload, nil
}

func nonNegative(ctx workflow.Context, i int) error {
	logger := workflow.GetLogger(ctx)
	if i < 0 {
		logger.Debug("Rejecting negative update", "addend", i)
		return fmt.Errorf("addend must be non-negative (%v)", i)
	}
	logger.Debug("Accepting update", "addend", i)
	return nil
}
