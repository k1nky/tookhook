package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultClientTimeout   = 15 * time.Second
	TelegramApiMessagesUrl = "https://api.telegram.org/bot%s/sendMessage"
)

type Adapter struct {
	Endpoint string
}

type Message struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func New() *Adapter {
	return &Adapter{
		Endpoint: TelegramApiMessagesUrl,
	}
}

func (a *Adapter) Send(token string, msg Message) ([]byte, error) {
	urlEndpoint := fmt.Sprintf(a.Endpoint, token)
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(msg); err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", urlEndpoint, buf)
	request.Header.Set("content-type", "application/json")
	if err != nil {
		return nil, err
	}
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
