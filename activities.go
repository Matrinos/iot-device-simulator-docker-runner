package main

import (
	"context"
	"fmt"

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
	containerName string, networkName string) (*container.ContainerCreateCreatedBody, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Running docker")
	data, err := runContainer(userName, password, imageName, containerName, port, networkName, false)

	if err != nil {
		logger.Error("Running simlation failed.", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

// TODO: port is not useful coz we are accessing via internal network
func startDeviceActivity(ctx context.Context, containerName string,
	port int, deviceJsonbytes []byte) ([]byte, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting Simulated Device")
	client := resty.New()

	// TODO: https??
	pingURL := fmt.Sprintf("http://%s:%d/ping", containerName, 8080)
	// TODO: Extract default timeout to config file
	_, err := pingSimulator(client, pingURL, 30, logger)

	if err != nil {
		logger.Info("Could not ping simulator before timeout.", zap.Error(err))
		return nil, err
	}

	// TODO: https??
	url := fmt.Sprintf("http://%s:%d/start", containerName, 8080)
	resp, err := postStartDevice(client, url, deviceJsonbytes)

	if err != nil {
		logger.Info("Failed to parse start device result", zap.Error(err))
		return nil, err
	}

	return resp.Body(), nil
}

// TODO: poll device status
