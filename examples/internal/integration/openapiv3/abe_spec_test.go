// abe_spec_test.go is a Tier 1 structural test for the OpenAPI 3.1 spec
// generated from a_bit_of_everything.proto — by far the gnarliest proto
// corpus in the tree (oneofs, additional_bindings, custom HTTP methods,
// dotted nested path params, multi-wildcard URL templates, deep schema
// references). It loads the checked-in spec and asserts shape facts that
// must hold for the spec to be a faithful description of the gateway.
//
// Drift surfaces as a test failure rather than as a runtime mystery. The
// matching Tier 2 oracle (round-tripping through an openapi-generator-cli
// client and a real gateway) is intentionally out of scope here — that
// test lives behind `make openapiv3-clients` and is opt-in.
package openapiv3_test

import (
	"encoding/json"
	"slices"
	"testing"

	examples "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
)

func loadABESpec(t *testing.T) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(examples.ABitOfEverythingOpenAPIV3Spec, &doc); err != nil {
		t.Fatalf("unmarshal embedded ABE spec: %v", err)
	}
	return doc
}

// TestABESpec_PathExpansion locks down the conversion of proto URL
// templates to OpenAPI URL templates: literal-prefix expansion, multi
// wildcard splitting, dotted nested path params, and verb suffixes.
func TestABESpec_PathExpansion(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	paths, _ := doc["paths"].(map[string]any)

	wantPresent := []string{
		// Single-wildcard literal expansion: {parent=publishers/*}.
		"/v1/publishers/{parent}/books",
		// Multi-wildcard: {book.name=publishers/*/books/*} → two synthetic
		// path params sharing the same proto field.
		"/v1/publishers/{book.name}/books/{book.name_1}",
		// Inline literal inside a longer path: {string_value=strprefix/*}.
		"/v1/example/a_bit_of_everything/{float_value}/{double_value}/{int64_value}/separator/{uint64_value}/{int32_value}/{fixed64_value}/{fixed32_value}/{bool_value}/strprefix/{string_value}/{uint32_value}/{sfixed32_value}/{sfixed64_value}/{sint32_value}/{sint64_value}/{nonConventionalNameValue}/{enum_value}/{path_enum_value}/{nested_path_enum_value}/{enum_value_annotation}",
		// Dotted nested path params on three different RPCs.
		"/v1/example/deep_path/{single_nested.name}",
		"/v1/example/a_bit_of_everything/params/get/{single_nested.name}",
		"/v1/example/a_bit_of_everything/params/get/nested_enum/{single_nested.ok}",
		"/v2/example/a_bit_of_everything/{abe.uuid}",
		// Custom :verb suffix preserved verbatim.
		"/v1/example/a_bit_of_everything/{uuid}:custom",
		"/v1/example/a_bit_of_everything/{uuid}:custom:custom",
		// Root-level :verb path.
		"/v2/{value}:check",
	}
	for _, p := range wantPresent {
		if _, ok := paths[p]; !ok {
			t.Errorf("missing path: %q", p)
		}
	}

	// The unexpanded forms must be absent — if they appear, the literal
	// stripping bug we fixed in convertPathTemplate has regressed.
	wantAbsent := []string{
		"/v1/{parent}/books",
		"/v1/{book.name}",
	}
	for _, p := range wantAbsent {
		if _, ok := paths[p]; ok {
			t.Errorf("path %q must be expanded, found unexpanded form", p)
		}
	}
}

// TestABESpec_PathItemSharing verifies that 6 distinct RPCs binding the
// same URL via different HTTP methods (Lookup/GET, Update/PUT, Delete/DELETE,
// CustomOptionsRequest/OPTIONS, Exists/HEAD, TraceRequest/TRACE) all land
// on the same PathItem. This exercises every branch of PathItem.SetOperation,
// including HEAD/OPTIONS/TRACE which the default openapi-generator targets
// often skip.
func TestABESpec_PathItemSharing(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	paths, _ := doc["paths"].(map[string]any)
	uuidPath, _ := paths["/v1/example/a_bit_of_everything/{uuid}"].(map[string]any)
	if uuidPath == nil {
		t.Fatalf("missing /v1/example/a_bit_of_everything/{uuid}")
	}
	wantOps := map[string]string{
		"get":     "ABitOfEverythingService_Lookup",
		"put":     "ABitOfEverythingService_Update",
		"delete":  "ABitOfEverythingService_Delete",
		"options": "ABitOfEverythingService_CustomOptionsRequest",
		"head":    "ABitOfEverythingService_Exists",
		"trace":   "ABitOfEverythingService_TraceRequest",
	}
	for method, wantID := range wantOps {
		op, _ := uuidPath[method].(map[string]any)
		if op == nil {
			t.Errorf("method %s missing on shared {uuid} path", method)
			continue
		}
		if got, _ := op["operationId"].(string); got != wantID {
			t.Errorf("method %s: operationId = %q, want %q", method, got, wantID)
		}
	}
}

// TestABESpec_AdditionalBindings asserts that RPCs with additional_bindings
// produce one operation per binding and that bindings beyond the first get
// an `_<idx>` suffix on the operationId so the spec stays valid.
func TestABESpec_AdditionalBindings(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	paths, _ := doc["paths"].(map[string]any)

	// Collect every operationId that appears anywhere in the document.
	seen := map[string]bool{}
	for _, item := range paths {
		m, _ := item.(map[string]any)
		for _, op := range m {
			o, _ := op.(map[string]any)
			if id, _ := o["operationId"].(string); id != "" {
				seen[id] = true
			}
		}
	}

	// Echo has 3 bindings (1 base + 2 additional).
	for _, id := range []string{"ABitOfEverythingService_Echo", "ABitOfEverythingService_Echo_1", "ABitOfEverythingService_Echo_2"} {
		if !seen[id] {
			t.Errorf("missing operationId %q", id)
		}
	}
	// UpdateV2 has 3 bindings (1 base + 2 additional).
	for _, id := range []string{"ABitOfEverythingService_UpdateV2", "ABitOfEverythingService_UpdateV2_1", "ABitOfEverythingService_UpdateV2_2"} {
		if !seen[id] {
			t.Errorf("missing operationId %q", id)
		}
	}
	// Custom has 2 bindings (1 base + 1 additional).
	if !seen["ABitOfEverythingService_Custom_1"] {
		t.Errorf("missing operationId ABitOfEverythingService_Custom_1")
	}
}

// TestABESpec_OneofShape asserts that the proto oneof on ABitOfEverything
// (oneof_value { Empty oneof_empty; string oneof_string }) is encoded as
// the spec-correct "at most one set" constraint: a top-level oneOf with a
// {type: object, not: {anyOf: [{required}, {required}]}} guard and one
// per-field option of the form {type: object, properties: {Fi: <schema>},
// required: [Fi]}. Both fields must also appear in the parent `properties`.
//
// The per-field options carry the variant's actual property schema, not
// just `required`.
func TestABESpec_OneofShape(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	schemas, _ := doc["components"].(map[string]any)
	abe, _ := schemas["schemas"].(map[string]any)["grpc.gateway.examples.internal.proto.examplepb.ABitOfEverything"].(map[string]any)
	if abe == nil {
		t.Fatal("missing ABitOfEverything component schema")
	}

	props, _ := abe["properties"].(map[string]any)
	for _, name := range []string{"oneofEmpty", "oneofString"} {
		if _, ok := props[name]; !ok {
			t.Errorf("oneof field %q missing from properties", name)
		}
	}

	oneOf, _ := abe["oneOf"].([]any)
	if len(oneOf) != 3 {
		t.Fatalf("ABitOfEverything: want oneOf with 3 options (none + 2 fields), got %d", len(oneOf))
	}
	none, _ := oneOf[0].(map[string]any)
	if got, _ := none["type"].(string); got != "object" {
		t.Errorf("none-set guard: want type=object, got %v", none["type"])
	}
	if not, _ := none["not"].(map[string]any); not == nil {
		t.Errorf("first oneOf option must be a {not: ...} guard, got %v", none)
	} else if _, ok := not["anyOf"]; !ok {
		t.Errorf("not guard must use anyOf, got %v", not)
	}
	wantPerField := map[string]string{
		"oneofEmpty":  "object", // google.protobuf.Empty WKT inlines as object
		"oneofString": "string",
	}
	requireds := []string{}
	for _, opt := range oneOf[1:] {
		o, _ := opt.(map[string]any)
		if got, _ := o["type"].(string); got != "object" {
			t.Errorf("per-field option: want type=object, got %v", o["type"])
		}
		req, _ := o["required"].([]any)
		if len(req) != 1 {
			t.Errorf("oneOf option must have exactly one required field, got %v", req)
			continue
		}
		fname := req[0].(string)
		requireds = append(requireds, fname)

		// Each branch must embed the variant's property schema with the
		// expected type, so UI viewers see "object with required <field>:
		// <typed-value>" rather than "object with any".
		branchProps, _ := o["properties"].(map[string]any)
		propSchema, _ := branchProps[fname].(map[string]any)
		if propSchema == nil {
			t.Errorf("branch for %q must embed the field's property schema, got %v", fname, branchProps)
			continue
		}
		if got, want := propSchema["type"], wantPerField[fname]; got != want {
			t.Errorf("branch for %q: properties.%s.type = %v, want %v", fname, fname, got, want)
		}
	}
	slices.Sort(requireds)
	if !slices.Equal(requireds, []string{"oneofEmpty", "oneofString"}) {
		t.Errorf("oneOf required options = %v, want [oneofEmpty oneofString]", requireds)
	}
}

// TestABESpec_ComponentSchemas asserts that the component graph reaches
// across imported proto files (pathenum, oneofenum, sub) and that the
// auto-injected google.rpc.Status default-error schema is present.
func TestABESpec_ComponentSchemas(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)

	want := []string{
		// Auto-injected default-error schema.
		"google.rpc.Status",
		// Top-level message and its nested type and nested enum.
		"grpc.gateway.examples.internal.proto.examplepb.ABitOfEverything",
		"grpc.gateway.examples.internal.proto.examplepb.ABitOfEverything.Nested",
		"grpc.gateway.examples.internal.proto.examplepb.ABitOfEverything.Nested.DeepEnum",
		// Reachable via field types: NumericEnum, Book, etc.
		"grpc.gateway.examples.internal.proto.examplepb.NumericEnum",
		"grpc.gateway.examples.internal.proto.examplepb.Book",
		// Reachable via REQUIRED-chained nested messages.
		"grpc.gateway.examples.internal.proto.examplepb.Foo",
		"grpc.gateway.examples.internal.proto.examplepb.Bar",
		// Reachable across imported proto files.
		"grpc.gateway.examples.internal.pathenum.PathEnum",
		"grpc.gateway.examples.internal.pathenum.MessagePathEnum.NestedPathEnum",
		"grpc.gateway.examples.internal.proto.oneofenum.ExampleEnum",
		"grpc.gateway.examples.internal.proto.sub.StringMessage",
	}
	for _, name := range want {
		if _, ok := schemas[name]; !ok {
			t.Errorf("missing component schema %q", name)
		}
	}
}

// TestABESpec_EnumEncoding verifies that enum components are emitted as
// plain string-typed JSON Schema enums. Dual-encoding (also accepting
// integer numbers) is intentionally not advertised in the spec — see
// ensureEnumSchema in protoc-gen-openapiv3 for the tooling reasons.
func TestABESpec_EnumEncoding(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	numeric, _ := schemas["grpc.gateway.examples.internal.proto.examplepb.NumericEnum"].(map[string]any)
	if numeric == nil {
		t.Fatal("missing NumericEnum schema")
	}
	if got, _ := numeric["type"].(string); got != "string" {
		t.Errorf("NumericEnum: want type=string, got %v", numeric["type"])
	}
	vals, _ := numeric["enum"].([]any)
	if len(vals) != 2 {
		t.Errorf("NumericEnum: want 2 enum values (ZERO, ONE), got %d (%v)", len(vals), vals)
	}
}

// TestABESpec_DefaultErrorResponse asserts that every operation in the
// document has a `default` response keyed to google.rpc.Status. This is
// the contract the generator promises and what consumers rely on for
// error handling.
func TestABESpec_DefaultErrorResponse(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	paths, _ := doc["paths"].(map[string]any)
	wantRef := "#/components/schemas/google.rpc.Status"
	for path, item := range paths {
		m, _ := item.(map[string]any)
		for method, op := range m {
			o, _ := op.(map[string]any)
			responses, _ := o["responses"].(map[string]any)
			def, _ := responses["default"].(map[string]any)
			if def == nil {
				t.Errorf("%s %s: missing default response", method, path)
				continue
			}
			content, _ := def["content"].(map[string]any)
			appJSON, _ := content["application/json"].(map[string]any)
			schema, _ := appJSON["schema"].(map[string]any)
			if ref, _ := schema["$ref"].(string); ref != wantRef {
				t.Errorf("%s %s: default schema $ref = %q, want %q", method, path, ref, wantRef)
			}
		}
	}
}

// TestABESpec_MultiWildcardParams verifies that the multi-wildcard URL
// template `{book.name=publishers/*/books/*}` produces one OpenAPI path
// parameter per wildcard, both flagged required and in:path, plus the
// remaining UpdateBookRequest fields as query params.
func TestABESpec_MultiWildcardParams(t *testing.T) {
	t.Parallel()
	doc := loadABESpec(t)
	paths, _ := doc["paths"].(map[string]any)
	updateBook, _ := paths["/v1/publishers/{book.name}/books/{book.name_1}"].(map[string]any)
	if updateBook == nil {
		t.Fatal("missing UpdateBook path")
	}
	patch, _ := updateBook["patch"].(map[string]any)
	if patch == nil {
		t.Fatal("UpdateBook: missing patch operation")
	}
	params, _ := patch["parameters"].([]any)
	wantPath := map[string]bool{"book.name": false, "book.name_1": false}
	wantQuery := map[string]bool{"updateMask": false, "allowMissing": false}
	for _, p := range params {
		m, _ := p.(map[string]any)
		name, _ := m["name"].(string)
		in, _ := m["in"].(string)
		switch in {
		case "path":
			if _, ok := wantPath[name]; ok {
				wantPath[name] = true
			}
			if req, _ := m["required"].(bool); !req {
				t.Errorf("path param %q must be required", name)
			}
		case "query":
			if _, ok := wantQuery[name]; ok {
				wantQuery[name] = true
			}
		}
	}
	for name, seen := range wantPath {
		if !seen {
			t.Errorf("missing path param %q", name)
		}
	}
	for name, seen := range wantQuery {
		if !seen {
			t.Errorf("missing query param %q", name)
		}
	}
}
