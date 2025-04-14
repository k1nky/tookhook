package rules

import (
	"errors"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/pipeline/transform"
	"github.com/k1nky/tookhook/pkg/thstrings/restrings"
)

type ActionType = string

const (
	ActionTypeTransform ActionType = "transform"
	ActionTypeDiscard   ActionType = "discard"
)

type Action struct {
	Type       ActionType         `yaml:"type"`
	Transforms transform.Pipeline `yaml:"transforms"`
	// On is a regexp, transformation will be applied if the regexp is matched.
	On restrings.String `yaml:"on"`
}

type Actions []*Action

func (a *Action) Compile() (err error) {
	if a.Type == "" {
		a.Type = ActionTypeTransform
	}
	if a.Type == ActionTypeTransform {
		if err = a.Transforms.Compile(); err != nil {
			return err
		}
	}
	if err = a.On.Compile(); err != nil {
		return err
	}
	return nil
}

func (a Action) Execute(data []byte) ([]byte, error) {
	if ok := a.On.Match(data); !ok {
		return data, fmt.Errorf("transform: %w", entity.ErrNotMatch)
	}
	if a.Type == ActionTypeDiscard {
		return nil, fmt.Errorf("transform: %w", entity.ErrDiscard)
	}
	return a.Transforms.Execute(data)
}

func (t Actions) Compile() error {
	for _, v := range t {
		if err := v.Compile(); err != nil {
			return err
		}
	}
	return nil
}

func (a Actions) Execute(data []byte) (transformed []byte, err error) {
	transformed = data
	for _, t := range a {
		transformed, err = t.Execute(data)
		if errors.Is(err, entity.ErrNotMatch) {
			continue
		}
		return
	}
	return
}
