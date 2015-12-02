package genswagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/gengo/grpc-gateway/protoc-gen-swagger/descriptor"
	"github.com/gengo/grpc-gateway/utilities"
	"github.com/golang/glog"
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
	glog.Infof("Applying template")
	//if err := headerTemplate.Execute(w, p); err != nil {
	//	return "", err
	//}
	//var methodSeen bool
	//for _, svc := range p.Services {
	//	for _, meth := range svc.Methods {
	//		glog.V(2).Infof("Processing %s.%s", svc.GetName(), meth.GetName())
	//		methodSeen = true
	//		for _, b := range meth.Bindings {
	//			if err := handlerTemplate.Execute(w, binding{Binding: b}); err != nil {
	//				return "", err
	//			}
	//		}
	//	}
	//}
	//if !methodSeen {
	//	return "", errNoTargetService
	//}
	//if err := trailerTemplate.Execute(w, p.Services); err != nil {
	//	return "", err
	//}

	//return string(""), nil
	return w.String(), nil
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`
// Code generated by protoc-gen-swagger
// source: {{.GetName}}
// DO NOT EDIT!

/*
Package {{.GoPkg.Name}} is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package {{.GoPkg.Name}}
import (
	{{range $i := .Imports}}{{if $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}

	{{range $i := .Imports}}{{if not $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}
)

var _ codes.Code
var _ io.Reader
var _ = runtime.String
var _ = json.Marshal
var _ = utilities.PascalFromSnake
`))

	handlerTemplate = template.Must(template.New("handler").Parse(`
{{if .Method.GetClientStreaming}}
{{template "client-streaming-request-func" .}}
{{else}}
{{template "client-rpc-request-func" .}}
{{end}}
`))

	_ = template.Must(handlerTemplate.New("request-func-signature").Parse(strings.Replace(`
{{if .Method.GetServerStreaming}}
func request_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}(ctx context.Context, client {{.Method.Service.GetName}}Client, req *http.Request, pathParams map[string]string) ({{.Method.Service.GetName}}_{{.Method.GetName}}Client, error)
{{else}}
func request_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}(ctx context.Context, client {{.Method.Service.GetName}}Client, req *http.Request, pathParams map[string]string) (proto.Message, error)
{{end}}`, "\n", "", -1)))

	_ = template.Must(handlerTemplate.New("client-streaming-request-func").Parse(`
{{template "request-func-signature" .}} {
	stream, err := client.{{.Method.GetName}}(ctx)
	if err != nil {
		glog.Errorf("Failed to start streaming: %v", err)
		return nil, err
	}
	dec := json.NewDecoder(req.Body)
	for {
		var protoReq {{.Method.RequestType.GoType .Method.Service.File.GoPkg.Path}}
		err = dec.Decode(&protoReq)
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Errorf("Failed to decode request: %v", err)
			return nil, grpc.Errorf(codes.InvalidArgument, "%v", err)
		}
		if err = stream.Send(&protoReq); err != nil {
			glog.Errorf("Failed to send request: %v", err)
			return nil, err
		}
	}
{{if .Method.GetServerStreaming}}
	if err = stream.CloseSend(); err != nil {
		glog.Errorf("Failed to terminate client stream: %v", err)
		return nil, err
	}
	return stream, nil
{{else}}
	return stream.CloseAndRecv()
{{end}}
}
`))

	_ = template.Must(handlerTemplate.New("client-rpc-request-func").Parse(`
{{if .HasQueryParam}}
var (
	filter_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}} = {{.QueryParamFilter}}
)
{{end}}
{{template "request-func-signature" .}} {
	var protoReq {{.Method.RequestType.GoType .Method.Service.File.GoPkg.Path}}
{{if .Body}}
	if err := json.NewDecoder(req.Body).Decode(&{{.Body.RHS "protoReq"}}); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}
{{end}}
{{if .PathParams}}
	var (
		val string
		ok bool
		err error
		_ = err
	)
	{{range $param := .PathParams}}
	val, ok = pathParams[{{$param | printf "%q"}}]
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing parameter %s", {{$param | printf "%q"}})
	}
{{if $param.IsNestedProto3 }}
	err = runtime.PopulateFieldFromPath(&protoReq, {{$param | printf "%q"}}, val)
{{else}}
	{{$param.RHS "protoReq"}}, err = {{$param.ConvertFuncExpr}}(val)
{{end}}
	if err != nil {
		return nil, err
	}
	{{end}}
{{end}}
{{if .HasQueryParam}}
	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}
{{end}}

	return client.{{.Method.GetName}}(ctx, &protoReq)
}`))

	trailerTemplate = template.Must(template.New("trailer").Parse(`
{{range $svc := .}}
// Register{{$svc.GetName}}HandlerFromEndpoint is same as Register{{$svc.GetName}}Handler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func Register{{$svc.GetName}}HandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string) (err error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				glog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				glog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return Register{{$svc.GetName}}Handler(ctx, mux, conn)
}

// Register{{$svc.GetName}}Handler registers the http handlers for service {{$svc.GetName}} to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func Register{{$svc.GetName}}Handler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	client := New{{$svc.GetName}}Client(conn)
	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, pattern_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		resp, err := request_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(runtime.AnnotateContext(ctx, req), client, req, pathParams)
		if err != nil {
			runtime.HTTPError(ctx, w, err)
			return
		}
		{{if $m.GetServerStreaming}}
		forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(ctx, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)
		{{else}}
		forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(ctx, w, req, resp, mux.GetForwardResponseOptions()...)
		{{end}}
	})
	{{end}}
	{{end}}
	return nil
}

var (
	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}
	pattern_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}} = runtime.MustPattern(runtime.NewPattern({{$b.PathTmpl.Version}}, {{$b.PathTmpl.OpCodes | printf "%#v"}}, {{$b.PathTmpl.Pool | printf "%#v"}}, {{$b.PathTmpl.Verb | printf "%q"}}))
	{{end}}
	{{end}}
)

var (
	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}
	forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}} = {{if $m.GetServerStreaming}}runtime.ForwardResponseStream{{else}}runtime.ForwardResponseMessage{{end}}
	{{end}}
	{{end}}
)
{{end}}`))
)
