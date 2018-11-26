package server

import (
	"context"

	examples "github.com/grpc-ecosystem/grpc-gateway/examples/proto/examplepb"
)

// Implements of ResponseBodyServiceServer

type responseBodyServer struct{}

func newResponseBodyServer() examples.ResponseBodyServiceServer {
	return new(responseBodyServer)
}

func (s *responseBodyServer) GetResponseBody(ctx context.Context, req *examples.ResponseBodyIn) (*examples.ResponseBodyOut, error) {
	return &examples.ResponseBodyOut{
		Response: &examples.ResponseBodyOut_Response{
			Data: req.Data,
		},
	}, nil
}
