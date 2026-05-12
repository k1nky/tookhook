package builtin

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newJSHandlerLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
}

func TestJsHandler_Name(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	assert.Equal(t, "~js", handler.Name())
}

func TestJsHandler_Execute_WithValidCode(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"John","age":30}`)

	opts := map[string]any{
		"code": "result = JSON.stringify({ greeting: 'Hello', name: input.name });",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Contains(t, string(result), "Hello")
	assert.Contains(t, string(result), "John")
}

func TestJsHandler_Execute_WithNoCode(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`test input`)

	opts := map[string]any{}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "test input", string(result))
}

func TestJsHandler_Execute_WithEmptyCode(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`test input`)

	opts := map[string]any{
		"code": "",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "test input", string(result))
}

func TestJsHandler_Execute_WithInvalidCode(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`test input`)

	opts := map[string]any{
		"code": "this is not valid javascript {",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "test input", string(result))
}

func TestJsHandler_Execute_WithJSONFunctions(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"value":100}`)

	opts := map[string]any{
		"code": `
			var parsed = JSON.parse(JSON.stringify(input));
			var result = { doubled: parsed.value * 2 };
		`,
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Contains(t, string(result), "doubled")
	assert.Contains(t, string(result), "200")
}

func TestJsHandler_Execute_WithOptionsAccess(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`test`)

	opts := map[string]any{
		"code":   "result = Options.prefix + input;",
		"prefix": "Result: ",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	// Result is JSON-encoded, so strings have quotes
	assert.Contains(t, string(result), "Result: test")
}

func TestJsHandler_Execute_WithStringInput(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`hello world`)

	opts := map[string]any{
		"code": "result = 'Output: ' + input;",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	// Result is JSON-encoded
	assert.Contains(t, string(result), "Output: hello world")
}

func TestJsHandler_Execute_WithLogFunction(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`test`)

	opts := map[string]any{
		"code": "log('test message'); result = 'done';",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, `"done"`, string(result))
}

func TestJsHandler_Execute_Concurrent(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"test":"data"}`)
	opts := map[string]any{
		"code": "result = input.test;",
	}

	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			result, err := handler.Execute(ctx, input, opts)
			assert.NoError(t, err)
			// Result is JSON-encoded string
			assert.Contains(t, string(result), "data")
			done <- true
		}()
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func TestJsHandler_Execute_WithResultVariable(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"Alice"}`)

	opts := map[string]any{
		"code": `
			result = {
				salutation: "Hello " + input.name + "!",
				success: true
			};
		`,
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Contains(t, string(result), "Hello")
	assert.Contains(t, string(result), "Alice")
}

func BenchmarkJsHandler_SimpleCode(b *testing.B) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"value":42}`)
	opts := map[string]any{
		"code": "result = input.value;",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Execute(ctx, input, opts)
		assert.NoError(b, err)
	}
}

func BenchmarkJsHandler_JSONOperations(b *testing.B) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"John","items":[{"id":1},{"id":2}]}`)
	opts := map[string]any{
		"code": `
			var data = JSON.parse(JSON.stringify(input));
			result = JSON.stringify({ count: data.items.length });
		`,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Execute(ctx, input, opts)
		assert.NoError(b, err)
	}
}

func BenchmarkJsHandler_OptionsAccess(b *testing.B) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`test`)
	opts := map[string]any{
		"code":   "result = Options.format.replace('{}', input);",
		"format": "Hello {}!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Execute(ctx, input, opts)
		assert.NoError(b, err)
	}
}

func TestJsHandler_Execute_HTTPRequest(t *testing.T) {
	// This test would require a mock HTTP server or real HTTP endpoint
	// For now, we just verify the HTTP object is accessible
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{}`)

	opts := map[string]any{
		"code": `
			// Verify HTTP object exists
			result = typeof HTTP !== 'undefined';
		`,
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "true", strings.TrimSpace(string(result)))
}

func TestJsHandler_Execute_NoExplicitResult(t *testing.T) {
	logger := newJSHandlerLogger()
	handler := NewJsHandler(logger)

	ctx := context.Background()
	input := []byte(`{"key":"value"}`)

	opts := map[string]any{
		"code": `
			// No result variable set, returns null
			var x = input.key;
		`,
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	// No explicit result returns null
	assert.Equal(t, "null", string(result))
}
