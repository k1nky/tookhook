package plugin

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/plugin/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.PluginClient }

func (m *GRPCClient) Forward(ctx context.Context, r entity.Receiver, data []byte) ([]byte, error) {
	_, err := m.client.Forward(ctx, &proto.ForwardRequest{
		Receiver: &proto.ReceiverSpec{
			Target:   r.Target,
			Token:    r.Token,
			Template: r.Template,
		},
		Data: &proto.Data{
			Data: data,
		},
	})
	return nil, err
}

func (m *GRPCClient) Validate(ctx context.Context, r entity.Receiver) error {
	_, err := m.client.Validate(ctx, &proto.ValidateRequest{
		Receiver: &proto.ReceiverSpec{
			Target:   r.Target,
			Token:    r.Token,
			Template: r.Template,
		},
	})
	return err
}

func (m *GRPCClient) Enrich(ctx context.Context, r entity.Ingest, data []byte) ([]byte, error) {
	_, err := m.client.Enrich(ctx, &proto.EnrichRequest{
		Ingest: &proto.IngestSpec{
			Endpoint: r.Endpoint,
			Token:    r.Token,
		},
		Data: &proto.Data{
			Data: data,
		},
	})
	return nil, err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	proto.UnimplementedPluginServer
	Impl Plugin
}

func (m *GRPCServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.Empty, error) {
	r := Receiver{
		Token:    req.Receiver.Token,
		Target:   req.Receiver.Target,
		Template: req.Receiver.Template,
	}
	return &proto.Empty{}, m.Impl.Validate(r)
}

func (m *GRPCServer) Forward(ctx context.Context, req *proto.ForwardRequest) (*proto.ForwardResponse, error) {
	r := Receiver{
		Token:    req.Receiver.Token,
		Target:   req.Receiver.Target,
		Template: req.Receiver.Template,
	}
	data := req.Data.Data
	v, err := m.Impl.Forward(r, data)
	return &proto.ForwardResponse{Data: v}, err
}

func (m *GRPCServer) Enrich(ctx context.Context, req *proto.EnrichRequest) (*proto.EnrichResponse, error) {
	r := IngestEndpoint{
		Endpoint: req.Ingest.Endpoint,
		Token:    req.Ingest.Token,
	}
	data := req.Data.Data
	v, err := m.Impl.Enrich(r, data)
	return &proto.EnrichResponse{Data: v}, err
}
