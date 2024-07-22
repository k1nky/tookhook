package main

import (
	"log"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type Plugin struct{}

func (p Plugin) Validate(r plugin.Receiver) error {
	log.Println("from validate")
	return nil
}

func (p Plugin) Health() error {
	return nil
}

func (p Plugin) Forward(r plugin.Receiver, data []byte) ([]byte, error) {
	token := r.Options["token"]
	chat := r.Options["chat"]
	target := strings.Split(chat, "/")
	pachca := NewPachca(token)
	m := MessagePayload{Message: Message{
		EntityType: target[0],
		EntityId:   target[1],
		Content:    string(data),
	}}
	response, err := pachca.Send(m)
	log.Println(chat, string(response))
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
