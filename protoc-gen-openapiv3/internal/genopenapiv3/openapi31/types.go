package openapi31

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/transform"
	"go.yaml.in/yaml/v3"
)

// Document is the OpenAPI 3.1.0 output document.
type Document struct {
	OpenAPI      string                        `json:"openapi" yaml:"openapi"`
	Info         *transform.Info               `json:"info" yaml:"info"`
	Servers      []*transform.Server           `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        *Paths                        `json:"paths,omitempty" yaml:"paths,omitempty"`
	Components   *Components                   `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []transform.SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []*transform.Tag              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *transform.ExternalDocs       `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Paths holds the relative paths to individual endpoints.
type Paths struct {
	Items map[string]*PathItem `json:"-" yaml:"-"`
}

// MarshalJSON outputs paths in sorted order.
func (p *Paths) MarshalJSON() ([]byte, error) {
	if p == nil || len(p.Items) == 0 {
		return []byte("{}"), nil
	}
	return marshalOrderedMap(p.Items)
}

// MarshalYAML outputs paths in sorted order.
func (p *Paths) MarshalYAML() (any, error) {
	if p == nil || len(p.Items) == 0 {
		return nil, nil
	}
	return marshalOrderedMapYAML(p.Items)
}

// PathItem describes operations available on a single path.
type PathItem struct {
	Ref         string              `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation          `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation          `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation          `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation          `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation          `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation          `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation          `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation          `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []*transform.Server `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []*ParameterOrRef   `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags         []string                        `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string                          `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                          `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *transform.ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string                          `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   []*ParameterOrRef               `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *RequestBodyOrRef               `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    *Responses                      `json:"responses,omitempty" yaml:"responses,omitempty"`
	Callbacks    map[string]*CallbackOrRef       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated   bool                            `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []transform.SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []*transform.Server             `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// ParameterOrRef is either an inline Parameter or a reference.
type ParameterOrRef struct {
	Ref   string     `json:"-" yaml:"-"`
	Value *Parameter `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the parameter fields.
func (p *ParameterOrRef) MarshalJSON() ([]byte, error) {
	if p.Ref != "" {
		return json.Marshal(map[string]string{"$ref": p.Ref})
	}
	if p.Value != nil {
		return json.Marshal(p.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the parameter fields.
func (p *ParameterOrRef) MarshalYAML() (any, error) {
	if p.Ref != "" {
		return map[string]string{"$ref": p.Ref}, nil
	}
	return p.Value, nil
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name            string                `json:"name" yaml:"name"`
	In              string                `json:"in" yaml:"in"`
	Description     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool                 `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *SchemaOrRef          `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// RequestBodyOrRef is either an inline RequestBody or a reference.
type RequestBodyOrRef struct {
	Ref   string       `json:"-" yaml:"-"`
	Value *RequestBody `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the request body fields.
func (r *RequestBodyOrRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return json.Marshal(map[string]string{"$ref": r.Ref})
	}
	if r.Value != nil {
		return json.Marshal(r.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the request body fields.
func (r *RequestBodyOrRef) MarshalYAML() (any, error) {
	if r.Ref != "" {
		return map[string]string{"$ref": r.Ref}, nil
	}
	return r.Value, nil
}

// RequestBody describes a single request body.
type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
}

// MediaType provides schema and examples for a media type.
type MediaType struct {
	Schema   *SchemaOrRef         `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples map[string]*Example  `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding defines encoding for a single property.
type Encoding struct {
	ContentType   string                  `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]*HeaderOrRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                  `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool                   `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                    `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// Responses is a container for expected responses.
type Responses struct {
	Default *ResponseOrRef            `json:"-" yaml:"-"`
	Codes   map[string]*ResponseOrRef `json:"-" yaml:"-"`
}

// MarshalJSON outputs responses with default first, then sorted codes.
func (r *Responses) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("{}"), nil
	}

	var buf bytes.Buffer
	buf.WriteString("{")
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
	for code := range r.Codes {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	for _, code := range codes {
		if !first {
			buf.WriteString(",")
		}
		keyBytes, _ := json.Marshal(code)
		buf.Write(keyBytes)
		buf.WriteString(":")
		data, err := json.Marshal(r.Codes[code])
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		first = false
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

// MarshalYAML outputs responses with default first, then sorted codes.
func (r *Responses) MarshalYAML() (any, error) {
	if r == nil {
		return nil, nil
	}

	node := &yaml.Node{Kind: yaml.MappingNode}

	if r.Default != nil {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "default"}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(r.Default); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valueNode)
	}

	codes := make([]string, 0, len(r.Codes))
	for code := range r.Codes {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	for _, code := range codes {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.DoubleQuotedStyle, Value: code}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(r.Codes[code]); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valueNode)
	}

	return node, nil
}

// ResponseOrRef is either an inline Response or a reference.
type ResponseOrRef struct {
	Ref   string    `json:"-" yaml:"-"`
	Value *Response `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the response fields.
func (r *ResponseOrRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return json.Marshal(map[string]string{"$ref": r.Ref})
	}
	if r.Value != nil {
		return json.Marshal(r.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the response fields.
func (r *ResponseOrRef) MarshalYAML() (any, error) {
	if r.Ref != "" {
		return map[string]string{"$ref": r.Ref}, nil
	}
	return r.Value, nil
}

// Response describes a single response from an API operation.
type Response struct {
	Description string                  `json:"description" yaml:"description"`
	Headers     map[string]*HeaderOrRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType   `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*LinkOrRef   `json:"links,omitempty" yaml:"links,omitempty"`
}

// HeaderOrRef is either an inline Header or a reference.
type HeaderOrRef struct {
	Ref   string  `json:"-" yaml:"-"`
	Value *Header `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the header fields.
func (h *HeaderOrRef) MarshalJSON() ([]byte, error) {
	if h.Ref != "" {
		return json.Marshal(map[string]string{"$ref": h.Ref})
	}
	if h.Value != nil {
		return json.Marshal(h.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the header fields.
func (h *HeaderOrRef) MarshalYAML() (any, error) {
	if h.Ref != "" {
		return map[string]string{"$ref": h.Ref}, nil
	}
	return h.Value, nil
}

// Header follows the same structure as Parameter but without "name" and "in".
type Header struct {
	Description     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool                 `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *SchemaOrRef          `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// LinkOrRef is either an inline Link or a reference.
type LinkOrRef struct {
	Ref   string `json:"-" yaml:"-"`
	Value *Link  `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the link fields.
func (l *LinkOrRef) MarshalJSON() ([]byte, error) {
	if l.Ref != "" {
		return json.Marshal(map[string]string{"$ref": l.Ref})
	}
	if l.Value != nil {
		return json.Marshal(l.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the link fields.
func (l *LinkOrRef) MarshalYAML() (any, error) {
	if l.Ref != "" {
		return map[string]string{"$ref": l.Ref}, nil
	}
	return l.Value, nil
}

// Link represents a possible design-time link for a response.
type Link struct {
	OperationRef string            `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string            `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]any    `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  any               `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *transform.Server `json:"server,omitempty" yaml:"server,omitempty"`
}

// CallbackOrRef is either an inline Callback or a reference.
type CallbackOrRef struct {
	Ref   string               `json:"-" yaml:"-"`
	Value map[string]*PathItem `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the callback map.
func (c *CallbackOrRef) MarshalJSON() ([]byte, error) {
	if c.Ref != "" {
		return json.Marshal(map[string]string{"$ref": c.Ref})
	}
	if c.Value != nil {
		return marshalOrderedMap(c.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the callback map.
func (c *CallbackOrRef) MarshalYAML() (any, error) {
	if c.Ref != "" {
		return map[string]string{"$ref": c.Ref}, nil
	}
	return marshalOrderedMapYAML(c.Value)
}

// SchemaOrRef is either an inline Schema or a reference.
// In OpenAPI 3.1.0, $ref can have sibling summary/description.
type SchemaOrRef struct {
	Ref         string  `json:"-" yaml:"-"`
	Summary     string  `json:"-" yaml:"-"` // 3.1.0: $ref can have summary
	Description string  `json:"-" yaml:"-"` // 3.1.0: $ref can have description
	Value       *Schema `json:"-" yaml:"-"`
}

// MarshalJSON outputs the schema with proper handling of $ref siblings.
func (s *SchemaOrRef) MarshalJSON() ([]byte, error) {
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
	if s.Value != nil {
		return json.Marshal(s.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs the schema with proper handling of $ref siblings.
func (s *SchemaOrRef) MarshalYAML() (any, error) {
	if s.Ref != "" {
		out := map[string]string{"$ref": s.Ref}
		if s.Summary != "" {
			out["summary"] = s.Summary
		}
		if s.Description != "" {
			out["description"] = s.Description
		}
		return out, nil
	}
	return s.Value, nil
}

// Schema represents a JSON Schema object for OpenAPI 3.1.0.
// Type is output as string or []string for nullable.
type Schema struct {
	// Type - can be string or []string for nullable (e.g., ["string", "null"])
	Type any `json:"type,omitempty" yaml:"type,omitempty"`

	Format       string                  `json:"format,omitempty" yaml:"format,omitempty"`
	Title        string                  `json:"title,omitempty" yaml:"title,omitempty"`
	Description  string                  `json:"description,omitempty" yaml:"description,omitempty"`
	Default      any                     `json:"default,omitempty" yaml:"default,omitempty"`
	Examples     []any                   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Deprecated   bool                    `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	ReadOnly     bool                    `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly    bool                    `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	ExternalDocs *transform.ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	// Numeric validation
	MultipleOf       *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`

	// String validation
	MinLength *uint64 `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength *uint64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Array validation
	MinItems    *uint64      `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems    *uint64      `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	UniqueItems bool         `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Items       *SchemaOrRef `json:"items,omitempty" yaml:"items,omitempty"`

	// Object validation
	MinProperties        *uint64                 `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	MaxProperties        *uint64                 `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	Required             []string                `json:"required,omitempty" yaml:"required,omitempty"`
	Properties           map[string]*SchemaOrRef `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *AdditionalProperties   `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

	// Composition
	AllOf         []*SchemaOrRef           `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf         []*SchemaOrRef           `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf         []*SchemaOrRef           `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not           *SchemaOrRef             `json:"not,omitempty" yaml:"not,omitempty"`
	Discriminator *transform.Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

	// Enum
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`
}

// AdditionalProperties can be a boolean or a schema.
type AdditionalProperties struct {
	Allowed bool         `json:"-" yaml:"-"`
	Schema  *SchemaOrRef `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a boolean or a schema.
func (a *AdditionalProperties) MarshalJSON() ([]byte, error) {
	if a.Schema != nil {
		return json.Marshal(a.Schema)
	}
	return json.Marshal(a.Allowed)
}

// MarshalYAML outputs either a boolean or a schema.
func (a *AdditionalProperties) MarshalYAML() (any, error) {
	if a.Schema != nil {
		return a.Schema, nil
	}
	return a.Allowed, nil
}

// Example represents an example value.
type Example struct {
	Ref           string `json:"-" yaml:"-"`
	Summary       string `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string `json:"description,omitempty" yaml:"description,omitempty"`
	Value         any    `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

// MarshalJSON handles $ref for examples.
func (e *Example) MarshalJSON() ([]byte, error) {
	if e.Ref != "" {
		return json.Marshal(map[string]string{"$ref": e.Ref})
	}
	type exampleAlias Example
	return json.Marshal((*exampleAlias)(e))
}

// MarshalYAML handles $ref for examples.
func (e *Example) MarshalYAML() (any, error) {
	if e.Ref != "" {
		return map[string]string{"$ref": e.Ref}, nil
	}
	type exampleAlias Example
	return (*exampleAlias)(e), nil
}

// Components holds reusable objects.
type Components struct {
	Schemas         map[string]*SchemaOrRef         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*ResponseOrRef       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*ParameterOrRef      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]*Example             `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBodyOrRef    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*HeaderOrRef         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeOrRef `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*LinkOrRef           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*CallbackOrRef       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	PathItems       map[string]*PathItem            `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
}

// SecuritySchemeOrRef is either an inline SecurityScheme or a reference.
type SecuritySchemeOrRef struct {
	Ref   string                    `json:"-" yaml:"-"`
	Value *transform.SecurityScheme `json:"-" yaml:"-"`
}

// MarshalJSON outputs either a $ref or the security scheme fields.
func (s *SecuritySchemeOrRef) MarshalJSON() ([]byte, error) {
	if s.Ref != "" {
		return json.Marshal(map[string]string{"$ref": s.Ref})
	}
	if s.Value != nil {
		return json.Marshal(s.Value)
	}
	return []byte("null"), nil
}

// MarshalYAML outputs either a $ref or the security scheme fields.
func (s *SecuritySchemeOrRef) MarshalYAML() (any, error) {
	if s.Ref != "" {
		return map[string]string{"$ref": s.Ref}, nil
	}
	return s.Value, nil
}


// Helper functions for ordered map marshaling

func marshalOrderedMap[V any](m map[string]V) ([]byte, error) {
	if len(m) == 0 {
		return []byte("{}"), nil
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteString("{")
	for i, key := range keys {
		if i > 0 {
			buf.WriteString(",")
		}
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteString(":")
		valBytes, err := json.Marshal(m[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func marshalOrderedMapYAML[V any](m map[string]V) (*yaml.Node, error) {
	if len(m) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	node := &yaml.Node{Kind: yaml.MappingNode}
	for _, key := range keys {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(m[key]); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valueNode)
	}
	return node, nil
}
