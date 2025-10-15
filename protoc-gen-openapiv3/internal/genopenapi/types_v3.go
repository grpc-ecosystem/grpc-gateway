package genopenapi

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/casing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

// OpenAPI v3 root document and all nested types with both JSON and YAML annotations.

type param struct {
	*descriptor.File
	reg *descriptor.Registry
}
type RawExample json.RawMessage

func (m RawExample) MarshalJSON() ([]byte, error) {
	return (json.RawMessage)(m).MarshalJSON()
}

func (m *RawExample) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(m).UnmarshalJSON(data)
}

// MarshalYAML implements yaml.Marshaler interface.
//
// It converts RawExample to one of yaml-supported types and returns it.
//
// From yaml.Marshaler docs: The Marshaler interface may be implemented
// by types to customize their behavior when being marshaled into a YAML
// document. The returned value is marshaled in place of the original
// value implementing Marshaler.
func (e RawExample) MarshalYAML() (interface{}, error) {
	// From docs, json.Unmarshal will store one of next types to data:
	// - bool, for JSON booleans;
	// - float64, for JSON numbers;
	// - string, for JSON strings;
	// - []interface{}, for JSON arrays;
	// - map[string]interface{}, for JSON objects;
	// - nil for JSON null.
	var data interface{}
	if err := json.Unmarshal(e, &data); err != nil {
		return nil, err
	}

	return data, nil
}

type OpenAPIV3Extensions = map[string]interface{}

type OpenAPIV3Document struct {
	OpenAPI      string                 `json:"openapi" yaml:"openapi"`
	Info         *OpenAPIV3Info         `json:"info" yaml:"info"`
	Servers      []OpenAPIV3Server      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        OpenAPIV3Paths         `json:"paths" yaml:"paths"`
	Components   *OpenAPIV3Components   `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []OpenAPIV3SecurityReq `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []OpenAPIV3Tag         `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *OpenAPIV3ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type OpenAPIV3Info struct {
	Title          string            `json:"title" yaml:"title"`
	Description    string            `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string            `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *OpenAPIV3Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *OpenAPIV3License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string            `json:"version" yaml:"version"`
}

type OpenAPIV3Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

type OpenAPIV3License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type OpenAPIV3Server struct {
	URL         string                        `json:"url" yaml:"url"`
	Description string                        `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]OpenAPIV3ServerVar `json:"variables,omitempty" yaml:"variables,omitempty"`
}

type OpenAPIV3ServerVar struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

type OpenAPIV3Paths map[string]*OpenAPIV3PathItem

type OpenAPIV3PathItem struct {
	Ref         string                  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string                  `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                  `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *OpenAPIV3Operation     `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *OpenAPIV3Operation     `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *OpenAPIV3Operation     `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *OpenAPIV3Operation     `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *OpenAPIV3Operation     `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *OpenAPIV3Operation     `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *OpenAPIV3Operation     `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *OpenAPIV3Operation     `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []OpenAPIV3Server       `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []OpenAPIV3ParameterRef `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type OpenAPIV3Operation struct {
	Tags                []string                        `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary             string                          `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description         string                          `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID         string                          `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters          []OpenAPIV3ParameterRef         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody         *OpenAPIV3RequestBodyRef        `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses           OpenAPIV3Responses              `json:"responses" yaml:"responses"`
	Callbacks           map[string]OpenAPIV3CallbackRef `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated          bool                            `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security            []OpenAPIV3SecurityReq          `json:"security,omitempty" yaml:"security,omitempty"`
	Servers             []OpenAPIV3Server               `json:"servers,omitempty" yaml:"servers,omitempty"`
	ExternalDocs        *OpenAPIV3ExternalDocs          `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OpenAPIV3Extensions `json:"-" yaml:"-"`
}

func (op OpenAPIV3Operation) MarshalJSON() ([]byte, error) {
	type Alias OpenAPIV3Operation
	// Marshal the operation without extensions
	b, err := json.Marshal(Alias(op))
	if err != nil {
		return nil, err
	}
	// Unmarshal into a map to add extensions
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	// Add extensions (assuming Extensions is a map[string]interface{})
	for k, v := range op.OpenAPIV3Extensions {
		m[k] = v
	}
	return json.Marshal(m)
}

type OpenAPIV3ParameterRef struct {
	Ref                 string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*OpenAPIV3Parameter `json:",omitempty" yaml:",omitempty"`
}

type OpenAPIV3Parameter struct {
	Name            string                         `json:"name" yaml:"name"`
	In              string                         `json:"in" yaml:"in"`
	Description     string                         `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                           `json:"required" yaml:"required"`
	Deprecated      bool                           `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                           `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                         `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool                          `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                           `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *OpenAPIV3SchemaRef            `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}                    `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]OpenAPIV3ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]OpenAPIV3MediaType  `json:"content,omitempty" yaml:"content,omitempty"`
}

type OpenAPIV3RequestBodyRef struct {
	Ref                   string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*OpenAPIV3RequestBody `json:",omitempty" yaml:",omitempty"`
}

type OpenAPIV3RequestBody struct {
	Description string                        `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]OpenAPIV3MediaType `json:"content" yaml:"content"`
	Required    bool                          `json:"required,omitempty" yaml:"required,omitempty"`
}

type OpenAPIV3MediaType struct {
	Schema   *OpenAPIV3SchemaRef            `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}                    `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]OpenAPIV3ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]OpenAPIV3Encoding   `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

type OpenAPIV3Encoding struct {
	ContentType   string                        `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]OpenAPIV3HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                        `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool                         `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                          `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

type OpenAPIV3Responses map[string]OpenAPIV3ResponseRef

type OpenAPIV3ResponseRef struct {
	Ref                string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*OpenAPIV3Response `json:",omitempty" yaml:",omitempty"`
}

type OpenAPIV3Response struct {
	Description string                        `json:"description" yaml:"description"`
	Headers     map[string]OpenAPIV3HeaderRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]OpenAPIV3MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]OpenAPIV3LinkRef   `json:"links,omitempty" yaml:"links,omitempty"`
}

type OpenAPIV3HeaderRef struct {
	Ref    string           `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Header *OpenAPIV3Header `json:"-" yaml:"-"`
}

type OpenAPIV3Header struct {
	Description     string                         `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                           `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                           `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                           `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                         `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool                          `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                           `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *OpenAPIV3SchemaRef            `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}                    `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]OpenAPIV3ExampleRef `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]OpenAPIV3MediaType  `json:"content,omitempty" yaml:"content,omitempty"`
}

type OpenAPIV3SchemaRef struct {
	Ref              string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*OpenAPIV3Schema `json:",omitempty" yaml:",omitempty"`
}

type OpenAPIV3Schema struct {
	Title                string                         `json:"title,omitempty" yaml:"title,omitempty"`
	MultipleOf           float64                        `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              float64                        `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool                           `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum              float64                        `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool                           `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength            uint64                         `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            uint64                         `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern              string                         `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             uint64                         `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             uint64                         `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          bool                           `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        uint64                         `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        uint64                         `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string                       `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []string                       `json:"enum,omitempty" yaml:"enum,omitempty"`
	Type                 string                         `json:"type,omitempty" yaml:"type,omitempty"`
	AllOf                []*OpenAPIV3SchemaRef          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []*OpenAPIV3SchemaRef          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []*OpenAPIV3SchemaRef          `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Not                  *OpenAPIV3SchemaRef            `json:"not,omitempty" yaml:"not,omitempty"`
	Items                *OpenAPIV3SchemaRef            `json:"items,omitempty" yaml:"items,omitempty"`
	Properties           map[string]*OpenAPIV3SchemaRef `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties interface{}                    `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Description          string                         `json:"description,omitempty" yaml:"description,omitempty"`
	Format               string                         `json:"format,omitempty" yaml:"format,omitempty"`
	Default              interface{}                    `json:"default,omitempty" yaml:"default,omitempty"`
	Nullable             bool                           `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Discriminator        *OpenAPIV3Discriminator        `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	ReadOnly             bool                           `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly            bool                           `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	Xml                  *OpenAPIV3XML                  `json:"xml,omitempty" yaml:"xml,omitempty"`
	ExternalDocs         *OpenAPIV3ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Example              RawExample                     `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated           bool                           `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	OpenAPIV3Extensions  `json:"-" yaml:"-"`
}

func (s *OpenAPIV3SchemaRef) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal(nil)
	}
	if s.Ref != "" {
		return json.Marshal(map[string]string{"$ref": s.Ref})
	}
	if s.OpenAPIV3Schema == nil {
		return json.Marshal(nil)
	}
	schema := *s.OpenAPIV3Schema
	schema.CamelCase()
	return json.Marshal(schema)
}

func (s *OpenAPIV3Schema) CamelCase() {
	if s == nil {
		return
	}
	if s.OneOf != nil {
		for _, schema := range s.OneOf {
			schema.CamelCase()
		}
		return
	}
	newProperties := make(map[string]*OpenAPIV3SchemaRef)
	newRequiredFields := make([]string, 0, len(s.Required))
	for _, requiredField := range s.Required {
		newRequiredFields = append(newRequiredFields, casing.JSONCamelCase(requiredField))
	}
	for k, v := range s.Properties {
		if v != nil {
			// Recursively call CamelCase on the nested schema
			v.CamelCase()
		}
		newProperties[casing.JSONCamelCase(k)] = v
	}
	if s.Items != nil {
		s.Items.CamelCase()
	}
	s.Properties = newProperties
	s.Required = newRequiredFields
}

type OpenAPIV3Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

type OpenAPIV3XML struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

type OpenAPIV3ExampleRef struct {
	Ref               string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*OpenAPIV3Example `json:",omitempty" yaml:",omitempty"`
}

type OpenAPIV3Example struct {
	Summary       string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string      `json:"description,omitempty" yaml:"description,omitempty"`
	Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

type OpenAPIV3LinkRef struct {
	Ref  string         `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Link *OpenAPIV3Link `json:"-" yaml:"-"`
}

type OpenAPIV3Link struct {
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *OpenAPIV3Server       `json:"server,omitempty" yaml:"server,omitempty"`
}

type OpenAPIV3CallbackRef struct {
	Ref      string         `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Callback OpenAPIV3Paths `json:"-" yaml:"-"`
}

type OpenAPIV3Components struct {
	Schemas         map[string]*OpenAPIV3SchemaRef        `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]OpenAPIV3ResponseRef       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]OpenAPIV3ParameterRef      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]OpenAPIV3ExampleRef        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]OpenAPIV3RequestBodyRef    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]OpenAPIV3HeaderRef         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]OpenAPIV3SecuritySchemeRef `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]OpenAPIV3LinkRef           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]OpenAPIV3CallbackRef       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

type OpenAPIV3SecurityReq map[string][]string

type OpenAPIV3SecuritySchemeRef struct {
	Ref            string                   `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	SecurityScheme *OpenAPIV3SecurityScheme `json:"-" yaml:"-"`
}

type OpenAPIV3SecurityScheme struct {
	Type             string               `json:"type" yaml:"type"`
	Description      string               `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string               `json:"name,omitempty" yaml:"name,omitempty"`
	In               string               `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string               `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string               `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OpenAPIV3OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string               `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

type OpenAPIV3OAuthFlows struct {
	Implicit          *OpenAPIV3OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OpenAPIV3OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OpenAPIV3OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OpenAPIV3OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

type OpenAPIV3OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"`
}

type OpenAPIV3Tag struct {
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *OpenAPIV3ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type OpenAPIV3ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}
