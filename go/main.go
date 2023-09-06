package main

import (
	"log"

	"com.github/dandavison/temporalio-features-lite/activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	var a *activities.Activities

	w := worker.New(c, "my-task-queue", worker.Options{})
	w.RegisterWorkflow(FakeOSSReleaseValidationWorkflow)
	w.RegisterWorkflow(WorkflowWithSignal)
	w.RegisterActivity(a.PipelineCompletionCallback)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
