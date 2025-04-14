package rules

import (
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/thstrings"
)

// Hook is the hook specification.
type Hook struct {
	// Incoming webhook request name.
	Income string `yaml:"income"`
	// List of handlers.
	Handlers []*Handler `yaml:"handlers"`
	// If true the hook will be skipped and the incoming request will be dropped.
	Disabled bool `yaml:"disabled"`
}

// Rules define how to process incoming webhooks.
type Rules struct {
	// Hooks are a list of rules by which webhooks will be processed.
	Hooks []Hook `yaml:"hooks"`
}

// Compile checks the rules common syntax and returns en error if there is one.
func (r *Rules) Compile() (err error) {
	for _, hook := range r.Hooks {
		if thstrings.IsEmpty(hook.Income) {
			return fmt.Errorf("income %w", entity.ErrEmptyValue)
		}
		for _, h := range hook.Handlers {
			if err := h.Compile(); err != nil {
				return fmt.Errorf("handler could not be compiled: %w", err)
			}
		}
	}
	return nil
}
