package merge

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

// minimal returns a minimal valid OpenAPI 3.1 document body with the given
// path and component schema. Used by the table-driven cases below to keep
// fixtures compact.
func minimal(title, pathKey, schemaName string) string {
	return `{
  "openapi": "3.1.0",
  "info": { "title": "` + title + `", "version": "1.0.0" },
  "paths": {
    "` + pathKey + `": {
      "get": { "operationId": "Op_` + schemaName + `", "responses": { "200": { "description": "ok" } } }
    }
  },
  "components": { "schemas": { "` + schemaName + `": { "type": "object" } } }
}`
}

func TestMerge_TwoDisjointInputs(t *testing.T) {
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(minimal("a", "/a", "pkg.A"))},
		{Name: "b.json", Data: []byte(minimal("b", "/b", "pkg.B"))},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	// info comes from the first input
	if title := got["info"].(map[string]any)["title"]; title != "a" {
		t.Errorf("info.title: got %q, want %q", title, "a")
	}
	paths := got["paths"].(map[string]any)
	if _, ok := paths["/a"]; !ok {
		t.Errorf("paths missing /a; got %v", paths)
	}
	if _, ok := paths["/b"]; !ok {
		t.Errorf("paths missing /b; got %v", paths)
	}
	schemas := got["components"].(map[string]any)["schemas"].(map[string]any)
	if _, ok := schemas["pkg.A"]; !ok {
		t.Errorf("schemas missing pkg.A")
	}
	if _, ok := schemas["pkg.B"]; !ok {
		t.Errorf("schemas missing pkg.B")
	}
}

func TestMerge_PathOrder(t *testing.T) {
	// Paths must appear in input order: a.json's paths first, b.json's after.
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(minimal("a", "/zzz", "pkg.A"))},
		{Name: "b.json", Data: []byte(minimal("b", "/aaa", "pkg.B"))},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	zzzIdx := bytes.Index(out, []byte(`"/zzz"`))
	aaaIdx := bytes.Index(out, []byte(`"/aaa"`))
	if zzzIdx < 0 || aaaIdx < 0 || zzzIdx > aaaIdx {
		t.Errorf("paths out of input order; /zzz at %d, /aaa at %d:\n%s", zzzIdx, aaaIdx, out)
	}
}

func TestMerge_ComponentsSorted(t *testing.T) {
	// components.schemas keys must come out sorted regardless of input order.
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(minimal("a", "/a", "z.Z"))},
		{Name: "b.json", Data: []byte(minimal("b", "/b", "a.A"))},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	aIdx := bytes.Index(out, []byte(`"a.A"`))
	zIdx := bytes.Index(out, []byte(`"z.Z"`))
	if aIdx < 0 || zIdx < 0 || aIdx > zIdx {
		t.Errorf("components not sorted; a.A at %d, z.Z at %d:\n%s", aIdx, zIdx, out)
	}
}

func TestMerge_PathCollisionConflictingValues(t *testing.T) {
	a := minimal("a", "/v1/echo", "pkg.A")
	b := minimal("b", "/v1/echo", "pkg.B")
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on path collision, got nil")
	}
	if !strings.Contains(err.Error(), "/v1/echo") {
		t.Errorf("error should mention the conflicting path; got: %v", err)
	}
}

func TestMerge_PathCollisionIdenticalValues(t *testing.T) {
	// Same path with identical content should merge cleanly.
	a := minimal("a", "/v1/echo", "pkg.X")
	b := minimal("b", "/v1/echo", "pkg.X")
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("expected identical paths to merge, got: %v", err)
	}
}

func TestMerge_ComponentCollisionConflictingValues(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "components": { "schemas": { "pkg.User": { "type": "object", "properties": { "id": { "type": "string" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "components": { "schemas": { "pkg.User": { "type": "object", "properties": { "name": { "type": "string" } } } } }
}`
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on component collision")
	}
	if !strings.Contains(err.Error(), "pkg.User") {
		t.Errorf("error should mention conflicting component; got: %v", err)
	}
}

func TestMerge_ComponentCollisionKeyOrderInsensitive(t *testing.T) {
	// Two schemas that differ only in key order must compare equal.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "components": { "schemas": { "pkg.User": { "type": "object", "description": "u" } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "components": { "schemas": { "pkg.User": { "description": "u", "type": "object" } } }
}`
	if _, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	}); err != nil {
		t.Errorf("expected canonical equality; got: %v", err)
	}
}

func TestMerge_TagDedup(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "tags": ["Users"], "responses": { "200": { "description": "ok" } } } } },
  "tags": [{ "name": "Users", "description": "manage users" }]
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "tags": ["Users"], "responses": { "200": { "description": "ok" } } } } },
  "tags": [{ "name": "Users", "description": "manage users" }]
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	tags := got["tags"].([]any)
	if len(tags) != 1 {
		t.Errorf("expected 1 deduplicated tag, got %d: %v", len(tags), tags)
	}
}

func TestMerge_TagConflict(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "tags": [{ "name": "Users", "description": "first" }]
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "tags": [{ "name": "Users", "description": "second" }]
}`
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on conflicting tag metadata")
	}
	if !strings.Contains(err.Error(), "Users") {
		t.Errorf("error should mention the conflicting tag; got: %v", err)
	}
}

func TestMerge_OpenAPIVersionMismatch(t *testing.T) {
	a := `{"openapi":"3.1.0","info":{"title":"a","version":"1.0.0"},"paths":{}}`
	b := `{"openapi":"3.2.0","info":{"title":"b","version":"1.0.0"},"paths":{}}`
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on openapi version mismatch")
	}
	if !strings.Contains(err.Error(), "openapi") {
		t.Errorf("error should mention the openapi field; got: %v", err)
	}
}

func TestMerge_InfoFirstWins(t *testing.T) {
	// Different info blocks across inputs must not error — the first input
	// wins quietly. protoc-gen-openapiv3 derives info.title from each file's
	// name, so disagreement is the common case.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "from-a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "from-b", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid output JSON: %v\n%s", err, out)
	}
	if title := got["info"].(map[string]any)["title"]; title != "from-a" {
		t.Errorf("info.title: got %q, want %q (first input)", title, "from-a")
	}
}

func TestMerge_InfoIdentical(t *testing.T) {
	// Same info on both inputs is the common case for protoc-gen-openapiv3
	// output sharing an `openapiv3_document` annotation.
	body := `{
  "openapi": "3.1.0",
  "info": { "title": "shared", "version": "1.0.0" },
  "paths": { "/PATH": { "get": { "operationId": "Op_X", "responses": { "200": { "description": "ok" } } } } }
}`
	a := strings.ReplaceAll(body, "PATH", "a")
	b := strings.ReplaceAll(body, "PATH", "b")
	if _, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	}); err != nil {
		t.Errorf("identical info should merge cleanly; got: %v", err)
	}
}

func TestMerge_NoInputs(t *testing.T) {
	if _, err := Merge(nil); err == nil {
		t.Error("expected error on empty input list")
	}
}

func TestMerge_SingleInputPassthrough(t *testing.T) {
	in := minimal("only", "/only", "pkg.Only")
	out, err := Merge([]Input{{Name: "only.json", Data: []byte(in)}})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid output JSON: %v\n%s", err, out)
	}
	if got["info"].(map[string]any)["title"] != "only" {
		t.Errorf("info.title not preserved")
	}
	if _, ok := got["paths"].(map[string]any)["/only"]; !ok {
		t.Errorf("paths not preserved")
	}
}

func TestMerge_ExtrasFirstWins(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "x-internal-id": "alpha",
  "x-team": "platform",
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "x-internal-id": "alpha",
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if !bytes.Contains(out, []byte(`"x-internal-id": "alpha"`)) {
		t.Errorf("x-internal-id not in output:\n%s", out)
	}
	if !bytes.Contains(out, []byte(`"x-team": "platform"`)) {
		t.Errorf("x-team not in output:\n%s", out)
	}
}

func TestMerge_ExtrasFirstWinsOnConflict(t *testing.T) {
	// Conflicting top-level extensions: first input's value wins silently.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "x-internal-id": "alpha",
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "x-internal-id": "beta",
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if id := got["x-internal-id"]; id != "alpha" {
		t.Errorf("x-internal-id: got %q, want %q (first input)", id, "alpha")
	}
}

func TestMerge_SecurityIdentical(t *testing.T) {
	// Identical security arrays merge cleanly under first-wins.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "security": [{ "bearer": [] }, { "oauth2": ["read"] }]
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "security": [{ "bearer": [] }, { "oauth2": ["read"] }]
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	sec := got["security"].([]any)
	if len(sec) != 2 {
		t.Errorf("expected 2 security requirements, got %d: %v", len(sec), sec)
	}
}

func TestMerge_SecurityFirstOnly(t *testing.T) {
	// First input's security is kept when later inputs don't declare any.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "security": [{ "bearer": [] }]
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	sec, ok := got["security"].([]any)
	if !ok || len(sec) != 1 {
		t.Errorf("expected first input's security to be preserved; got %v", got["security"])
	}
}

func TestMerge_SecurityLaterOnly(t *testing.T) {
	// Later inputs may introduce security when the first didn't declare any.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "security": [{ "bearer": [] }]
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	sec, ok := got["security"].([]any)
	if !ok || len(sec) != 1 {
		t.Errorf("expected later input's security to land in output; got %v", got["security"])
	}
}

func TestMerge_SecurityConflict(t *testing.T) {
	// Two inputs declaring non-empty, non-identical security arrays is an
	// error: the root security list governs the whole API and silently
	// keeping one would change what's allowed.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "security": [{ "bearer": [] }]
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "security": [{ "bearer": [] }, { "oauth2": ["read"] }]
}`
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on conflicting security definitions")
	}
	if !strings.Contains(err.Error(), "security") {
		t.Errorf("error should mention security; got: %v", err)
	}
}

func TestMerge_TopLevelFieldOrder(t *testing.T) {
	// Output must be in OpenAPI 3.1 declaration order regardless of input
	// order. Decode the top-level keys to compare against the expected order.
	in := `{
  "openapi": "3.1.0",
  "info": { "title": "x", "version": "1" },
  "tags": [{ "name": "T" }],
  "paths": { "/x": { "get": { "operationId": "Op_x", "tags": ["T"], "responses": { "200": { "description": "ok" } } } } },
  "components": { "schemas": { "Foo": { "type": "object" } } }
}`
	out, err := Merge([]Input{{Name: "x.json", Data: []byte(in)}})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	keys := topLevelKeys(t, out)
	want := []string{"openapi", "info", "paths", "components", "tags"}
	if !reflectDeepEqualStrings(keys, want) {
		t.Errorf("top-level field order:\n got: %v\nwant: %v\nout:\n%s", keys, want, out)
	}
}

// topLevelKeys decodes a JSON object and returns its top-level keys in
// source order.
func topLevelKeys(t *testing.T, data []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		t.Fatalf("expected JSON object, got %v", tok)
	}
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("decode key: %v", err)
		}
		keys = append(keys, tok.(string))
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			t.Fatalf("decode value: %v", err)
		}
	}
	return keys
}

func reflectDeepEqualStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMerge_WebhooksUnioned(t *testing.T) {
	// `webhooks` follow the same rules as `paths`: unioned in input order,
	// conflicting keys are an error, identical keys merge.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "webhooks": { "userCreated": { "post": { "operationId": "Hook_UserCreated", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "webhooks": { "orderShipped": { "post": { "operationId": "Hook_OrderShipped", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	hooks, ok := got["webhooks"].(map[string]any)
	if !ok {
		t.Fatalf("webhooks missing from output: %v", got)
	}
	if _, ok := hooks["userCreated"]; !ok {
		t.Errorf("webhooks missing userCreated")
	}
	if _, ok := hooks["orderShipped"]; !ok {
		t.Errorf("webhooks missing orderShipped")
	}
}

func TestMerge_WebhooksConflict(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "webhooks": { "evt": { "post": { "operationId": "Hook_v1", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "webhooks": { "evt": { "post": { "operationId": "Hook_v2", "responses": { "200": { "description": "ok" } } } } }
}`
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on conflicting webhook entries")
	}
	if !strings.Contains(err.Error(), "webhooks") || !strings.Contains(err.Error(), "evt") {
		t.Errorf("error should mention webhooks.evt; got: %v", err)
	}
}

func TestMerge_ExternalDocsFirstWins(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "externalDocs": { "url": "https://example.com/a", "description": "from a" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "externalDocs": { "url": "https://example.com/b", "description": "from b" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	ext, ok := got["externalDocs"].(map[string]any)
	if !ok {
		t.Fatalf("externalDocs missing from output: %v", got)
	}
	if url := ext["url"]; url != "https://example.com/a" {
		t.Errorf("externalDocs.url: got %q, want first input's value", url)
	}
}

func TestMerge_ServersFirstWins(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "servers": [{ "url": "https://api-a.example.com" }],
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "servers": [{ "url": "https://api-b.example.com" }],
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	servers, ok := got["servers"].([]any)
	if !ok || len(servers) != 1 {
		t.Fatalf("servers missing or wrong length: %v", got["servers"])
	}
	if url := servers[0].(map[string]any)["url"]; url != "https://api-a.example.com" {
		t.Errorf("servers[0].url: got %q, want first input's value", url)
	}
}

func TestMerge_MissingInfo(t *testing.T) {
	// OpenAPI 3.1 requires `info`. Reject inputs lacking it.
	in := `{
  "openapi": "3.1.0",
  "paths": { "/x": { "get": { "operationId": "Op_x", "responses": { "200": { "description": "ok" } } } } }
}`
	_, err := Merge([]Input{{Name: "noinfo.json", Data: []byte(in)}})
	if err == nil {
		t.Fatal("expected error when info is missing")
	}
	if !strings.Contains(err.Error(), "info") {
		t.Errorf("error should mention info; got: %v", err)
	}
}

func TestMerge_ComponentsSecuritySchemes(t *testing.T) {
	// Components sub-maps other than `schemas` must merge with the same
	// rules: union, sorted keys on output, conflicts rejected.
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "components": {
    "securitySchemes": {
      "bearer": { "type": "http", "scheme": "bearer" }
    }
  }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "components": {
    "securitySchemes": {
      "apiKey": { "type": "apiKey", "in": "header", "name": "X-API-Key" }
    }
  }
}`
	out, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	schemes := got["components"].(map[string]any)["securitySchemes"].(map[string]any)
	if _, ok := schemes["bearer"]; !ok {
		t.Errorf("securitySchemes missing bearer: %v", schemes)
	}
	if _, ok := schemes["apiKey"]; !ok {
		t.Errorf("securitySchemes missing apiKey: %v", schemes)
	}
	// Keys must come out sorted (apiKey before bearer).
	apiIdx := bytes.Index(out, []byte(`"apiKey"`))
	bearerIdx := bytes.Index(out, []byte(`"bearer"`))
	if apiIdx < 0 || bearerIdx < 0 || apiIdx > bearerIdx {
		t.Errorf("securitySchemes keys not sorted; apiKey at %d, bearer at %d", apiIdx, bearerIdx)
	}
}

func TestMerge_ComponentsSecuritySchemesConflict(t *testing.T) {
	a := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/a": { "get": { "operationId": "Op_a", "responses": { "200": { "description": "ok" } } } } },
  "components": { "securitySchemes": { "bearer": { "type": "http", "scheme": "bearer" } } }
}`
	b := `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": { "/b": { "get": { "operationId": "Op_b", "responses": { "200": { "description": "ok" } } } } },
  "components": { "securitySchemes": { "bearer": { "type": "http", "scheme": "basic" } } }
}`
	_, err := Merge([]Input{
		{Name: "a.json", Data: []byte(a)},
		{Name: "b.json", Data: []byte(b)},
	})
	if err == nil {
		t.Fatal("expected error on conflicting securitySchemes")
	}
	if !strings.Contains(err.Error(), "securitySchemes") || !strings.Contains(err.Error(), "bearer") {
		t.Errorf("error should mention securitySchemes.bearer; got: %v", err)
	}
}

// TestMerge_Golden exercises the merger against a realistic fixture and
// compares against a golden file. Run with -update to refresh the golden.
func TestMerge_Golden(t *testing.T) {
	inputs := []Input{}
	for _, name := range []string{"library.openapi.json", "users.openapi.json"} {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		inputs = append(inputs, Input{Name: name, Data: data})
	}
	got, err := Merge(inputs)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	goldenPath := filepath.Join("testdata", "merged.openapi.json")
	if *update {
		if err := os.WriteFile(goldenPath, got, 0644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if !bytes.Equal(bytes.TrimRight(got, "\n"), bytes.TrimRight(want, "\n")) {
		t.Errorf("golden mismatch; run with -update to refresh.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
