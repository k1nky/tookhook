package entity

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/k1nky/tookhook/pkg/plugin"
)

// Hook is the hook specification.
type Hook struct {
	// Incoming webhook request name.
	Income string `yaml:"income"`
	// List of receivers.
	Outcome []Receiver `yaml:"outcome"`
	// If true the hook will be skipped and the incoming request will be dropped.
	Disabled bool `yaml:"disabled"`
}

// Receiver is the component that will receive data from the webhook.
type Receiver struct {
	// Type is actually plugin name that will process incoming data.
	// TODO: rename to `plugin`
	Type    string            `yaml:"type"`
	Options map[string]string `yaml:"options"`
	// List of template that will be executed before being passed to the plugin.
	Template Templates `yaml:"template"`
	// If true the receiver will be skipped.
	Disabled bool `yaml:"disabled"`
}

// Rules define how to process incoming webhooks.
type Rules struct {
	// Hooks are a list of rules by which webhooks will be processed.
	Hooks []Hook `yaml:"hooks"`
	// Templates is named collection of templates. Нook rules can reference them.
	Templates map[string]Templates `yaml:"templates"`
}

type Template struct {
	RegExp   string `yaml:"regexp"`
	Template string `yaml:"template"`
	On       string `yaml:"on"`
}

type Templates []Template

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func (r *Rules) Validate() error {
	for _, hook := range r.Hooks {
		if isEmpty(hook.Income) {
			return fmt.Errorf("income %w", ErrEmptyValue)
		}
		for _, outcome := range hook.Outcome {
			if isEmpty(outcome.Type) {
				return fmt.Errorf("outcome type %w", ErrEmptyValue)
			}
		}
	}
	return nil
}

func (t Templates) Execute(data []byte) ([]byte, error) {
	for _, t := range t {
		if ok, _ := regexp.Match(t.On, data); ok {
			if !isEmpty(t.RegExp) {
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

func (r Receiver) AsPluginReceiver() plugin.Receiver {
	return plugin.Receiver{
		Options: r.Options,
	}
}
