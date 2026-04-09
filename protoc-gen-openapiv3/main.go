// Command protoc-gen-openapiv3 implements an OpenAPI 3.1.0 generator for
// proto files annotated with google.api.http rules.
//
// Status: alpha. The emitted JSON shape is not yet stable — encodings for
// oneofs, wrapper types, enums, and path-template expansion may change
// between minor releases while the mapping rules settle in response to
// real-world feedback. For a production-stable OpenAPI pipeline today,
// use protoc-gen-openapiv2.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/codegenerator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapi"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Write log output to standard error to prevent clobbering standard out.
	// Drop the date/time prefix the standard logger adds by default — buf
	// surfaces these messages interactively, so a timestamp adds noise
	// without context.
	log.SetOutput(os.Stderr)
	log.SetPrefix("protoc-gen-openapiv3: ")
	log.SetFlags(0)

	if err := run(); err != nil {
		emitError(err)
		os.Exit(1)
	}
}

func run() error {
	req, err := codegenerator.ParseRequest(os.Stdin)
	if err != nil {
		return err
	}

	reg := descriptor.NewRegistry()
	if err := reg.Load(req); err != nil {
		return err
	}

	var targets []*descriptor.File
	for _, name := range req.FileToGenerate {
		f, err := reg.LookupFile(name)
		if err != nil {
			return err
		}
		targets = append(targets, f)
	}

	out, err := genopenapi.Generate(reg, targets)
	if err != nil {
		return err
	}

	emitFiles(out)
	return nil
}

func emitFiles(files []*pluginpb.CodeGeneratorResponse_File) {
	resp := &pluginpb.CodeGeneratorResponse{File: files}
	codegenerator.SetSupportedFeaturesOnCodeGeneratorResponse(resp)
	emitResp(resp)
}

func emitError(err error) {
	// Echo to stderr in addition to the proto response
	grpclog.Infoln(err)
	emitResp(&pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitResp(resp *pluginpb.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
