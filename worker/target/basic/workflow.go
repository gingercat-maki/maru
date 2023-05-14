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
	"time"

	"github.com/temporalio/maru"
	"github.com/temporalio/maru/internal"
	"go.temporal.io/sdk/workflow"
)

func WorkflowWithNormalActivity(ctx workflow.Context, request internal.NormalActivityRequest) (string, error) {

	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		TaskQueue: maru.ActivityQueue,
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
			Payload:                   request.ActivityPayload,
			ResultPayload:             request.ActivityResultPayload,
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
	return request.ActivityPayload, nil
}
