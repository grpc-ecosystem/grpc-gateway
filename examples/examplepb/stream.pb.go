// Code generated by protoc-gen-go.
// source: examples/examplepb/stream.proto
// DO NOT EDIT!

package examplepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"
import grpc_gateway_examples_sub "github.com/grpc-ecosystem/grpc-gateway/examples/sub"

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for StreamService service

type StreamServiceClient interface {
	BulkCreate(ctx context.Context, opts ...grpc.CallOption) (StreamService_BulkCreateClient, error)
	List(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (StreamService_ListClient, error)
	BulkEcho(ctx context.Context, opts ...grpc.CallOption) (StreamService_BulkEchoClient, error)
}

type streamServiceClient struct {
	cc *grpc.ClientConn
}

func NewStreamServiceClient(cc *grpc.ClientConn) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) BulkCreate(ctx context.Context, opts ...grpc.CallOption) (StreamService_BulkCreateClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_StreamService_serviceDesc.Streams[0], c.cc, "/grpc.gateway.examples.examplepb.StreamService/BulkCreate", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceBulkCreateClient{stream}
	return x, nil
}

type StreamService_BulkCreateClient interface {
	Send(*ABitOfEverything) error
	CloseAndRecv() (*google_protobuf1.Empty, error)
	grpc.ClientStream
}

type streamServiceBulkCreateClient struct {
	grpc.ClientStream
}

func (x *streamServiceBulkCreateClient) Send(m *ABitOfEverything) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamServiceBulkCreateClient) CloseAndRecv() (*google_protobuf1.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(google_protobuf1.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamServiceClient) List(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (StreamService_ListClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_StreamService_serviceDesc.Streams[1], c.cc, "/grpc.gateway.examples.examplepb.StreamService/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_ListClient interface {
	Recv() (*ABitOfEverything, error)
	grpc.ClientStream
}

type streamServiceListClient struct {
	grpc.ClientStream
}

func (x *streamServiceListClient) Recv() (*ABitOfEverything, error) {
	m := new(ABitOfEverything)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamServiceClient) BulkEcho(ctx context.Context, opts ...grpc.CallOption) (StreamService_BulkEchoClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_StreamService_serviceDesc.Streams[2], c.cc, "/grpc.gateway.examples.examplepb.StreamService/BulkEcho", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceBulkEchoClient{stream}
	return x, nil
}

type StreamService_BulkEchoClient interface {
	Send(*grpc_gateway_examples_sub.StringMessage) error
	Recv() (*grpc_gateway_examples_sub.StringMessage, error)
	grpc.ClientStream
}

type streamServiceBulkEchoClient struct {
	grpc.ClientStream
}

func (x *streamServiceBulkEchoClient) Send(m *grpc_gateway_examples_sub.StringMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamServiceBulkEchoClient) Recv() (*grpc_gateway_examples_sub.StringMessage, error) {
	m := new(grpc_gateway_examples_sub.StringMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for StreamService service

type StreamServiceServer interface {
	BulkCreate(StreamService_BulkCreateServer) error
	List(*google_protobuf1.Empty, StreamService_ListServer) error
	BulkEcho(StreamService_BulkEchoServer) error
}

func RegisterStreamServiceServer(s *grpc.Server, srv StreamServiceServer) {
	s.RegisterService(&_StreamService_serviceDesc, srv)
}

func _StreamService_BulkCreate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServiceServer).BulkCreate(&streamServiceBulkCreateServer{stream})
}

type StreamService_BulkCreateServer interface {
	SendAndClose(*google_protobuf1.Empty) error
	Recv() (*ABitOfEverything, error)
	grpc.ServerStream
}

type streamServiceBulkCreateServer struct {
	grpc.ServerStream
}

func (x *streamServiceBulkCreateServer) SendAndClose(m *google_protobuf1.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamServiceBulkCreateServer) Recv() (*ABitOfEverything, error) {
	m := new(ABitOfEverything)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StreamService_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(google_protobuf1.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).List(m, &streamServiceListServer{stream})
}

type StreamService_ListServer interface {
	Send(*ABitOfEverything) error
	grpc.ServerStream
}

type streamServiceListServer struct {
	grpc.ServerStream
}

func (x *streamServiceListServer) Send(m *ABitOfEverything) error {
	return x.ServerStream.SendMsg(m)
}

func _StreamService_BulkEcho_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServiceServer).BulkEcho(&streamServiceBulkEchoServer{stream})
}

type StreamService_BulkEchoServer interface {
	Send(*grpc_gateway_examples_sub.StringMessage) error
	Recv() (*grpc_gateway_examples_sub.StringMessage, error)
	grpc.ServerStream
}

type streamServiceBulkEchoServer struct {
	grpc.ServerStream
}

func (x *streamServiceBulkEchoServer) Send(m *grpc_gateway_examples_sub.StringMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamServiceBulkEchoServer) Recv() (*grpc_gateway_examples_sub.StringMessage, error) {
	m := new(grpc_gateway_examples_sub.StringMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _StreamService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gateway.examples.examplepb.StreamService",
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
	},
	Metadata: "examples/examplepb/stream.proto",
}

func init() { proto.RegisterFile("examples/examplepb/stream.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0xbf, 0x4a, 0x43, 0x31,
	0x14, 0xc6, 0xb9, 0x2a, 0xa2, 0x11, 0x97, 0x0c, 0x0e, 0x51, 0x28, 0x16, 0xc1, 0x2a, 0x92, 0xb4,
	0xba, 0xb9, 0x59, 0xe9, 0xa6, 0x38, 0x74, 0x73, 0x29, 0xc9, 0xe5, 0x34, 0x0d, 0xbd, 0xf7, 0x26,
	0x24, 0xe7, 0x56, 0x0b, 0x4e, 0x8e, 0xae, 0x7d, 0x11, 0xdf, 0xc5, 0x57, 0xf0, 0x41, 0xa4, 0xf7,
	0xdf, 0xd4, 0xd2, 0xba, 0x25, 0x9c, 0x2f, 0xf9, 0x7e, 0xe7, 0x47, 0x5a, 0xf0, 0x2e, 0x53, 0x97,
	0x40, 0x10, 0xd5, 0xc1, 0x29, 0x11, 0xd0, 0x83, 0x4c, 0xb9, 0xf3, 0x16, 0x2d, 0x6d, 0x69, 0xef,
	0x62, 0xae, 0x25, 0xc2, 0x9b, 0x9c, 0xf3, 0x3a, 0xcd, 0x9b, 0x34, 0x3b, 0xd3, 0xd6, 0xea, 0x04,
	0x84, 0x74, 0x46, 0xc8, 0x2c, 0xb3, 0x28, 0xd1, 0xd8, 0x2c, 0x94, 0xcf, 0xd9, 0x69, 0x35, 0x2d,
	0x6e, 0x2a, 0x1f, 0x0b, 0x48, 0x1d, 0xce, 0xab, 0xe1, 0xcd, 0x8a, 0x72, 0x39, 0x52, 0x06, 0x47,
	0x76, 0x3c, 0x82, 0x19, 0xf8, 0x39, 0x4e, 0x4c, 0xa6, 0xab, 0x34, 0x6b, 0xd2, 0x21, 0x57, 0x22,
	0x85, 0x10, 0xa4, 0x86, 0x72, 0x76, 0xfb, 0xbd, 0x4b, 0x8e, 0x87, 0x05, 0xf6, 0x10, 0xfc, 0xcc,
	0xc4, 0x40, 0xbf, 0x22, 0x42, 0xfa, 0x79, 0x32, 0x7d, 0xf4, 0x20, 0x11, 0x68, 0x8f, 0x6f, 0xd8,
	0x83, 0x3f, 0xf4, 0x0d, 0xbe, 0x8c, 0x07, 0x4d, 0x2b, 0x3b, 0xe1, 0x25, 0x3b, 0xaf, 0xd9, 0xf9,
	0x60, 0xc9, 0xde, 0x16, 0x9f, 0x3f, 0xbf, 0x8b, 0x9d, 0xab, 0xf6, 0x85, 0x98, 0xf5, 0x6a, 0xf0,
	0x55, 0xd8, 0x42, 0xe5, 0xc9, 0xf4, 0x3e, 0xba, 0xee, 0x44, 0xf4, 0x83, 0xec, 0x3d, 0x99, 0x80,
	0x74, 0xcd, 0x97, 0xec, 0xff, 0x74, 0xed, 0xcb, 0x82, 0xe2, 0x9c, 0xb6, 0x36, 0x50, 0x74, 0x23,
	0xba, 0x88, 0xc8, 0xc1, 0x52, 0xc5, 0x20, 0x9e, 0x58, 0xda, 0x59, 0x53, 0x15, 0x72, 0xc5, 0x87,
	0xe8, 0x4d, 0xa6, 0x9f, 0x4b, 0xb3, 0x6c, 0xeb, 0xe4, 0xf6, 0x46, 0x20, 0x9e, 0xd8, 0xc2, 0x48,
	0x37, 0xea, 0x1f, 0xbd, 0x1e, 0x36, 0xeb, 0xa9, 0xfd, 0x42, 0xc8, 0xdd, 0x5f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xbc, 0x52, 0x49, 0x85, 0x8f, 0x02, 0x00, 0x00,
}
