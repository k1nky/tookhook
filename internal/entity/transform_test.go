package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransforms_Execute(t *testing.T) {
	tests := []struct {
		name       string
		tf         Transforms
		data       []byte
		wantError  error
		wantResult []byte
	}{
		{
			name:       "JSON with Map Keys",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("It is Select 1 at 2024-07-11T12:40:31.574Z on hostname"),
			tf: Transforms{
				&Transform{
					Template: `It is {{ .message }} at {{ index . "@timestamp" }} on {{ index . "host.name" }}`,
				},
			},
		},
		{
			name: "JSON with Embedded Keys",
			data: []byte(`{
				"uuid" : "fe6aed0c-b672-43c9-a9d8-eb3f81215ab3",
				"timestamp" : "2024-07-10 16:59:19 +0300",
				"notification" : "Deployment Finished Notification",
				"deployment" : {
				  "status" : "Successful",
				  "environmentName": "Production Docker"
				}
			  }
			`),
			wantResult: []byte("It is Production Docker"),
			tf: Transforms{
				&Transform{
					Template: "It is {{ .deployment.environmentName }}",
				},
			},
		},
		{
			name:       "JSON with not exist key",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("Select 1 and <no value>"),
			tf: Transforms{
				&Transform{
					Template: "{{ .message }} and {{ .message2 }}",
				},
			},
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", `),
			wantError: ErrFailedExecution,
			tf: Transforms{
				&Transform{
					Template: "{{ .message }}",
				},
			},
		},
		{
			name:       "Regexp Not Match",
			data:       []byte(`{"name": "Name", "data": "My Data"}`),
			wantResult: []byte(`{"name": "Name", "data": "My Data"}`),
			tf: Transforms{
				&Transform{
					Template: "Got {{ index . 1 }}",
					RegExp:   `not_match\":\s*\"([^\"]+)`,
				},
			},
		},
		{
			name:       "Regexp Not Match",
			data:       []byte(`{"name": "Name", "data": "My Data"}`),
			wantResult: []byte(`Got My Data`),
			tf: Transforms{
				&Transform{
					Template: "Got {{ index . 1 }}",
					RegExp:   `data\":\s*\"([^\"]+)`,
				},
			},
		},
		{
			name:       "Multiple templates",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("Select 1"),
			tf: Transforms{
				&Transform{
					Template: `{"msg": "{{ .message }}"}`,
				},
				&Transform{
					Template: `{{ .msg }}`,
				},
			},
		},
		{
			name:       "Multiple templates 2",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "uri": "https://httpbin.io/json"}`),
			wantResult: []byte("Sample Slide Show"),
			tf: Transforms{
				&Transform{
					Template: `{{ httpGet .uri }}`,
				},
				&Transform{
					Template: `{{ .slideshow.title }}`,
				},
			},
		},
	}
	//
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tf.Compile()
			got, err := tt.tf.Execute(tt.data)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestTransform_Compile(t *testing.T) {
	tests := []struct {
		name      string
		tf        Transforms
		wantError error
	}{
		{
			name: "Empty template",
			tf: Transforms{
				&Transform{
					Template: "",
				},
			},
		},
		{
			name: "Invalid template",
			tf: Transforms{
				&Transform{
					Template: "{{ .value",
				},
			},
			wantError: ErrCompile,
		},
		{
			name: "Valid template",
			tf: Transforms{
				&Transform{
					Template: "{{ .value }}",
				},
			},
		},
		{
			name: "Invalid second template",
			tf: Transforms{
				&Transform{
					Template: "{{ .value }}",
				},
				&Transform{
					Template: "{{ .value",
				},
			},
			wantError: ErrCompile,
		},
		{
			name: "Invalid RegExp value",
			tf: Transforms{
				&Transform{
					Template: "invalid",
					RegExp:   ")",
				},
			},
			wantError: ErrCompile,
		},
		{
			name: "Valid RegExp value",
			tf: Transforms{
				&Transform{
					Template: "valid",
					RegExp:   ".*",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tf.Compile()
			assert.ErrorIs(t, got, tt.wantError)
		})
	}

}
