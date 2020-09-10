package server

import (
	"context"

	"github.com/golang/glog"
	examples "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
)

// Implements NonStandardServiceServer

type nonStandardServer struct {
}

func newNonStandardServer() *examples.NonStandardServiceService {
	var service nonStandardServer
	// to ensure everything is implemented, unstable is a unhappy name :-(
	var _ examples.UnstableNonStandardServiceService = service
	return examples.NewNonStandardServiceService(service)
}

func (s nonStandardServer) Update(ctx context.Context, msg *examples.NonStandardUpdateRequest) (*examples.NonStandardMessage, error) {
	glog.Info(msg)

	newMsg := &examples.NonStandardMessage{
		Thing: &examples.NonStandardMessage_Thing{SubThing: &examples.NonStandardMessage_Thing_SubThing{}}, // The fieldmask_helper doesn't generate nested structs if they are nil
	}
	applyFieldMask(newMsg, msg.Body, msg.UpdateMask)

	glog.Info(newMsg)
	return newMsg, nil
}

func (s nonStandardServer) UpdateWithJSONNames(ctx context.Context, msg *examples.NonStandardWithJSONNamesUpdateRequest) (*examples.NonStandardMessageWithJSONNames, error) {
	glog.Info(msg)

	newMsg := &examples.NonStandardMessageWithJSONNames{
		Thing: &examples.NonStandardMessageWithJSONNames_Thing{SubThing: &examples.NonStandardMessageWithJSONNames_Thing_SubThing{}}, // The fieldmask_helper doesn't generate nested structs if they are nil
	}
	applyFieldMask(newMsg, msg.Body, msg.UpdateMask)

	glog.Info(newMsg)
	return newMsg, nil
}
