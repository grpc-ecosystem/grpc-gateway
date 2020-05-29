package server

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	examples "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
)

// Implements of ResponseBodyServiceServer

type responseBodyServer struct {
	examples.UnimplementedResponseBodyServiceServer
}

func newResponseBodyServer() examples.ResponseBodyServiceServer {
	return new(responseBodyServer)
}

func (s *responseBodyServer) GetResponseBody(ctx context.Context, req *examples.ResponseBodyIn) (*examples.ResponseBodyOut, error) {
	glog.Info(req)
	return &examples.ResponseBodyOut{
		Response: &examples.ResponseBodyOut_Response{
			Data: req.Data,
		},
	}, nil
}

func (s *responseBodyServer) ListResponseBodies(ctx context.Context, req *examples.ResponseBodyIn) (*examples.RepeatedResponseBodyOut, error) {
	glog.Info(req)
	return &examples.RepeatedResponseBodyOut{
		Response: []*examples.RepeatedResponseBodyOut_Response{
			{
				Data: req.Data,
			},
		},
	}, nil
}

func (s *responseBodyServer) ListResponseStrings(ctx context.Context, req *examples.ResponseBodyIn) (*examples.RepeatedResponseStrings, error) {
	glog.Info(req)
	if req.Data == "empty" {
		return &examples.RepeatedResponseStrings{
			Values: []string{},
		}, nil
	}
	return &examples.RepeatedResponseStrings{
		Values: []string{"hello", req.Data},
	}, nil
}

func (s *responseBodyServer) GetResponseBodyStream(req *examples.ResponseBodyIn, stream examples.ResponseBodyService_GetResponseBodyStreamServer) error {
	glog.Info(req)
	if err := stream.Send(&examples.ResponseBodyOut{
		Response: &examples.ResponseBodyOut_Response{
			Data: fmt.Sprintf("first %s", req.Data),
		},
	}); err != nil {
		return err
	}

	return stream.Send(&examples.ResponseBodyOut{
		Response: &examples.ResponseBodyOut_Response{
			Data: fmt.Sprintf("second %s", req.Data),
		},
	})
}
