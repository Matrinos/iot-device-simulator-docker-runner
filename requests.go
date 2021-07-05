package main

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

func PostStartDevice(client *resty.Client, url string, body []byte) (*resty.Response, error) {
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}

func HttpGet(client *resty.Client, url string, result interface{}) error {
	resp, err := client.R().
		Get(url)
	if err != nil {
		return err
	}
	json.Unmarshal(resp.Body(), result)
	return nil
}
