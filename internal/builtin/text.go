package builtin

import (
	"context"
	"log/slog"
	"sync"

	"github.com/k1nky/tookhook/internal/domain/template"
)

// TextHandler is a built-in handler that processes input with a text template.
type TextHandler struct {
	logger    *slog.Logger
	templates sync.Map // cache of compiled templates: template string -> *template.String
}

// NewTextHandler creates a new text handler.
func NewTextHandler(logger *slog.Logger) *TextHandler {
	return &TextHandler{
		logger: logger,
	}
}

// Name returns the handler name.
func (h *TextHandler) Name() string {
	return "~text"
}

// Execute processes the input data with the template option and returns the result.
func (h *TextHandler) Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error) {
	// Get template from options
	tplStr, ok := opts["template"].(string)
	if !ok || tplStr == "" {
		h.logger.Warn("text handler: no template provided, returning input unchanged")
		return input, nil
	}

	// Get or compile template
	tpl, err := h.getOrCompileTemplate(tplStr)
	if err != nil {
		h.logger.Error("text handler: failed to compile template", "error", err)
		return input, nil
	}

	// Execute template with input data
	result, err := tpl.Execute(input)
	if err != nil {
		h.logger.Error("text handler: failed to execute template", "error", err)
		return input, nil
	}

	return []byte(result), nil
}

// getOrCompileTemplate returns a cached compiled template or compiles and caches it.
func (h *TextHandler) getOrCompileTemplate(tplStr string) (*template.String, error) {
	// Try to get from cache first
	if val, ok := h.templates.Load(tplStr); ok {
		return val.(*template.String), nil
	}

	// Compile new template
	tpl := template.New(tplStr)
	if err := tpl.Compile(); err != nil {
		return nil, err
	}

	// Store in cache
	h.templates.Store(tplStr, tpl)
	return tpl, nil
}
