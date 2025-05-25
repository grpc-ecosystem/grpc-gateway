package genopenapiv3

import (
	"fmt"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	_ "google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var wktSchemas = map[string]*openapi3.Schema{
	".google.protobuf.FieldMask": {
		Type: &openapi3.Types{"string"},
	},
	".google.protobuf.Timestamp": {
		Type:   &openapi3.Types{"string"},
		Format: "date-time",
	},
	".google.protobuf.Duration": {
		Type: &openapi3.Types{"string"},
	},
	".google.protobuf.StringValue": {
		Type: &openapi3.Types{"string"},
	},
	".google.protobuf.BytesValue": {
		Type:   &openapi3.Types{"string"},
		Format: "byte",
	},
	".google.protobuf.Int32Value": {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	".google.protobuf.UInt32Value": {
		Type:   &openapi3.Types{"integer"},
		Format: "int64",
	},
	".google.protobuf.Int64Value": {
		Type:   &openapi3.Types{"string"},
		Format: "int64",
	},
	".google.protobuf.UInt64Value": {
		Type:   &openapi3.Types{"string"},
		Format: "uint64",
	},
	".google.protobuf.FloatValue": {
		Type:   &openapi3.Types{"number"},
		Format: "float",
	},
	".google.protobuf.DoubleValue": {
		Type:   &openapi3.Types{"number"},
		Format: "double",
	},
	".google.protobuf.BoolValue": {
		Type: &openapi3.Types{"boolean"},
	},
	".google.protobuf.Empty": {
		Type: &openapi3.Types{"object"},
	},
	".google.protobuf.Struct": {
		Type: &openapi3.Types{"object"},
	},
	".google.protobuf.Value": {},
	".google.protobuf.ListValue": {
		Type: &openapi3.Types{"array"},
		Items: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"object"},
			},
		},
	},
	".google.protobuf.NullValue": {
		Type: &openapi3.Types{"string"},
	},
}

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
		for _, msg := range t.Messages {
			g.addMessageToSchemes(msg, components.Schemas)
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

func (g *generator) addMessageToSchemes(msg *descriptor.Message, schemas openapi3.Schemas) {
	if scheme, ok := wktSchemas[msg.FQMN()]; ok {
		schemas[msg.FQMN()] = &openapi3.SchemaRef{
			Value: scheme,
		}
		return
	}

	// TODO: Implement schema generation for custom messages
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

func primitiveTypeSchema(t descriptorpb.FieldDescriptorProto_Type) (*openapi3.Schema, bool) {
	switch t {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return &openapi3.Schema{
			Type:   &openapi3.Types{"number"},
			Format: "double",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return &openapi3.Schema{
			Type:   &openapi3.Types{"number"},
			Format: "float",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		// 64bit integer types are marshaled as string in the default JSONPb marshaler.
		// This maintains compatibility with JSON's limited number precision.
		return &openapi3.Schema{
			Type:   &openapi3.Types{"string"},
			Format: "int64",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		// 64bit integer types are marshaled as string in the default JSONPb marshaler.
		// TODO(yugui) Add an option to declare 64bit integers as int64.
		//
		// NOTE: uint64 is not a standard format in OpenAPI spec.
		// So we cannot expect that uint64 is commonly supported by OpenAPI processors.
		return &openapi3.Schema{
			Type:   &openapi3.Types{"string"},
			Format: "uint64",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		return &openapi3.Schema{
			Type:   &openapi3.Types{"integer"},
			Format: "int32",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		// 64bit types marshaled as string for JSON compatibility
		return &openapi3.Schema{
			Type:   &openapi3.Types{"string"},
			Format: "uint64",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		// Fixed 32-bit unsigned integer
		return &openapi3.Schema{
			Type:   &openapi3.Types{"integer"},
			Format: "int32",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		// NOTE: In OpenAPI v3 specification, format should be empty on boolean type
		return &openapi3.Schema{
			Type: &openapi3.Types{"boolean"},
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		// NOTE: In OpenAPI v3 specification, format can be empty on string type
		return &openapi3.Schema{
			Type: &openapi3.Types{"string"},
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		// Base64 encoded string representation
		return &openapi3.Schema{
			Type:   &openapi3.Types{"string"},
			Format: "byte",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		// 32-bit unsigned integer
		return &openapi3.Schema{
			Type:   &openapi3.Types{"integer"},
			Format: "int32",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return &openapi3.Schema{
			Type:   &openapi3.Types{"integer"},
			Format: "int32",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		// 64bit types marshaled as string for JSON compatibility
		return &openapi3.Schema{
			Type:   &openapi3.Types{"string"},
			Format: "int64",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return &openapi3.Schema{
			Type:   &openapi3.Types{"integer"},
			Format: "int32",
		}, true
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		// 64bit types marshaled as string for JSON compatibility
		return &openapi3.Schema{
			Type:   &openapi3.Types{"string"},
			Format: "int64",
		}, true
	default:
		return nil, false
	}
}
