package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/encoded"
	"go.uber.org/cadence/testsuite"
	"go.uber.org/zap"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.env.RegisterWorkflow(simulatorStartingWorkflow)
	s.env.RegisterWorkflow(simulatorStatusWorkflow)
	s.env.RegisterActivityWithOptions(runSimulationActivity, activity.RegisterOptions{
		Name: runSimulationActivityName,
	})
	s.env.RegisterActivityWithOptions(startDeviceActivity, activity.RegisterOptions{
		Name: startDeviceActivityName,
	})
	s.env.RegisterActivityWithOptions(getSimulatorStatusActivity, activity.RegisterOptions{
		Name: getSimulatorStatusActivityName,
	})

}

func (s *UnitTestSuite) TearDownTest() {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_RunDockerProcessingWorkflow() {
	expectedCall := []string{
		"runSimulationActivityName",
		"startDeviceActivityName",
	}

	var activityCalled []string
	s.env.SetOnActivityStartedListener(
		func(activityInfo *activity.Info, ctx context.Context, args encoded.Values) {
			activityType := activityInfo.ActivityType.Name
			if strings.HasPrefix(activityType, "internalSession") {
				return
			}
			activityCalled = append(activityCalled, activityType)
			//TODO: verify args here.
			// switch activityType {
			// case expectedCall[0]:
			// 	var input string
			// 	s.NoError(args.Get(&input))
			// 	s.Equal(fileID, input)
			// default:
			// 	panic("unexpected activity call")
			// }
		})

	old := runContainer
	defer func() { runContainer = old }()

	runContainer = func(userName string,
		password string,
		imageName string,
		containerName string,
		port string,
		networkName string,
		shouldWait bool) (
		container.ContainerCreateCreatedBody, error) {
		// This will be called, do whatever you want to,
		// return whatever you want to
		return container.ContainerCreateCreatedBody{}, nil
	}

	originalPingSimulator := pingSimulator
	defer func() { pingSimulator = originalPingSimulator }()

	pingSimulator = func(client *resty.Client, pingUrl string, durationSeconds int, logger *zap.Logger) (bool, error) {
		return true, nil
	}

	originalPostDevice := postStartDevice
	defer func() { postStartDevice = originalPostDevice }()

	postStartDevice = func(client *resty.Client, url string, body []byte) (*resty.Response, error) {
		return &resty.Response{}, nil
	}

	var deviceJsonBytes = []byte("")

	s.env.ExecuteWorkflow(simulatorStartingWorkflow, deviceJsonBytes)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	s.Equal(expectedCall, activityCalled)
}

func (s *UnitTestSuite) TestGetSimulatorStatusWorkflowSuccessful() {
	containerName := "mockContainerName"

	expectedCall := []string{
		getSimulatorStatusActivityName,
	}

	var activityCalled []string
	s.env.SetOnActivityStartedListener(
		func(activityInfo *activity.Info, ctx context.Context, args encoded.Values) {
			activityType := activityInfo.ActivityType.Name
			if strings.HasPrefix(activityType, "internalSession") {
				return
			}
			activityCalled = append(activityCalled, activityType)
			switch activityType {
			case expectedCall[0]:
				var input string
				s.NoError(args.Get(&input))
				s.Equal(containerName, input)
			default:
				panic("unexpected activity call")
			}
		})

	originalHttpGet := httpGet
	defer func() { httpGet = originalHttpGet }()

	httpGet = func(client *resty.Client, url string, res interface{}) error {
		err := json.Unmarshal([]byte("{\"status\":\"running\"}"), &res)
		if err != nil {
			panic(err)
		}
		return nil
	}

	s.env.ExecuteWorkflow(simulatorStatusWorkflow, containerName)
	var result SimulatorStatusResult
	err := s.env.GetWorkflowResult(&result)
	if err != nil {
		panic(err)
	}
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	s.Equal(expectedCall, activityCalled)
	s.Equal(Running, result.Status)
}

func (s *UnitTestSuite) TestGetSimulatorStatusWorkflowError() {
	containerName := "mockContainerName"

	expectedCall := []string{
		getSimulatorStatusActivityName,
	}

	var activityCalled []string
	s.env.SetOnActivityStartedListener(
		func(activityInfo *activity.Info, ctx context.Context, args encoded.Values) {
			activityType := activityInfo.ActivityType.Name
			if strings.HasPrefix(activityType, "internalSession") {
				return
			}
			activityCalled = append(activityCalled, activityType)
			switch activityType {
			case expectedCall[0]:
				var input string
				s.NoError(args.Get(&input))
				s.Equal(containerName, input)
			default:
				panic("unexpected activity call")
			}
		})

	originalHttpGet := httpGet
	defer func() { httpGet = originalHttpGet }()

	httpGet = func(client *resty.Client, url string, res interface{}) error {
		return errors.New("Error")
	}

	s.env.ExecuteWorkflow(simulatorStatusWorkflow, containerName)
	var result SimulatorStatusResult
	err := s.env.GetWorkflowResult(&result)
	if err != nil {
		panic(err)
	}
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	s.Equal(expectedCall, activityCalled)
	s.Equal(Error, result.Status)
}
