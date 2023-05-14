package internal

import (
	"context"
	"log"

	"github.com/temporalio/samples-go/update"
	"go.temporal.io/sdk/client"
)

// trigger
func TriggerUpdates(c client.Client, we client.WorkflowRun) {
	mult := 1
	for i := 0; i < 10; i++ {
		addend := mult * i
		mult *= -1 // flip addend between negative and positive for each iteration
		handle, err := c.UpdateWorkflow(context.Background(), we.GetID(), we.GetRunID(), update.FetchAndAdd, addend)
		if err != nil {
			log.Fatal("error issuing update request", err)
		}
		var result int
		err = handle.Get(context.Background(), &result)
		if err != nil {
			log.Printf("fetch_and_add with addend %v failed: %v", addend, err)
		} else {
			log.Printf("fetch_and_add with addend %v succeeded: %v", addend, result)
		}
	}

}
