package main

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"

	"io.matrinos/docker/runner/docker"
)

/**
 * The activities used by running simulation workflow.
 */
const (
	runSimulationActivityName = "runSimulationActivityName"
)

func runSimulationActivity(ctx context.Context,
	userName string, password string, port string, imageName string,
	containerName string) (*container.ContainerCreateCreatedBody, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Running docker")
	data, err := docker.RunContainer(userName, password, imageName, containerName, port, false)

	if err != nil {
		logger.Error("Running simlation failed.", zap.Error(err))
		return nil, err
	}

	return &data, nil
}
