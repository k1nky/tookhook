// Package config provides configuration loading from YAML files.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	// Endpoints is a list of webhook endpoint configurations.
	Endpoints []EndpointConfig `yaml:"endpoints"`
}

// EndpointConfig represents a webhook endpoint configuration.
type EndpointConfig struct {
	// Name is the unique identifier for the endpoint.
	Name string `yaml:"name"`
	// Chains is a list of processing chains.
	Chains []ChainConfig `yaml:"chains"`
	// Disabled indicates whether the endpoint is disabled.
	Disabled bool `yaml:"disabled"`
}

// ChainConfig represents a processing chain configuration.
type ChainConfig struct {
	// Handlers is a list of handlers in the chain.
	Handlers []HandlerConfig `yaml:"handlers"`
	// Disabled indicates whether the chain is disabled.
	Disabled bool `yaml:"disabled"`
	// On is the CEL condition for the chain.
	// If specified, the chain is only executed if the condition evaluates to true.
	On string `yaml:"on"`
}

// HandlerConfig represents a handler configuration.
type HandlerConfig struct {
	// Type identifies the handler type (e.g., "~log", "~http").
	Type string `yaml:"type"`
	// Options contains handler-specific options.
	Options map[string]any `yaml:"options"`
	// Disabled indicates whether the handler is disabled.
	Disabled bool `yaml:"disabled"`
	// On is the CEL condition for the handler.
	// If specified, the handler is only executed if the condition evaluates to true.
	On string `yaml:"on"`
}

// Load reads configuration from a YAML file.
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
