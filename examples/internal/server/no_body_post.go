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
	noBodyPost_contextChRPC    = make(chan context.Context)
	noBodyPost_contextChStream = make(chan context.Context)
)

func NoBodyPostServer_RetrieveContextRPC() context.Context {
	return <-noBodyPost_contextChRPC
}

func NoBodyPostServer_RetrieveContextStream() context.Context {
	return <-noBodyPost_contextChStream
}

type noBodyPostServer struct{}

func newNoBodyPostServer() examples.NoBodyPostServiceServer {
	return &noBodyPostServer{}
}

func (s noBodyPostServer) RpcEmptyRpc(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	noBodyPost_contextChRPC <- ctx
	<-ctx.Done()
	return nil, status.Error(codes.Canceled, "context canceled")
}

func (s noBodyPostServer) RpcEmptyStream(req *emptypb.Empty, stream grpc.ServerStreamingServer[emptypb.Empty]) error {
	noBodyPost_contextChStream <- stream.Context()
	<-stream.Context().Done()
	return status.Error(codes.Canceled, "context canceled")
}
