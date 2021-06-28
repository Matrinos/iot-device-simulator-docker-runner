package main

import (
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func PingSimulator(client *resty.Client, pingUrl string, logger *zap.Logger) bool {
	pong := false
	sleepSeconds := 1
	for !pong {
		_, err := client.R().
			Get(pingUrl)
		if err != nil {
			logger.Info("Ping simulator failed", zap.Error(err))
			time.Sleep(time.Second * time.Duration(sleepSeconds))
			sleepSeconds = sleepSeconds * 2
		} else {
			logger.Info("Simulator device online!")
			pong = true
		}
	}
	return pong
}
