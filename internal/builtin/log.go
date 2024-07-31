package builtin

import (
	"context"

	"github.com/k1nky/tookhook/pkg/plugin"
)

type LogHandler struct {
	Logger logger
}

func NewLogHandler(log logger) *LogHandler {
	return &LogHandler{
		Logger: log,
	}
}

func (h *LogHandler) Validate(ctx context.Context, r plugin.Receiver) error {
	return nil
}

func (h *LogHandler) Forward(ctx context.Context, r plugin.Receiver, data []byte) ([]byte, error) {
	h.Logger.Infof("%s", data)
	return nil, nil
}

func (h *LogHandler) Health(ctx context.Context) error {
	return nil
}
