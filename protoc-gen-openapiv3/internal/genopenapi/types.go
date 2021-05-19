package genopenapi

import (
	"regexp"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

// https://swagger.io/specification/#openapi-object
type Openapi struct {
	ExtensionProps
	OpenAPI      string               `json:"openapi" yaml:"openapi"` // Required
	Components   Components           `json:"components,omitempty" yaml:"components,omitempty"`
	Info         Info                `json:"info" yaml:"info"`   // Required
	Paths        Paths                `json:"paths" yaml:"paths"` // Required
	Security     SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      Servers              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Tags         Tags                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocs        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// https://swagger.io/specification/#components-object
type Components struct {
	ExtensionProps
	Schemas         Schemas         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Parameters      ParametersMap   `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Headers         Headers         `json:"headers,omitempty" yaml:"headers,omitempty"`
	RequestBodies   RequestBodies   `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Responses       Responses       `json:"responses,omitempty" yaml:"responses,omitempty"`
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Examples        Examples        `json:"examples,omitempty" yaml:"examples,omitempty"`
	Links           Links           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       Callbacks       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

// https://swagger.io/specification/#info-object
type Info struct {
	ExtensionProps
	Title          string   `json:"title" yaml:"title"` // Required
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"` // Required
}

// https://swagger.io/specification/#contact-object
type Contact struct {
	ExtensionProps
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// https://swagger.io/specification/#license-object
type License struct {
	ExtensionProps
	Name string `json:"name" yaml:"name"` // Required
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type Paths map[string]*PathItem

// https://swagger.io/specification/#path-item-object
type PathItem struct {
	ExtensionProps
	Ref         string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Connect     *Operation `json:"connect,omitempty" yaml:"connect,omitempty"`
	Delete      *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Get         *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Head        *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Options     *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	Patch       *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Post        *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Put         *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Trace       *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     Servers    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// https://swagger.io/specification/#operation-object
type Operation struct {
	ExtensionProps

	// Optional tags for documentation.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Optional short summary.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`

	// Optional description. Should use CommonMark syntax.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Optional operation ID.
	OperationID string `json:"operationId,omitempty" yaml:"operationId,omitempty"`

	// Optional parameters.
	Parameters Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	// Optional body parameter.
	RequestBody *RequestBodyRef `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`

	// Responses.
	Responses Responses `json:"responses" yaml:"responses"` // Required

	// Optional callbacks
	Callbacks Callbacks `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`

	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Optional security requirements that overrides top-level security.
	Security *SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`

	// Optional servers that overrides top-level servers.
	Servers *Servers `json:"servers,omitempty" yaml:"servers,omitempty"`

	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type SecuritySchemeRef struct {
	Ref   string
	Value *SecurityScheme
}

type SecurityRequirements []SecurityRequirement
// https://swagger.io/specification/#security-requirement-object
type SecurityRequirement map[string][]string

type SecuritySchemes map[string]*SecuritySchemeRef
// https://swagger.io/specification/#security-scheme-object
type SecurityScheme struct {
	ExtensionProps

	Type             string      `json:"type,omitempty" yaml:"type,omitempty"`
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty" yaml:"name,omitempty"`
	In               string      `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// https://swagger.io/specification/#oauth-flows-object
type OAuthFlows struct {
	ExtensionProps
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// https://swagger.io/specification/#oauth-flow-object
type OAuthFlow struct {
	ExtensionProps
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"`
}

type Servers []*Server

// https://swagger.io/specification/#server-object
type Server struct {
	ExtensionProps
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// https://swagger.io/specification/#server-variable-object
type ServerVariable struct {
	ExtensionProps
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default,omitempty" yaml:"default,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

type Tags []*Tag

// https://swagger.io/specification/#tag-object
type Tag struct {
	ExtensionProps
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type Schemas map[string]*SchemaRef

type SchemaRef struct {
	Ref   string
	Value *Schema
}

func (value *SchemaRef) MarshalJSON() ([]byte, error) {
	return marshalRef(value.Ref, value.Value)
}

func (value *SchemaRef) UnmarshalJSON(data []byte) error {
	return unmarshalRef(data, &value.Ref, &value.Value)
}

type SchemaRefs []*SchemaRef

// https://swagger.io/specification/#schema-object
type Schema struct {
	ExtensionProps

	OneOf        SchemaRefs    `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf        SchemaRefs    `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	AllOf        SchemaRefs    `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	Not          *SchemaRef    `json:"not,omitempty" yaml:"not,omitempty"`
	Type         string        `json:"type,omitempty" yaml:"type,omitempty"`
	Title        string        `json:"title,omitempty" yaml:"title,omitempty"`
	Format       string        `json:"format,omitempty" yaml:"format,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	Enum         []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default      interface{}   `json:"default,omitempty" yaml:"default,omitempty"`
	Example      interface{}   `json:"example,omitempty" yaml:"example,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	// Object-related, here for struct compactness
	AdditionalPropertiesAllowed *bool `json:"-" multijson:"additionalProperties,omitempty" yaml:"-"`
	// Array-related, here for struct compactness
	UniqueItems bool `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	// Number-related, here for struct compactness
	ExclusiveMinimum bool `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum bool `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	// Properties
	Nullable        bool        `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnly        bool        `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly       bool        `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	XML             interface{} `json:"xml,omitempty" yaml:"xml,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Number
	Minimum    float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum    float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MultipleOf float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`

	// String
	MinLength       uint64  `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength       uint64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern         string  `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	compiledPattern *regexp.Regexp

	// Array
	MinItems uint64     `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems uint64    `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	Items    *SchemaRef `json:"items,omitempty" yaml:"items,omitempty"`

	// Object
	Required             []string       `json:"required,omitempty" yaml:"required,omitempty"`
	Properties           Schemas        `json:"properties,omitempty" yaml:"properties,omitempty"`
	MinProperties        uint64         `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	MaxProperties        uint64         `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	AdditionalProperties *SchemaRef     `json:"-" multijson:"additionalProperties,omitempty" yaml:"-"`
	Discriminator        *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
}

// https://swagger.io/specification/#external-documentation-object
type ExternalDocs struct {
	ExtensionProps

	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// https://swagger.io/specification/#discriminator-object
type Discriminator struct {
	ExtensionProps
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

type ParametersMap map[string]*ParameterRef

type Parameters []*ParameterRef

type ParameterRef struct {
	Ref   string
	Value *Parameter
}

func (value *ParameterRef) MarshalJSON() ([]byte, error) {
	return marshalRef(value.Ref, value.Value)
}

func (value *ParameterRef) UnmarshalJSON(data []byte) error {
	return unmarshalRef(data, &value.Ref, &value.Value)
}

type Headers map[string]*HeaderRef

type HeaderRef struct {
	Ref   string
	Value *Header
}

func (value *HeaderRef) MarshalJSON() ([]byte, error) {
	return marshalRef(value.Ref, value.Value)
}

func (value *HeaderRef) UnmarshalJSON(data []byte) error {
	return unmarshalRef(data, &value.Ref, &value.Value)
}

// https://swagger.io/specification/#header-object
type Header struct {
	ExtensionProps

	// Optional description. Should use CommonMark syntax.
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Deprecated  bool        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Required    bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      *SchemaRef  `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty" yaml:"example,omitempty"`
	Examples    Examples    `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content     Content     `json:"content,omitempty" yaml:"content,omitempty"`
}

type RequestBodies map[string]*RequestBodyRef

type RequestBodyRef struct {
	Ref   string
	Value *RequestBody
}

func (value *RequestBodyRef) MarshalJSON() ([]byte, error) {
	return marshalRef(value.Ref, value.Value)
}

func (value *RequestBodyRef) UnmarshalJSON(data []byte) error {
	return unmarshalRef(data, &value.Ref, &value.Value)
}

// https://swagger.io/specification/#request-body-object
type RequestBody struct {
	ExtensionProps
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Content     Content `json:"content,omitempty" yaml:"content,omitempty"`
}

type Responses map[string]*ResponseRef

type ResponseRef struct {
	Ref   string
	Value *Response
}

func (value *ResponseRef) MarshalJSON() ([]byte, error) {
	return marshalRef(value.Ref, value.Value)
}

func (value *ResponseRef) UnmarshalJSON(data []byte) error {
	return unmarshalRef(data, &value.Ref, &value.Value)
}

// https://swagger.io/specification/#responses-object
type Response struct {
	ExtensionProps
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	Headers     Headers `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     Content `json:"content,omitempty" yaml:"content,omitempty"`
	Links       Links   `json:"links,omitempty" yaml:"links,omitempty"`
}

type Links map[string]*LinkRef

type LinkRef struct {
	Ref   string
	Value *Link
}

func (value *LinkRef) MarshalJSON() ([]byte, error) {
	return marshalRef(value.Ref, value.Value)
}

func (value *LinkRef) UnmarshalJSON(data []byte) error {
	return unmarshalRef(data, &value.Ref, &value.Value)
}

// https://swagger.io/specification/#link-object
type Link struct {
	ExtensionProps
	OperationID  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Server       *Server                `json:"server,omitempty" yaml:"server,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
}

type CallbackRef struct {
	Ref   string
	Value *Callback
}

type Callbacks map[string]*CallbackRef

// https://swagger.io/specification/#callback-object
type Callback map[string]*PathItem

type Examples map[string]*ExampleRef

type ExampleRef struct {
	Ref   string
	Value *Example
}

// https://swagger.io/specification/#example-object
type Example struct {
	ExtensionProps

	Summary       string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string      `json:"description,omitempty" yaml:"description,omitempty"`
	Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

type Content map[string]*MediaType

// https://swagger.io/specification/#media-type-object
type MediaType struct {
	ExtensionProps

	Schema   *SchemaRef           `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}          `json:"example,omitempty" yaml:"example,omitempty"`
	Examples Examples             `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// https://swagger.io/specification/#encoding-object
type Encoding struct {
	ExtensionProps

	ContentType   string  `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       Headers `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string  `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool   `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool    `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// https://swagger.io/specification/#parameter-object
type Parameter struct {
	ExtensionProps
	Name            string      `json:"name,omitempty" yaml:"name,omitempty"`
	In              string      `json:"in,omitempty" yaml:"in,omitempty"`
	Description     string      `json:"description,omitempty" yaml:"description,omitempty"`
	Style           string      `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool       `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Required        bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Schema          *SchemaRef  `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{} `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        Examples    `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         Content     `json:"content,omitempty" yaml:"content,omitempty"`
}

// Internal type mapping from FQMN to descriptor.Message. Used as a set by the
// findServiceMessages function.
type messageMap map[string]*descriptor.Message

// Internal type mapping from FQEN to descriptor.Enum. Used as a set by the
// findServiceMessages function.
type enumMap map[string]*descriptor.Enum

// Internal type to store used references.
type refMap map[string]struct{}

// ExtensionProps provides support for OpenAPI extensions.
// It reads/writes all properties that begin with "x-".
type ExtensionProps struct {
	Extensions map[string]interface{} `json:"-" yaml:"-"`
}
