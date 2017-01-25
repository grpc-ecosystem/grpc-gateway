package gengateway

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/golang/glog"
	"github.com/shilkin/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/shilkin/grpc-gateway/utilities"
)

type param struct {
	*descriptor.File
	Imports           []descriptor.GoPackage
	UseRequestContext bool
}

type binding struct {
	*descriptor.Binding
}

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

type trailerParams struct {
	Services          []*descriptor.Service
	UseRequestContext bool
}

func applyTemplate(p param) (string, error) {
	w := bytes.NewBuffer(nil)
	if err := headerTemplate.Execute(w, p); err != nil {
		return "", err
	}
	var targetServices []*descriptor.Service
	for _, svc := range p.Services {
		var methodWithBindingsSeen bool
		for _, meth := range svc.Methods {
			glog.V(2).Infof("Processing %s.%s", svc.GetName(), meth.GetName())
			for _, b := range meth.Bindings {
				methodWithBindingsSeen = true
				if err := handlerTemplate.Execute(w, binding{Binding: b}); err != nil {
					return "", err
				}
			}
		}
		if methodWithBindingsSeen {
			targetServices = append(targetServices, svc)
		}
	}
	if len(targetServices) == 0 {
		return "", errNoTargetService
	}

	tp := trailerParams{
		Services:          targetServices,
		UseRequestContext: p.UseRequestContext,
	}
	if err := trailerTemplate.Execute(w, tp); err != nil {
		return "", err
	}
	return w.String(), nil
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`
// Code generated by protoc-gen-grpc-gateway
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
var _ = utilities.NewDoubleArray
`))

	handlerTemplate = template.Must(template.New("handler").Parse(`
{{if and .Method.GetClientStreaming .Method.GetServerStreaming}}
{{template "bidi-streaming-request-func" .}}
{{else if .Method.GetClientStreaming}}
{{template "client-streaming-request-func" .}}
{{else}}
{{template "client-rpc-request-func" .}}
{{end}}
`))

	_ = template.Must(handlerTemplate.New("request-func-signature").Parse(strings.Replace(`
{{if .Method.GetServerStreaming}}
func request_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}(ctx context.Context, marshaler runtime.Marshaler, client {{.Method.Service.GetName}}Client, req *http.Request, pathParams map[string]string) ({{.Method.Service.GetName}}_{{.Method.GetName}}Client, runtime.ServerMetadata, error)
{{else}}
func request_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}(ctx context.Context, marshaler runtime.Marshaler, client {{.Method.Service.GetName}}Client, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error)
{{end}}`, "\n", "", -1)))

	_ = template.Must(handlerTemplate.New("client-streaming-request-func").Parse(`
{{template "request-func-signature" .}} {
	var metadata runtime.ServerMetadata
	stream, err := client.{{.Method.GetName}}(ctx)
	if err != nil {
		grpclog.Printf("Failed to start streaming: %v", err)
		return nil, metadata, err
	}
	dec := marshaler.NewDecoder(req.Body)
	for {
		var protoReq {{.Method.RequestType.GoType .Method.Service.File.GoPkg.Path}}
		err = dec.Decode(&protoReq)
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Printf("Failed to decode request: %v", err)
			return nil, metadata, grpc.Errorf(codes.InvalidArgument, "%v", err)
		}
		if err = stream.Send(&protoReq); err != nil {
			grpclog.Printf("Failed to send request: %v", err)
			return nil, metadata, err
		}
	}

	if err := stream.CloseSend(); err != nil {
		grpclog.Printf("Failed to terminate client stream: %v", err)
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		grpclog.Printf("Failed to get header from client: %v", err)
		return nil, metadata, err
	}
	metadata.HeaderMD = header
{{if .Method.GetServerStreaming}}
	return stream, metadata, nil
{{else}}
	msg, err := stream.CloseAndRecv()
	metadata.TrailerMD = stream.Trailer()
	return msg, metadata, err
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
	var metadata runtime.ServerMetadata
{{if .Body}}
	if err := marshaler.NewDecoder(req.Body).Decode(&{{.Body.RHS "protoReq"}}); err != nil {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "%v", err)
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
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "missing parameter %s", {{$param | printf "%q"}})
	}
{{if $param.IsNestedProto3 }}
	err = runtime.PopulateFieldFromPath(&protoReq, {{$param | printf "%q"}}, val)
{{else}}
	{{$param.RHS "protoReq"}}, err = {{$param.ConvertFuncExpr}}(val)
{{end}}
	if err != nil {
		return nil, metadata, err
	}
	{{end}}
{{end}}
{{if .HasQueryParam}}
	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}); err != nil {
		return nil, metadata, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}
{{end}}
{{if .Method.GetServerStreaming}}
	stream, err := client.{{.Method.GetName}}(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil
{{else}}
	msg, err := client.{{.Method.GetName}}(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
{{end}}
}`))

	_ = template.Must(handlerTemplate.New("bidi-streaming-request-func").Parse(`
{{template "request-func-signature" .}} {
	var metadata runtime.ServerMetadata
	stream, err := client.{{.Method.GetName}}(ctx)
	if err != nil {
		grpclog.Printf("Failed to start streaming: %v", err)
		return nil, metadata, err
	}
	dec := marshaler.NewDecoder(req.Body)
	handleSend := func() error {
		var protoReq {{.Method.RequestType.GoType .Method.Service.File.GoPkg.Path}}
		err = dec.Decode(&protoReq)
		if err == io.EOF {
			return err
		}
		if err != nil {
			grpclog.Printf("Failed to decode request: %v", err)
			return err
		}
		if err = stream.Send(&protoReq); err != nil {
			grpclog.Printf("Failed to send request: %v", err)
			return err
		}
		return nil
	}
	if err := handleSend(); err != nil {
		if cerr := stream.CloseSend(); cerr != nil {
			grpclog.Printf("Failed to terminate client stream: %v", cerr)
		}
		if err == io.EOF {
			return stream, metadata, nil
		}
		return nil, metadata, err
	}
	go func() {
		for {
			if err := handleSend(); err != nil {
				break
			}
		}
		if err := stream.CloseSend(); err != nil {
			grpclog.Printf("Failed to terminate client stream: %v", err)
		}
	}()
	header, err := stream.Header()
	if err != nil {
		grpclog.Printf("Failed to get header from client: %v", err)
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil
}
`))

	trailerTemplate = template.Must(template.New("trailer").Parse(`
{{$UseRequestContext := .UseRequestContext}}
{{range $svc := .Services}}
// Register{{$svc.GetName}}HandlerFromEndpoint is same as Register{{$svc.GetName}}Handler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func Register{{$svc.GetName}}HandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	middleware := map[string]runtime.Middleware{}
	return Register{{$svc.GetName}}HandlerFromEndpointWithMiddleware(ctx, mux, middleware, endpoint, opts)
}

// Register{{$svc.GetName}}HandlerFromEndpointWithMiddleware is same as Register{{$svc.GetName}}HandlerWithMiddleware but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
// It receives a map of additional http handlres and wraps final handler with chain of middleware.
// Each middleware has a name.
// To include middleware in chain of calls you need to specify this name in rpc method option.
func Register{{$svc.GetName}}HandlerFromEndpointWithMiddleware(ctx context.Context, mux *runtime.ServeMux, middleware map[string]runtime.Middleware, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return Register{{$svc.GetName}}HandlerWithMiddleware(ctx, mux, middleware, conn)
}

// Register{{$svc.GetName}}Handler registers the http handlers for service {{$svc.GetName}} to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func Register{{$svc.GetName}}Handler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	middleware := map[string]runtime.Middleware{}
	return Register{{$svc.GetName}}HandlerWithMiddleware(ctx, mux, middleware, conn)
}

// Register{{$svc.GetName}}HandlerWithMiddleware registers the http handlers for service {{$svc.GetName}} to "mux".
// It receives a map of additional http handlres and wraps final handler with chain of middleware.
// Each middleware has a name.
// To include middleware in chain of calls you need to specify this name in rpc method option.
// The handlers forward requests to the grpc endpoint over "conn".
func Register{{$svc.GetName}}HandlerWithMiddleware(ctx context.Context, mux *runtime.ServeMux, middleware map[string]runtime.Middleware, conn *grpc.ClientConn) error {
	client := New{{$svc.GetName}}Client(conn)
	var handler runtime.HandlerFunc
	var mw []string

	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}

	mw = []string{ {{range $name := $b.Middleware}}
			"{{$name}}",
			{{end}} }

	handler = func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
  {{- if $UseRequestContext }}
		ctx, cancel := context.WithCancel(req.Context())
	{{- else -}}
		ctx, cancel := context.WithCancel(ctx)
	{{- end }}
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, req)
		if err != nil {
			runtime.HTTPError(ctx, outboundMarshaler, w, req, err)
		}
		resp, md, err := request_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, outboundMarshaler, w, req, err)
			return
		}
		{{if $m.GetServerStreaming}}
		forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(ctx, outboundMarshaler, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)
		{{else}}
		forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(ctx, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
		{{end}}
	}

	for _, name := range mw {
		if m, ok := middleware[name]; ok {
			handler = m(handler)
		}
	}

	mux.Handle({{$b.HTTPMethod | printf "%q"}}, pattern_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}, handler)
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
