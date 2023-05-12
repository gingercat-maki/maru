package internal

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

var logger zap.Logger*

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

// trigger methods
func TriggerAwaitSignal(c client.Client, we client.WorkflowRun) {
	signals := []int{1, 2, 3}
	rand.Shuffle(len(signals), func(i, j int) { signals[i], signals[j] = signals[j], signals[i] })
	for _, signal := range signals {
		signalName := fmt.Sprintf("Signal%d", signal)
		err := c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), signalName, nil)
		if err != nil {
			logger.Fatal("Unable to signals workflow", zap.Error(err))
		}
		time.Sleep(2 * time.Second)
	}
}

// workflow patterns: no extra config needed
