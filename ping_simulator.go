package main

import (
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var getPingSimulator = HttpGet

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
			// Delay before first ping. It is most likely the
			// sim is still not ready at this time.
			time.Sleep(time.Second * time.Duration(sleepSeconds))
			// var result interface{}
			err = getPingSimulator(client, pingUrl, nil)
			if err != nil {
				logger.Info("Ping simulator failed", zap.Error(err))
				sleepSeconds = sleepSeconds * 2
			} else {
				logger.Info("Simulator device online!")
				pong = true
			}
		}
	}
	return pong, nil
}
