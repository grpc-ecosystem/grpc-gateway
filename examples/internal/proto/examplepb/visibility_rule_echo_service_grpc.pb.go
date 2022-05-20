// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package examplepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VisibilityRuleEchoServiceClient is the client API for VisibilityRuleEchoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VisibilityRuleEchoServiceClient interface {
	// Echo method receives a simple message and returns it.
	// It should always be visible in the open API output.
	Echo(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error)
	// EchoInternal is an internal API that should only be visible in the OpenAPI spec
	// if `visibility_restriction_selectors` includes "INTERNAL".
	EchoInternal(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error)
	// EchoPreview is a preview API that should only be visible in the OpenAPI spec
	// if `visibility_restriction_selectors` includes "PREVIEW".
	EchoPreview(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error)
	// EchoInternalAndPreview is a internal and preview API that should only be visible in the OpenAPI spec
	// if `visibility_restriction_selectors` includes "PREVIEW" or "INTERNAL".
	EchoInternalAndPreview(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error)
	// EchoInternalMessage method is always visible
	EchoInternalMessage(ctx context.Context, in *HiddenMessage, opts ...grpc.CallOption) (*HiddenMessage, error)
}

type visibilityRuleEchoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVisibilityRuleEchoServiceClient(cc grpc.ClientConnInterface) VisibilityRuleEchoServiceClient {
	return &visibilityRuleEchoServiceClient{cc}
}

func (c *visibilityRuleEchoServiceClient) Echo(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error) {
	out := new(VisibilityRuleSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visibilityRuleEchoServiceClient) EchoInternal(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error) {
	out := new(VisibilityRuleSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoInternal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visibilityRuleEchoServiceClient) EchoPreview(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error) {
	out := new(VisibilityRuleSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoPreview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visibilityRuleEchoServiceClient) EchoInternalAndPreview(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error) {
	out := new(VisibilityRuleSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoInternalAndPreview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visibilityRuleEchoServiceClient) EchoInternalMessage(ctx context.Context, in *HiddenMessage, opts ...grpc.CallOption) (*HiddenMessage, error) {
	out := new(HiddenMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoInternalMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VisibilityRuleEchoServiceServer is the server API for VisibilityRuleEchoService service.
// All implementations should embed UnimplementedVisibilityRuleEchoServiceServer
// for forward compatibility
type VisibilityRuleEchoServiceServer interface {
	// Echo method receives a simple message and returns it.
	// It should always be visible in the open API output.
	Echo(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error)
	// EchoInternal is an internal API that should only be visible in the OpenAPI spec
	// if `visibility_restriction_selectors` includes "INTERNAL".
	EchoInternal(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error)
	// EchoPreview is a preview API that should only be visible in the OpenAPI spec
	// if `visibility_restriction_selectors` includes "PREVIEW".
	EchoPreview(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error)
	// EchoInternalAndPreview is a internal and preview API that should only be visible in the OpenAPI spec
	// if `visibility_restriction_selectors` includes "PREVIEW" or "INTERNAL".
	EchoInternalAndPreview(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error)
	// EchoInternalMessage method is always visible
	EchoInternalMessage(context.Context, *HiddenMessage) (*HiddenMessage, error)
}

// UnimplementedVisibilityRuleEchoServiceServer should be embedded to have forward compatible implementations.
type UnimplementedVisibilityRuleEchoServiceServer struct {
}

func (UnimplementedVisibilityRuleEchoServiceServer) Echo(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedVisibilityRuleEchoServiceServer) EchoInternal(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EchoInternal not implemented")
}
func (UnimplementedVisibilityRuleEchoServiceServer) EchoPreview(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EchoPreview not implemented")
}
func (UnimplementedVisibilityRuleEchoServiceServer) EchoInternalAndPreview(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EchoInternalAndPreview not implemented")
}
func (UnimplementedVisibilityRuleEchoServiceServer) EchoInternalMessage(context.Context, *HiddenMessage) (*HiddenMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EchoInternalMessage not implemented")
}

// UnsafeVisibilityRuleEchoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VisibilityRuleEchoServiceServer will
// result in compilation errors.
type UnsafeVisibilityRuleEchoServiceServer interface {
	mustEmbedUnimplementedVisibilityRuleEchoServiceServer()
}

func RegisterVisibilityRuleEchoServiceServer(s grpc.ServiceRegistrar, srv VisibilityRuleEchoServiceServer) {
	s.RegisterService(&VisibilityRuleEchoService_ServiceDesc, srv)
}

func _VisibilityRuleEchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisibilityRuleSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisibilityRuleEchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisibilityRuleEchoServiceServer).Echo(ctx, req.(*VisibilityRuleSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisibilityRuleEchoService_EchoInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisibilityRuleSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisibilityRuleEchoServiceServer).EchoInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoInternal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisibilityRuleEchoServiceServer).EchoInternal(ctx, req.(*VisibilityRuleSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisibilityRuleEchoService_EchoPreview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisibilityRuleSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisibilityRuleEchoServiceServer).EchoPreview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoPreview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisibilityRuleEchoServiceServer).EchoPreview(ctx, req.(*VisibilityRuleSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisibilityRuleEchoService_EchoInternalAndPreview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisibilityRuleSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisibilityRuleEchoServiceServer).EchoInternalAndPreview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoInternalAndPreview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisibilityRuleEchoServiceServer).EchoInternalAndPreview(ctx, req.(*VisibilityRuleSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisibilityRuleEchoService_EchoInternalMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HiddenMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisibilityRuleEchoServiceServer).EchoInternalMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService/EchoInternalMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisibilityRuleEchoServiceServer).EchoInternalMessage(ctx, req.(*HiddenMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// VisibilityRuleEchoService_ServiceDesc is the grpc.ServiceDesc for VisibilityRuleEchoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VisibilityRuleEchoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleEchoService",
	HandlerType: (*VisibilityRuleEchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _VisibilityRuleEchoService_Echo_Handler,
		},
		{
			MethodName: "EchoInternal",
			Handler:    _VisibilityRuleEchoService_EchoInternal_Handler,
		},
		{
			MethodName: "EchoPreview",
			Handler:    _VisibilityRuleEchoService_EchoPreview_Handler,
		},
		{
			MethodName: "EchoInternalAndPreview",
			Handler:    _VisibilityRuleEchoService_EchoInternalAndPreview_Handler,
		},
		{
			MethodName: "EchoInternalMessage",
			Handler:    _VisibilityRuleEchoService_EchoInternalMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/internal/proto/examplepb/visibility_rule_echo_service.proto",
}

// VisibilityRuleInternalEchoServiceClient is the client API for VisibilityRuleInternalEchoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VisibilityRuleInternalEchoServiceClient interface {
	// Echo method receives a simple message and returns it.
	// It should not be visible in the open API output.
	Echo(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error)
}

type visibilityRuleInternalEchoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVisibilityRuleInternalEchoServiceClient(cc grpc.ClientConnInterface) VisibilityRuleInternalEchoServiceClient {
	return &visibilityRuleInternalEchoServiceClient{cc}
}

func (c *visibilityRuleInternalEchoServiceClient) Echo(ctx context.Context, in *VisibilityRuleSimpleMessage, opts ...grpc.CallOption) (*VisibilityRuleSimpleMessage, error) {
	out := new(VisibilityRuleSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleInternalEchoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VisibilityRuleInternalEchoServiceServer is the server API for VisibilityRuleInternalEchoService service.
// All implementations should embed UnimplementedVisibilityRuleInternalEchoServiceServer
// for forward compatibility
type VisibilityRuleInternalEchoServiceServer interface {
	// Echo method receives a simple message and returns it.
	// It should not be visible in the open API output.
	Echo(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error)
}

// UnimplementedVisibilityRuleInternalEchoServiceServer should be embedded to have forward compatible implementations.
type UnimplementedVisibilityRuleInternalEchoServiceServer struct {
}

func (UnimplementedVisibilityRuleInternalEchoServiceServer) Echo(context.Context, *VisibilityRuleSimpleMessage) (*VisibilityRuleSimpleMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}

// UnsafeVisibilityRuleInternalEchoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VisibilityRuleInternalEchoServiceServer will
// result in compilation errors.
type UnsafeVisibilityRuleInternalEchoServiceServer interface {
	mustEmbedUnimplementedVisibilityRuleInternalEchoServiceServer()
}

func RegisterVisibilityRuleInternalEchoServiceServer(s grpc.ServiceRegistrar, srv VisibilityRuleInternalEchoServiceServer) {
	s.RegisterService(&VisibilityRuleInternalEchoService_ServiceDesc, srv)
}

func _VisibilityRuleInternalEchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisibilityRuleSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisibilityRuleInternalEchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleInternalEchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisibilityRuleInternalEchoServiceServer).Echo(ctx, req.(*VisibilityRuleSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// VisibilityRuleInternalEchoService_ServiceDesc is the grpc.ServiceDesc for VisibilityRuleInternalEchoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VisibilityRuleInternalEchoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gateway.examples.internal.proto.examplepb.VisibilityRuleInternalEchoService",
	HandlerType: (*VisibilityRuleInternalEchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _VisibilityRuleInternalEchoService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/internal/proto/examplepb/visibility_rule_echo_service.proto",
}
