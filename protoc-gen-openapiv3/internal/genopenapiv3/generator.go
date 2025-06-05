package genopenapiv3

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"sync"

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
		doc := g.generateFileDoc(t)

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

func (g *generator) generateFileDoc(file *descriptor.File) *openapi3.T {
	doc, ok := g.extractFileOptions(file)
	if !ok {
		doc = &openapi3.T{
			OpenAPI: "3.0.2",
		}
	}

	if doc.Paths == nil {
		doc.Paths = &openapi3.Paths{}
	}

	for _, svc := range file.Services {
		g.generateServiceDoc(svc, doc)
	}

	components := openapi3.NewComponents()
	components.Schemas = make(openapi3.Schemas)
	doc.Components = &components

	for _, msg := range file.Messages {
		msgName := g.getMessageName(msg.FQMN())
		schemaRef := g.generateMessageSchema(msg, components.Schemas)
		components.Schemas[msgName] = schemaRef
	}

	return doc
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
	tempOneOfsProperties := make(map[int32]openapi3.Schemas)
	for _, field := range msg.Fields {
		fieldDoc := g.generateFieldDoc(field, schemas)
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

	// remove single oneof values
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
		switch g.reg.GetOneOfStrategy() {
		case "oneOf":
			return &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					OneOf: g.generateMessageWithOneOfsSchemas(allOneOfsProperties, properties, msg.GetOneofDecl()),
				},
			}
		default:
			grpclog.Fatal("unknown oneof strategy")
		}
	}

	return &openapi3.SchemaRef{
		Value: schema,
	}
}

func (g *generator) generateMessageWithOneOfsSchemas(allOneOfsProperties map[int32]openapi3.Schemas, properties openapi3.Schemas,
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
		res = append(res, g.generateMessageWithOneOfsSchemas(newAllOneOfsProperties, newProperties, newOneOfs)...)
	}

	return res
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
			Ref: "#/components/schemas/" + openAPIRef,
		}
	default:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{ft.String()}, Format: "UNKNOWN"},
		}
	}
}

func (g *generator) getMessageName(fqmn string) string {
	// TODO: have different naming strategies
	return fqmn[1:]
}

func (g *generator) generateServiceDoc(svc *descriptor.Service, doc *openapi3.T) {
	for _, meth := range svc.Methods {
		g.generateMethodDoc(meth, svc.GetName(), doc)
	}
}

func (g *generator) generateMethodDoc(meth *descriptor.Method, serviceName string, doc *openapi3.T) {
	for bindingIdx, binding := range meth.Bindings {
		// Extract OpenAPI operation options from method descriptor
		opOpts, err := extractOperationOptionFromMethodDescriptor(meth.MethodDescriptorProto)
		if err != nil {
			grpclog.Errorf("error extracting method %s operations: %v", meth.GetName(), err)
		}

		operation := &openapi3.Operation{
			Tags:        []string{serviceName},
			Summary:     opOpts.GetSummary(),
			Description: opOpts.GetDescription(),
			OperationID: g.getOperationName(serviceName, meth.GetName(), bindingIdx),
			Parameters:  openapi3.Parameters{},
			Responses:   openapi3.NewResponses(),
		}

		// Add path parameters for all HTTP methods
		pathParams, err := g.generatePathParameters(binding.PathParams)
		if err != nil {
			grpclog.Errorf("error generating path parameters for method %s: %v", meth.GetName(), err)
		} else {
			operation.Parameters = append(operation.Parameters, pathParams...)
		}

		// Add parameters based on HTTP method
		if meth.RequestType != nil {
			switch binding.HTTPMethod {
			case "GET", "DELETE":
				// For GET and DELETE, add query parameters
				queryParams, err := g.messageToQueryParameters(meth.RequestType, binding.PathParams, binding.Body, binding.HTTPMethod)
				if err != nil {
					grpclog.Errorf("error generating query parameters for method %s: %v", meth.GetName(), err)
				} else {
					operation.Parameters = append(operation.Parameters, queryParams...)
				}
			case "POST", "PUT", "PATCH":
				// For POST, PUT, PATCH, add request body
				operation.RequestBody = g.extractBodyRequest(binding, meth.RequestType)

				// Also add query parameters for fields not in the body
				queryParams, err := g.messageToQueryParameters(meth.RequestType, binding.PathParams, binding.Body, binding.HTTPMethod)
				if err != nil {
					grpclog.Errorf("error generating query parameters for method %s: %v", meth.GetName(), err)
				} else {
					operation.Parameters = append(operation.Parameters, queryParams...)
				}
			}
		}

		// Add responses
		g.addMethodResponses(operation, meth)

		// Add security if needed
		if opOpts.GetSecurity() != nil {
			operation.Security = g.convertSecurity(opOpts.GetSecurity())
		}

		// Get or create PathItem
		pathTemplate := g.convertPathTemplate(binding.PathTmpl.Template)
		pathItem := doc.Paths.Find(pathTemplate)
		if pathItem == nil {
			pathItem = &openapi3.PathItem{}
			doc.Paths.Set(pathTemplate, pathItem)
		}

		// Assign operation to the appropriate HTTP method
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
}

func (g *generator) generateResponseSchema(responseType *descriptor.Message) *openapi3.SchemaRef {
	if responseType == nil {
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeObject}},
		}
	}

	openAPIRef, ok := g.fullyQualifiedNameToOpenAPIName(responseType.FQMN())
	if !ok {
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeObject}},
		}
	}

	return &openapi3.SchemaRef{
		Ref: "#/components/schemas/" + openAPIRef,
	}
}

func (g *generator) getOperationName(serviceName, methodName string, bindingIdx int) string {
	if bindingIdx == 0 {
		return fmt.Sprintf("%s_%s", serviceName, methodName)
	}
	return fmt.Sprintf("%s_%s_%d", serviceName, methodName, bindingIdx)
}

func (g *generator) messageToQueryParameters(message *descriptor.Message, pathParams []descriptor.Parameter, body *descriptor.Body, httpMethod string) ([]*openapi3.ParameterRef, error) {
	var params []*openapi3.ParameterRef

	for _, field := range message.Fields {
		// When body is set to oneof field, we want to skip other fields in the oneof group.
		if isBodySameOneOf(body, field) {
			continue
		}

		if g.reg.GetAllowPatchFeature() && field.GetTypeName() == ".google.protobuf.FieldMask" && field.GetName() == "update_mask" && httpMethod == "PATCH" && len(body.FieldPath) != 0 {
			continue
		}

		p, err := g.queryParams(message, field, "", pathParams, body, g.reg.GetRecursiveDepth())
		if err != nil {
			return nil, err
		}
		params = append(params, p...)
	}
	return params, nil
}

func isBodySameOneOf(body *descriptor.Body, field *descriptor.Field) bool {
	if field.OneofIndex == nil {
		return false
	}

	if body == nil || len(body.FieldPath) == 0 {
		return false
	}

	if body.FieldPath[0].Target.OneofIndex == nil {
		return false
	}

	return *body.FieldPath[0].Target.OneofIndex == *field.OneofIndex
}

// queryParams converts a field to a list of OpenAPI query parameters recursively through the use of nestedQueryParams.
func (g *generator) queryParams(message *descriptor.Message, field *descriptor.Field, prefix string, pathParams []descriptor.Parameter, body *descriptor.Body, recursiveCount int) ([]*openapi3.ParameterRef, error) {
	return g.nestedQueryParams(message, field, prefix, pathParams, body, newCycleChecker(recursiveCount))
}

type cycleChecker struct {
	m     map[string]int
	count int
}

func newCycleChecker(recursive int) *cycleChecker {
	return &cycleChecker{
		m:     make(map[string]int),
		count: recursive,
	}
}

// Check returns whether name is still within recursion toleration
func (c *cycleChecker) Check(name string) bool {
	count, ok := c.m[name]
	count += 1
	isCycle := count > c.count

	if isCycle {
		return false
	}

	// provision map entry if not available
	if !ok {
		c.m[name] = 1
		return true
	}

	c.m[name] = count
	return true
}

func (c *cycleChecker) Branch() *cycleChecker {
	copy := &cycleChecker{
		count: c.count,
		m:     make(map[string]int, len(c.m)),
	}

	copy.m = maps.Clone(c.m)

	return copy
}

// nestedQueryParams converts a field to a list of OpenAPI query parameters recursively.
func (g *generator) nestedQueryParams(message *descriptor.Message, field *descriptor.Field, prefix string, pathParams []descriptor.Parameter, body *descriptor.Body, cycle *cycleChecker) ([]*openapi3.ParameterRef, error) {
	// make sure the parameter is not already listed as a path parameter
	for _, pathParam := range pathParams {
		if pathParam.Target == field {
			return nil, nil
		}
	}

	// make sure the parameter is not already listed as a body parameter
	if body != nil && body.FieldPath != nil {
		for _, fieldPath := range body.FieldPath {
			if fieldPath.Target == field {
				return nil, nil
			}
		}
	}

	fieldType := field.GetTypeName()
	isEnum := field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM
	isRepeated := field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED

	// Check if field has a primitive type schema
	if primitiveSchema, ok := primitiveTypeSchemas[field.GetType()]; ok {
		paramSchema := primitiveSchema
		if isRepeated {
			paramSchema = &openapi3.Schema{
				Type:  &openapi3.Types{openapi3.TypeArray},
				Items: &openapi3.SchemaRef{Value: primitiveSchema},
			}
		}

		param := &openapi3.Parameter{
			Name:     prefix + g.reg.FieldName(field),
			In:       "query",
			Required: false, // Query parameters are typically optional
			Schema:   &openapi3.SchemaRef{Value: paramSchema},
		}

		return []*openapi3.ParameterRef{{Value: param}}, nil
	}

	// Handle message types
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		location := ""
		if ix := strings.LastIndex(field.Message.FQMN(), "."); ix > 0 {
			location = field.Message.FQMN()[0:ix]
		}

		if m, err := g.reg.LookupMsg(location, field.GetTypeName()); err == nil {
			// Handle map types
			if opt := m.GetOptions(); opt != nil && opt.MapEntry != nil && *opt.MapEntry {
				k := m.GetField()[0]
				kType, err := getMapParamKey(k.GetType())
				if err != nil {
					return nil, err
				}
				// This will generate a query in the format map_name[key_type]
				fieldName := fmt.Sprintf("%s[%s]", field.GetName(), kType)
				return []*openapi3.ParameterRef{
					{
						Value: &openapi3.Parameter{
							Name:     prefix + fieldName,
							In:       "query",
							Required: false,
							Schema:   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
						},
					},
				}, nil
			}

			// Handle arrays of objects - currently not supported
			if isRepeated {
				return nil, nil // TODO: currently, mapping object in query parameter is not supported
			}
		}
	}

	// Handle enum types
	if isEnum {
		enum, err := g.reg.LookupEnum("", fieldType)
		if err != nil {
			return nil, fmt.Errorf("unknown enum type %s", fieldType)
		}

		var paramSchema *openapi3.Schema
		if isRepeated { // array of enums
			paramSchema = &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeArray},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
						Enum: listEnumNames(enum),
					},
				},
			}
			if g.reg.GetEnumsAsInts() {
				paramSchema.Items.Value.Type = &openapi3.Types{openapi3.TypeInteger}
				paramSchema.Items.Value.Enum = listEnumNumbers(enum)
			}
		} else {
			paramSchema = &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
				Enum: listEnumNames(enum),
			}
			if g.reg.GetEnumsAsInts() {
				paramSchema.Type = &openapi3.Types{openapi3.TypeInteger}
				paramSchema.Enum = listEnumNumbers(enum)
			}
		}

		param := &openapi3.Parameter{
			Name:     prefix + g.reg.FieldName(field),
			In:       "query",
			Required: false, // Query parameters are typically optional
			Schema:   &openapi3.SchemaRef{Value: paramSchema},
		}

		return []*openapi3.ParameterRef{{Value: param}}, nil
	}

	// nested type, recurse
	msg, err := g.reg.LookupMsg("", fieldType)
	if err != nil {
		return nil, fmt.Errorf("unknown message type %s", fieldType)
	}

	// Check for cyclical message reference
	if ok := cycle.Check(msg.GetName()); !ok {
		return nil, fmt.Errorf("exceeded recursive count (%d) for query parameter %q", cycle.count, fieldType)
	}

	touchedOut := cycle.Branch()
	var params []*openapi3.ParameterRef

	for _, nestedField := range msg.Fields {
		fieldName := g.reg.FieldName(field)
		p, err := g.nestedQueryParams(msg, nestedField, prefix+fieldName+".", pathParams, body, touchedOut)
		if err != nil {
			return nil, err
		}
		params = append(params, p...)
	}
	return params, nil
}

// convertPathTemplate converts gRPC gateway path template to OpenAPI path format
func (g *generator) convertPathTemplate(template string) string {
	// For now, return the template as-is since gRPC gateway templates
	// are already compatible with OpenAPI path format
	return template
}

func (g *generator) extractBodyRequest(binding *descriptor.Binding, requestType *descriptor.Message) *openapi3.RequestBodyRef {
	if binding.Body != nil && binding.Body.FieldPath == nil {
		// Use the entire message as request body
		if requestType != nil {
			return &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: g.generateResponseSchema(requestType),
						},
					},
				},
			}
		}
	} else if binding.Body != nil && len(binding.Body.FieldPath) > 0 {
		// Use specific field as request body
		fieldSchema := g.generateFieldTypeSchema(binding.Body.FieldPath[0].Target.FieldDescriptorProto, make(openapi3.Schemas))
		return &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: fieldSchema,
					},
				},
			},
		}
	}
	return nil
}

// Helper functions that need to be implemented based on your specific needs:

func getMapParamKey(fieldType descriptorpb.FieldDescriptorProto_Type) (string, error) {
	switch fieldType {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "string", nil
	case descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_INT64:
		return "integer", nil
	default:
		return "", fmt.Errorf("unsupported map key type: %v", fieldType)
	}
}

func listEnumNames(enum *descriptor.Enum) []any {
	var names []any
	for _, value := range enum.GetValue() {
		names = append(names, value.GetName())
	}
	return names
}

func listEnumNumbers(enum *descriptor.Enum) []any {
	var numbers []any
	for _, value := range enum.GetValue() {
		numbers = append(numbers, value.GetNumber())
	}
	return numbers
}

// extractOperationOptionFromMethodDescriptor extracts the message of type
// options.Operation from a given proto method's descriptor.
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

// generatePathParameters generates OpenAPI path parameters from gRPC gateway path parameters
func (g *generator) generatePathParameters(pathParams []descriptor.Parameter) ([]*openapi3.ParameterRef, error) {
	var params []*openapi3.ParameterRef

	for _, pathParam := range pathParams {
		if pathParam.Target == nil {
			continue
		}

		schema := g.generateFieldTypeSchema(pathParam.Target.FieldDescriptorProto, make(openapi3.Schemas))

		param := &openapi3.Parameter{
			Name:     pathParam.FieldPath.String(),
			In:       "path",
			Required: true, // Path parameters are always required
			Schema:   schema,
		}

		params = append(params, &openapi3.ParameterRef{Value: param})
	}

	return params, nil
}

// addMethodResponses adds standard HTTP responses to the operation
func (g *generator) addMethodResponses(operation *openapi3.Operation, meth *descriptor.Method) {
	// Add successful response (200)
	operation.Responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: proto.String("Successful response"),
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: g.generateResponseSchema(meth.ResponseType),
				},
			},
		},
	})
}

// convertSecurity converts security requirements from protobuf options to OpenAPI security requirements
func (g *generator) convertSecurity(securityReqs []*options.SecurityRequirement) *openapi3.SecurityRequirements {
	if len(securityReqs) == 0 {
		return nil
	}

	var requirements openapi3.SecurityRequirements

	for _, secReq := range securityReqs {
		requirement := make(openapi3.SecurityRequirement)

		// For now, just add an empty requirement since the exact structure
		// of SecurityRequirement in options may vary
		// This can be expanded based on the actual protobuf definition
		_ = secReq // Use the parameter to avoid unused variable error

		requirements = append(requirements, requirement)
	}

	return &requirements
}
