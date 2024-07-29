package options

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrEmptyValue = errors.New("could not be empty")
)

//go:generate easyjson options.go
//easyjson:json
type PluginOptions struct {
	Shell string   `json:"shell"`
	Args  []string `json:"args"`
}

func New(encoded []byte) (PluginOptions, error) {
	po := &PluginOptions{}
	err := json.Unmarshal(encoded, po)
	if len(po.Shell) == 0 {
		po.Shell = "/bin/sh"
		po.Args = append([]string{"-c"}, po.Args...)
	}
	return *po, err
}

func (po PluginOptions) Validate() error {
	if len(po.Args) == 0 {
		return fmt.Errorf("cmd: %w", ErrEmptyValue)
	}
	return nil
}
