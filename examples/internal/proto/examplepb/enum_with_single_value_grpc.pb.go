// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: examples/internal/proto/examplepb/enum_with_single_value.proto

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

// EnumWithSingleValueServiceClient is the client API for EnumWithSingleValueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnumWithSingleValueServiceClient interface {
	Echo(ctx context.Context, in *EnumWithSingleValueEchoRequest, opts ...grpc.CallOption) (*EnumWithSingleValueEchoResponse, error)
}

type enumWithSingleValueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEnumWithSingleValueServiceClient(cc grpc.ClientConnInterface) EnumWithSingleValueServiceClient {
	return &enumWithSingleValueServiceClient{cc}
}

func (c *enumWithSingleValueServiceClient) Echo(ctx context.Context, in *EnumWithSingleValueEchoRequest, opts ...grpc.CallOption) (*EnumWithSingleValueEchoResponse, error) {
	out := new(EnumWithSingleValueEchoResponse)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnumWithSingleValueServiceServer is the server API for EnumWithSingleValueService service.
// All implementations should embed UnimplementedEnumWithSingleValueServiceServer
// for forward compatibility
type EnumWithSingleValueServiceServer interface {
	Echo(context.Context, *EnumWithSingleValueEchoRequest) (*EnumWithSingleValueEchoResponse, error)
}

// UnimplementedEnumWithSingleValueServiceServer should be embedded to have forward compatible implementations.
type UnimplementedEnumWithSingleValueServiceServer struct {
}

func (UnimplementedEnumWithSingleValueServiceServer) Echo(context.Context, *EnumWithSingleValueEchoRequest) (*EnumWithSingleValueEchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}

// UnsafeEnumWithSingleValueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnumWithSingleValueServiceServer will
// result in compilation errors.
type UnsafeEnumWithSingleValueServiceServer interface {
	mustEmbedUnimplementedEnumWithSingleValueServiceServer()
}

func RegisterEnumWithSingleValueServiceServer(s grpc.ServiceRegistrar, srv EnumWithSingleValueServiceServer) {
	s.RegisterService(&EnumWithSingleValueService_ServiceDesc, srv)
}

func _EnumWithSingleValueService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnumWithSingleValueEchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnumWithSingleValueServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnumWithSingleValueServiceServer).Echo(ctx, req.(*EnumWithSingleValueEchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EnumWithSingleValueService_ServiceDesc is the grpc.ServiceDesc for EnumWithSingleValueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EnumWithSingleValueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueService",
	HandlerType: (*EnumWithSingleValueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EnumWithSingleValueService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/internal/proto/examplepb/enum_with_single_value.proto",
}
