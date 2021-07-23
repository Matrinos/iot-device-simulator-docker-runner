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

func main() {

	var h core.WorkflowHelper
	h.SetupServiceConfig()

	h.RegisterWorkflow(simulatorStatusWorkflow)
	h.RegisterWorkflow(simulatorStartingWorkflow)
	h.RegisterWorkflow(simulatorStopWorkflow)

	h.RegisterActivityWithAlias(runSimulationActivity, runSimulationActivityName)
	h.RegisterActivityWithAlias(startDeviceActivity, startDeviceActivityName)
	h.RegisterActivityWithAlias(getSimulatorStatusActivity, getSimulatorStatusActivityName)
	h.RegisterActivityWithAlias(stopDeviceActivity, stopDeviceActivityName)

	startWorkers(&h)

	// The workers are supposed to be long running process that should not exit.
	// Use select{} to block indefinitely for samples, you can quit by CMD+C.
	select {}
}
