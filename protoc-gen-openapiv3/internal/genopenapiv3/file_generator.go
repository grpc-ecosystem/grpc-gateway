package genopenapiv3

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

type fileGenerator struct {
	*generator

	spec *openapi3.T
}

func (fg *fileGenerator) generateFileSpec(file *descriptor.File) *openapi3.T {

	fg.spec = convertFileOptions(file)

	fg.spec.Components = &openapi3.Components{}
	fg.spec.Components.Schemas = make(openapi3.Schemas)
	fg.spec.Components.RequestBodies = make(openapi3.RequestBodies)

	if fg.spec.Paths == nil {
		fg.spec.Paths = &openapi3.Paths{}
	}

	for _, svc := range file.Services {
		if fg.IsServiceVisible(svc) {
			err := fg.generateServiceSpec(svc)
			if err != nil {
				grpclog.Errorf("could not generate service document: %v", err)
			}
		}
	}

	for _, msg := range file.Messages {
		if fg.IsMessageVisible(msg) {
			fg.generateMessageSchemaRef(msg, map[string]bool{})
		}
	}

	for _, enum := range file.Enums {
		if fg.IsEnumVisible(enum) {
			fg.generateEnumSchemaRef(enum)
		}
	}

	return fg.spec
}

func (fg *fileGenerator) generateMessageSchemaRef(msg *descriptor.Message, generationPath map[string]bool) *openapi3.SchemaRef {
	name := fg.resolveName(msg.FQMN())
	resultRef := openapi3.NewSchemaRef(fmt.Sprintf("#/components/schemas/%s", name), nil)

	_, ok := fg.spec.Components.Schemas[name]
	if ok {
		return resultRef
	}

	_, ok = generationPath[msg.FQMN()]
	if ok {
		return resultRef
	}

	generationPath[msg.FQMN()] = true

	fg.spec.Components.Schemas[name] = fg.generateMessageSchema(msg, nil, generationPath).NewRef()

	generationPath[msg.FQMN()] = false

	return resultRef
}

func (fg *fileGenerator) generateMessageSchema(msg *descriptor.Message, excludeFields []string, generationPath map[string]bool) *openapi3.Schema {
	if scheme, ok := wktSchemas[msg.FQMN()]; ok {
		return scheme
	}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeObject},
	}

	properties := make(openapi3.Schemas)
	tempOneOfsProperties := make(map[int32]openapi3.Schemas)
	for _, field := range msg.Fields {
		if slices.Contains(excludeFields, field.FQFN()) {
			continue
		}

		fieldDoc := fg.generateFieldDoc(field, generationPath)
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
				OneOf: fg.generateMessageWithOneOfsSchemas(allOneOfsProperties, properties, msg.GetOneofDecl(), ""),
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
	oneOfs []*descriptorpb.OneofDescriptorProto, namePrefix string) openapi3.SchemaRefs {
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
		res = append(res, fg.generateMessageWithOneOfsSchemas(newAllOneOfsProperties, newProperties, newOneOfs, namePrefix+fieldName)...)
	}

	return res
}

func (fg *fileGenerator) generateFieldDoc(field *descriptor.Field, generationPath map[string]bool) *openapi3.SchemaRef {
	location := fg.fqmnToLocation(field.Message.FQMN())
	if m, err := fg.reg.LookupMsg(location, field.GetTypeName()); err == nil {
		if opt := m.GetOptions(); opt != nil && opt.MapEntry != nil && *opt.MapEntry {
			FieldDesc := m.GetField()[1]

			return &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeObject},
					AdditionalProperties: openapi3.AdditionalProperties{
						Schema: fg.generateFieldTypeSchema(FieldDesc, location, generationPath),
					},
				},
			}
		}
	}

	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:  &openapi3.Types{openapi3.TypeArray},
				Items: fg.generateFieldTypeSchema(field.FieldDescriptorProto, location, generationPath),
			},
		}
	}

	return fg.generateFieldTypeSchema(field.FieldDescriptorProto, location, generationPath)
}

func (fg *fileGenerator) generateFieldTypeSchema(fd *descriptorpb.FieldDescriptorProto, location string, generationPath map[string]bool) *openapi3.SchemaRef {
	if schema, ok := primitiveTypeSchemas[fd.GetType()]; ok {
		return schema.NewRef()
	}

	switch ft := fd.GetType(); ft {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		openAPIRef := fg.resolveType(fd.GetTypeName())
		if schema, ok := fg.spec.Components.Schemas[openAPIRef]; ok {
			return schema
		} else {
			if fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
				fieldTypeEnum, err := fg.reg.LookupEnum(location, fd.GetTypeName())
				if err != nil {
					panic(err)
				}

				return fg.generateEnumSchemaRef(fieldTypeEnum)
			} else {
				fieldTypeMsg, err := fg.reg.LookupMsg(location, fd.GetTypeName())
				if err != nil {
					panic(err)
				}
				return fg.generateMessageSchemaRef(fieldTypeMsg, generationPath)
			}
		}
	default:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{ft.String()}, Format: "UNKNOWN"},
		}
	}
}

func (fg *fileGenerator) generateEnumSchemaRef(enum *descriptor.Enum) *openapi3.SchemaRef {
	name := fg.resolveName(enum.FQEN())

	schemaRef, ok := fg.spec.Components.Schemas[name]
	if ok {
		return schemaRef
	}

	schemaRef = fg.generateEnumSchema(enum).NewRef()
	fg.spec.Components.Schemas[name] = schemaRef

	return schemaRef
}

func (fg *fileGenerator) generateEnumSchema(enum *descriptor.Enum) *openapi3.Schema {
	var enumValues []any
	for _, value := range enum.GetValue() {
		enumValues = append(enumValues, value.GetName())
	}

	return &openapi3.Schema{
		Type: &openapi3.Types{openapi3.TypeString},
		Enum: enumValues,
	}
}

func (fg *fileGenerator) generateServiceSpec(svc *descriptor.Service) error {
	for _, meth := range svc.Methods {
		if fg.IsMethodVisible(meth) {
			err := fg.generateMethodDoc(meth)
			if err != nil {
				return fmt.Errorf("could not generate method %s doc: %w", meth.GetName(), err)
			}
		}
	}

	return nil
}

func (fg *fileGenerator) generateMethodDoc(meth *descriptor.Method) error {
	for bindingIdx, binding := range meth.Bindings {
		opOpts, err := extractOperationOptionFromMethodDescriptor(meth.MethodDescriptorProto)
		if err != nil {
			return fmt.Errorf("error extracting method %s operations: %v", meth.GetName(), err)
		}

		var params openapi3.Parameters
		var requestBody *openapi3.RequestBodyRef

		if meth.RequestType != nil {
			tmpParams, err := fg.messageToParameters(meth.RequestType, binding.PathParams, binding.Body,
				binding.HTTPMethod, "")
			if err != nil {
				grpclog.Errorf("error generating query parameters for method %s: %v", meth.GetName(), err)
			} else {
				params = append(params, tmpParams...)
			}

			pathParamsFQFNs := make([]string, len(binding.PathParams))
			for i, param := range binding.PathParams {
				pathParamsFQFNs[i] = param.Target.FQFN()
			}

			switch binding.HTTPMethod {
			case "POST", "PUT", "PATCH":
				// For POST, PUT, PATCH, add request body
				if len(pathParamsFQFNs) > 0 {
					messageSchema := fg.generateMessageSchema(meth.RequestType, pathParamsFQFNs, map[string]bool{})

					name := fg.resolveName(meth.RequestType.FQMN())
					resultRef := fmt.Sprintf("#/components/requestBodies/%s", name)

					fg.spec.Components.RequestBodies[name] = &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithContent(
						openapi3.NewContentWithJSONSchemaRef(messageSchema.NewRef()))}

					requestBody = &openapi3.RequestBodyRef{Ref: resultRef}
				} else {
					requestBody = &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().
						WithJSONSchemaRef(fg.generateMessageSchemaRef(meth.RequestType, map[string]bool{}))}
				}
			}
		}

		defaultResponse, err := fg.defaultResponse()
		if err != nil {
			return fmt.Errorf("could not get default response: %w", err)
		}

		successResponseSchema := openapi3.NewResponse().
			WithJSONSchemaRef(fg.generateMessageSchemaRef(meth.ResponseType, map[string]bool{}))

		defaultResponseSchema := openapi3.NewResponse().
			WithJSONSchemaRef(fg.generateMessageSchemaRef(defaultResponse, map[string]bool{}))

		responses := openapi3.NewResponses(openapi3.WithStatus(200, &openapi3.ResponseRef{Value: successResponseSchema}),
			openapi3.WithName("default", defaultResponseSchema))

		operation := &openapi3.Operation{
			Tags:        []string{meth.Service.GetName()},
			Summary:     opOpts.GetSummary(),
			Description: opOpts.GetDescription(),
			OperationID: fg.getOperationName(meth.Service.GetName(), meth.GetName(), bindingIdx),
			RequestBody: requestBody,
			Responses:   responses,
			Parameters:  params,
		}

		if opOpts.GetSecurity() != nil {
			operation.Security = convertSecurityRequiremnt(opOpts.GetSecurity())
		}

		path := fg.convertPathTemplate(binding.PathTmpl.Template)
		pathItem := fg.spec.Paths.Find(path)
		if pathItem == nil {
			pathItem = &openapi3.PathItem{}
			fg.spec.Paths.Set(path, pathItem)
		}

		switch binding.HTTPMethod {
		case "GET":
			pathItem.Get = operation
		case "POST":
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
	}

	return nil
}

func (fg *fileGenerator) convertPathTemplate(template string) string {
	// TODO: handle /{arg=foo/*} and /{arg=foo/**}
	return template
}

func (fg *fileGenerator) messageToParameters(message *descriptor.Message,
	pathParams []descriptor.Parameter, body *descriptor.Body,
	httpMethod string, paramPrefix string) (openapi3.Parameters, error) {

	// TODO: remove this after handling nested path parameters
	_ = paramPrefix

	params := openapi3.NewParameters()
	for _, field := range message.Fields {
		paramType, isParam := fg.getParamType(field, pathParams, body, message, httpMethod)
		if !isParam {
			// TODO: handle nested path parameter reference
			continue
		}

		schema := fg.generateFieldTypeSchema(field.FieldDescriptorProto, fg.fqmnToLocation(field.FQFN()), map[string]bool{})

		switch paramType {
		case openapi3.ParameterInPath:
			param := openapi3.NewPathParameter(field.GetName())
			param.Schema = schema
			params = append(params, &openapi3.ParameterRef{
				Value: param,
			})
		case openapi3.ParameterInQuery:
			param := openapi3.NewQueryParameter(field.GetName())
			param.Schema = schema
			params = append(params, &openapi3.ParameterRef{
				Value: param,
			})
		}
	}

	return params, nil
}

func (fg *fileGenerator) getParamType(field *descriptor.Field, pathParams []descriptor.Parameter, body *descriptor.Body,
	message *descriptor.Message, httpMethod string) (string, bool) {

	for _, pathParam := range pathParams {
		if pathParam.Target.FQFN() == field.FQFN() {
			return openapi3.ParameterInPath, true
		}

		if strings.HasSuffix(pathParam.Target.FQFN(), message.FQMN()) {
			return "", false
		}
	}

	if httpMethod == "GET" || httpMethod == "DELETE" {
		return openapi3.ParameterInQuery, true
	}

	if body == nil {
		return openapi3.ParameterInQuery, true
	}

	if len(body.FieldPath) == 0 {
		return "", false
	}

	if body.FieldPath[len(body.FieldPath)-1].Target.FQFN() == field.FQFN() {
		return "", false
	}

	return openapi3.ParameterInQuery, true
}
