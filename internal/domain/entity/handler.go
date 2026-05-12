package entity

// Handler represents a single processing step in a chain.
// It contains the handler type and its configuration options.
type Handler struct {
	// Type identifies the handler to use.
	// Built-in handlers are prefixed with ~ (e.g., ~log, ~http, ~exec).
	Type string `yaml:"type"`
	// Options contains handler-specific configuration.
	Options map[string]any `yaml:"options"`
	// Disabled indicates whether the handler is disabled.
	// If true, the handler is skipped during processing.
	Disabled bool `yaml:"disabled"`
	// On is the CEL condition for the handler.
	// If specified, the handler is only executed if the condition evaluates to true.
	On *Condition `yaml:"on"`
}

// Validate checks if the handler configuration is valid.
func (h *Handler) Validate() error {
	if h.Type == "" {
		return ErrEmptyHandlerType
	}
	if h.On != nil {
		if err := h.On.Compile(); err != nil {
			return err
		}
	}
	return nil
}
