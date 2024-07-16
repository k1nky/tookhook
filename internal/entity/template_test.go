package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateMapKeys(t *testing.T) {
	templ := `It is {{ .message }} at {{ index . "@timestamp" }} on {{ index . "host.name" }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`)
	expected := []byte("It is Select 1 at 2024-07-11T12:40:31.574Z on hostname")
	got, err := ExecuteTemplateByJson(templ, data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestTemplateEmbeddedKeys(t *testing.T) {
	templ := `It is {{ .deployment.environmentName }}`
	data := []byte(`{
		"uuid" : "fe6aed0c-b672-43c9-a9d8-eb3f81215ab3",
		"timestamp" : "2024-07-10 16:59:19 +0300",
		"notification" : "Deployment Finished Notification",
		"webhook" : {
		  "webhookTemplatedId" : "39780354",
		  "webhookTemplatedName" : "Deploy webhook"
		},
		"deployment" : {
		  "deploymentResultId" : "47061801",
		  "status" : "Successful",
		  "deploymentProjectId" : "43155482",
		  "environmentId" : "43253797",
		  "environmentName": "Production Docker",
		  "deploymentVersionId" : "46864395",
		  "deploymentVersionName" : "release-78",
	  
		  "startedAt" : "2024-07-10 16:58:29 +0300",
		  "finishedAt" : "2024-07-10 16:59:20 +0300",
		  "agentId" : "18251777",
	  
		  "triggerReason" : "Manual build",
		  "triggerSentence" : "was manually triggered by user"
		}
	  }
	`)
	expected := []byte("It is Production Docker")
	got, err := ExecuteTemplateByJson(templ, data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestLostValues(t *testing.T) {
	templ := `{{ .message }} and{{ .message2 }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", "type": "app_log"}`)
	expected := []byte("Select 1 and")
	got, err := ExecuteTemplateByJson(templ, data)
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}

func TestInvalidJSON(t *testing.T) {
	templ := `{{ .message }}`
	data := []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "message": "Select 1", `)
	got, err := ExecuteTemplateByJson(templ, data)
	assert.Nil(t, got)
	assert.Error(t, err)
}
