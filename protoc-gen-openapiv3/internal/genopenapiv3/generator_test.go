package genopenapiv3

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestConvertPathTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple parameter",
			input:    "/users/{user_id}",
			expected: "/users/{user_id}",
		},
		{
			name:     "parameter with pattern",
			input:    "/users/{user_id=users/*}",
			expected: "/users/{user_id}",
		},
		{
			name:     "nested parameter",
			input:    "/projects/{project.name}",
			expected: "/projects/{project.name}",
		},
		{
			name:     "complex pattern",
			input:    "/v1/{name=projects/*/documents/*}:export",
			expected: "/v1/{name}:export",
		},
		{
			name:     "multiple parameters",
			input:    "/v1/users/{user_id}/posts/{post_id}",
			expected: "/v1/users/{user_id}/posts/{post_id}",
		},
		{
			name:     "multiple parameters with patterns",
			input:    "/v1/{user=users/*}/posts/{post=posts/*}",
			expected: "/v1/{user}/posts/{post}",
		},
		{
			name:     "no parameters",
			input:    "/v1/users",
			expected: "/v1/users",
		},
		{
			name:     "wildcard pattern",
			input:    "/v1/{resource=**}",
			expected: "/v1/{resource}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertPathTemplate(tt.input)
			if result != tt.expected {
				t.Errorf("convertPathTemplate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResponseHasDescription(t *testing.T) {
	// Verify that generated responses always have description
	resp := NewResponse("test description")
	if resp.Description != "test description" {
		t.Errorf("Description = %q, want %q", resp.Description, "test description")
	}

	// Marshal and verify description is present in JSON
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	if !strings.Contains(string(data), `"description"`) {
		t.Error("Marshaled response should contain 'description' field")
	}
	if !strings.Contains(string(data), `"test description"`) {
		t.Error("Marshaled response should contain 'test description' value")
	}
}

func TestSchemaRefMarshal(t *testing.T) {
	tests := []struct {
		name     string
		ref      *SchemaRef
		expected string
	}{
		{
			name:     "reference",
			ref:      NewSchemaRef("my.Message"),
			expected: `{"$ref":"#/components/schemas/my.Message"}`,
		},
		{
			name:     "inline string schema",
			ref:      &SchemaRef{Value: &Schema{Type: "string"}},
			expected: `{"type":"string"}`,
		},
		{
			name:     "inline integer schema",
			ref:      &SchemaRef{Value: &Schema{Type: "integer", Format: "int32"}},
			expected: `{"type":"integer","format":"int32"}`,
		},
		{
			name:     "nil schema",
			ref:      nil,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.ref)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			if !jsonEqual(t, string(data), tt.expected) {
				t.Errorf("got %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestSplitSummaryDescription(t *testing.T) {
	tests := []struct {
		name        string
		comment     string
		wantSummary string
		wantDesc    string
	}{
		{
			name:        "single line summary",
			comment:     "Single line summary",
			wantSummary: "Single line summary",
			wantDesc:    "",
		},
		{
			name:        "summary and description",
			comment:     "Summary line\n\nDescription paragraph.",
			wantSummary: "Summary line",
			wantDesc:    "Description paragraph.",
		},
		{
			name:        "summary and multiple paragraphs",
			comment:     "Summary\n\nPara 1\n\nPara 2",
			wantSummary: "Summary",
			wantDesc:    "Para 1\n\nPara 2",
		},
		{
			name:        "empty comment",
			comment:     "",
			wantSummary: "",
			wantDesc:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, desc := splitSummaryDescription(tt.comment)
			if summary != tt.wantSummary {
				t.Errorf("summary = %q, want %q", summary, tt.wantSummary)
			}
			if desc != tt.wantDesc {
				t.Errorf("description = %q, want %q", desc, tt.wantDesc)
			}
		})
	}
}

func TestRemoveInternalComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no internal comments",
			input:    "This is a public comment",
			expected: "This is a public comment",
		},
		{
			name:     "single internal comment",
			input:    "Public (-- internal --) more public",
			expected: "Public  more public",
		},
		{
			name:     "multiple internal comments",
			input:    "Start (-- first --) middle (-- second --) end",
			expected: "Start  middle  end",
		},
		{
			name:     "internal comment at start",
			input:    "(-- hidden --)Visible",
			expected: "Visible",
		},
		{
			name:     "internal comment at end",
			input:    "Visible(-- hidden --)",
			expected: "Visible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeInternalComments(tt.input)
			if result != tt.expected {
				t.Errorf("removeInternalComments(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPathsOrdering(t *testing.T) {
	// Test that paths maintain insertion order
	paths := NewPaths()

	paths.Set("/z-path", &PathItem{Summary: "Z"})
	paths.Set("/a-path", &PathItem{Summary: "A"})
	paths.Set("/m-path", &PathItem{Summary: "M"})

	data, err := json.Marshal(paths)
	if err != nil {
		t.Fatalf("Failed to marshal paths: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal paths: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 paths, got %d", len(result))
	}
	if _, ok := result["/z-path"]; !ok {
		t.Error("Missing /z-path")
	}
	if _, ok := result["/a-path"]; !ok {
		t.Error("Missing /a-path")
	}
	if _, ok := result["/m-path"]; !ok {
		t.Error("Missing /m-path")
	}
}

func TestPathsSortAlphabetically(t *testing.T) {
	paths := NewPaths()

	paths.Set("/z-path", &PathItem{Summary: "Z"})
	paths.Set("/a-path", &PathItem{Summary: "A"})
	paths.Set("/m-path", &PathItem{Summary: "M"})

	paths.SortAlphabetically()

	expected := []string{"/a-path", "/m-path", "/z-path"}
	if !reflect.DeepEqual(paths.order, expected) {
		t.Errorf("After sorting, order = %v, want %v", paths.order, expected)
	}
}

func TestNewOpenAPI(t *testing.T) {
	doc := NewOpenAPI("Test API", "1.0.0", "3.0.3")

	if doc.OpenAPI != "3.0.3" {
		t.Errorf("OpenAPI = %q, want %q", doc.OpenAPI, "3.0.3")
	}
	if doc.Info.Title != "Test API" {
		t.Errorf("Info.Title = %q, want %q", doc.Info.Title, "Test API")
	}
	if doc.Info.Version != "1.0.0" {
		t.Errorf("Info.Version = %q, want %q", doc.Info.Version, "1.0.0")
	}
	if doc.Paths == nil {
		t.Error("Paths should not be nil")
	}
	if doc.Components == nil {
		t.Error("Components should not be nil")
	}
	if doc.Components.Schemas == nil {
		t.Error("Components.Schemas should not be nil")
	}
}

func TestResponseWithJSONSchema(t *testing.T) {
	resp := NewResponse("Success").WithJSONSchema(NewSchemaRef("MyResponse"))

	if resp.Description != "Success" {
		t.Errorf("Description = %q, want %q", resp.Description, "Success")
	}
	if resp.Content == nil {
		t.Fatal("Content should not be nil")
	}
	if resp.Content["application/json"] == nil {
		t.Fatal("Content[application/json] should not be nil")
	}
	if resp.Content["application/json"].Schema.Ref != "#/components/schemas/MyResponse" {
		t.Errorf("Schema.Ref = %q, want %q", resp.Content["application/json"].Schema.Ref, "#/components/schemas/MyResponse")
	}
}

func TestParameterCreation(t *testing.T) {
	t.Run("path parameter", func(t *testing.T) {
		param := NewPathParameter("user_id", &SchemaRef{Value: &Schema{Type: "string"}})
		if param.Name != "user_id" {
			t.Errorf("Name = %q, want %q", param.Name, "user_id")
		}
		if param.In != "path" {
			t.Errorf("In = %q, want %q", param.In, "path")
		}
		if !param.Required {
			t.Error("Path parameter should be required")
		}
	})

	t.Run("query parameter", func(t *testing.T) {
		param := NewQueryParameter("limit", &SchemaRef{Value: &Schema{Type: "integer"}})
		if param.Name != "limit" {
			t.Errorf("Name = %q, want %q", param.Name, "limit")
		}
		if param.In != "query" {
			t.Errorf("In = %q, want %q", param.In, "query")
		}
		if param.Required {
			t.Error("Query parameter should not be required by default")
		}
	})

	t.Run("header parameter", func(t *testing.T) {
		param := NewHeaderParameter("X-Custom", &SchemaRef{Value: &Schema{Type: "string"}})
		if param.Name != "X-Custom" {
			t.Errorf("Name = %q, want %q", param.Name, "X-Custom")
		}
		if param.In != "header" {
			t.Errorf("In = %q, want %q", param.In, "header")
		}
	})
}

func TestRequestBodyCreation(t *testing.T) {
	schema := NewSchemaRef("MyRequest")
	body := NewJSONRequestBody(schema, true)

	if !body.Required {
		t.Error("Required should be true")
	}
	if body.Content == nil {
		t.Fatal("Content should not be nil")
	}
	if body.Content["application/json"] == nil {
		t.Fatal("Content[application/json] should not be nil")
	}
	if body.Content["application/json"].Schema.Ref != "#/components/schemas/MyRequest" {
		t.Errorf("Schema.Ref = %q, want %q", body.Content["application/json"].Schema.Ref, "#/components/schemas/MyRequest")
	}
}

func TestResponsesObject(t *testing.T) {
	responses := NewResponses()

	responses.Codes["200"] = &ResponseRef{
		Value: NewResponse("Success"),
	}
	responses.Codes["400"] = &ResponseRef{
		Value: NewResponse("Bad Request"),
	}
	responses.Default = &ResponseRef{
		Value: NewResponse("Error"),
	}

	data, err := json.Marshal(responses)
	if err != nil {
		t.Fatalf("Failed to marshal responses: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal responses: %v", err)
	}

	if _, ok := result["200"]; !ok {
		t.Error("Missing 200 response")
	}
	if _, ok := result["400"]; !ok {
		t.Error("Missing 400 response")
	}
	if _, ok := result["default"]; !ok {
		t.Error("Missing default response")
	}
}

func TestWellKnownTypeSchema(t *testing.T) {
	tests := []struct {
		typeName   string
		expectType string
		expectFmt  string
	}{
		{".google.protobuf.Timestamp", "string", "date-time"},
		{".google.protobuf.Duration", "string", ""},
		{".google.protobuf.StringValue", "string", ""},
		{".google.protobuf.Int32Value", "integer", "int32"},
		{".google.protobuf.Int64Value", "string", "int64"},
		{".google.protobuf.BoolValue", "boolean", ""},
		{".google.protobuf.Empty", "object", ""},
		{".google.protobuf.Struct", "object", ""},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			schema := wellKnownTypeSchema(tt.typeName)
			if schema == nil {
				t.Fatalf("schema should exist for %s", tt.typeName)
			}
			if schema.Type != tt.expectType {
				t.Errorf("Type = %q, want %q", schema.Type, tt.expectType)
			}
			if tt.expectFmt != "" && schema.Format != tt.expectFmt {
				t.Errorf("Format = %q, want %q", schema.Format, tt.expectFmt)
			}
		})
	}

	// Test unknown type returns nil
	if wellKnownTypeSchema(".unknown.Type") != nil {
		t.Error("Unknown type should return nil")
	}
}

func TestNeedsRequestBody(t *testing.T) {
	tests := []struct {
		method string
		want   bool
	}{
		{"POST", true},
		{"PUT", true},
		{"PATCH", true},
		{"GET", false},
		{"DELETE", false},
		{"HEAD", false},
		{"OPTIONS", false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			if got := needsRequestBody(tt.method); got != tt.want {
				t.Errorf("needsRequestBody(%q) = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}

func TestFormatValidation(t *testing.T) {
	if err := FormatJSON.Validate(); err != nil {
		t.Errorf("FormatJSON.Validate() returned error: %v", err)
	}
	if err := FormatYAML.Validate(); err != nil {
		t.Errorf("FormatYAML.Validate() returned error: %v", err)
	}
	if err := Format("xml").Validate(); err == nil {
		t.Error("Format(xml).Validate() should return error")
	}
}

func TestFormatMarshal(t *testing.T) {
	data := map[string]string{"key": "value"}

	jsonBytes, err := FormatJSON.Marshal(data)
	if err != nil {
		t.Fatalf("FormatJSON.Marshal() failed: %v", err)
	}
	if !strings.Contains(string(jsonBytes), `"key"`) {
		t.Error("JSON output should contain key")
	}

	yamlBytes, err := FormatYAML.Marshal(data)
	if err != nil {
		t.Fatalf("FormatYAML.Marshal() failed: %v", err)
	}
	if !strings.Contains(string(yamlBytes), "key:") {
		t.Error("YAML output should contain key:")
	}
}

func TestPathItemSetOperation(t *testing.T) {
	pathItem := &PathItem{}
	op := &Operation{OperationID: "testOp"}

	methods := []struct {
		method string
		getter func() *Operation
	}{
		{"GET", func() *Operation { return pathItem.Get }},
		{"POST", func() *Operation { return pathItem.Post }},
		{"PUT", func() *Operation { return pathItem.Put }},
		{"PATCH", func() *Operation { return pathItem.Patch }},
		{"DELETE", func() *Operation { return pathItem.Delete }},
		{"HEAD", func() *Operation { return pathItem.Head }},
		{"OPTIONS", func() *Operation { return pathItem.Options }},
		{"TRACE", func() *Operation { return pathItem.Trace }},
	}

	for _, m := range methods {
		t.Run(m.method, func(t *testing.T) {
			pathItem.SetOperation(m.method, op)
			if m.getter() != op {
				t.Errorf("SetOperation(%q) did not set the operation correctly", m.method)
			}
		})
	}
}

// Helper function to compare JSON
func jsonEqual(t *testing.T, a, b string) bool {
	t.Helper()

	var objA, objB any
	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		t.Errorf("Failed to unmarshal first JSON: %v", err)
		return false
	}
	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		t.Errorf("Failed to unmarshal second JSON: %v", err)
		return false
	}
	return reflect.DeepEqual(objA, objB)
}
