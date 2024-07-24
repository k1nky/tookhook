package main

import (
	"context"
	"log"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin"
	"github.com/k1nky/tookhook/plugins/pachca/internal/options"
	"github.com/k1nky/tookhook/plugins/pachca/internal/pachca"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Plugin struct{}

func (p Plugin) Validate(ctx context.Context, r plugin.Receiver) error {
	opts, err := options.New(r.Options)
	if err != nil {
		log.Println(err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := opts.Validate(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (p Plugin) Health(ctx context.Context) error {
	return nil
}

func (p Plugin) Forward(ctx context.Context, r plugin.Receiver, data []byte) ([]byte, error) {
	opts, err := options.New(r.Options)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	chat := strings.Split(opts.Chat, "/")
	pch := pachca.NewPachca(opts.Token)
	m := pachca.MessagePayload{Message: pachca.Message{
		EntityType: chat[0],
		EntityId:   chat[1],
		Content:    string(data),
	}}
	response, err := pch.Send(m)
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
