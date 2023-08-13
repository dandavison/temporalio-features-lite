package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "my-task-queue", worker.Options{})
	w.RegisterWorkflow(FakeOSSReleaseValidationWorkflow)
	w.RegisterWorkflow(WorkflowWithSignal)
	w.RegisterActivity(NotifyWorkflowComplete)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
