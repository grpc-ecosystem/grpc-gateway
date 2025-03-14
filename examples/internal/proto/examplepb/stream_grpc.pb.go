// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: examples/internal/proto/examplepb/stream.proto

package examplepb

import (
	context "context"
	sub "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/sub"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	StreamService_BulkCreate_FullMethodName       = "/grpc.gateway.examples.internal.proto.examplepb.StreamService/BulkCreate"
	StreamService_List_FullMethodName             = "/grpc.gateway.examples.internal.proto.examplepb.StreamService/List"
	StreamService_BulkEcho_FullMethodName         = "/grpc.gateway.examples.internal.proto.examplepb.StreamService/BulkEcho"
	StreamService_BulkEchoDuration_FullMethodName = "/grpc.gateway.examples.internal.proto.examplepb.StreamService/BulkEchoDuration"
	StreamService_Download_FullMethodName         = "/grpc.gateway.examples.internal.proto.examplepb.StreamService/Download"
)

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Defines some more operations to be added to ABitOfEverythingService
type StreamServiceClient interface {
	BulkCreate(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[ABitOfEverything, emptypb.Empty], error)
	List(ctx context.Context, in *Options, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ABitOfEverything], error)
	BulkEcho(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[sub.StringMessage, sub.StringMessage], error)
	BulkEchoDuration(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[durationpb.Duration, durationpb.Duration], error)
	Download(ctx context.Context, in *Options, opts ...grpc.CallOption) (grpc.ServerStreamingClient[httpbody.HttpBody], error)
}

type streamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamServiceClient(cc grpc.ClientConnInterface) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) BulkCreate(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[ABitOfEverything, emptypb.Empty], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[0], StreamService_BulkCreate_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ABitOfEverything, emptypb.Empty]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_BulkCreateClient = grpc.ClientStreamingClient[ABitOfEverything, emptypb.Empty]

func (c *streamServiceClient) List(ctx context.Context, in *Options, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ABitOfEverything], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[1], StreamService_List_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Options, ABitOfEverything]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_ListClient = grpc.ServerStreamingClient[ABitOfEverything]

func (c *streamServiceClient) BulkEcho(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[sub.StringMessage, sub.StringMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[2], StreamService_BulkEcho_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[sub.StringMessage, sub.StringMessage]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_BulkEchoClient = grpc.BidiStreamingClient[sub.StringMessage, sub.StringMessage]

func (c *streamServiceClient) BulkEchoDuration(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[durationpb.Duration, durationpb.Duration], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[3], StreamService_BulkEchoDuration_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[durationpb.Duration, durationpb.Duration]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_BulkEchoDurationClient = grpc.BidiStreamingClient[durationpb.Duration, durationpb.Duration]

func (c *streamServiceClient) Download(ctx context.Context, in *Options, opts ...grpc.CallOption) (grpc.ServerStreamingClient[httpbody.HttpBody], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[4], StreamService_Download_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Options, httpbody.HttpBody]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_DownloadClient = grpc.ServerStreamingClient[httpbody.HttpBody]

// StreamServiceServer is the server API for StreamService service.
// All implementations should embed UnimplementedStreamServiceServer
// for forward compatibility.
//
// Defines some more operations to be added to ABitOfEverythingService
type StreamServiceServer interface {
	BulkCreate(grpc.ClientStreamingServer[ABitOfEverything, emptypb.Empty]) error
	List(*Options, grpc.ServerStreamingServer[ABitOfEverything]) error
	BulkEcho(grpc.BidiStreamingServer[sub.StringMessage, sub.StringMessage]) error
	BulkEchoDuration(grpc.BidiStreamingServer[durationpb.Duration, durationpb.Duration]) error
	Download(*Options, grpc.ServerStreamingServer[httpbody.HttpBody]) error
}

// UnimplementedStreamServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStreamServiceServer struct{}

func (UnimplementedStreamServiceServer) BulkCreate(grpc.ClientStreamingServer[ABitOfEverything, emptypb.Empty]) error {
	return status.Errorf(codes.Unimplemented, "method BulkCreate not implemented")
}
func (UnimplementedStreamServiceServer) List(*Options, grpc.ServerStreamingServer[ABitOfEverything]) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedStreamServiceServer) BulkEcho(grpc.BidiStreamingServer[sub.StringMessage, sub.StringMessage]) error {
	return status.Errorf(codes.Unimplemented, "method BulkEcho not implemented")
}
func (UnimplementedStreamServiceServer) BulkEchoDuration(grpc.BidiStreamingServer[durationpb.Duration, durationpb.Duration]) error {
	return status.Errorf(codes.Unimplemented, "method BulkEchoDuration not implemented")
}
func (UnimplementedStreamServiceServer) Download(*Options, grpc.ServerStreamingServer[httpbody.HttpBody]) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedStreamServiceServer) testEmbeddedByValue() {}

// UnsafeStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServiceServer will
// result in compilation errors.
type UnsafeStreamServiceServer interface {
	mustEmbedUnimplementedStreamServiceServer()
}

func RegisterStreamServiceServer(s grpc.ServiceRegistrar, srv StreamServiceServer) {
	// If the following call pancis, it indicates UnimplementedStreamServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StreamService_ServiceDesc, srv)
}

func _StreamService_BulkCreate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServiceServer).BulkCreate(&grpc.GenericServerStream[ABitOfEverything, emptypb.Empty]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_BulkCreateServer = grpc.ClientStreamingServer[ABitOfEverything, emptypb.Empty]

func _StreamService_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Options)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).List(m, &grpc.GenericServerStream[Options, ABitOfEverything]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_ListServer = grpc.ServerStreamingServer[ABitOfEverything]

func _StreamService_BulkEcho_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServiceServer).BulkEcho(&grpc.GenericServerStream[sub.StringMessage, sub.StringMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_BulkEchoServer = grpc.BidiStreamingServer[sub.StringMessage, sub.StringMessage]

func _StreamService_BulkEchoDuration_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServiceServer).BulkEchoDuration(&grpc.GenericServerStream[durationpb.Duration, durationpb.Duration]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_BulkEchoDurationServer = grpc.BidiStreamingServer[durationpb.Duration, durationpb.Duration]

func _StreamService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Options)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).Download(m, &grpc.GenericServerStream[Options, httpbody.HttpBody]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamService_DownloadServer = grpc.ServerStreamingServer[httpbody.HttpBody]

// StreamService_ServiceDesc is the grpc.ServiceDesc for StreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gateway.examples.internal.proto.examplepb.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BulkCreate",
			Handler:       _StreamService_BulkCreate_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "List",
			Handler:       _StreamService_List_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BulkEcho",
			Handler:       _StreamService_BulkEcho_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "BulkEchoDuration",
			Handler:       _StreamService_BulkEchoDuration_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _StreamService_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "examples/internal/proto/examplepb/stream.proto",
}
