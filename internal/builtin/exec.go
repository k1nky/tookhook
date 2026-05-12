package builtin

import (
	"context"
	"log/slog"
)

// ExecHandler is a built-in handler that executes commands.
type ExecHandler struct {
	logger *slog.Logger
}

// NewExecHandler creates a new exec handler.
func NewExecHandler(logger *slog.Logger) *ExecHandler {
	return &ExecHandler{
		logger: logger,
	}
}

// Name returns the handler name.
func (h *ExecHandler) Name() string {
	return "~exec"
}

// Execute runs a command with the input data.
// TODO: Implement actual command execution.
func (h *ExecHandler) Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error) {
	command, _ := opts["command"].(string)

	h.logger.Debug("handler exec stub",
		"command", command,
		"payload_size", len(input),
	)

	// Stub: return input unchanged
	return input, nil
}
