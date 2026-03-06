package genopenapiv3

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/pluginpb"
)

var errNoTargetService = errors.New("no target service defined in the file")

// generator implements the Generator interface for OpenAPI v3.
type generator struct {
	reg            *descriptor.Registry
	format         Format
	openapiVersion string
}

// wrapper wraps the OpenAPI document with its source file name.
type wrapper struct {
	fileName string
	openapi  *OpenAPI
}

// New returns a new generator which generates OpenAPI v3 files.
func New(reg *descriptor.Registry, format Format, openapiVersion string) gen.Generator {
	return &generator{
		reg:            reg,
		format:         format,
		openapiVersion: openapiVersion,
	}
}

// Generate implements gen.Generator.
func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile

	// Handle merge mode
	if g.reg.IsAllowMerge() {
		var mergedTarget *descriptor.File
		for _, f := range targets {
			if mergedTarget == nil {
				mergedTarget = f
			} else {
				mergedTarget.Enums = append(mergedTarget.Enums, f.Enums...)
				mergedTarget.Messages = append(mergedTarget.Messages, f.Messages...)
				mergedTarget.Services = append(mergedTarget.Services, f.Services...)
			}
		}
		targets = []*descriptor.File{mergedTarget}
	}

	var openapis []*wrapper
	for _, file := range targets {
		if grpclog.V(1) {
			grpclog.Infof("Processing %s", file.GetName())
		}

		doc, err := g.generateFileSpec(file)
		if errors.Is(err, errNoTargetService) {
			if grpclog.V(1) {
				grpclog.Infof("%s: %v", file.GetName(), err)
			}
			continue
		}
		if err != nil {
			return nil, err
		}

		openapis = append(openapis, &wrapper{
			fileName: file.GetName(),
			openapi:  doc,
		})
	}

	if g.reg.IsAllowMerge() && len(openapis) > 0 {
		targetOpenAPI := g.mergeTargetFile(openapis, g.reg.GetMergeFileName())
		if !g.reg.IsPreserveRPCOrder() {
			targetOpenAPI.openapi.Paths.SortAlphabetically()
		}
		f, err := g.encodeOpenAPI(targetOpenAPI)
		if err != nil {
			return nil, fmt.Errorf("failed to encode OpenAPI for %s: %w", g.reg.GetMergeFileName(), err)
		}
		files = append(files, f)
	} else {
		for _, file := range openapis {
			if !g.reg.IsPreserveRPCOrder() {
				file.openapi.Paths.SortAlphabetically()
			}
			f, err := g.encodeOpenAPI(file)
			if err != nil {
				return nil, fmt.Errorf("failed to encode OpenAPI for %s: %w", file.fileName, err)
			}
			files = append(files, f)
		}
	}

	return files, nil
}

// generateFileSpec generates an OpenAPI v3 document for a single proto file.
func (g *generator) generateFileSpec(file *descriptor.File) (*OpenAPI, error) {
	// Check if file has services
	if len(file.Services) == 0 {
		return nil, errNoTargetService
	}

	// Check if any service has HTTP bindings (unless generating unbound methods)
	hasBindings := false
	for _, svc := range file.Services {
		for _, method := range svc.Methods {
			if len(method.Bindings) > 0 {
				hasBindings = true
				break
			}
		}
		if hasBindings {
			break
		}
	}

	if !hasBindings && !g.reg.GetGenerateUnboundMethods() {
		return nil, errNoTargetService
	}

	// Create root document
	doc := NewOpenAPI(file.GetName(), "version not set", g.openapiVersion)

	// Track referenced schemas
	referencedSchemas := make(map[string]bool)

	// Process each service
	for _, svc := range file.Services {
		g.generateServicePaths(doc, svc, referencedSchemas)
	}

	// Generate schema definitions for referenced messages
	for schemaName := range referencedSchemas {
		if _, exists := doc.Components.Schemas[schemaName]; exists {
			continue
		}

		// Try to look up the message
		msg, err := g.reg.LookupMsg("", "."+schemaName)
		if err != nil {
			// Try looking up as enum
			enum, err := g.reg.LookupEnum("", "."+schemaName)
			if err != nil {
				continue // Skip if not found
			}
			g.generateEnumSchema(doc, enum)
			continue
		}
		g.generateMessageSchema(doc, msg, make(map[string]bool))
	}

	// Add default error schema if not disabled
	if !g.reg.GetDisableDefaultErrors() {
		g.addErrorSchema(doc)
	}

	// Apply file-level annotations
	g.applyFileAnnotation(doc, file)

	return doc, nil
}

// generateServicePaths generates paths for all methods in a service.
func (g *generator) generateServicePaths(doc *OpenAPI, svc *descriptor.Service, referencedSchemas map[string]bool) {
	// Add service as a tag (unless disabled)
	if !g.reg.GetDisableServiceTags() {
		tagName := svc.GetName()
		if g.reg.IsIncludePackageInTags() && svc.File.GetPackage() != "" {
			tagName = svc.File.GetPackage() + "." + tagName
		}

		svcComment := serviceComments(g.reg, svc)
		tag := &Tag{
			Name:        tagName,
			Description: svcComment,
		}

		// Apply service-level annotations
		g.applyServiceAnnotation(tag, svc)

		doc.Tags = append(doc.Tags, tag)
	}

	// Generate paths for each method
	for _, method := range svc.Methods {
		g.generateMethodPaths(doc, svc, method, referencedSchemas)
	}
}

// generateMethodPaths generates paths for all HTTP bindings of a method.
func (g *generator) generateMethodPaths(doc *OpenAPI, svc *descriptor.Service, method *descriptor.Method, referencedSchemas map[string]bool) {
	// If no bindings and not generating unbound methods, skip
	if len(method.Bindings) == 0 {
		if !g.reg.GetGenerateUnboundMethods() {
			return
		}
		// Generate a default binding for unbound methods
		g.generateUnboundMethodPath(doc, svc, method, referencedSchemas)
		return
	}

	// Generate path for each binding
	for bindingIdx, binding := range method.Bindings {
		g.generateBindingPath(doc, svc, method, binding, bindingIdx, referencedSchemas)
	}
}

// generateBindingPath generates a single path for an HTTP binding.
func (g *generator) generateBindingPath(doc *OpenAPI, svc *descriptor.Service, method *descriptor.Method,
	binding *descriptor.Binding, bindingIdx int, referencedSchemas map[string]bool) {

	// Convert proto path template to OpenAPI path
	path := convertPathTemplate(binding.PathTmpl.Template)

	// Get or create PathItem for this path
	pathItem := doc.Paths.Get(path)
	if pathItem == nil {
		pathItem = &PathItem{}
		doc.Paths.Set(path, pathItem)
	}

	// Build operation
	op := g.buildOperation(svc, method, binding, bindingIdx, referencedSchemas)

	// Set operation for HTTP method
	pathItem.SetOperation(binding.HTTPMethod, op)
}

// generateUnboundMethodPath generates a path for a method without HTTP bindings.
func (g *generator) generateUnboundMethodPath(doc *OpenAPI, svc *descriptor.Service, method *descriptor.Method, referencedSchemas map[string]bool) {
	// Generate a POST endpoint with the method's full name as the path
	path := fmt.Sprintf("/%s/%s", svc.GetName(), method.GetName())

	pathItem := &PathItem{}
	doc.Paths.Set(path, pathItem)

	// Create a dummy binding for the operation
	op := g.buildUnboundOperation(svc, method, referencedSchemas)
	pathItem.SetOperation("POST", op)
}

// buildOperation builds an Operation object from a method binding.
func (g *generator) buildOperation(svc *descriptor.Service, method *descriptor.Method,
	binding *descriptor.Binding, bindingIdx int, referencedSchemas map[string]bool) *Operation {

	// Extract comments
	comment := methodComments(g.reg, method)
	summary, description := splitSummaryDescription(comment)

	// Build operation ID
	opID := g.buildOperationID(svc, method, bindingIdx)

	// Build tags
	var tags []string
	if !g.reg.GetDisableServiceTags() {
		tagName := svc.GetName()
		if g.reg.IsIncludePackageInTags() && svc.File.GetPackage() != "" {
			tagName = svc.File.GetPackage() + "." + tagName
		}
		tags = []string{tagName}
	}

	// Build parameters
	params := g.buildParameters(method, binding, referencedSchemas)

	// Build request body (for POST, PUT, PATCH)
	var requestBody *RequestBodyRef
	if needsRequestBody(binding.HTTPMethod) && binding.Body != nil {
		requestBody = g.buildRequestBody(method, binding, referencedSchemas)
	}

	// Build responses
	responses := g.buildResponses(method, referencedSchemas)

	// Check deprecation
	deprecated := false
	if g.reg.GetEnableRpcDeprecation() && method.GetOptions() != nil {
		deprecated = method.GetOptions().GetDeprecated()
	}

	op := &Operation{
		Tags:        tags,
		Summary:     summary,
		Description: description,
		OperationID: opID,
		Parameters:  params,
		RequestBody: requestBody,
		Responses:   responses,
		Deprecated:  deprecated,
	}

	// Apply method-level annotations
	g.applyOperationAnnotation(op, method)

	return op
}

// buildUnboundOperation builds an Operation for methods without HTTP bindings.
func (g *generator) buildUnboundOperation(svc *descriptor.Service, method *descriptor.Method, referencedSchemas map[string]bool) *Operation {
	comment := methodComments(g.reg, method)
	summary, description := splitSummaryDescription(comment)

	opID := g.buildOperationID(svc, method, 0)

	var tags []string
	if !g.reg.GetDisableServiceTags() {
		tagName := svc.GetName()
		if g.reg.IsIncludePackageInTags() && svc.File.GetPackage() != "" {
			tagName = svc.File.GetPackage() + "." + tagName
		}
		tags = []string{tagName}
	}

	// Request body is the entire request message
	var requestBody *RequestBodyRef
	if method.RequestType != nil {
		schemaName := g.messageSchemaName(method.RequestType)
		referencedSchemas[schemaName] = true
		requestBody = &RequestBodyRef{
			Value: NewJSONRequestBody(NewSchemaRef(schemaName), true),
		}
	}

	responses := g.buildResponses(method, referencedSchemas)

	deprecated := false
	if g.reg.GetEnableRpcDeprecation() && method.GetOptions() != nil {
		deprecated = method.GetOptions().GetDeprecated()
	}

	op := &Operation{
		Tags:        tags,
		Summary:     summary,
		Description: description,
		OperationID: opID,
		RequestBody: requestBody,
		Responses:   responses,
		Deprecated:  deprecated,
	}

	// Apply method-level annotations
	g.applyOperationAnnotation(op, method)

	return op
}

// buildOperationID generates a unique operation ID.
func (g *generator) buildOperationID(svc *descriptor.Service, method *descriptor.Method, bindingIdx int) string {
	var opID string
	if g.reg.GetSimpleOperationIDs() {
		opID = method.GetName()
	} else {
		opID = fmt.Sprintf("%s_%s", svc.GetName(), method.GetName())
	}
	if bindingIdx > 0 {
		opID = fmt.Sprintf("%s_%d", opID, bindingIdx)
	}
	return opID
}

// buildParameters builds parameters for an operation.
func (g *generator) buildParameters(method *descriptor.Method, binding *descriptor.Binding, referencedSchemas map[string]bool) []*ParameterRef {
	var params []*ParameterRef

	// Add path parameters
	for _, pathParam := range binding.PathParams {
		schema := g.fieldToSchemaRef(pathParam.Target, referencedSchemas)
		param := NewPathParameter(pathParam.FieldPath.String(), schema)

		// Add description from field comments
		comment := fieldComments(g.reg, pathParam.Target)
		if comment != "" {
			param.Description = comment
		}

		params = append(params, &ParameterRef{Value: param})
	}

	// Add query parameters (fields not in path or body)
	if method.RequestType != nil {
		for _, field := range method.RequestType.Fields {
			// Skip path parameters
			if isPathParam(field, binding.PathParams) {
				continue
			}
			// Skip body fields
			if isBodyField(field, binding.Body) {
				continue
			}

			schema := g.fieldToSchemaRef(field, referencedSchemas)
			param := NewQueryParameter(g.fieldName(field), schema)

			// Add description from field comments
			comment := fieldComments(g.reg, field)
			if comment != "" {
				param.Description = comment
			}

			// Check deprecation
			if g.reg.GetEnableFieldDeprecation() && field.GetOptions() != nil {
				param.Deprecated = field.GetOptions().GetDeprecated()
			}

			// Mark as required when using proto3 field semantics
			if g.reg.GetUseProto3FieldSemantics() && g.isFieldRequired(field) {
				param.Required = true
			}

			params = append(params, &ParameterRef{Value: param})
		}
	}

	return params
}

// buildRequestBody builds a request body for an operation.
func (g *generator) buildRequestBody(method *descriptor.Method, binding *descriptor.Binding, referencedSchemas map[string]bool) *RequestBodyRef {
	if method.RequestType == nil {
		return nil
	}

	// If body="*", use the entire request message minus path params
	if len(binding.Body.FieldPath) == 0 {
		// Build schema for body fields only
		schemaName := g.messageSchemaName(method.RequestType)
		referencedSchemas[schemaName] = true

		return &RequestBodyRef{
			Value: NewJSONRequestBody(NewSchemaRef(schemaName), true),
		}
	}

	// If body="field_name", use just that field
	bodyField := binding.Body.FieldPath[len(binding.Body.FieldPath)-1].Target
	schema := g.fieldToSchemaRef(bodyField, referencedSchemas)

	return &RequestBodyRef{
		Value: NewJSONRequestBody(schema, true),
	}
}

// buildResponses builds responses for an operation.
func (g *generator) buildResponses(method *descriptor.Method, referencedSchemas map[string]bool) *Responses {
	responses := NewResponses()

	// Success response (200)
	if method.ResponseType != nil {
		schemaName := g.messageSchemaName(method.ResponseType)
		referencedSchemas[schemaName] = true

		successComment := messageComments(g.reg, method.ResponseType)
		if successComment == "" {
			successComment = "A successful response"
		}

		responses.Codes["200"] = &ResponseRef{
			Value: NewResponse(successComment).WithJSONSchema(NewSchemaRef(schemaName)),
		}
	}

	// Default error response (unless disabled)
	if !g.reg.GetDisableDefaultResponses() {
		responses.Default = &ResponseRef{
			Value: NewResponse("An unexpected error response").WithJSONSchema(NewSchemaRef("google.rpc.Status")),
		}
		referencedSchemas["google.rpc.Status"] = true
	}

	return responses
}

// oneofGroup holds fields belonging to a single oneof declaration.
type oneofGroup struct {
	name   string
	fields []*descriptor.Field
}

// groupFieldsByOneof separates regular fields from oneof fields and groups oneof fields.
// Returns: regular fields, oneof groups
// Note: proto3 optional fields use synthetic oneofs, which we skip.
func (g *generator) groupFieldsByOneof(msg *descriptor.Message) ([]*descriptor.Field, []oneofGroup) {
	var regularFields []*descriptor.Field
	oneofMap := make(map[int32]*oneofGroup)

	for _, field := range msg.Fields {
		oneofIndex := field.GetOneofIndex()
		// Check if field is part of a oneof (and not a proto3 optional synthetic oneof)
		if field.OneofIndex != nil && !field.GetProto3Optional() {
			idx := oneofIndex
			if _, exists := oneofMap[idx]; !exists {
				oneofDecl := msg.GetOneofDecl()[idx]
				oneofMap[idx] = &oneofGroup{
					name:   oneofDecl.GetName(),
					fields: []*descriptor.Field{},
				}
			}
			oneofMap[idx].fields = append(oneofMap[idx].fields, field)
		} else {
			regularFields = append(regularFields, field)
		}
	}

	// Convert map to slice, preserving order by index
	var oneofGroups []oneofGroup
	for i := int32(0); i < int32(len(msg.GetOneofDecl())); i++ {
		if group, exists := oneofMap[i]; exists {
			oneofGroups = append(oneofGroups, *group)
		}
	}

	return regularFields, oneofGroups
}

// generateMessageSchema generates a schema definition for a message.
func (g *generator) generateMessageSchema(doc *OpenAPI, msg *descriptor.Message, visited map[string]bool) {
	schemaName := g.messageSchemaName(msg)

	// Prevent infinite recursion for self-referential messages
	if visited[schemaName] {
		return
	}
	visited[schemaName] = true
	defer func() { visited[schemaName] = false }()

	// Check if already generated
	if _, exists := doc.Components.Schemas[schemaName]; exists {
		return
	}

	// Handle well-known types
	if wktSchema := wellKnownTypeSchema(msg.FQMN()); wktSchema != nil {
		doc.Components.Schemas[schemaName] = &SchemaRef{Value: wktSchema}
		return
	}

	// Group fields by oneof
	regularFields, oneofGroups := g.groupFieldsByOneof(msg)

	// Build schema for this message
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]*SchemaRef),
	}

	// Add description from comments
	comment := messageComments(g.reg, msg)
	if comment != "" {
		schema.Description = comment
	}

	// Apply message-level annotations
	g.applySchemaAnnotation(schema, msg)

	// Process regular fields (not in oneof)
	for _, field := range regularFields {
		g.addFieldToSchema(doc, schema, field, visited)

		// Track required fields when using proto3 field semantics
		if g.reg.GetUseProto3FieldSemantics() && g.isFieldRequired(field) {
			schema.Required = append(schema.Required, g.fieldName(field))
		}
	}

	// Process oneof groups - add all oneof fields as properties too
	// (they need to be in properties for JSON serialization)
	// Note: oneof fields are not required since only one can be set
	for _, group := range oneofGroups {
		for _, field := range group.fields {
			g.addFieldToSchema(doc, schema, field, visited)
		}
	}

	// Generate oneOf constraint for each oneof group
	if len(oneofGroups) > 0 {
		schema.OneOf = g.generateOneOfSchemas(schema, oneofGroups)
	}

	doc.Components.Schemas[schemaName] = &SchemaRef{Value: schema}
}

// addFieldToSchema adds a single field to the schema's properties.
func (g *generator) addFieldToSchema(doc *OpenAPI, schema *Schema, field *descriptor.Field, visited map[string]bool) {
	fieldSchemaRef := g.fieldToSchemaRef(field, nil)
	fieldName := g.fieldName(field)
	schema.Properties[fieldName] = fieldSchemaRef

	// Add description from field comments
	fieldComment := fieldComments(g.reg, field)
	if fieldComment != "" && fieldSchemaRef.Value != nil {
		fieldSchemaRef.Value.Description = fieldComment
	}

	// Apply field-level annotations
	if fieldSchemaRef.Value != nil {
		g.applyFieldAnnotation(fieldSchemaRef.Value, field)
	}

	// If field references another message, generate that too
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		if refMsg, err := g.reg.LookupMsg("", field.GetTypeName()); err == nil {
			g.generateMessageSchema(doc, refMsg, visited)
		}
	}

	// If field references an enum, generate that too
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
		if refEnum, err := g.reg.LookupEnum("", field.GetTypeName()); err == nil {
			g.generateEnumSchema(doc, refEnum)
		}
	}
}

// generateOneOfSchemas generates oneOf schemas for oneof field groups.
// Each oneof option becomes a schema requiring exactly that field.
func (g *generator) generateOneOfSchemas(parentSchema *Schema, groups []oneofGroup) []*SchemaRef {
	var oneOfSchemas []*SchemaRef

	for _, group := range groups {
		// For each field in the oneof, create a schema that requires only that field
		for _, field := range group.fields {
			fieldName := g.fieldName(field)

			// Create a schema that requires this specific field from the oneof
			optionSchema := &Schema{
				Type: "object",
				Properties: map[string]*SchemaRef{
					fieldName: parentSchema.Properties[fieldName],
				},
				Required: []string{fieldName},
			}

			// Add a title to identify which oneof option this is
			optionSchema.Title = fmt.Sprintf("%s.%s", group.name, fieldName)

			oneOfSchemas = append(oneOfSchemas, &SchemaRef{Value: optionSchema})
		}
	}

	return oneOfSchemas
}

// generateEnumSchema generates a schema definition for an enum.
func (g *generator) generateEnumSchema(doc *OpenAPI, enum *descriptor.Enum) {
	schemaName := g.enumSchemaName(enum)

	// Check if already generated
	if _, exists := doc.Components.Schemas[schemaName]; exists {
		return
	}

	// Collect enum values
	var enumValues []any
	if g.reg.GetEnumsAsInts() {
		for _, v := range enum.GetValue() {
			if g.reg.GetOmitEnumDefaultValue() && v.GetNumber() == 0 {
				continue
			}
			enumValues = append(enumValues, int(v.GetNumber()))
		}
	} else {
		for _, v := range enum.GetValue() {
			if g.reg.GetOmitEnumDefaultValue() && v.GetNumber() == 0 {
				continue
			}
			enumValues = append(enumValues, v.GetName())
		}
	}

	schema := &Schema{
		Type: "string",
		Enum: enumValues,
	}

	if g.reg.GetEnumsAsInts() {
		schema.Type = "integer"
	}

	// Add description from comments
	comment := enumComments(g.reg, enum)
	if comment != "" {
		schema.Description = comment
	}

	// Apply enum-level annotations
	g.applyEnumAnnotation(schema, enum)

	doc.Components.Schemas[schemaName] = &SchemaRef{Value: schema}
}

// addErrorSchema adds the google.rpc.Status schema for error responses.
func (g *generator) addErrorSchema(doc *OpenAPI) {
	// google.rpc.Status schema
	statusSchema := &Schema{
		Type: "object",
		Properties: map[string]*SchemaRef{
			"code": {
				Value: &Schema{
					Type:        "integer",
					Format:      "int32",
					Description: "The status code, which should be an enum value of google.rpc.Code.",
				},
			},
			"message": {
				Value: &Schema{
					Type:        "string",
					Description: "A developer-facing error message, which should be in English.",
				},
			},
			"details": {
				Value: &Schema{
					Type: "array",
					Items: &SchemaRef{
						Value: &Schema{
							Type: "object",
							Properties: map[string]*SchemaRef{
								"@type": {
									Value: &Schema{
										Type:        "string",
										Description: "A URL/resource name that uniquely identifies the type of the serialized protocol buffer message.",
									},
								},
							},
							AdditionalProperties: &SchemaRef{Value: &Schema{}},
						},
					},
					Description: "A list of messages that carry the error details.",
				},
			},
		},
		Description: "The `Status` type defines a logical error model that is suitable for different programming environments.",
	}

	doc.Components.Schemas["google.rpc.Status"] = &SchemaRef{Value: statusSchema}
}

// fieldToSchemaRef converts a proto field to a SchemaRef.
func (g *generator) fieldToSchemaRef(field *descriptor.Field, referencedSchemas map[string]bool) *SchemaRef {
	// Handle repeated fields (arrays)
	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		itemSchema := g.fieldTypeToSchema(field, referencedSchemas)
		return &SchemaRef{
			Value: &Schema{
				Type:  "array",
				Items: itemSchema,
			},
		}
	}

	return g.fieldTypeToSchema(field, referencedSchemas)
}

// fieldTypeToSchema converts a field type to a SchemaRef.
func (g *generator) fieldTypeToSchema(field *descriptor.Field, referencedSchemas map[string]bool) *SchemaRef {
	// Check for well-known types first
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		typeName := field.GetTypeName()
		if wktSchema := wellKnownTypeSchema(typeName); wktSchema != nil {
			return &SchemaRef{Value: wktSchema}
		}
	}

	switch field.GetType() {
	// String types
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return &SchemaRef{Value: &Schema{Type: "string"}}

	// Integer types
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return &SchemaRef{Value: &Schema{Type: "integer", Format: "int32"}}

	// 64-bit integers (represented as strings in JSON)
	case descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return &SchemaRef{Value: &Schema{Type: "string", Format: "int64"}}

	// Unsigned integers
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return &SchemaRef{Value: &Schema{Type: "integer", Format: "int64"}}

	// 64-bit unsigned integers (represented as strings in JSON)
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return &SchemaRef{Value: &Schema{Type: "string", Format: "uint64"}}

	// Floating point types
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return &SchemaRef{Value: &Schema{Type: "number", Format: "float"}}

	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return &SchemaRef{Value: &Schema{Type: "number", Format: "double"}}

	// Boolean type
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return &SchemaRef{Value: &Schema{Type: "boolean"}}

	// Bytes type
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return &SchemaRef{Value: &Schema{Type: "string", Format: "byte"}}

	// Message type - create reference
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		msg, err := g.reg.LookupMsg("", field.GetTypeName())
		if err != nil {
			return &SchemaRef{Value: &Schema{Type: "object"}}
		}
		schemaName := g.messageSchemaName(msg)
		if referencedSchemas != nil {
			referencedSchemas[schemaName] = true
		}
		return NewSchemaRef(schemaName)

	// Enum type
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum, err := g.reg.LookupEnum("", field.GetTypeName())
		if err != nil {
			return &SchemaRef{Value: &Schema{Type: "string"}}
		}
		schemaName := g.enumSchemaName(enum)
		if referencedSchemas != nil {
			referencedSchemas[schemaName] = true
		}
		return NewSchemaRef(schemaName)

	default:
		return &SchemaRef{Value: &Schema{Type: "string"}}
	}
}

// registriesSeen holds memoized OpenAPI name mappings per registry.
// This ensures consistent naming across the entire generation run.
var (
	registriesSeen      = map[*descriptor.Registry]map[string]string{}
	registriesSeenMutex sync.Mutex
)

// resolveOpenAPIName resolves a fully-qualified name to an OpenAPI name
// using the configured naming strategy. Results are memoized per registry.
func (g *generator) resolveOpenAPIName(fqn string) string {
	registriesSeenMutex.Lock()
	defer registriesSeenMutex.Unlock()

	if mapping, present := registriesSeen[g.reg]; present {
		if name, ok := mapping[fqn]; ok {
			return name
		}
	}

	// Collect all FQMNs and FQENs
	allNames := append(g.reg.GetAllFQMNs(), g.reg.GetAllFQENs()...)
	strategy := g.reg.GetOpenAPINamingStrategy()
	if strategy == "" {
		strategy = "legacy" // Default matches v2 behavior
	}

	mapping := resolveFullyQualifiedNameToOpenAPINames(allNames, strategy)
	registriesSeen[g.reg] = mapping

	return mapping[fqn]
}

// messageSchemaName returns the schema name for a message.
func (g *generator) messageSchemaName(msg *descriptor.Message) string {
	return g.resolveOpenAPIName(msg.FQMN())
}

// enumSchemaName returns the schema name for an enum.
func (g *generator) enumSchemaName(enum *descriptor.Enum) string {
	return g.resolveOpenAPIName(enum.FQEN())
}

// fieldName returns the JSON field name for a proto field.
func (g *generator) fieldName(field *descriptor.Field) string {
	if g.reg.GetUseJSONNamesForFields() {
		return field.GetJsonName()
	}
	return field.GetName()
}

// mergeTargetFile merges multiple OpenAPI files into one.
func (g *generator) mergeTargetFile(targets []*wrapper, mergeFileName string) *wrapper {
	var mergedTarget *wrapper
	for _, f := range targets {
		if mergedTarget == nil {
			mergedTarget = &wrapper{
				fileName: mergeFileName,
				openapi:  f.openapi,
			}
		} else {
			// Merge schemas
			for k, v := range f.openapi.Components.Schemas {
				mergedTarget.openapi.Components.Schemas[k] = v
			}
			// Merge paths
			for _, path := range f.openapi.Paths.order {
				item := f.openapi.Paths.paths[path]
				mergedTarget.openapi.Paths.Set(path, item)
			}
			// Merge tags
			mergedTarget.openapi.Tags = append(mergedTarget.openapi.Tags, f.openapi.Tags...)
		}
	}
	return mergedTarget
}

// encodeOpenAPI converts OpenAPI file obj to ResponseFile.
func (g *generator) encodeOpenAPI(file *wrapper) (*descriptor.ResponseFile, error) {
	var contentBuf bytes.Buffer
	enc, err := g.format.NewEncoder(&contentBuf)
	if err != nil {
		return nil, err
	}

	if err := enc.Encode(file.openapi); err != nil {
		return nil, err
	}

	name := file.fileName
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.openapi." + string(g.format), base)

	return &descriptor.ResponseFile{
		CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(contentBuf.String()),
		},
	}, nil
}

// AddErrorDefs adds google.rpc.Status and google.protobuf.Any to registry.
func AddErrorDefs(reg *descriptor.Registry) error {
	anyProto := protodesc.ToFileDescriptorProto((&anypb.Any{}).ProtoReflect().Descriptor().ParentFile())
	anyProto.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	status := protodesc.ToFileDescriptorProto((&statuspb.Status{}).ProtoReflect().Descriptor().ParentFile())
	status.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	return reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			anyProto,
			status,
		},
	})
}

// pathParamPattern matches {param} or {param=pattern}.
var pathParamPattern = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_.]*)(=[^}]*)?\}`)

// convertPathTemplate converts proto path template to OpenAPI format.
// Input:  /v1/users/{user_id}/posts/{post.id=posts/*}
// Output: /v1/users/{user_id}/posts/{post.id}
func convertPathTemplate(template string) string {
	return pathParamPattern.ReplaceAllString(template, "{$1}")
}

// Helper functions

func needsRequestBody(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}

func isPathParam(field *descriptor.Field, pathParams []descriptor.Parameter) bool {
	for _, p := range pathParams {
		if p.Target == field {
			return true
		}
	}
	return false
}

func isBodyField(field *descriptor.Field, body *descriptor.Body) bool {
	if body == nil {
		return false
	}
	// body="*" means entire message is body
	if len(body.FieldPath) == 0 {
		return true
	}
	// Check if field is in body field path
	for _, fp := range body.FieldPath {
		if fp.Target == field {
			return true
		}
	}
	return false
}

// isFieldRequired determines if a field should be marked as required.
// In proto3, a field is required if:
// - It's NOT a proto3 optional field (which explicitly tracks presence)
// - It's NOT part of a oneof (handled separately, only one can be set)
// - It's NOT a repeated field (empty array is equivalent to absent)
func (g *generator) isFieldRequired(field *descriptor.Field) bool {
	// Proto3 optional fields explicitly track presence, so can be absent
	if field.GetProto3Optional() {
		return false
	}
	// Oneof fields are mutually exclusive, not required individually
	if field.OneofIndex != nil {
		return false
	}
	// Repeated fields default to empty, so are effectively optional
	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return false
	}
	return true
}
