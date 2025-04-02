package builtin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/k1nky/tookhook/pkg/httpclient"
	"github.com/k1nky/tookhook/pkg/plugin"
)

var (
	HttpHandlerAllowrdMethods = []string{http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodTrace}
)

type HttpHandler struct {
	builtinPlugin
}

//go:generate easyjson http.go
//easyjson:json
type HttpHandlerOptions struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

func NewHttpHandlerOptions(encoded []byte) (HttpHandlerOptions, error) {
	opts := &HttpHandlerOptions{}
	err := json.Unmarshal(encoded, opts)
	return *opts, err
}

func NewHttpHandler(log logger) *HttpHandler {
	return &HttpHandler{
		builtinPlugin: builtinPlugin{
			Logger: log,
		},
	}
}

func (opts *HttpHandlerOptions) Validate() error {
	return nil
}

func (h *HttpHandler) Validate(ctx context.Context, r plugin.Handler) error {
	opts, err := NewHttpHandlerOptions(r.Options)
	if err != nil {
		return err
	}
	return opts.Validate()
}

func (h *HttpHandler) Forward(ctx context.Context, r plugin.Handler, data []byte) ([]byte, error) {
	opts, err := NewHttpHandlerOptions(r.Options)
	if err != nil {
		return nil, err
	}
	h.Logger.Debugf("http handler: %s %s", opts.Method, opts.URL)
	body, err := httpclient.SendRequest(ctx, opts.Method, opts.URL, data)
	if err != nil {
		return nil, err
	}
	return body, nil
}
