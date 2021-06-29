package main

import "github.com/go-resty/resty/v2"

func PostStartDevice(client *resty.Client, url string, body []byte) (*resty.Response, error) {
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}

func GetPingSimulator(client *resty.Client, url string) (*resty.Response, error) {
	return client.R().
		Get(url)
}
