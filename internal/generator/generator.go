// Package generator provides an abstract interface to code generators.
package generator

import (
	pluginpb "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptorpb "github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

// Generator is an abstraction of code generators.
type Generator interface {
	// Generate generates output files from input .proto files.
	Generate(targets []*descriptorpb.File) ([]*pluginpb.CodeGeneratorResponse_File, error)
}
