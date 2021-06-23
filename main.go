package main

import (
	"github.com/Matrinos/iot-cadence-go-core/core"
	"go.uber.org/cadence/worker"
)

// This needs to be done as part of a bootstrap step when the process starts.
// The workers are supposed to be long running.
func startWorkers(h *core.WorkflowHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope:          h.WorkerMetricScope,
		Logger:                h.Logger,
		EnableLoggingInReplay: true,
		EnableSessionWorker:   true,
	}
	h.StartWorkers(h.Config.DomainName, ApplicationName, workerOptions)

	// Host Specific activities processing case
	workerOptions.DisableWorkflowWorker = true
	h.StartWorkers(h.Config.DomainName, HostID, workerOptions)
}

// func startWorkflow(h *core.WorkflowHelper, fileID string) {
// 	workflowOptions := client.StartWorkflowOptions{
// 		ID:                              "fileprocessing_" + uuid.New(),
// 		TaskList:                        ApplicationName,
// 		ExecutionStartToCloseTimeout:    time.Minute,
// 		DecisionTaskStartToCloseTimeout: time.Minute,
// 	}
// 	h.StartWorkflow(workflowOptions, "main.simulatorStartingWorkflow", fileID)
// }

func main() {
	// var mode string
	// flag.StringVar(&mode, "m", "trigger", "Mode is worker or trigger.")
	// flag.Parse()

	var h core.WorkflowHelper
	h.SetupServiceConfig()

	// switch mode {
	// case "worker":
	h.RegisterWorkflow(simulatorStartingWorkflow)
	h.RegisterActivityWithAlias(runSimulationActivity, runSimulationActivityName)
	// h.RegisterActivityWithAlias(processFileActivity, processFileActivityName)
	// h.RegisterActivityWithAlias(uploadFileActivity, uploadFileActivityName)
	startWorkers(&h)

	// The workers are supposed to be long running process that should not exit.
	// Use select{} to block indefinitely for samples, you can quit by CMD+C.
	select {}
	// case "trigger":
	// 	startWorkflow(&h, uuid.New())
	// }
}
