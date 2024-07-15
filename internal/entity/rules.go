package entity

import (
	"fmt"
	"strings"
)

const (
	ReceiverTypePachca = "pachca"
	ReceiverTypeNull   = "null"
)

type Rules struct {
	Hooks []Hook `yaml:"hooks"`
}

type Hook struct {
	Income   string     `yaml:"income"`
	Ingest   Ingest     `yaml:"ingest"`
	Outcome  []Receiver `yaml:"outcome"`
	Disabled bool       `yaml:"disabled"`
}

type Ingest struct {
	Type     string `yaml:"type"`
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
}

type Receiver struct {
	Type     string `yaml:"type"`
	Token    string `yaml:"token"`
	Target   string `yaml:"target"`
	Template string `yaml:"template"`
	Disabled bool   `yaml:"disabled"`
}

func trimLength(s string) int {
	return len(strings.TrimSpace(s))
}

func (r *Rules) Validate() error {
	for _, hook := range r.Hooks {
		if trimLength(hook.Income) == 0 {
			return fmt.Errorf("income %w", ErrEmptyValue)
		}
		for _, outcome := range hook.Outcome {
			if trimLength(outcome.Type) == 0 {
				return fmt.Errorf("outcome type %w", ErrEmptyValue)
			}
		}
	}
	return nil
}

func (r Receiver) Content(data []byte) ([]byte, error) {
	if trimLength(r.Template) == 0 {
		return data, nil
	}
	return ExecuteTemplate(r.Template, data)
}
