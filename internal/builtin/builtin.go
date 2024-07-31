package builtin

import (
	"strings"

	"github.com/k1nky/tookhook/pkg/plugin"
)

const (
	LogHandlerName  = "~log"
	ExecHandlerName = "~exec"
)

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
