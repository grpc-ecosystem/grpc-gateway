package genopenapi

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/pluginpb"
)

// TestGenerate_Library exercises a richer shape: query parameters, repeated
// fields, an enum, and a referenced message in the response. We don't golden
// the whole document — we just assert structural facts that lock down the
// behaviors we care about, so the test stays useful as the generator grows.
func TestGenerate_Library(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/library.prototext")
	got := runGenerator(t, req)

	var doc map[string]any
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, string(got))
	}

	// 1. The list endpoint exists and has the expected shape.
	paths, _ := doc["paths"].(map[string]any)
	listOp, ok := paths["/v1/books"].(map[string]any)
	if !ok {
		t.Fatalf("expected /v1/books path, got: %v", paths)
	}
	get, _ := listOp["get"].(map[string]any)
	if get == nil {
		t.Fatalf("expected GET operation under /v1/books")
	}

	// 2. Query parameters are emitted for the request fields, with nested
	//    message fields flattened to dotted names. The fixture's
	//    ListBooksRequest carries:
	//      - scalar `pageSize`, `pageToken`
	//      - nested `filter` (BookFilter) → flattens to `filter.kind`,
	//        `filter.minPages`, and recursively `filter.range.start`/`.end`
	//      - WKT `updatedAfter` (Timestamp) → emitted inline as a single
	//        string parameter, NOT flattened into seconds/nanos
	//      - repeated scalar `tags` → emitted as a single array-typed param
	params, _ := get["parameters"].([]any)
	gotParams := map[string]map[string]any{}
	for _, p := range params {
		m, _ := p.(map[string]any)
		if m["in"] != "query" {
			t.Errorf("parameter %v: want in=query, got %v", m["name"], m["in"])
		}
		if name, ok := m["name"].(string); ok {
			gotParams[name] = m
		}
	}
	for _, name := range []string{
		"pageSize", "pageToken",
		"filter.kind", "filter.minPages",
		"filter.range.start", "filter.range.end",
		"filter.labels[string]",
		"updatedAfter", "tags",
	} {
		if _, ok := gotParams[name]; !ok {
			t.Errorf("missing query parameter %q", name)
		}
	}
	// `filter` itself must NOT appear as its own parameter — it has been
	// flattened away. Same for `filter.range`.
	for _, name := range []string{"filter", "filter.range"} {
		if _, ok := gotParams[name]; ok {
			t.Errorf("nested message %q should be flattened, not emitted as a parameter", name)
		}
	}
	// Spot-check the flattened parameter shapes:
	//   filter.kind reuses the Book.Kind enum schema via $ref.
	if schema, _ := gotParams["filter.kind"]["schema"].(map[string]any); schema["$ref"] != "#/components/schemas/lib.v1.Book.Kind" {
		t.Errorf("filter.kind: want schema $ref to Book.Kind, got %v", schema)
	}
	//   updatedAfter is the WKT inline form (string/date-time), not a $ref
	//   into Timestamp's component, and not exploded into seconds/nanos.
	if schema, _ := gotParams["updatedAfter"]["schema"].(map[string]any); schema["type"] != "string" || schema["format"] != "date-time" {
		t.Errorf("updatedAfter: want inline string/date-time WKT schema, got %v", schema)
	}
	//   tags is a repeated scalar — single param with array schema.
	if schema, _ := gotParams["tags"]["schema"].(map[string]any); schema["type"] != "array" {
		t.Errorf("tags: want array schema, got %v", schema)
	}
	//   filter.labels is a map<string,string> — emitted as a single
	//   `filter.labels[string]` parameter whose schema is the value type.
	//   The `[string]` token documents the expected key type; the runtime
	//   accepts any real key value between the brackets.
	if schema, _ := gotParams["filter.labels[string]"]["schema"].(map[string]any); schema["type"] != "string" {
		t.Errorf("filter.labels[string]: want value schema type=string, got %v", schema)
	}

	// 3. The Book and ListBooksResponse schemas are emitted as components.
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	for _, name := range []string{"lib.v1.Book", "lib.v1.ListBooksResponse", "lib.v1.Book.Kind", "google.rpc.Status"} {
		if _, ok := schemas[name]; !ok {
			t.Errorf("missing component schema %q", name)
		}
	}

	// 4. The enum schema is a string with the right values. Dual-encoding
	// (also accepting integer numbers) is intentionally not represented
	// here — see ensureEnumSchema for why.
	kind, _ := schemas["lib.v1.Book.Kind"].(map[string]any)
	if kind["type"] != "string" {
		t.Errorf("Book.Kind: want type=string, got %v", kind["type"])
	}
	enumVals, _ := kind["enum"].([]any)
	if len(enumVals) != 3 {
		t.Errorf("Book.Kind: want 3 enum values, got %v", enumVals)
	}

	// 5. The Book.tags field is an array of strings (repeated scalar).
	book, _ := schemas["lib.v1.Book"].(map[string]any)
	props, _ := book["properties"].(map[string]any)
	tags, _ := props["tags"].(map[string]any)
	if tags["type"] != "array" {
		t.Errorf("Book.tags: want type=array, got %v", tags["type"])
	}
	items, _ := tags["items"].(map[string]any)
	if items["type"] != "string" {
		t.Errorf("Book.tags.items: want type=string, got %v", items["type"])
	}

	// 6. The Book.kind field is a $ref to the enum component.
	kindProp, _ := props["kind"].(map[string]any)
	if kindProp["$ref"] != "#/components/schemas/lib.v1.Book.Kind" {
		t.Errorf("Book.kind: want $ref to enum, got %v", kindProp)
	}

	// 7. The ListBooksResponse.books field is a $ref-bearing array.
	listResp, _ := schemas["lib.v1.ListBooksResponse"].(map[string]any)
	listProps, _ := listResp["properties"].(map[string]any)
	books, _ := listProps["books"].(map[string]any)
	if books["type"] != "array" {
		t.Errorf("ListBooksResponse.books: want type=array, got %v", books["type"])
	}
	bookItems, _ := books["items"].(map[string]any)
	if bookItems["$ref"] != "#/components/schemas/lib.v1.Book" {
		t.Errorf("ListBooksResponse.books.items: want $ref to Book, got %v", bookItems)
	}

	// 8. Book.published_at is a WKT Timestamp — inlined as string+date-time,
	// not a $ref. This is the property the WKT fast-path guards.
	published, _ := props["publishedAt"].(map[string]any)
	if published["type"] != "string" || published["format"] != "date-time" {
		t.Errorf("Book.publishedAt: want inline string/date-time, got %v", published)
	}

	// 9. Book.sequel is a self-reference. The cycle guard in
	// ensureMessageSchema should emit a $ref to Book without recursing.
	sequel, _ := props["sequel"].(map[string]any)
	if sequel["$ref"] != "#/components/schemas/lib.v1.Book" {
		t.Errorf("Book.sequel: want $ref to Book, got %v", sequel)
	}

	// 10. CreateBook has body="*" AND a {shelf} path param. The synthesized
	// body schema must include `book` but NOT `shelf` — this is the exact
	// case the v2 generator and earlier v3 PR get wrong.
	createOp, _ := paths["/v1/shelves/{shelf}/books"].(map[string]any)
	if createOp == nil {
		t.Fatalf("expected /v1/shelves/{shelf}/books path")
	}
	createPost, _ := createOp["post"].(map[string]any)
	body, _ := createPost["requestBody"].(map[string]any)
	content, _ := body["content"].(map[string]any)
	appJSON, _ := content["application/json"].(map[string]any)
	bodySchema, _ := appJSON["schema"].(map[string]any)
	bodyProps, _ := bodySchema["properties"].(map[string]any)
	if _, ok := bodyProps["book"]; !ok {
		t.Errorf("CreateBook body: missing 'book' property, got %v", bodyProps)
	}
	if _, ok := bodyProps["shelf"]; ok {
		t.Errorf("CreateBook body: 'shelf' must be excluded (it's a path param), got %v", bodyProps)
	}
	// And the path parameter must still be emitted as a path param.
	createParams, _ := createPost["parameters"].([]any)
	sawShelf := false
	for _, p := range createParams {
		m, _ := p.(map[string]any)
		if m["name"] == "shelf" && m["in"] == "path" {
			sawShelf = true
		}
	}
	if !sawShelf {
		t.Errorf("CreateBook: missing 'shelf' path parameter, got %v", createParams)
	}

	// 11. DeleteBook returns google.protobuf.Empty. The grpc-gateway runtime
	// writes `{}` with a 200 regardless of the response type, so the spec
	// emits a 200 with an empty object schema (via the WKT fast-path).
	// DeleteBook is also marked `deprecated = true` on the method itself,
	// which must surface as deprecated: true on the operation.
	deleteOp, _ := paths["/v1/books/{id}"].(map[string]any)
	if deleteOp == nil {
		t.Fatalf("expected /v1/books/{id} path")
	}
	deleteMethod, _ := deleteOp["delete"].(map[string]any)
	deleteResponses, _ := deleteMethod["responses"].(map[string]any)
	if _, ok := deleteResponses["200"]; !ok {
		t.Errorf("DeleteBook: want 200 response, got %v", deleteResponses)
	}
	if _, ok := deleteResponses["204"]; ok {
		t.Errorf("DeleteBook: spec must not promise 204 — runtime writes 200, got %v", deleteResponses["204"])
	}
	if got, _ := deleteMethod["deprecated"].(bool); !got {
		t.Errorf("DeleteBook: method-level deprecated must surface as deprecated:true on the operation, got %v", deleteMethod["deprecated"])
	}

	// 12. Book has two oneof groups (`format` and `provenance`). All four
	// oneof fields must appear in `properties`, and the message must carry
	// an allOf with one oneOf per group. Each oneof has a {none} guard plus
	// one per-field option that embeds the variant's actual property schema.
	for _, name := range []string{"paperbackIsbn", "ebookUrl", "boughtAt", "borrowedFrom"} {
		if _, ok := props[name]; !ok {
			t.Errorf("Book.%s: missing from properties", name)
		}
	}
	allOf, _ := book["allOf"].([]any)
	if len(allOf) != 2 {
		t.Fatalf("Book: want allOf with 2 oneOf groups, got %v", book["allOf"])
	}
	wantGroups := []map[string]bool{
		{"paperbackIsbn": true, "ebookUrl": true},
		{"boughtAt": true, "borrowedFrom": true},
	}
	for i, entry := range allOf {
		m, _ := entry.(map[string]any)
		opts, _ := m["oneOf"].([]any)
		// 1 "none" option + one option per field.
		if got, want := len(opts), len(wantGroups[i])+1; got != want {
			t.Errorf("Book oneof group %d: want %d options, got %d", i, want, got)
			continue
		}
		// First option is the "none of the fields are set" guard. It must
		// pin type=object.
		none, _ := opts[0].(map[string]any)
		if got, _ := none["type"].(string); got != "object" {
			t.Errorf("Book oneof group %d: guard type = %v, want object", i, none["type"])
		}
		not, _ := none["not"].(map[string]any)
		if _, ok := not["anyOf"]; !ok {
			t.Errorf("Book oneof group %d: first option should be {not: {anyOf:...}}, got %v", i, none)
		}
		// Remaining options each require exactly one field of the group
		// AND embed that field's typed property schema.
		seen := map[string]bool{}
		for _, opt := range opts[1:] {
			o, _ := opt.(map[string]any)
			if got, _ := o["type"].(string); got != "object" {
				t.Errorf("Book oneof group %d: branch type = %v, want object", i, o["type"])
			}
			req, _ := o["required"].([]any)
			if len(req) != 1 {
				t.Errorf("Book oneof group %d: required must have 1 entry, got %v", i, req)
				continue
			}
			fname := req[0].(string)
			seen[fname] = true
			branchProps, _ := o["properties"].(map[string]any)
			propSchema, _ := branchProps[fname].(map[string]any)
			if propSchema == nil {
				t.Errorf("Book oneof group %d: branch for %q must embed a property schema, got %v", i, fname, branchProps)
				continue
			}
			// All four fields in this test fixture are strings.
			if got, _ := propSchema["type"].(string); got != "string" {
				t.Errorf("Book oneof group %d: branch %q properties.%s.type = %v, want string", i, fname, fname, got)
			}
		}
		for f := range wantGroups[i] {
			if !seen[f] {
				t.Errorf("Book oneof group %d: missing required option for %q", i, f)
			}
		}
	}
	// A single-group case would be hoisted to top-level oneOf instead of
	// allOf; the multi-group case must NOT have a top-level oneOf.
	if _, ok := book["oneOf"]; ok {
		t.Errorf("Book: multi-group message must not have top-level oneOf")
	}
}

// TestGenerate_DeprecatedFileCascade exercises the file-level cascade rule
// for the `deprecated` option: when the proto file itself is marked
// `option deprecated = true`, every method, message, enum, and field
// defined in that file must end up flagged `deprecated: true` in the
// generated spec, without needing per-element opt-in.
func TestGenerate_DeprecatedFileCascade(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/deprecated.prototext")
	got := runGenerator(t, req)

	var doc map[string]any
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, string(got))
	}

	// Method cascade: LegacyService.LegacyRead is not self-deprecated but
	// its enclosing file is, so the operation must be flagged.
	paths, _ := doc["paths"].(map[string]any)
	legacyGet, _ := paths["/v1/legacy"].(map[string]any)
	op, _ := legacyGet["get"].(map[string]any)
	if op == nil {
		t.Fatal("missing GET /v1/legacy operation")
	}
	if got, _ := op["deprecated"].(bool); !got {
		t.Errorf("LegacyRead: want deprecated:true via file cascade, got %v", op["deprecated"])
	}

	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)

	// Message cascade: LegacyMessage is not self-deprecated.
	msg, _ := schemas["legacy.v1.LegacyMessage"].(map[string]any)
	if msg == nil {
		t.Fatal("missing LegacyMessage component schema")
	}
	if got, _ := msg["deprecated"].(bool); !got {
		t.Errorf("LegacyMessage: want deprecated:true via file cascade, got %v", msg["deprecated"])
	}

	// Field cascade: legacy_field and legacy_enum inside LegacyMessage
	// are not self-deprecated. The enum reference triggers an allOf
	// wrapper so the deprecated flag has somewhere to land next to the
	// $ref; the scalar string lands directly on the inline schema.
	msgProps, _ := msg["properties"].(map[string]any)
	scalarProp, _ := msgProps["legacyField"].(map[string]any)
	if got, _ := scalarProp["deprecated"].(bool); !got {
		t.Errorf("LegacyMessage.legacyField: want deprecated:true via file cascade, got %v", scalarProp["deprecated"])
	}
	enumProp, _ := msgProps["legacyEnum"].(map[string]any)
	// For a $ref, fieldDeprecated triggers an allOf wrap carrying the flag.
	if _, isRef := enumProp["$ref"]; isRef {
		t.Errorf("LegacyMessage.legacyEnum: expected allOf wrapper to carry deprecated flag, got bare $ref: %v", enumProp)
	}
	allOfWrap, _ := enumProp["allOf"].([]any)
	if len(allOfWrap) == 0 {
		t.Errorf("LegacyMessage.legacyEnum: expected allOf wrap for deprecated $ref, got %v", enumProp)
	}
	if got, _ := enumProp["deprecated"].(bool); !got {
		t.Errorf("LegacyMessage.legacyEnum: want deprecated:true on allOf wrapper, got %v", enumProp["deprecated"])
	}

	// Enum cascade: LegacyEnum is not self-deprecated.
	enum, _ := schemas["legacy.v1.LegacyEnum"].(map[string]any)
	if enum == nil {
		t.Fatal("missing LegacyEnum component schema")
	}
	if got, _ := enum["deprecated"].(bool); !got {
		t.Errorf("LegacyEnum: want deprecated:true via file cascade, got %v", enum["deprecated"])
	}
}

// TestGenerate_SimpleEcho is the end-to-end smoke test: load a hand-rolled
// CodeGeneratorRequest from prototext, run the generator, and compare the
// JSON output against a golden.
//
// Comparison is semantic (parse both sides into any), not byte-equal: that
// keeps the test stable while we iterate on map iteration order. The next
// step beyond this PoC would be to lock down ordering for components.schemas
// and switch to byte equality.
func TestGenerate_SimpleEcho(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/simple_echo.prototext")
	got := runGenerator(t, req)

	want, err := os.ReadFile("testdata/simple_echo.openapi.json")
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}

	assertJSONEqual(t, got, want)
}

// loadRequest reads a prototext-encoded CodeGeneratorRequest from disk.
func loadRequest(t *testing.T, path string) *pluginpb.CodeGeneratorRequest {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	req := new(pluginpb.CodeGeneratorRequest)
	if err := prototext.Unmarshal(body, req); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	return req
}

// runGenerator runs Generate against a request, returning the JSON content of
// the single emitted file. Assumes exactly one output file; the test fails if
// the generator emits zero or multiple.
func runGenerator(t *testing.T, req *pluginpb.CodeGeneratorRequest) []byte {
	t.Helper()
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(true)
	if err := reg.Load(req); err != nil {
		t.Fatalf("registry load: %v", err)
	}
	var targets []*descriptor.File
	for _, name := range req.FileToGenerate {
		f, err := reg.LookupFile(name)
		if err != nil {
			t.Fatalf("lookup %s: %v", name, err)
		}
		targets = append(targets, f)
	}
	out, err := Generate(reg, targets)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 output file, got %d", len(out))
	}
	return []byte(out[0].GetContent())
}

// assertJSONEqual parses both sides as generic JSON and reports a structural
// diff if they differ. We do not require key order to match.
func assertJSONEqual(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("unmarshal got: %v\n--- got ---\n%s", err, string(got))
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if diff := cmp.Diff(w, g); diff != "" {
		t.Errorf("json mismatch (-want +got):\n%s\n--- raw got ---\n%s", diff, string(got))
	}
}
