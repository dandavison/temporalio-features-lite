package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

	CallbackData struct {
		Method  string            `json:"method"`
		Url     string            `json:"url"`
		Headers map[string]string `json:"headers"`
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

	callbackData, err := getCallbackData(input)
	if err != nil {
		return Output{}, err
	}

	var status = InProgress

	output := Output{
		Description: []StringDatum{{Name: "Run", Value: "Fake OSS Server validation run 1"}, {Name: "Server version", Value: "fake-server-release-tag"}},
	}
	err = workflow.SetQueryHandler(ctx, "getStatus", func() (Status, error) { return status, nil })
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

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})
	if err = workflow.ExecuteActivity(ctx, NotifyWorkflowComplete, callbackData).Get(ctx, nil); err != nil {
		return output, err
	}

	return output, nil
}

func NotifyWorkflowComplete(callbackData CallbackData) error {
	client := &http.Client{}
	req, err := http.NewRequest(callbackData.Method, callbackData.Url, nil)
	if err != nil {
		return err
	}
	curlCmd := fmt.Sprintf("curl -L -X %s", callbackData.Method)
	for k, v := range callbackData.Headers {
		req.Header.Add(k, v)
		curlCmd += fmt.Sprintf(` -H "%s: %s"`, k, v)
	}
	curlCmd += " " + callbackData.Url
	fmt.Println(curlCmd)

	client.Do(req)
	return nil
}

func getCallbackData(input Input) (CallbackData, error) {
	value, ok := input.Values["GithubCallbackData"]
	if !ok {
		return CallbackData{}, errors.New("GithubCallbackData not found in input")
	}

	s, err := json.Marshal(value)
	if err != nil {
		return CallbackData{}, errors.New("failed to marshal GithubCallbackData to JSON")
	}

	var callbackData CallbackData
	if err := json.Unmarshal(s, &callbackData); err != nil {
		return CallbackData{}, errors.New("failed to parse GithubCallbackData JSON")
	}
	return callbackData, nil
}
