package builtin

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/k1nky/tookhook/internal/domain/template"
)

// HttpHandler is a built-in handler that sends HTTP requests.
type HttpHandler struct {
	logger *slog.Logger
}

// NewHttpHandler creates a new HTTP handler.
func NewHttpHandler(logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		logger: logger,
	}
}

// Name returns the handler name.
func (h *HttpHandler) Name() string {
	return "~http"
}

// Execute sends an HTTP request with the input data.
func (h *HttpHandler) Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error) {
	// Get URL from options
	urlStr, ok := opts["url"].(string)
	if !ok || urlStr == "" {
		h.logger.Error("http handler: url option is required")
		return input, nil
	}

	// Compile template for URL if it contains template syntax
	tpl := template.New(urlStr)
	if err := tpl.Compile(); err != nil {
		h.logger.Error("http handler: failed to compile URL template", "error", err)
		return input, nil
	}

	// Execute URL template with input data
	url, err := tpl.Execute(input)
	if err != nil {
		h.logger.Error("http handler: failed to execute URL template", "error", err)
		return input, nil
	}

	// Get HTTP method from options
	method, _ := opts["method"].(string)
	if method == "" {
		method = "POST"
	}
	method = strings.ToUpper(method)

	// Get headers from options
	headers, _ := opts["headers"].(map[string]any)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, io.NopCloser(bytes.NewReader(input)))
	if err != nil {
		h.logger.Error("http handler: failed to create request", "error", err)
		return input, nil
	}

	// Set Content-Type header if not provided
	if _, hasContentType := headers["Content-Type"]; !hasContentType {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers
	for key, value := range headers {
		if strVal, ok := value.(string); ok {
			req.Header.Set(key, strVal)
		}
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Error("http handler: failed to send request", "error", err)
		return input, nil
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("http handler: failed to read response body", "error", err)
		return input, nil
	}

	h.logger.Debug("http handler: request completed",
		"status", resp.Status,
		"status_code", resp.StatusCode,
		"response_size", len(body),
	)

	return body, nil
}
