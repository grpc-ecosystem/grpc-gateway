package genopenapiv3

import (
	"bytes"
	"encoding/json"
	"sort"
)

// OpenAPI is the root document structure for OpenAPI 3.x
// See: https://spec.openapis.org/oas/v3.0.3#openapi-object
type OpenAPI struct {
	OpenAPI      string                  `json:"openapi" yaml:"openapi"`                           // REQUIRED: "3.0.3" or "3.1.0"
	Info         *Info                   `json:"info" yaml:"info"`                                 // REQUIRED
	Servers      []*Server               `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        *Paths                  `json:"paths,omitempty" yaml:"paths,omitempty"`
	Components   *Components             `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []SecurityRequirement   `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []*Tag                  `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation  `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// NewOpenAPI creates a new OpenAPI v3 document with required fields.
func NewOpenAPI(title, version, openapiVersion string) *OpenAPI {
	return &OpenAPI{
		OpenAPI: openapiVersion,
		Info: &Info{
			Title:   title,
			Version: version,
		},
		Paths: NewPaths(),
		Components: &Components{
			Schemas: make(map[string]*SchemaRef),
		},
	}
}

// Info provides metadata about the API.
// See: https://spec.openapis.org/oas/v3.0.3#info-object
type Info struct {
	Title          string   `json:"title" yaml:"title"`                                     // REQUIRED
	Summary        string   `json:"summary,omitempty" yaml:"summary,omitempty"`             // v3.1 only
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`                                 // REQUIRED

	extensions []extension `json:"-" yaml:"-"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	Name       string `json:"name" yaml:"name"` // REQUIRED
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"` // v3.1 only, SPDX identifier
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Server represents a server.
// See: https://spec.openapis.org/oas/v3.0.3#server-object
type Server struct {
	URL         string                     `json:"url" yaml:"url"` // REQUIRED
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable represents a server variable for server URL template substitution.
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"` // REQUIRED
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// Paths holds the relative paths to the individual endpoints and their operations.
// Maintains insertion order for consistent output.
type Paths struct {
	paths map[string]*PathItem
	order []string
}

// NewPaths creates a new Paths object.
func NewPaths() *Paths {
	return &Paths{
		paths: make(map[string]*PathItem),
		order: []string{},
	}
}

// Set adds or updates a path item.
func (p *Paths) Set(path string, item *PathItem) {
	if _, exists := p.paths[path]; !exists {
		p.order = append(p.order, path)
	}
	p.paths[path] = item
}

// Get retrieves a path item by path.
func (p *Paths) Get(path string) *PathItem {
	if p == nil || p.paths == nil {
		return nil
	}
	return p.paths[path]
}

// Len returns the number of paths.
func (p *Paths) Len() int {
	if p == nil {
		return 0
	}
	return len(p.paths)
}

// SortAlphabetically sorts paths alphabetically.
func (p *Paths) SortAlphabetically() {
	if p == nil {
		return
	}
	sort.Strings(p.order)
}

// MarshalJSON outputs paths in insertion order.
func (p *Paths) MarshalJSON() ([]byte, error) {
	if p == nil || len(p.paths) == 0 {
		return []byte("{}"), nil
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	first := true
	for _, path := range p.order {
		item := p.paths[path]
		if item == nil {
			continue
		}
		if !first {
			buf.WriteString(",")
		}
		first = false

		key, err := json.Marshal(path)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")

		val, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// MarshalYAML implements yaml.Marshaler for Paths.
func (p *Paths) MarshalYAML() (any, error) {
	if p == nil || len(p.paths) == 0 {
		return map[string]any{}, nil
	}

	ordered := make(map[string]*PathItem)
	for _, path := range p.order {
		ordered[path] = p.paths[path]
	}
	return ordered, nil
}

// PathItem describes the operations available on a single path.
// See: https://spec.openapis.org/oas/v3.0.3#path-item-object
type PathItem struct {
	Ref         string          `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string          `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string          `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation      `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation      `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation      `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation      `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation      `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation      `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation      `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation      `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []*Server       `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []*ParameterRef `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// SetOperation sets the operation for the given HTTP method.
func (p *PathItem) SetOperation(method string, op *Operation) {
	switch method {
	case "GET":
		p.Get = op
	case "POST":
		p.Post = op
	case "PUT":
		p.Put = op
	case "PATCH":
		p.Patch = op
	case "DELETE":
		p.Delete = op
	case "HEAD":
		p.Head = op
	case "OPTIONS":
		p.Options = op
	case "TRACE":
		p.Trace = op
	}
}

// Operation describes a single API operation on a path.
// See: https://spec.openapis.org/oas/v3.0.3#operation-object
type Operation struct {
	Tags         []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   []*ParameterRef        `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *RequestBodyRef        `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    *Responses             `json:"responses,omitempty" yaml:"responses,omitempty"` // REQUIRED
	Callbacks    map[string]*CallbackRef `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []SecurityRequirement  `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []*Server              `json:"servers,omitempty" yaml:"servers,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// Parameter describes a single operation parameter.
// See: https://spec.openapis.org/oas/v3.0.3#parameter-object
type Parameter struct {
	Name            string      `json:"name" yaml:"name"`       // REQUIRED
	In              string      `json:"in" yaml:"in"`           // REQUIRED: query, header, path, cookie
	Description     string      `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool        `json:"required,omitempty" yaml:"required,omitempty"` // REQUIRED if in="path"
	Deprecated      bool        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string      `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool       `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *SchemaRef  `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         any         `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType  `json:"content,omitempty" yaml:"content,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// NewPathParameter creates a path parameter (always required).
func NewPathParameter(name string, schema *SchemaRef) *Parameter {
	return &Parameter{
		Name:     name,
		In:       "path",
		Required: true,
		Schema:   schema,
	}
}

// NewQueryParameter creates a query parameter.
func NewQueryParameter(name string, schema *SchemaRef) *Parameter {
	return &Parameter{
		Name:   name,
		In:     "query",
		Schema: schema,
	}
}

// NewHeaderParameter creates a header parameter.
func NewHeaderParameter(name string, schema *SchemaRef) *Parameter {
	return &Parameter{
		Name:   name,
		In:     "header",
		Schema: schema,
	}
}

// ParameterRef can be a reference or inline parameter.
type ParameterRef struct {
	Ref   string     `json:"-" yaml:"-"`
	Value *Parameter `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (p *ParameterRef) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}
	if p.Ref != "" {
		return json.Marshal(map[string]string{"$ref": p.Ref})
	}
	return json.Marshal(p.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (p *ParameterRef) MarshalYAML() (any, error) {
	if p == nil {
		return nil, nil
	}
	if p.Ref != "" {
		return map[string]string{"$ref": p.Ref}, nil
	}
	return p.Value, nil
}

// RequestBody describes a single request body.
// See: https://spec.openapis.org/oas/v3.0.3#request-body-object
type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content" yaml:"content"` // REQUIRED
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// NewJSONRequestBody creates a request body with JSON content.
func NewJSONRequestBody(schema *SchemaRef, required bool) *RequestBody {
	return &RequestBody{
		Content: map[string]*MediaType{
			"application/json": {Schema: schema},
		},
		Required: required,
	}
}

// RequestBodyRef can be a reference or inline request body.
type RequestBodyRef struct {
	Ref   string       `json:"-" yaml:"-"`
	Value *RequestBody `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (r *RequestBodyRef) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	if r.Ref != "" {
		return json.Marshal(map[string]string{"$ref": r.Ref})
	}
	return json.Marshal(r.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (r *RequestBodyRef) MarshalYAML() (any, error) {
	if r == nil {
		return nil, nil
	}
	if r.Ref != "" {
		return map[string]string{"$ref": r.Ref}, nil
	}
	return r.Value, nil
}

// MediaType provides schema and examples for a media type.
// See: https://spec.openapis.org/oas/v3.0.3#media-type-object
type MediaType struct {
	Schema   *SchemaRef             `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  any                    `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding   `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding defines how to serialize a parameter.
type Encoding struct {
	ContentType   string                `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]*HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool                 `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                  `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// Responses contains the expected responses for an operation.
// See: https://spec.openapis.org/oas/v3.0.3#responses-object
type Responses struct {
	Default *ResponseRef            `json:"-" yaml:"-"`
	Codes   map[string]*ResponseRef `json:"-" yaml:"-"`
}

// NewResponses creates a new Responses object.
func NewResponses() *Responses {
	return &Responses{
		Codes: make(map[string]*ResponseRef),
	}
}

// MarshalJSON implements json.Marshaler for Responses.
func (r *Responses) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("{}"), nil
	}

	m := make(map[string]any)
	if r.Default != nil {
		m["default"] = r.Default
	}
	for code, resp := range r.Codes {
		m[code] = resp
	}
	return json.Marshal(m)
}

// MarshalYAML implements yaml.Marshaler for Responses.
func (r *Responses) MarshalYAML() (any, error) {
	if r == nil {
		return map[string]any{}, nil
	}

	m := make(map[string]any)
	if r.Default != nil {
		m["default"] = r.Default
	}
	for code, resp := range r.Codes {
		m[code] = resp
	}
	return m, nil
}

// Response describes a single response from an API operation.
// See: https://spec.openapis.org/oas/v3.0.3#response-object
// IMPORTANT: description is REQUIRED per OpenAPI spec!
type Response struct {
	Description string                 `json:"description" yaml:"description"` // REQUIRED - no omitempty!
	Headers     map[string]*HeaderRef  `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType  `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*LinkRef    `json:"links,omitempty" yaml:"links,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// NewResponse creates a response with required description.
func NewResponse(description string) *Response {
	return &Response{Description: description}
}

// WithJSONSchema adds JSON content with the given schema.
func (r *Response) WithJSONSchema(schema *SchemaRef) *Response {
	if r.Content == nil {
		r.Content = make(map[string]*MediaType)
	}
	r.Content["application/json"] = &MediaType{Schema: schema}
	return r
}

// ResponseRef can be a reference or inline response.
type ResponseRef struct {
	Ref   string    `json:"-" yaml:"-"`
	Value *Response `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (r *ResponseRef) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	if r.Ref != "" {
		return json.Marshal(map[string]string{"$ref": r.Ref})
	}
	return json.Marshal(r.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (r *ResponseRef) MarshalYAML() (any, error) {
	if r == nil {
		return nil, nil
	}
	if r.Ref != "" {
		return map[string]string{"$ref": r.Ref}, nil
	}
	return r.Value, nil
}

// Schema represents a JSON Schema object.
// See: https://spec.openapis.org/oas/v3.0.3#schema-object
type Schema struct {
	// Type information
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`     // object, array, string, number, integer, boolean
	Format string `json:"format,omitempty" yaml:"format,omitempty"` // int32, int64, float, double, byte, date, date-time, etc.

	// Documentation
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Validation - general
	Default any   `json:"default,omitempty" yaml:"default,omitempty"`
	Enum    []any `json:"enum,omitempty" yaml:"enum,omitempty"`

	// Validation - numeric
	Minimum          *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // boolean in v3.0, number in v3.1
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`

	// Validation - string
	MinLength *uint64 `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength *uint64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Object properties
	Properties           map[string]*SchemaRef `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *SchemaRef            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Required             []string              `json:"required,omitempty" yaml:"required,omitempty"`
	MinProperties        *uint64               `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	MaxProperties        *uint64               `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`

	// Array items
	Items       *SchemaRef `json:"items,omitempty" yaml:"items,omitempty"`
	MinItems    *uint64    `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems    *uint64    `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	UniqueItems bool       `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`

	// Composition (v3 feature!)
	OneOf []*SchemaRef `json:"oneOf,omitempty" yaml:"oneOf,omitempty"` // Exactly one must match
	AnyOf []*SchemaRef `json:"anyOf,omitempty" yaml:"anyOf,omitempty"` // At least one must match
	AllOf []*SchemaRef `json:"allOf,omitempty" yaml:"allOf,omitempty"` // All must match
	Not   *SchemaRef   `json:"not,omitempty" yaml:"not,omitempty"`     // Must not match

	// Access control
	ReadOnly  bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly bool `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	Nullable  bool `json:"nullable,omitempty" yaml:"nullable,omitempty"` // v3.0 feature, removed in v3.1 (use type: ["string", "null"])

	// Deprecation
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Example
	Example any `json:"example,omitempty" yaml:"example,omitempty"`

	// External docs
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// SchemaRef can be either a reference or an inline schema.
// This pattern allows: {"$ref": "#/..."} OR {"type": "string", ...}
type SchemaRef struct {
	Ref   string  `json:"-" yaml:"-"`
	Value *Schema `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (s *SchemaRef) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	if s.Ref != "" {
		return json.Marshal(map[string]string{"$ref": s.Ref})
	}
	return json.Marshal(s.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (s *SchemaRef) MarshalYAML() (any, error) {
	if s == nil {
		return nil, nil
	}
	if s.Ref != "" {
		return map[string]string{"$ref": s.Ref}, nil
	}
	return s.Value, nil
}

// NewSchemaRef creates a reference to a schema in components.
func NewSchemaRef(name string) *SchemaRef {
	return &SchemaRef{
		Ref: "#/components/schemas/" + name,
	}
}

// NewInlineSchema creates an inline schema (not a reference).
func NewInlineSchema(schema *Schema) *SchemaRef {
	return &SchemaRef{Value: schema}
}

// Components holds reusable objects.
// See: https://spec.openapis.org/oas/v3.0.3#components-object
type Components struct {
	Schemas         map[string]*SchemaRef         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*ResponseRef       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*ParameterRef      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]*ExampleRef        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBodyRef    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*HeaderRef         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeRef `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*LinkRef           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*CallbackRef       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

// SecurityScheme defines a security scheme for operations.
// See: https://spec.openapis.org/oas/v3.0.3#security-scheme-object
type SecurityScheme struct {
	Type             string      `json:"type" yaml:"type"`                                         // REQUIRED: apiKey, http, oauth2, openIdConnect
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty" yaml:"name,omitempty"`                     // REQUIRED for apiKey
	In               string      `json:"in,omitempty" yaml:"in,omitempty"`                         // REQUIRED for apiKey: query, header, cookie
	Scheme           string      `json:"scheme,omitempty" yaml:"scheme,omitempty"`                 // REQUIRED for http: basic, bearer, etc.
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`                   // REQUIRED for oauth2
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"` // REQUIRED for openIdConnect

	extensions []extension `json:"-" yaml:"-"`
}

// SecuritySchemeRef can be a reference or inline security scheme.
type SecuritySchemeRef struct {
	Ref   string          `json:"-" yaml:"-"`
	Value *SecurityScheme `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (s *SecuritySchemeRef) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	if s.Ref != "" {
		return json.Marshal(map[string]string{"$ref": s.Ref})
	}
	return json.Marshal(s.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (s *SecuritySchemeRef) MarshalYAML() (any, error) {
	if s == nil {
		return nil, nil
	}
	if s.Ref != "" {
		return map[string]string{"$ref": s.Ref}, nil
	}
	return s.Value, nil
}

// OAuthFlows allows configuration of OAuth flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// OAuthFlow configuration details for a specific OAuth flow.
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"` // REQUIRED

	extensions []extension `json:"-" yaml:"-"`
}

// SecurityRequirement maps scheme names to required scopes.
type SecurityRequirement map[string][]string

// Tag adds metadata to a single tag.
// See: https://spec.openapis.org/oas/v3.0.3#tag-object
type Tag struct {
	Name         string                 `json:"name" yaml:"name"` // REQUIRED
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// ExternalDocumentation allows referencing external resource for extended documentation.
type ExternalDocumentation struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"` // REQUIRED
}

// Header describes a single header.
// See: https://spec.openapis.org/oas/v3.0.3#header-object
type Header struct {
	Description     string     `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool       `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool       `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool       `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string     `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool      `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool       `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *SchemaRef `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         any        `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType  `json:"content,omitempty" yaml:"content,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// HeaderRef can be a reference or inline header.
type HeaderRef struct {
	Ref   string  `json:"-" yaml:"-"`
	Value *Header `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (h *HeaderRef) MarshalJSON() ([]byte, error) {
	if h == nil {
		return []byte("null"), nil
	}
	if h.Ref != "" {
		return json.Marshal(map[string]string{"$ref": h.Ref})
	}
	return json.Marshal(h.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (h *HeaderRef) MarshalYAML() (any, error) {
	if h == nil {
		return nil, nil
	}
	if h.Ref != "" {
		return map[string]string{"$ref": h.Ref}, nil
	}
	return h.Value, nil
}

// Example object.
type Example struct {
	Summary       string `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string `json:"description,omitempty" yaml:"description,omitempty"`
	Value         any    `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// ExampleRef can be a reference or inline example.
type ExampleRef struct {
	Ref   string   `json:"-" yaml:"-"`
	Value *Example `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExampleRef) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	if e.Ref != "" {
		return json.Marshal(map[string]string{"$ref": e.Ref})
	}
	return json.Marshal(e.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExampleRef) MarshalYAML() (any, error) {
	if e == nil {
		return nil, nil
	}
	if e.Ref != "" {
		return map[string]string{"$ref": e.Ref}, nil
	}
	return e.Value, nil
}

// Link represents a possible design-time link for a response.
type Link struct {
	OperationRef string            `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationId  string            `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]any    `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  any               `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *Server           `json:"server,omitempty" yaml:"server,omitempty"`

	extensions []extension `json:"-" yaml:"-"`
}

// LinkRef can be a reference or inline link.
type LinkRef struct {
	Ref   string `json:"-" yaml:"-"`
	Value *Link  `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (l *LinkRef) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("null"), nil
	}
	if l.Ref != "" {
		return json.Marshal(map[string]string{"$ref": l.Ref})
	}
	return json.Marshal(l.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (l *LinkRef) MarshalYAML() (any, error) {
	if l == nil {
		return nil, nil
	}
	if l.Ref != "" {
		return map[string]string{"$ref": l.Ref}, nil
	}
	return l.Value, nil
}

// Callback is a map of runtime expressions to PathItems.
type Callback map[string]*PathItem

// CallbackRef can be a reference or inline callback.
type CallbackRef struct {
	Ref   string    `json:"-" yaml:"-"`
	Value *Callback `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler.
func (c *CallbackRef) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	if c.Ref != "" {
		return json.Marshal(map[string]string{"$ref": c.Ref})
	}
	return json.Marshal(c.Value)
}

// MarshalYAML implements yaml.Marshaler.
func (c *CallbackRef) MarshalYAML() (any, error) {
	if c == nil {
		return nil, nil
	}
	if c.Ref != "" {
		return map[string]string{"$ref": c.Ref}, nil
	}
	return c.Value, nil
}

// Discriminator used when request bodies or response payloads may be one of several types.
type Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"` // REQUIRED
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

// extension holds a custom extension key-value pair.
type extension struct {
	key   string
	value json.RawMessage
}
