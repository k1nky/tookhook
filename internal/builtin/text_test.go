package builtin

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
}

func TestTextHandler_Name(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	assert.Equal(t, "~text", handler.Name())
}

func TestTextHandler_Execute_WithTemplate(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"John","age":30}`)

	opts := map[string]any{
		"template": "Hello {{ .name }}, you are {{ .age }} years old",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "Hello John, you are 30 years old", string(result))
}

func TestTextHandler_Execute_WithNoTemplate(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"John"}`)

	opts := map[string]any{}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, `{"name":"John"}`, string(result))
}

func TestTextHandler_Execute_WithEmptyTemplate(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`test input`)

	opts := map[string]any{
		"template": "",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "test input", string(result))
}

func TestTextHandler_Execute_WithInvalidTemplate(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`test input`)

	opts := map[string]any{
		"template": "{{ .name }", // Missing closing brace
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "test input", string(result))
}

func TestTextHandler_Execute_WithInvalidData(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`invalid json {`)

	opts := map[string]any{
		"template": "Input: {{ . }}",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(result), "Input: invalid json"), "should handle invalid JSON as string")
}

func TestTextHandler_Execute_WithTextInput(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`hello world`)

	opts := map[string]any{
		"template": "Output: {{ . }}",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "Output: hello world", string(result))
}

func TestTextHandler_Execute_WithBuiltInFunctions(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"john","description":"hello world"}`)

	opts := map[string]any{
		"template": "Name: {{ .name | title }}, Description: {{ .description | toUpper }}",
	}

	result, err := handler.Execute(ctx, input, opts)

	assert.NoError(t, err)
	assert.Equal(t, "Name: John, Description: HELLO WORLD", string(result))
}

func TestTextHandler_Execute_Concurrent(t *testing.T) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"test":"data"}`)
	opts := map[string]any{
		"template": "{{ .test }}",
	}

	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			result, err := handler.Execute(ctx, input, opts)
			assert.NoError(t, err)
			assert.Equal(t, "data", string(result))
			done <- true
		}()
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func BenchmarkTextHandler_WithTemplate(b *testing.B) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"John","age":30,"city":"New York"}`)
	opts := map[string]any{
		"template": "Name: {{ .name }}, Age: {{ .age }}, City: {{ .city }}",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Execute(ctx, input, opts)
		assert.NoError(b, err)
	}
}

func BenchmarkTextHandler_WithNoTemplate(b *testing.B) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"name":"John"}`)
	opts := map[string]any{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Execute(ctx, input, opts)
		assert.NoError(b, err)
	}
}

func BenchmarkTextHandler_WithComplexTemplate(b *testing.B) {
	logger := newTestLogger()
	handler := NewTextHandler(logger)

	ctx := context.Background()
	input := []byte(`{"items":[{"name":"item1","price":10},{"name":"item2","price":20}],"total":30}`)
	opts := map[string]any{
		"template": `Items: {{ range .items }}{{ .name }}: ${{ .price }} {{ end }}, Total: {{ .total }}`,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Execute(ctx, input, opts)
		assert.NoError(b, err)
	}
}
