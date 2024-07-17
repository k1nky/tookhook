package main

import (
	"log"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type Plugin struct{}

func (f Plugin) Validate(r plugin.Receiver) error {
	log.Println("from validate")
	return nil
}

func (f Plugin) Forward(r plugin.Receiver, data []byte) ([]byte, error) {
	target := strings.Split(r.Target, "/")
	p := NewPachca(r.Token)
	m := MessagePayload{Message: Message{
		EntityType: target[0],
		EntityId:   target[1],
		Content:    string(data),
	}}
	response, err := p.Send(m)
	log.Println(r.Target, string(response))
	return response, err
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
