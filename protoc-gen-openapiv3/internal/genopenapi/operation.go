package genopenapi

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/descriptorpb"
)

// statusSchemaName is the component name we use for the auto-injected
// google.rpc.Status default error response schema.
const statusSchemaName = "google.rpc.Status"

// buildOperation produces an OpenAPI Operation for one HTTP binding of an
// RPC method, registering any referenced schemas with the schema builder.
//
// `bindingIdx` disambiguates operationId when a method has multiple bindings;
// the first binding gets the bare ID, subsequent bindings append `_<idx>`.
// `pathParams` is the synthetic OpenAPI path parameter list produced by
// convertPathTemplate; it may be longer than binding.PathParams when a
// single proto field expanded into multiple wildcards.
func buildOperation(b *schemaBuilder, svc *descriptor.Service, m *descriptor.Method, binding *descriptor.Binding, bindingIdx int, pathParams []pathParam) (*Operation, error) {
	summary, description := splitSummaryDescription(methodComments(m))

	op := &Operation{
		Summary:     summary,
		Description: description,
		OperationID: operationID(svc, m, bindingIdx),
		Tags:        []string{svc.GetName()},
		Deprecated:  methodDeprecated(m),
	}

	op.Parameters = buildParameters(b, m, binding, pathParams)

	if needsRequestBody(binding.HTTPMethod) && binding.Body != nil {
		op.RequestBody = buildRequestBody(b, m, binding)
	}

	op.Responses = buildResponses(b, m)
	if o, ok := methodOperationAnnotation(m); ok {
		if err := applyOperationOverride(op, o); err != nil {
			return nil, fmt.Errorf("openapiv3 operation %s.%s: %w", svc.GetName(), m.GetName(), err)
		}
	}
	return op, nil
}

// buildParameters returns the path and query parameters for an operation.
// Body fields are excluded; path-bound fields are emitted as `in: path` and
// everything else from the request type is emitted as `in: query`.
//
// Path parameters are driven by `pathParams` (the synthetic OpenAPI list from
// convertPathTemplate) rather than by `binding.PathParams` directly, so a
// single proto field that expanded into multiple URL wildcards yields one
// OpenAPI parameter per wildcard. The proto field is looked up via the
// dotted name on each pathParam so the schema, comments, and deprecation
// flag come from the original descriptor.
func buildParameters(b *schemaBuilder, m *descriptor.Method, binding *descriptor.Binding, pathParams []pathParam) []*ParameterRef {
	var params []*ParameterRef

	byField := make(map[string]descriptor.Parameter, len(binding.PathParams))
	for _, p := range binding.PathParams {
		byField[p.FieldPath.String()] = p
	}

	for _, pp := range pathParams {
		proto, ok := byField[pp.fieldName]
		if !ok {
			// Should not happen: the path template referenced a field that
			// the binding did not record. Emit an open string parameter so
			// the operation is at least syntactically valid.
			params = append(params, &ParameterRef{Value: &Parameter{
				Name:     pp.openAPIName,
				In:       "path",
				Required: true,
				Schema:   &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}}},
			}})
			grpclog.Infof("unexpected path parameter not found in binding: %q", pp.fieldName)
			continue
		}
		param := &Parameter{
			Name:       pp.openAPIName,
			In:         "path",
			Required:   true,
			Schema:     b.fieldSchema(proto.Target),
			Deprecated: fieldDeprecated(proto.Target),
		}
		if desc := fieldComments(b.reg, proto.Target); desc != "" {
			param.Description = desc
		}
		params = append(params, &ParameterRef{Value: param})
	}

	if m.RequestType == nil {
		return params
	}
	cycle := newQueryCycleChecker(b.reg.GetRecursiveDepth())
	for _, field := range m.RequestType.Fields {
		if !isVisible(fieldVisibility(field), b.reg) {
			// Field is hidden by visibility rules, skip.
			continue
		}
		if isPathParam(field, binding.PathParams) {
			// Already handled as a path parameter, skip.
			continue
		}
		if isBodyField(field, binding.Body) {
			// This field is part of the request body, skip.
			continue
		}
		params = append(params, b.queryParameters(field, "", false, cycle)...)
	}
	return params
}

// queryParameters expands a single request field into one or more query
// parameters:
//
//   - Scalars, enums, well-known types, and repeated scalars/enums/WKTs
//     produce a single parameter using the field's own schema (repeated
//     scalars become an array-typed parameter, matching the runtime's
//     `?tag=a&tag=b` form).
//   - Non-WKT message fields recurse into their own fields with a dotted
//     prefix ("filter.kind", "filter.range.start", ...).
//   - Map fields emit a single parameter named `<name>[<keyType>]`,
//     matching the runtime's `?labels[foo]=bar` form; see mapQueryParameter.
//   - Repeated messages are skipped — they have no natural query-string
//     representation.
//
// Cycles are bounded by the registry's recursion depth: if recursing into a
// message would exceed the limit on the current path, the recursion is
// truncated at that point and a log line is emitted. Truncating beats failing
// because a partially-flattened spec is still useful.
//
// parentDeprecated propagates deprecation from an ancestor field: if a
// message-typed field is deprecated, all flattened child parameters inherit
// the flag even if the nested fields themselves are not marked deprecated.
func (b *schemaBuilder) queryParameters(field *descriptor.Field, prefix string, parentDeprecated bool, cycle *queryCycleChecker) []*ParameterRef {
	name := prefix + jsonName(field)
	deprecated := parentDeprecated || fieldDeprecated(field)

	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE &&
		wellKnownTypeSchema(field.GetTypeName()) == nil {
		// Non-WKT message: distinguish maps, repeated messages, and regular
		// nested messages. An unresolved message falls through to the
		// single-parameter path so the operation stays syntactically valid.
		msg, err := b.reg.LookupMsg("", field.GetTypeName())
		switch {
		case err != nil:
			grpclog.Infof("protoc-gen-openapiv3: cannot resolve message %q for query parameter %q: %v; emitting opaque parameter",
				field.GetTypeName(), name, err)
		case msg.GetOptions().GetMapEntry():
			// Map fields: emit a single parameter using the runtime's
			// `name[key]` form, where the bracketed token documents the
			// expected key type. The parameter's schema is the map value's
			// schema. Unsupported key types (float, double, bytes) are
			// dropped with a log line — the runtime can't key URLs by them.
			return mapQueryParameter(b, field, name, deprecated, msg)
		case field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED:
			// Repeated messages cannot be naturally represented in a query string.
			return nil
		default:
			if !cycle.enter(msg.FQMN()) {
				grpclog.Infof("protoc-gen-openapiv3: query parameter %q exceeds recursion depth (%d); truncating",
					name, cycle.limit)
				return nil
			}
			defer cycle.leave(msg.FQMN())
			var out []*ParameterRef
			for _, nested := range msg.Fields {
				if !isVisible(fieldVisibility(nested), b.reg) {
					// Field is hidden by visibility rules, skip.
					continue
				}
				out = append(out, b.queryParameters(nested, name+".", deprecated, cycle)...)
			}
			return out
		}
	}

	param := &Parameter{
		Name:   name,
		In:     "query",
		Schema: b.fieldSchema(field),
	}
	if desc := fieldComments(b.reg, field); desc != "" {
		param.Description = desc
	}
	if deprecated {
		param.Deprecated = true
	}
	return []*ParameterRef{{Value: param}}
}

// mapQueryParameter constructs a single query parameter for a proto map
// field. The parameter name is `<jsonName>[<keyType>]`, matching the
// `name[key]` URL form the runtime parses (see
// runtime.PopulateQueryParameters — any real key is accepted between the
// brackets; the token in the spec just documents the expected key type).
// The parameter schema is the map value's schema.
//
// Unsupported key types — float, double, bytes — are dropped with a log
// line: the runtime can't key URL parameters by them.
func mapQueryParameter(b *schemaBuilder, field *descriptor.Field, name string, deprecated bool, entry *descriptor.Message) []*ParameterRef {
	var keyField, valueField *descriptor.Field
	for _, f := range entry.Fields {
		switch f.GetNumber() {
		case 1:
			keyField = f
		case 2:
			valueField = f
		}
	}
	if keyField == nil || valueField == nil {
		return nil
	}
	keyType, ok := queryMapKeyType(keyField.GetType())
	if !ok {
		grpclog.Infof("protoc-gen-openapiv3: map query parameter %q has unsupported key type %s; skipping",
			name, keyField.GetType())
		return nil
	}
	param := &Parameter{
		Name:   fmt.Sprintf("%s[%s]", name, keyType),
		In:     "query",
		Schema: b.fieldSchema(valueField),
	}
	if desc := fieldComments(b.reg, field); desc != "" {
		param.Description = desc
	}
	if deprecated {
		param.Deprecated = true
	}
	return []*ParameterRef{{Value: param}}
}

// queryMapKeyType returns the JSON Schema type name used to document the
// expected key type of a map query parameter. float/double/bytes keys are
// intentionally unsupported: they don't survive round-tripping through a URL
// query string cleanly enough to be useful.
func queryMapKeyType(t descriptorpb.FieldDescriptorProto_Type) (string, bool) {
	switch t {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "string", true
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return "integer", true
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "boolean", true
	}
	return "", false
}

// queryCycleChecker bounds recursion depth when flattening nested message
// fields into query parameters. It tracks how many times each message has
// been entered on the current path; the same message may appear multiple
// times across sibling branches but only up to `limit` times along any one
// chain of recursive calls.
type queryCycleChecker struct {
	depth map[string]int
	limit int
}

func newQueryCycleChecker(limit int) *queryCycleChecker {
	return &queryCycleChecker{depth: make(map[string]int), limit: limit}
}

// enter records a recursion into the named message. It returns false (and
// does not record) if entering would exceed the configured limit, so the
// caller can stop without an unbalanced leave.
func (c *queryCycleChecker) enter(fqmn string) bool {
	if c.depth[fqmn] >= c.limit {
		return false
	}
	c.depth[fqmn]++
	return true
}

func (c *queryCycleChecker) leave(fqmn string) {
	c.depth[fqmn]--
}

// buildRequestBody constructs the requestBody for a binding. There are two
// shapes:
//
//   - body="*": the body is the entire request message minus path parameters.
//     We synthesize an inline object schema rather than referencing the
//     request component, because the component still includes the path
//     fields.
//   - body="some_field": the body is just that field's type.
func buildRequestBody(b *schemaBuilder, m *descriptor.Method, binding *descriptor.Binding) *RequestBodyRef {
	if m.RequestType == nil {
		return nil
	}

	if len(binding.Body.FieldPath) == 0 {
		// body="*": synthesize a body-only inline schema.
		bodySchema := &Schema{Type: SchemaType{"object"}}
		for _, field := range m.RequestType.Fields {
			if !isVisible(fieldVisibility(field), b.reg) {
				// Field is hidden by visibility rules, skip.
				continue
			}
			if isPathParam(field, binding.PathParams) {
				continue
			}
			b.addProperty(bodySchema, field)
		}
		return &RequestBodyRef{Value: &RequestBody{
			Content:  map[string]*MediaType{"application/json": {Schema: &SchemaOrRef{Value: bodySchema}}},
			Required: true,
		}}
	}

	// body="field_name": the body is the type of that single field.
	bodyField := binding.Body.FieldPath[len(binding.Body.FieldPath)-1].Target
	return &RequestBodyRef{Value: &RequestBody{
		Content:  map[string]*MediaType{"application/json": {Schema: b.fieldSchema(bodyField)}},
		Required: true,
	}}
}

// buildResponses constructs the responses map for an RPC: a 200 with the
// response message schema, and (unless disabled) a default google.rpc.Status
// error response.
//
// RPCs returning google.protobuf.Empty get a 200 with an empty object
// schema rather than the HTTP-conventional 204 No Content. The grpc-gateway
// runtime writes `{}` on success regardless of the response type, so the
// spec has to match that or generated clients will reject valid responses.
func buildResponses(b *schemaBuilder, m *descriptor.Method) *Responses {
	resp := NewResponses()

	if m.ResponseType != nil {
		desc := messageComments(b.reg, m.ResponseType)
		if desc == "" {
			desc = "A successful response."
		}
		var schema *SchemaOrRef
		if wkt := wellKnownTypeSchema(m.ResponseType.FQMN()); wkt != nil {
			schema = &SchemaOrRef{Value: wkt}
		} else {
			b.ensureMessageSchema(m.ResponseType)
			schema = NewSchemaRef(schemaName(m.ResponseType.FQMN()))
		}
		resp.Codes["200"] = &ResponseRef{Value: NewResponse(desc).WithJSONSchema(schema)}
	}

	if !b.reg.GetDisableDefaultErrors() {
		ensureStatusSchema(b.doc)
		resp.Default = &ResponseRef{Value: NewResponse("An unexpected error response.").WithJSONSchema(NewSchemaRef(statusSchemaName))}
	}
	return resp
}

// ensureStatusSchema makes sure the google.rpc.Status component schema is
// present.
func ensureStatusSchema(doc *Document) {
	if _, ok := doc.Components.Schemas[statusSchemaName]; ok {
		return
	}
	t := true
	doc.Components.Schemas[statusSchemaName] = &SchemaOrRef{Value: &Schema{
		Type:        SchemaType{"object"},
		Description: "The Status type defines a logical error model suitable for different programming environments.",
		Properties: map[string]*SchemaOrRef{
			"code": {Value: &Schema{
				Type:        SchemaType{"integer"},
				Format:      "int32",
				Description: "The status code, which should be an enum value of google.rpc.Code.",
			}},
			"message": {Value: &Schema{
				Type:        SchemaType{"string"},
				Description: "A developer-facing error message.",
			}},
			"details": {Value: &Schema{
				Type:        SchemaType{"array"},
				Description: "A list of messages that carry the error details.",
				Items: &SchemaOrRef{Value: &Schema{
					Type: SchemaType{"object"},
					Properties: map[string]*SchemaOrRef{
						"@type": {Value: &Schema{Type: SchemaType{"string"}}},
					},
					AdditionalProperties: &AdditionalProperties{Bool: &t},
				}},
			}},
		},
	}}
}

// operationID returns the OpenAPI operationId for a method binding. We use
// "<Service>_<Method>" with a numeric suffix on bindings beyond the first.
func operationID(svc *descriptor.Service, m *descriptor.Method, bindingIdx int) string {
	id := fmt.Sprintf("%s_%s", svc.GetName(), m.GetName())
	if bindingIdx > 0 {
		id = fmt.Sprintf("%s_%d", id, bindingIdx)
	}
	return id
}

func needsRequestBody(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH":
		return true
	}
	return false
}

func isPathParam(field *descriptor.Field, params []descriptor.Parameter) bool {
	for _, p := range params {
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
	if len(body.FieldPath) == 0 {
		return true // body="*"
	}
	for _, fp := range body.FieldPath {
		if fp.Target == field {
			return true
		}
	}
	return false
}
