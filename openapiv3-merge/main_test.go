package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const docA = `{
  "openapi": "3.1.0",
  "info": { "title": "a", "version": "1.0.0" },
  "paths": {
    "/v1/a": {
      "get": { "operationId": "Op_A", "responses": { "200": { "description": "ok" } } }
    }
  }
}`

const docB = `{
  "openapi": "3.1.0",
  "info": { "title": "b", "version": "1.0.0" },
  "paths": {
    "/v1/b": {
      "get": { "operationId": "Op_B", "responses": { "200": { "description": "ok" } } }
    }
  }
}`

func writeTemp(t *testing.T, name, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

func TestRun_MergesTwoFiles(t *testing.T) {
	a := writeTemp(t, "a.openapi.json", docA)
	b := writeTemp(t, "b.openapi.json", docB)

	var buf bytes.Buffer
	if err := run([]string{a, b}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, buf.String())
	}
	paths := got["paths"].(map[string]any)
	if _, ok := paths["/v1/a"]; !ok {
		t.Errorf("missing /v1/a in output")
	}
	if _, ok := paths["/v1/b"]; !ok {
		t.Errorf("missing /v1/b in output")
	}
	if !bytes.HasSuffix(buf.Bytes(), []byte("\n")) {
		t.Errorf("output should end with a newline")
	}
}

func TestRun_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	err := run(nil, &buf)
	if err == nil {
		t.Fatal("expected error with no arguments")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("error should be a usage message; got: %v", err)
	}
}

func TestRun_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := run([]string{"/no/such/file.json"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing input file")
	}
	if !strings.Contains(err.Error(), "/no/such/file.json") {
		t.Errorf("error should mention the missing path; got: %v", err)
	}
}

func TestRun_PropagatesMergeError(t *testing.T) {
	// Two files with conflicting path definitions: the merge package's
	// strictness should surface through run().
	a := writeTemp(t, "a.openapi.json", `{
  "openapi": "3.1.0",
  "info": { "title": "x", "version": "1" },
  "paths": { "/x": { "get": { "operationId": "Op_x1", "responses": { "200": { "description": "ok" } } } } }
}`)
	b := writeTemp(t, "b.openapi.json", `{
  "openapi": "3.1.0",
  "info": { "title": "x", "version": "1" },
  "paths": { "/x": { "get": { "operationId": "Op_x2", "responses": { "200": { "description": "ok" } } } } }
}`)
	var buf bytes.Buffer
	if err := run([]string{a, b}, &buf); err == nil {
		t.Fatal("expected error from conflicting path definitions")
	}
}
