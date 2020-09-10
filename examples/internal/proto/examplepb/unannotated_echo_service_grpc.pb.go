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
const _ = grpc.SupportPackageIsVersion7

// UnannotatedEchoServiceClient is the client API for UnannotatedEchoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UnannotatedEchoServiceClient interface {
	// Echo method receives a simple message and returns it.
	//
	// The message posted as the id parameter will also be
	// returned.
	Echo(ctx context.Context, in *UnannotatedSimpleMessage, opts ...grpc.CallOption) (*UnannotatedSimpleMessage, error)
	// EchoBody method receives a simple message and returns it.
	EchoBody(ctx context.Context, in *UnannotatedSimpleMessage, opts ...grpc.CallOption) (*UnannotatedSimpleMessage, error)
	// EchoDelete method receives a simple message and returns it.
	EchoDelete(ctx context.Context, in *UnannotatedSimpleMessage, opts ...grpc.CallOption) (*UnannotatedSimpleMessage, error)
}

type unannotatedEchoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUnannotatedEchoServiceClient(cc grpc.ClientConnInterface) UnannotatedEchoServiceClient {
	return &unannotatedEchoServiceClient{cc}
}

var unannotatedEchoServiceEchoStreamDesc = &grpc.StreamDesc{
	StreamName: "Echo",
}

func (c *unannotatedEchoServiceClient) Echo(ctx context.Context, in *UnannotatedSimpleMessage, opts ...grpc.CallOption) (*UnannotatedSimpleMessage, error) {
	out := new(UnannotatedSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var unannotatedEchoServiceEchoBodyStreamDesc = &grpc.StreamDesc{
	StreamName: "EchoBody",
}

func (c *unannotatedEchoServiceClient) EchoBody(ctx context.Context, in *UnannotatedSimpleMessage, opts ...grpc.CallOption) (*UnannotatedSimpleMessage, error) {
	out := new(UnannotatedSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService/EchoBody", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var unannotatedEchoServiceEchoDeleteStreamDesc = &grpc.StreamDesc{
	StreamName: "EchoDelete",
}

func (c *unannotatedEchoServiceClient) EchoDelete(ctx context.Context, in *UnannotatedSimpleMessage, opts ...grpc.CallOption) (*UnannotatedSimpleMessage, error) {
	out := new(UnannotatedSimpleMessage)
	err := c.cc.Invoke(ctx, "/grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService/EchoDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UnannotatedEchoServiceService is the service API for UnannotatedEchoService service.
// Fields should be assigned to their respective handler implementations only before
// RegisterUnannotatedEchoServiceService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type UnannotatedEchoServiceService struct {
	// Echo method receives a simple message and returns it.
	//
	// The message posted as the id parameter will also be
	// returned.
	Echo func(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	// EchoBody method receives a simple message and returns it.
	EchoBody func(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	// EchoDelete method receives a simple message and returns it.
	EchoDelete func(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
}

func (s *UnannotatedEchoServiceService) echo(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnannotatedSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Echo(ctx, req.(*UnannotatedSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *UnannotatedEchoServiceService) echoBody(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnannotatedSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.EchoBody(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService/EchoBody",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.EchoBody(ctx, req.(*UnannotatedSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *UnannotatedEchoServiceService) echoDelete(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnannotatedSimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.EchoDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService/EchoDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.EchoDelete(ctx, req.(*UnannotatedSimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterUnannotatedEchoServiceService registers a service implementation with a gRPC server.
func RegisterUnannotatedEchoServiceService(s grpc.ServiceRegistrar, srv *UnannotatedEchoServiceService) {
	srvCopy := *srv
	if srvCopy.Echo == nil {
		srvCopy.Echo = func(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
		}
	}
	if srvCopy.EchoBody == nil {
		srvCopy.EchoBody = func(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error) {
			return nil, status.Errorf(codes.Unimplemented, "method EchoBody not implemented")
		}
	}
	if srvCopy.EchoDelete == nil {
		srvCopy.EchoDelete = func(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error) {
			return nil, status.Errorf(codes.Unimplemented, "method EchoDelete not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "grpc.gateway.examples.internal.proto.examplepb.UnannotatedEchoService",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Echo",
				Handler:    srvCopy.echo,
			},
			{
				MethodName: "EchoBody",
				Handler:    srvCopy.echoBody,
			},
			{
				MethodName: "EchoDelete",
				Handler:    srvCopy.echoDelete,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "examples/internal/proto/examplepb/unannotated_echo_service.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewUnannotatedEchoServiceService creates a new UnannotatedEchoServiceService containing the
// implemented methods of the UnannotatedEchoService service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewUnannotatedEchoServiceService(s interface{}) *UnannotatedEchoServiceService {
	ns := &UnannotatedEchoServiceService{}
	if h, ok := s.(interface {
		Echo(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	}); ok {
		ns.Echo = h.Echo
	}
	if h, ok := s.(interface {
		EchoBody(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	}); ok {
		ns.EchoBody = h.EchoBody
	}
	if h, ok := s.(interface {
		EchoDelete(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	}); ok {
		ns.EchoDelete = h.EchoDelete
	}
	return ns
}

// UnstableUnannotatedEchoServiceService is the service API for UnannotatedEchoService service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableUnannotatedEchoServiceService interface {
	// Echo method receives a simple message and returns it.
	//
	// The message posted as the id parameter will also be
	// returned.
	Echo(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	// EchoBody method receives a simple message and returns it.
	EchoBody(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
	// EchoDelete method receives a simple message and returns it.
	EchoDelete(context.Context, *UnannotatedSimpleMessage) (*UnannotatedSimpleMessage, error)
}
