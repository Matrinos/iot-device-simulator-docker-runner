package main

import (
	"context"
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
	s.env.RegisterActivityWithOptions(runSimulationActivity, activity.RegisterOptions{
		Name: runSimulationActivityName,
	})
	s.env.RegisterActivityWithOptions(startDeviceActivity, activity.RegisterOptions{
		Name: startDeviceActivityName,
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
		shouldWait bool) (
		container.ContainerCreateCreatedBody, error) {
		// This will be called, do whatever you want to,
		// return whatever you want to
		return container.ContainerCreateCreatedBody{}, nil
	}

	originalPingSimulator := pingSimulator
	defer func() { pingSimulator = originalPingSimulator }()

	pingSimulator = func(client *resty.Client, pingUrl string, logger *zap.Logger) bool {
		return true
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
