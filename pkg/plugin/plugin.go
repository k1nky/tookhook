package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin/proto"
	"google.golang.org/grpc"
)

type Handler struct {
	Options []byte
}

type Plugin interface {
	Forward(ctx context.Context, h Handler, data []byte) ([]byte, error)
	Health(ctx context.Context) error
	Validate(ctx context.Context, h Handler) error
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TOOKHOOK_PLUGIN",
	MagicCookieValue: "hello",
}

var PluginMap = map[string]plugin.Plugin{
	"grpc": &GRPCPlugin{},
}

type GRPCPlugin struct {
	plugin.Plugin
	Impl Plugin
}

func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewPluginClient(c)}, nil
}
