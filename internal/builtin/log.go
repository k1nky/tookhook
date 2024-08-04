package builtin

import (
	"context"

	"github.com/k1nky/tookhook/pkg/plugin"
)

type LogHandler struct {
	builtinPlugin
}

func NewLogHandler(log logger) *LogHandler {
	return &LogHandler{
		builtinPlugin: builtinPlugin{
			Logger: log,
		},
	}
}

func (h *LogHandler) Validate(ctx context.Context, r plugin.Handler) error {
	return nil
}

func (h *LogHandler) Forward(ctx context.Context, r plugin.Handler, data []byte) ([]byte, error) {
	h.Logger.Infof("%s", data)
	return nil, nil
}
