package tstrings

import (
	"testing"

	"github.com/k1nky/tookhook/pkg/thstrings"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestString_Compile(t *testing.T) {
	tests := []struct {
		name      string
		s         *String
		wantError error
	}{
		{
			name: "Empty template",
			s: &String{
				Template: "",
			},
		},
		{
			name: "Invalid template",
			s: &String{
				Template: "{{ .value",
			},
			wantError: thstrings.ErrCompile,
		},
		{
			name: "Valid template",
			s: &String{
				Template: "{{ .value }}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Compile()
			assert.ErrorIs(t, got, tt.wantError)
		})
	}

}

func TestString_Execute(t *testing.T) {
	tests := []struct {
		name       string
		s          *String
		data       []byte
		wantError  error
		wantResult []byte
	}{
		{
			name:       "JSON with Map Keys",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("It is Select 1 at 2024-07-11T12:40:31.574Z on hostname"),
			s: &String{
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
			s: &String{
				Template: "It is {{ .deployment.environmentName }}",
			},
		},
		{
			name:       "JSON with not exist key",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", "type": "app_log"}`),
			wantResult: []byte("Select 1 and <no value>"),
			s: &String{
				Template: "{{ .message }} and {{ .message2 }}",
			},
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", `),
			wantError: thstrings.ErrFailedExecution,
			s: &String{
				Template: "{{ .message }}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Compile()
			got, err := tt.s.Execute(tt.data)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestString_MarshalYAML(t *testing.T) {
	s := String{Template: "abc"}
	got, err := yaml.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, []byte("abc\n"), got)
}

func TestString_UnmarshalYAML(t *testing.T) {
	s := String{}
	err := yaml.Unmarshal([]byte("abc"), &s)
	assert.NoError(t, err)
	assert.Equal(t, String{Template: "abc"}, s)
}

func TestString_BuiltinFuncs(t *testing.T) {
	tests := []struct {
		name       string
		s          *String
		data       any
		wantError  error
		wantResult string
	}{
		{
			name:       "title",
			s:          New(`Hello {{ . | title }}`),
			data:       `andrew`,
			wantError:  nil,
			wantResult: "Hello Andrew",
		},
		{
			name:       "reFindAll",
			s:          New(`{{ $a := (. | reFindAll "name\":\\s*\"([^\"]+)") }}Hello {{ index $a 1 }}`),
			data:       `"name": "Name", "data": "My Data"`,
			wantError:  nil,
			wantResult: "Hello Name",
		},
		{
			name:       "json__1",
			s:          New(`{{ $a := (. | json) }}Hello {{ $a.name | title }}`),
			data:       `{"name": "name", "data": "My Data"}`,
			wantError:  nil,
			wantResult: "Hello Name",
		},
		{
			name:       "json__2",
			s:          New(`Hello {{ (. | json).name }}`),
			data:       `{"name": "Name", "data": "My Data"}`,
			wantError:  nil,
			wantResult: "Hello Name",
		},
	}
	for _, tt := range tests {
		err := tt.s.Compile()
		assert.NoError(t, err)
		result, err := tt.s.Execute(tt.data)
		print(string(result))
		assert.ErrorIs(t, err, tt.wantError, tt.name)
		assert.Equal(t, tt.wantResult, string(result), tt.name)
	}
}
