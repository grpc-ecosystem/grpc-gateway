package genopenapiv3

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
)

// ToCanonical converts the internal OpenAPI types to the canonical model.
// This allows us to use the adapter pattern for version-specific serialization
// without rewriting the entire generator.

func (o *OpenAPI) ToCanonical() *model.Document {
	if o == nil {
		return nil
	}

	doc := &model.Document{
		OpenAPIVersion: o.OpenAPI,
		Info:           toCanonicalInfo(o.Info),
		Servers:        toCanonicalServers(o.Servers),
		Paths:          toCanonicalPaths(o.Paths),
		Components:     toCanonicalComponents(o.Components),
		Security:       toCanonicalSecurityRequirements(o.Security),
		Tags:           toCanonicalTags(o.Tags),
		ExternalDocs:   toCanonicalExternalDocs(o.ExternalDocs),
	}

	return doc
}

func toCanonicalInfo(info *Info) *model.Info {
	if info == nil {
		return nil
	}
	return &model.Info{
		Title:          info.Title,
		Description:    info.Description,
		TermsOfService: info.TermsOfService,
		Contact:        toCanonicalContact(info.Contact),
		License:        toCanonicalLicense(info.License),
		Version:        info.Version,
		Summary:        info.Summary,
	}
}

func toCanonicalContact(c *Contact) *model.Contact {
	if c == nil {
		return nil
	}
	return &model.Contact{
		Name:  c.Name,
		URL:   c.URL,
		Email: c.Email,
	}
}

func toCanonicalLicense(l *License) *model.License {
	if l == nil {
		return nil
	}
	return &model.License{
		Name:       l.Name,
		URL:        l.URL,
		Identifier: l.Identifier,
	}
}

func toCanonicalServers(servers []*Server) []*model.Server {
	if len(servers) == 0 {
		return nil
	}
	result := make([]*model.Server, len(servers))
	for i, s := range servers {
		result[i] = toCanonicalServer(s)
	}
	return result
}

func toCanonicalServer(s *Server) *model.Server {
	if s == nil {
		return nil
	}
	var vars map[string]*model.ServerVariable
	if len(s.Variables) > 0 {
		vars = make(map[string]*model.ServerVariable)
		for name, v := range s.Variables {
			vars[name] = &model.ServerVariable{
				Enum:        v.Enum,
				Default:     v.Default,
				Description: v.Description,
			}
		}
	}
	return &model.Server{
		URL:         s.URL,
		Description: s.Description,
		Variables:   vars,
	}
}

func toCanonicalPaths(paths *Paths) *model.Paths {
	if paths == nil || len(paths.paths) == 0 {
		return nil
	}
	result := &model.Paths{
		Items: make(map[string]*model.PathItem),
	}
	for _, path := range paths.order {
		item := paths.paths[path]
		if item != nil {
			result.Items[path] = toCanonicalPathItem(item)
		}
	}
	return result
}

func toCanonicalPathItem(item *PathItem) *model.PathItem {
	if item == nil {
		return nil
	}
	return &model.PathItem{
		Ref:         item.Ref,
		Summary:     item.Summary,
		Description: item.Description,
		Get:         toCanonicalOperation(item.Get),
		Put:         toCanonicalOperation(item.Put),
		Post:        toCanonicalOperation(item.Post),
		Delete:      toCanonicalOperation(item.Delete),
		Options:     toCanonicalOperation(item.Options),
		Head:        toCanonicalOperation(item.Head),
		Patch:       toCanonicalOperation(item.Patch),
		Trace:       toCanonicalOperation(item.Trace),
		Servers:     toCanonicalServers(item.Servers),
		Parameters:  toCanonicalParameterRefs(item.Parameters),
	}
}

func toCanonicalOperation(op *Operation) *model.Operation {
	if op == nil {
		return nil
	}
	return &model.Operation{
		Tags:         op.Tags,
		Summary:      op.Summary,
		Description:  op.Description,
		ExternalDocs: toCanonicalExternalDocs(op.ExternalDocs),
		OperationID:  op.OperationID,
		Parameters:   toCanonicalParameterRefs(op.Parameters),
		RequestBody:  toCanonicalRequestBodyRef(op.RequestBody),
		Responses:    toCanonicalResponses(op.Responses),
		Callbacks:    toCanonicalCallbacks(op.Callbacks),
		Deprecated:   op.Deprecated,
		Security:     toCanonicalSecurityRequirements(op.Security),
		Servers:      toCanonicalServers(op.Servers),
	}
}

func toCanonicalExternalDocs(ed *ExternalDocumentation) *model.ExternalDocs {
	if ed == nil {
		return nil
	}
	return &model.ExternalDocs{
		Description: ed.Description,
		URL:         ed.URL,
	}
}

func toCanonicalParameterRefs(params []*ParameterRef) []*model.ParameterOrRef {
	if len(params) == 0 {
		return nil
	}
	result := make([]*model.ParameterOrRef, len(params))
	for i, p := range params {
		result[i] = toCanonicalParameterRef(p)
	}
	return result
}

func toCanonicalParameterRef(p *ParameterRef) *model.ParameterOrRef {
	if p == nil {
		return nil
	}
	if p.Ref != "" {
		return &model.ParameterOrRef{Ref: p.Ref}
	}
	if p.Value == nil {
		return nil
	}
	return &model.ParameterOrRef{
		Value: &model.Parameter{
			Name:            p.Value.Name,
			In:              p.Value.In,
			Description:     p.Value.Description,
			Required:        p.Value.Required,
			Deprecated:      p.Value.Deprecated,
			AllowEmptyValue: p.Value.AllowEmptyValue,
			Style:           p.Value.Style,
			Explode:         p.Value.Explode,
			AllowReserved:   p.Value.AllowReserved,
			Schema:          toCanonicalSchemaOrReference(p.Value.Schema),
			Examples:        toCanonicalExampleRefs(p.Value.Examples),
			Content:         toCanonicalContent(p.Value.Content),
		},
	}
}

func toCanonicalRequestBodyRef(rb *RequestBodyRef) *model.RequestBodyOrRef {
	if rb == nil {
		return nil
	}
	if rb.Ref != "" {
		return &model.RequestBodyOrRef{Ref: rb.Ref}
	}
	if rb.Value == nil {
		return nil
	}
	return &model.RequestBodyOrRef{
		Value: &model.RequestBody{
			Description: rb.Value.Description,
			Content:     toCanonicalContent(rb.Value.Content),
			Required:    rb.Value.Required,
		},
	}
}

func toCanonicalContent(content map[string]*MediaType) map[string]*model.MediaType {
	if len(content) == 0 {
		return nil
	}
	result := make(map[string]*model.MediaType)
	for mt, media := range content {
		result[mt] = toCanonicalMediaType(media)
	}
	return result
}

func toCanonicalMediaType(mt *MediaType) *model.MediaType {
	if mt == nil {
		return nil
	}
	return &model.MediaType{
		Schema:   toCanonicalSchemaOrReference(mt.Schema),
		Examples: toCanonicalExampleRefs(mt.Examples),
		Encoding: toCanonicalEncoding(mt.Encoding),
	}
}

func toCanonicalEncoding(encoding map[string]*Encoding) map[string]*model.Encoding {
	if len(encoding) == 0 {
		return nil
	}
	result := make(map[string]*model.Encoding)
	for name, enc := range encoding {
		result[name] = &model.Encoding{
			ContentType:   enc.ContentType,
			Headers:       toCanonicalHeaderRefs(enc.Headers),
			Style:         enc.Style,
			Explode:       enc.Explode,
			AllowReserved: enc.AllowReserved,
		}
	}
	return result
}

func toCanonicalResponses(r *Responses) *model.Responses {
	if r == nil {
		return nil
	}
	result := &model.Responses{
		Default: toCanonicalResponseRef(r.Default),
		Codes:   make(map[string]*model.ResponseOrRef),
	}
	for code, resp := range r.Codes {
		result.Codes[code] = toCanonicalResponseRef(resp)
	}
	return result
}

func toCanonicalResponseRef(r *ResponseRef) *model.ResponseOrRef {
	if r == nil {
		return nil
	}
	if r.Ref != "" {
		return &model.ResponseOrRef{Ref: r.Ref}
	}
	if r.Value == nil {
		return nil
	}
	return &model.ResponseOrRef{
		Value: &model.Response{
			Description: r.Value.Description,
			Headers:     toCanonicalHeaderRefs(r.Value.Headers),
			Content:     toCanonicalContent(r.Value.Content),
			Links:       toCanonicalLinkRefs(r.Value.Links),
		},
	}
}

func toCanonicalHeaderRefs(headers map[string]*HeaderRef) map[string]*model.HeaderOrRef {
	if len(headers) == 0 {
		return nil
	}
	result := make(map[string]*model.HeaderOrRef)
	for name, h := range headers {
		result[name] = toCanonicalHeaderRef(h)
	}
	return result
}

func toCanonicalHeaderRef(h *HeaderRef) *model.HeaderOrRef {
	if h == nil {
		return nil
	}
	if h.Ref != "" {
		return &model.HeaderOrRef{Ref: h.Ref}
	}
	if h.Value == nil {
		return nil
	}
	return &model.HeaderOrRef{
		Value: &model.Header{
			Description:     h.Value.Description,
			Required:        h.Value.Required,
			Deprecated:      h.Value.Deprecated,
			AllowEmptyValue: h.Value.AllowEmptyValue,
			Style:           h.Value.Style,
			Explode:         h.Value.Explode,
			AllowReserved:   h.Value.AllowReserved,
			Schema:          toCanonicalSchemaOrReference(h.Value.Schema),
			Examples:        toCanonicalExampleRefs(h.Value.Examples),
			Content:         toCanonicalContent(h.Value.Content),
		},
	}
}

func toCanonicalLinkRefs(links map[string]*LinkRef) map[string]*model.LinkOrRef {
	if len(links) == 0 {
		return nil
	}
	result := make(map[string]*model.LinkOrRef)
	for name, l := range links {
		result[name] = toCanonicalLinkRef(l)
	}
	return result
}

func toCanonicalLinkRef(l *LinkRef) *model.LinkOrRef {
	if l == nil {
		return nil
	}
	if l.Ref != "" {
		return &model.LinkOrRef{Ref: l.Ref}
	}
	if l.Value == nil {
		return nil
	}
	return &model.LinkOrRef{
		Value: &model.Link{
			OperationRef: l.Value.OperationRef,
			OperationID:  l.Value.OperationId,
			Parameters:   l.Value.Parameters,
			RequestBody:  l.Value.RequestBody,
			Description:  l.Value.Description,
			Server:       toCanonicalServer(l.Value.Server),
		},
	}
}

func toCanonicalCallbacks(callbacks map[string]*CallbackRef) map[string]*model.CallbackOrRef {
	if len(callbacks) == 0 {
		return nil
	}
	result := make(map[string]*model.CallbackOrRef)
	for name, cb := range callbacks {
		result[name] = toCanonicalCallbackRef(cb)
	}
	return result
}

func toCanonicalCallbackRef(cb *CallbackRef) *model.CallbackOrRef {
	if cb == nil {
		return nil
	}
	if cb.Ref != "" {
		return &model.CallbackOrRef{Ref: cb.Ref}
	}
	if cb.Value == nil {
		return nil
	}
	pathItems := make(map[string]*model.PathItem)
	for expr, item := range *cb.Value {
		pathItems[expr] = toCanonicalPathItem(item)
	}
	return &model.CallbackOrRef{Value: pathItems}
}

func toCanonicalSchemaOrReference(s *SchemaOrReference) *model.SchemaOrRef {
	if s == nil {
		return nil
	}
	if s.Reference != nil {
		return &model.SchemaOrRef{
			Ref:         s.Reference.Ref,
			Summary:     s.Reference.Summary,
			Description: s.Reference.Description,
		}
	}
	if s.Schema == nil {
		return nil
	}
	return &model.SchemaOrRef{
		Value: toCanonicalSchema(s.Schema),
	}
}

func toCanonicalSchema(s *Schema) *model.Schema {
	if s == nil {
		return nil
	}

	// Determine type and nullability from SchemaType
	var typ string
	var isNullable bool
	if len(s.Type) > 0 {
		typ = s.Type[0]
		// Check if "null" is in the type array (3.1.0 style)
		for _, t := range s.Type {
			if t == "null" {
				isNullable = true
				break
			}
		}
	}
	// Also check the nullable boolean (3.0.x style)
	if s.Nullable {
		isNullable = true
	}

	return &model.Schema{
		Type:                 typ,
		Format:               s.Format,
		IsNullable:           isNullable,
		Title:                s.Title,
		Description:          s.Description,
		Default:              s.Default,
		Examples:             s.Examples, // JSON Schema: array of example values
		Deprecated:           s.Deprecated,
		ReadOnly:             s.ReadOnly,
		WriteOnly:            s.WriteOnly,
		ExternalDocs:         toCanonicalExternalDocs(s.ExternalDocs),
		MultipleOf:           s.MultipleOf,
		Minimum:              s.Minimum,
		Maximum:              s.Maximum,
		ExclusiveMinimum:     s.ExclusiveMinimum,
		ExclusiveMaximum:     s.ExclusiveMaximum,
		MinLength:            s.MinLength,
		MaxLength:            s.MaxLength,
		Pattern:              s.Pattern,
		MinItems:             s.MinItems,
		MaxItems:             s.MaxItems,
		UniqueItems:          s.UniqueItems,
		Items:                toCanonicalSchemaOrReference(s.Items),
		MinProperties:        s.MinProperties,
		MaxProperties:        s.MaxProperties,
		Required:             s.Required,
		Properties:           toCanonicalSchemaProperties(s.Properties),
		AdditionalProperties: toCanonicalAdditionalProperties(s.AdditionalProperties),
		AllOf:                toCanonicalSchemaOrReferences(s.AllOf),
		AnyOf:                toCanonicalSchemaOrReferences(s.AnyOf),
		OneOf:                toCanonicalSchemaOrReferences(s.OneOf),
		Not:                  toCanonicalSchemaOrReference(s.Not),
		Discriminator:        toCanonicalDiscriminator(s.Discriminator),
		Enum:                 s.Enum,
	}
}

func toCanonicalSchemaProperties(props map[string]*SchemaOrReference) map[string]*model.SchemaOrRef {
	if len(props) == 0 {
		return nil
	}
	result := make(map[string]*model.SchemaOrRef)
	for name, prop := range props {
		result[name] = toCanonicalSchemaOrReference(prop)
	}
	return result
}

func toCanonicalAdditionalProperties(ap *SchemaOrReference) *model.AdditionalProperties {
	if ap == nil {
		return nil
	}
	// If it's a schema reference or inline schema
	if ap.Reference != nil || ap.Schema != nil {
		return &model.AdditionalProperties{
			Schema: toCanonicalSchemaOrReference(ap),
		}
	}
	return nil
}

func toCanonicalSchemaOrReferences(schemas []*SchemaOrReference) []*model.SchemaOrRef {
	if len(schemas) == 0 {
		return nil
	}
	result := make([]*model.SchemaOrRef, len(schemas))
	for i, s := range schemas {
		result[i] = toCanonicalSchemaOrReference(s)
	}
	return result
}

func toCanonicalDiscriminator(d *Discriminator) *model.Discriminator {
	if d == nil {
		return nil
	}
	return &model.Discriminator{
		PropertyName: d.PropertyName,
		Mapping:      d.Mapping,
	}
}

func toCanonicalExampleRefs(examples map[string]*ExampleRef) map[string]*model.Example {
	if len(examples) == 0 {
		return nil
	}
	result := make(map[string]*model.Example)
	for name, ex := range examples {
		if ex == nil {
			continue
		}
		if ex.Ref != "" {
			result[name] = &model.Example{Ref: ex.Ref}
		} else if ex.Value != nil {
			result[name] = &model.Example{
				Summary:       ex.Value.Summary,
				Description:   ex.Value.Description,
				Value:         ex.Value.Value,
				ExternalValue: ex.Value.ExternalValue,
			}
		}
	}
	return result
}

func toCanonicalComponents(c *Components) *model.Components {
	if c == nil {
		return nil
	}
	return &model.Components{
		Schemas:         toCanonicalComponentSchemas(c.Schemas),
		Responses:       toCanonicalComponentResponses(c.Responses),
		Parameters:      toCanonicalComponentParameters(c.Parameters),
		Examples:        toCanonicalComponentExamples(c.Examples),
		RequestBodies:   toCanonicalComponentRequestBodies(c.RequestBodies),
		Headers:         toCanonicalComponentHeaders(c.Headers),
		SecuritySchemes: toCanonicalComponentSecuritySchemes(c.SecuritySchemes),
		Links:           toCanonicalComponentLinks(c.Links),
		Callbacks:       toCanonicalComponentCallbacks(c.Callbacks),
	}
}

func toCanonicalComponentSchemas(schemas map[string]*SchemaOrReference) map[string]*model.SchemaOrRef {
	if len(schemas) == 0 {
		return nil
	}
	result := make(map[string]*model.SchemaOrRef)
	for name, s := range schemas {
		result[name] = toCanonicalSchemaOrReference(s)
	}
	return result
}

func toCanonicalComponentResponses(responses map[string]*ResponseRef) map[string]*model.ResponseOrRef {
	if len(responses) == 0 {
		return nil
	}
	result := make(map[string]*model.ResponseOrRef)
	for name, r := range responses {
		result[name] = toCanonicalResponseRef(r)
	}
	return result
}

func toCanonicalComponentParameters(params map[string]*ParameterRef) map[string]*model.ParameterOrRef {
	if len(params) == 0 {
		return nil
	}
	result := make(map[string]*model.ParameterOrRef)
	for name, p := range params {
		result[name] = toCanonicalParameterRef(p)
	}
	return result
}

func toCanonicalComponentExamples(examples map[string]*ExampleRef) map[string]*model.Example {
	if len(examples) == 0 {
		return nil
	}
	result := make(map[string]*model.Example)
	for name, ex := range examples {
		if ex == nil {
			continue
		}
		if ex.Ref != "" {
			result[name] = &model.Example{Ref: ex.Ref}
		} else if ex.Value != nil {
			result[name] = &model.Example{
				Summary:       ex.Value.Summary,
				Description:   ex.Value.Description,
				Value:         ex.Value.Value,
				ExternalValue: ex.Value.ExternalValue,
			}
		}
	}
	return result
}

func toCanonicalComponentRequestBodies(bodies map[string]*RequestBodyRef) map[string]*model.RequestBodyOrRef {
	if len(bodies) == 0 {
		return nil
	}
	result := make(map[string]*model.RequestBodyOrRef)
	for name, rb := range bodies {
		result[name] = toCanonicalRequestBodyRef(rb)
	}
	return result
}

func toCanonicalComponentHeaders(headers map[string]*HeaderRef) map[string]*model.HeaderOrRef {
	if len(headers) == 0 {
		return nil
	}
	result := make(map[string]*model.HeaderOrRef)
	for name, h := range headers {
		result[name] = toCanonicalHeaderRef(h)
	}
	return result
}

func toCanonicalComponentSecuritySchemes(schemes map[string]*SecuritySchemeRef) map[string]*model.SecuritySchemeOrRef {
	if len(schemes) == 0 {
		return nil
	}
	result := make(map[string]*model.SecuritySchemeOrRef)
	for name, ss := range schemes {
		result[name] = toCanonicalSecuritySchemeRef(ss)
	}
	return result
}

func toCanonicalSecuritySchemeRef(ss *SecuritySchemeRef) *model.SecuritySchemeOrRef {
	if ss == nil {
		return nil
	}
	if ss.Ref != "" {
		return &model.SecuritySchemeOrRef{Ref: ss.Ref}
	}
	if ss.Value == nil {
		return nil
	}
	return &model.SecuritySchemeOrRef{
		Value: &model.SecurityScheme{
			Type:             ss.Value.Type,
			Description:      ss.Value.Description,
			Name:             ss.Value.Name,
			In:               ss.Value.In,
			Scheme:           ss.Value.Scheme,
			BearerFormat:     ss.Value.BearerFormat,
			Flows:            toCanonicalOAuthFlows(ss.Value.Flows),
			OpenIDConnectURL: ss.Value.OpenIdConnectUrl,
		},
	}
}

func toCanonicalOAuthFlows(flows *OAuthFlows) *model.OAuthFlows {
	if flows == nil {
		return nil
	}
	return &model.OAuthFlows{
		Implicit:          toCanonicalOAuthFlow(flows.Implicit),
		Password:          toCanonicalOAuthFlow(flows.Password),
		ClientCredentials: toCanonicalOAuthFlow(flows.ClientCredentials),
		AuthorizationCode: toCanonicalOAuthFlow(flows.AuthorizationCode),
	}
}

func toCanonicalOAuthFlow(flow *OAuthFlow) *model.OAuthFlow {
	if flow == nil {
		return nil
	}
	return &model.OAuthFlow{
		AuthorizationURL: flow.AuthorizationURL,
		TokenURL:         flow.TokenURL,
		RefreshURL:       flow.RefreshURL,
		Scopes:           flow.Scopes,
	}
}

func toCanonicalComponentLinks(links map[string]*LinkRef) map[string]*model.LinkOrRef {
	if len(links) == 0 {
		return nil
	}
	result := make(map[string]*model.LinkOrRef)
	for name, l := range links {
		result[name] = toCanonicalLinkRef(l)
	}
	return result
}

func toCanonicalComponentCallbacks(callbacks map[string]*CallbackRef) map[string]*model.CallbackOrRef {
	if len(callbacks) == 0 {
		return nil
	}
	result := make(map[string]*model.CallbackOrRef)
	for name, cb := range callbacks {
		result[name] = toCanonicalCallbackRef(cb)
	}
	return result
}

func toCanonicalSecurityRequirements(reqs []SecurityRequirement) []model.SecurityRequirement {
	if len(reqs) == 0 {
		return nil
	}
	result := make([]model.SecurityRequirement, len(reqs))
	for i, req := range reqs {
		result[i] = model.SecurityRequirement(req)
	}
	return result
}

func toCanonicalTags(tags []*Tag) []*model.Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]*model.Tag, len(tags))
	for i, t := range tags {
		result[i] = &model.Tag{
			Name:         t.Name,
			Description:  t.Description,
			ExternalDocs: toCanonicalExternalDocs(t.ExternalDocs),
		}
	}
	return result
}
