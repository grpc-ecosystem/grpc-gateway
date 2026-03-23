package v31

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
	"go.yaml.in/yaml/v3"
)

// Format represents output format.
type Format int

const (
	FormatJSON Format = iota
	FormatYAML
)

// Adapter converts canonical model to OpenAPI 3.1.0 output.
//
// OpenAPI 3.1.0 specifics handled by this adapter:
// - nullable via type arrays: ["string", "null"]
// - examples as map, not singular example
// - $ref can have sibling properties (summary, description)
// - webhooks support
// - JSON Schema draft 2020-12 alignment
type Adapter struct{}

// New creates a new OpenAPI 3.1.0 adapter.
func New() *Adapter {
	return &Adapter{}
}

// Version returns "3.1.0".
func (a *Adapter) Version() string {
	return "3.1.0"
}

// Adapt converts a canonical Document to OpenAPI 3.1.0 JSON or YAML.
func (a *Adapter) Adapt(doc *model.Document, format Format) ([]byte, error) {
	output := a.adaptDocument(doc)

	switch format {
	case FormatYAML:
		return yaml.Marshal(output)
	default:
		return json.MarshalIndent(output, "", "  ")
	}
}

// adaptDocument converts canonical Document to output map.
func (a *Adapter) adaptDocument(doc *model.Document) map[string]any {
	if doc == nil {
		return nil
	}

	out := map[string]any{
		"openapi": "3.1.0",
	}

	if doc.Info != nil {
		out["info"] = a.adaptInfo(doc.Info)
	}

	if len(doc.Servers) > 0 {
		servers := make([]any, len(doc.Servers))
		for i, s := range doc.Servers {
			servers[i] = a.adaptServer(s)
		}
		out["servers"] = servers
	}

	if doc.Paths != nil && len(doc.Paths.Items) > 0 {
		paths := make(map[string]any)
		for path, item := range doc.Paths.Items {
			paths[path] = a.adaptPathItem(item)
		}
		out["paths"] = paths
	}

	if comp := a.adaptComponents(doc.Components); len(comp) > 0 {
		out["components"] = comp
	}

	if len(doc.Security) > 0 {
		out["security"] = doc.Security
	}

	if len(doc.Tags) > 0 {
		tags := make([]any, len(doc.Tags))
		for i, t := range doc.Tags {
			tags[i] = a.adaptTag(t)
		}
		out["tags"] = tags
	}

	if doc.ExternalDocs != nil {
		out["externalDocs"] = a.adaptExternalDocs(doc.ExternalDocs)
	}

	if doc.JSONSchemaDialect != "" {
		out["jsonSchemaDialect"] = doc.JSONSchemaDialect
	}

	return out
}

func (a *Adapter) adaptInfo(info *model.Info) map[string]any {
	if info == nil {
		return nil
	}

	out := map[string]any{}
	setIfNotEmpty(out, "title", info.Title)
	setIfNotEmpty(out, "description", info.Description)
	setIfNotEmpty(out, "termsOfService", info.TermsOfService)
	setIfNotEmpty(out, "version", info.Version)
	setIfNotEmpty(out, "summary", info.Summary) // 3.1.0 only

	if info.Contact != nil {
		out["contact"] = a.adaptContact(info.Contact)
	}
	if info.License != nil {
		out["license"] = a.adaptLicense(info.License)
	}

	return out
}

func (a *Adapter) adaptContact(c *model.Contact) map[string]any {
	if c == nil {
		return nil
	}
	out := map[string]any{}
	setIfNotEmpty(out, "name", c.Name)
	setIfNotEmpty(out, "url", c.URL)
	setIfNotEmpty(out, "email", c.Email)
	return out
}

func (a *Adapter) adaptLicense(l *model.License) map[string]any {
	if l == nil {
		return nil
	}
	out := map[string]any{}
	setIfNotEmpty(out, "name", l.Name)
	setIfNotEmpty(out, "url", l.URL)
	setIfNotEmpty(out, "identifier", l.Identifier) // 3.1.0 only
	return out
}

func (a *Adapter) adaptServer(s *model.Server) map[string]any {
	if s == nil {
		return nil
	}
	out := map[string]any{
		"url": s.URL,
	}
	setIfNotEmpty(out, "description", s.Description)

	if len(s.Variables) > 0 {
		vars := make(map[string]any)
		for name, v := range s.Variables {
			vars[name] = map[string]any{
				"default":     v.Default,
				"description": v.Description,
				"enum":        v.Enum,
			}
		}
		out["variables"] = vars
	}
	return out
}

func (a *Adapter) adaptPathItem(item *model.PathItem) map[string]any {
	if item == nil {
		return nil
	}
	out := map[string]any{}

	setIfNotEmpty(out, "$ref", item.Ref)
	setIfNotEmpty(out, "summary", item.Summary)
	setIfNotEmpty(out, "description", item.Description)

	if item.Get != nil {
		out["get"] = a.adaptOperation(item.Get)
	}
	if item.Put != nil {
		out["put"] = a.adaptOperation(item.Put)
	}
	if item.Post != nil {
		out["post"] = a.adaptOperation(item.Post)
	}
	if item.Delete != nil {
		out["delete"] = a.adaptOperation(item.Delete)
	}
	if item.Options != nil {
		out["options"] = a.adaptOperation(item.Options)
	}
	if item.Head != nil {
		out["head"] = a.adaptOperation(item.Head)
	}
	if item.Patch != nil {
		out["patch"] = a.adaptOperation(item.Patch)
	}
	if item.Trace != nil {
		out["trace"] = a.adaptOperation(item.Trace)
	}

	if len(item.Servers) > 0 {
		servers := make([]any, len(item.Servers))
		for i, s := range item.Servers {
			servers[i] = a.adaptServer(s)
		}
		out["servers"] = servers
	}

	if len(item.Parameters) > 0 {
		params := make([]any, len(item.Parameters))
		for i, p := range item.Parameters {
			params[i] = a.adaptParameterOrRef(p)
		}
		out["parameters"] = params
	}

	return out
}

func (a *Adapter) adaptOperation(op *model.Operation) map[string]any {
	if op == nil {
		return nil
	}
	out := map[string]any{}

	if len(op.Tags) > 0 {
		out["tags"] = op.Tags
	}
	setIfNotEmpty(out, "summary", op.Summary)
	setIfNotEmpty(out, "description", op.Description)
	setIfNotEmpty(out, "operationId", op.OperationID)

	if op.ExternalDocs != nil {
		out["externalDocs"] = a.adaptExternalDocs(op.ExternalDocs)
	}

	if len(op.Parameters) > 0 {
		params := make([]any, len(op.Parameters))
		for i, p := range op.Parameters {
			params[i] = a.adaptParameterOrRef(p)
		}
		out["parameters"] = params
	}

	if op.RequestBody != nil {
		out["requestBody"] = a.adaptRequestBodyOrRef(op.RequestBody)
	}

	if op.Responses != nil {
		out["responses"] = a.adaptResponses(op.Responses)
	}

	if op.Deprecated {
		out["deprecated"] = true
	}

	if len(op.Security) > 0 {
		out["security"] = op.Security
	}

	if len(op.Servers) > 0 {
		servers := make([]any, len(op.Servers))
		for i, s := range op.Servers {
			servers[i] = a.adaptServer(s)
		}
		out["servers"] = servers
	}

	return out
}

func (a *Adapter) adaptExternalDocs(ed *model.ExternalDocs) map[string]any {
	if ed == nil {
		return nil
	}
	out := map[string]any{
		"url": ed.URL,
	}
	setIfNotEmpty(out, "description", ed.Description)
	return out
}

func (a *Adapter) adaptParameterOrRef(p *model.ParameterOrRef) map[string]any {
	if p == nil {
		return nil
	}
	if p.Ref != "" {
		return map[string]any{"$ref": p.Ref}
	}
	if p.Value == nil {
		return nil
	}

	out := map[string]any{
		"name": p.Value.Name,
		"in":   p.Value.In,
	}
	setIfNotEmpty(out, "description", p.Value.Description)

	if p.Value.Required {
		out["required"] = true
	}
	if p.Value.Deprecated {
		out["deprecated"] = true
	}
	if p.Value.AllowEmptyValue {
		out["allowEmptyValue"] = true
	}
	setIfNotEmpty(out, "style", p.Value.Style)
	if p.Value.Explode != nil {
		out["explode"] = *p.Value.Explode
	}
	if p.Value.AllowReserved {
		out["allowReserved"] = true
	}

	if p.Value.Schema != nil {
		out["schema"] = a.adaptSchemaOrRef(p.Value.Schema)
	}

	if examples := a.adaptExamples(p.Value.Examples); len(examples) > 0 {
		out["examples"] = examples
	}

	return out
}

func (a *Adapter) adaptRequestBodyOrRef(rb *model.RequestBodyOrRef) map[string]any {
	if rb == nil {
		return nil
	}
	if rb.Ref != "" {
		return map[string]any{"$ref": rb.Ref}
	}
	if rb.Value == nil {
		return nil
	}

	out := map[string]any{}
	setIfNotEmpty(out, "description", rb.Value.Description)

	if rb.Value.Required {
		out["required"] = true
	}

	if len(rb.Value.Content) > 0 {
		content := make(map[string]any)
		for mt, media := range rb.Value.Content {
			content[mt] = a.adaptMediaType(media)
		}
		out["content"] = content
	}

	return out
}

func (a *Adapter) adaptMediaType(mt *model.MediaType) map[string]any {
	if mt == nil {
		return nil
	}
	out := map[string]any{}

	if mt.Schema != nil {
		out["schema"] = a.adaptSchemaOrRef(mt.Schema)
	}

	if examples := a.adaptExamples(mt.Examples); len(examples) > 0 {
		out["examples"] = examples
	}

	return out
}

func (a *Adapter) adaptResponses(r *model.Responses) map[string]any {
	if r == nil {
		return nil
	}
	out := make(map[string]any)

	if r.Default != nil {
		out["default"] = a.adaptResponseOrRef(r.Default)
	}

	for code, resp := range r.Codes {
		out[code] = a.adaptResponseOrRef(resp)
	}

	return out
}

func (a *Adapter) adaptResponseOrRef(r *model.ResponseOrRef) map[string]any {
	if r == nil {
		return nil
	}
	if r.Ref != "" {
		return map[string]any{"$ref": r.Ref}
	}
	if r.Value == nil {
		return nil
	}

	out := map[string]any{
		"description": r.Value.Description,
	}

	if len(r.Value.Headers) > 0 {
		headers := make(map[string]any)
		for name, h := range r.Value.Headers {
			headers[name] = a.adaptHeaderOrRef(h)
		}
		out["headers"] = headers
	}

	if len(r.Value.Content) > 0 {
		content := make(map[string]any)
		for mt, media := range r.Value.Content {
			content[mt] = a.adaptMediaType(media)
		}
		out["content"] = content
	}

	return out
}

func (a *Adapter) adaptHeaderOrRef(h *model.HeaderOrRef) map[string]any {
	if h == nil {
		return nil
	}
	if h.Ref != "" {
		return map[string]any{"$ref": h.Ref}
	}
	if h.Value == nil {
		return nil
	}

	out := map[string]any{}
	setIfNotEmpty(out, "description", h.Value.Description)

	if h.Value.Required {
		out["required"] = true
	}
	if h.Value.Deprecated {
		out["deprecated"] = true
	}

	if h.Value.Schema != nil {
		out["schema"] = a.adaptSchemaOrRef(h.Value.Schema)
	}

	if examples := a.adaptExamples(h.Value.Examples); len(examples) > 0 {
		out["examples"] = examples
	}

	return out
}

func (a *Adapter) adaptSchemaOrRef(s *model.SchemaOrRef) map[string]any {
	if s == nil {
		return nil
	}

	// Handle reference - 3.1.0 allows $ref with summary/description siblings
	if s.Ref != "" {
		out := map[string]any{"$ref": s.Ref}
		setIfNotEmpty(out, "summary", s.Summary)
		setIfNotEmpty(out, "description", s.Description)
		return out
	}

	if s.Value == nil {
		return nil
	}

	return a.adaptSchema(s.Value)
}

func (a *Adapter) adaptSchema(schema *model.Schema) map[string]any {
	if schema == nil {
		return nil
	}

	out := map[string]any{}

	// 3.1.0: Nullable via type array ["type", "null"], otherwise single string
	if schema.Type != "" {
		if schema.IsNullable {
			out["type"] = []string{schema.Type, "null"}
		} else {
			out["type"] = schema.Type
		}
	}

	setIfNotEmpty(out, "format", schema.Format)
	setIfNotEmpty(out, "title", schema.Title)
	setIfNotEmpty(out, "description", schema.Description)

	if schema.Default != nil {
		out["default"] = schema.Default
	}

	// 3.1.0: examples as map
	if examples := a.adaptExamples(schema.Examples); len(examples) > 0 {
		out["examples"] = examples
	}

	if schema.Deprecated {
		out["deprecated"] = true
	}
	if schema.ReadOnly {
		out["readOnly"] = true
	}
	if schema.WriteOnly {
		out["writeOnly"] = true
	}

	if schema.ExternalDocs != nil {
		out["externalDocs"] = a.adaptExternalDocs(schema.ExternalDocs)
	}

	// Numeric validation
	if schema.MultipleOf != nil {
		out["multipleOf"] = *schema.MultipleOf
	}
	if schema.Minimum != nil {
		out["minimum"] = *schema.Minimum
	}
	if schema.Maximum != nil {
		out["maximum"] = *schema.Maximum
	}
	// 3.1.0: exclusiveMinimum/Maximum are numeric, not boolean
	if schema.ExclusiveMinimum != nil {
		out["exclusiveMinimum"] = *schema.ExclusiveMinimum
	}
	if schema.ExclusiveMaximum != nil {
		out["exclusiveMaximum"] = *schema.ExclusiveMaximum
	}

	// String validation
	if schema.MinLength != nil {
		out["minLength"] = *schema.MinLength
	}
	if schema.MaxLength != nil {
		out["maxLength"] = *schema.MaxLength
	}
	setIfNotEmpty(out, "pattern", schema.Pattern)

	// Array validation
	if schema.MinItems != nil {
		out["minItems"] = *schema.MinItems
	}
	if schema.MaxItems != nil {
		out["maxItems"] = *schema.MaxItems
	}
	if schema.UniqueItems {
		out["uniqueItems"] = true
	}
	if schema.Items != nil {
		out["items"] = a.adaptSchemaOrRef(schema.Items)
	}

	// Object validation
	if schema.MinProperties != nil {
		out["minProperties"] = *schema.MinProperties
	}
	if schema.MaxProperties != nil {
		out["maxProperties"] = *schema.MaxProperties
	}
	if len(schema.Required) > 0 {
		out["required"] = schema.Required
	}

	if len(schema.Properties) > 0 {
		props := make(map[string]any)
		for name, prop := range schema.Properties {
			props[name] = a.adaptSchemaOrRef(prop)
		}
		out["properties"] = props
	}

	if schema.AdditionalProperties != nil {
		if schema.AdditionalProperties.Schema != nil {
			out["additionalProperties"] = a.adaptSchemaOrRef(schema.AdditionalProperties.Schema)
		} else if schema.AdditionalProperties.Allowed {
			out["additionalProperties"] = true
		} else {
			out["additionalProperties"] = false
		}
	}

	// Composition
	if len(schema.AllOf) > 0 {
		allOf := make([]any, len(schema.AllOf))
		for i, s := range schema.AllOf {
			allOf[i] = a.adaptSchemaOrRef(s)
		}
		out["allOf"] = allOf
	}
	if len(schema.AnyOf) > 0 {
		anyOf := make([]any, len(schema.AnyOf))
		for i, s := range schema.AnyOf {
			anyOf[i] = a.adaptSchemaOrRef(s)
		}
		out["anyOf"] = anyOf
	}
	if len(schema.OneOf) > 0 {
		oneOf := make([]any, len(schema.OneOf))
		for i, s := range schema.OneOf {
			oneOf[i] = a.adaptSchemaOrRef(s)
		}
		out["oneOf"] = oneOf
	}
	if schema.Not != nil {
		out["not"] = a.adaptSchemaOrRef(schema.Not)
	}

	if schema.Discriminator != nil {
		disc := map[string]any{
			"propertyName": schema.Discriminator.PropertyName,
		}
		if len(schema.Discriminator.Mapping) > 0 {
			disc["mapping"] = schema.Discriminator.Mapping
		}
		out["discriminator"] = disc
	}

	if len(schema.Enum) > 0 {
		out["enum"] = schema.Enum
	}

	return out
}

func (a *Adapter) adaptExamples(examples []*model.Example) map[string]any {
	if len(examples) == 0 {
		return nil
	}
	out := make(map[string]any)
	for i, ex := range examples {
		name := ex.Name
		if name == "" {
			if i == 0 {
				name = "example"
			} else {
				name = "example_" + string(rune('0'+i))
			}
		}
		exOut := map[string]any{}
		setIfNotEmpty(exOut, "summary", ex.Summary)
		setIfNotEmpty(exOut, "description", ex.Description)
		if ex.Value != nil {
			exOut["value"] = ex.Value
		}
		setIfNotEmpty(exOut, "externalValue", ex.ExternalValue)
		out[name] = exOut
	}
	return out
}

func (a *Adapter) adaptComponents(c *model.Components) map[string]any {
	if c == nil {
		return nil
	}
	out := map[string]any{}

	if len(c.Schemas) > 0 {
		schemas := make(map[string]any)
		for name, s := range c.Schemas {
			schemas[name] = a.adaptSchemaOrRef(s)
		}
		out["schemas"] = schemas
	}

	if len(c.Responses) > 0 {
		responses := make(map[string]any)
		for name, r := range c.Responses {
			responses[name] = a.adaptResponseOrRef(r)
		}
		out["responses"] = responses
	}

	if len(c.Parameters) > 0 {
		params := make(map[string]any)
		for name, p := range c.Parameters {
			params[name] = a.adaptParameterOrRef(p)
		}
		out["parameters"] = params
	}

	if len(c.Examples) > 0 {
		examples := make(map[string]any)
		for name, ex := range c.Examples {
			if ex.Ref != "" {
				examples[name] = map[string]any{"$ref": ex.Ref}
			} else if ex.Value != nil {
				exOut := map[string]any{}
				setIfNotEmpty(exOut, "summary", ex.Value.Summary)
				setIfNotEmpty(exOut, "description", ex.Value.Description)
				if ex.Value.Value != nil {
					exOut["value"] = ex.Value.Value
				}
				setIfNotEmpty(exOut, "externalValue", ex.Value.ExternalValue)
				examples[name] = exOut
			}
		}
		out["examples"] = examples
	}

	if len(c.RequestBodies) > 0 {
		bodies := make(map[string]any)
		for name, rb := range c.RequestBodies {
			bodies[name] = a.adaptRequestBodyOrRef(rb)
		}
		out["requestBodies"] = bodies
	}

	if len(c.Headers) > 0 {
		headers := make(map[string]any)
		for name, h := range c.Headers {
			headers[name] = a.adaptHeaderOrRef(h)
		}
		out["headers"] = headers
	}

	if len(c.SecuritySchemes) > 0 {
		schemes := make(map[string]any)
		for name, ss := range c.SecuritySchemes {
			schemes[name] = a.adaptSecurityScheme(ss)
		}
		out["securitySchemes"] = schemes
	}

	if len(c.PathItems) > 0 {
		items := make(map[string]any)
		for name, pi := range c.PathItems {
			items[name] = a.adaptPathItem(pi)
		}
		out["pathItems"] = items
	}

	return out
}

func (a *Adapter) adaptSecurityScheme(ss *model.SecuritySchemeOrRef) map[string]any {
	if ss == nil {
		return nil
	}
	if ss.Ref != "" {
		return map[string]any{"$ref": ss.Ref}
	}
	if ss.Value == nil {
		return nil
	}

	out := map[string]any{}
	setIfNotEmpty(out, "type", ss.Value.Type)
	setIfNotEmpty(out, "description", ss.Value.Description)
	setIfNotEmpty(out, "name", ss.Value.Name)
	setIfNotEmpty(out, "in", ss.Value.In)
	setIfNotEmpty(out, "scheme", ss.Value.Scheme)
	setIfNotEmpty(out, "bearerFormat", ss.Value.BearerFormat)
	setIfNotEmpty(out, "openIdConnectUrl", ss.Value.OpenIDConnectURL)

	if ss.Value.Flows != nil {
		out["flows"] = a.adaptOAuthFlows(ss.Value.Flows)
	}

	return out
}

func (a *Adapter) adaptOAuthFlows(flows *model.OAuthFlows) map[string]any {
	if flows == nil {
		return nil
	}
	out := map[string]any{}

	if flows.Implicit != nil {
		out["implicit"] = a.adaptOAuthFlow(flows.Implicit)
	}
	if flows.Password != nil {
		out["password"] = a.adaptOAuthFlow(flows.Password)
	}
	if flows.ClientCredentials != nil {
		out["clientCredentials"] = a.adaptOAuthFlow(flows.ClientCredentials)
	}
	if flows.AuthorizationCode != nil {
		out["authorizationCode"] = a.adaptOAuthFlow(flows.AuthorizationCode)
	}

	return out
}

func (a *Adapter) adaptOAuthFlow(flow *model.OAuthFlow) map[string]any {
	if flow == nil {
		return nil
	}
	out := map[string]any{}
	setIfNotEmpty(out, "authorizationUrl", flow.AuthorizationURL)
	setIfNotEmpty(out, "tokenUrl", flow.TokenURL)
	setIfNotEmpty(out, "refreshUrl", flow.RefreshURL)

	if len(flow.Scopes) > 0 {
		out["scopes"] = flow.Scopes
	}

	return out
}

func (a *Adapter) adaptTag(t *model.Tag) map[string]any {
	if t == nil {
		return nil
	}
	out := map[string]any{
		"name": t.Name,
	}
	setIfNotEmpty(out, "description", t.Description)

	if t.ExternalDocs != nil {
		out["externalDocs"] = a.adaptExternalDocs(t.ExternalDocs)
	}

	return out
}

// setIfNotEmpty sets a key in the map only if the value is not empty.
func setIfNotEmpty(m map[string]any, key, value string) {
	if value != "" {
		m[key] = value
	}
}
