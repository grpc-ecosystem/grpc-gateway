package genopenapiv3

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/grpc/grpclog"
)

// parseExampleValue attempts to parse a string as JSON.
// If parsing succeeds, returns the parsed value (which could be a map, slice, number, bool, or null).
// If parsing fails, returns the original string (for simple string examples like "hello").
// This ensures examples are properly typed in the generated OpenAPI spec:
//   - `"{\"id\": 123}"` becomes `{"id": 123}` (object)
//   - `"[1, 2, 3]"` becomes `[1, 2, 3]` (array)
//   - `"42"` becomes `42` (number)
//   - `"true"` becomes `true` (boolean)
//   - `"hello"` stays `"hello"` (string, since it's not valid JSON)
func parseExampleValue(s string) any {
	if s == "" {
		return nil
	}

	var result any
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		// Not valid JSON, return as string
		return s
	}
	return result
}

// makeExamplesArray converts a single example value to the JSON Schema examples array format.
// In OpenAPI 3.1.0, Schema Objects use JSON Schema which has `examples` as an array.
func makeExamplesArray(value any) []any {
	if value == nil {
		return nil
	}
	return []any{value}
}

// makeExamplesMap converts a single example value to the OpenAPI examples map format.
// Used for Parameter, Header, and MediaType objects (not Schema objects).
func makeExamplesMap(value any) map[string]*ExampleRef {
	if value == nil {
		return nil
	}
	return map[string]*ExampleRef{
		"example": {
			Value: &Example{
				Value: value,
			},
		},
	}
}

// convertExamplesMap converts a proto examples map to OpenAPI examples map.
// This handles the plural `examples` field from proto options.
// Supports both inline examples and references (ExampleOrReference).
func convertExamplesMap(protoExamples map[string]*options.ExampleOrReference) map[string]*ExampleRef {
	if len(protoExamples) == 0 {
		return nil
	}
	result := make(map[string]*ExampleRef, len(protoExamples))
	for name, exOrRef := range protoExamples {
		result[name] = convertExampleOrReference(exOrRef)
	}
	return result
}

// convertExampleOrReference converts a proto ExampleOrReference to an ExampleRef.
func convertExampleOrReference(exOrRef *options.ExampleOrReference) *ExampleRef {
	if exOrRef == nil {
		return nil
	}

	switch v := exOrRef.GetOneof().(type) {
	case *options.ExampleOrReference_Reference:
		return &ExampleRef{
			Ref: v.Reference.GetRef(),
		}
	case *options.ExampleOrReference_Example:
		ex := v.Example
		return &ExampleRef{
			Value: &Example{
				Summary:       ex.GetSummary(),
				Description:   ex.GetDescription(),
				Value:         parseExampleValue(ex.GetValue()),
				ExternalValue: ex.GetExternalValue(),
			},
		}
	default:
		return nil
	}
}

// applyFileAnnotation applies file-level OpenAPI v3 annotations to the document.
func (g *generator) applyFileAnnotation(doc *OpenAPI, file *descriptor.File) {
	opts := getFileAnnotation(file)
	if opts == nil {
		return
	}

	// Apply OpenAPI version if specified
	if opts.GetOpenapi() != "" {
		doc.OpenAPI = opts.GetOpenapi()
	}

	// Apply Info
	if info := opts.GetInfo(); info != nil {
		g.applyInfoAnnotation(doc.Info, info)
	}

	// Apply Servers
	for _, s := range opts.GetServers() {
		doc.Servers = append(doc.Servers, convertServer(s))
	}

	// Apply Security
	for _, sec := range opts.GetSecurity() {
		doc.Security = append(doc.Security, convertSecurityRequirement(sec))
	}

	// Apply Tags (prepend to existing tags)
	for _, t := range opts.GetTags() {
		doc.Tags = append([]*Tag{convertTag(t)}, doc.Tags...)
	}

	// Apply External Docs
	if extDocs := opts.GetExternalDocs(); extDocs != nil {
		doc.ExternalDocs = convertExternalDocs(extDocs)
	}

	// Apply Components (security schemes, responses, etc.)
	if comp := opts.GetComponents(); comp != nil {
		g.applyComponentsAnnotation(doc.Components, comp)
	}
}

// applyInfoAnnotation applies info-level annotations.
func (g *generator) applyInfoAnnotation(info *Info, opts *options.Info) {
	if opts.GetTitle() != "" {
		info.Title = opts.GetTitle()
	}
	if opts.GetSummary() != "" {
		info.Summary = opts.GetSummary()
	}
	if opts.GetDescription() != "" {
		info.Description = opts.GetDescription()
	}
	if opts.GetTermsOfService() != "" {
		info.TermsOfService = opts.GetTermsOfService()
	}
	if opts.GetVersion() != "" {
		info.Version = opts.GetVersion()
	}
	if c := opts.GetContact(); c != nil {
		info.Contact = &Contact{
			Name:  c.GetName(),
			URL:   c.GetUrl(),
			Email: c.GetEmail(),
		}
	}
	if l := opts.GetLicense(); l != nil {
		license := &License{
			Name: l.GetName(),
		}
		// In OpenAPI 3.1.0, identifier and url are mutually exclusive.
		// Prefer identifier if both are set.
		if l.GetIdentifier() != "" {
			license.Identifier = l.GetIdentifier()
		} else if l.GetUrl() != "" {
			license.URL = l.GetUrl()
		}
		info.License = license
	}
}

// applyComponentsAnnotation applies components-level annotations.
func (g *generator) applyComponentsAnnotation(comp *Components, opts *options.Components) {
	// Apply Security Schemes
	for name, scheme := range opts.GetSecuritySchemes() {
		if comp.SecuritySchemes == nil {
			comp.SecuritySchemes = make(map[string]*SecuritySchemeRef)
		}
		comp.SecuritySchemes[name] = &SecuritySchemeRef{Value: convertSecurityScheme(scheme)}
	}

	// Apply Responses
	for _, namedResp := range opts.GetResponses() {
		if comp.Responses == nil {
			comp.Responses = make(map[string]*ResponseRef)
		}
		comp.Responses[namedResp.GetName()] = g.convertResponseOrReference(namedResp.GetValue())
	}

	// Apply Parameters
	for name, param := range opts.GetParameters() {
		if comp.Parameters == nil {
			comp.Parameters = make(map[string]*ParameterRef)
		}
		comp.Parameters[name] = &ParameterRef{Value: g.convertParameter(param)}
	}

	// Apply Request Bodies
	for name, body := range opts.GetRequestBodies() {
		if comp.RequestBodies == nil {
			comp.RequestBodies = make(map[string]*RequestBodyRef)
		}
		comp.RequestBodies[name] = &RequestBodyRef{Value: g.convertRequestBody(body)}
	}

	// Apply Headers
	for name, headerOrRef := range opts.GetHeaders() {
		if comp.Headers == nil {
			comp.Headers = make(map[string]*HeaderRef)
		}
		comp.Headers[name] = g.convertHeaderOrReference(headerOrRef)
	}
}

// applyOperationAnnotation applies method-level annotations to an operation.
func (g *generator) applyOperationAnnotation(op *Operation, method *descriptor.Method) {
	opts := getMethodAnnotation(method)
	if opts == nil {
		return
	}

	// Override summary if provided
	if opts.GetSummary() != "" {
		op.Summary = opts.GetSummary()
	}

	// Override description if provided
	if opts.GetDescription() != "" {
		op.Description = opts.GetDescription()
	}

	// Override operation ID if provided
	if opts.GetOperationId() != "" {
		op.OperationID = opts.GetOperationId()
	}

	// Override or append tags
	if len(opts.GetTags()) > 0 {
		op.Tags = opts.GetTags()
	}

	// Apply deprecated flag
	if opts.GetDeprecated() {
		op.Deprecated = true
	}

	// Apply external docs
	if extDocs := opts.GetExternalDocs(); extDocs != nil {
		op.ExternalDocs = convertExternalDocs(extDocs)
	}

	// Apply security requirements
	for _, sec := range opts.GetSecurity() {
		op.Security = append(op.Security, convertSecurityRequirement(sec))
	}

	// Apply servers
	for _, s := range opts.GetServers() {
		op.Servers = append(op.Servers, convertServer(s))
	}

	// Apply additional responses (merge with existing)
	if responses := opts.GetResponses(); responses != nil {
		if op.Responses == nil {
			op.Responses = NewResponses()
		}
		// Apply default response - merges inline, overwrites for reference
		if defaultResp := responses.GetDefault(); defaultResp != nil {
			op.Responses.Default = g.applyResponseOrReference(op.Responses.Default, defaultResp)
		}
		// Apply status code specific responses - merges inline, overwrites for reference
		for _, namedResp := range responses.GetResponseOrReference() {
			code := namedResp.GetName()
			if code == "200" {
				methodName := method.GetName()
				if method.Service != nil {
					methodName = method.Service.GetName() + "." + methodName
				}
				grpclog.Warningf("Annotation overrides 200 response for method %s - the success response schema should be derived from proto return type", methodName)
			}
			existing := op.Responses.Codes[code]
			op.Responses.Codes[code] = g.applyResponseOrReference(existing, namedResp.GetValue())
		}
	}

	// Apply request body annotation - merges with existing (preserves content if not specified)
	if reqBody := opts.GetRequestBody(); reqBody != nil {
		op.RequestBody = g.applyRequestBody(op.RequestBody, reqBody)
	}

	// Apply custom parameters (headers and cookies)
	if params := opts.GetParameters(); params != nil {
		// Add header parameters
		for _, headerParam := range params.GetHeaders() {
			if paramRef := g.convertHeaderParameterOrReference(headerParam); paramRef != nil {
				op.Parameters = append(op.Parameters, paramRef)
			}
		}
		// Add cookie parameters
		for _, cookieParam := range params.GetCookies() {
			if paramRef := g.convertCookieParameterOrReference(cookieParam); paramRef != nil {
				op.Parameters = append(op.Parameters, paramRef)
			}
		}
	}
}

// applySchemaAnnotation applies message-level annotations to a schema.
func (g *generator) applySchemaAnnotation(schema *Schema, msg *descriptor.Message) {
	opts := getMessageAnnotation(msg)
	if opts == nil {
		return
	}

	// Apply title
	if opts.GetTitle() != "" {
		schema.Title = opts.GetTitle()
	}

	// Apply description
	if opts.GetDescription() != "" {
		schema.Description = opts.GetDescription()
	}

	// Apply required fields
	if len(opts.GetRequired()) > 0 {
		schema.Required = opts.GetRequired()
	}

	// Apply example (using examples map for OpenAPI 3.1.0 compliance)
	if opts.GetExample() != "" {
		schema.Examples = makeExamplesArray(parseExampleValue(opts.GetExample()))
	}

	// Apply read only
	if opts.GetReadOnly() {
		schema.ReadOnly = true
	}

	// Apply write only
	if opts.GetWriteOnly() {
		schema.WriteOnly = true
	}

	// Apply nullable using version-appropriate method
	if opts.GetNullable() {
		g.applyNullable(schema)
	}

	// Apply deprecated
	if opts.GetDeprecated() {
		schema.Deprecated = true
	}

	// Apply external docs
	if extDocs := opts.GetExternalDocs(); extDocs != nil {
		schema.ExternalDocs = convertExternalDocs(extDocs)
	}

	// Apply composition types (OpenAPI v3 specific)

	// Apply allOf
	if len(opts.GetAllOf()) > 0 {
		for _, allOfSchema := range opts.GetAllOf() {
			schema.AllOf = append(schema.AllOf, g.convertSchemaOrReference(allOfSchema))
		}
	}

	// Apply anyOf
	if len(opts.GetAnyOf()) > 0 {
		for _, anyOfSchema := range opts.GetAnyOf() {
			schema.AnyOf = append(schema.AnyOf, g.convertSchemaOrReference(anyOfSchema))
		}
	}

	// Apply oneOf (appends to auto-detected oneOf from proto oneof fields)
	if len(opts.GetOneOf()) > 0 {
		for _, oneOfSchema := range opts.GetOneOf() {
			schema.OneOf = append(schema.OneOf, g.convertSchemaOrReference(oneOfSchema))
		}
	}

	// Apply not
	if notSchema := opts.GetNot(); notSchema != nil {
		schema.Not = g.convertSchemaOrReference(notSchema)
	}

	// Apply discriminator
	if disc := opts.GetDiscriminator(); disc != nil {
		schema.Discriminator = &Discriminator{
			PropertyName: disc.GetPropertyName(),
			Mapping:      disc.GetMapping(),
		}
	}
}

// applyFieldAnnotation applies field-level annotations to a field schema.
func (g *generator) applyFieldAnnotation(schema *Schema, field *descriptor.Field) {
	opts := getFieldAnnotation(field)
	if opts == nil {
		return
	}

	// Apply title
	if opts.GetTitle() != "" {
		schema.Title = opts.GetTitle()
	}

	// Apply description
	if opts.GetDescription() != "" {
		schema.Description = opts.GetDescription()
	}

	// Apply default value
	if opts.GetDefault() != "" {
		schema.Default = opts.GetDefault()
	}

	// Apply example (using examples map for OpenAPI 3.1.0 compliance)
	if opts.GetExample() != "" {
		schema.Examples = makeExamplesArray(parseExampleValue(opts.GetExample()))
	}

	// Apply format
	if opts.GetFormat() != "" {
		schema.Format = opts.GetFormat()
	}

	// Apply pattern
	if opts.GetPattern() != "" {
		schema.Pattern = opts.GetPattern()
	}

	// Apply min/max length
	if opts.GetMinLength() > 0 {
		minLen := opts.GetMinLength()
		schema.MinLength = &minLen
	}
	if opts.GetMaxLength() > 0 {
		maxLen := opts.GetMaxLength()
		schema.MaxLength = &maxLen
	}

	// Apply min/max items
	if opts.GetMinItems() > 0 {
		minItems := opts.GetMinItems()
		schema.MinItems = &minItems
	}
	if opts.GetMaxItems() > 0 {
		maxItems := opts.GetMaxItems()
		schema.MaxItems = &maxItems
	}

	// Apply unique items
	if opts.GetUniqueItems() {
		schema.UniqueItems = true
	}

	// Apply min/max properties
	if opts.GetMinProperties() > 0 {
		minProps := opts.GetMinProperties()
		schema.MinProperties = &minProps
	}
	if opts.GetMaxProperties() > 0 {
		maxProps := opts.GetMaxProperties()
		schema.MaxProperties = &maxProps
	}

	// Apply numeric constraints
	// Using pointer checks to correctly handle zero values (0 is valid for min/max)
	if opts.MultipleOf != nil {
		multipleOf := opts.GetMultipleOf()
		schema.MultipleOf = &multipleOf
	}
	if opts.Minimum != nil {
		min := opts.GetMinimum()
		schema.Minimum = &min
	}
	if opts.Maximum != nil {
		max := opts.GetMaximum()
		schema.Maximum = &max
	}
	if opts.ExclusiveMinimum != nil {
		exclusiveMin := opts.GetExclusiveMinimum()
		schema.ExclusiveMinimum = &exclusiveMin
	}
	if opts.ExclusiveMaximum != nil {
		exclusiveMax := opts.GetExclusiveMaximum()
		schema.ExclusiveMaximum = &exclusiveMax
	}

	// Apply read only
	if opts.GetReadOnly() {
		schema.ReadOnly = true
	}

	// Apply write only
	if opts.GetWriteOnly() {
		schema.WriteOnly = true
	}

	// Apply nullable using version-appropriate method
	if opts.GetNullable() {
		g.applyNullable(schema)
	}

	// Apply deprecated
	if opts.GetDeprecated() {
		schema.Deprecated = true
	}

	// Apply external docs
	if extDocs := opts.GetExternalDocs(); extDocs != nil {
		schema.ExternalDocs = convertExternalDocs(extDocs)
	}

	// Apply composition types (OpenAPI v3 specific)

	// Apply allOf
	if len(opts.GetAllOf()) > 0 {
		for _, allOfSchema := range opts.GetAllOf() {
			schema.AllOf = append(schema.AllOf, g.convertSchemaOrReference(allOfSchema))
		}
	}

	// Apply anyOf
	if len(opts.GetAnyOf()) > 0 {
		for _, anyOfSchema := range opts.GetAnyOf() {
			schema.AnyOf = append(schema.AnyOf, g.convertSchemaOrReference(anyOfSchema))
		}
	}

	// Apply oneOf
	if len(opts.GetOneOf()) > 0 {
		for _, oneOfSchema := range opts.GetOneOf() {
			schema.OneOf = append(schema.OneOf, g.convertSchemaOrReference(oneOfSchema))
		}
	}

	// Apply not
	if notSchema := opts.GetNot(); notSchema != nil {
		schema.Not = g.convertSchemaOrReference(notSchema)
	}

	// Apply discriminator
	if disc := opts.GetDiscriminator(); disc != nil {
		schema.Discriminator = &Discriminator{
			PropertyName: disc.GetPropertyName(),
			Mapping:      disc.GetMapping(),
		}
	}
}

// applyServiceAnnotation applies service-level annotations to a tag.
func (g *generator) applyServiceAnnotation(tag *Tag, svc *descriptor.Service) {
	opts := getServiceAnnotation(svc)
	if opts == nil {
		return
	}

	// Override name if provided
	if opts.GetName() != "" {
		tag.Name = opts.GetName()
	}

	// Override description if provided
	if opts.GetDescription() != "" {
		tag.Description = opts.GetDescription()
	}

	// Apply external docs
	if extDocs := opts.GetExternalDocs(); extDocs != nil {
		tag.ExternalDocs = convertExternalDocs(extDocs)
	}
}

// applyEnumAnnotation applies enum-level annotations to an enum schema.
func (g *generator) applyEnumAnnotation(schema *Schema, enum *descriptor.Enum) {
	opts := getEnumAnnotation(enum)
	if opts == nil {
		return
	}

	// Apply title
	if opts.GetTitle() != "" {
		schema.Title = opts.GetTitle()
	}

	// Apply description
	if opts.GetDescription() != "" {
		schema.Description = opts.GetDescription()
	}

	// Apply default value
	if opts.GetDefault() != "" {
		schema.Default = opts.GetDefault()
	}

	// Apply example (using examples map for OpenAPI 3.1.0 compliance)
	if opts.GetExample() != "" {
		schema.Examples = makeExamplesArray(parseExampleValue(opts.GetExample()))
	}

	// Apply deprecated
	if opts.GetDeprecated() {
		schema.Deprecated = true
	}

	// Apply external docs
	if extDocs := opts.GetExternalDocs(); extDocs != nil {
		schema.ExternalDocs = convertExternalDocs(extDocs)
	}
}

// Conversion helpers

func convertServer(s *options.Server) *Server {
	server := &Server{
		URL:         s.GetUrl(),
		Description: s.GetDescription(),
	}
	if len(s.GetVariables()) > 0 {
		server.Variables = make(map[string]*ServerVariable)
		for name, v := range s.GetVariables() {
			server.Variables[name] = &ServerVariable{
				Enum:        v.GetEnum(),
				Default:     v.GetDefault(),
				Description: v.GetDescription(),
			}
		}
	}
	return server
}

func convertSecurityRequirement(sec *options.SecurityRequirement) SecurityRequirement {
	result := make(SecurityRequirement)
	for name, val := range sec.GetSecurityRequirement() {
		scopes := val.GetScope()
		if scopes == nil {
			scopes = []string{} // OpenAPI requires empty array, not null
		}
		result[name] = scopes
	}
	return result
}

func convertTag(t *options.Tag) *Tag {
	tag := &Tag{
		Name:        t.GetName(),
		Description: t.GetDescription(),
	}
	if extDocs := t.GetExternalDocs(); extDocs != nil {
		tag.ExternalDocs = convertExternalDocs(extDocs)
	}
	return tag
}

func convertExternalDocs(extDocs *options.ExternalDocumentation) *ExternalDocumentation {
	return &ExternalDocumentation{
		Description: extDocs.GetDescription(),
		URL:         extDocs.GetUrl(),
	}
}

func convertSecurityScheme(scheme *options.SecurityScheme) *SecurityScheme {
	ss := &SecurityScheme{
		Description:      scheme.GetDescription(),
		Name:             scheme.GetName(),
		Scheme:           scheme.GetScheme(),
		BearerFormat:     scheme.GetBearerFormat(),
		OpenIdConnectUrl: scheme.GetOpenIdConnectUrl(),
	}

	// Convert type
	switch scheme.GetType() {
	case options.SecurityScheme_TYPE_API_KEY:
		ss.Type = "apiKey"
	case options.SecurityScheme_TYPE_HTTP:
		ss.Type = "http"
	case options.SecurityScheme_TYPE_OAUTH2:
		ss.Type = "oauth2"
	case options.SecurityScheme_TYPE_OPEN_ID_CONNECT:
		ss.Type = "openIdConnect"
	}

	// Convert in
	switch scheme.GetIn() {
	case options.SecurityScheme_IN_QUERY:
		ss.In = "query"
	case options.SecurityScheme_IN_HEADER:
		ss.In = "header"
	case options.SecurityScheme_IN_COOKIE:
		ss.In = "cookie"
	}

	// Convert flows
	if flows := scheme.GetFlows(); flows != nil {
		ss.Flows = &OAuthFlows{}
		if f := flows.GetImplicit(); f != nil {
			ss.Flows.Implicit = convertOAuthFlow(f)
		}
		if f := flows.GetPassword(); f != nil {
			ss.Flows.Password = convertOAuthFlow(f)
		}
		if f := flows.GetClientCredentials(); f != nil {
			ss.Flows.ClientCredentials = convertOAuthFlow(f)
		}
		if f := flows.GetAuthorizationCode(); f != nil {
			ss.Flows.AuthorizationCode = convertOAuthFlow(f)
		}
	}

	return ss
}

func convertOAuthFlow(flow *options.OAuthFlow) *OAuthFlow {
	return &OAuthFlow{
		AuthorizationURL: flow.GetAuthorizationUrl(),
		TokenURL:         flow.GetTokenUrl(),
		RefreshURL:       flow.GetRefreshUrl(),
		Scopes:           flow.GetScopes(),
	}
}

func (g *generator) convertResponse(resp *options.Response) *Response {
	r := &Response{
		Description: resp.GetDescription(),
	}
	if len(resp.GetHeaders()) > 0 {
		r.Headers = make(map[string]*HeaderRef)
		for name, h := range resp.GetHeaders() {
			r.Headers[name] = g.convertHeaderOrReference(h)
		}
	}
	if len(resp.GetContent()) > 0 {
		r.Content = make(map[string]*MediaType)
		for mediaType, mt := range resp.GetContent() {
			r.Content[mediaType] = g.convertMediaType(mt)
		}
	}
	return r
}

func (g *generator) convertParameter(param *options.Parameter) *Parameter {
	p := &Parameter{
		Name:            param.GetName(),
		In:              param.GetIn(),
		Description:     param.GetDescription(),
		Required:        param.GetRequired(),
		Deprecated:      param.GetDeprecated(),
		AllowEmptyValue: param.GetAllowEmptyValue(),
		Style:           param.GetStyle(),
		AllowReserved:   param.GetAllowReserved(),
	}
	if param.GetExplode() {
		explode := true
		p.Explode = &explode
	}
	// Prefer plural examples map over singular example field (OpenAPI 3.1.0 compliance)
	if len(param.GetExamples()) > 0 {
		p.Examples = convertExamplesMap(param.GetExamples())
	} else if param.GetExample() != "" {
		p.Examples = makeExamplesMap(parseExampleValue(param.GetExample()))
	}
	if schema := param.GetSchema(); schema != nil {
		p.Schema = g.convertSchemaOrReference(schema)
	}
	return p
}

func (g *generator) convertRequestBody(body *options.RequestBody) *RequestBody {
	rb := &RequestBody{
		Description: body.GetDescription(),
		Required:    body.GetRequired(),
	}
	if len(body.GetContent()) > 0 {
		rb.Content = make(map[string]*MediaType)
		for mediaType, mt := range body.GetContent() {
			rb.Content[mediaType] = g.convertMediaType(mt)
		}
	}
	return rb
}

func (g *generator) convertHeader(header *options.Header) *Header {
	h := &Header{
		Description:     header.GetDescription(),
		Required:        header.GetRequired(),
		Deprecated:      header.GetDeprecated(),
		AllowEmptyValue: header.GetAllowEmptyValue(),
		Style:           header.GetStyle(),
		AllowReserved:   header.GetAllowReserved(),
	}
	if header.GetExplode() {
		explode := true
		h.Explode = &explode
	}
	// Prefer plural examples map over singular example field (OpenAPI 3.1.0 compliance)
	if len(header.GetExamples()) > 0 {
		h.Examples = convertExamplesMap(header.GetExamples())
	} else if header.GetExample() != "" {
		h.Examples = makeExamplesMap(parseExampleValue(header.GetExample()))
	}
	if schema := header.GetSchema(); schema != nil {
		h.Schema = g.convertSchemaOrReference(schema)
	}
	return h
}

// convertHeaderOrReference converts a proto HeaderOrReference to a HeaderRef.
func (g *generator) convertHeaderOrReference(hor *options.HeaderOrReference) *HeaderRef {
	if hor == nil {
		return nil
	}

	switch v := hor.GetOneof().(type) {
	case *options.HeaderOrReference_Reference:
		return &HeaderRef{
			Ref: v.Reference.GetRef(),
		}
	case *options.HeaderOrReference_Header:
		return &HeaderRef{
			Value: g.convertHeader(v.Header),
		}
	default:
		return nil
	}
}

// convertResponseOrReference converts a proto ResponseOrReference to a ResponseRef.
func (g *generator) convertResponseOrReference(ror *options.ResponseOrReference) *ResponseRef {
	if ror == nil {
		return nil
	}

	switch v := ror.GetOneof().(type) {
	case *options.ResponseOrReference_Reference:
		return &ResponseRef{
			Ref: v.Reference.GetRef(),
		}
	case *options.ResponseOrReference_Response:
		return &ResponseRef{
			Value: g.convertResponse(v.Response),
		}
	default:
		return nil
	}
}

// applyResponseOrReference applies a ResponseOrReference annotation to an existing ResponseRef.
// For reference annotations ($ref): overwrites entirely.
// For inline annotations: merges with existing (preserves content/headers if not specified in annotation).
func (g *generator) applyResponseOrReference(existing *ResponseRef, annotation *options.ResponseOrReference) *ResponseRef {
	if annotation == nil {
		return existing
	}

	switch v := annotation.GetOneof().(type) {
	case *options.ResponseOrReference_Reference:
		// Reference: overwrite entirely
		return &ResponseRef{
			Ref: v.Reference.GetRef(),
		}
	case *options.ResponseOrReference_Response:
		// Inline: merge with existing
		annotationResp := v.Response
		result := &Response{}

		// Description: annotation wins if specified
		if annotationResp.GetDescription() != "" {
			result.Description = annotationResp.GetDescription()
		} else if existing != nil && existing.Value != nil {
			result.Description = existing.Value.Description
		}

		// Content: keep existing if annotation doesn't specify
		if len(annotationResp.GetContent()) > 0 {
			result.Content = make(map[string]*MediaType)
			for mediaType, mt := range annotationResp.GetContent() {
				result.Content[mediaType] = g.convertMediaType(mt)
			}
		} else if existing != nil && existing.Value != nil && len(existing.Value.Content) > 0 {
			result.Content = existing.Value.Content
		}

		// Headers: merge maps (existing first, then annotation overrides/adds)
		if (existing != nil && existing.Value != nil && len(existing.Value.Headers) > 0) || len(annotationResp.GetHeaders()) > 0 {
			result.Headers = make(map[string]*HeaderRef)
			// Copy existing headers first
			if existing != nil && existing.Value != nil {
				for name, header := range existing.Value.Headers {
					result.Headers[name] = header
				}
			}
			// Add/override with annotation headers
			for name, headerOrRef := range annotationResp.GetHeaders() {
				result.Headers[name] = g.convertHeaderOrReference(headerOrRef)
			}
		}

		return &ResponseRef{Value: result}
	default:
		return existing
	}
}

// applyRequestBody applies a RequestBody annotation to an existing RequestBodyRef.
// Merges with existing (preserves content if not specified in annotation).
// Note: References are not supported for Operation.request_body because the
// request body schema is derived from the proto input message.
func (g *generator) applyRequestBody(existing *RequestBodyRef, annotation *options.RequestBody) *RequestBodyRef {
	if annotation == nil {
		return existing
	}

	result := &RequestBody{}

	// Description: annotation wins if specified
	if annotation.GetDescription() != "" {
		result.Description = annotation.GetDescription()
	} else if existing != nil && existing.Value != nil {
		result.Description = existing.Value.Description
	}

	// Content: keep existing if annotation doesn't specify
	if len(annotation.GetContent()) > 0 {
		result.Content = make(map[string]*MediaType)
		for mediaType, mt := range annotation.GetContent() {
			result.Content[mediaType] = g.convertMediaType(mt)
		}
	} else if existing != nil && existing.Value != nil && len(existing.Value.Content) > 0 {
		result.Content = existing.Value.Content
	}

	// Required: annotation wins if true, otherwise preserve existing
	if annotation.GetRequired() {
		result.Required = true
	} else if existing != nil && existing.Value != nil {
		result.Required = existing.Value.Required
	}

	return &RequestBodyRef{Value: result}
}

// convertCookieParameter converts a proto CookieParameter to a Parameter with in="cookie".
// CookieParameter includes the name field directly, unlike Cookie.
func (g *generator) convertCookieParameter(cookie *options.CookieParameter) *Parameter {
	if cookie == nil {
		return nil
	}
	p := &Parameter{
		Name:        cookie.GetName(),
		In:          "cookie",
		Description: cookie.GetDescription(),
		Required:    cookie.GetRequired(),
		Deprecated:  cookie.GetDeprecated(),
	}
	// Prefer plural examples map over singular example field (OpenAPI 3.1.0 compliance)
	if len(cookie.GetExamples()) > 0 {
		p.Examples = convertExamplesMap(cookie.GetExamples())
	} else if cookie.GetExample() != "" {
		p.Examples = makeExamplesMap(parseExampleValue(cookie.GetExample()))
	}
	// Convert schema or default to string type
	if schema := cookie.GetSchema(); schema != nil {
		p.Schema = g.convertSchemaOrReference(schema)
	} else {
		// Default to string type if no schema specified
		p.Schema = &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}}
	}
	return p
}

// convertHeaderParameter converts a proto HeaderParameter to a Parameter with in="header".
// HeaderParameter includes the name field directly, unlike Header.
func (g *generator) convertHeaderParameter(header *options.HeaderParameter) *Parameter {
	if header == nil {
		return nil
	}
	p := &Parameter{
		Name:            header.GetName(),
		In:              "header",
		Description:     header.GetDescription(),
		Required:        header.GetRequired(),
		Deprecated:      header.GetDeprecated(),
		AllowEmptyValue: header.GetAllowEmptyValue(),
		Style:           header.GetStyle(),
		AllowReserved:   header.GetAllowReserved(),
	}
	if header.GetExplode() {
		explode := true
		p.Explode = &explode
	}
	// Prefer plural examples map over singular example field (OpenAPI 3.1.0 compliance)
	if len(header.GetExamples()) > 0 {
		p.Examples = convertExamplesMap(header.GetExamples())
	} else if header.GetExample() != "" {
		p.Examples = makeExamplesMap(parseExampleValue(header.GetExample()))
	}
	// Convert schema or default to string type
	if schema := header.GetSchema(); schema != nil {
		p.Schema = g.convertSchemaOrReference(schema)
	} else {
		// Default to string type if no schema specified
		p.Schema = &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}}
	}
	return p
}

// convertHeaderParameterOrReference converts a HeaderParameterOrReference to a ParameterRef.
// Headers are converted to parameters with in="header".
func (g *generator) convertHeaderParameterOrReference(param *options.HeaderParameterOrReference) *ParameterRef {
	if param == nil {
		return nil
	}
	switch v := param.GetOneof().(type) {
	case *options.HeaderParameterOrReference_Reference:
		return &ParameterRef{
			Ref: v.Reference.GetRef(),
		}
	case *options.HeaderParameterOrReference_Header:
		return &ParameterRef{
			Value: g.convertHeaderParameter(v.Header),
		}
	default:
		return nil
	}
}

// convertCookieParameterOrReference converts a CookieParameterOrReference to a ParameterRef.
// Cookies are converted to parameters with in="cookie".
func (g *generator) convertCookieParameterOrReference(param *options.CookieParameterOrReference) *ParameterRef {
	if param == nil {
		return nil
	}
	switch v := param.GetOneof().(type) {
	case *options.CookieParameterOrReference_Reference:
		return &ParameterRef{
			Ref: v.Reference.GetRef(),
		}
	case *options.CookieParameterOrReference_Cookie:
		return &ParameterRef{
			Value: g.convertCookieParameter(v.Cookie),
		}
	default:
		return nil
	}
}

func (g *generator) convertMediaType(mt *options.MediaType) *MediaType {
	result := &MediaType{}
	// Prefer plural examples map over singular example field (OpenAPI 3.1.0 compliance)
	if len(mt.GetExamples()) > 0 {
		result.Examples = convertExamplesMap(mt.GetExamples())
	} else if mt.GetExample() != "" {
		result.Examples = makeExamplesMap(parseExampleValue(mt.GetExample()))
	}
	if schema := mt.GetSchema(); schema != nil {
		result.Schema = g.convertSchemaOrReference(schema)
	}
	return result
}

// convertSchemaOrReference converts a proto SchemaOrReference to a Go SchemaRef.
// This handles the discriminated union of inline schema vs reference.
func (g *generator) convertSchemaOrReference(sor *options.SchemaOrReference) *SchemaOrReference {
	if sor == nil {
		return nil
	}

	switch v := sor.GetOneof().(type) {
	case *options.SchemaOrReference_Reference:
		return &SchemaOrReference{
			Reference: &Reference{
				Ref:         v.Reference.GetRef(),
				Summary:     v.Reference.GetSummary(),
				Description: v.Reference.GetDescription(),
			},
		}
	case *options.SchemaOrReference_Value:
		return g.convertSchema(v.Value)
	default:
		return nil
	}
}

func (g *generator) convertSchema(schema *options.Schema) *SchemaOrReference {
	s := &Schema{
		Type:        SchemaType(schema.GetType()),
		Format:      schema.GetFormat(),
		Title:       schema.GetTitle(),
		Description: schema.GetDescription(),
		ReadOnly:    schema.GetReadOnly(),
		WriteOnly:   schema.GetWriteOnly(),
		Deprecated:  schema.GetDeprecated(),
		Pattern:     schema.GetPattern(),
		Required:    schema.GetRequired(),
		UniqueItems: schema.GetUniqueItems(),
	}

	// Apply nullable - version-specific output is handled by the adapter
	if schema.GetNullable() {
		s.Nullable = true
	}

	// Apply default
	if schema.GetDefault() != "" {
		s.Default = schema.GetDefault()
	}

	// Apply example (using examples map for OpenAPI 3.1.0 compliance)
	if schema.GetExample() != "" {
		s.Examples = makeExamplesArray(parseExampleValue(schema.GetExample()))
	}

	// Apply enum values
	if len(schema.GetEnum()) > 0 {
		for _, e := range schema.GetEnum() {
			s.Enum = append(s.Enum, e)
		}
	}

	// Apply numeric constraints
	if schema.GetMultipleOf() != 0 {
		multipleOf := schema.GetMultipleOf()
		s.MultipleOf = &multipleOf
	}
	// Use pointer checks to verify presence instead of value checks,
	// so that explicitly setting minimum: 0 or maximum: 0 is preserved.
	if schema.Minimum != nil {
		min := schema.GetMinimum()
		s.Minimum = &min
	}
	if schema.Maximum != nil {
		max := schema.GetMaximum()
		s.Maximum = &max
	}

	// Apply string constraints
	if schema.GetMinLength() > 0 {
		minLen := schema.GetMinLength()
		s.MinLength = &minLen
	}
	if schema.GetMaxLength() > 0 {
		maxLen := schema.GetMaxLength()
		s.MaxLength = &maxLen
	}

	// Apply array constraints
	if schema.GetMinItems() > 0 {
		minItems := schema.GetMinItems()
		s.MinItems = &minItems
	}
	if schema.GetMaxItems() > 0 {
		maxItems := schema.GetMaxItems()
		s.MaxItems = &maxItems
	}

	// Apply object constraints
	if schema.GetMinProperties() > 0 {
		minProps := schema.GetMinProperties()
		s.MinProperties = &minProps
	}
	if schema.GetMaxProperties() > 0 {
		maxProps := schema.GetMaxProperties()
		s.MaxProperties = &maxProps
	}

	// Apply external docs
	if extDocs := schema.GetExternalDocs(); extDocs != nil {
		s.ExternalDocs = convertExternalDocs(extDocs)
	}

	// Schema composition fields (OpenAPI v3 features)

	// Apply allOf
	if len(schema.GetAllOf()) > 0 {
		for _, allOfSchema := range schema.GetAllOf() {
			s.AllOf = append(s.AllOf, g.convertSchemaOrReference(allOfSchema))
		}
	}

	// Apply anyOf
	if len(schema.GetAnyOf()) > 0 {
		for _, anyOfSchema := range schema.GetAnyOf() {
			s.AnyOf = append(s.AnyOf, g.convertSchemaOrReference(anyOfSchema))
		}
	}

	// Apply oneOf
	if len(schema.GetOneOf()) > 0 {
		for _, oneOfSchema := range schema.GetOneOf() {
			s.OneOf = append(s.OneOf, g.convertSchemaOrReference(oneOfSchema))
		}
	}

	// Apply not
	if notSchema := schema.GetNot(); notSchema != nil {
		s.Not = g.convertSchemaOrReference(notSchema)
	}

	// Apply discriminator
	if disc := schema.GetDiscriminator(); disc != nil {
		s.Discriminator = &Discriminator{
			PropertyName: disc.GetPropertyName(),
			Mapping:      disc.GetMapping(),
		}
	}

	// Apply items (for array schemas)
	if items := schema.GetItems(); items != nil {
		s.Items = g.convertSchemaOrReference(items)
	}

	// Apply properties (for object schemas) - now uses NamedSchemaOrReference for ordering
	if len(schema.GetProperties()) > 0 {
		s.Properties = make(map[string]*SchemaOrReference)
		for _, namedProp := range schema.GetProperties() {
			s.Properties[namedProp.GetName()] = g.convertSchemaOrReference(namedProp.GetValue())
		}
	}

	// Apply additionalProperties
	if addProps := schema.GetAdditionalProperties(); addProps != nil {
		switch kind := addProps.GetKind().(type) {
		case *options.AdditionalPropertiesItem_Allows:
			// Boolean true allows any additional properties
			if kind.Allows {
				s.AdditionalProperties = &SchemaOrReference{Schema: &Schema{}}
			}
			// Boolean false - we don't set anything (default is no additional properties)
		case *options.AdditionalPropertiesItem_SchemaOrReference:
			s.AdditionalProperties = g.convertSchemaOrReference(kind.SchemaOrReference)
		}
	}

	// Note: prefixItems, propertyNames, and patternProperties are JSON Schema
	// draft 2020-12 features not directly supported in OpenAPI 3.0.x.
	// They are defined in the proto but not converted here.

	return &SchemaOrReference{Schema: s}
}
