package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/go-resty/resty/v2"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"

	"io.matrinos/docker/runner/docker"
)

var runContainer = docker.RunContainer
var pingSimulator = PingSimulator
var postStartDevice = PostStartDevice

/**
 * The activities used by running simulation workflow.
 */
const (
	runSimulationActivityName        = "runSimulationActivityName"
	startDeviceActivityName          = "startDeviceActivityName"
	pingSimulationDeviceActivityName = "pingSimulationDeviceActivityName"
)

func runSimulationActivity(ctx context.Context,
	userName string, password string, port string, imageName string,
	containerName string) (*container.ContainerCreateCreatedBody, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Running docker")
	data, err := runContainer(userName, password, imageName, containerName, port, false)

	if err != nil {
		logger.Error("Running simlation failed.", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func startDeviceActivity(ctx context.Context, port int, deviceJsonbytes []byte) ([]byte, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting Simulated Device")
	client := resty.New()

	pingUrl := fmt.Sprintf("%s:%d/ping", os.Getenv("CONTAINER_HOST"), port)
	// TODO: Extract default timeout to config file
	_, err := pingSimulator(client, pingUrl, 10, logger)

	if err != nil {
		logger.Info("Could not ping simulator before timeout.", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf("%s:%d/start", os.Getenv("CONTAINER_HOST"), port)
	resp, err := postStartDevice(client, url, deviceJsonbytes)

	if err != nil {
		logger.Info("Failed to parse start device result", zap.Error(err))
		return nil, err
	}

	return resp.Body(), nil
}

// TODO: poll device status
