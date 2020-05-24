---
title: v2 migration guide
category: documentation
---

# gRPC-Gateway v2 migration guide

This guide is supposed to help users of the gateway migrate from v1 to v2.
See https://github.com/grpc-ecosystem/grpc-gateway/issues/1223 for detailed
information on all changes that were made specifically to v2.

The following behavioural defaults have been changed:

## protoc-gen-swagger has been renamed protoc-gen-openapiv2

See
[the original issue](https://github.com/grpc-ecosystem/grpc-gateway/issues/675)
for more information. Apart from the new name, the only real
difference to users will be a slightly different proto annotation:

```protobuf
import "protoc-gen-openapiv2/options/annotations.proto";
```
instead of
```protobuf
import "protoc-gen-swagger/options/annotations.proto";
```
and
```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
```
instead of
```protobuf
option (grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger) = {
```

The bazel rule has been renamed `protoc_gen_openapiv2`.

## We now use the camelCase JSON names by default
See
[the original issue](https://github.com/grpc-ecosystem/grpc-gateway/issues/375)
and
[original pull request](https://github.com/grpc-ecosystem/grpc-gateway/issues/375)
for more information.

If you want to revert to the old behaviour, configure a custom marshaler with
`UseProtoNames: true`:
```go
mux := runtime.NewServeMux(
	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}),
)
```

To change the OpenAPI generator behaviour to match, set `json_names_for_fields=false` when generating:

```shell
--openapiv2_out=json_names_for_fields=false:./gen/openapiv2 path/to/my/proto/v1/myproto.proto
```

If using the Bazel rule, set `json_names_for_fields=False`.

## We now emit default vaules for all fields

See [the original issue](https://github.com/grpc-ecosystem/grpc-gateway/issues/233)
for more information.

If you want to revert to the old behaviour, configure a custom marshaler with
`EmitUnpopulated: false`:
```go
mux := runtime.NewServeMux(
	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: false,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}),
)
```

## We now support google.api.HttpBody message types by default

The `runtime.SetHTTPBodyMarshaler` function has disappeared, and is now
enabled by default. If you for some reason don't want `HttpBody` messages to be
respected, you can disable it by overwriting the default marshaler:

```go
mux := runtime.NewServeMux(
	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}),
)
```

## runtime.DisallowUnknownFields has been removed

All marshalling settings are now inherited from the configured marshaler. If you wish
to disallow unknown fields, configure a custom marshaler:

```go
mux := runtime.NewServeMux(
	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: false,
		},
	}),
)
```

## WithLastMatchWins and allow_colon_final_segments=true is now default behaviour

If you were previously specifying these, please remove them, as this is now
the default behaviour. See
[the original issue](https://github.com/grpc-ecosystem/grpc-gateway/issues/224)
for more information.

There is no workaround for this, as we considered it a correct interpretation of the spec.
If this breaks your application, carefully consider the order in which you define your
services.

## Error handling configuration has been overhauled

`runtime.HTTPError`, `runtime.OtherErrorHandler`, `runtime.GlobalHTTPErrorHandler`,
`runtime.WithProtoErrorHandler` are all gone. Error handling is rewritten around the
use of gRPCs Status types. If you wish to configure how the gateway handles errors,
please use `runtime.WithErrorHandler` and `runtime.WithStreamErrorHandler`.
