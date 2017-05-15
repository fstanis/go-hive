package hive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	requestOrigin  = "https://my.hivehome.com"
	headerKeyToken = "Authorization"
)

var (
	customHeaders = [][2]string{
		[2]string{"Content-Type", "application/json"},
		[2]string{"Accept", "*/*"},
		[2]string{"Origin", requestOrigin},
	}
)

type httpClient struct {
	client http.Client
}

func (c *httpClient) prepareRequest(method string, url string, jsonStr []byte) (*http.Request, error) {
	body := io.Reader(nil)
	if jsonStr != nil {
		body = bytes.NewBuffer(jsonStr)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for _, header := range customHeaders {
		req.Header.Set(header[0], header[1])
	}
	return req, nil
}

func (c *httpClient) sendRequest(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return checkError(resp.StatusCode, body)
}

func checkError(status int, body []byte) ([]byte, error) {
	if status == http.StatusOK {
		return body, nil
	}

	var errorData jsonError
	json.Unmarshal(body, &errorData)
	if errorData.ErrorText == nil {
		return body, fmt.Errorf("got HTTP status %s", http.StatusText(status))
	}
	return body, fmt.Errorf("got HTTP status %s, error text %q", http.StatusText(status), *errorData.ErrorText)
}

func (c *httpClient) PostJSON(url string, jsonStr []byte, token string) ([]byte, error) {
	req, err := c.prepareRequest("POST", url, jsonStr)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set(headerKeyToken, token)
	}
	return c.sendRequest(req)
}

func (c *httpClient) Get(url string, token string) ([]byte, error) {
	req, err := c.prepareRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set(headerKeyToken, token)
	}
	return c.sendRequest(req)
}
