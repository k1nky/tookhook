package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// JsHandler is a built-in handler that executes JavaScript code.
type JsHandler struct {
	logger   *slog.Logger
	pool     *sync.Pool
	programs sync.Map // cache of compiled programs: code string -> *goja.Program
}

// NewJsHandler creates a new JavaScript handler.
func NewJsHandler(logger *slog.Logger) *JsHandler {
	return &JsHandler{
		logger: logger,
		pool: &sync.Pool{
			New: func() any {
				return newRuntime(logger)
			},
		},
	}
}

func newRuntime(logger *slog.Logger) *goja.Runtime {
	vm := goja.New()

	// Export HTTP object for JavaScript
	vm.Set("HTTP", map[string]any{
		"request": func(method, url string, body any, headers map[string]any) (any, error) {
			// Get body as string
			var bodyReader io.Reader
			if body != nil {
				var bodyStr string
				if b, ok := body.(string); ok {
					bodyStr = b
				} else {
					data, _ := json.Marshal(body)
					bodyStr = string(data)
				}
				bodyReader = bytes.NewBufferString(bodyStr)
			}

			// Create HTTP request context
			reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(reqCtx, method, url, bodyReader)
			if err != nil {
				return nil, err
			}

			// Set Content-Type if not provided and body exists
			if bodyReader != nil && req.Header.Get("Content-Type") == "" {
				req.Header.Set("Content-Type", "application/json")
			}

			// Set custom headers
			for key, value := range headers {
				if strVal, ok := value.(string); ok {
					req.Header.Set(key, strVal)
				}
			}

			// Use native Go HTTP client
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			// Read response body
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			// Try to parse response as JSON
			var result any
			if err := json.Unmarshal(respBody, &result); err == nil {
				return result, nil
			}
			return string(respBody), nil
		},
	})

	// Export log functions
	vm.Set("log", func(args ...any) {
		logger.Info(fmt.Sprint(args...))
	})
	vm.Set("error", func(args ...any) {
		logger.Error(fmt.Sprint(args...))
	})
	vm.Set("warn", func(args ...any) {
		logger.Warn(fmt.Sprint(args...))
	})
	vm.Set("debug", func(args ...any) {
		logger.Debug(fmt.Sprint(args...))
	})

	return vm
}

// Name returns the handler name.
func (h *JsHandler) Name() string {
	return "~js"
}

// Execute executes JavaScript code with the input data and options.
func (h *JsHandler) Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error) {
	code, ok := opts["code"].(string)
	if !ok || code == "" {
		h.logger.Error("js handler: code option is required")
		return input, nil
	}

	// Get runtime from pool
	vm := h.pool.Get().(*goja.Runtime)
	defer h.pool.Put(vm)

	// Set timeout for script execution
	timeout := time.AfterFunc(5*time.Second, func() {
		vm.Interrupt("execution timeout")
	})
	defer timeout.Stop()

	// Reset the runtime by clearing all exported values
	vm.Set("input", nil)
	vm.Set("Options", nil)
	vm.Set("result", nil)

	// Make input available to JavaScript
	var inputJSON any
	if err := json.Unmarshal(input, &inputJSON); err == nil {
		vm.Set("input", inputJSON)
	} else {
		vm.Set("input", string(input))
	}

	// Make options available to JavaScript
	vm.Set("Options", opts)

	// Get or compile the program
	program, err := h.getOrCompileProgram(code)
	if err != nil {
		h.logger.Error("js handler: failed to compile script", "error", err)
		return input, nil
	}

	// Run the compiled program
	_, err = vm.RunProgram(program)
	if err != nil {
		h.logger.Error("js handler: failed to run script", "error", err)
		return input, nil
	}

	// Get the result from the last statement or return value
	result := vm.Get("result")
	if result == nil || goja.IsUndefined(result) {
		// If no explicit result, return input
		return input, nil
	}

	// Convert result to JSON
	jsonResult, err := json.Marshal(result.Export())
	if err != nil {
		h.logger.Error("js handler: failed to marshal result", "error", err)
		return input, nil
	}

	return jsonResult, nil
}

// getOrCompileProgram returns a cached compiled program or compiles and caches it.
func (h *JsHandler) getOrCompileProgram(code string) (*goja.Program, error) {
	// Try to get from cache first
	if val, ok := h.programs.Load(code); ok {
		return val.(*goja.Program), nil
	}

	// Compile new program
	program, err := goja.Compile("", code, false)
	if err != nil {
		return nil, err
	}

	// Store in cache
	h.programs.Store(code, program)
	return program, nil
}
