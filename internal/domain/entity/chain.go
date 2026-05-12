package entity

// Chain represents a sequence of handlers that process data sequentially.
// Each handler in the chain receives the output of the previous handler (pipeline pattern).
type Chain struct {
	// Handlers is a list of handlers to execute sequentially.
	// The first handler receives the original webhook data,
	// and each subsequent handler receives the output of the previous one.
	Handlers []Handler `yaml:"handlers"`
	// Disabled indicates whether the chain is disabled.
	// If true, the chain is skipped during processing.
	Disabled bool `yaml:"disabled"`
	// On is the CEL condition for the chain.
	// If specified, the chain is only executed if the condition evaluates to true.
	On *Condition `yaml:"on"`
}

// Validate checks if the chain configuration is valid.
func (c *Chain) Validate() error {
	if len(c.Handlers) == 0 {
		return ErrEmptyChain
	}
	if c.On != nil {
		if err := c.On.Compile(); err != nil {
			return err
		}
	}
	for i := range c.Handlers {
		if err := c.Handlers[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}
