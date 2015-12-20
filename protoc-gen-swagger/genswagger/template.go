package genswagger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/gengo/grpc-gateway/utilities"
)

type param struct {
	*descriptor.File
	Imports []descriptor.GoPackage
}

type binding struct {
	*descriptor.Binding
}

type swaggerInfoObject struct {
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
	Get  swaggerOperationObject `json:"get,omitempty"`
	Post swaggerOperationObject `json:"post,omitempty"`
}

type swaggerOperationObject struct {
	Summary   string                 `json:"summary"`
	Responses swaggerResponsesObject `json:"response"`
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

// HasQueryParam determines if the binding needs parameters in query string.
//
// It sometimes returns true even though actually the binding does not need.
// But it is not serious because it just results in a small amount of extra codes generated.
func (b binding) HasQueryParam() bool {
	if b.Body != nil && len(b.Body.FieldPath) == 0 {
		return false
	}
	fields := make(map[string]bool)
	for _, f := range b.Method.RequestType.Fields {
		fields[f.GetName()] = true
	}
	if b.Body != nil {
		delete(fields, b.Body.FieldPath.String())
	}
	for _, p := range b.PathParams {
		delete(fields, p.FieldPath.String())
	}
	return len(fields) > 0
}

func (b binding) QueryParamFilter() queryParamFilter {
	var seqs [][]string
	if b.Body != nil {
		seqs = append(seqs, strings.Split(b.Body.FieldPath.String(), "."))
	}
	for _, p := range b.PathParams {
		seqs = append(seqs, strings.Split(p.FieldPath.String(), "."))
	}
	return queryParamFilter{utilities.NewDoubleArray(seqs)}
}

// queryParamFilter is a wrapper of utilities.DoubleArray which provides String() to output DoubleArray.Encoding in a stable and predictable format.
type queryParamFilter struct {
	*utilities.DoubleArray
}

func (f queryParamFilter) String() string {
	encodings := make([]string, len(f.Encoding))
	for str, enc := range f.Encoding {
		encodings[enc] = fmt.Sprintf("%q: %d", str, enc)
	}
	e := strings.Join(encodings, ", ")
	return fmt.Sprintf("&utilities.DoubleArray{Encoding: map[string]int{%s}, Base: %#v, Check: %#v}", e, f.Base, f.Check)
}

// This function is called with a param which contains the entire definition of a method.
func applyTemplate(p param) (string, error) {
	s := swaggerObject{
		// Swagger 2.0 is the version of this document
		Swagger:  "2.0",
		Schemes:  []string{"http", "https"},
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Paths:    make(swaggerPathsObject),
	}

	for _, svc := range p.Services {
		for _, meth := range svc.Methods {
			if meth.GetClientStreaming() || meth.GetServerStreaming() {
				return "", errors.New(`Service uses streaming, which is not currently supported. Maybe you would like to implement it? It wouldn't be that hard and we don't bite. Why don't you send a pull request to https://github.com/gengo/grpc-gateway?`)
			}
			for _, b := range meth.Bindings {
				pathItemObject := swaggerPathItemObject{}
				operationObject := swaggerOperationObject{
					Summary: fmt.Sprintf("Generated for %s.%s - %s", svc.GetName(), meth.GetName(), b.PathTmpl.Verb),
					Responses: swaggerResponsesObject{
						"default": swaggerResponseObject{
							Description: "Description",
							Schema: swaggerSchemaObject{
								Ref: fmt.Sprintf("#/definitions/%s", meth.Service.File.GetMessageType()[0].GetName()),
							},
						},
					},
				}
				switch b.HTTPMethod {
				case "GET":
					pathItemObject.Get = operationObject
					break
				case "POST":
					pathItemObject.Post = operationObject
					break
				}
				s.Paths[b.PathTmpl.Template] = pathItemObject
			}
		}
	}

	// Now we parse the definitions
	s.Definitions = swaggerDefinitionsObject{}
	for _, msg := range p.Messages {
		swaggerSchemaObject := swaggerSchemaObject{
			Properties: map[string]swaggerSchemaPropertyObject{},
		}
		for _, field := range msg.Fields {
			var fieldType, fieldFormat string
			primitive := true
			// Field type and format from http://swagger.io/specification/ in the
			// "Data Types" table
			switch field.FieldDescriptorProto.Type.String() {
			case "TYPE_DOUBLE":
				fieldType = "number"
				fieldFormat = "double"
				break
			case "TYPE_FLOAT":
				fieldType = "number"
				fieldFormat = "float"
				break
			case "TYPE_INT64":
				fieldType = "integer"
				fieldFormat = "int64"
				break
			case "TYPE_UINT64":
				fieldType = "integer"
				fieldFormat = "int64"
				break
			case "TYPE_INT32":
				fieldType = "integer"
				fieldFormat = "int32"
				break
			case "TYPE_FIXED64":
				fieldType = "integer"
				fieldFormat = "int64"
				break
			case "TYPE_FIXED32":
				fieldType = "integer"
				fieldFormat = "int32"
				break
			case "TYPE_BOOL":
				fieldType = "boolean"
				fieldFormat = "boolean"
				break
			case "TYPE_STRING":
				fieldType = "string"
				fieldFormat = "string"
				break
			case "TYPE_GROUP":
				// WTF is a group? is this sufficient?
				primitive = false
				break
			case "TYPE_MESSAGE":
				// Check in here if it is the special date/datetime proto and
				// serialize as a primitive date object
				primitive = false
				fieldType = ""
				fieldFormat = ""
				break
			case "TYPE_BYTES":
				fieldType = "string"
				fieldFormat = "byte"
				break
			case "TYPE_UINT32":
				fieldType = "integer"
				fieldFormat = "int64"
				break
			case "TYPE_ENUM":
				fieldType = "integer"
				fieldFormat = "int64"
				break
			case "TYPE_SFIXED32":
				fieldType = "integer"
				fieldFormat = "int32"
				break
			case "TYPE_SFIXED64":
				fieldType = "integer"
				fieldFormat = "int32"
				break
			case "TYPE_SINT32":
				fieldType = "integer"
				fieldFormat = "int32"
				break
			case "TYPE_SINT64":
				fieldType = "integer"
				fieldFormat = "int64"
				break
			default:
				fieldType = field.FieldDescriptorProto.Type.String()
				fieldFormat = "UNKNOWN"
			}

			if primitive {
				swaggerSchemaObject.Properties[*field.FieldDescriptorProto.Name] = swaggerSchemaPropertyObject{
					//Name:   string(*field.FieldDescriptorProto.Name),
					Type:   fieldType,
					Format: fieldFormat,
				}
			} else {
				nestedType := strings.Split(*field.FieldDescriptorProto.TypeName, ".")
				swaggerSchemaObject.Properties[*field.FieldDescriptorProto.Name] = swaggerSchemaPropertyObject{
					Ref: "#/definitions/" + nestedType[len(nestedType)-1],
				}
				//_ = "breakpoint" // Pauses the program.
			}

			//_ = "breakpoint" // Pauses the program.
			//}
		}
		s.Definitions[*msg.Name] = swaggerSchemaObject
	}

	w := bytes.NewBuffer(nil)
	enc := json.NewEncoder(w)
	enc.Encode(&s)

	return w.String(), nil
}
