// Package builtin provides built-in handler implementations.
package builtin

import (
	"context"
	"log/slog"

	"github.com/k1nky/tookhook/internal/domain/template"
)

// LogHandler is a built-in handler that logs data.
type LogHandler struct {
	logger *slog.Logger
	level  slog.Level
}

// NewLogHandler creates a new log handler.
func NewLogHandler(logger *slog.Logger) *LogHandler {
	return &LogHandler{
		logger: logger,
		level:  slog.LevelInfo,
	}
}

// Name returns the handler name.
func (h *LogHandler) Name() string {
	return "~log"
}

// Execute logs the input data and returns it unchanged.
func (h *LogHandler) Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error) {
	level := h.level
	if lvlStr, ok := opts["level"].(string); ok {
		switch lvlStr {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	// Check if custom message template is provided
	if msgTpl, ok := opts["message"].(string); ok && msgTpl != "" {
		tpl := template.New(msgTpl)
		if err := tpl.Compile(); err != nil {
			h.logger.Error("failed to compile message template", "error", err)
			return input, nil
		}
		message, err := tpl.Execute(input)
		if err != nil {
			h.logger.Error("failed to execute message template", "error", err)
			return input, nil
		}
		h.logger.Log(ctx, level, message)
	} else {
		h.logger.Log(ctx, level, "handler log", "payload", string(input))
	}

	return input, nil
}
