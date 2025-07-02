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

type Request struct {
	Method  string
	URL     string
	Timeout time.Duration
}

func SendRequest(ctx context.Context, r Request, data []byte) ([]byte, error) {
	if r.Timeout == 0 {
		r.Timeout = DefaultRequestTimeout
	}
	client := http.Client{
		Timeout:   r.Timeout,
		Transport: http.DefaultTransport,
	}
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, buf)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if response != nil {
		defer response.Body.Close()
		buf.Reset()
		buf.ReadFrom(response.Body)
	}
	return buf.Bytes(), err
}
