package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type ExecHandler struct {
	builtinPlugin
}

//go:generate easyjson exec.go
//easyjson:json
type ExecHandlerOptions struct {
	Shell string   `json:"shell"`
	Args  []string `json:"args"`
}

func NewExecHandlerOptions(encoded []byte) (ExecHandlerOptions, error) {
	opts := &ExecHandlerOptions{}
	err := json.Unmarshal(encoded, opts)
	if len(opts.Shell) == 0 {
		opts.Shell = "/bin/sh"
		opts.Args = append([]string{"-c"}, opts.Args...)
	}
	return *opts, err
}

func (opts ExecHandlerOptions) Validate() error {
	if len(opts.Args) == 0 {
		return fmt.Errorf("cmd: %w", entity.ErrEmptyValue)
	}
	return nil
}

func NewExecHandler(log logger) *ExecHandler {
	return &ExecHandler{
		builtinPlugin: builtinPlugin{
			Logger: log,
		},
	}
}

func (h *ExecHandler) Validate(ctx context.Context, r plugin.Receiver) error {
	opts, err := NewExecHandlerOptions(r.Options)
	if err != nil {
		return err
	}
	return opts.Validate()
}

func (h *ExecHandler) Forward(ctx context.Context, r plugin.Receiver, data []byte) ([]byte, error) {
	opts, err := NewExecHandlerOptions(r.Options)
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, opts.Shell, opts.Args...)
	cmd.Env = []string{
		fmt.Sprintf("PLUGIN_EXEC_DATA=%s", data),
	}
	stdout := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	response := stdout.Bytes()
	return response, err
}
