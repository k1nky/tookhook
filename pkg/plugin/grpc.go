package plugin

import (
	"context"

	"github.com/k1nky/tookhook/pkg/plugin/proto"
)

type GRPCClient struct{ client proto.PluginClient }

func (m *GRPCClient) Forward(ctx context.Context, r Receiver, data []byte) ([]byte, error) {
	_, err := m.client.Forward(ctx, &proto.ForwardRequest{
		Receiver: &proto.ReceiverSpec{
			Target: r.Target,
			Token:  r.Token,
		},
		Data: &proto.Data{
			Data: data,
		},
	})
	return nil, err
}

func (m *GRPCClient) Health(ctx context.Context) error {
	_, err := m.client.Health(ctx, nil)
	return err
}

type GRPCServer struct {
	proto.UnimplementedPluginServer
	Impl Plugin
}

func (m *GRPCServer) Forward(ctx context.Context, req *proto.ForwardRequest) (*proto.ForwardResponse, error) {
	r := Receiver{
		Token:  req.Receiver.Token,
		Target: req.Receiver.Target,
	}
	data := req.Data.Data
	v, err := m.Impl.Forward(r, data)
	return &proto.ForwardResponse{Data: v}, err
}

func (m *GRPCServer) Health(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	err := m.Impl.Health()
	return nil, err
}
