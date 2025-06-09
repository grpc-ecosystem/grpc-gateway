package genopenapiv3

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
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
	respFiles := make([]*descriptor.ResponseFile, len(targets))
	for i, t := range targets {
		fileGenerator := &fileGenerator{generator: g, doc: &openapi3.T{}}
		doc := fileGenerator.generateFileDoc(t)

		contentBytes, err := g.format.MarshalOpenAPIDoc(doc)
		if err != nil {
			return nil, err
		}

		base := filepath.Base(t.GetName())
		ext := filepath.Ext(base)
		fileName := fmt.Sprintf("%s.openapiv3.%s", base[:len(base)-len(ext)], g.format)

		respFiles[i] = &descriptor.ResponseFile{
			GoPkg: t.GoPkg,
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(fileName),
				Content: proto.String(string(contentBytes)),
			},
		}
	}

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
	grpclog.Infof("resolveFQMN: %s", fqmn)
	return fqmn
}

func (g *generator) resolveType(typeName string) string {
	grpclog.Infof("resolveType: %s", typeName)
	return typeName
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
