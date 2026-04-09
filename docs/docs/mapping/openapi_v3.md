---
layout: default
title: OpenAPI 3.1 Output
nav_order: 7
parent: Mapping
---

# OpenAPI 3.1 Output

{: .warning }
**Alpha — output is not yet stable.** `protoc-gen-openapiv3` is a new
generator and the exact JSON shape it emits may change between minor
releases as the mapping rules settle in response to user feedback.
Specifically expect the encodings for oneofs, wrapper types, enums, and
path-template expansion to evolve, and expect tooling-compatibility
compromises (documented inline in the generator source) to be revisited
as consumer OpenAPI 3.1 support matures. If you need a production-stable
OpenAPI pipeline today, use [`protoc-gen-openapiv2`](./customizing_openapi_output.md).

{: .note }
`protoc-gen-openapiv3` emits OpenAPI 3.1.0 JSON directly from
`google.api.http` annotations. It is intentionally smaller in scope than
[`protoc-gen-openapiv2`](./customizing_openapi_output.md) and does **not**
consume the `grpc.gateway.protoc_gen_openapiv2.options` annotations. If
you rely on those options, continue to use the v2 generator.

## What it does

`protoc-gen-openapiv3` walks every proto file it is asked to generate, finds
services with HTTP bindings, and emits **one OpenAPI 3.1.0 JSON document per
proto file** (`foo.proto` → `foo.openapi.json`). Files with no HTTP-bound
services produce no output.

The output is deterministic and byte-stable across runs: paths are emitted in
RPC declaration order, component schemas are sorted alphabetically, and
response codes are emitted default-first then sorted.

## What it does not do

The generator is deliberately opinionated and has no configuration flags. The
following are **not** supported today:

- The `grpc.gateway.protoc_gen_openapiv2.options` annotation set.
- OpenAPI 2.0 / Swagger output (use `protoc-gen-openapiv2` for that).
- YAML output.
- Alternative naming strategies. Component names are always the
  fully-qualified proto name with the leading dot stripped
  (e.g. `lib.v1.Book`).
- Integer enums. Enums are always rendered as strings.
- Configurable error response schemas. A `google.rpc.Status` default
  response is injected automatically on every operation.

Features can be added back as concrete needs emerge — if you want one of the
above, please open an issue describing your use case.

## Installation

```sh
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3@latest
```

## Usage with buf

Add the plugin to your `buf.gen.yaml`:

```yaml
version: v2
plugins:
  - local: protoc-gen-openapiv3
    out: .
```

Then run `buf generate` as usual. A `.openapi.json` file will appear next to
each proto file that declares HTTP bindings.

## Usage with protoc

```sh
protoc -I. \
  --openapiv3_out=. \
  path/to/your/service.proto
```

## Mapping rules

The mapping from proto constructs to OpenAPI is fixed and summarised below.

### Services and methods

| Proto                    | OpenAPI                                           |
| ------------------------ | ------------------------------------------------- |
| Service                  | Tag (`tags[]`) plus `tags` on each operation      |
| Method                   | Operation under a `PathItem`                      |
| Method leading comment   | Operation `summary` / `description`               |
| Method `deprecated=true` | Operation `deprecated: true`                      |
| Additional bindings      | Extra operations; `operationId` suffixed `_<idx>` |

The default `operationId` is `<Service>_<Method>`. When a method has multiple
HTTP bindings, the first uses the bare ID and subsequent bindings append
`_1`, `_2`, ... so the spec remains valid.

### Parameters

- Path parameters from the `google.api.http` URL template become `in: path`
  parameters.
- Constrained path templates like `{name=shelves/*}` are **expanded** into
  the OpenAPI URL verbatim: the literal prefix stays in the URL and each
  wildcard becomes its own OpenAPI parameter. This is required because the
  proto field binds to the full matched substring (e.g. `shelves/abc`), so
  a client following the OpenAPI URL must actually send `/shelves/abc`
  rather than just `/abc`. Examples:

  | Proto template                 | OpenAPI URL                         |
  | ------------------------------ | ----------------------------------- |
  | `/v1/{name}`                   | `/v1/{name}`                        |
  | `/v1/{name=shelves/*}`         | `/v1/shelves/{name}`                |
  | `/v1/{name=shelves/*/books/*}` | `/v1/shelves/{name}/books/{name_1}` |
  | `/v1/{name=files/**}`          | `/v1/files/{name}`                  |

  When a single proto field expands into multiple wildcards, the second and
  subsequent parameters are suffixed `_1`, `_2`, ... so the OpenAPI URL
  stays valid. All synthetic parameters point back at the same underlying
  proto field, which the grpc-gateway runtime reconstructs from the full
  matched substring at request time.

- Fields that are neither path parameters nor part of the request body
  become `in: query` parameters.
- `body="*"` synthesises an inline object schema that includes every request
  field **except** path parameters.
- `body="field"` uses that single field's type as the request body.

### Schemas

- Messages become component schemas under `#/components/schemas/`.
- Enums become string-valued component schemas.
- `repeated` fields produce `type: array`.
- Map fields produce `type: object` with `additionalProperties`.
- 64-bit integer fields are emitted as JSON strings, matching protojson.
- Well-known types are inlined, not referenced (see below).
- Field descriptions are populated from leading proto comments.
- `google.api.field_behavior` maps to `required`, `readOnly`, and
  `writeOnly` on the parent schema.
- `deprecated = true` on a field sets `deprecated: true` on the property.

### Oneofs

Every field inside a `oneof` appears as a normal property on the component
schema, and each `oneof` group contributes an "at most one set" constraint
on top of those properties. The constraint is expressed as a JSON Schema
`oneOf` whose options are:

1. a "none of the fields are set" guard of the form
   `{"not": {"anyOf": [{"required": ["fieldA"]}, {"required": ["fieldB"]}, ...]}}`;
2. one `{"required": ["fieldN"]}` option per field in the group.

Because `oneOf` matches exactly one sub-schema:

- zero fields set → only option 1 matches → passes,
- exactly one field set → only that field's `required` option matches → passes,
- two or more set → multiple `required` options match → fails validation.

This matches proto3 `oneof` semantics, which allow either zero or one
member to be set.

Messages with a single `oneof` group hoist the constraint directly onto
the component schema's `oneOf`. Messages with multiple groups wrap their
group constraints in `allOf` so each group is enforced independently.

Single-field `oneof` groups produce a trivially-true constraint and are
skipped for output hygiene. Synthetic proto3-optional oneofs are treated
as regular optional fields (the proto compiler models `optional string
foo = 1;` as a one-element synthetic oneof).

### Well-known types

Well-known types are **inlined** as their protojson representation:

| WKT                           | Schema                          |
| ----------------------------- | ------------------------------- |
| `google.protobuf.Timestamp`   | `string`, `format: date-time`   |
| `google.protobuf.Duration`    | `string`                        |
| `google.protobuf.FieldMask`   | `string`                        |
| `google.protobuf.StringValue` | `string`                        |
| `google.protobuf.BytesValue`  | `string`, `format: byte`        |
| `google.protobuf.Int32Value`  | `integer`, `format: int32`      |
| `google.protobuf.UInt32Value` | `integer`, `format: int64`      |
| `google.protobuf.Int64Value`  | `string`, `format: int64`       |
| `google.protobuf.UInt64Value` | `string`, `format: uint64`      |
| `google.protobuf.FloatValue`  | `number`, `format: float`       |
| `google.protobuf.DoubleValue` | `number`, `format: double`      |
| `google.protobuf.BoolValue`   | `boolean`                       |
| `google.protobuf.Empty`       | `object`                        |
| `google.protobuf.Struct`      | `object`                        |
| `google.protobuf.Value`       | unconstrained                   |
| `google.protobuf.ListValue`   | `array` of unconstrained        |
| `google.protobuf.NullValue`   | `null`                          |
| `google.protobuf.Any`         | `object` with `@type` plus open |

**Wrapper types are intentionally not marked nullable.** The strictly-correct
JSON Schema 2020-12 form for `google.protobuf.StringValue` would be
`type: ["string", "null"]`, but no Go OpenAPI generator in the ecosystem
handles that form today: `openapi-generator-cli`'s Go target produces
unbuildable code, and `oapi-codegen` (including its experimental 3.1 fork)
errors with `unhandled Schema type: &[string null]` on request bodies.
Until tooling catches up, the generator describes wrappers as their
underlying primitive. The grpc-gateway runtime still accepts a JSON `null`
on the wire regardless, because `protojson` treats wrappers as optional.

RPCs returning `google.protobuf.Empty` emit a `200` with an `object` schema
body (`{}`), matching what the `grpc-gateway` runtime actually writes on
the wire. The HTTP-conventional `204 No Content` would be more idiomatic
but conflicts with runtime behaviour, so generated clients would reject
valid responses.

### Error responses

Every operation automatically gets a `default` response keyed to a
component schema named `google.rpc.Status` with the standard `code`,
`message`, and `details[]` fields. This matches what `grpc-gateway` returns
on error paths.

## Example

Given:

```protobuf
syntax = "proto3";

package example.v1;

import "google/api/annotations.proto";

service EchoService {
  rpc Echo(EchoRequest) returns (EchoResponse) {
    option (google.api.http) = {
      post: "/v1/echo"
      body: "*"
    };
  }
}

message EchoRequest  { string message = 1; }
message EchoResponse { string message = 1; }
```

The generator emits `echo.openapi.json` with one `POST /v1/echo` operation,
an `EchoResponse` component schema, and the auto-injected
`google.rpc.Status` default error response. See
[`protoc-gen-openapiv3/internal/genopenapi/testdata/simple_echo.openapi.json`](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-openapiv3/internal/genopenapi/testdata/simple_echo.openapi.json)
for the exact output.

## Consuming the output

The generated spec validates as OpenAPI 3.1.0 and can be fed to any 3.1-aware
tool, including [`openapi-generator-cli`](https://github.com/OpenAPITools/openapi-generator)
for client generation. An end-to-end example that runs a real grpc-gateway
behind a generated Go client lives at
[`examples/internal/integration/openapiv3`](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/examples/internal/integration/openapiv3).
