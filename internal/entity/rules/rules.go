package rules

import (
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/thstrings"
)

// Endpoint is the hook specification.
type Endpoint struct {
	// Incoming webhook request name.
	Name string `yaml:"name"`
	// List of handlers.
	Handlers []*Handler `yaml:"handlers"`
	// If true the hook will be skipped and the incoming request will be dropped.
	Disabled bool `yaml:"disabled"`
}

// Rules define how to process incoming webhooks.
type Rules struct {
	// Endpoints are a list of rules by which webhooks will be processed.
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Compile checks the rules common syntax and returns en error if there is one.
func (r *Rules) Compile() (err error) {
	for _, hook := range r.Endpoints {
		if thstrings.IsEmpty(hook.Name) {
			return fmt.Errorf("name %w", entity.ErrEmptyValue)
		}
		for _, h := range hook.Handlers {
			if err := h.Compile(); err != nil {
				return fmt.Errorf("handler could not be compiled: %w", err)
			}
		}
	}
	return nil
}
