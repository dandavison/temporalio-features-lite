package main

import (
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

	Input struct {
		CallbackData CallbackData `json:"callback-data"`
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

	for i := 0; i < 10; i++ {
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
	if err = workflow.ExecuteActivity(ctx, NotifyWorkflowComplete, input.CallbackData).Get(ctx, nil); err != nil {
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
