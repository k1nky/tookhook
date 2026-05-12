package config

import (
	"fmt"

	"github.com/k1nky/tookhook/internal/domain/entity"
)

// ToEntities converts configuration to domain entities.
func (c *Config) ToEntities() ([]entity.Endpoint, error) {
	endpoints := make([]entity.Endpoint, 0, len(c.Endpoints))

	for _, ec := range c.Endpoints {
		ep, err := ec.toEntity()
		if err != nil {
			return nil, fmt.Errorf("endpoint %s: %w", ec.Name, err)
		}
		endpoints = append(endpoints, *ep)
	}

	return endpoints, nil
}

// toEntity converts EndpointConfig to entity.Endpoint.
func (ec *EndpointConfig) toEntity() (*entity.Endpoint, error) {
	chains := make([]entity.Chain, 0, len(ec.Chains))
	for i, cc := range ec.Chains {
		chain, err := cc.toEntity()
		if err != nil {
			return nil, fmt.Errorf("chain %d: %w", i, err)
		}
		chains = append(chains, *chain)
	}

	return &entity.Endpoint{
		Name:     ec.Name,
		Chains:   chains,
		Disabled: ec.Disabled,
	}, nil
}

// toEntity converts ChainConfig to entity.Chain.
func (cc *ChainConfig) toEntity() (*entity.Chain, error) {
	handlers := make([]entity.Handler, 0, len(cc.Handlers))
	for i, hc := range cc.Handlers {
		handler, err := hc.toEntity()
		if err != nil {
			return nil, fmt.Errorf("handler %d: %w", i, err)
		}
		handlers = append(handlers, *handler)
	}

	return &entity.Chain{
		Handlers: handlers,
		Disabled: cc.Disabled,
	}, nil
}

// toEntity converts HandlerConfig to entity.Handler.
func (hc *HandlerConfig) toEntity() (*entity.Handler, error) {
	return &entity.Handler{
		Type:     hc.Type,
		Options:  hc.Options,
		Disabled: hc.Disabled,
	}, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	for _, ec := range c.Endpoints {
		if ec.Name == "" {
			return fmt.Errorf("endpoint name cannot be empty")
		}
		if len(ec.Chains) == 0 {
			return fmt.Errorf("endpoint %s: at least one chain is required", ec.Name)
		}
		for i, cc := range ec.Chains {
			if len(cc.Handlers) == 0 {
				return fmt.Errorf("endpoint %s, chain %d: at least one handler is required", ec.Name, i)
			}
			for j, hc := range cc.Handlers {
				if hc.Type == "" {
					return fmt.Errorf("endpoint %s, chain %d, handler %d: type cannot be empty", ec.Name, i, j)
				}
			}
		}
	}
	return nil
}
