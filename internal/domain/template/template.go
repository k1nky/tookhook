// Package template provides text template processing with data binding.
package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// String represents a template string that can be executed with data.
type String struct {
	template *template.Template
	// Template is the raw template string.
	Template string
}

// builtinFuncs returns template functions available in templates.
func builtinFuncs() template.FuncMap {
	return template.FuncMap{
		"title":     cases.Title(language.Und).String,
		"toUpper":   strings.ToUpper,
		"toLower":   strings.ToLower,
		"trimSpace": strings.TrimSpace,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"match": regexp.MatchString,
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"reFindAll": func(pattern string, text string) []string {
			re := regexp.MustCompile(pattern)
			found := re.FindAllStringSubmatch(text, -1)
			if len(found) == 0 {
				return nil
			}
			return found[0]
		},
		"date": func(fmt string, t time.Time) string {
			return t.Format(fmt)
		},
		"tz": func(name string, t time.Time) (time.Time, error) {
			loc, err := time.LoadLocation(name)
			if err != nil {
				return time.Time{}, err
			}
			return t.In(loc), nil
		},
		"since": time.Since,
		"now":   time.Now,
	}
}

// New creates a new template string.
func New(s string) *String {
	return &String{
		Template: s,
	}
}

// Compile parses the template string and prepares it for execution.
func (t *String) Compile() error {
	if t.Template == "" {
		return nil
	}
	templ := template.New("")
	templ.Funcs(builtinFuncs())
	var err error
	t.template, err = templ.Parse(t.Template)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}
	return nil
}

// unmarshalData attempts to parse data as JSON, falling back to string.
func (t *String) unmarshalData(data []byte) (any, error) {
	// Try to parse as JSON object
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err == nil {
		return obj, nil
	}
	// Try to parse as JSON array
	var arr []any
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}
	// Fall back to raw string
	return string(data), nil
}

// Execute processes the template with the given data and returns the result.
func (t *String) Execute(data []byte) (string, error) {
	if t.template == nil {
		return string(data), nil
	}
	d, err := t.unmarshalData(data)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	if err := t.template.Execute(buf, d); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}
	return buf.String(), nil
}

// IsEmpty returns true if the template is empty.
func (t *String) IsEmpty() bool {
	return t.Template == ""
}
