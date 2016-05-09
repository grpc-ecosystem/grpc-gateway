package genswagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	pbdescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// findServicesMessagesAndEnumerations discovers all messages and enums defined in the RPC methods of the service.
func findServicesMessagesAndEnumerations(s []*descriptor.Service, reg *descriptor.Registry, m messageMap, e enumMap) {
	for _, svc := range s {
		for _, meth := range svc.Methods {
			m[fullyQualifiedNameToSwaggerName(meth.RequestType.FQMN(), reg)] = meth.RequestType
			findNestedMessagesAndEnumerations(meth.RequestType, reg, m, e)
			m[fullyQualifiedNameToSwaggerName(meth.ResponseType.FQMN(), reg)] = meth.ResponseType
			findNestedMessagesAndEnumerations(meth.ResponseType, reg, m, e)
		}
	}
}

// findNestedMessagesAndEnumerations those can be generated by the services.
func findNestedMessagesAndEnumerations(message *descriptor.Message, reg *descriptor.Registry, m messageMap, e enumMap) {
	// Iterate over all the fields that
	for _, t := range message.Fields {
		fieldType := t.GetTypeName()
		// If the type is an empty string then it is a proto primitive
		if fieldType != "" {
			if _, ok := m[fieldType]; !ok {
				msg, err := reg.LookupMsg("", fieldType)
				if err != nil {
					enum, err := reg.LookupEnum("", fieldType)
					if err != nil {
						panic(err)
					}
					e[fieldType] = enum
					continue
				}
				m[fieldType] = msg
				findNestedMessagesAndEnumerations(msg, reg, m, e)
			}
		}
	}
}

func renderMessagesAsDefinition(messages messageMap, d swaggerDefinitionsObject, reg *descriptor.Registry) {
	for _, msg := range messages {
		if opt := msg.GetOptions(); opt != nil && opt.MapEntry != nil && *opt.MapEntry {
			continue
		}
		schema := swaggerSchemaObject{
			schemaCore: schemaCore{
				Type: "object",
			},
			Properties: make(map[string]swaggerSchemaObject),
		}
		for _, f := range msg.Fields {
			schema.Properties[f.GetName()] = schemaOfField(f, reg)
		}
		d[fullyQualifiedNameToSwaggerName(msg.FQMN(), reg)] = schema
	}
}

// schemaOfField returns a swagger Schema Object for a protobuf field.
func schemaOfField(f *descriptor.Field, reg *descriptor.Registry) swaggerSchemaObject {
	const (
		singular = 0
		array    = 1
		object   = 2
	)
	var (
		core      schemaCore
		aggregate int
	)

	fd := f.FieldDescriptorProto
	if m, err := reg.LookupMsg("", f.GetTypeName()); err == nil {
		if opt := m.GetOptions(); opt != nil && opt.MapEntry != nil && *opt.MapEntry {
			fd = m.GetField()[1]
			aggregate = object
		}
	}
	if fd.GetLabel() == pbdescriptor.FieldDescriptorProto_LABEL_REPEATED {
		aggregate = array
	}

	switch ft := fd.GetType(); ft {
	case pbdescriptor.FieldDescriptorProto_TYPE_ENUM, pbdescriptor.FieldDescriptorProto_TYPE_MESSAGE, pbdescriptor.FieldDescriptorProto_TYPE_GROUP:
		core = schemaCore{
			Ref: "#/definitions/" + fullyQualifiedNameToSwaggerName(fd.GetTypeName(), reg),
		}
	default:
		ftype, format, ok := primitiveSchema(ft)
		if ok {
			core = schemaCore{Type: ftype, Format: format}
		} else {
			core = schemaCore{Type: ft.String(), Format: "UNKNOWN"}
		}
	}
	switch aggregate {
	case array:
		return swaggerSchemaObject{
			schemaCore: schemaCore{
				Type: "array",
			},
			Items: (*swaggerItemsObject)(&core),
		}
	case object:
		return swaggerSchemaObject{
			schemaCore: schemaCore{
				Type: "object",
			},
			AdditionalProperties: &swaggerSchemaObject{schemaCore: core},
		}
	default:
		return swaggerSchemaObject{schemaCore: core}
	}
}

// primitiveSchema returns a pair of "Type" and "Format" in JSON Schema for
// the given primitive field type.
// The last return parameter is true iff the field type is actually primitive.
func primitiveSchema(t pbdescriptor.FieldDescriptorProto_Type) (ftype, format string, ok bool) {
	switch t {
	case pbdescriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "number", "double", true
	case pbdescriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "number", "float", true
	case pbdescriptor.FieldDescriptorProto_TYPE_INT64:
		return "integer", "int64", true
	case pbdescriptor.FieldDescriptorProto_TYPE_UINT64:
		return "integer", "int64", true
	case pbdescriptor.FieldDescriptorProto_TYPE_INT32:
		return "integer", "int32", true
	case pbdescriptor.FieldDescriptorProto_TYPE_FIXED64:
		return "integer", "int64", true
	case pbdescriptor.FieldDescriptorProto_TYPE_FIXED32:
		return "integer", "int32", true
	case pbdescriptor.FieldDescriptorProto_TYPE_BOOL:
		return "boolean", "boolean", true
	case pbdescriptor.FieldDescriptorProto_TYPE_STRING:
		return "string", "string", true
	case pbdescriptor.FieldDescriptorProto_TYPE_BYTES:
		return "string", "byte", true
	case pbdescriptor.FieldDescriptorProto_TYPE_UINT32:
		return "integer", "int64", true
	case pbdescriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return "integer", "int32", true
	case pbdescriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return "integer", "int32", true
	case pbdescriptor.FieldDescriptorProto_TYPE_SINT32:
		return "integer", "int32", true
	case pbdescriptor.FieldDescriptorProto_TYPE_SINT64:
		return "integer", "int64", true
	default:
		return "", "", false
	}
}

// renderEnumerationsAsDefinition inserts enums into the definitions object.
func renderEnumerationsAsDefinition(enums enumMap, d swaggerDefinitionsObject, reg *descriptor.Registry) {
	for _, enum := range enums {
		var enumNames []string
		// it may be necessary to sort the result of the GetValue function.
		var defaultValue string
		for _, value := range enum.GetValue() {
			enumNames = append(enumNames, value.GetName())
			if defaultValue == "" && value.GetNumber() == 0 {
				defaultValue = value.GetName()
			}
		}

		d[fullyQualifiedNameToSwaggerName(enum.FQEN(), reg)] = swaggerSchemaObject{
			schemaCore: schemaCore{
				Type: "string",
			},
			Enum:    enumNames,
			Default: defaultValue,
		}
	}
}

// Take in a FQMN or FQEN and return a swagger safe version of the FQMN
func fullyQualifiedNameToSwaggerName(fqn string, reg *descriptor.Registry) string {
	return resolveFullyQualifiedNameToSwaggerName(fqn, append(reg.GetAllFQMNs(), reg.GetAllFQENs()...))
}

// Take the names of every proto and "uniq-ify" them. The idea is to produce a
// set of names that meet a couple of conditions. They must be stable, they
// must be unique, and they must be shorter than the FQN.
//
// This likely could be made better. This will always generate the same names
// but may not always produce optimal names. This is a reasonably close
// approximation of what they should look like in most cases.
func resolveFullyQualifiedNameToSwaggerName(fqn string, messages []string) string {
	packagesByDepth := make(map[int][][]string)
	uniqueNames := make(map[string]string)

	hierarchy := func(pkg string) []string {
		return strings.Split(pkg, ".")
	}

	for _, p := range messages {
		h := hierarchy(p)
		for depth := range h {
			if _, ok := packagesByDepth[depth]; !ok {
				packagesByDepth[depth] = make([][]string, 0)
			}
			packagesByDepth[depth] = append(packagesByDepth[depth], h[len(h)-depth:])
		}
	}

	count := func(list [][]string, item []string) int {
		i := 0
		for _, element := range list {
			if reflect.DeepEqual(element, item) {
				i++
			}
		}
		return i
	}

	for _, p := range messages {
		h := hierarchy(p)
		for depth := 0; depth < len(h); depth++ {
			if count(packagesByDepth[depth], h[len(h)-depth:]) == 1 {
				uniqueNames[p] = strings.Join(h[len(h)-depth-1:], "")
				break
			}
			if depth == len(h)-1 {
				uniqueNames[p] = strings.Join(h, "")
			}
		}
	}
	return uniqueNames[fqn]
}

// Swagger expects paths of the form /path/{string_value} but grpc-gateway paths are expected to be of the form /path/{string_value=strprefix/*}. This should reformat it correctly.
func templateToSwaggerPath(path string) string {
	// It seems like the right thing to do here is to just use
	// strings.Split(path, "/") but that breaks badly when you hit a url like
	// /{my_field=prefix/*}/ and end up with 2 sections representing my_field.
	// Instead do the right thing and write a small pushdown (counter) automata
	// for it.
	var parts []string
	depth := 0
	buffer := ""
	for _, char := range path {
		switch char {
		case '{':
			// Push on the stack
			depth += 1
			buffer += string(char)
			break
		case '}':
			if depth == 0 {
				panic("Encountered } without matching { before it.")
			}
			// Pop from the stack
			depth -= 1
			buffer += "}"
		case '/':
			if depth == 0 {
				parts = append(parts, buffer)
				buffer = ""
				// Since the stack was empty when we hit the '/' we are done with this
				// section.
				continue
			}
		default:
			buffer += string(char)
			break
		}
	}

	// Now append the last element to parts
	parts = append(parts, buffer)

	// Parts is now an array of segments of the path. Interestingly, since the
	// syntax for this subsection CAN be handled by a regexp since it has no
	// memory.
	re := regexp.MustCompile("{([a-zA-Z][a-zA-Z0-9_.]*).*}")
	for index, part := range parts {
		parts[index] = re.ReplaceAllString(part, "{$1}")
	}

	return strings.Join(parts, "/")
}

func renderServices(services []*descriptor.Service, paths swaggerPathsObject, reg *descriptor.Registry) error {
	for _, svc := range services {
		for _, meth := range svc.Methods {
			if meth.GetClientStreaming() || meth.GetServerStreaming() {
				return fmt.Errorf(`service uses streaming, which is not currently supported. Maybe you would like to implement it? It wouldn't be that hard and we don't bite. Why don't you send a pull request to https://github.com/gengo/grpc-gateway?`)
			}
			for _, b := range meth.Bindings {
				// Iterate over all the swagger parameters
				parameters := swaggerParametersObject{}
				for _, parameter := range b.PathParams {

					var paramType, paramFormat string
					switch pt := parameter.Target.GetType(); pt {
					case pbdescriptor.FieldDescriptorProto_TYPE_GROUP, pbdescriptor.FieldDescriptorProto_TYPE_MESSAGE:
						return fmt.Errorf("only primitive types are allowed in path parameters")
					case pbdescriptor.FieldDescriptorProto_TYPE_ENUM:
						paramType = fullyQualifiedNameToSwaggerName(parameter.Target.GetTypeName(), reg)
						paramFormat = ""
					default:
						var ok bool
						paramType, paramFormat, ok = primitiveSchema(pt)
						if !ok {
							return fmt.Errorf("unknown field type %v", pt)
						}
					}

					parameters = append(parameters, swaggerParameterObject{
						Name:     parameter.String(),
						In:       "path",
						Required: true,
						// Parameters in gRPC-Gateway can only be strings?
						Type:   paramType,
						Format: paramFormat,
					})
				}
				// Now check if there is a body parameter
				if b.Body != nil {
					parameters = append(parameters, swaggerParameterObject{
						Name:     "body",
						In:       "body",
						Required: true,
						Schema: &swaggerSchemaObject{
							schemaCore: schemaCore{
								Ref: fmt.Sprintf("#/definitions/%s", fullyQualifiedNameToSwaggerName(meth.RequestType.FQMN(), reg)),
							},
						},
					})
				}

				pathItemObject, ok := paths[templateToSwaggerPath(b.PathTmpl.Template)]
				if !ok {
					pathItemObject = swaggerPathItemObject{}
				}
				operationObject := &swaggerOperationObject{
					Summary:     fmt.Sprintf("%s.%s", svc.GetName(), meth.GetName()),
					Tags:        []string{svc.GetName()},
					OperationId: fmt.Sprintf("%s", meth.GetName()),
					Parameters:  parameters,
					Responses: swaggerResponsesObject{
						"default": swaggerResponseObject{
							Description: "Description",
							Schema: swaggerSchemaObject{
								schemaCore: schemaCore{
									Ref: fmt.Sprintf("#/definitions/%s", fullyQualifiedNameToSwaggerName(meth.ResponseType.FQMN(), reg)),
								},
							},
						},
					},
				}

				switch b.HTTPMethod {
				case "DELETE":
					pathItemObject.Delete = operationObject
					break
				case "GET":
					pathItemObject.Get = operationObject
					break
				case "POST":
					pathItemObject.Post = operationObject
					break
				case "PUT":
					pathItemObject.Put = operationObject
					break
				}
				paths[templateToSwaggerPath(b.PathTmpl.Template)] = pathItemObject
			}
		}
	}

	// Success! return nil on the error object
	return nil
}

// This function is called with a param which contains the entire definition of a method.
func applyTemplate(p param) (string, error) {
	// Create the basic template object. This is the object that everything is
	// defined off of.
	s := swaggerObject{
		// Swagger 2.0 is the version of this document
		Swagger:     "2.0",
		Schemes:     []string{"http", "https"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Paths:       make(swaggerPathsObject),
		Definitions: make(swaggerDefinitionsObject),
	}

	// Loops through all the services and their exposed GET/POST/PUT/DELETE definitions
	// and create entries for all of them.
	renderServices(p.Services, s.Paths, p.reg)

	// Find all the service's messages and enumerations that are defined (recursively) and then
	// write their request and response types out as definition objects.
	m := messageMap{}
	e := enumMap{}
	findServicesMessagesAndEnumerations(p.Services, p.reg, m, e)
	renderMessagesAsDefinition(m, s.Definitions, p.reg)
	renderEnumerationsAsDefinition(e, s.Definitions, p.reg)

	// We now have rendered the entire swagger object. Write the bytes out to a
	// string so it can be written to disk.
	var w bytes.Buffer
	enc := json.NewEncoder(&w)
	enc.Encode(&s)

	return w.String(), nil
}
