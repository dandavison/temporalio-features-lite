package main

import (
	"fmt"
	"time"

	"com.github/dandavison/temporalio-features-lite/activities"
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

	Values map[string]any

	Input struct {
		KubeConfig  string
		KubeContext string

		// User provided values for the helm charts
		// Key is the chart name
		// Value is the values yaml content
		Values Values

		EnableFaultInjection bool

		// BenchTestDuration specifies how long the bench test is expected to run
		// If specified, the bench test will run multiple times until the specified duration is reached.
		// NOTE: The duration is only checked when starting a new bench test run, and not enforced during the run.
		// This means that the bench test will run longer than the specified duration and could take up to
		// the specified duration + the duration need by one run of the bench test.
		// If not specified, the bench test will only run once.
		// The format is the same as the one used by time.ParseDuration.
		BenchTestDuration string

		// BenchTestParallelism specifies how many parallel bench tests should be run.
		BenchTestParallelism int

		CallbackInput activities.PipelineCompletionCallbackInput
	}
)

const (
	NotStarted Status = "NotStarted"
	InProgress Status = "InProgress"
	Done       Status = "Done"
	Error      Status = "Error"
)

func FakeOSSReleaseValidationWorkflow(ctx workflow.Context, input Input) (Output, error) {
	fmt.Printf("Wf: Starting. Input: %v\n", input)

	var status = InProgress

	output := Output{
		Description: []StringDatum{{Name: "Run", Value: "Fake OSS Server validation run 1"}, {Name: "Server version", Value: "fake-server-release-tag"}},
	}
	err := workflow.SetQueryHandler(ctx, "getStatus", func() (Status, error) { return status, nil })
	if err != nil {
		return output, err
	}

	for i := 0; i < 4; i++ {
		fmt.Printf("Wf: Sleeping %d\n", i)
		workflow.Sleep(ctx, 10*time.Second)
	}
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

	// *** Notify CI service that pipeline is near complete ***
	// We do not allow failure of this activity to fail the workflow.
	activities.CallPipelineCompletionCallback(ctx, input.CallbackInput)

	return output, nil
}
