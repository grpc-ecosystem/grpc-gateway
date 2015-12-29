package genswagger

import (
	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

type param struct {
	*descriptor.File
	reg *descriptor.Registry
}

type binding struct {
	*descriptor.Binding
}

type swaggerInfoObject struct {
	Version string `json:"version"`
	Title   string `json:"title"`
}

type swaggerObject struct {
	Swagger     string                   `json:"swagger"`
	Info        swaggerInfoObject        `json:"info"`
	Host        string                   `json:"host,omitempty"`
	BasePath    string                   `json:"basePath,omitempty"`
	Schemes     []string                 `json:"schemes"`
	Consumes    []string                 `json:"consumes"`
	Produces    []string                 `json:"produces"`
	Paths       swaggerPathsObject       `json:"paths"`
	Definitions swaggerDefinitionsObject `json:"definitions"`
}

type swaggerPathsObject map[string]swaggerPathItemObject
type swaggerPathItemObject struct {
	Get    *swaggerOperationObject `json:"get,omitempty"`
	Delete *swaggerOperationObject `json:"delete,omitempty"`
	Post   *swaggerOperationObject `json:"post,omitempty"`
	Put    *swaggerOperationObject `json:"put,omitempty"`
}

type swaggerOperationObject struct {
	Summary    string                  `json:"summary"`
	Responses  swaggerResponsesObject  `json:"responses"`
	Parameters swaggerParametersObject `json:"parameters,omitempty"`
}

type swaggerParametersObject []swaggerParameterObject
type swaggerParameterObject struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	In          string             `json:"in,omitempty"`
	Required    bool               `json:"required"`
	Type        string             `json:"type,omitempty"`
	Format      string             `json:"format,omitempty"`
	Items       swaggerItemsObject `json:"items,omitempty"`

	// Or you can explicitly refer to another type. If this is defined all
	// other fields should be empty
	Schema *swaggerSchemaObject `json:"schema,omitempty"`
}

type swaggerItemsObject struct {
	Ref string `json:"$ref,omitempty"`
}

type swaggerResponsesObject map[string]swaggerResponseObject

type swaggerResponseObject struct {
	Description string              `json:"description"`
	Schema      swaggerSchemaObject `json:"schema"`
}

type swaggerSchemaObject struct {
	Ref        string                                 `json:"$ref,omitempty"`
	Type       string                                 `json:"type,omitempty"`
	Properties map[string]swaggerSchemaPropertyObject `json:"properties,omitempty"`
}

type swaggerSchemaPropertyObject struct {
	//Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`

	// Or you can explicitly refer to another type. If this is defined all
	// other fields should be empty
	Ref string `json:"$ref,omitempty"`
}

type swaggerReferenceObject struct {
	Ref string `json:"$ref"`
}

type swaggerDefinitionsObject map[string]swaggerSchemaObject

type messageMap map[string]*descriptor.Message
