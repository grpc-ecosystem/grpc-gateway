package gengateway

import (
	"strings"
	"text/template"
)

var (
	localHandlerTemplate = template.Must(template.New("local-handler").Parse(`
{{if and .Method.GetClientStreaming .Method.GetServerStreaming}}
// TODO bidi-streaming-request
{{else if .Method.GetClientStreaming}}
// TODO client-streaming-request
{{else if .Method.GetServerStreaming}}
// TODO server-streaming-request
{{else}}
{{template "local-rpc-request-func" .}}
{{end}}
`))

	_ = template.Must(localHandlerTemplate.New("local-request-func-signature").Parse(strings.Replace(`
{{if .Method.GetServerStreaming}}
// TODO
{{else}}
func local_request_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}(ctx context.Context, marshaler runtime.Marshaler, server {{.Method.Service.GetName}}Server, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error)
{{end}}`, "\n", "", -1)))

	_ = template.Must(localHandlerTemplate.New("local-rpc-request-func").Parse(`
{{$AllowPatchFeature := .AllowPatchFeature}}
{{template "local-request-func-signature" .}} {
	var protoReq {{.Method.RequestType.GoType .Method.Service.File.GoPkg.Path}}
	var metadata runtime.ServerMetadata
{{if .Body}}
	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&{{.Body.AssignableExpr "protoReq"}}); err != nil && err != io.EOF  {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	{{- if and $AllowPatchFeature (and (eq (.HTTPMethod) "PATCH") (.FieldMaskField))}}
	if protoReq.{{.FieldMaskField}} != nil && len(protoReq.{{.FieldMaskField}}.GetPaths()) > 0 {
		runtime.CamelCaseFieldMask(protoReq.{{.FieldMaskField}})
	} {{if not (eq "*" .GetBodyFieldPath)}} else {
			if fieldMask, err := runtime.FieldMaskFromRequestBody(newReader()); err != nil {
				return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
			} else {
				protoReq.{{.FieldMaskField}} = fieldMask
			}		
	} {{end}}		
	{{end}}
{{end}}
{{if .PathParams}}
	var (
		val string
{{- if .HasEnumPathParam}}
		e int32
{{- end}}
{{- if .HasRepeatedEnumPathParam}}
		es []int32
{{- end}}
		ok bool
		err error
		_ = err
	)
	{{$binding := .}}
	{{range $param := .PathParams}}
	{{$enum := $binding.LookupEnum $param}}
	val, ok = pathParams[{{$param | printf "%q"}}]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", {{$param | printf "%q"}})
	}
{{if $param.IsNestedProto3}}
	err = runtime.PopulateFieldFromPath(&protoReq, {{$param | printf "%q"}}, val)
{{else if $enum}}
	e{{if $param.IsRepeated}}s{{end}}, err = {{$param.ConvertFuncExpr}}(val{{if $param.IsRepeated}}, {{$binding.Registry.GetRepeatedPathParamSeparator | printf "%c" | printf "%q"}}{{end}}, {{$enum.GoType $param.Target.Message.File.GoPkg.Path}}_value)
{{else}}
	{{$param.AssignableExpr "protoReq"}}, err = {{$param.ConvertFuncExpr}}(val{{if $param.IsRepeated}}, {{$binding.Registry.GetRepeatedPathParamSeparator | printf "%c" | printf "%q"}}{{end}})
{{end}}
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", {{$param | printf "%q"}}, err)
	}
{{if and $enum $param.IsRepeated}}
	s := make([]{{$enum.GoType $param.Target.Message.File.GoPkg.Path}}, len(es))
	for i, v := range es {
		s[i] = {{$enum.GoType $param.Target.Message.File.GoPkg.Path}}(v)
	}
	{{$param.AssignableExpr "protoReq"}} = s
{{else if $enum}}
	{{$param.AssignableExpr "protoReq"}} = {{$enum.GoType $param.Target.Message.File.GoPkg.Path}}(e)
{{end}}
	{{end}}
{{end}}
{{if .HasQueryParam}}
	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_{{.Method.Service.GetName}}_{{.Method.GetName}}_{{.Index}}); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
{{end}}
{{if .Method.GetServerStreaming}}
	// TODO
{{else}}
	msg, err := server.{{.Method.GetName}}(ctx, &protoReq)
	return msg, metadata, err
{{end}}
}`))

	localTrailerTemplate = template.Must(template.New("local-trailer").Parse(`
{{$UseRequestContext := .UseRequestContext}}
{{range $svc := .Services}}
// {{$svc.GetName}} local server register
func Register{{$svc.GetName}}{{$.RegisterFuncSuffix}}Server(ctx context.Context, mux *runtime.ServeMux, server {{$svc.GetName}}Server) error {
	{{range $m := $svc.Methods}}
	{{range $b := $m.Bindings}}
	mux.Handle({{$b.HTTPMethod | printf "%q"}}, pattern_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
	{{- if $UseRequestContext }}
		ctx, cancel := context.WithCancel(req.Context())
	{{- else -}}
		ctx, cancel := context.WithCancel(ctx)
	{{- end }}
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		{{if $m.GetServerStreaming}}
		// TODO
		{{else}}
		{{ if $b.ResponseBody }}
		forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(ctx, mux, outboundMarshaler, w, req, response_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}{resp}, mux.GetForwardResponseOptions()...)
		{{ else }}
		forward_{{$svc.GetName}}_{{$m.GetName}}_{{$b.Index}}(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
		{{end}}
		{{end}}
	})
	{{end}}
	{{end}}
	return nil
}
{{end}}`))
)
