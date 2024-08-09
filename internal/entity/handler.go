package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/k1nky/tookhook/pkg/plugin"
)

type handler struct {
	options []byte
	on      *regexp.Regexp
}

// Handler is the component that will receive data from the webhook.
type Handler struct {
	handler
	// Type is actually plugin name that will process incoming data.
	Type string `yaml:"type"`
	// Options will be passed to the plugin.
	Options map[string]interface{} `yaml:"options"`
	// On contains a regular expression string. The data will be passed to the receiver
	// if the regexp matches.
	On string `yaml:"on"`
	// List of transformations that will be executed before being passed to the plugin.
	// The first one that matches the condition `On` is applied.
	PreTransform Transforms `yaml:"pre"`
	// If true the handler will be skipped.
	Disabled bool `yaml:"disabled"`
}

// AsPluginHandler returns plugin.Handler instance.
func (h Handler) AsPluginHandler() plugin.Handler {
	return plugin.Handler{
		Options: h.handler.options,
	}
}

// Content applies transformations and returns processed data.
// The handler must be pre-compiled by `Compile`.
func (h Handler) Content(data []byte) ([]byte, error) {
	if len(h.PreTransform) == 0 {
		return data, nil
	}
	return h.PreTransform.Execute(data)
}

// Compile validates the handler definition and compiles it.
func (h *Handler) Compile() (err error) {
	if isEmpty(h.Type) {
		return fmt.Errorf("handler type %w", ErrEmptyValue)
	}

	// serialize the options to a json string
	// because the options are changed only on reload and always passed to a plugin
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(h.Options); err != nil {
		// should not happen
		return err
	}
	h.handler.options = buf.Bytes()
	// compile `on` condition
	// TODO: if On is empty
	if h.handler.on, err = regexp.Compile(h.On); err != nil {
		return err
	}
	// compile transformations
	if err := h.PreTransform.Compile(); err != nil {
		return err
	}
	return nil
}

// Match returns true if the handler should be called on the data.
func (h Handler) Match(data []byte) bool {
	if h.handler.on == nil {
		return true
	}
	return h.handler.on.Match(data)
}
