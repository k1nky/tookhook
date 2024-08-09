package entity

import (
	"fmt"
	"strings"
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

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Compile checks the rules common syntax and returns en error if there is one.
func (r *Rules) Compile() (err error) {
	for _, hook := range r.Hooks {
		if isEmpty(hook.Income) {
			return fmt.Errorf("income %w", ErrEmptyValue)
		}
		for _, h := range hook.Handlers {
			if err := h.Compile(); err != nil {
				return fmt.Errorf("handler could not be compiled: %w", err)
			}
		}
	}
	return nil
}
