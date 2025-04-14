package transform

import (
	"testing"

	"github.com/k1nky/tookhook/pkg/thstrings"
	"github.com/k1nky/tookhook/pkg/thstrings/tstrings"
	"github.com/stretchr/testify/assert"
)

func TestPipeline_Compile(t *testing.T) {
	tests := []struct {
		name      string
		pipeline  Pipeline
		wantError error
	}{
		{
			name:     "Empty",
			pipeline: Pipeline{},
		},
		{
			name: "Invalid stage",
			pipeline: Pipeline{
				&Stage{
					Type:              StageTypeTemplate,
					TemplateTransform: &TemplateStage{*tstrings.New(`{{ .invalid`)},
				},
			},
			wantError: thstrings.ErrCompile,
		},
		{
			name: "Valid template",
			pipeline: Pipeline{
				&Stage{
					Type:              StageTypeTemplate,
					TemplateTransform: &TemplateStage{*tstrings.New(`{{ . }}`)},
				},
			},
		},
		{
			name: "Invalid second template",
			pipeline: Pipeline{
				&Stage{
					Type:              StageTypeTemplate,
					TemplateTransform: &TemplateStage{*tstrings.New(`{{ . }}`)},
				},
				&Stage{
					Type:              StageTypeTemplate,
					TemplateTransform: &TemplateStage{*tstrings.New(`{{ . `)},
				},
			},
			wantError: thstrings.ErrCompile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pipeline.Compile()
			assert.ErrorIs(t, got, tt.wantError)
		})
	}

}

func TestPipeline_Execute(t *testing.T) {
	tests := []struct {
		name       string
		pipeline   Pipeline
		data       []byte
		wantError  error
		wantResult []byte
	}{
		{
			name:       "JSON",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("It is Select 1 at 2024-07-11T12:40:31.574Z on hostname"),
			pipeline: Pipeline{
				&Stage{
					Type: StageTypeTemplate,
					TemplateTransform: &TemplateStage{
						*tstrings.New(`It is {{ .message }} at {{ index . "@timestamp" }} on {{ index . "host.name" }}`),
					},
				},
			},
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", `),
			wantError: thstrings.ErrFailedExecution,
			pipeline: Pipeline{
				&Stage{
					Type: StageTypeTemplate,
					TemplateTransform: &TemplateStage{
						*tstrings.New(`{{ .message }}`),
					},
				},
			},
		},
		{
			name:       "Multiple templates",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("Select 1"),
			pipeline: Pipeline{
				&Stage{
					Type: StageTypeTemplate,
					TemplateTransform: &TemplateStage{
						*tstrings.New(`{"msg": "{{ .message }}"}`),
					},
				},
				&Stage{
					Type: StageTypeTemplate,
					TemplateTransform: &TemplateStage{
						*tstrings.New(`{{ .msg }}`),
					},
				},
			},
		},
		{
			name:       "Multiple templates 2",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "uri": "https://httpbin.io/json"}`),
			wantResult: []byte("Sample Slide Show"),
			pipeline: Pipeline{
				&Stage{
					Type: StageTypeHTTP,
					HTTPTransform: &HTTPStage{
						URL: tstrings.String{Template: `{{ .uri }}`},
					},
				},
				&Stage{
					Type:              StageTypeTemplate,
					TemplateTransform: &TemplateStage{*tstrings.New(`{{ .slideshow.title }}`)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pipeline.Compile()
			assert.NoError(t, err)
			got, err := tt.pipeline.Execute(tt.data)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}
