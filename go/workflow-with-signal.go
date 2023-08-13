package main

import (
	"fmt"

	"go.temporal.io/sdk/workflow"
)

const (
	SignalName = "MySignal"
)

func WorkflowWithSignal(ctx workflow.Context) error {
	fmt.Printf("Wf: Starting")
	state := 0
	signalCh := workflow.GetSignalChannel(ctx, "MySignal")
	signalCh.Receive(ctx, &state)
	fmt.Printf("Received signal: state = %d\n", state)
	signalCh.Receive(ctx, &state)
	fmt.Printf("Received signal: state = %d\n", state)
	signalCh.Receive(ctx, &state)
	fmt.Printf("Received signal: state = %d\n", state)
	return nil
}
