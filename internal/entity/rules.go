package entity

import (
	"bytes"
	"encoding/json"
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
	Handlers []Receiver `yaml:"handlers"`
	// If true the hook will be skipped and the incoming request will be dropped.
	Disabled bool `yaml:"disabled"`
}

// Receiver is the component that will receive data from the webhook.
type Receiver struct {
	// Type is actually plugin name that will process incoming data.
	Type string `yaml:"type"`
	// Options will be passed to the plugin.
	Options map[string]interface{} `yaml:"options"`
	// List of template that will be executed before being passed to the plugin.
	Template Templates `yaml:"template"`
	// If true the receiver will be skipped.
	Disabled bool `yaml:"disabled"`
	// Serialized value of `Options`
	options []byte
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

// Validate checks the rules and returns en error if there is one.
func (r *Rules) Validate() error {
	for _, hook := range r.Hooks {
		if isEmpty(hook.Income) {
			return fmt.Errorf("income %w", ErrEmptyValue)
		}
		for i, outcome := range hook.Handlers {
			if isEmpty(outcome.Type) {
				return fmt.Errorf("outcome type %w", ErrEmptyValue)
			}
			buf := bytes.NewBuffer(nil)
			if err := json.NewEncoder(buf).Encode(outcome.Options); err != nil {
				// should not happen
				return fmt.Errorf("outcome options has %w", err)
			}
			hook.Handlers[i].options = buf.Bytes()
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
		Options: r.options,
	}
}
