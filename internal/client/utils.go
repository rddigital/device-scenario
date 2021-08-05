package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

// SendRequest will make a request with raw data to the specified URL.
// It returns the body as a byte array if successful and an error otherwise.
func SendRequest(baseUrl string, api string, method string, data []byte) (response []byte, err error) {
	url := baseUrl + "/" + api
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create a http request %s", err.Error())
	}
	req.Header.Set(ContentType, ContentTypeJSON)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send a http request %s", err.Error())
	}

	if resp == nil {
		return nil, fmt.Errorf("the response should not be a nil")
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode <= http.StatusMultiStatus {
		return bodyBytes, nil
	}
	return nil, fmt.Errorf("request failed, status code: %d, err: %s", resp.StatusCode, string(bodyBytes))
}
