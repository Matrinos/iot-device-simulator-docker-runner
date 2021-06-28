package main

import (
	"errors"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type PingSimulatorSuite struct {
	suite.Suite

	logger *zap.Logger
	client *resty.Client
}

func (s *PingSimulatorSuite) SetupTest() {
	s.logger, _ = zap.NewProduction()
	s.client = resty.New()
}

func (s *PingSimulatorSuite) TestPingDevice() {
	originalGetPingSimulator := getPingSimulator
	defer func() { getPingSimulator = originalGetPingSimulator }()
	getPingSimulator = func(client *resty.Client, url string) (*resty.Response, error) {
		return &resty.Response{}, nil
	}

	result, _ := PingSimulator(s.client, "", 10, s.logger)
	s.Equal(true, result)
}

func (s *PingSimulatorSuite) TestPingDeviceTimeout() {
	originalGetPingSimulator := getPingSimulator
	defer func() { getPingSimulator = originalGetPingSimulator }()
	getPingSimulator = func(client *resty.Client, url string) (*resty.Response, error) {
		return nil, errors.New("error")
	}

	result, _ := PingSimulator(s.client, "", 1, s.logger)
	s.Equal(false, result)
}

func TestPingSimulatorSuite(t *testing.T) {
	suite.Run(t, new(PingSimulatorSuite))
}
