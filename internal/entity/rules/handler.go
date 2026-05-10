package rules

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
	"github.com/k1nky/tookhook/pkg/plugin"
	"github.com/k1nky/tookhook/pkg/thstrings"
	"github.com/k1nky/tookhook/pkg/thstrings/tstrings"
)

type handler struct {
	options []byte
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
	On tstrings.String `yaml:"on"`
	// List of transformations that will be executed before being passed to the plugin.
	// The first one that matches the condition `On` is applied.
	PreActions Actions `yaml:"pre"`
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
func (h Handler) Execute(r *hooks.Hook) ([]byte, error) {
	if len(h.PreActions) == 0 {
		return r.RawBody, nil
	}
	return h.PreActions.Execute(r)
}

// Compile validates the handler definition and compiles it.
func (h *Handler) Compile() (err error) {
	if thstrings.IsEmpty(h.Type) {
		return fmt.Errorf("handler type %w", entity.ErrEmptyValue)
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
	if err = h.On.Compile(); err != nil {
		return err
	}
	// compile transformations
	if err := h.PreActions.Compile(); err != nil {
		return err
	}
	return nil
}
