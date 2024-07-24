package plugin

import (
	"context"

	"github.com/k1nky/tookhook/pkg/plugin/proto"
)

type GRPCClient struct{ client proto.PluginClient }

func (m *GRPCClient) Forward(ctx context.Context, r Receiver, data []byte) ([]byte, error) {
	_, err := m.client.Forward(ctx, &proto.ForwardRequest{
		Receiver: &proto.ReceiverSpec{
			Options: &proto.PluginOptions{
				Value: r.Options,
			},
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

func (m *GRPCClient) Validate(ctx context.Context, r Receiver) error {
	_, err := m.client.Validate(ctx, &proto.ValidateRequest{
		PluginOptions: &proto.PluginOptions{
			Value: r.Options,
		},
	})
	return err
}

type GRPCServer struct {
	proto.UnimplementedPluginServer
	Impl Plugin
}

func (m *GRPCServer) Forward(ctx context.Context, req *proto.ForwardRequest) (*proto.ForwardResponse, error) {
	r := Receiver{
		Options: req.Receiver.Options.Value,
	}
	data := req.Data.Data
	v, err := m.Impl.Forward(ctx, r, data)
	return &proto.ForwardResponse{Data: v}, err
}

func (m *GRPCServer) Health(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	err := m.Impl.Health(ctx)
	return nil, err
}

func (m *GRPCServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	r := Receiver{
		Options: req.PluginOptions.Value,
	}
	err := m.Impl.Validate(ctx, r)
	return nil, err
}
