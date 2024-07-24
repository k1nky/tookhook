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
	Chat  string `json:"chat"`
	Token string `json:"token"`
}

func New(encoded []byte) (PluginOptions, error) {
	po := &PluginOptions{}
	err := json.Unmarshal(encoded, po)
	return *po, err
}

func (po PluginOptions) Validate() error {
	if len(po.Chat) == 0 {
		return fmt.Errorf("chat: %w", ErrEmptyValue)
	}
	if len(po.Token) == 0 {
		return fmt.Errorf("token: %w", ErrEmptyValue)
	}
	return nil
}
