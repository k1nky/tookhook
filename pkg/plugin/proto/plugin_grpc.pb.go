// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.2
// source: pkg/plugin/proto/plugin.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Plugin_Forward_FullMethodName  = "/protocol.grpc.Plugin/Forward"
	Plugin_Health_FullMethodName   = "/protocol.grpc.Plugin/Health"
	Plugin_Validate_FullMethodName = "/protocol.grpc.Plugin/Validate"
)

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginClient interface {
	Forward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*ForwardResponse, error)
	Health(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) Forward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*ForwardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ForwardResponse)
	err := c.cc.Invoke(ctx, Plugin_Forward_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Health(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, Plugin_Health_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, Plugin_Validate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
// All implementations must embed UnimplementedPluginServer
// for forward compatibility
type PluginServer interface {
	Forward(context.Context, *ForwardRequest) (*ForwardResponse, error)
	Health(context.Context, *Empty) (*Empty, error)
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
	mustEmbedUnimplementedPluginServer()
}

// UnimplementedPluginServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (UnimplementedPluginServer) Forward(context.Context, *ForwardRequest) (*ForwardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Forward not implemented")
}
func (UnimplementedPluginServer) Health(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedPluginServer) Validate(context.Context, *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedPluginServer) mustEmbedUnimplementedPluginServer() {}

// UnsafePluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServer will
// result in compilation errors.
type UnsafePluginServer interface {
	mustEmbedUnimplementedPluginServer()
}

func RegisterPluginServer(s grpc.ServiceRegistrar, srv PluginServer) {
	s.RegisterService(&Plugin_ServiceDesc, srv)
}

func _Plugin_Forward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Forward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_Forward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Forward(ctx, req.(*ForwardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Health(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_Validate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Plugin_ServiceDesc is the grpc.ServiceDesc for Plugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Plugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.grpc.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Forward",
			Handler:    _Plugin_Forward_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _Plugin_Health_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _Plugin_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/plugin/proto/plugin.proto",
}
