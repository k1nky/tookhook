package entity

import (
	"fmt"
	"regexp"
	"strings"
)

type Hook struct {
	Income   string     `yaml:"income"`
	Outcome  []Receiver `yaml:"outcome"`
	Disabled bool       `yaml:"disabled"`
}

type Receiver struct {
	Type     string    `yaml:"type"`
	Token    string    `yaml:"token"`
	Target   string    `yaml:"target"`
	Template Templates `yaml:"template"`
	Disabled bool      `yaml:"disabled"`
}

type Rules struct {
	Hooks     []Hook               `yaml:"hooks"`
	Templates map[string]Templates `yaml:"templates"`
}

type Template struct {
	RegExp   string `yaml:"regexp"`
	Template string `yaml:"template"`
	On       string `yaml:"on"`
}

type Templates []Template

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

func (t Templates) Execute(data []byte) ([]byte, error) {
	for _, t := range t {
		if ok, _ := regexp.Match(t.On, data); ok {
			if t.RegExp != "" {
				re, err := regexp.Compile(t.RegExp)
				if err != nil {
					return data, err
				}
				found := re.FindAllStringSubmatch(string(data), -1)
				if len(found) > 0 {
					return ExecuteTemplate(t.Template, found[0])
				}
				return data, nil
			}
			return ExecuteTemplateByJson(t.Template, data)
		}
	}
	return data, nil
}

func (r Receiver) Content(data []byte) ([]byte, error) {
	if len(r.Template) == 0 {
		return data, nil
	}
	return r.Template.Execute(data)
}
