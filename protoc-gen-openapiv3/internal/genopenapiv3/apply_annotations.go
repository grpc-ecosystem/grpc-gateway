package genopenapiv3

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
)

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
		info.License = &License{
			Name:       l.GetName(),
			Identifier: l.GetIdentifier(),
			URL:        l.GetUrl(),
		}
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
	for name, resp := range opts.GetResponses() {
		if comp.Responses == nil {
			comp.Responses = make(map[string]*ResponseRef)
		}
		comp.Responses[name] = &ResponseRef{Value: convertResponse(resp)}
	}

	// Apply Parameters
	for name, param := range opts.GetParameters() {
		if comp.Parameters == nil {
			comp.Parameters = make(map[string]*ParameterRef)
		}
		comp.Parameters[name] = &ParameterRef{Value: convertParameter(param)}
	}

	// Apply Request Bodies
	for name, body := range opts.GetRequestBodies() {
		if comp.RequestBodies == nil {
			comp.RequestBodies = make(map[string]*RequestBodyRef)
		}
		comp.RequestBodies[name] = &RequestBodyRef{Value: convertRequestBody(body)}
	}

	// Apply Headers
	for name, header := range opts.GetHeaders() {
		if comp.Headers == nil {
			comp.Headers = make(map[string]*HeaderRef)
		}
		comp.Headers[name] = &HeaderRef{Value: convertHeader(header)}
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
	for code, resp := range opts.GetResponses() {
		if op.Responses == nil {
			op.Responses = NewResponses()
		}
		op.Responses.Codes[code] = &ResponseRef{Value: convertResponse(resp)}
	}

	// Apply request body override if provided
	if reqBody := opts.GetRequestBody(); reqBody != nil {
		op.RequestBody = &RequestBodyRef{Value: convertRequestBody(reqBody)}
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

	// Apply example
	if opts.GetExample() != "" {
		schema.Example = opts.GetExample()
	}

	// Apply read only
	if opts.GetReadOnly() {
		schema.ReadOnly = true
	}

	// Apply write only
	if opts.GetWriteOnly() {
		schema.WriteOnly = true
	}

	// Apply nullable
	if opts.GetNullable() {
		schema.Nullable = true
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

	// Apply example
	if opts.GetExample() != "" {
		schema.Example = opts.GetExample()
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
	if opts.GetMultipleOf() != 0 {
		multipleOf := opts.GetMultipleOf()
		schema.MultipleOf = &multipleOf
	}
	if opts.GetMinimum() != 0 {
		min := opts.GetMinimum()
		schema.Minimum = &min
	}
	if opts.GetMaximum() != 0 {
		max := opts.GetMaximum()
		schema.Maximum = &max
	}
	if opts.GetExclusiveMinimum() {
		// In OpenAPI 3.0, exclusiveMinimum is boolean, but in 3.1 it's the actual value
		// For now, we set it to minimum value if exclusiveMinimum is true
		if schema.Minimum != nil {
			schema.ExclusiveMinimum = schema.Minimum
			schema.Minimum = nil
		}
	}
	if opts.GetExclusiveMaximum() {
		if schema.Maximum != nil {
			schema.ExclusiveMaximum = schema.Maximum
			schema.Maximum = nil
		}
	}

	// Apply read only
	if opts.GetReadOnly() {
		schema.ReadOnly = true
	}

	// Apply write only
	if opts.GetWriteOnly() {
		schema.WriteOnly = true
	}

	// Apply nullable
	if opts.GetNullable() {
		schema.Nullable = true
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

	// Apply example
	if opts.GetExample() != "" {
		schema.Example = opts.GetExample()
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
		result[name] = val.GetScope()
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

func convertResponse(resp *options.Response) *Response {
	r := &Response{
		Description: resp.GetDescription(),
	}
	if len(resp.GetHeaders()) > 0 {
		r.Headers = make(map[string]*HeaderRef)
		for name, h := range resp.GetHeaders() {
			r.Headers[name] = &HeaderRef{Value: convertHeader(h)}
		}
	}
	if len(resp.GetContent()) > 0 {
		r.Content = make(map[string]*MediaType)
		for mediaType, mt := range resp.GetContent() {
			r.Content[mediaType] = convertMediaType(mt)
		}
	}
	return r
}

func convertParameter(param *options.Parameter) *Parameter {
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
	if param.GetExample() != "" {
		p.Example = param.GetExample()
	}
	if schema := param.GetSchema(); schema != nil {
		p.Schema = convertSchema(schema)
	}
	return p
}

func convertRequestBody(body *options.RequestBody) *RequestBody {
	rb := &RequestBody{
		Description: body.GetDescription(),
		Required:    body.GetRequired(),
	}
	if len(body.GetContent()) > 0 {
		rb.Content = make(map[string]*MediaType)
		for mediaType, mt := range body.GetContent() {
			rb.Content[mediaType] = convertMediaType(mt)
		}
	}
	return rb
}

func convertHeader(header *options.Header) *Header {
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
	if header.GetExample() != "" {
		h.Example = header.GetExample()
	}
	if schema := header.GetSchema(); schema != nil {
		h.Schema = convertSchema(schema)
	}
	return h
}

func convertMediaType(mt *options.MediaType) *MediaType {
	result := &MediaType{}
	if mt.GetExample() != "" {
		result.Example = mt.GetExample()
	}
	if schema := mt.GetSchema(); schema != nil {
		result.Schema = convertSchema(schema)
	}
	return result
}

func convertSchema(schema *options.Schema) *SchemaRef {
	// If it's a reference, return a reference
	if schema.GetRef() != "" {
		return &SchemaRef{Ref: schema.GetRef()}
	}

	s := &Schema{
		Type:        schema.GetType(),
		Format:      schema.GetFormat(),
		Title:       schema.GetTitle(),
		Description: schema.GetDescription(),
		ReadOnly:    schema.GetReadOnly(),
		WriteOnly:   schema.GetWriteOnly(),
		Nullable:    schema.GetNullable(),
		Deprecated:  schema.GetDeprecated(),
		Pattern:     schema.GetPattern(),
		Required:    schema.GetRequired(),
		UniqueItems: schema.GetUniqueItems(),
	}

	// Apply default
	if schema.GetDefault() != "" {
		s.Default = schema.GetDefault()
	}

	// Apply example
	if schema.GetExample() != "" {
		s.Example = schema.GetExample()
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
	if schema.GetMinimum() != 0 {
		min := schema.GetMinimum()
		s.Minimum = &min
	}
	if schema.GetMaximum() != 0 {
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
			s.AllOf = append(s.AllOf, convertSchema(allOfSchema))
		}
	}

	// Apply anyOf
	if len(schema.GetAnyOf()) > 0 {
		for _, anyOfSchema := range schema.GetAnyOf() {
			s.AnyOf = append(s.AnyOf, convertSchema(anyOfSchema))
		}
	}

	// Apply oneOf
	if len(schema.GetOneOf()) > 0 {
		for _, oneOfSchema := range schema.GetOneOf() {
			s.OneOf = append(s.OneOf, convertSchema(oneOfSchema))
		}
	}

	// Apply not
	if notSchema := schema.GetNot(); notSchema != nil {
		s.Not = convertSchema(notSchema)
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
		s.Items = convertSchema(items)
	}

	// Apply properties (for object schemas)
	if len(schema.GetProperties()) > 0 {
		s.Properties = make(map[string]*SchemaRef)
		for name, propSchema := range schema.GetProperties() {
			s.Properties[name] = convertSchema(propSchema)
		}
	}

	// Apply additionalProperties
	if addProps := schema.GetAdditionalProperties(); addProps != nil {
		switch kind := addProps.GetKind().(type) {
		case *options.AdditionalPropertiesItem_Allows:
			// Boolean true allows any additional properties
			if kind.Allows {
				s.AdditionalProperties = &SchemaRef{Value: &Schema{}}
			}
			// Boolean false - we don't set anything (default is no additional properties)
		case *options.AdditionalPropertiesItem_Schema:
			s.AdditionalProperties = convertSchema(kind.Schema)
		}
	}

	// Apply prefixItems (tuple validation)
	if len(schema.GetPrefixItems()) > 0 {
		// Note: prefixItems is a JSON Schema draft 2020-12 feature
		// For OpenAPI 3.0.x, this is represented differently
		// We'll store as items for now since OpenAPI 3.0 doesn't support prefixItems
	}

	// Apply propertyNames
	if propNames := schema.GetPropertyNames(); propNames != nil {
		// Note: propertyNames is not directly supported in OpenAPI 3.0.x
		// but is part of JSON Schema
	}

	// Apply patternProperties
	if len(schema.GetPatternProperties()) > 0 {
		// Note: patternProperties is not directly supported in OpenAPI 3.0.x
		// but is part of JSON Schema
	}

	return &SchemaRef{Value: s}
}
