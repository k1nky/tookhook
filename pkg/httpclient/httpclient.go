package httpclient

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

const (
	DefaultRequestTimeout = 10 * time.Second
)

func SendRequest(ctx context.Context, method string, url string, data []byte) ([]byte, error) {
	client := http.Client{
		Timeout: DefaultRequestTimeout,
	}
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	buf.Reset()
	buf.ReadFrom(response.Body)
	return buf.Bytes(), nil
}
