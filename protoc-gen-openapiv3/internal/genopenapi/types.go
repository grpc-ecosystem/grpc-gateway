package genopenapi

import (
	"bytes"
	"encoding/json"
	"maps"
	"slices"
	"sort"
)

// Document is the root OpenAPI 3.1.0 object.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#openapi-object
type Document struct {
	OpenAPI      string                `json:"openapi"`
	Info         *Info                 `json:"info"`
	Servers      []*Server             `json:"servers,omitempty"`
	Paths        *Paths                `json:"paths,omitempty"`
	Components   *Components           `json:"components,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Tags         []*Tag                `json:"tags,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty"`
}

// NewDocument returns a Document populated with the required fields.
func NewDocument(title, version string) *Document {
	return &Document{
		OpenAPI: "3.1.0",
		Info: &Info{
			Title:   title,
			Version: version,
		},
		Paths: NewPaths(),
		Components: &Components{
			Schemas: make(map[string]*SchemaOrRef),
		},
	}
}

// Info is metadata about the API.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#info-object
type Info struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
}

// Contact is the API contact info.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#contact-object
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License is the API license info. In 3.1.0 identifier and url are mutually
// exclusive; this generator does not enforce that — set whichever is correct.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#license-object
type License struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier,omitempty"`
	URL        string `json:"url,omitempty"`
}

// Server represents an API server.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#server-object
type Server struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

// ServerVariable is a server URL template variable.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#server-variable-object
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// Paths is an insertion-ordered map from URL template to PathItem. The
// ordering is exposed via a custom JSON marshaler so the generated output
// preserves RPC declaration order.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#paths-object
type Paths struct {
	items map[string]*PathItem
	order []string
}

// NewPaths constructs an empty Paths.
func NewPaths() *Paths {
	return &Paths{items: make(map[string]*PathItem)}
}

// Set inserts or replaces a PathItem at the given URL template.
func (p *Paths) Set(path string, item *PathItem) {
	if _, ok := p.items[path]; !ok {
		p.order = append(p.order, path)
	}
	p.items[path] = item
}

// Get retrieves a PathItem. The boolean is false if no item exists at path.
func (p *Paths) Get(path string) (*PathItem, bool) {
	if p == nil {
		return nil, false
	}
	item, ok := p.items[path]
	return item, ok
}

// Len reports how many paths are stored.
func (p *Paths) Len() int {
	if p == nil {
		return 0
	}
	return len(p.items)
}

// MarshalJSON serializes paths in insertion order.
func (p *Paths) MarshalJSON() ([]byte, error) {
	if p == nil || len(p.items) == 0 {
		return []byte("{}"), nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, path := range p.order {
		if i > 0 {
			buf.WriteByte(',')
		}
		key, err := json.Marshal(path)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteByte(':')
		val, err := json.Marshal(p.items[path])
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// PathItem describes the operations available on a single URL template.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#path-item-object
type PathItem struct {
	Summary     string          `json:"summary,omitempty"`
	Description string          `json:"description,omitempty"`
	Get         *Operation      `json:"get,omitempty"`
	Put         *Operation      `json:"put,omitempty"`
	Post        *Operation      `json:"post,omitempty"`
	Delete      *Operation      `json:"delete,omitempty"`
	Options     *Operation      `json:"options,omitempty"`
	Head        *Operation      `json:"head,omitempty"`
	Patch       *Operation      `json:"patch,omitempty"`
	Trace       *Operation      `json:"trace,omitempty"`
	Servers     []*Server       `json:"servers,omitempty"`
	Parameters  []*ParameterRef `json:"parameters,omitempty"`
}

// SetOperation places op under the slot for the given HTTP method.
func (p *PathItem) SetOperation(method string, op *Operation) {
	switch method {
	case "GET":
		p.Get = op
	case "PUT":
		p.Put = op
	case "POST":
		p.Post = op
	case "DELETE":
		p.Delete = op
	case "OPTIONS":
		p.Options = op
	case "HEAD":
		p.Head = op
	case "PATCH":
		p.Patch = op
	case "TRACE":
		p.Trace = op
	}
}

// Operation describes a single API operation on a path.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#operation-object
type Operation struct {
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty"`
	Parameters   []*ParameterRef       `json:"parameters,omitempty"`
	RequestBody  *RequestBodyRef       `json:"requestBody,omitempty"`
	Responses    *Responses            `json:"responses,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Servers      []*Server             `json:"servers,omitempty"`
}

// Parameter describes one operation parameter.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#parameter-object
type Parameter struct {
	Name        string       `json:"name"`
	In          string       `json:"in"`
	Description string       `json:"description,omitempty"`
	Required    bool         `json:"required,omitempty"`
	Deprecated  bool         `json:"deprecated,omitempty"`
	Schema      *SchemaOrRef `json:"schema,omitempty"`
	Example     any          `json:"example,omitempty"`
}

// ParameterRef is either an inline Parameter or a $ref to a component
// parameter.
//
// Spec (reference form): https://spec.openapis.org/oas/v3.1.0#reference-object
type ParameterRef struct {
	Ref   string
	Value *Parameter
}

// MarshalJSON dispatches to either {"$ref": ...} or the inline value.
func (p *ParameterRef) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}
	if p.Ref != "" {
		return json.Marshal(map[string]string{"$ref": p.Ref})
	}
	return json.Marshal(p.Value)
}

// RequestBody describes a single request body.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#request-body-object
type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Content     map[string]*MediaType `json:"content"`
	Required    bool                  `json:"required,omitempty"`
}

// RequestBodyRef is either an inline body or a reference.
//
// Spec (reference form): https://spec.openapis.org/oas/v3.1.0#reference-object
type RequestBodyRef struct {
	Ref   string
	Value *RequestBody
}

// MarshalJSON dispatches between $ref and inline value.
func (r *RequestBodyRef) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	if r.Ref != "" {
		return json.Marshal(map[string]string{"$ref": r.Ref})
	}
	return json.Marshal(r.Value)
}

// MediaType describes one media-type entry under content.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaType struct {
	Schema  *SchemaOrRef `json:"schema,omitempty"`
	Example any          `json:"example,omitempty"`
}

// Responses is the responses map. It is encoded with "default" first, then
// status codes in lexicographic order.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#responses-object
type Responses struct {
	Default *ResponseRef
	Codes   map[string]*ResponseRef
}

// NewResponses constructs an empty Responses container.
func NewResponses() *Responses {
	return &Responses{Codes: make(map[string]*ResponseRef)}
}

// MarshalJSON serializes responses with default first, then sorted codes.
func (r *Responses) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("{}"), nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	if r.Default != nil {
		buf.WriteString(`"default":`)
		data, err := json.Marshal(r.Default)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		first = false
	}
	codes := make([]string, 0, len(r.Codes))
	for c := range r.Codes {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	for _, c := range codes {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		key, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteByte(':')
		data, err := json.Marshal(r.Codes[c])
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// Response describes one response from an operation.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#response-object
type Response struct {
	Description string                `json:"description"`
	Headers     map[string]*Header    `json:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

// NewResponse constructs a Response with the required description set.
func NewResponse(description string) *Response {
	return &Response{Description: description}
}

// WithJSONSchema attaches a JSON content media type with the given schema.
func (r *Response) WithJSONSchema(schema *SchemaOrRef) *Response {
	if r.Content == nil {
		r.Content = make(map[string]*MediaType)
	}
	r.Content["application/json"] = &MediaType{Schema: schema}
	return r
}

// ResponseRef is either an inline Response or a reference.
//
// Spec (reference form): https://spec.openapis.org/oas/v3.1.0#reference-object
type ResponseRef struct {
	Ref   string
	Value *Response
}

// MarshalJSON dispatches between $ref and inline value.
func (r *ResponseRef) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	if r.Ref != "" {
		return json.Marshal(map[string]string{"$ref": r.Ref})
	}
	return json.Marshal(r.Value)
}

// Header is a response header description (a Parameter without name/in).
//
// Spec: https://spec.openapis.org/oas/v3.1.0#header-object
type Header struct {
	Description string       `json:"description,omitempty"`
	Schema      *SchemaOrRef `json:"schema,omitempty"`
}

// SchemaType is JSON Schema's type field. In 3.1.0 it can be a string ("foo")
// or an array of strings (["foo","null"]) for nullable types. We use a slice
// internally; MarshalJSON collapses single-element slices to a string.
//
// Spec: https://json-schema.org/draft/2020-12/json-schema-validation#name-type
type SchemaType []string

// MarshalJSON outputs a single string when len==1, an array otherwise.
func (t SchemaType) MarshalJSON() ([]byte, error) {
	if len(t) == 0 {
		return []byte("null"), nil
	}
	if len(t) == 1 {
		return json.Marshal(t[0])
	}
	return json.Marshal([]string(t))
}

// Schema is a JSON Schema 2020-12 object as used by OpenAPI 3.1.0.
//
// Only the fields the generator actually emits today are defined.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#schema-object
// Underlying dialect: https://json-schema.org/draft/2020-12/json-schema-core
type Schema struct {
	Type   SchemaType `json:"type,omitempty"`
	Format string     `json:"format,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Default  any   `json:"default,omitempty"`
	Examples []any `json:"examples,omitempty"`
	Enum     []any `json:"enum,omitempty"`

	// Object validation
	Properties           map[string]*SchemaOrRef `json:"properties,omitempty"`
	Required             []string                `json:"required,omitempty"`
	AdditionalProperties *AdditionalProperties   `json:"additionalProperties,omitempty"`

	// Array validation
	Items *SchemaOrRef `json:"items,omitempty"`

	// Composition
	AllOf []*SchemaOrRef `json:"allOf,omitempty"`
	OneOf []*SchemaOrRef `json:"oneOf,omitempty"`
	AnyOf []*SchemaOrRef `json:"anyOf,omitempty"`
	Not   *SchemaOrRef   `json:"not,omitempty"`

	// Access / status
	ReadOnly   bool `json:"readOnly,omitempty"`
	WriteOnly  bool `json:"writeOnly,omitempty"`
	Deprecated bool `json:"deprecated,omitempty"`
}

// AdditionalProperties is the union of a boolean and a schema. JSON Schema
// 2020-12 distinguishes "no additional properties" (false) from "any
// additional properties" (true) from "additional properties matching X".
//
// Spec: https://json-schema.org/draft/2020-12/json-schema-core#name-additionalproperties
type AdditionalProperties struct {
	Bool   *bool
	Schema *SchemaOrRef
}

// MarshalJSON outputs either a boolean or a schema, never both.
func (a *AdditionalProperties) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	if a.Schema != nil {
		return json.Marshal(a.Schema)
	}
	if a.Bool != nil {
		return json.Marshal(*a.Bool)
	}
	return []byte("null"), nil
}

// SchemaOrRef is either an inline Schema or a $ref. In 3.1.0 a $ref may
// carry sibling description / summary properties; this is the simplest path
// for attaching documentation to a referenced schema, and is preferred over
// allOf wrapping when only documentation is needed.
//
// Spec (reference form): https://spec.openapis.org/oas/v3.1.0#reference-object
type SchemaOrRef struct {
	Ref         string
	Summary     string
	Description string
	Value       *Schema
}

// MarshalJSON outputs a $ref (with optional siblings) or the inline schema.
func (s *SchemaOrRef) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	if s.Ref != "" {
		out := map[string]string{"$ref": s.Ref}
		if s.Summary != "" {
			out["summary"] = s.Summary
		}
		if s.Description != "" {
			out["description"] = s.Description
		}
		return json.Marshal(out)
	}
	return json.Marshal(s.Value)
}

// NewSchemaRef returns a SchemaOrRef pointing at #/components/schemas/<name>.
func NewSchemaRef(name string) *SchemaOrRef {
	return &SchemaOrRef{Ref: "#/components/schemas/" + name}
}

// Components is the components/* section. Schemas is the only sub-map this
// generator populates today.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#components-object
type Components struct {
	Schemas         map[string]*SchemaOrRef    `json:"-"`
	SecuritySchemes map[string]*SecurityScheme `json:"-"`
}

// Empty reports whether the components section has no content to emit.
func (c *Components) Empty() bool {
	return c == nil || (len(c.Schemas) == 0 && len(c.SecuritySchemes) == 0)
}

// MarshalJSON emits components sub-maps with keys in lexicographic order so
// the whole document is byte-stable across runs.
func (c *Components) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	if err := writeSortedKeys(&buf, &first, "schemas", c.Schemas); err != nil {
		return nil, err
	}
	if err := writeSortedKeys(&buf, &first, "securitySchemes", c.SecuritySchemes); err != nil {
		return nil, err
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// writeSortedKeys emits a `"field":{...}` JSON object into buf with the
// entries of m sorted by key. It writes a leading comma when *first is
// false, then sets *first to false. If m is empty, nothing is written and
// *first is unchanged — callers can chain calls without tracking emptiness
// themselves.
func writeSortedKeys[V any](buf *bytes.Buffer, first *bool, field string, m map[string]V) error {
	if len(m) == 0 {
		return nil
	}
	if !*first {
		buf.WriteByte(',')
	}
	*first = false
	buf.WriteByte('"')
	buf.WriteString(field)
	buf.WriteString(`":{`)
	for i, k := range slices.Sorted(maps.Keys(m)) {
		if i > 0 {
			buf.WriteByte(',')
		}
		key, err := json.Marshal(k)
		if err != nil {
			return err
		}
		buf.Write(key)
		buf.WriteByte(':')
		data, err := json.Marshal(m[k])
		if err != nil {
			return err
		}
		buf.Write(data)
	}
	buf.WriteByte('}')
	return nil
}

// SecurityScheme defines an authentication mechanism.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecurityScheme struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty"`
}

// OAuthFlows declares the OAuth2 flows supported by a security scheme.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow describes one OAuth2 flow. Per spec, scopes is required, even
// when empty (output as {}); the missing omitempty here is deliberate.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#oauth-flow-object
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// SecurityRequirement is a map of scheme name to required scopes.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#security-requirement-object
type SecurityRequirement map[string][]string

// Tag is metadata for a top-level tag in the document.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#tag-object
type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs is a link to external documentation.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#external-documentation-object
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}
