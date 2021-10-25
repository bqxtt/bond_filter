package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpClient struct {
	client *http.Client
}

var Client = &HttpClient{client: &http.Client{}}

func (client *HttpClient) Request(method string, url string, header http.Header, payload string) (string, error) {
	//var requestBody *strings.Reader
	//if payload != "" {
	requestBody := strings.NewReader(payload)

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return "", err
	}
	if header != nil {
		req.Header = header
	}
	res, err := client.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	ResponseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(ResponseBody), nil
}
