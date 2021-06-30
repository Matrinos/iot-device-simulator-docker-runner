package main

import (
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var getPingSimulator = GetPingSimulator

func PingSimulator(client *resty.Client, pingUrl string, timeoutSeconds int, logger *zap.Logger) (pong bool, err error) {
	pong = false
	sleepSeconds := 5
	timer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	for !pong {
		select {
		case <-timer.C:
			// time's up
			return pong, errors.New("ping device timed out")
		default:
			_, err := getPingSimulator(client, pingUrl)
			if err != nil {
				logger.Info("Ping simulator failed", zap.Error(err))
				time.Sleep(time.Second * time.Duration(sleepSeconds))
				// sleepSeconds = sleepSeconds * 2
			} else {
				logger.Info("Simulator device online!")
				pong = true
			}
		}
	}
	return pong, nil
}
