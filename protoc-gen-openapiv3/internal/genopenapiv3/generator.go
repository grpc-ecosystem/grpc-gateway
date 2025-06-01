package genopenapiv3

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	_ "google.golang.org/grpc/grpclog"
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
		doc, ok := g.extractFileOptions(t)
		if !ok {
			doc = &openapi3.T{}
		}

		for _, svc := range t.Services {
			g.generateServiceDoc(svc)
		}

		components := openapi3.NewComponents()
		components.Schemas = make(openapi3.Schemas)
		doc.Components = &components

		for _, msg := range t.Messages {
			msgName := g.getMessageName(msg.FQMN())
			schemaRef := g.generateMessageSchema(msg, components.Schemas)
			components.Schemas[msgName] = schemaRef
		}

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

func (g *generator) generateMessageSchema(msg *descriptor.Message, schemas openapi3.Schemas) *openapi3.SchemaRef {
	msgName := g.getMessageName(msg.FQMN())
	if scheme, ok := wktSchemas[msgName]; ok {
		return &openapi3.SchemaRef{
			Value: scheme,
		}
	}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeObject},
	}

	properties := make(openapi3.Schemas)

	for _, field := range msg.Fields {
		properties[field.GetName()] = g.generateFieldDoc(field, schemas)
	}

	schema.Properties = properties

	return &openapi3.SchemaRef{
		Value: schema,
	}
}

func (g *generator) generateFieldDoc(field *descriptor.Field, schemas openapi3.Schemas) *openapi3.SchemaRef {
	fd := field.FieldDescriptorProto
	location := ""
	if ix := strings.LastIndex(field.Message.FQMN(), "."); ix > 0 {
		location = field.Message.FQMN()[0:ix]
	}

	if m, err := g.reg.LookupMsg(location, field.GetTypeName()); err == nil {
		if opt := m.GetOptions(); opt != nil && opt.MapEntry != nil && *opt.MapEntry {
			// Generate Map<k, v> schema
			return &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					AdditionalProperties: openapi3.AdditionalProperties{
						Schema: g.generateFieldTypeSchema(m.GetField()[1], schemas),
					},
				},
			}
		}
	}

	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:  &openapi3.Types{openapi3.TypeArray},
				Items: g.generateFieldTypeSchema(fd, schemas),
			},
		}
	}

	return g.generateFieldTypeSchema(fd, schemas)
}

func (g *generator) generateFieldTypeSchema(fd *descriptorpb.FieldDescriptorProto, schemas openapi3.Schemas) *openapi3.SchemaRef {
	if schema, ok := primitiveTypeSchemas[fd.GetType()]; ok {
		return &openapi3.SchemaRef{
			Value: schema,
		}
	}

	switch ft := fd.GetType(); ft {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		openAPIRef, ok := g.fullyQualifiedNameToOpenAPIName(fd.GetTypeName())
		if !ok {
			panic(fmt.Sprintf("can't resolve OpenAPI ref from typename %q", fd.GetTypeName()))
		}

		return &openapi3.SchemaRef{
			Ref: "#/definitions/" + openAPIRef,
		}
	default:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{ft.String()}, Format: "UNKNOWN"},
		}
	}

}

func (g *generator) getMessageName(fqmn string) string {
	// TODO: have different naming stratgies
	return fqmn[1:]
}

func (g *generator) generateServiceDoc(svc *descriptor.Service) {
	for _, meth := range svc.Methods {
		g.generateMethodDoc(meth, svc.GetName())
	}
}

func (g *generator) generateMethodDoc(meth *descriptor.Method, serviceName string) {
	var pathItems []*openapi3.PathItem
	for bindingIdx, binding := range meth.Bindings {
		// Extract OpenAPI operation options from method descriptor
		opOpts, err := extractOperationOptionFromMethodDescriptor(meth.MethodDescriptorProto)
		if err != nil {
			// Log error but continue with default values
			// TODO: Add proper logging
			continue
		}

		operation := &openapi3.Operation{
			Tags:        []string{serviceName},
			Summary:     opOpts.GetSummary(),
			Description: opOpts.GetDescription(),
			OperationID: fmt.Sprintf("%s_%s_%d", serviceName, meth.GetName(), bindingIdx),
			Parameters:  openapi3.Parameters{},
			Responses:   &openapi3.Responses{},
		}

		// Create PathItem and assign the operation to the appropriate HTTP method
		pathItem := &openapi3.PathItem{}

		switch binding.HTTPMethod {
		case "GET":
			pathItem.Get = operation
		case "POST":
			operation.RequestBody = g.extractBodyRequest(binding)
			pathItem.Post = operation
		case "PUT":
			pathItem.Put = operation
		case "PATCH":
			pathItem.Patch = operation
		case "DELETE":
			pathItem.Delete = operation
		case "HEAD":
			pathItem.Head = operation
		case "OPTIONS":
			pathItem.Options = operation
		}
		// TODO: Extract and add path parameters to operation.Parameters
		// TODO: Extract and add query parameters to operation.Parameters
		// TODO: Add request body schema if applicable (for POST, PUT, PATCH)
		// TODO: Add response schemas
		// TODO: Add the PathItem to the OpenAPI document paths
		pathItems = append(pathItems, pathItem)
	}
}

// convertPathTemplate converts gRPC gateway path template to OpenAPI path format
// Example: "/v1/users/{user_id}" -> "/v1/users/{user_id}"
func (g *generator) convertPathTemplate(template string) string {
	// For now, return the template as-is since gRPC gateway templates
	// are already compatible with OpenAPI path format
	return template
}

func (g *generator) extractBodyRequest(binding *descriptor.Binding) *openapi3.RequestBodyRef {
	if binding.Body != nil && binding.Body.FieldPath == nil {
		// TODO: Create request body schema for the entire message
		return &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{},
						},
					},
				},
			},
		}
	}
	return nil
}

// extractOperationOptionFromMethodDescriptor extracts the message of type
// options.Operation from a given proto method's descriptor.
func extractOperationOptionFromMethodDescriptor(meth *descriptorpb.MethodDescriptorProto) (*options.Operation, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, options.E_Openapiv3Operation) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, options.E_Openapiv3Operation)
	opts, ok := ext.(*options.Operation)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an Operation", ext)
	}
	return opts, nil
}

// Take in a FQMN or FQEN and return a OpenAPI safe version of the FQMN and
// a boolean indicating if FQMN was properly resolved.
func (g *generator) fullyQualifiedNameToOpenAPIName(fqn string) (string, bool) {
	registriesSeenMutex.Lock()
	defer registriesSeenMutex.Unlock()
	if mapping, present := registriesSeen[g.reg]; present {
		ret, ok := mapping[fqn]
		return ret, ok
	}

	mapping := g.resolveFullyQualifiedNameToOpenAPINames(append(g.reg.GetAllFQMNs(), append(g.reg.GetAllFQENs(), g.reg.GetAllFQMethNs()...)...), g.reg.GetOpenAPINamingStrategy())
	registriesSeen[g.reg] = mapping
	ret, ok := mapping[fqn]

	return ret, ok
}

// Lookup message type by location.name and return an openapiv2-safe version
// of its FQMN.
func (g *generator) lookupMsgAndOpenAPIName(location, name string) (*descriptor.Message, string, error) {
	msg, err := g.reg.LookupMsg(location, name)
	if err != nil {
		return nil, "", err
	}
	swgName, ok := g.fullyQualifiedNameToOpenAPIName(msg.FQMN())
	if !ok {
		return nil, "", fmt.Errorf("can't map OpenAPI name from FQMN %q", msg.FQMN())
	}
	return msg, swgName, nil
}

// registriesSeen is used to memoise calls to resolveFullyQualifiedNameToOpenAPINames so
// we don't repeat it unnecessarily, since it can take some time.
var (
	registriesSeen      = map[*descriptor.Registry]map[string]string{}
	registriesSeenMutex sync.Mutex
)

// Take the names of every proto message and generate a unique reference for each, according to the given strategy.
func (g *generator) resolveFullyQualifiedNameToOpenAPINames(messages []string, _ string) map[string]string {
	strategyFn := func(messages []string) map[string]string {
		res := make(map[string]string)
		for _, msg := range messages {
			res[msg] = g.getMessageName(msg)
		}

		return res
	}

	return strategyFn(messages)
}
