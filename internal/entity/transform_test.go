package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformationExecuteByJsonMapKeys(t *testing.T) {
	templ := `It is {{ .message }} at {{ index . "@timestamp" }} on {{ index . "host.name" }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`)
	expected := []byte("It is Select 1 at 2024-07-11T12:40:31.574Z on hostname")
	tf := &Transform{
		Template: templ,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestTransformationExecuteByJsonEmbeddedKeys(t *testing.T) {
	templ := `It is {{ .deployment.environmentName }}`
	data := []byte(`{
		"uuid" : "fe6aed0c-b672-43c9-a9d8-eb3f81215ab3",
		"timestamp" : "2024-07-10 16:59:19 +0300",
		"notification" : "Deployment Finished Notification",
		"deployment" : {
		  "status" : "Successful",
		  "environmentName": "Production Docker"
		}
	  }
	`)
	expected := []byte("It is Production Docker")
	tf := &Transform{
		Template: templ,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestTransformationExecuteByJsonLostValues(t *testing.T) {
	templ := `{{ .message }} and{{ .message2 }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", "type": "app_log"}`)
	expected := []byte("Select 1 and")
	tf := &Transform{
		Template: templ,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestTransformationExecuteByJsonInvalidJSON(t *testing.T) {
	templ := `{{ .message }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", `)
	tf := &Transform{
		Template: templ,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestTransformationExecuteByRegexp(t *testing.T) {
	data := []byte(`{"name": "Name", "data": "My Data"}`)
	expected := []byte("Got My Data")
	tf := &Transform{
		Template: "Got {{ index . 1 }}",
		RegExp:   `data\":\s*\"([^\"]+)`,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestTransformationExecuteByRegexpNotMatch(t *testing.T) {
	data := []byte(`{"name": "Name", "data": "My Data"}`)
	tf := &Transform{
		Template: "Got {{ index . 1 }}",
		RegExp:   `not_match\":\s*\"([^\"]+)`,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, data, got)
	assert.NoError(t, err)
}

func TestTransformationExecuteNotMatch(t *testing.T) {
	templ := `It is {{ .message }} at {{ index . "@timestamp" }} on {{ index . "host.name" }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`)
	tf := &Transform{
		Template: templ,
		On:       `NOT_MATCH`,
	}
	tf.Compile()
	got, err := tf.Execute(data)
	assert.Equal(t, data, got)
	assert.NoError(t, err)
}

func TestTransformCompile(t *testing.T) {
	tests := []struct {
		name string
		tf   *Transform
	}{
		{
			name: "empty template",
			tf: &Transform{
				Template: "",
			},
		},
		{
			name: "invalid on",
			tf: &Transform{
				Template: "invalid",
				On:       ")",
			},
		},
		{
			name: "invalid regexp",
			tf: &Transform{
				Template: "invalid",
				RegExp:   ")",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tf.Compile()
			assert.Error(t, got, tt.name)
		})
	}

}
