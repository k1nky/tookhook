package rules

import (
	"errors"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
	"github.com/k1nky/tookhook/pkg/thstrings/tstrings"
)

type Action struct {
	Pipeline Pipeline `yaml:"pipeline"`
	// On is a regexp, action will be applied if the regexp is matched.
	On tstrings.String `yaml:"on"`
}

type Actions []*Action

func (a *Action) Compile() (err error) {
	if err = a.Pipeline.Compile(); err != nil {
		return err
	}
	if err = a.On.Compile(); err != nil {
		return err
	}
	return nil
}

func (a Action) Execute(r *hooks.Hook) ([]byte, error) {
	if ok, err := a.On.Match(r); !ok || err != nil {
		return r.RawBody, fmt.Errorf("transform: %s: %w", err, entity.ErrNotMatch)
	}
	return a.Pipeline.Execute(r)
}

func (t Actions) Compile() error {
	for _, v := range t {
		if err := v.Compile(); err != nil {
			return err
		}
	}
	return nil
}

func (a Actions) Execute(r *hooks.Hook) (transformed []byte, err error) {
	transformed = r.RawBody
	for _, action := range a {
		transformed, err = action.Execute(r)
		if errors.Is(err, entity.ErrNotMatch) {
			continue
		}
		return
	}
	return
}
