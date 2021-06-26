package main

import (
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func PingSimulator(client *resty.Client, pingUrl string, logger *zap.Logger) {
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
}
