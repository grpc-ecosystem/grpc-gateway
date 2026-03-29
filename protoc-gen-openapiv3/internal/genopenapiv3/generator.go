package genopenapiv3

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	openapioptions "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/visibility"
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
	adapter        Adapter
}

// wrapper wraps the OpenAPI document with its source file name.
type wrapper struct {
	fileName string
	openapi  *OpenAPI
}

// New returns a new generator which generates OpenAPI v3 files.
// Returns an error if the requested OpenAPI version is not supported.
func New(reg *descriptor.Registry, format Format, openapiVersion string) (gen.Generator, error) {
	registry := NewAdapterRegistry()
	adapter, err := registry.Get(openapiVersion)
	if err != nil {
		return nil, fmt.Errorf("unsupported OpenAPI version %q: %w (supported: %v)", openapiVersion, err, registry.SupportedVersions())
	}
	return &generator{
		reg:            reg,
		format:         format,
		openapiVersion: adapter.Version(),
		adapter:        adapter,
	}, nil
}

// applyNullable marks a schema as nullable.
// The actual output format (type array vs nullable boolean) is handled by the adapter.
func (g *generator) applyNullable(schema *Schema) {
	if schema == nil {
		return
	}
	// Just set the nullable flag - the converter and adapter handle version-specific output
	schema.Nullable = true
}

// Generate implements gen.Generator.
func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile

	if g.reg.IsAllowMerge() {
		// Collect all file-level annotations first
		var documentAnnotations []*descriptor.File
		for _, f := range targets {
			if proto.HasExtension(f.Options, openapioptions.E_Openapiv3Document) {
				documentAnnotations = append(documentAnnotations, f)
			}
		}

		// Merge all files
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

		// Apply annotations from the annotated file(s) to merged target
		// This separates concerns: annotations vs content
		if len(documentAnnotations) > 0 {
			// Use the first annotated file's options
			mergedTarget.Options = documentAnnotations[0].Options
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
		if isVisible(getServiceVisibilityOption(svc), g.reg) {
			for _, method := range svc.Methods {
				if len(method.Bindings) > 0 && isVisible(getMethodVisibilityOption(method), g.reg) {
					hasBindings = true
					break
				}
			}
			if hasBindings {
				break
			}
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
	// referencedSchemas contains FQMNs (with leading dot) or special names like "google.rpc.Status"
	for fqn := range referencedSchemas {
		// Try to look up the message
		msg, err := g.reg.LookupMsg("", fqn)
		if err != nil {
			// Try looking up as enum
			enum, err := g.reg.LookupEnum("", fqn)
			if err != nil {
				continue // Skip if not found (e.g., google.rpc.Status handled separately)
			}
			g.generateEnumSchema(doc, enum)
			continue
		}

		schemaName := g.messageSchemaName(msg)
		if _, exists := doc.Components.Schemas[schemaName]; exists {
			continue
		}
		g.generateMessageSchema(doc, msg, make(map[string]bool))
	}

	// Generate schemas for all messages and enums in the file if include_all_messages is enabled
	if g.reg.GetIncludeAllMessages() {
		g.generateAllMessageSchemas(doc, file)
	}

	// Add default error schema if not disabled
	if !g.reg.GetDisableDefaultErrors() {
		g.addErrorSchema(doc)
	}

	// Apply file-level annotations
	g.applyFileAnnotation(doc, file)

	return doc, nil
}

// generateAllMessageSchemas generates schemas for all messages and enums defined in the file,
// including nested messages. The file.Messages already contains flattened nested messages.
func (g *generator) generateAllMessageSchemas(doc *OpenAPI, file *descriptor.File) {
	visited := make(map[string]bool)

	// Generate schemas for all messages in the file (includes nested messages)
	for _, msg := range file.Messages {
		// Skip synthetic map entry messages
		if msg.GetOptions() != nil && msg.GetOptions().GetMapEntry() {
			continue
		}

		schemaName := g.messageSchemaName(msg)

		// Skip if already generated
		if _, exists := doc.Components.Schemas[schemaName]; exists {
			continue
		}

		// Generate the schema for this message
		g.generateMessageSchema(doc, msg, visited)
	}

	// Generate schemas for all enums in the file (includes nested enums)
	for _, enum := range file.Enums {
		g.generateEnumSchema(doc, enum)
	}
}

// generateServicePaths generates paths for all methods in a service.
func (g *generator) generateServicePaths(doc *OpenAPI, svc *descriptor.Service, referencedSchemas map[string]bool) {
	if !isVisible(getServiceVisibilityOption(svc), g.reg) {
		return
	}

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
	if !isVisible(getMethodVisibilityOption(method), g.reg) {
		return
	}
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
	binding *descriptor.Binding, bindingIdx int, referencedSchemas map[string]bool,
) {
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
	binding *descriptor.Binding, bindingIdx int, referencedSchemas map[string]bool,
) *Operation {
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
		referencedSchemas[method.RequestType.FQMN()] = true
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
		if !isVisible(getFieldVisibilityOption(pathParam.Target), g.reg) {
			continue
		}
		schema := g.fieldToSchemaRef(pathParam.Target, referencedSchemas)
		param := NewPathParameter(pathParam.String(), schema)

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
			if !isVisible(getFieldVisibilityOption(field), g.reg) {
				continue
			}
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

			// Mark as required based on field_behavior annotations
			if g.getFieldRequiredFromBehavior(field) {
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
		referencedSchemas[method.RequestType.FQMN()] = true

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
	if !g.reg.GetDisableDefaultResponses() && method.ResponseType != nil {
		schemaName := g.messageSchemaName(method.ResponseType)
		referencedSchemas[method.ResponseType.FQMN()] = true

		successComment := messageComments(g.reg, method.ResponseType)
		if successComment == "" {
			successComment = "A successful response"
		}

		responses.Codes["200"] = &ResponseRef{
			Value: NewResponse(successComment).WithJSONSchema(NewSchemaRef(schemaName)),
		}
	}

	// Default error response (unless disabled)
	if !g.reg.GetDisableDefaultErrors() {
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
	// Skip synthetic map entry messages - they're handled inline by tryMapFieldSchema
	if msg.GetOptions() != nil && msg.GetOptions().GetMapEntry() {
		return
	}

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
		doc.Components.Schemas[schemaName] = &SchemaOrReference{Schema: wktSchema}
		return
	}

	// Group fields by oneof
	regularFields, oneofGroups := g.groupFieldsByOneof(msg)

	// Build schema for this message
	schema := &Schema{
		Type:       SchemaType{"object"},
		Properties: make(map[string]*SchemaOrReference),
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
		if !isVisible(getFieldVisibilityOption(field), g.reg) {
			continue
		}
		g.addFieldToSchema(doc, schema, field, visited)
	}

	// Generate oneOf constraints for oneof groups
	// For messages with ONLY oneof fields, use pure oneOf without top-level properties
	// For messages with both regular and oneof fields, add oneof fields to properties too
	if len(oneofGroups) > 0 {
		hasRegularFields := len(schema.Properties) > 0

		if hasRegularFields {
			// Add oneof fields to properties when there are also regular fields
			for _, group := range oneofGroups {
				for _, field := range group.fields {
					if !isVisible(getFieldVisibilityOption(field), g.reg) {
						continue
					}
					g.addFieldToSchema(doc, schema, field, visited)
				}
			}
		}

		oneOf, allOf := g.generateOneOfConstraints(doc, schema, oneofGroups, hasRegularFields, visited)
		if oneOf != nil {
			schema.OneOf = oneOf // Single group
		}
		if allOf != nil {
			schema.AllOf = allOf // Multiple groups, each wrapped in oneOf
		}

		// For pure oneof messages, clear type and properties
		if !hasRegularFields && len(oneofGroups) == 1 {
			schema.Type = nil
			schema.Properties = nil
		}
	}

	doc.Components.Schemas[schemaName] = &SchemaOrReference{Schema: schema}
}

// addFieldToSchema adds a single field to the schema's properties.
// It handles field annotations and field_behavior, mutating schema.Required directly.
func (g *generator) addFieldToSchema(doc *OpenAPI, schema *Schema, field *descriptor.Field, visited map[string]bool) {
	fieldSchemaRef := g.fieldToSchemaRef(field, nil)
	fieldName := g.fieldName(field)
	schema.Properties[fieldName] = fieldSchemaRef

	// Add description from field comments
	fieldComment := fieldComments(g.reg, field)
	if fieldComment != "" && fieldSchemaRef.Schema != nil {
		fieldSchemaRef.Schema.Description = fieldComment
	}

	// Apply field-level annotations
	if fieldSchemaRef.Schema != nil {
		g.applyFieldAnnotation(fieldSchemaRef.Schema, field)
	}

	// Apply field_behavior annotations (mutates schema.Required directly)
	if fieldSchemaRef.Schema != nil {
		g.applyFieldBehaviorToSchema(schema, fieldSchemaRef.Schema, field)
	}

	// Apply proto3 optional nullable using version-appropriate method
	if g.reg.GetProto3OptionalNullable() && field.GetProto3Optional() {
		if fieldSchemaRef.Schema != nil {
			g.applyNullable(fieldSchemaRef.Schema)
		}
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

// generateOneOfConstraints generates oneOf constraints for oneof groups.
// For a single group: returns oneOf directly for use in schema.OneOf
// For multiple groups: each group gets its own oneOf, combined via allOf
// This ensures multiple oneof groups are independent (can set one field from each group).
// When hasRegularFields is false (pure oneof message), generates complete field schemas inline.
func (g *generator) generateOneOfConstraints(doc *OpenAPI, parentSchema *Schema, groups []oneofGroup, hasRegularFields bool, visited map[string]bool) (oneOf []*SchemaOrReference, allOf []*SchemaOrReference) {
	if len(groups) == 0 {
		return nil, nil
	}

	// Generate oneOf constraint for each group
	var groupConstraints []*SchemaOrReference
	for _, group := range groups {
		groupOneOf := g.generateSingleGroupOneOf(doc, parentSchema, group, hasRegularFields, visited)
		if len(groupOneOf) > 0 {
			groupConstraints = append(groupConstraints, &SchemaOrReference{
				Schema: &Schema{OneOf: groupOneOf},
			})
		}
	}

	// Single group: use oneOf directly
	// Multiple groups: wrap each in allOf so they're independent
	// No visible groups: return nil for both
	switch len(groupConstraints) {
	case 0:
		return nil, nil
	case 1:
		return groupConstraints[0].Schema.OneOf, nil
	default:
		return nil, groupConstraints
	}
}

// generateSingleGroupOneOf generates oneOf options for a single oneof group.
// When hasRegularFields is true, references properties from parentSchema.
// When hasRegularFields is false, generates complete field schemas inline.
// Includes a "none" option since protobuf oneofs are optional.
func (g *generator) generateSingleGroupOneOf(doc *OpenAPI, parentSchema *Schema, group oneofGroup, hasRegularFields bool, visited map[string]bool) []*SchemaOrReference {
	var options []*SchemaOrReference
	var fieldNames []string

	for _, field := range group.fields {
		if !isVisible(getFieldVisibilityOption(field), g.reg) {
			continue
		}
		fieldName := g.fieldName(field)
		fieldNames = append(fieldNames, fieldName)

		var fieldSchemaRef *SchemaOrReference
		if hasRegularFields {
			// Use existing property from parent schema
			fieldSchemaRef = parentSchema.Properties[fieldName]
		} else {
			// Generate field schema directly for pure oneof messages
			fieldSchemaRef = g.fieldToSchemaRef(field, nil)

			// Add description from field comments
			fieldComment := fieldComments(g.reg, field)
			if fieldComment != "" && fieldSchemaRef.Schema != nil {
				fieldSchemaRef.Schema.Description = fieldComment
			}

			// Apply field-level annotations
			if fieldSchemaRef.Schema != nil {
				g.applyFieldAnnotation(fieldSchemaRef.Schema, field)
			}

			// Generate referenced message/enum schemas
			if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
				if refMsg, err := g.reg.LookupMsg("", field.GetTypeName()); err == nil {
					g.generateMessageSchema(doc, refMsg, visited)
				}
			}
			if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
				if refEnum, err := g.reg.LookupEnum("", field.GetTypeName()); err == nil {
					g.generateEnumSchema(doc, refEnum)
				}
			}
		}

		// Create a schema that requires this specific field from the oneof
		optionSchema := &Schema{
			Type: SchemaType{"object"},
			Properties: map[string]*SchemaOrReference{
				fieldName: fieldSchemaRef,
			},
			Required: []string{fieldName},
		}

		// Add a title to identify which oneof option this is
		optionSchema.Title = fmt.Sprintf("%s.%s", group.name, fieldName)

		options = append(options, &SchemaOrReference{Schema: optionSchema})
	}

	// Add "none" option - oneofs are optional in protobuf
	if len(fieldNames) > 0 {
		neitherSchema := g.buildNeitherSetSchema(group.name, fieldNames)
		options = append(options, &SchemaOrReference{Schema: neitherSchema})
	}

	return options
}

// buildNeitherSetSchema creates a schema that matches when none of the oneof fields are set.
// Uses the pattern: { "not": { "anyOf": [{ "required": ["field1"] }, { "required": ["field2"] }] } }
func (g *generator) buildNeitherSetSchema(groupName string, fieldNames []string) *Schema {
	var anyOfSchemas []*SchemaOrReference
	for _, fieldName := range fieldNames {
		anyOfSchemas = append(anyOfSchemas, &SchemaOrReference{
			Schema: &Schema{
				Required: []string{fieldName},
			},
		})
	}

	return &Schema{
		Title: fmt.Sprintf("%s.none", groupName),
		Not: &SchemaOrReference{
			Schema: &Schema{
				AnyOf: anyOfSchemas,
			},
		},
	}
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
			if !isVisible(getEnumValueVisibilityOption(v), g.reg) || (g.reg.GetOmitEnumDefaultValue() && v.GetNumber() == 0) {
				continue
			}
			enumValues = append(enumValues, int(v.GetNumber()))
		}
	} else {
		for _, v := range enum.GetValue() {
			if !isVisible(getEnumValueVisibilityOption(v), g.reg) || (g.reg.GetOmitEnumDefaultValue() && v.GetNumber() == 0) {
				continue
			}
			enumValues = append(enumValues, v.GetName())
		}
	}

	schema := &Schema{
		Type: SchemaType{"string"},
		Enum: enumValues,
	}

	if g.reg.GetEnumsAsInts() {
		schema.Type = SchemaType{"integer"}
	}

	// Add description from comments
	comment := enumComments(g.reg, enum)
	if comment != "" {
		schema.Description = comment
	}

	// Apply enum-level annotations
	g.applyEnumAnnotation(schema, enum)

	doc.Components.Schemas[schemaName] = &SchemaOrReference{Schema: schema}
}

// addErrorSchema adds the google.rpc.Status schema for error responses.
func (g *generator) addErrorSchema(doc *OpenAPI) {
	// google.rpc.Status schema
	statusSchema := &Schema{
		Type: SchemaType{"object"},
		Properties: map[string]*SchemaOrReference{
			"code": {
				Schema: &Schema{
					Type:        SchemaType{"integer"},
					Format:      "int32",
					Description: "The status code, which should be an enum value of google.rpc.Code.",
				},
			},
			"message": {
				Schema: &Schema{
					Type:        SchemaType{"string"},
					Description: "A developer-facing error message, which should be in English.",
				},
			},
			"details": {
				Schema: &Schema{
					Type: SchemaType{"array"},
					Items: &SchemaOrReference{
						Schema: &Schema{
							Type: SchemaType{"object"},
							Properties: map[string]*SchemaOrReference{
								"@type": {
									Schema: &Schema{
										Type:        SchemaType{"string"},
										Description: "A URL/resource name that uniquely identifies the type of the serialized protocol buffer message.",
									},
								},
							},
							AdditionalProperties: &SchemaOrReference{Schema: &Schema{}},
						},
					},
					Description: "A list of messages that carry the error details.",
				},
			},
		},
		Description: "The `Status` type defines a logical error model that is suitable for different programming environments.",
	}

	doc.Components.Schemas["google.rpc.Status"] = &SchemaOrReference{Schema: statusSchema}
}

// fieldToSchemaRef converts a proto field to a SchemaRef.
func (g *generator) fieldToSchemaRef(field *descriptor.Field, referencedSchemas map[string]bool) *SchemaOrReference {
	// Handle repeated fields (arrays or maps)
	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		// Check if this is a map field (repeated message with map_entry=true)
		if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			if mapSchema := g.tryMapFieldSchema(field, referencedSchemas); mapSchema != nil {
				return mapSchema
			}
		}

		// Regular array
		itemSchema := g.fieldTypeToSchema(field, referencedSchemas)
		return &SchemaOrReference{
			Schema: &Schema{
				Type:  SchemaType{"array"},
				Items: itemSchema,
			},
		}
	}

	return g.fieldTypeToSchema(field, referencedSchemas)
}

// tryMapFieldSchema checks if a repeated message field is a map and returns
// the appropriate object schema with additionalProperties. Returns nil if
// the field is not a map.
func (g *generator) tryMapFieldSchema(field *descriptor.Field, referencedSchemas map[string]bool) *SchemaOrReference {
	// Look up the message type
	msg, err := g.reg.LookupMsg("", field.GetTypeName())
	if err != nil {
		return nil
	}

	// Check if the message has map_entry=true option
	if msg.GetOptions() == nil || !msg.GetOptions().GetMapEntry() {
		return nil
	}

	// This is a map field. Map entries have two fields: key (field 1) and value (field 2)
	// We need to get the value field's type for additionalProperties
	fields := msg.GetField()
	if len(fields) != 2 {
		return nil
	}

	// Find the value field (should be field number 2)
	var valueField *descriptorpb.FieldDescriptorProto
	for _, f := range fields {
		if f.GetNumber() == 2 {
			valueField = f
			break
		}
	}
	if valueField == nil {
		return nil
	}

	// Create a descriptor.Field wrapper for the value field to reuse fieldTypeToSchema
	valueDescField := &descriptor.Field{
		FieldDescriptorProto: valueField,
	}

	// Get the schema for the value type
	valueSchema := g.fieldTypeToSchema(valueDescField, referencedSchemas)

	// Return object schema with additionalProperties
	return &SchemaOrReference{
		Schema: &Schema{
			Type:                 SchemaType{"object"},
			AdditionalProperties: valueSchema,
		},
	}
}

// fieldTypeToSchema converts a field type to a SchemaRef.
func (g *generator) fieldTypeToSchema(field *descriptor.Field, referencedSchemas map[string]bool) *SchemaOrReference {
	// Check for well-known types first
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		typeName := field.GetTypeName()
		if wktSchema := wellKnownTypeSchema(typeName); wktSchema != nil {
			return &SchemaOrReference{Schema: wktSchema}
		}
	}

	switch field.GetType() {
	// String types
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}}

	// Integer types
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"integer"}, Format: "int32"}}

	// 64-bit integers (represented as strings in JSON)
	case descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}, Format: "int64"}}

	// Unsigned integers
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"integer"}, Format: "int64"}}

	// 64-bit unsigned integers (represented as strings in JSON)
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}, Format: "uint64"}}

	// Floating point types
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"number"}, Format: "float"}}

	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"number"}, Format: "double"}}

	// Boolean type
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"boolean"}}}

	// Bytes type
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}, Format: "byte"}}

	// Message type - create reference
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		msg, err := g.reg.LookupMsg("", field.GetTypeName())
		if err != nil {
			return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"object"}}}
		}
		schemaName := g.messageSchemaName(msg)
		if referencedSchemas != nil {
			referencedSchemas[msg.FQMN()] = true
		}
		return NewSchemaRef(schemaName)

	// Enum type
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum, err := g.reg.LookupEnum("", field.GetTypeName())
		if err != nil {
			return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}}
		}
		schemaName := g.enumSchemaName(enum)
		if referencedSchemas != nil {
			referencedSchemas[enum.FQEN()] = true
		}
		return NewSchemaRef(schemaName)

	default:
		return &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}}
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

// encodeOpenAPI converts OpenAPI file obj to ResponseFile using the adapter.
func (g *generator) encodeOpenAPI(file *wrapper) (*descriptor.ResponseFile, error) {
	// Convert to canonical model
	canonicalDoc := file.openapi.ToCanonical()

	// Use adapter to serialize
	content, err := g.adapter.Adapt(canonicalDoc, g.format)
	if err != nil {
		return nil, fmt.Errorf("adapter serialization failed: %w", err)
	}

	name := file.fileName
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.openapi."+string(g.format), base)

	return &descriptor.ResponseFile{
		CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(string(content)),
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

// extractFieldBehavior extracts google.api.field_behavior from a field descriptor.
// This matches the pattern from openapiv2's template.go.
func extractFieldBehavior(fd *descriptorpb.FieldDescriptorProto) []annotations.FieldBehavior {
	if fd.Options == nil {
		return nil
	}
	if !proto.HasExtension(fd.Options, annotations.E_FieldBehavior) {
		return nil
	}
	ext := proto.GetExtension(fd.Options, annotations.E_FieldBehavior)
	opts, ok := ext.([]annotations.FieldBehavior)
	if !ok {
		return nil
	}
	return opts
}

// getFieldBehavior returns field behavior annotations for a descriptor.Field.
func getFieldBehavior(field *descriptor.Field) []annotations.FieldBehavior {
	return extractFieldBehavior(field.FieldDescriptorProto)
}

// applyFieldBehaviorToSchema applies field_behavior annotations to a schema.
// It mutates parentSchema.Required directly when the field should be required.
// This matches the pattern from openapiv2's updateSwaggerObjectFromFieldBehavior.
func (g *generator) applyFieldBehaviorToSchema(parentSchema, fieldSchema *Schema, field *descriptor.Field) {
	// Start with proto3 semantics as default (if enabled)
	required := false
	if g.reg.GetUseProto3FieldSemantics() {
		required = !field.GetProto3Optional() && field.OneofIndex == nil
	}

	// Apply field_behavior annotations (these take precedence)
	behaviors := getFieldBehavior(field)
	for _, fb := range behaviors {
		switch fb {
		case annotations.FieldBehavior_REQUIRED:
			required = true
		case annotations.FieldBehavior_OUTPUT_ONLY:
			fieldSchema.ReadOnly = true
		case annotations.FieldBehavior_OPTIONAL:
			required = false
		case annotations.FieldBehavior_INPUT_ONLY:
			fieldSchema.WriteOnly = true // OpenAPI v3 supports this!
		case annotations.FieldBehavior_IMMUTABLE:
			// No direct OpenAPI mapping
		case annotations.FieldBehavior_FIELD_BEHAVIOR_UNSPECIFIED:
			// No action
		}
	}

	// Set required on parent schema directly
	if required {
		parentSchema.Required = append(parentSchema.Required, g.fieldName(field))
	}
}

// getFieldRequiredFromBehavior determines if a field should be required based on
// field_behavior annotations and proto3 semantics.
func (g *generator) getFieldRequiredFromBehavior(field *descriptor.Field) bool {
	// Start with proto3 semantics as default (if enabled)
	required := false
	if g.reg.GetUseProto3FieldSemantics() {
		required = !field.GetProto3Optional() && field.OneofIndex == nil
	}

	// Apply field_behavior annotations (these take precedence)
	for _, fb := range getFieldBehavior(field) {
		switch fb {
		case annotations.FieldBehavior_REQUIRED:
			required = true
		case annotations.FieldBehavior_OPTIONAL:
			required = false
		}
	}

	return required
}

// Visibility helpers
func getFieldVisibilityOption(fd *descriptor.Field) *visibility.VisibilityRule {
	if fd.Options == nil {
		return nil
	}
	if !proto.HasExtension(fd.Options, visibility.E_FieldVisibility) {
		return nil
	}
	ext := proto.GetExtension(fd.Options, visibility.E_FieldVisibility)
	opts, ok := ext.(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

func getServiceVisibilityOption(fd *descriptor.Service) *visibility.VisibilityRule {
	if fd.Options == nil {
		return nil
	}
	if !proto.HasExtension(fd.Options, visibility.E_ApiVisibility) {
		return nil
	}
	ext := proto.GetExtension(fd.Options, visibility.E_ApiVisibility)
	opts, ok := ext.(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

func getMethodVisibilityOption(fd *descriptor.Method) *visibility.VisibilityRule {
	if fd.Options == nil {
		return nil
	}
	if !proto.HasExtension(fd.Options, visibility.E_MethodVisibility) {
		return nil
	}
	ext := proto.GetExtension(fd.Options, visibility.E_MethodVisibility)
	opts, ok := ext.(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

func getEnumValueVisibilityOption(fd *descriptorpb.EnumValueDescriptorProto) *visibility.VisibilityRule {
	if fd.Options == nil {
		return nil
	}
	if !proto.HasExtension(fd.Options, visibility.E_ValueVisibility) {
		return nil
	}
	ext := proto.GetExtension(fd.Options, visibility.E_ValueVisibility)
	opts, ok := ext.(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

// isVisible checks if a field/RPC is visible based on the visibility restriction
// combined with the `visibility_restriction_selectors`.
// Elements with an overlap on `visibility_restriction_selectors` are visible, those without are not visible.
// Elements without `google.api.VisibilityRule` annotations entirely are always visible.
func isVisible(r *visibility.VisibilityRule, reg *descriptor.Registry) bool {
	if r == nil {
		return true
	}

	// Trim the restriction first - empty/whitespace-only means no restriction
	trimmedRestriction := strings.TrimSpace(r.Restriction)
	if trimmedRestriction == "" {
		return true
	}

	restrictions := strings.Split(trimmedRestriction, ",")
	for _, restriction := range restrictions {
		if reg.GetVisibilityRestrictionSelectors()[strings.TrimSpace(restriction)] {
			return true
		}
	}

	return false
}
