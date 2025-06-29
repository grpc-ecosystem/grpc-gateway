package genopenapiv3

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	reg    *descriptor.Registry
	format Format
}

func NewGenerator(reg *descriptor.Registry, format Format) gen.Generator {
	return &generator{
		reg:    reg,
		format: format,
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	err := g.loadPrequisiteProtos()
	if err != nil {
		return nil, fmt.Errorf("could not load prequisite proto files in registry: %w", err)
	}

	respFiles := make([]*descriptor.ResponseFile, 0, len(targets))
	docs := make([]*openapi3.T, 0, len(targets))
	for _, t := range targets {
		fileGenerator := &fileGenerator{generator: g, doc: &openapi3.T{}}
		doc := fileGenerator.generateFileDoc(t)
		docs = append(docs, doc)

		contentBytes, err := g.format.MarshalOpenAPIDoc(doc)
		if err != nil {
			return nil, err
		}

		base := filepath.Base(t.GetName())
		ext := filepath.Ext(base)
		fileName := fmt.Sprintf("%s.openapiv3.%s", base[:len(base)-len(ext)], g.format)

		respFiles = append(respFiles, &descriptor.ResponseFile{
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(fileName),
				Content: proto.String(string(contentBytes)),
			},
		})

	}

	mergedDocs, err := MergeOpenAPISpecs(docs...)
	if err != nil {
		return nil, fmt.Errorf("could not merge docs: %w", err)
	}

	contentBytes, err := g.format.MarshalOpenAPIDoc(mergedDocs)
	if err != nil {
		return nil, err
	}

	respFiles = append(respFiles, &descriptor.ResponseFile{
		CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("merged.openapiv3.%s", g.format)),
			Content: proto.String(string(contentBytes)),
		},
	})

	return respFiles, nil
}

func (g *generator) getOperationName(serviceName, methodName string, bindingIdx int) string {
	if bindingIdx == 0 {
		return fmt.Sprintf("%s_%s", serviceName, methodName)
	}
	return fmt.Sprintf("%s_%s_%d", serviceName, methodName, bindingIdx)
}

func (g *generator) fqmnToLocation(fqmn string) string {
	location := ""
	if ix := strings.LastIndex(fqmn, "."); ix > 0 {
		location = fqmn[0:ix]
	}
	return location
}

func (g *generator) resolveName(fqmn string) string {
	return fqmn
}

func (g *generator) resolveType(typeName string) string {
	return typeName
}

func (g *generator) loadPrequisiteProtos() error {
	any := protodesc.ToFileDescriptorProto((&anypb.Any{}).ProtoReflect().Descriptor().ParentFile())
	any.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	status := protodesc.ToFileDescriptorProto((&statuspb.Status{}).ProtoReflect().Descriptor().ParentFile())
	status.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	return g.reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			any,
			status,
		},
	})
}

func (g *generator) defaultResponse() (*descriptor.Message, error) {
	return g.reg.LookupMsg("", statusProtoFQMN)
}

func extractOperationOptionFromMethodDescriptor(meth *descriptorpb.MethodDescriptorProto) (*options.Operation, error) {
	if meth.Options == nil {
		return &options.Operation{}, nil
	}
	if !proto.HasExtension(meth.Options, options.E_Openapiv3Operation) {
		return &options.Operation{}, nil
	}
	ext := proto.GetExtension(meth.Options, options.E_Openapiv3Operation)
	opts, ok := ext.(*options.Operation)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an Operation", ext)
	}
	return opts, nil
}
