package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform_Execute(t *testing.T) {
	tests := []struct {
		name       string
		tf         *Transform
		data       []byte
		wantError  error
		wantResult []byte
	}{
		{
			name:       "JSON with Map Keys",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("It is Select 1 at 2024-07-11T12:40:31.574Z on hostname"),
			tf: &Transform{
				Template: `It is {{ .message }} at {{ index . "@timestamp" }} on {{ index . "host.name" }}`,
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
			tf: &Transform{
				Template: "It is {{ .deployment.environmentName }}",
			},
		},
		{
			name:       "JSON with not exist key",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("Select 1 and"),
			tf: &Transform{
				Template: "{{ .message }} and{{ .message2 }}",
			},
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", `),
			wantError: ErrFailedExecution,
			tf: &Transform{
				Template: "{{ .message }}",
			},
		},
		{
			name:       "Not match",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantError:  ErrNotMatch,
			wantResult: []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			tf: &Transform{
				On: "NOT_MATCH",
			},
		},
		{
			name:       "Regexp Not Match",
			data:       []byte(`{"name": "Name", "data": "My Data"}`),
			wantResult: []byte(`{"name": "Name", "data": "My Data"}`),
			tf: &Transform{
				Template: "Got {{ index . 1 }}",
				RegExp:   `not_match\":\s*\"([^\"]+)`,
			},
		},
		{
			name:       "Regexp Not Match",
			data:       []byte(`{"name": "Name", "data": "My Data"}`),
			wantResult: []byte(`Got My Data`),
			tf: &Transform{
				Template: "Got {{ index . 1 }}",
				RegExp:   `data\":\s*\"([^\"]+)`,
			},
		},
		{
			name:      "Discard",
			data:      []byte(`{"name": "Name", "data": "My Data"}`),
			wantError: ErrDiscard,
			tf: &Transform{
				Action: "discard",
			},
		},
	}
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
		name         string
		tf           *Transform
		wantError    error
		wantTemplate bool
		wantRegExp   bool
		wantOn       bool
	}{
		{
			name: "Empty template",
			tf: &Transform{
				Template: "",
			},
		},
		{
			name: "Invalid template",
			tf: &Transform{
				Template: "{{ .value",
			},
			wantError: ErrCompile,
		},
		{
			name: "Valid template",
			tf: &Transform{
				Template: "{{ .value }}",
			},
			wantTemplate: true,
		},
		{
			name: "Invalid On value",
			tf: &Transform{
				Template: "invalid",
				On:       ")",
			},
			wantError:    ErrCompile,
			wantTemplate: true,
		},
		{
			name: "Valid On value",
			tf: &Transform{
				Template: "valid",
				On:       ".*",
			},
			wantTemplate: true,
			wantOn:       true,
		},
		{
			name: "Invalid RegExp value",
			tf: &Transform{
				Template: "invalid",
				RegExp:   ")",
			},
			wantError:    ErrCompile,
			wantTemplate: true,
		},
		{
			name: "Valid RegExp value",
			tf: &Transform{
				Template: "valid",
				RegExp:   ".*",
			},
			wantTemplate: true,
			wantRegExp:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tf.Compile()
			assert.ErrorIs(t, got, tt.wantError)
			assert.Equal(t, tt.wantTemplate, tt.tf.template != nil, "unexpected template")
			assert.Equal(t, tt.wantOn, tt.tf.on != nil, "unexpected on")
			assert.Equal(t, tt.wantRegExp, tt.tf.regexp != nil, "unexpected regexp")
		})
	}

}

func TestTransforms_Compile(t *testing.T) {
	templ := `{{ .message }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`)
	tf := &Transforms{
		&Transform{
			Template: templ,
			On:       `NOT_MATCH`,
		},
		&Transform{
			Template: templ,
			On:       `INFO`,
		},
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, []byte("Select 1"), got)
	assert.NoError(t, err)
}
