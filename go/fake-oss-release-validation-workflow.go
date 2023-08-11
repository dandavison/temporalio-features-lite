package main

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	Status string

	RunTestsInput struct {
		TemporalAddress   string
		BenchTestDuration string
	}

	StringDatum struct {
		Name  string
		Value string
	}

	Description []StringDatum

	Link struct {
		Text string
		Url  string
	}

	RunBenchTestsActivityOutput struct {
		Links []Link
	}

	TestOutput struct {
		Description Description
		Output      RunBenchTestsActivityOutput
	}

	RunTestsOutput struct {
		TestsOutput []TestOutput
	}

	Output struct {
		Description []StringDatum
		Output      RunTestsOutput
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

	output := Output{
		Description: []StringDatum{{Name: "Run", Value: "Fake OSS Server validation run 1"}, {Name: "Server version", Value: "fake-server-release-tag"}},
	}
	err := workflow.SetQueryHandler(ctx, "getStatus", func() (Status, error) { return status, nil })
	if err != nil {
		return output, err
	}

	fmt.Println("Wf: Sleeping")
	workflow.Sleep(ctx, 10*time.Second)
	fmt.Println("Wf: Done")

	output.Output = RunTestsOutput{
		TestsOutput: []TestOutput{
			{
				Description: []StringDatum{{Name: "Test name", Value: "Test 1"}},
				Output: RunBenchTestsActivityOutput{
					Links: []Link{
						{
							Text: "Log link 1",
							Url:  "https://www.popsci.com/uploads/2023/06/06/lumber.png?auto=webp&width=1440&height=810",
						},
						{
							Text: "Log link 2",
							Url:  "https://www.perc.org/wp-content/uploads/2023/04/logging-1.jpg",
						},
					},
				},
			},
		},
	}

	status = Done
	return output, nil
}
