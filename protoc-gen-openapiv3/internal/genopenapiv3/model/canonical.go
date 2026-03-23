package model

// Canonical Model - Version-agnostic internal representation of OpenAPI documents.
// The generator builds this model, and adapters convert it to version-specific output.
//
// Design principles:
// - Simple, idiomatic Go types (no JSON tags - those belong in output types)
// - Rich enough to capture all features from OpenAPI 3.0.x and 3.1.0
// - Uses pointers for optional fields to distinguish "not set" from "zero value"
// - No version-specific concepts (e.g., uses IsNullable bool, not type arrays)

// Document is the canonical representation of an OpenAPI document.
type Document struct {
	// OpenAPI version to generate (e.g., "3.1.0", "3.0.3")
	OpenAPIVersion string

	Info          *Info
	Servers       []*Server
	Paths         *Paths
	Components    *Components
	Security      []SecurityRequirement
	Tags              []*Tag
	ExternalDocs      *ExternalDocs
	JSONSchemaDialect string // 3.1.0 only
}

// Info provides metadata about the API.
type Info struct {
	Title          string
	Description    string
	TermsOfService string
	Contact        *Contact
	License        *License
	Version        string
	Summary        string // 3.1.0 only
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
	URL        string
	Identifier string // 3.1.0 only (SPDX identifier)
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
	Ref         string // $ref to another PathItem
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
	Ref   string // If set, this is a $ref
	Value *Parameter
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name            string
	In              string // "query", "header", "path", "cookie"
	Description     string
	Required        bool
	Deprecated      bool
	AllowEmptyValue bool
	Style           string
	Explode         *bool
	AllowReserved   bool
	Schema          *SchemaOrRef
	Examples        []*Example // Canonical: list of examples
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
	Examples []*Example // Canonical: list of examples
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
	Codes   map[string]*ResponseOrRef // HTTP status codes
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
	Examples        []*Example
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
type SchemaOrRef struct {
	Ref         string
	Summary     string // 3.1.0: $ref can have summary
	Description string // 3.1.0: $ref can have description
	Value       *Schema
}

// Schema represents a JSON Schema object (version-agnostic).
type Schema struct {
	// Type information
	Type   string // Single type: "string", "integer", "object", etc.
	Format string

	// Nullability (version-agnostic - adapters convert appropriately)
	IsNullable bool

	// Metadata
	Title       string
	Description string
	Default     any
	Examples    []*Example // Canonical: list of examples
	Deprecated  bool
	ReadOnly    bool
	WriteOnly   bool
	ExternalDocs *ExternalDocs

	// Validation - Numeric
	MultipleOf       *float64
	Minimum          *float64
	Maximum          *float64
	ExclusiveMinimum *float64 // 3.1.0: numeric, 3.0.x: boolean (adapter handles)
	ExclusiveMaximum *float64

	// Validation - String
	MinLength *uint64
	MaxLength *uint64
	Pattern   string

	// Validation - Array
	MinItems    *uint64
	MaxItems    *uint64
	UniqueItems bool
	Items       *SchemaOrRef

	// Validation - Object
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
	Allowed bool         // If true, any additional properties allowed
	Schema  *SchemaOrRef // If set, additional properties must match this schema
}

// Discriminator helps with polymorphism.
type Discriminator struct {
	PropertyName string
	Mapping      map[string]string
}

// Example represents an example value.
type Example struct {
	Name          string // Optional name for named examples
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
	Examples        map[string]*ExampleOrRef
	RequestBodies   map[string]*RequestBodyOrRef
	Headers         map[string]*HeaderOrRef
	SecuritySchemes map[string]*SecuritySchemeOrRef
	Links           map[string]*LinkOrRef
	Callbacks       map[string]*CallbackOrRef
	PathItems       map[string]*PathItem // 3.1.0 only
}

// ExampleOrRef is either an inline Example or a reference.
type ExampleOrRef struct {
	Ref   string
	Value *Example
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
			Examples:        make(map[string]*ExampleOrRef),
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
