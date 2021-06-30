package main

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/pborman/uuid"
	"github.com/phayes/freeport"
	"github.com/teris-io/shortid"
	"go.uber.org/cadence"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// type (
// 	fileInfo struct {
// 		FileName string
// 		HostID   string
// 	}
// )

// ApplicationName is the task list for this sample
const ApplicationName = "SimulatorRunningGroup"

// HostID - Use a new uuid just for demo so we can run 2 host specific activity workers on same machine.
// In real world case, you would use a hostname or ip address as HostID.
var HostID = ApplicationName + "_" + uuid.New()

//sampleFileProcessingWorkflow workflow decider
func simulatorStartingWorkflow(ctx workflow.Context, deviceJsonBytes []byte) (err error) {
	// step 1: download resource file
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Second * 5,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 300, // need debug to understand the right timeout setting.
		RetryPolicy: &cadence.RetryPolicy{
			InitialInterval:          time.Second * 30,
			BackoffCoefficient:       2.0,
			MaximumInterval:          time.Minute,
			ExpirationInterval:       time.Minute * 10,
			NonRetriableErrorReasons: []string{"bad-error"},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	port, err := freeport.GetFreePort()
	if err != nil {
		return err
	}

	sid, err := shortid.New(1, shortid.DefaultABC, uint64(time.Now().UnixNano()))
	if err != nil {
		return err
	}

	shortID, _ := sid.Generate()
	containerName := "matrinos.sim-" + shortID

	// Retry the whole sequence from the first activity on any error
	// to retry it on a different host. In a real application it might be reasonable to
	// retry individual activities and the whole sequence discriminating between different types of errors.
	// See the retryactivity sample for a more sophisticated retry implementation.
	// TODO enable loop retry
	// for i := 1; i < 5; i++ {
	err = runDocker(ctx, strconv.Itoa(port), containerName)
	if err != nil {
		workflow.GetLogger(ctx).Error("Workflow failed.", zap.String("Error", err.Error()))
		return err
	}
	// }

	// call the start end point with device parameter
	err = StartDevice(ctx, containerName, port, deviceJsonBytes)

	if err != nil {
		workflow.GetLogger(ctx).Error("Workflow failed.", zap.String("Error", err.Error()))
		return
	}

	workflow.GetLogger(ctx).Info("Workflow completed.")
	return nil
}

func runDocker(ctx workflow.Context, port string, containerName string) (err error) {
	var containerResponse *container.ContainerCreateCreatedBody
	so := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: time.Minute,
	}

	sessionCtx, err := workflow.CreateSession(ctx, so)
	if err != nil {
		return err
	}
	defer workflow.CompleteSession(sessionCtx)
	userName := os.Getenv("DOCKERHUB_USERNAME")

	if userName == "" {
		return errors.New("please specify the username for pulling docker image")
	}

	password := os.Getenv("DOCKERHUB_TOKEN")
	if password == "" {
		return errors.New("please specify the username for pulling docker image")
	}

	networkName := os.Getenv("DOCKER_NETWORK")
	if networkName == "" {
		return errors.New("please specify the networkName for running the docker")
	}

	// TODO: move docker image name to config?
	err = workflow.ExecuteActivity(sessionCtx, runSimulationActivityName,
		userName, password,
		port, "matrinos/iot-smart-device-simulator",
		containerName, networkName).Get(sessionCtx, &containerResponse)

	if err != nil {
		return err
	}

	// var fInfoProcessed *fileInfo
	// err = workflow.ExecuteActivity(sessionCtx, processFileActivityName, *fInfo).Get(sessionCtx, &fInfoProcessed)
	// if err != nil {
	// 	return err
	// }

	// err = workflow.ExecuteActivity(sessionCtx, uploadFileActivityName, *fInfoProcessed).Get(sessionCtx, nil)
	// return err
	return nil
}

func StartDevice(ctx workflow.Context, containerName string,
	port int, deviceJsonBytes []byte) (err error) {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	if err != nil {
		return err
	}

	var res []byte

	err = workflow.ExecuteActivity(ctx,
		startDeviceActivityName,
		containerName,
		port, deviceJsonBytes).Get(ctx, &res)

	if err != nil {
		return err
	}

	return nil
}
