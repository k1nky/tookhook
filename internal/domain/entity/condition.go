// Package entity contains domain entities for the webhook processing service.
package entity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/cel-go/cel"
)

// Condition represents a CEL-based condition for filtering chains and handlers.
// If the condition evaluates to true, the chain/handler is executed.
type Condition struct {
	// Expr is the CEL expression string.
	Expr string `yaml:"on"`
	// program is the compiled CEL program.
	program cel.Program
	// env is the CEL environment.
	env *cel.Env
}

// NewCondition creates a new condition from a CEL expression string.
func NewCondition(expr string) *Condition {
	return &Condition{
		Expr: expr,
	}
}

// Compile parses and compiles the CEL expression.
// Must be called before Match.
func (c *Condition) Compile() error {
	if c.Expr == "" {
		// Empty condition always matches
		return nil
	}

	// Create CEL environment with common variables
	env, err := cel.NewEnv(
		cel.Variable("input", cel.DynType),
		cel.Variable("payload", cel.DynType),
		cel.Variable("headers", cel.MapType(cel.StringType, cel.ListType(cel.StringType))),
		// String functions
		cel.Function("startsWith", cel.MemberOverload(
			"string_starts_with",
			[]*cel.Type{cel.StringType, cel.StringType},
			cel.BoolType,
		)),
		cel.Function("endsWith", cel.MemberOverload(
			"string_ends_with",
			[]*cel.Type{cel.StringType, cel.StringType},
			cel.BoolType,
		)),
		cel.Function("contains", cel.MemberOverload(
			"string_contains",
			[]*cel.Type{cel.StringType, cel.StringType},
			cel.BoolType,
		)),
		// Size/length function
		cel.Function("size", cel.Overload(
			"size",
			[]*cel.Type{cel.DynType},
			cel.IntType,
		)),
		// Type conversion helpers
		cel.Function("bytesToString", cel.Overload(
			"bytes_to_string",
			[]*cel.Type{cel.BytesType},
			cel.StringType,
		)),
		cel.Function("stringToBytes", cel.Overload(
			"string_to_bytes",
			[]*cel.Type{cel.StringType},
			cel.BytesType,
		)),
	)
	if err != nil {
		return fmt.Errorf("create CEL env: %w", err)
	}
	c.env = env

	// Parse the expression
	ast, issues := env.Parse(c.Expr)
	if issues != nil && issues.Err() != nil {
		return fmt.Errorf("parse CEL expression: %w", issues.Err())
	}

	// Check the expression
	if _, issues := env.Check(ast); issues != nil && issues.Err() != nil {
		return fmt.Errorf("check CEL expression: %w", issues.Err())
	}

	// Compile the program
	program, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("compile CEL expression: %w", err)
	}
	c.program = program

	return nil
}

// Match evaluates the condition against the input data.
// Returns true if the condition matches or if there's no condition.
func (c *Condition) Match(input []byte, contentType string, headers map[string][]string) (bool, error) {
	// Empty condition always matches
	if c.program == nil {
		return true, nil
	}

	// Parse input based on content type
	data := c.parseInput(input, contentType)

	// Prepare evaluation context
	ctx := map[string]any{
		"input":   data,
		"payload": data,
		"headers": headers,
	}

	// Evaluate the CEL expression
	result, _, err := c.program.Eval(ctx)
	if err != nil {
		return false, fmt.Errorf("evaluate CEL expression: %w", err)
	}

	return result.Value().(bool), nil
}

// parseInput parses the input bytes into a structured value based on content type.
func (c *Condition) parseInput(input []byte, contentType string) any {
	// Empty input
	if len(input) == 0 {
		return ""
	}

	// JSON content
	if strings.Contains(contentType, "application/json") {
		var data any
		if err := json.Unmarshal(input, &data); err == nil {
			return data
		}
	}

	// Form data
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		values, err := url.ParseQuery(string(input))
		if err == nil {
			// Convert to map[string]any for CEL
			result := make(map[string]any)
			for k, v := range values {
				if len(v) == 1 {
					result[k] = v[0]
				} else {
					result[k] = v
				}
			}
			return result
		}
	}

	// Plain text or unknown - return as string
	return string(input)
}
