// Package service provides application services.
package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/k1nky/tookhook/internal/builtin"
	"github.com/k1nky/tookhook/internal/domain/entity"
)

// Processor handles webhook task processing with pipeline execution.
type Processor struct {
	registry *builtin.Registry
	logger   *slog.Logger
}

// NewProcessor creates a new processor.
func NewProcessor(registry *builtin.Registry, logger *slog.Logger) *Processor {
	return &Processor{
		registry: registry,
		logger:   logger,
	}
}

// ProcessWebhook processes a webhook task through its endpoint chains.
func (p *Processor) ProcessWebhook(ctx context.Context, endpoint *entity.Endpoint, task *entity.WebhookTask) error {
	// Check if endpoint is disabled
	if endpoint.Disabled {
		p.logger.Debug("endpoint disabled, skipping", "endpoint", endpoint.Name)
		return nil
	}

	// Process each chain in parallel using errgroup semantics
	for chainIdx := range endpoint.Chains {
		chain := &endpoint.Chains[chainIdx]
		if chain.Disabled {
			p.logger.Debug("chain disabled, skipping",
				"endpoint", endpoint.Name,
				"chain", chainIdx,
			)
			continue
		}

		// Check chain condition
		if chain.On != nil {
			matched, err := chain.On.Match(task.Payload, task.ContentType, task.Headers)
			if err != nil {
				p.logger.Error("chain condition evaluation failed",
					"endpoint", endpoint.Name,
					"chain", chainIdx,
					"error", err,
				)
				continue
			}
			if !matched {
				p.logger.Debug("chain condition not matched, skipping",
					"endpoint", endpoint.Name,
					"chain", chainIdx,
				)
				continue
			}
		}

		// Execute chain synchronously (can be parallelized later)
		start := time.Now()
		if err := p.executeChain(ctx, chain, task); err != nil {
			p.logger.Error("chain execution failed",
				"endpoint", endpoint.Name,
				"chain", chainIdx,
				"error", err,
				"task_id", task.ID,
			)
			// Continue processing other chains even if one fails
		} else {
			p.logger.Info("chain execution finished",
				"endpoint", endpoint.Name,
				"chain", chainIdx,
				"task_id", task.ID,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		}
	}

	return nil
}

// executeChain executes a single chain of handlers sequentially (pipeline).
// func (p *Processor) executeChain(ctx context.Context, chain *entity.Chain, input []byte, contentType string, headers map[string][]string) error {
func (p *Processor) executeChain(ctx context.Context, chain *entity.Chain, task *entity.WebhookTask) error {
	data := task.Payload

	for handlerIdx := range chain.Handlers {
		h := &chain.Handlers[handlerIdx]

		if h.Disabled {
			p.logger.Debug("handler disabled, skipping",
				"type", h.Type,
			)
			continue
		}

		// Check handler condition
		if h.On != nil {
			matched, err := h.On.Match(data, task.ContentType, task.Headers)
			if err != nil {
				p.logger.Error("handler condition evaluation failed",
					"type", h.Type,
					"error", err,
				)
				return err
			}
			if !matched {
				p.logger.Debug("handler condition not matched, skipping",
					"type", h.Type,
				)
				continue
			}
		}

		// Get handler from registry
		executor, err := p.registry.Get(h.Type)
		if err != nil {
			p.logger.Error("handler not found",
				"type", h.Type,
				"error", err,
				"task_id", task.ID,
			)
			return err
		}

		// Execute handler
		output, err := executor.Execute(ctx, data, h.Options)
		if err != nil {
			p.logger.Error("handler execution failed",
				"type", h.Type,
				"error", err,
				"task_id", task.ID,
			)
			return err
		}

		// Pass output to next handler (pipeline)
		if output != nil {
			data = output
		}
	}

	return nil
}
