package rules

import (
	"encoding/json"
	"testing"

	"github.com/k1nky/tookhook/internal/entity/pipeline/transform"
	"github.com/k1nky/tookhook/pkg/thstrings/restrings"
	"github.com/k1nky/tookhook/pkg/thstrings/tstrings"
	"github.com/stretchr/testify/assert"
)

func TestHandlerContentWithoutTransform(t *testing.T) {
	h := Handler{}
	data := []byte(`some data`)
	got, err := h.Content(data)
	assert.NoError(t, err)
	assert.Equal(t, data, got)
}

func TestHandlerContentWithTransform(t *testing.T) {
	h := Handler{
		PreActions: Actions{
			&Action{
				Transforms: transform.Pipeline{
					&transform.Stage{
						Type:              transform.StageTypeTemplate,
						TemplateTransform: &transform.TemplateStage{String: *tstrings.New(`{{ .data }}`)},
					},
				},
			},
		},
	}
	err := h.PreActions.Compile()
	assert.NoError(t, err)
	got, err := h.Content([]byte(`{"data": 123}`))
	assert.NoError(t, err)
	assert.Equal(t, []byte("123"), got)
}

func TestHandlerCompileWithType(t *testing.T) {
	h := &Handler{
		Type: "handler1",
	}
	err := h.Compile()
	assert.NoError(t, err)
	assert.NotNil(t, h.On)
	assert.NotNil(t, h.options)
}

func TestHandlerCompileBadTypeValue(t *testing.T) {
	h := &Handler{
		Type: "",
	}
	err := h.Compile()
	assert.Error(t, err)
}

func TestHandlerCompileBadOnValue(t *testing.T) {
	h := &Handler{
		Type: "handler1",
		On:   *restrings.New("(abc"),
	}
	err := h.Compile()
	assert.Error(t, err)
}

func TestHandlerCompileWithOptions(t *testing.T) {
	h := &Handler{
		Type: "handler1",
		Options: map[string]interface{}{
			"option1": "value1",
		},
	}
	err := h.Compile()
	assert.NoError(t, err)
	assert.NotNil(t, h.options)
	options := make(map[string]interface{})
	err = json.Unmarshal(h.options, &options)
	assert.NoError(t, err)
}
