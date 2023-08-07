package main

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	Status string
	Output struct {
		LogUrls []string
		Success bool
	}
)

const (
	NotStarted Status = "NotStarted"
	InProgress Status = "InProgress"
	Done       Status = "Done"
	Error      Status = "Error"
)

func FakeOSSReleaseValidationWorkflow(ctx workflow.Context) (Output, error) {
	fmt.Println("Wf: Starting")
	var status = InProgress

	err := workflow.SetQueryHandler(ctx, "getStatus", func() (Status, error) { return status, nil })
	if err != nil {
		return Output{
			LogUrls: []string{},
			Success: false,
		}, err
	}

	fmt.Println("Wf: Sleeping")
	workflow.Sleep(ctx, 30*time.Second)
	fmt.Println("Wf: Done")

	status = Done
	return Output{
		LogUrls: []string{"https://www.popsci.com/uploads/2023/06/06/lumber.png?auto=webp&width=1440&height=810"},
		Success: true,
	}, nil
}
