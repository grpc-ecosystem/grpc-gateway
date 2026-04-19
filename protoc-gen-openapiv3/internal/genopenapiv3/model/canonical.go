package model

// Canonical Model - Internal representation of OpenAPI documents.
// The generator builds this model, and adapters transform it to version-specific output.
//
// Design principles:
// - Pure data storage, no serialization logic
// - Version-agnostic (no JSON/YAML tags)
// - Adapters handle all version-specific transformations and serialization
// - Uses pointers for optional fields to distinguish "not set" from "zero value"

// Document is the canonical representation of an OpenAPI document.
type Document struct {
	OpenAPIVersion    string
	Info              *Info
	JSONSchemaDialect string
	Servers           []*Server
	Paths             *Paths
	Components        *Components
	Security          []SecurityRequirement
	Tags              []*Tag
	ExternalDocs      *ExternalDocs
}

// Info provides metadata about the API.
type Info struct {
	Title          string
	Summary        string
	Description    string
	TermsOfService string
	Contact        *Contact
	License        *License
	Version        string
}

// Contact information for the API.
type Contact struct {
	Name  string
	URL   string
	Email string
}

// License information for the API.
type License struct {
	Name       string
	Identifier string
	URL        string
}

// Server represents a server.
type Server struct {
	URL         string
	Description string
	Variables   map[string]*ServerVariable
}

// ServerVariable represents a server variable for URL template substitution.
type ServerVariable struct {
	Enum        []string
	Default     string
	Description string
}

// Paths holds the relative paths to individual endpoints.
type Paths struct {
	Items map[string]*PathItem
}

// PathItem describes operations available on a single path.
type PathItem struct {
	Ref         string
	Summary     string
	Description string
	Get         *Operation
	Put         *Operation
	Post        *Operation
	Delete      *Operation
	Options     *Operation
	Head        *Operation
	Patch       *Operation
	Trace       *Operation
	Servers     []*Server
	Parameters  []*ParameterOrRef
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags         []string
	Summary      string
	Description  string
	ExternalDocs *ExternalDocs
	OperationID  string
	Parameters   []*ParameterOrRef
	RequestBody  *RequestBodyOrRef
	Responses    *Responses
	Callbacks    map[string]*CallbackOrRef
	Deprecated   bool
	Security     []SecurityRequirement
	Servers      []*Server
}

// ExternalDocs allows referencing external documentation.
type ExternalDocs struct {
	Description string
	URL         string
}

// ParameterOrRef is either an inline Parameter or a reference.
type ParameterOrRef struct {
	Ref   string
	Value *Parameter
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name            string
	In              string
	Description     string
	Required        bool
	Deprecated      bool
	AllowEmptyValue bool
	Style           string
	Explode         *bool
	AllowReserved   bool
	Schema          *SchemaOrRef
	Examples        map[string]*Example
	Content         map[string]*MediaType
}

// RequestBodyOrRef is either an inline RequestBody or a reference.
type RequestBodyOrRef struct {
	Ref   string
	Value *RequestBody
}

// RequestBody describes a single request body.
type RequestBody struct {
	Description string
	Content     map[string]*MediaType
	Required    bool
}

// MediaType provides schema and examples for a media type.
type MediaType struct {
	Schema   *SchemaOrRef
	Examples map[string]*Example
	Encoding map[string]*Encoding
}

// Encoding defines encoding for a single property.
type Encoding struct {
	ContentType   string
	Headers       map[string]*HeaderOrRef
	Style         string
	Explode       *bool
	AllowReserved bool
}

// Responses is a container for expected responses.
type Responses struct {
	Default *ResponseOrRef
	Codes   map[string]*ResponseOrRef
}

// ResponseOrRef is either an inline Response or a reference.
type ResponseOrRef struct {
	Ref   string
	Value *Response
}

// Response describes a single response from an API operation.
type Response struct {
	Description string
	Headers     map[string]*HeaderOrRef
	Content     map[string]*MediaType
	Links       map[string]*LinkOrRef
}

// HeaderOrRef is either an inline Header or a reference.
type HeaderOrRef struct {
	Ref   string
	Value *Header
}

// Header follows the same structure as Parameter but without "name" and "in".
type Header struct {
	Description     string
	Required        bool
	Deprecated      bool
	AllowEmptyValue bool
	Style           string
	Explode         *bool
	AllowReserved   bool
	Schema          *SchemaOrRef
	Examples        map[string]*Example
	Content         map[string]*MediaType
}

// LinkOrRef is either an inline Link or a reference.
type LinkOrRef struct {
	Ref   string
	Value *Link
}

// Link represents a possible design-time link for a response.
type Link struct {
	OperationRef string
	OperationID  string
	Parameters   map[string]any
	RequestBody  any
	Description  string
	Server       *Server
}

// CallbackOrRef is either an inline Callback or a reference.
type CallbackOrRef struct {
	Ref   string
	Value map[string]*PathItem
}

// SchemaOrRef is either an inline Schema or a reference.
// In 3.1.0, $ref can have sibling summary/description; in 3.0.x it cannot.
type SchemaOrRef struct {
	Ref         string
	Summary     string // For $ref overrides (3.1.0+)
	Description string // For $ref overrides (3.1.0+)
	Value       *Schema
}

// Schema represents a JSON Schema object.
// Type and IsNullable are version-agnostic; adapters convert to version-specific format:
// - 3.0.x: type + nullable: true
// - 3.1.0: type array ["string", "null"]
type Schema struct {
	Type       string
	IsNullable bool

	Format       string
	Title        string
	Description  string
	Default      any
	Examples     []any
	Deprecated   bool
	ReadOnly     bool
	WriteOnly    bool
	ExternalDocs *ExternalDocs

	// Numeric validation
	MultipleOf       *float64
	Minimum          *float64
	Maximum          *float64
	ExclusiveMinimum *float64
	ExclusiveMaximum *float64

	// String validation
	MinLength *uint64
	MaxLength *uint64
	Pattern   string

	// Array validation
	MinItems    *uint64
	MaxItems    *uint64
	UniqueItems bool
	Items       *SchemaOrRef

	// Object validation
	MinProperties        *uint64
	MaxProperties        *uint64
	Required             []string
	Properties           map[string]*SchemaOrRef
	AdditionalProperties *AdditionalProperties

	// Composition
	AllOf         []*SchemaOrRef
	AnyOf         []*SchemaOrRef
	OneOf         []*SchemaOrRef
	Not           *SchemaOrRef
	Discriminator *Discriminator

	// Enum
	Enum []any

	// Extensions
	Extensions map[string]any
}

// AdditionalProperties can be a boolean or a schema.
type AdditionalProperties struct {
	Allowed bool
	Schema  *SchemaOrRef
}

// Discriminator helps with polymorphism.
type Discriminator struct {
	PropertyName string
	Mapping      map[string]string
}

// Example represents an example value.
type Example struct {
	Ref           string
	Summary       string
	Description   string
	Value         any
	ExternalValue string
}

// Components holds reusable objects.
type Components struct {
	Schemas         map[string]*SchemaOrRef
	Responses       map[string]*ResponseOrRef
	Parameters      map[string]*ParameterOrRef
	Examples        map[string]*Example
	RequestBodies   map[string]*RequestBodyOrRef
	Headers         map[string]*HeaderOrRef
	SecuritySchemes map[string]*SecuritySchemeOrRef
	Links           map[string]*LinkOrRef
	Callbacks       map[string]*CallbackOrRef
	PathItems       map[string]*PathItem
}

// SecuritySchemeOrRef is either an inline SecurityScheme or a reference.
type SecuritySchemeOrRef struct {
	Ref   string
	Value *SecurityScheme
}

// SecurityScheme defines a security scheme.
type SecurityScheme struct {
	Type             string
	Description      string
	Name             string
	In               string
	Scheme           string
	BearerFormat     string
	Flows            *OAuthFlows
	OpenIDConnectURL string
}

// OAuthFlows allows configuration of supported OAuth Flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow
	Password          *OAuthFlow
	ClientCredentials *OAuthFlow
	AuthorizationCode *OAuthFlow
}

// OAuthFlow configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

// SecurityRequirement lists required security schemes.
type SecurityRequirement map[string][]string

// Tag adds metadata to a single tag.
type Tag struct {
	Name         string
	Description  string
	ExternalDocs *ExternalDocs
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
