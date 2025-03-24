package entity

import (
	"errors"
	"fmt"
	"regexp"
)

type ActionType = string

const (
	ActionTypeTransform ActionType = "transform"
	ActionTypeDiscard   ActionType = "discard"
)

type action struct {
	on *regexp.Regexp
}

type Action struct {
	action
	Type       ActionType `yaml:"type"`
	Transforms Transforms `yaml:"transforms"`
	// On is a regexp, transformation will be applied if the regexp is matched.
	On string `yaml:"on"`
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
	if a.action.on, err = compileRegExp(a.On); err != nil {
		return fmt.Errorf("invalid on value: %w %w", err, ErrCompile)
	}
	return nil
}

func (a Action) Execute(data []byte) ([]byte, error) {
	if a.on != nil {
		if ok := a.on.Match(data); !ok {
			return data, fmt.Errorf("transform: %w", ErrNotMatch)
		}
	}
	if a.Type == ActionTypeDiscard {
		return nil, fmt.Errorf("transform: %w", ErrDiscard)
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
		if errors.Is(err, ErrNotMatch) {
			continue
		}
		return
	}
	return
}
