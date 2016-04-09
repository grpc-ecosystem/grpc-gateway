// Package generator provides an abstract interface to code generators.
package generator

import (
	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Generator is an abstraction of code generators.
type Generator interface {
	// Generate generates output files from input .proto files.
	Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error)
}
