package main

import (
	"os"
	"os/signal"

	"go.temporal.io/temporal/client"
	"go.temporal.io/temporal/worker"
	"go.uber.org/zap"

	child_workflow "github.com/temporalio/temporal-go-samples/child-workflow"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   logger,
	})
	if err != nil {
		logger.Fatal("Unable to create client", zap.Error(err))
	}
	defer c.CloseConnection()

	w := worker.New(c, "child-workflow", worker.Options{})
	defer w.Stop()

	w.RegisterWorkflow(child_workflow.SampleParentWorkflow)
	w.RegisterWorkflow(child_workflow.SampleChildWorkflow)

	err = w.Start()
	if err != nil {
		logger.Fatal("Unable to start worker", zap.Error(err))
	}

	// The workers are supposed to be long running process that should not exit.
	waitCtrlC()
}

func waitCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
