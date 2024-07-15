package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultClientTimeout = 5 * time.Second
)

type Adapter struct {
	endpoint string
	token    string
}

func NewBamboo(token string, endpoint string) *Adapter {
	return &Adapter{
		endpoint: endpoint,
		token:    token,
	}
}

func (a *Adapter) GetDeployProjectById(projectID string) ([]byte, error) {
	requestUrl, _ := url.JoinPath(a.endpoint, "/deploy/project", projectID)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+a.token)

	client := &http.Client{
		Timeout: DefaultClientTimeout,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return body, fmt.Errorf("unexpected response with code %d", response.StatusCode)
	}
	return body, err
}
