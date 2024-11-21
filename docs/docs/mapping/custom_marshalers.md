---
layout: default
title: Custom marshalers
nav_order: 6
parent: Mapping
---

# Custom marshalers

[`Marshaler`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime?tab=doc#Marshaler)
implementations can implement optional additional methods to customize their
behaviour beyond the methods required by the core interface.

## Stream delimiters

By default, a streamed response delimits each response body with a single
newline (`"\n"`). You can change this delimiter by having your marshaler
implement
[`Delimited`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime#Delimited).

For example, to separate each entry with a pipe (`"|"`) instead:

```go
type YourMarshaler struct {
  // ...
}

// ...

func (*YourMarshaler) Delimiter() []byte {
  return []byte("|")
}
```

## Stream content type

By default, a streamed response emits a `Content-Type` header that is the same
for a unary response, from the `ContentType()` method of the
[`Marshaler`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime?tab=doc#Marshaler)
interface.

If you require the server to declare a distinct content type for stream
responses versus unary responses, the marshaler must implement
[`StreamContentType`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime#StreamContentType).
This provides the MIME type when specifically responding to a streaming
response.

For example, by default the
[`JSONPb`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime#JSONPb)
marshaler results in `application/json` for its `Content-Type` response header,
irrespective of unary versus streaming. This can be changed for streaming
endpoints by wrapping the marshaler with a custom marshaler that implements
[`StreamContentType`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime#StreamContentType)
to return the [NDJSON](https://github.com/ndjson/ndjson-spec) MIME type for
streaming response endpoints:

```go
type CustomJSONPb struct {
  runtime.JSONPb
}

func (*CustomJSONPb) Delimiter() []byte {
  // Strictly speaking this is already the default delimiter for JSONPb, but
  // providing it here for completeness with an NDJSON marshaler all in one
  // place.
  return []byte("\n")
}

func (*CustomJSONPb) StreamContentType(interface{}) string {
  return "application/x-ndjson"
}
```
