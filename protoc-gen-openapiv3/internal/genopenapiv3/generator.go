package genopenapiv3

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/genproto/googleapis/api/visibility"
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
		fileGenerator := &fileGenerator{generator: g, spec: &openapi3.T{}}
		doc := fileGenerator.generateFileSpec(t)
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

func (g *generator) IsMessageVisible(msg *descriptor.Message) bool {
	if !proto.HasExtension(msg, visibility.E_MessageVisibility) {
		return true
	}

	ext := proto.GetExtension(msg.DescriptorProto, visibility.E_MessageVisibility)
	visibilityOpt, ok := ext.(*visibility.VisibilityRule)
	if ok {
		return g.isVisible(visibilityOpt)
	}

	return true
}

func (g *generator) IsFieldVisible(field *descriptor.Field) bool {
	if !proto.HasExtension(field, visibility.E_MessageVisibility) {
		return true
	}

	ext := proto.GetExtension(field.FieldDescriptorProto, visibility.E_MessageVisibility)
	visibilityOpt, ok := ext.(*visibility.VisibilityRule)
	if ok {
		return g.isVisible(visibilityOpt)
	}

	return true
}
func (g *generator) IsEnumVisible(enum *descriptor.Enum) bool {

	if !proto.HasExtension(enum, visibility.E_MessageVisibility) {
		return true
	}

	ext := proto.GetExtension(enum.EnumDescriptorProto, visibility.E_EnumVisibility)
	visibilityOpt, ok := ext.(*visibility.VisibilityRule)
	if ok {
		return g.isVisible(visibilityOpt)
	}

	return true
}

func (g *generator) IsMethodVisible(meth *descriptor.Method) bool {
	if !proto.HasExtension(meth, visibility.E_MessageVisibility) {
		return true
	}

	ext := proto.GetExtension(meth.MethodDescriptorProto, visibility.E_MethodVisibility)
	visibilityOpt, ok := ext.(*visibility.VisibilityRule)
	if ok {
		return g.isVisible(visibilityOpt)
	}

	return true
}

func (g *generator) IsServiceVisible(svc *descriptor.Service) bool {
	if !proto.HasExtension(svc, visibility.E_ApiVisibility) {
		return true
	}

	ext := proto.GetExtension(svc.ServiceDescriptorProto, visibility.E_ApiVisibility)
	visibilityOpt, ok := ext.(*visibility.VisibilityRule)
	if ok {
		return g.isVisible(visibilityOpt)
	}

	return true
}

func (g *generator) isVisible(rule *visibility.VisibilityRule) bool {
	if rule == nil {
		return true
	}

	return g.reg.GetVisibilityRestrictionSelectors()[rule.GetRestriction()]
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
