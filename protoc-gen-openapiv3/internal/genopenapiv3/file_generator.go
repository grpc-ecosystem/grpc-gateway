package genopenapiv3

import (
	"maps"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/descriptorpb"
)

type fileGenerator struct {
	*generator

	doc *openapi3.T
}

func (fg *fileGenerator) generateFileDoc(file *descriptor.File) *openapi3.T {
	var ok bool
	fg.doc, ok = fg.extractFileOptions(file)
	if !ok {
		fg.doc = &openapi3.T{
			OpenAPI: OpenAPIVersion,
		}
	}

	fg.doc.Components = new(openapi3.Components)
	fg.doc.Components.Schemas = make(openapi3.Schemas)

	if fg.doc.Paths == nil {
		fg.doc.Paths = &openapi3.Paths{}
	}

	for _, svc := range file.Services {
		fg.generateServiceDoc(svc)
	}

	for _, msg := range file.Messages {
		fg.getMessageSchema(msg)
	}

	for _, enum := range file.Enums {
		fg.getEnumSchema(enum)
	}

	return fg.doc
}

func (fg *fileGenerator) getMessageSchema(msg *descriptor.Message) *openapi3.SchemaRef {
	name := fg.resolveName(msg.FQMN())
	schemaRef, ok := fg.doc.Components.Schemas[name]
	if ok {
		return schemaRef.Value.NewRef()
	}

	schemaRef = fg.generateMessageSchema(msg).NewRef()
	fg.doc.Components.Schemas[name] = schemaRef

	return schemaRef
}

func (fg *fileGenerator) generateMessageSchema(msg *descriptor.Message) *openapi3.Schema {
	msgName := fg.resolveName(msg.FQMN())
	if scheme, ok := wktSchemas[msgName]; ok {
		return scheme
	}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeObject},
	}

	properties := make(openapi3.Schemas)
	tempOneOfsProperties := make(map[int32]openapi3.Schemas)
	for _, field := range msg.Fields {
		fieldDoc := fg.generateFieldDoc(field)
		if field.OneofIndex != nil {
			if tempOneOfsProperties[*field.OneofIndex] == nil {
				tempOneOfsProperties[*field.OneofIndex] = make(openapi3.Schemas)
			}

			tempOneOfsProperties[*field.OneofIndex][field.GetName()] = fieldDoc
		} else {
			properties[field.GetName()] = fieldDoc
		}
	}

	allOneOfsProperties := make(map[int32]openapi3.Schemas)

	for oneOfKey, oneOfProperties := range tempOneOfsProperties {
		if len(oneOfProperties) == 1 {
			keys := slices.Collect(maps.Keys(oneOfProperties))
			key := keys[0]
			value := oneOfProperties[key]
			properties[key] = value
		} else {
			allOneOfsProperties[oneOfKey] = oneOfProperties
		}
	}

	if len(allOneOfsProperties) == 0 {
		schema.Properties = properties
	} else {
		switch fg.reg.GetOneOfStrategy() {
		case "oneOf":
			return &openapi3.Schema{
				OneOf: fg.generateMessageWithOneOfsSchemas(allOneOfsProperties, properties, msg.GetOneofDecl()),
			}
		default:
			grpclog.Fatal("unknown oneof strategy")
		}
	}

	return schema
}

/*
this type of oneof strategy, creates a oneof object for every possible combination of oneof fields
e.g.: if you have a proto like this:

	message sample{
		oneof one {
			string field_one = 1;
			string field_two = 2;
		}

		oneof two {
			string field_three = 3;
			string field_four = 4;
		}
	}

2 * 2 = 4, object schemas will be generate for each combination of set {field_one, field_two} and {field_three, field_four}
*/
func (fg *fileGenerator) generateMessageWithOneOfsSchemas(allOneOfsProperties map[int32]openapi3.Schemas, properties openapi3.Schemas,
	oneOfs []*descriptorpb.OneofDescriptorProto) openapi3.SchemaRefs {
	if len(oneOfs) == 0 {
		return openapi3.SchemaRefs{&openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:       &openapi3.Types{openapi3.TypeObject},
				Properties: properties,
			},
		}}
	}

	oneOfId := len(oneOfs) - 1
	newOneOfs := oneOfs[:oneOfId]

	oneOfProperties := allOneOfsProperties[int32(oneOfId)]

	newAllOneOfsProperties := maps.Clone(allOneOfsProperties)
	delete(newAllOneOfsProperties, int32(oneOfId))

	var res openapi3.SchemaRefs

	for fieldName, fieldSchema := range oneOfProperties {
		newProperties := maps.Clone(properties)
		newProperties[fieldName] = fieldSchema
		res = append(res, fg.generateMessageWithOneOfsSchemas(newAllOneOfsProperties, newProperties, newOneOfs)...)
	}

	return res
}

func (fg *fileGenerator) generateFieldDoc(field *descriptor.Field) *openapi3.SchemaRef {
	location := fg.fqmnToLocation(field.Message.FQMN())
	if m, err := fg.reg.LookupMsg(location, field.GetTypeName()); err == nil {
		if opt := m.GetOptions(); opt != nil && opt.MapEntry != nil && *opt.MapEntry {
			FieldDesc := m.GetField()[1]

			return &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					AdditionalProperties: openapi3.AdditionalProperties{
						Schema: fg.generateFieldTypeSchema(FieldDesc, location),
					},
				},
			}
		}
	}

	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:  &openapi3.Types{openapi3.TypeArray},
				Items: fg.generateFieldTypeSchema(field.FieldDescriptorProto, location),
			},
		}
	}

	return fg.generateFieldTypeSchema(field.FieldDescriptorProto, location)
}

func (fg *fileGenerator) generateFieldTypeSchema(fd *descriptorpb.FieldDescriptorProto, location string) *openapi3.SchemaRef {
	if schema, ok := primitiveTypeSchemas[fd.GetType()]; ok {
		return &openapi3.SchemaRef{
			Value: schema,
		}
	}

	switch ft := fd.GetType(); ft {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		openAPIRef := fg.resolveType(fd.GetTypeName())
		if schema, ok := fg.doc.Components.Schemas[openAPIRef]; ok {
			return schema
		} else {
			if fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
				fieldTypeEnum, err := fg.reg.LookupEnum(location, fd.GetTypeName())
				if err != nil {
					panic(err)
				}

				return fg.getEnumSchema(fieldTypeEnum)
			} else {
				fieldTypeMsg, err := fg.reg.LookupMsg(location, fd.GetTypeName())
				if err != nil {
					panic(err)
				}
				return fg.getMessageSchema(fieldTypeMsg)
			}
		}
	default:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{ft.String()}, Format: "UNKNOWN"},
		}
	}
}

func (fg *fileGenerator) getEnumSchema(enum *descriptor.Enum) *openapi3.SchemaRef {
	name := fg.resolveName(enum.FQEN())

	schemaRef, ok := fg.doc.Components.Schemas[name]
	if ok {
		return schemaRef
	}

	schemaRef = fg.generateEnumSchema(enum).NewRef()
	fg.doc.Components.Schemas[name] = schemaRef

	return schemaRef
}

func (fg *fileGenerator) generateEnumSchema(enum *descriptor.Enum) *openapi3.Schema {
	var enumValues []interface{}
	for _, value := range enum.GetValue() {
		enumValues = append(enumValues, value.GetName())
	}

	return &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeString},
		Enum: enumValues,
	}
}

func (fg *fileGenerator) generateServiceDoc(svc *descriptor.Service) {
	for _, meth := range svc.Methods {
		fg.generateMethodDoc(meth)
	}
}

func (fg *fileGenerator) generateMethodDoc(meth *descriptor.Method) error {
	return nil
	// for bindingIdx, binding := range meth.Bindings {
	// 	opOpts, err := extractOperationOptionFromMethodDescriptor(meth.MethodDescriptorProto)
	// 	if err != nil {
	// 		return fmt.Errorf("error extracting method %s operations: %v", meth.GetName(), err)
	// 	}
	//
	// 	pathParams, err := fg.generatePathParameters(binding.PathParams)
	// 	if err != nil {
	// 		return fmt.Errorf("error generating path parameters for method %s: %v", meth.GetName(), err)
	// 	}
	//
	// 	if meth.RequestType != nil {
	// 		switch binding.HTTPMethod {
	// 		case "GET", "DELETE":
	// 			queryParams, err := fg.messageToQueryParameters(meth.RequestType, binding.PathParams, binding.Body, binding.HTTPMethod)
	// 			if err != nil {
	// 				grpclog.Errorf("error generating query parameters for method %s: %v", meth.GetName(), err)
	// 			} else {
	// 				pathParams = append(pathParams, queryParams...)
	// 			}
	// 		case "POST", "PUT", "PATCH":
	// 			// For POST, PUT, PATCH, add request body
	// 			operation.RequestBody = fg.generateRequestBody(binding, meth.RequestType)
	//
	// 			queryParams, err := fg.messageToQueryParameters(meth.RequestType, binding.PathParams, binding.Body, binding.HTTPMethod)
	// 			if err != nil {
	// 				grpclog.Errorf("error generating query parameters for method %s: %v", meth.GetName(), err)
	// 			} else {
	// 				pathParams = append(pathParams, queryParams...)
	// 			}
	// 		}
	// 	}
	//
	// 	var responses *openapi3.Responses
	//
	// 	operation := &openapi3.Operation{
	// 		Tags:        []string{meth.Service.GetName()},
	// 		Summary:     opOpts.GetSummary(),
	// 		Description: opOpts.GetDescription(),
	// 		OperationID: fg.getOperationName(meth.Service.GetName(), meth.GetName(), bindingIdx),
	// 		Parameters:  pathParams,
	// 		Responses:   openapi3.NewResponses(),
	// 	}
	//
	// 	fg.addMethodResponses(operation, meth)
	//
	// 	if opOpts.GetSecurity() != nil {
	// 		operation.Security = fg.convertSecurity(opOpts.GetSecurity())
	// 	}
	//
	// 	pathTemplate := fg.convertPathTemplate(binding.PathTmpl.Template)
	// 	pathItem := doc.Paths.Find(pathTemplate)
	// 	if pathItem == nil {
	// 		pathItem = &openapi3.PathItem{}
	// 		doc.Paths.Set(pathTemplate, pathItem)
	// 	}
	//
	// 	switch binding.HTTPMethod {
	// 	case "GET":
	// 		pathItem.Get = operation
	// 	case "POST":
	// 		pathItem.Post = operation
	// 	case "PUT":
	// 		pathItem.Put = operation
	// 	case "PATCH":
	// 		pathItem.Patch = operation
	// 	case "DELETE":
	// 		pathItem.Delete = operation
	// 	case "HEAD":
	// 		pathItem.Head = operation
	// 	case "OPTIONS":
	// 		pathItem.Options = operation
	// 	}
	// }
	//
	// return nil
}

func (fg *fileGenerator) generatePathParameters(params []descriptor.Parameter) (any, error) {
	panic("unimplemented")
}

func (fg *fileGenerator) generateResponseSchema(responseType *descriptor.Message) *openapi3.SchemaRef {
	if responseType == nil {
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeObject}},
		}
	}

	return fg.getMessageSchema(responseType)
}
