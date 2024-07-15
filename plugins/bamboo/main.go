package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Plugin struct{}
type mapString map[string]string

func getDataType(data []byte) (string, string, error) {
	m := make(mapString)
	if err := json.Unmarshal(data, &m); err != nil {
		return "", "", err
	}
	if v, ok := m["deployment"]; ok {
		return "deployment", v, nil
	}
	if v, ok := m["build"]; ok {
		return "build", v, nil
	}
	return "", "", errors.New("unsupported data type")
}

func (f Plugin) enrich(r plugin.IngestEndpoint, data []byte) ([]byte, error) {
	bamboo := NewBamboo(r.Token, r.Endpoint)
	t, v, err := getDataType(data)
	if err != nil {
		return nil, err
	}
	switch t {
	case "deployment":
		deployment := make(mapString)
		if err := json.Unmarshal([]byte(v), &deployment); err != nil {
			return nil, err
		}
		return bamboo.GetDeployProjectById(deployment["deploymentProjectId"])
	}
	return nil, nil
}

func (f Plugin) Validate(r plugin.Receiver) error {
	return status.Error(codes.Unimplemented, "")
}

func (f Plugin) Forward(r plugin.Receiver, data []byte) ([]byte, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (f Plugin) Enrich(r plugin.IngestEndpoint, data []byte) ([]byte, error) {
	extraData, err := f.enrich(r, data)
	if err != nil {
		return data, err
	}
	if extraData == nil {
		return data, nil
	}
	response := make(map[string][]byte)
	response["enrich"] = extraData
	response["original"] = data
	log.Println(r.Endpoint, response)
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(response); err != nil {
		return data, err
	}
	return buf.Bytes(), err
}

func main() {
	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: map[string]hcplugin.Plugin{
			"grpc": &plugin.GRPCPlugin{Impl: &Plugin{}},
		},

		GRPCServer: hcplugin.DefaultGRPCServer,
	})
}
