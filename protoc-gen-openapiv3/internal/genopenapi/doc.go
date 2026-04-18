// Package genopenapi generates OpenAPI 3.1.0 JSON documents from grpc-gateway
// proto descriptors. It is the implementation backing the protoc-gen-openapiv3
// plugin.
//
// # Design
//
// There is no canonical intermediate model. The types in types.go are both
// the builder's write target and the encoding/json source, which keeps the
// pipeline short and avoids adapter layers.
//
// # Determinism
//
// Output is byte-stable across runs. Paths are emitted in RPC declaration
// order via a custom MarshalJSON on Paths; component maps are emitted in
// lexicographic key order via a custom MarshalJSON on Components; response
// codes are emitted default-first, then sorted. Regular Go maps inside
// Schema are serialized by encoding/json, which sorts keys alphabetically.
//
// # Schema references
//
// Message and enum types are emitted as component schemas and referenced via
// $ref. Descriptions on a $ref use the 3.1.0 sibling form; per-occurrence
// flags (deprecated, readOnly, writeOnly) require a real schema body, so
// those are wrapped in allOf. Cycles terminate because ensureMessageSchema
// reserves the component slot before recursing into fields, so a
// self-reference short-circuits on the exists check.
//
// # File layout
//
//   - generator.go — entry point; one Document per input file.
//   - operation.go — one HTTP binding → one OpenAPI Operation, including the
//     path/query parameter split, body="*" synthesis, and the auto-injected
//     google.rpc.Status default error response.
//   - schema.go    — field → schema, message → component. Handles oneofs,
//     maps, field_behavior, and cycle-safe recursion.
//   - wkt.go       — inline schemas for google.protobuf well-known types,
//     matching protojson's wire representation.
//   - types.go     — the OpenAPI 3.1.0 data model plus custom MarshalJSON
//     shims that drive deterministic output.
//   - comments.go  — source-code-info path construction for leading comments
//     on services, methods, messages, fields, and enums.
//   - path.go      — proto URL template → OpenAPI URL template conversion.
//   - naming.go    — FQN → component schema name.
package genopenapi
