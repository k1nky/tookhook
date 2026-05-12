// Package entity contains domain entities for the webhook processing service.
package entity

// Endpoint represents a webhook endpoint with its processing chains.
// Each endpoint has a unique name and can have multiple processing chains
// that are executed in parallel when a webhook is received.
type Endpoint struct {
	// Name is the unique identifier for the endpoint.
	// Webhooks are received at /hook/{name}
	Name string `yaml:"name"`
	// Chains is a list of processing chains.
	// Each chain is executed independently in parallel.
	Chains []Chain `yaml:"chains"`
	// Disabled indicates whether the endpoint is disabled.
	// If true, the endpoint is skipped and 404 is returned.
	Disabled bool `yaml:"disabled"`
}

// Validate checks if the endpoint configuration is valid.
func (e *Endpoint) Validate() error {
	if e.Name == "" {
		return ErrEmptyEndpointName
	}
	for i := range e.Chains {
		if err := e.Chains[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}
