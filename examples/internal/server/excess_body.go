package server

import (
	"context"

	examples "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	excessBody_contextChRPC    = make(chan context.Context)
	excessBody_contextChStream = make(chan context.Context)
)

func ExcessBodyServer_RetrieveContextRPC() context.Context {
	return <-excessBody_contextChRPC
}

func ExcessBodyServer_RetrieveContextStream() context.Context {
	return <-excessBody_contextChStream
}

type excessBodyServer struct{}

func newExcessBodyServer() examples.ExcessBodyServiceServer {
	return &excessBodyServer{}
}

func (s excessBodyServer) NoBodyRpc(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	excessBody_contextChRPC <- ctx
	<-ctx.Done()
	return nil, status.Error(codes.Canceled, "context canceled")
}

func (s excessBodyServer) NoBodyServerStream(req *emptypb.Empty, stream grpc.ServerStreamingServer[emptypb.Empty]) error {
	excessBody_contextChStream <- stream.Context()
	<-stream.Context().Done()
	return status.Error(codes.Canceled, "context canceled")
}

func (s excessBodyServer) WithBodyRpc(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	excessBody_contextChRPC <- ctx
	<-ctx.Done()
	return nil, status.Error(codes.Canceled, "context canceled")
}

func (s excessBodyServer) WithBodyServerStream(req *emptypb.Empty, stream grpc.ServerStreamingServer[emptypb.Empty]) error {
	excessBody_contextChStream <- stream.Context()
	<-stream.Context().Done()
	return status.Error(codes.Canceled, "context canceled")
}
