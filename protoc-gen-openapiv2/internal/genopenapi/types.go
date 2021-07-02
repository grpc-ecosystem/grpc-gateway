package genopenapi

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

type param struct {
	*descriptor.File
	reg *descriptor.Registry
}

// http://swagger.io/specification/#infoObject
type openapiInfoObject struct {
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
	Version        string `json:"version"`

	Contact *openapiContactObject `json:"contact,omitempty"`
	License *openapiLicenseObject `json:"license,omitempty"`

	extensions []extension
}

// https://swagger.io/specification/#tagObject
type openapiTagObject struct {
	Name         string                              `json:"name"`
	Description  string                              `json:"description,omitempty"`
	ExternalDocs *openapiExternalDocumentationObject `json:"externalDocs,omitempty"`
}

// http://swagger.io/specification/#contactObject
type openapiContactObject struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// http://swagger.io/specification/#licenseObject
type openapiLicenseObject struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// http://swagger.io/specification/#externalDocumentationObject
type openapiExternalDocumentationObject struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type extension struct {
	key   string
	value json.RawMessage
}

// http://swagger.io/specification/#swaggerObject
type openapiSwaggerObject struct {
	Swagger             string                              `json:"swagger"`
	Info                openapiInfoObject                   `json:"info"`
	Tags                []openapiTagObject                  `json:"tags,omitempty"`
	Host                string                              `json:"host,omitempty"`
	BasePath            string                              `json:"basePath,omitempty"`
	Schemes             []string                            `json:"schemes,omitempty"`
	Consumes            []string                            `json:"consumes"`
	Produces            []string                            `json:"produces"`
	Paths               openapiPathsObject                  `json:"paths"`
	Definitions         openapiDefinitionsObject            `json:"definitions"`
	SecurityDefinitions openapiSecurityDefinitionsObject    `json:"securityDefinitions,omitempty"`
	Security            []openapiSecurityRequirementObject  `json:"security,omitempty"`
	ExternalDocs        *openapiExternalDocumentationObject `json:"externalDocs,omitempty"`

	extensions []extension
}

// http://swagger.io/specification/#securityDefinitionsObject
type openapiSecurityDefinitionsObject map[string]openapiSecuritySchemeObject

// http://swagger.io/specification/#securitySchemeObject
type openapiSecuritySchemeObject struct {
	Type             string              `json:"type"`
	Description      string              `json:"description,omitempty"`
	Name             string              `json:"name,omitempty"`
	In               string              `json:"in,omitempty"`
	Flow             string              `json:"flow,omitempty"`
	AuthorizationURL string              `json:"authorizationUrl,omitempty"`
	TokenURL         string              `json:"tokenUrl,omitempty"`
	Scopes           openapiScopesObject `json:"scopes,omitempty"`

	extensions []extension
}

// http://swagger.io/specification/#scopesObject
type openapiScopesObject map[string]string

// http://swagger.io/specification/#securityRequirementObject
type openapiSecurityRequirementObject map[string][]string

// http://swagger.io/specification/#pathsObject
type openapiPathsObject map[string]openapiPathItemObject

// http://swagger.io/specification/#pathItemObject
type openapiPathItemObject struct {
	Get    *openapiOperationObject `json:"get,omitempty"`
	Delete *openapiOperationObject `json:"delete,omitempty"`
	Post   *openapiOperationObject `json:"post,omitempty"`
	Put    *openapiOperationObject `json:"put,omitempty"`
	Patch  *openapiOperationObject `json:"patch,omitempty"`
}

// http://swagger.io/specification/#operationObject
type openapiOperationObject struct {
	Summary     string                  `json:"summary,omitempty"`
	Description string                  `json:"description,omitempty"`
	OperationID string                  `json:"operationId"`
	Responses   openapiResponsesObject  `json:"responses"`
	Parameters  openapiParametersObject `json:"parameters,omitempty"`
	Tags        []string                `json:"tags,omitempty"`
	Deprecated  bool                    `json:"deprecated,omitempty"`
	Produces    []string                `json:"produces,omitempty"`

	Security     *[]openapiSecurityRequirementObject `json:"security,omitempty"`
	ExternalDocs *openapiExternalDocumentationObject `json:"externalDocs,omitempty"`

	extensions []extension
}

type openapiParametersObject []openapiParameterObject

// http://swagger.io/specification/#parameterObject
type openapiParameterObject struct {
	Name             string              `json:"name"`
	Description      string              `json:"description,omitempty"`
	In               string              `json:"in,omitempty"`
	Required         bool                `json:"required"`
	Type             string              `json:"type,omitempty"`
	Format           string              `json:"format,omitempty"`
	Items            *openapiItemsObject `json:"items,omitempty"`
	Enum             []string            `json:"enum,omitempty"`
	CollectionFormat string              `json:"collectionFormat,omitempty"`
	Default          string              `json:"default,omitempty"`
	MinItems         *int                `json:"minItems,omitempty"`

	// Or you can explicitly refer to another type. If this is defined all
	// other fields should be empty
	Schema *openapiSchemaObject `json:"schema,omitempty"`
}

// core part of schema, which is common to itemsObject and schemaObject.
// http://swagger.io/specification/v2/#itemsObject
// The OAS3 spec (https://swagger.io/specification/#schemaObject) defines the
// `nullable` field as part of a Schema Object. This behavior has been
// "back-ported" to OAS2 as the Specification Extension `x-nullable`, and is
// supported by generation tools such as swagger-codegen and go-swagger.
// For protoc-gen-openapiv3, we'd want to add `nullable` instead.
type schemaCore struct {
	Type      string          `json:"type,omitempty"`
	Format    string          `json:"format,omitempty"`
	Ref       string          `json:"$ref,omitempty"`
	XNullable bool            `json:"x-nullable,omitempty"`
	Example   json.RawMessage `json:"example,omitempty"`

	Items *openapiItemsObject `json:"items,omitempty"`

	// If the item is an enumeration include a list of all the *NAMES* of the
	// enum values.  I'm not sure how well this will work but assuming all enums
	// start from 0 index it will be great. I don't think that is a good assumption.
	Enum    []string `json:"enum,omitempty"`
	Default string   `json:"default,omitempty"`
}

func (s *schemaCore) setRefFromFQN(ref string, reg *descriptor.Registry) error {
	name, ok := fullyQualifiedNameToOpenAPIName(ref, reg)
	if !ok {
		return fmt.Errorf("setRefFromFQN: can't resolve OpenAPI name from '%v'", ref)
	}
	s.Ref = fmt.Sprintf("#/definitions/%s", name)
	return nil
}

type openapiItemsObject schemaCore

// http://swagger.io/specification/#responsesObject
type openapiResponsesObject map[string]openapiResponseObject

// http://swagger.io/specification/#responseObject
type openapiResponseObject struct {
	Description string                 `json:"description"`
	Schema      openapiSchemaObject    `json:"schema"`
	Examples    map[string]interface{} `json:"examples,omitempty"`
	Headers     openapiHeadersObject   `json:"headers,omitempty"`

	extensions []extension
}

type openapiHeadersObject map[string]openapiHeaderObject

// http://swagger.io/specification/#headerObject
type openapiHeaderObject struct {
	Description string          `json:"description,omitempty"`
	Type        string          `json:"type,omitempty"`
	Format      string          `json:"format,omitempty"`
	Default     json.RawMessage `json:"default,omitempty"`
	Pattern     string          `json:"pattern,omitempty"`
}

type keyVal struct {
	Key   string
	Value interface{}
}

type openapiSchemaObjectProperties []keyVal

func (op openapiSchemaObjectProperties) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, kv := range op {
		if i != 0 {
			buf.WriteString(",")
		}
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")
		val, err := json.Marshal(kv.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

// http://swagger.io/specification/#schemaObject
type openapiSchemaObject struct {
	schemaCore
	// Properties can be recursively defined
	Properties           *openapiSchemaObjectProperties `json:"properties,omitempty"`
	AdditionalProperties *openapiSchemaObject           `json:"additionalProperties,omitempty"`

	Description string `json:"description,omitempty"`
	Title       string `json:"title,omitempty"`

	ExternalDocs *openapiExternalDocumentationObject `json:"externalDocs,omitempty"`

	ReadOnly         bool     `json:"readOnly,omitempty"`
	MultipleOf       float64  `json:"multipleOf,omitempty"`
	Maximum          float64  `json:"maximum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
	Minimum          float64  `json:"minimum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`
	MaxLength        uint64   `json:"maxLength,omitempty"`
	MinLength        uint64   `json:"minLength,omitempty"`
	Pattern          string   `json:"pattern,omitempty"`
	MaxItems         uint64   `json:"maxItems,omitempty"`
	MinItems         uint64   `json:"minItems,omitempty"`
	UniqueItems      bool     `json:"uniqueItems,omitempty"`
	MaxProperties    uint64   `json:"maxProperties,omitempty"`
	MinProperties    uint64   `json:"minProperties,omitempty"`
	Required         []string `json:"required,omitempty"`
}

// http://swagger.io/specification/#definitionsObject
type openapiDefinitionsObject map[string]openapiSchemaObject

// Internal type mapping from FQMN to descriptor.Message. Used as a set by the
// findServiceMessages function.
type messageMap map[string]*descriptor.Message

// Internal type mapping from FQEN to descriptor.Enum. Used as a set by the
// findServiceMessages function.
type enumMap map[string]*descriptor.Enum

// Internal type to store used references.
type refMap map[string]struct{}
