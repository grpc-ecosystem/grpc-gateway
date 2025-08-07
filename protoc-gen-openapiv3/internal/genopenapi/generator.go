package genopenapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct{}

// New returns a new generator which generates OpenAPI v3 files.
func New() *generator {
	return &generator{}
}

// Generate generates OpenAPI v3 files from the given CodeGeneratorRequest.
func (g *generator) Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	reg := descriptor.NewRegistry()
	if err := reg.Load(req); err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	var files []*pluginpb.CodeGeneratorResponse_File
	for _, file := range req.FileToGenerate {
		f, err := reg.LookupFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup file %q: %w", file, err)
		}

		out, err := g.generate(f)
		if err != nil {
			return nil, fmt.Errorf("failed to generate openapi v3 for file %q: %w", file, err)
		}
		files = append(files, out)
	}

	return &pluginpb.CodeGeneratorResponse{
		File: files,
	}, nil
}

func (g *generator) generate(file *descriptor.File) (*pluginpb.CodeGeneratorResponse_File, error) {
	openapi := &OpenAPI{
		OpenAPI: "3.0.0",
		Info: &Info{
			Title:   file.GetName(),
			Version: "0.0.1",
		},
		Paths: &Paths{
			PathItems: make(map[string]*PathItem),
		},
		Components: &Components{
			Schemas: make(map[string]*Schema),
		},
	}

	g.renderMessagesAsDefinition(file, openapi.Components.Schemas)

	for _, svc := range file.Services {
		for _, meth := range svc.Methods {
			for _, b := range meth.Bindings {
				pathItem := &PathItem{}
				op := &Operation{
					Summary:   meth.GetName(),
					Responses: &Responses{},
				}

				var params []*Parameter
				for _, pathParam := range b.PathParams {
					params = append(params, &Parameter{
						Name:     pathParam.FieldPath.String(),
						In:       "path",
						Required: true,
						Schema:   schemaOfField(pathParam.Target),
					})
				}

				if b.Body != nil {
					var schema *Schema
					if len(b.Body.FieldPath) > 0 {
						schema = schemaOfField(b.Body.FieldPath[0].Target)
					} else {
						// TODO(ivucica): This should be a reference to a schema in components/schemas
						schema = &Schema{Type: "object"}
					}
					op.RequestBody = &RequestBody{
						Content: map[string]*MediaType{
							"application/json": {
								Schema: schema,
							},
						},
					}
				} else {
					for _, field := range meth.RequestType.Fields {
						if isPathParameter(field, b.PathParams) {
							continue
						}
						params = append(params, &Parameter{
							Name:   field.GetName(),
							In:     "query",
							Schema: schemaOfField(field),
						})
					}
				}
				op.Parameters = params
				op.Responses = &Responses{
					Responses: map[string]*Response{
						"200": {
							Description: "A successful response.",
							Content: map[string]*MediaType{
								"application/json": {
									Schema: &Schema{
										Ref: "#/components/schemas/" + meth.ResponseType.GetName(),
									},
								},
							},
						},
					},
				}

				switch b.HTTPMethod {
				case "GET":
					pathItem.Get = op
				case "POST":
					pathItem.Post = op
				case "PUT":
					pathItem.Put = op
				case "DELETE":
					pathItem.Delete = op
				case "PATCH":
					pathItem.Patch = op
				}

				openapi.Paths.PathItems[b.PathTmpl.Template] = pathItem
			}
		}
	}

	b, err := json.MarshalIndent(openapi, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openapi v3: %w", err)
	}

	name := file.GetName()
	ext := ".swagger.json"
	if strings.HasSuffix(name, ".proto") {
		name = name[:len(name)-len(".proto")]
	}
	name += ext

	return &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(name),
		Content: proto.String(string(b)),
	}, nil
}

func isPathParameter(field *descriptor.Field, pathParams []descriptor.Parameter) bool {
	for _, p := range pathParams {
		if p.FieldPath.String() == field.GetName() {
			return true
		}
	}
	return false
}

func (g *generator) renderMessagesAsDefinition(file *descriptor.File, schemas map[string]*Schema) {
	for _, msg := range file.Messages {
		schemas[msg.GetName()] = g.schemaOfMessage(msg)
	}
}

func (g *generator) schemaOfMessage(msg *descriptor.Message) *Schema {
	s := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
	}
	for _, field := range msg.Fields {
		s.Properties[field.GetName()] = schemaOfField(field)
	}
	return s
}

func schemaOfField(field *descriptor.Field) *Schema {
	s := &Schema{}
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		s.Type = "number"
		s.Format = "double"
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		s.Type = "number"
		s.Format = "float"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		s.Type = "string"
		s.Format = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		s.Type = "string"
		s.Format = "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		s.Type = "integer"
		s.Format = "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		s.Type = "string"
		s.Format = "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		s.Type = "integer"
		s.Format = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		s.Type = "boolean"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		s.Type = "string"
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		s.Type = "string"
		s.Format = "byte"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		s.Type = "integer"
		s.Format = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		s.Type = "integer"
		s.Format = "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		s.Type = "string"
		s.Format = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		s.Type = "integer"
		s.Format = "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		s.Type = "string"
		s.Format = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		// TODO(ivucica): handle message types
		s.Type = "object"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		// TODO(ivucica): handle enum types
		s.Type = "string"
	}
	return s
}
