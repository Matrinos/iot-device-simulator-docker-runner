package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/go-resty/resty/v2"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"

	"io.matrinos/docker/runner/docker"
)

var runContainer = docker.RunContainer

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

	pingUrl := fmt.Sprintf("%s:%d/start", os.Getenv("CONTAINER_HOST"), port)

	url := fmt.Sprintf("%s:%d/start", os.Getenv("CONTAINER_HOST"), port)

	pong := false
	sleep := 1
	for !pong {
		_, err := client.R().
			Get(pingUrl)
		if err != nil {
			logger.Info("Ping simulator failed", zap.Error(err))
			time.Sleep(time.Second * time.Duration(sleep))
			sleep = sleep * 2
		} else {
			logger.Info("Simulator device online!")
			pong = true
		}
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(deviceJsonbytes).
		Post(url)

	if err != nil {
		logger.Info("Failed to parse start device result", zap.Error(err))
		return nil, err
	}
	return resp.Body(), nil

}

// TODO: poll device status
