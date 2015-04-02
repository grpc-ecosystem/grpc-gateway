package main

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gengo/grpc-gateway/options"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var (
	msgTbl = make(map[string]*descriptor.DescriptorProto)
)

func registerMsg(location string, msgs []*descriptor.DescriptorProto) {
	for _, m := range msgs {
		name := fmt.Sprintf("%s.%s", location, m.GetName())
		msgTbl[name] = m
		registerMsg(name, m.GetNestedType())
	}
}

func lookupMsg(name string) (*descriptor.DescriptorProto, error) {
	m, ok := msgTbl[name]
	if !ok {
		return nil, fmt.Errorf("no such message: %s", name)
	}
	return m, nil
}

func generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	targets := make(map[string]bool)
	for _, fname := range req.GetFileToGenerate() {
		targets[fname] = true
	}
	for _, file := range req.GetProtoFile() {
		pkg := file.GetPackage()
		registerMsg(pkg, file.GetMessageType())
	}

	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range req.GetProtoFile() {
		if !targets[file.GetName()] {
			glog.V(1).Infof("Skip non-target file: %s", file.GetName())
			continue
		}
		glog.V(1).Infof("Processing %s", file.GetName())
		code, err := generateSingleFile(file)
		if err != nil {
			return &plugin.CodeGeneratorResponse{
				Error: proto.String(err.Error()),
			}
		}
		name := file.GetName()
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s.pb.gw.go", base)),
			Content: proto.String(code),
		})
	}
	return &plugin.CodeGeneratorResponse{
		File: files,
	}
}

func goPackage(d *descriptor.FileDescriptorProto) string {
	if d.Options != nil && d.Options.GoPackage != nil {
		return d.Options.GetGoPackage()
	}
	if d.Package == nil {
		return filepath.Base(d.GetName())
	}
	return strings.NewReplacer("-", "_", ".", "_").Replace(d.GetPackage())
}

var (
	upperPattern = regexp.MustCompile("[A-Z]")
)

func toSnake(str string) string {
	str = upperPattern.ReplaceAllStringFunc(str, func(c string) string {
		return "_" + strings.ToLower(c)
	})
	return strings.TrimPrefix(str, "_")
}

func toCamel(str string) string {
	var components []string
	for _, c := range strings.Split(str, "_") {
		components = append(components, strings.Title(c))
	}
	return strings.Join(components, "")
}

func getAPIOptions(meth *descriptor.MethodDescriptorProto) (*options.ApiMethodOptions, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, options.E_ApiMethodOptions_ApiOptions) {
		return nil, nil
	}
	ext, err := proto.GetExtension(meth.Options, options.E_ApiMethodOptions_ApiOptions)
	if err != nil {
		return nil, err
	}
	opts, ok := ext.(*options.ApiMethodOptions)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an ApiMethodOptions", ext)
	}
	return opts, nil
}

var (
	convertFuncs = map[descriptor.FieldDescriptorProto_Type]string{
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "Float64",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:   "Float32",
		descriptor.FieldDescriptorProto_TYPE_INT64:   "Int64",
		descriptor.FieldDescriptorProto_TYPE_UINT64:  "Uint64",
		descriptor.FieldDescriptorProto_TYPE_INT32:   "Int32",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "Uint64",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "Uint32",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "Bool",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "String",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(yugui) Handle bytes
		descriptor.FieldDescriptorProto_TYPE_UINT32: "Uint32",
		// FieldDescriptorProto_TYPE_ENUM
		// TODO(yugui) Handle Enum
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "Int32",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "Int64",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "Int32",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "Int64",
	}
)

func pathParams(msg *descriptor.DescriptorProto, opts *options.ApiMethodOptions) (prefix string, params []pathParamDesc, err error) {
	components := strings.Split(opts.GetPath(), "/")
	var firstParam int
	for i, c := range components {
		if !strings.HasPrefix(c, ":") {
			if firstParam > 0 {
				// TODO(yugui) Relax this restriction?
				return "", nil, fmt.Errorf("fixed component after positional parameter: %q", opts.Path)
			}
			continue
		}
		if firstParam == 0 {
			firstParam = i
		}
		name := strings.TrimPrefix(c, ":")
		fd := lookupField(msg, name)
		if fd == nil {
			return "", nil, fmt.Errorf("field %q not found in %s", name, msg.GetName())
		}
		conv, ok := convertFuncs[fd.GetType()]
		if !ok {
			return "", nil, fmt.Errorf("unsupported path parameter type %s in %s", fd.GetType(), msg.GetName())
		}
		params = append(params, pathParamDesc{
			ProtoName:   name,
			PathIndex:   i,
			ConvertFunc: conv,
		})
	}
	if firstParam == 0 {
		return opts.GetPath(), nil, nil
	}
	return strings.Join(components[:firstParam], "/"), params, nil
}

func lookupField(msg *descriptor.DescriptorProto, name string) *descriptor.FieldDescriptorProto {
	for _, f := range msg.GetField() {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func generateSingleFile(file *descriptor.FileDescriptorProto) (string, error) {
	pkg := file.GetPackage()
	buf := bytes.NewBuffer(nil)
	if err := headerTemplate.Execute(buf, pkg); err != nil {
		return "", err
	}
	var svcDescs []serviceDesc
	for _, svc := range file.GetService() {
		sd := serviceDesc{
			Name: svc.GetName(),
		}
		for _, meth := range svc.GetMethod() {
			opts, err := getAPIOptions(meth)
			if err != nil {
				glog.Errorf("Failed to extract ApiMethodOptions: %v", err)
				return "", err
			}
			input, err := lookupMsg(meth.GetInputType())
			if err != nil {
				return "", err
			}
			fields := make(map[string]bool)
			for _, f := range input.Field {
				fields[f.GetName()] = true
			}
			prefix, params, err := pathParams(input, opts)
			if err != nil {
				return "", err
			}
			for _, p := range params {
				delete(fields, p.ProtoName)
			}
			needsBody := len(fields) != 0
			if needsBody && (opts.GetMethod() == "GET" || opts.GetMethod() == "DELETE") {
				return "", fmt.Errorf("needs request body even though http method is %s: %s", opts.Method, meth.GetName())
			}
			md := methodDesc{
				ServiceName: svc.GetName(),
				Name:        meth.GetName(),
				Method:      opts.GetMethod(),
				RequestType: input.GetName(),
				PathParams:  params,
				Prefix:      prefix,
				NeedsBody:   needsBody,
			}
			sd.Methods = append(sd.Methods, md)
			if err = handlerTemplate.Execute(buf, md); err != nil {
				return "", err
			}
		}
		if len(sd.Methods) == 0 {
			continue
		}
		svcDescs = append(svcDescs, sd)
	}
	if err := trailerTemplate.Execute(buf, svcDescs); err != nil {
		return "", err
	}

	code, err := format.Source(buf.Bytes())
	if err != nil {
		glog.Errorf("Failed to gofmt: %s: %v", buf.String(), err)
		return "", err
	}
	return string(code), nil
}

type pathParamDesc struct {
	ProtoName   string
	PathIndex   int
	ConvertFunc string
}

type queryParamDesc struct {
	ProtoName   string
	ConvertFunc string
}

type methodDesc struct {
	ServiceName string
	Name        string
	Method      string
	RequestType string
	QueryParams []queryParamDesc
	PathParams  []pathParamDesc
	Prefix      string
	NeedsBody   bool
}

type serviceDesc struct {
	Name    string
	Methods []methodDesc
}

func (d serviceDesc) EndpointFlag() string {
	return fmt.Sprintf("%s_endpoint", toSnake(d.Name))
}

func (d methodDesc) NeedsPathParam() bool {
	return len(d.PathParams) != 0
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`
package {{.}}
import (
	"net/http"
	"encoding/json"

	"google.golang.org/grpc"
	"github.com/golang/protobuf/proto"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)
`))
	handlerTemplate = template.Must(template.New("handler").Parse(`
func handle_{{.ServiceName}}_{{.Name}}(ctx context.Context, c *{{.ServiceName}}Client, req *http.Request) (proto.Message, error) {
	protoReq := new({{.RequestType}})
{{if .NeedsBody}}
	if err = json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return err
	}
{{end}}
	{{range $desc := .QueryParams}}
	protoReq.{{$desc.ProtoName}} = proto.{{$desc.ConvertFunc}}(req.FormValue({{$desc.ProtoName | printf "%q"}}))
	{{end}}
{{if .NeedsPathParam}}
	components := strings.Split(req.URL.Path, "/")
	{{range $desc := .PathParams}}
	protoReq.{{$desc.ProtoName}} = proto.{{$desc.ConvertFunc}}(components[{{$desc.PathIndex}}])
	{{end}}
{{end}}
	return c.{{.Name}}(ctx, protoReq)
}
`))

	trailerTemplate = template.Must(template.New("trailer").Parse(`
var (
	{{range $svc := .}}
	endpoint{{$svc.Name}} string
	{{end}}
)

func init() {
	{{range $svc := .}}
	flag.StringVar(&endpoint{{$svc.Name}}, {{$svc.EndpointFlag | printf "%q"}}, "", "endpoint host:port of {{$svc.Name}}")
	{{end}}
}

type handler struct {
	mux http.ServeMux
	conns map[string]*grpc.ClientConn
}

func (h *handler) Close() error {
	var err error
	for svc, conn := range h.conns {
		cerr := conn.Close()
		if err == nil {
			err = cerr
		}
		if cerr != nil {
			glog.Errorf("Failed to close gRPC connection to %s: %v", svc, err)
		}
	}
	return err
}

func NewHandler(ctx context.Context) (http.Handler, error) {
	h := &handler{
		conn: make(map[string]*grpc.ClientConn),
	}
	var err error
	defer func() {
		if err != nil {
			h.Close()
		}
	}()
	{{range $svc := .}}
	err = func() error {
		conn, err := grpc.Dial(endpoint{{$svc.Name}})
		if err != nil {
			return err
		}
		h.conn[{{$svc.Name | printf "%q"}}] = conn
		client := New{{$svc.Name}}Client(conn)
		{{range $m := $svc.Methods}}
		mux.HandleFunc({{$m.Prefix | printf "%q"}}, func(w http.ResponseWriter, req *http.Request) {
			resp, err := handle_{{$m.ServiceName}}_{{$m.Name}}(ctx, client, req)
			if err != nil {
				glog.Errorf("RPC error: %v", err)
				http.Error(w, err.String(), http.StatusInternalServerError)
				return
			}
			buf, err := proto.Marshal(resp)
			if err != nil {
				glog.Errorf("Marshal error: %v", err)
				http.Error(w, err.String(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err = w.Write(buf); err != nil {
				glog.Errorf("Failed to write response: %v", err)
			}
		})
		{{end}}
	}()
	if err != nil {
		return nil, err
	}
	{{end}}
	return h, nil
}
`))
)
