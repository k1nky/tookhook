package pachca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultClientTimeout = 5 * time.Second
	PachcaApiMessagesUrl = "https://api.pachca.com/api/shared/v1/messages"
)

type Adapter struct {
	token    string
	Endpoint string
}

type MessagePayload struct {
	Message Message `json:"message"`
}

type Message struct {
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
	Content    string `json:"content"`
}

func NewPachca(token string) *Adapter {
	return &Adapter{
		token:    token,
		Endpoint: PachcaApiMessagesUrl,
	}
}

func (a *Adapter) Send(payload MessagePayload) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", a.Endpoint, buf)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
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
	if response.StatusCode != http.StatusCreated {
		return body, fmt.Errorf("unexpected response with code %d", response.StatusCode)
	}
	return body, err
}
