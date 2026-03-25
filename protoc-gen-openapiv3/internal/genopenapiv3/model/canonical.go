package model

import (
	"bytes"
	"encoding/json"
	"sort"

	"go.yaml.in/yaml/v3"
)

// Canonical Model - Internal representation of OpenAPI documents.
// The generator builds this model, and adapters handle version-specific output.
//
// Design principles:
// - Fields ordered per OpenAPI specification for deterministic output
// - JSON/YAML tags for direct serialization
// - Custom marshaling only where version-specific behavior is needed
// - Uses pointers for optional fields to distinguish "not set" from "zero value"

// Document is the canonical representation of an OpenAPI document.
// Field order matches OpenAPI specification.
type Document struct {
	OpenAPIVersion    string                `json:"openapi" yaml:"openapi"`
	Info              *Info                 `json:"info" yaml:"info"`
	JSONSchemaDialect string                `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	Servers           []*Server             `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths             *Paths                `json:"paths,omitempty" yaml:"paths,omitempty"`
	Components        *Components           `json:"components,omitempty" yaml:"components,omitempty"`
	Security          []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Tags              []*Tag                `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs      *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Info provides metadata about the API.
// Field order: title, summary, description, termsOfService, contact, license, version
type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Summary        string   `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`
}

// Contact information for the API.
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License information for the API.
type License struct {
	Name       string `json:"name" yaml:"name"`
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Server represents a server.
type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable represents a server variable for URL template substitution.
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// Paths holds the relative paths to individual endpoints.
// Uses custom marshaling for sorted, deterministic output.
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
	Ref         string            `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string            `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation        `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation        `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation        `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation        `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation        `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation        `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation        `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation        `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []*Server         `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []*ParameterOrRef `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags         []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   []*ParameterOrRef     `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *RequestBodyOrRef     `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    *Responses            `json:"responses,omitempty" yaml:"responses,omitempty"`
	Callbacks    map[string]*CallbackOrRef `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []*Server             `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// ExternalDocs allows referencing external documentation.
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

// ParameterOrRef is either an inline Parameter or a reference.
// Uses custom marshaling to flatten the structure.
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
	Schema   *SchemaOrRef          `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding  `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding defines encoding for a single property.
type Encoding struct {
	ContentType   string                 `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]*HeaderOrRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                 `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool                  `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                   `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// Responses is a container for expected responses.
// Uses custom marshaling for sorted output.
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
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: code}
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
	OperationRef string         `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string         `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *Server        `json:"server,omitempty" yaml:"server,omitempty"`
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
// In 3.1.0, $ref can have sibling summary/description.
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

// Schema represents a JSON Schema object.
// Uses custom marshaling for version-specific type handling (nullable).
type Schema struct {
	// Type and nullability - handled specially in marshaling
	Type       string `json:"-" yaml:"-"`
	IsNullable bool   `json:"-" yaml:"-"`

	Format       string        `json:"format,omitempty" yaml:"format,omitempty"`
	Title        string        `json:"title,omitempty" yaml:"title,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	Default      any           `json:"default,omitempty" yaml:"default,omitempty"`
	Examples     map[string]*Example `json:"examples,omitempty" yaml:"examples,omitempty"`
	Deprecated   bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	ReadOnly     bool          `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly    bool          `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

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
	MinProperties        *uint64                  `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	MaxProperties        *uint64                  `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	Required             []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	Properties           map[string]*SchemaOrRef  `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *AdditionalProperties    `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

	// Composition
	AllOf         []*SchemaOrRef `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf         []*SchemaOrRef `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf         []*SchemaOrRef `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not           *SchemaOrRef   `json:"not,omitempty" yaml:"not,omitempty"`
	Discriminator *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

	// Enum
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`

	// Extensions
	Extensions map[string]any `json:"-" yaml:"-"`
}

// schemaAlias is used to avoid infinite recursion in MarshalJSON/MarshalYAML.
type schemaAlias Schema

// MarshalJSON handles type and nullable specially for 3.1.0 format.
func (s *Schema) MarshalJSON() ([]byte, error) {
	type schemaWithType struct {
		Type any `json:"type,omitempty"`
		*schemaAlias
	}

	out := schemaWithType{schemaAlias: (*schemaAlias)(s)}

	if s.Type != "" {
		if s.IsNullable {
			out.Type = []string{s.Type, "null"}
		} else {
			out.Type = s.Type
		}
	}

	return json.Marshal(out)
}

// MarshalYAML handles type and nullable specially for 3.1.0 format.
func (s *Schema) MarshalYAML() (any, error) {
	// Build a map manually to control field order and handle type specially
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Type first
	if s.Type != "" {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "type"}
		valueNode := &yaml.Node{}
		if s.IsNullable {
			if err := valueNode.Encode([]string{s.Type, "null"}); err != nil {
				return nil, err
			}
		} else {
			valueNode.Kind = yaml.ScalarNode
			valueNode.Value = s.Type
		}
		node.Content = append(node.Content, keyNode, valueNode)
	}

	// Use reflection-free approach: encode as schemaAlias and merge
	alias := (*schemaAlias)(s)
	aliasNode := &yaml.Node{}
	if err := aliasNode.Encode(alias); err != nil {
		return nil, err
	}

	// Append alias content (skip if empty)
	if aliasNode.Kind == yaml.MappingNode {
		node.Content = append(node.Content, aliasNode.Content...)
	}

	return node, nil
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

// Discriminator helps with polymorphism.
type Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
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
// Field order matches OpenAPI specification.
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

// ExampleOrRef is either an inline Example or a reference.
// Deprecated: Use Example directly with Ref field.
type ExampleOrRef struct {
	Ref   string
	Value *Example
}

// SecuritySchemeOrRef is either an inline SecurityScheme or a reference.
type SecuritySchemeOrRef struct {
	Ref   string          `json:"-" yaml:"-"`
	Value *SecurityScheme `json:"-" yaml:"-"`
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

// SecurityScheme defines a security scheme.
type SecurityScheme struct {
	Type             string      `json:"type" yaml:"type"`
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty" yaml:"name,omitempty"`
	In               string      `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlows allows configuration of supported OAuth Flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlow configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

// SecurityRequirement lists required security schemes.
type SecurityRequirement map[string][]string

// Tag adds metadata to a single tag.
type Tag struct {
	Name         string        `json:"name" yaml:"name"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// NewDocument creates a new canonical Document with defaults.
func NewDocument(version string) *Document {
	return &Document{
		OpenAPIVersion: version,
		Info:           &Info{Version: "1.0.0"},
		Paths:          &Paths{Items: make(map[string]*PathItem)},
		Components: &Components{
			Schemas:         make(map[string]*SchemaOrRef),
			Responses:       make(map[string]*ResponseOrRef),
			Parameters:      make(map[string]*ParameterOrRef),
			Examples:        make(map[string]*Example),
			RequestBodies:   make(map[string]*RequestBodyOrRef),
			Headers:         make(map[string]*HeaderOrRef),
			SecuritySchemes: make(map[string]*SecuritySchemeOrRef),
			Links:           make(map[string]*LinkOrRef),
			Callbacks:       make(map[string]*CallbackOrRef),
		},
	}
}

// NewResponses creates a new Responses container.
func NewResponses() *Responses {
	return &Responses{
		Codes: make(map[string]*ResponseOrRef),
	}
}

// NewSchema creates a new Schema with the given type.
func NewSchema(typ string) *Schema {
	return &Schema{Type: typ}
}

// NewSchemaRef creates a schema reference.
func NewSchemaRef(ref string) *SchemaOrRef {
	return &SchemaOrRef{Ref: ref}
}

// NewInlineSchema wraps a schema as SchemaOrRef.
func NewInlineSchema(s *Schema) *SchemaOrRef {
	return &SchemaOrRef{Value: s}
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
