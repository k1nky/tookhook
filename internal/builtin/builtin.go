package builtin

import (
	"context"
	"strings"

	"github.com/k1nky/tookhook/pkg/plugin"
)

const (
	LogHandlerName  = "~log"
	ExecHandlerName = "~exec"
)

type builtinPlugin struct {
	Logger logger
}

func IsBuiltin(name string) bool {
	return strings.HasPrefix(name, "~")
}

func NewHandler(name string, log logger) plugin.Plugin {
	switch name {
	case LogHandlerName:
		return NewLogHandler(log)
	case ExecHandlerName:
		return NewExecHandler(log)
	}
	return nil
}

func (bp *builtinPlugin) Health(ctx context.Context) error {
	return nil
}
