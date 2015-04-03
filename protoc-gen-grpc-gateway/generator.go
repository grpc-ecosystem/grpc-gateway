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
		glog.V(1).Infof("register name: %s", name)
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
		if !strings.HasPrefix(pkg, ".") {
			pkg = fmt.Sprintf(".%s", pkg)
		}
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
		base := filepath.Base(d.GetName())
		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)
	}
	return strings.NewReplacer("-", "_", ".", "_").Replace(d.GetPackage())
}

var (
	upperPattern = regexp.MustCompile("[A-Z]")
)

func toCamel(str string) string {
	var components []string
	for _, c := range strings.Split(str, "_") {
		components = append(components, strings.Title(strings.ToLower(c)))
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
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "convert.Float64",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:   "convert.Float32",
		descriptor.FieldDescriptorProto_TYPE_INT64:   "convert.Int64",
		descriptor.FieldDescriptorProto_TYPE_UINT64:  "convert.Uint64",
		descriptor.FieldDescriptorProto_TYPE_INT32:   "convert.Int32",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "convert.Uint64",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "convert.Uint32",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "convert.Bool",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "convert.String",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(yugui) Handle bytes
		descriptor.FieldDescriptorProto_TYPE_UINT32: "convert.Uint32",
		// FieldDescriptorProto_TYPE_ENUM
		// TODO(yugui) Handle Enum
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "convert.Int32",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "convert.Int64",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "convert.Int32",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "convert.Int64",
	}
)

func pathParams(msg *descriptor.DescriptorProto, opts *options.ApiMethodOptions) ([]paramDesc, error) {
	var params []paramDesc
	components := strings.Split(opts.GetPath(), "/")
	for _, c := range components {
		if !strings.HasPrefix(c, ":") {
			continue
		}
		name := strings.TrimPrefix(c, ":")
		fd := lookupField(msg, name)
		if fd == nil {
			return nil, fmt.Errorf("field %q not found in %s", name, msg.GetName())
		}
		conv, ok := convertFuncs[fd.GetType()]
		if !ok {
			return nil, fmt.Errorf("unsupported path parameter type %s in %s", fd.GetType(), msg.GetName())
		}
		params = append(params, paramDesc{
			ProtoName:   name,
			ConvertFunc: conv,
		})
	}
	return params, nil
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
	buf := bytes.NewBuffer(nil)
	var svcDescs []serviceDesc
	var anyNeedsBody, anyNeedsPathParam bool
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
			params, err := pathParams(input, opts)
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
				Path:        opts.GetPath(),
				RequestType: input.GetName(),
				PathParams:  params,
				NeedsBody:   needsBody,
			}
			if md.NeedsBody {
				anyNeedsBody = true
			}
			if md.NeedsPathParam() {
				anyNeedsPathParam = true
			}
			sd.Methods = append(sd.Methods, md)
		}
		if len(sd.Methods) == 0 {
			continue
		}
		svcDescs = append(svcDescs, sd)
	}

	err := headerTemplate.Execute(buf, headerParams{
		Pkg:            goPackage(file),
		NeedsBody:      anyNeedsBody,
		NeedsPathParam: anyNeedsPathParam,
	})
	if err != nil {
		return "", err
	}
	for _, sd := range svcDescs {
		for _, md := range sd.Methods {
			if err = handlerTemplate.Execute(buf, md); err != nil {
				return "", err
			}
		}
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

type headerParams struct {
	Pkg                       string
	NeedsBody, NeedsPathParam bool
}

type paramDesc struct {
	ProtoName   string
	ConvertFunc string
}

func (d paramDesc) GoName() string {
	return toCamel(d.ProtoName)
}

type methodDesc struct {
	ServiceName string
	Name        string
	Method      string
	Path        string
	RequestType string
	QueryParams []paramDesc
	PathParams  []paramDesc
	NeedsBody   bool
}

func (d methodDesc) MuxRegistererName() string {
	return toCamel(d.Method)
}

func (d methodDesc) NeedsPathParam() bool {
	return len(d.PathParams) != 0
}

type serviceDesc struct {
	Name    string
	Methods []methodDesc
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`
package {{.Pkg}}
import (
	{{if .NeedsBody}}
	"encoding/json"
	{{end}}
	{{if .NeedsPathParam}}
	"fmt"
	{{end}}
	"net/http"

	"google.golang.org/grpc"
	{{if (or .NeedsBody .NeedsPathParam)}}
	"github.com/gengo/grpc-gateway/convert"
	{{end}}
	"github.com/golang/protobuf/proto"
	"github.com/golang/glog"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)
`))
	handlerTemplate = template.Must(template.New("handler").Parse(`
func handle_{{.ServiceName}}_{{.Name}}(ctx context.Context, c web.C, client {{.ServiceName}}Client, req *http.Request) (msg proto.Message, err error) {
	protoReq := new({{.RequestType}})
{{if .NeedsBody}}
	if err = json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}
{{end}}
	{{range $desc := .QueryParams}}
	protoReq.{{$desc.ProtoName}}, err = {{$desc.ConvertFunc}}(req.FormValue({{$desc.ProtoName | printf "%q"}}))
	if err != nil {
		return nil, err
	}
	{{end}}
{{if .NeedsPathParam}}
	var val string
	var ok bool
	{{range $desc := .PathParams}}
	val, ok = c.URLParams[{{$desc.ProtoName | printf "%q"}}]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", {{$desc.ProtoName | printf "%q"}})
	}
	protoReq.{{$desc.GoName}}, err = {{$desc.ConvertFunc}}(val)
	if err != nil {
		return nil, err
	}
	{{end}}
{{end}}
	return client.{{.Name}}(ctx, protoReq)
}
`))

	trailerTemplate = template.Must(template.New("trailer").Parse(`
{{range $svc := .}}
func Register{{$svc.Name}}HandlerFromEndpoint(ctx context.Context, mux *web.Mux, endpoint string) (err error) {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				glog.Error("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				glog.Error("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return Register{{$svc.Name}}Handler(ctx, mux, conn)
}

func Register{{$svc.Name}}Handler(ctx context.Context, mux *web.Mux, conn *grpc.ClientConn) error {
	client := New{{$svc.Name}}Client(conn)
	{{range $m := $svc.Methods}}
	mux.{{$m.MuxRegistererName}}({{$m.Path | printf "%q"}}, func(c web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_{{$m.ServiceName}}_{{$m.Name}}(ctx, c, client, req)
		if err != nil {
			glog.Errorf("RPC error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf, err := proto.Marshal(resp)
		if err != nil {
			glog.Errorf("Marshal error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(buf); err != nil {
			glog.Errorf("Failed to write response: %v", err)
		}
	})
	{{end}}
	return nil
}
{{end}}`))
)
