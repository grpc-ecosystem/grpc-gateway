package server

import (
	"context"

	examples "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
	"google.golang.org/grpc/grpclog"
)

// Implements NonStandardServiceServer

type nonStandardServer struct{}

func newNonStandardServer() examples.NonStandardServiceServer {
	return new(nonStandardServer)
}

func (s *nonStandardServer) Update(ctx context.Context, msg *examples.NonStandardUpdateRequest) (*examples.NonStandardMessage, error) {
	grpclog.Info(msg)

	newMsg := &examples.NonStandardMessage{
		Thing: &examples.NonStandardMessage_Thing{SubThing: &examples.NonStandardMessage_Thing_SubThing{}}, // The fieldmask_helper doesn't generate nested structs if they are nil
	}
	applyFieldMask(newMsg, msg.Body, msg.UpdateMask)

	grpclog.Info(newMsg)
	return newMsg, nil
}

func (s *nonStandardServer) UpdateWithJSONNames(ctx context.Context, msg *examples.NonStandardWithJSONNamesUpdateRequest) (*examples.NonStandardMessageWithJSONNames, error) {
	grpclog.Info(msg)

	newMsg := &examples.NonStandardMessageWithJSONNames{
		Thing: &examples.NonStandardMessageWithJSONNames_Thing{SubThing: &examples.NonStandardMessageWithJSONNames_Thing_SubThing{}}, // The fieldmask_helper doesn't generate nested structs if they are nil
	}
	applyFieldMask(newMsg, msg.Body, msg.UpdateMask)

	grpclog.Info(newMsg)
	return newMsg, nil
}
