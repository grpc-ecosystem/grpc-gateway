package genopenapi

import (
	"encoding/json"
	"os"
	"strings"
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
	// Description on a wrapped $ref must sit on the OUTER schema, next to
	// the deprecated flag, not on the inner $ref entry. Locks down the
	// placement choice so a future refactor doesn't silently move it.
	if enumProp["description"] != "Outer-wrapper description on a deprecated $ref." {
		t.Errorf("LegacyMessage.legacyEnum.description: want annotated value on outer wrapper, got %v", enumProp["description"])
	}
	if len(allOfWrap) > 0 {
		inner, _ := allOfWrap[0].(map[string]any)
		if _, hasInnerDesc := inner["description"]; hasInnerDesc {
			t.Errorf("LegacyMessage.legacyEnum.allOf[0]: description must not be duplicated on inner $ref entry, got %v", inner)
		}
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

// TestGenerate_CustomAnnotations exercises the four openapiv3 custom
// annotations (document, operation, message schema, field schema). It
// locks down that annotation values take precedence over comment-derived
// and default values, and that field-level title force an
// allOf wrapper when the field type is a $ref.
func TestGenerate_CustomAnnotations(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/annotations.prototext")
	got := runGenerator(t, req)

	var doc map[string]any
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, string(got))
	}

	// 1. File-level openapiv3_document replaces the default Info fields and
	// installs servers and externalDocs.
	info, _ := doc["info"].(map[string]any)
	for k, want := range map[string]string{
		"title":          "Annotated API",
		"summary":        "An API demonstrating openapiv3 annotations.",
		"description":    "Used by the generator test suite.",
		"version":        "2.0.0",
		"termsOfService": "https://example.com/tos",
	} {
		if got, _ := info[k].(string); got != want {
			t.Errorf("info.%s: want %q, got %q", k, want, got)
		}
	}
	contact, _ := info["contact"].(map[string]any)
	if contact["name"] != "API Team" || contact["email"] != "team@example.com" {
		t.Errorf("info.contact mismatch: %v", contact)
	}
	license, _ := info["license"].(map[string]any)
	if license["name"] != "Apache-2.0" || license["identifier"] != "Apache-2.0" {
		t.Errorf("info.license mismatch: %v", license)
	}
	servers, _ := doc["servers"].([]any)
	if len(servers) != 1 {
		t.Fatalf("servers: want 1, got %v", servers)
	}
	server, _ := servers[0].(map[string]any)
	if server["url"] != "https://api.example.com/v1" || server["description"] != "Production" {
		t.Errorf("servers[0] mismatch: %v", server)
	}
	extDocs, _ := doc["externalDocs"].(map[string]any)
	if extDocs["url"] != "https://example.com/external" {
		t.Errorf("externalDocs mismatch: %v", extDocs)
	}

	// 2. Method-level openapiv3_operation overrides summary/description,
	// replaces tags with the annotated list, sets operationId and externalDocs,
	// and forces deprecated=true.
	paths, _ := doc["paths"].(map[string]any)
	pi, _ := paths["/v1/widgets/{id}"].(map[string]any)
	op, _ := pi["get"].(map[string]any)
	if op == nil {
		t.Fatal("missing GET /v1/widgets/{id} operation")
	}
	if op["summary"] != "Fetch a widget." {
		t.Errorf("operation summary: want annotated value, got %v", op["summary"])
	}
	if op["description"] != "Loads a widget by its unique id." {
		t.Errorf("operation description: want annotated value, got %v", op["description"])
	}
	if op["operationId"] != "getWidgetById" {
		t.Errorf("operationId: want annotated value, got %v", op["operationId"])
	}
	if got, _ := op["deprecated"].(bool); !got {
		t.Errorf("deprecated: want true via annotation, got %v", op["deprecated"])
	}
	tags, _ := op["tags"].([]any)
	if len(tags) != 2 || tags[0] != "Widgets" || tags[1] != "Inventory" {
		t.Errorf("tags: want [Widgets Inventory] from annotation, got %v", tags)
	}
	opExtDocs, _ := op["externalDocs"].(map[string]any)
	if opExtDocs["url"] != "https://example.com/docs/get-widget" {
		t.Errorf("operation externalDocs: want annotated URL, got %v", opExtDocs)
	}
	opServers, _ := op["servers"].([]any)
	if len(opServers) != 1 {
		t.Fatalf("operation.servers: want 1 entry from annotation, got %v", op["servers"])
	}
	opServer, _ := opServers[0].(map[string]any)
	if opServer["url"] != "https://api.example.com/widgets" || opServer["description"] != "Widget shard." {
		t.Errorf("operation.servers[0] mismatch: %v", opServer)
	}

	// 3. Message-level openapiv3_schema sets title and description on the
	// Widget component schema.
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	widget, _ := schemas["ann.v1.Widget"].(map[string]any)
	if widget == nil {
		t.Fatal("missing ann.v1.Widget component schema")
	}
	if widget["title"] != "Widget" {
		t.Errorf("Widget.title: want annotated value, got %v", widget["title"])
	}
	if widget["description"] != "A widget, annotated." {
		t.Errorf("Widget.description: want annotated value, got %v", widget["description"])
	}

	detail, _ := schemas["ann.v1.WidgetDetail"].(map[string]any)
	if detail["title"] != "Widget Detail" {
		t.Errorf("WidgetDetail.title: want annotated value, got %v", detail["title"])
	}
	if got, _ := detail["deprecated"].(bool); !got {
		t.Errorf("WidgetDetail.deprecated: want true via annotation, got %v", detail["deprecated"])
	}

	// 4. Field-level openapiv3_field annotations attach description and
	// title to the inline property schema.
	widgetProps, _ := widget["properties"].(map[string]any)
	idProp, _ := widgetProps["id"].(map[string]any)
	if idProp["description"] != "Unique widget identifier." {
		t.Errorf("Widget.id.description: want annotated value, got %v", idProp["description"])
	}
	countProp, _ := widgetProps["count"].(map[string]any)
	if countProp["title"] != "Count" {
		t.Errorf("Widget.count.title: want annotated value, got %v", countProp["title"])
	}

	// 5. Field annotation on a $ref-typed field (detail) carries description
	// via the $ref sibling; title forces an allOf wrapper.
	detailProp, _ := widgetProps["detail"].(map[string]any)
	if _, isDirectRef := detailProp["$ref"]; isDirectRef {
		t.Errorf("Widget.detail: want allOf wrapper since annotation sets title, got bare $ref: %v", detailProp)
	}
	if detailProp["title"] != "Detail" {
		t.Errorf("Widget.detail.title: want annotated value on wrapper, got %v", detailProp["title"])
	}
	if detailProp["description"] != "Override for the referenced detail." {
		t.Errorf("Widget.detail.description: want annotated value on wrapper, got %v", detailProp["description"])
	}
	allOfWrap, _ := detailProp["allOf"].([]any)
	if len(allOfWrap) != 1 {
		t.Fatalf("Widget.detail.allOf: want 1 entry, got %v", detailProp["allOf"])
	}
	allOfEntry, _ := allOfWrap[0].(map[string]any)
	if allOfEntry["$ref"] != "#/components/schemas/ann.v1.WidgetDetail" {
		t.Errorf("Widget.detail.allOf[0]: want $ref to WidgetDetail, got %v", allOfEntry)
	}

	// 6. A $ref-typed field without an annotation remains a plain $ref.
	plainProp, _ := widgetProps["plain"].(map[string]any)
	if plainProp["$ref"] != "#/components/schemas/ann.v1.WidgetDetail" {
		t.Errorf("Widget.plain: want bare $ref, got %v", plainProp)
	}

	// 7. A $ref-typed field whose annotation sets only `description` must
	// stay a plain $ref with description as a sibling — description is
	// allowed next to $ref in OpenAPI 3.1.0, so no allOf wrapper is needed.
	refDescProp, _ := widgetProps["refDescOnly"].(map[string]any)
	if refDescProp["$ref"] != "#/components/schemas/ann.v1.WidgetDetail" {
		t.Errorf("Widget.refDescOnly: want plain $ref with description sibling, got %v", refDescProp)
	}
	if refDescProp["description"] != "Description-only override for a referenced field." {
		t.Errorf("Widget.refDescOnly.description: want annotated value as $ref sibling, got %v", refDescProp["description"])
	}
	if _, hasAllOf := refDescProp["allOf"]; hasAllOf {
		t.Errorf("Widget.refDescOnly: description-only annotation must not force allOf wrapper, got %v", refDescProp)
	}

	// 7a. A field annotation with `deprecated: true` on an inline scalar
	// flips the inline schema's deprecated flag without forcing a wrapper.
	legacyProp, _ := widgetProps["legacy"].(map[string]any)
	if got, _ := legacyProp["deprecated"].(bool); !got {
		t.Errorf("Widget.legacy.deprecated: want true via annotation, got %v", legacyProp["deprecated"])
	}
	if _, hasAllOf := legacyProp["allOf"]; hasAllOf {
		t.Errorf("Widget.legacy: deprecated-on-scalar annotation must not force allOf wrapper, got %v", legacyProp)
	}

	// 7b. A field annotation with `deprecated: true` on a $ref-typed field
	// forces an allOf wrapper carrying the flag, even with no title.
	legacyRefProp, _ := widgetProps["legacyRef"].(map[string]any)
	if _, isDirectRef := legacyRefProp["$ref"]; isDirectRef {
		t.Errorf("Widget.legacyRef: want allOf wrapper since annotation sets deprecated, got bare $ref: %v", legacyRefProp)
	}
	if got, _ := legacyRefProp["deprecated"].(bool); !got {
		t.Errorf("Widget.legacyRef.deprecated: want true on allOf wrapper, got %v", legacyRefProp["deprecated"])
	}
	legacyRefAllOf, _ := legacyRefProp["allOf"].([]any)
	if len(legacyRefAllOf) != 1 {
		t.Fatalf("Widget.legacyRef.allOf: want 1 entry, got %v", legacyRefProp["allOf"])
	}
	legacyRefEntry, _ := legacyRefAllOf[0].(map[string]any)
	if legacyRefEntry["$ref"] != "#/components/schemas/ann.v1.WidgetDetail" {
		t.Errorf("Widget.legacyRef.allOf[0]: want $ref to WidgetDetail, got %v", legacyRefEntry)
	}

	// 8. Document-level tags merge with service-derived tags. The fixture
	// declares three tags ("Widgets", "Inventory", "WidgetService"); the
	// last collides with the service name and must suppress the default
	// service tag so the annotation's description wins. The operation's
	// own tag override ("Widgets", "Inventory") is satisfied by the
	// top-level entries, so no tag referenced by the operation is orphaned.
	topTags, _ := doc["tags"].([]any)
	byName := map[string]map[string]any{}
	for _, tg := range topTags {
		m, _ := tg.(map[string]any)
		if n, _ := m["name"].(string); n != "" {
			byName[n] = m
		}
	}
	for name, want := range map[string]string{
		"Widgets":       "Widget-related operations.",
		"Inventory":     "Inventory tracking.",
		"WidgetService": "Override for the service-derived tag.",
	} {
		got, ok := byName[name]
		if !ok {
			t.Errorf("doc.tags: missing %q, got %v", name, byName)
			continue
		}
		if got["description"] != want {
			t.Errorf("doc.tags[%q].description: want %q, got %v", name, want, got["description"])
		}
	}
	// Each name must appear exactly once — the WidgetService default must
	// not duplicate the annotation entry.
	counts := map[string]int{}
	for _, tg := range topTags {
		m, _ := tg.(map[string]any)
		if n, _ := m["name"].(string); n != "" {
			counts[n]++
		}
	}
	for name, c := range counts {
		if c != 1 {
			t.Errorf("doc.tags[%q]: want 1 entry, got %d", name, c)
		}
	}
	// The "Widgets" tag carries externalDocs from the annotation.
	widgetsTag := byName["Widgets"]
	widgetsExt, _ := widgetsTag["externalDocs"].(map[string]any)
	if widgetsExt["url"] != "https://example.com/docs/widgets" {
		t.Errorf("doc.tags[Widgets].externalDocs: want annotated URL, got %v", widgetsExt)
	}
}

// TestGenerate_AnnotationErrors covers the spec-required-field validations
// and the cross-operation invariants (operationId uniqueness, tag references
// resolving). One sub-case per failure mode keeps the error messages
// readable when one of these fixtures drifts.
func TestGenerate_AnnotationErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		fixture string
		// substring expected in the returned error message.
		wantErr string
	}{
		{
			name:    "license missing name",
			fixture: "testdata/license_missing_name.prototext",
			wantErr: "license: name is required",
		},
		{
			name:    "document server missing url",
			fixture: "testdata/server_missing_url.prototext",
			wantErr: "servers[0]: server: url is required",
		},
		{
			name:    "document external_docs missing url",
			fixture: "testdata/external_docs_missing_url.prototext",
			wantErr: "external_docs: url is required",
		},
		{
			name:    "document tag missing name",
			fixture: "testdata/tag_missing_name.prototext",
			wantErr: "tags[0]: name is required",
		},
		{
			name:    "duplicate operationId",
			fixture: "testdata/duplicate_operation_id.prototext",
			wantErr: `duplicate operationId "duplicateId"`,
		},
		{
			name:    "operation references undeclared tag",
			fixture: "testdata/orphan_operation_tag.prototext",
			wantErr: `references undeclared tag "Undeclared"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := loadRequest(t, tc.fixture)
			gotErr := runGeneratorExpectError(t, req)
			if !strings.Contains(gotErr, tc.wantErr) {
				t.Errorf("error: want substring %q, got %q", tc.wantErr, gotErr)
			}
		})
	}
}

// TestGenerate_EmptyAnnotations locks down the "all sub-fields empty means
// no-op" semantics for each of the four openapiv3 extensions. An empty
// annotation is indistinguishable from no annotation at all as far as the
// generated document is concerned — every value must come from the
// generator defaults.
func TestGenerate_EmptyAnnotations(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/empty_annotations.prototext")
	got := runGenerator(t, req)

	var doc map[string]any
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, string(got))
	}

	// Document: empty openapiv3_document leaves Info.title = file basename
	// and version = "1.0.0" (the defaults from NewDocument).
	info, _ := doc["info"].(map[string]any)
	if info["title"] != "empty" {
		t.Errorf("info.title: want file-derived default %q, got %v", "empty", info["title"])
	}
	if info["version"] != "1.0.0" {
		t.Errorf("info.version: want default %q, got %v", "1.0.0", info["version"])
	}
	if _, hasLicense := info["license"]; hasLicense {
		t.Errorf("info.license: empty annotation must not emit a license, got %v", info["license"])
	}
	if _, hasServers := doc["servers"]; hasServers {
		t.Errorf("servers: empty annotation must not emit servers, got %v", doc["servers"])
	}

	// Operation: empty openapiv3_operation leaves the generator-derived
	// operationId, the service-name tag, and deprecated=false (omitted).
	paths, _ := doc["paths"].(map[string]any)
	pi, _ := paths["/v1/things/{id}"].(map[string]any)
	op, _ := pi["get"].(map[string]any)
	if op == nil {
		t.Fatal("missing GET /v1/things/{id} operation")
	}
	if op["operationId"] != "ThingService_GetThing" {
		t.Errorf("operationId: want generator default, got %v", op["operationId"])
	}
	tags, _ := op["tags"].([]any)
	if len(tags) != 1 || tags[0] != "ThingService" {
		t.Errorf("operation.tags: want [ThingService] from service default, got %v", tags)
	}
	if _, hasDeprecated := op["deprecated"]; hasDeprecated {
		t.Errorf("operation.deprecated: empty annotation must not set the flag, got %v", op["deprecated"])
	}
	if _, hasSummary := op["summary"]; hasSummary {
		t.Errorf("operation.summary: want empty (no comments, empty annotation), got %v", op["summary"])
	}

	// Message schema: empty openapiv3_schema leaves title/description unset.
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	thing, _ := schemas["empty.v1.Thing"].(map[string]any)
	if thing == nil {
		t.Fatal("missing empty.v1.Thing component schema")
	}
	if _, hasTitle := thing["title"]; hasTitle {
		t.Errorf("Thing.title: empty annotation must not set title, got %v", thing["title"])
	}
	if _, hasDesc := thing["description"]; hasDesc {
		t.Errorf("Thing.description: empty annotation must not set description, got %v", thing["description"])
	}

	// Field schema: empty openapiv3_field on Thing.id leaves the scalar
	// inline schema untouched (no title, no description, no allOf wrap).
	thingProps, _ := thing["properties"].(map[string]any)
	idProp, _ := thingProps["id"].(map[string]any)
	if idProp["type"] != "string" {
		t.Errorf("Thing.id.type: want string, got %v", idProp["type"])
	}
	if _, hasTitle := idProp["title"]; hasTitle {
		t.Errorf("Thing.id.title: empty annotation must not set title, got %v", idProp["title"])
	}
	if _, hasDesc := idProp["description"]; hasDesc {
		t.Errorf("Thing.id.description: empty annotation must not set description, got %v", idProp["description"])
	}
	if _, hasAllOf := idProp["allOf"]; hasAllOf {
		t.Errorf("Thing.id: empty annotation must not force allOf wrap, got %v", idProp)
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

// TestGenerate_DisableDefaultErrors verifies that when DisableDefaultErrors is
// set on the registry, no default error response (google.rpc.Status) is
// injected into the operations, and the google.rpc.Status component schema is
// also absent from the document.
func TestGenerate_DisableDefaultErrors(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/simple_echo.prototext")

	reg := descriptor.NewRegistry()
	if err := reg.Load(req); err != nil {
		t.Fatalf("registry load: %v", err)
	}
	reg.SetDisableDefaultErrors(true)

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

	var doc map[string]any
	if err := json.Unmarshal([]byte(out[0].GetContent()), &doc); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	// No operation should have a default response when disable_default_errors=true.
	paths, _ := doc["paths"].(map[string]any)
	for urlPath, item := range paths {
		m, _ := item.(map[string]any)
		for method, op := range m {
			o, _ := op.(map[string]any)
			responses, _ := o["responses"].(map[string]any)
			if _, ok := responses["default"]; ok {
				t.Errorf("%s %s: unexpected default response when disable_default_errors=true", method, urlPath)
			}
		}
	}

	// The google.rpc.Status component schema must not be present either.
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	if _, ok := schemas["google.rpc.Status"]; ok {
		t.Errorf("google.rpc.Status component schema must be absent when disable_default_errors=true")
	}
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

// runGeneratorExpectError runs Generate expecting it to fail, returning the
// error string for assertion.
func runGeneratorExpectError(t *testing.T, req *pluginpb.CodeGeneratorRequest) string {
	t.Helper()
	reg := descriptor.NewRegistry()
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
	_, err := Generate(reg, targets)
	if err == nil {
		t.Fatalf("generate: want error, got nil")
	}
	return err.Error()
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
