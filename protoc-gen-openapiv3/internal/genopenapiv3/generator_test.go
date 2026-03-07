package genopenapiv3

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
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
		{
			name:     "multi-line internal comment",
			input:    "Public comment\n(-- api-linter: core::0131::http-body=disabled\n    api-linter: core::0131::http-method=disabled --)\nMore public",
			expected: "Public comment\n\nMore public",
		},
		{
			name:     "buf lint ignore directive",
			input:    "Some description\nbuf:lint:ignore ENUM_VALUE_PREFIX\nbuf:lint:ignore ENUM_VALUE_UPPER_SNAKE_CASE\nMore text",
			expected: "Some description\n\nMore text",
		},
		{
			name:     "buf lint ignore only",
			input:    "buf:lint:ignore FIELD_LOWER_SNAKE_CASE",
			expected: "",
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

func TestGenerateOneOfConstraintsSingleGroup(t *testing.T) {
	// Create a generator with a registry that uses JSON names
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(true)
	g := &generator{reg: reg}

	// Create a parent schema with properties (using JSON names)
	parentSchema := &Schema{
		Type: "object",
		Properties: map[string]*SchemaRef{
			"stringValue":  {Value: &Schema{Type: "string"}},
			"intValue":     {Value: &Schema{Type: "integer"}},
			"regularField": {Value: &Schema{Type: "boolean"}},
		},
	}

	// Create oneof groups
	groups := []oneofGroup{
		{
			name: "value",
			fields: []*descriptor.Field{
				{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     stringPtr("string_value"),
					JsonName: stringPtr("stringValue"),
				}},
				{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     stringPtr("int_value"),
					JsonName: stringPtr("intValue"),
				}},
			},
		},
	}

	oneOf, allOf := g.generateOneOfConstraints(parentSchema, groups)

	// Single group should use oneOf directly, not allOf
	if allOf != nil {
		t.Error("Single group should not use allOf")
	}

	// Should have 3 oneOf options (2 fields + neither)
	if len(oneOf) != 3 {
		t.Fatalf("Expected 3 oneOf schemas (2 fields + neither), got %d", len(oneOf))
	}

	// First option should be for stringValue
	opt1 := oneOf[0].Value
	if opt1.Type != "object" {
		t.Errorf("Option 1 Type = %q, want %q", opt1.Type, "object")
	}
	if len(opt1.Properties) != 1 {
		t.Errorf("Option 1 should have 1 property, got %d", len(opt1.Properties))
	}
	if _, ok := opt1.Properties["stringValue"]; !ok {
		t.Error("Option 1 should have 'stringValue' property")
	}
	if len(opt1.Required) != 1 || opt1.Required[0] != "stringValue" {
		t.Errorf("Option 1 Required = %v, want [stringValue]", opt1.Required)
	}
	if opt1.Title != "value.stringValue" {
		t.Errorf("Option 1 Title = %q, want %q", opt1.Title, "value.stringValue")
	}

	// Second option should be for intValue
	opt2 := oneOf[1].Value
	if _, ok := opt2.Properties["intValue"]; !ok {
		t.Error("Option 2 should have 'intValue' property")
	}
	if len(opt2.Required) != 1 || opt2.Required[0] != "intValue" {
		t.Errorf("Option 2 Required = %v, want [intValue]", opt2.Required)
	}
	if opt2.Title != "value.intValue" {
		t.Errorf("Option 2 Title = %q, want %q", opt2.Title, "value.intValue")
	}

	// Third option should be "neither" - allows none of the oneof fields to be set
	neitherOpt := oneOf[2].Value
	if neitherOpt.Title != "value.none" {
		t.Errorf("Neither option Title = %q, want %q", neitherOpt.Title, "value.none")
	}
	if neitherOpt.Not == nil {
		t.Fatal("Neither option should have 'not' schema")
	}
	if neitherOpt.Not.Value == nil {
		t.Fatal("Neither option 'not' should have value")
	}
	if neitherOpt.Not.Value.AnyOf == nil || len(neitherOpt.Not.Value.AnyOf) != 2 {
		t.Fatalf("Neither option 'not.anyOf' should have 2 entries, got %v", neitherOpt.Not.Value.AnyOf)
	}
}

func TestGenerateOneOfConstraintsMultipleGroups(t *testing.T) {
	// Create a generator with a registry that uses JSON names
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(true)
	g := &generator{reg: reg}

	// Create a parent schema with properties for multiple oneofs
	parentSchema := &Schema{
		Type: "object",
		Properties: map[string]*SchemaRef{
			"createEvent": {Value: &Schema{Type: "object"}},
			"updateEvent": {Value: &Schema{Type: "object"}},
			"error":       {Value: &Schema{Type: "object"}},
			"success":     {Value: &Schema{Type: "object"}},
		},
	}

	// Create multiple oneof groups
	groups := []oneofGroup{
		{
			name: "event",
			fields: []*descriptor.Field{
				{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     stringPtr("create_event"),
					JsonName: stringPtr("createEvent"),
				}},
				{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     stringPtr("update_event"),
					JsonName: stringPtr("updateEvent"),
				}},
			},
		},
		{
			name: "result",
			fields: []*descriptor.Field{
				{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     stringPtr("error"),
					JsonName: stringPtr("error"),
				}},
				{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     stringPtr("success"),
					JsonName: stringPtr("success"),
				}},
			},
		},
	}

	oneOf, allOf := g.generateOneOfConstraints(parentSchema, groups)

	// Multiple groups should use allOf, not oneOf directly
	if oneOf != nil {
		t.Error("Multiple groups should not return direct oneOf")
	}

	// Should have 2 allOf entries (one per group), each containing a oneOf
	if len(allOf) != 2 {
		t.Fatalf("Expected 2 allOf schemas (one per group), got %d", len(allOf))
	}

	// Each allOf entry should wrap a oneOf with 3 options (2 fields + neither)
	for groupIdx, allOfEntry := range allOf {
		if allOfEntry.Value == nil {
			t.Errorf("allOf[%d] should have value", groupIdx)
			continue
		}
		groupOneOf := allOfEntry.Value.OneOf
		if len(groupOneOf) != 3 {
			t.Errorf("Group %d should have 3 oneOf options (2 fields + neither), got %d", groupIdx, len(groupOneOf))
		}
	}

	// Verify group 1 (event) has correct options
	group1OneOf := allOf[0].Value.OneOf
	group1Titles := make(map[string]bool)
	for _, schema := range group1OneOf {
		group1Titles[schema.Value.Title] = true
	}
	for _, expected := range []string{"event.createEvent", "event.updateEvent", "event.none"} {
		if !group1Titles[expected] {
			t.Errorf("Group 1 missing expected oneOf option with title %q", expected)
		}
	}

	// Verify group 2 (result) has correct options
	group2OneOf := allOf[1].Value.OneOf
	group2Titles := make(map[string]bool)
	for _, schema := range group2OneOf {
		group2Titles[schema.Value.Title] = true
	}
	for _, expected := range []string{"result.error", "result.success", "result.none"} {
		if !group2Titles[expected] {
			t.Errorf("Group 2 missing expected oneOf option with title %q", expected)
		}
	}
}

func TestBuildNeitherSetSchema(t *testing.T) {
	g := &generator{}

	schema := g.buildNeitherSetSchema("myGroup", []string{"field1", "field2", "field3"})

	// Should have correct title
	if schema.Title != "myGroup.none" {
		t.Errorf("Title = %q, want %q", schema.Title, "myGroup.none")
	}

	// Should have "not" with "anyOf"
	if schema.Not == nil {
		t.Fatal("Schema should have 'not'")
	}
	if schema.Not.Value == nil {
		t.Fatal("Not should have value")
	}
	if schema.Not.Value.AnyOf == nil {
		t.Fatal("Not should have anyOf")
	}

	// AnyOf should contain required entry for each field
	if len(schema.Not.Value.AnyOf) != 3 {
		t.Errorf("anyOf should have 3 entries, got %d", len(schema.Not.Value.AnyOf))
	}

	requiredFields := make(map[string]bool)
	for _, ref := range schema.Not.Value.AnyOf {
		if ref.Value != nil && len(ref.Value.Required) == 1 {
			requiredFields[ref.Value.Required[0]] = true
		}
	}
	for _, field := range []string{"field1", "field2", "field3"} {
		if !requiredFields[field] {
			t.Errorf("Missing required entry for field %q", field)
		}
	}
}

func TestGenerateOneOfConstraintsEmptyGroups(t *testing.T) {
	// Test edge case: empty groups slice should return nil for both oneOf and allOf
	reg := descriptor.NewRegistry()
	g := &generator{reg: reg}

	parentSchema := &Schema{
		Type:       "object",
		Properties: map[string]*SchemaRef{},
	}

	// Empty groups
	oneOf, allOf := g.generateOneOfConstraints(parentSchema, []oneofGroup{})

	if oneOf != nil {
		t.Errorf("Expected oneOf to be nil for empty groups, got %v", oneOf)
	}
	if allOf != nil {
		t.Errorf("Expected allOf to be nil for empty groups, got %v", allOf)
	}
}

func stringPtr(s string) *string {
	return &s
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

func jsonEqualOrdered(t *testing.T, a, b []byte) bool {
	var bufA, bufB bytes.Buffer
	if err := json.Compact(&bufA, a); err != nil {
		t.Errorf("Failed to compact first JSON: %v", err)
		return false
	}
	if err := json.Compact(&bufB, b); err != nil {
		t.Errorf("Failed to compact second JSON: %v", err)
		return false
	}
	return bytes.Equal(bufA.Bytes(), bufB.Bytes())
}

func TestExtractFieldBehavior(t *testing.T) {
	tests := []struct {
		name     string
		options  *descriptorpb.FieldOptions
		expected int // expected number of behaviors
	}{
		{
			name:     "nil options",
			options:  nil,
			expected: 0,
		},
		{
			name:     "empty options",
			options:  &descriptorpb.FieldOptions{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fd := &descriptorpb.FieldDescriptorProto{
				Name:    stringPtr("test_field"),
				Options: tt.options,
			}
			behaviors := extractFieldBehavior(fd)
			if len(behaviors) != tt.expected {
				t.Errorf("extractFieldBehavior() returned %d behaviors, want %d", len(behaviors), tt.expected)
			}
		})
	}
}

func TestGetFieldBehavior(t *testing.T) {
	// Test with nil options
	field := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:    stringPtr("test_field"),
			Options: nil,
		},
	}
	behaviors := getFieldBehavior(field)
	if len(behaviors) != 0 {
		t.Errorf("getFieldBehavior() returned %d behaviors, want 0", len(behaviors))
	}
}

func TestApplyFieldBehaviorToSchema(t *testing.T) {
	tests := []struct {
		name               string
		useProto3Semantics bool
		proto3Optional     bool
		oneofIndex         *int32
		fieldBehaviors     []annotations.FieldBehavior
		wantRequired       bool
		wantReadOnly       bool
		wantWriteOnly      bool
	}{
		{
			name:               "proto3 semantics disabled",
			useProto3Semantics: false,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     nil,
			wantRequired:       false,
			wantReadOnly:       false,
			wantWriteOnly:      false,
		},
		{
			name:               "proto3 semantics enabled, regular field",
			useProto3Semantics: true,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     nil,
			wantRequired:       true,
			wantReadOnly:       false,
			wantWriteOnly:      false,
		},
		{
			name:               "proto3 semantics enabled, optional field",
			useProto3Semantics: true,
			proto3Optional:     true,
			oneofIndex:         nil,
			fieldBehaviors:     nil,
			wantRequired:       false,
			wantReadOnly:       false,
			wantWriteOnly:      false,
		},
		{
			name:               "proto3 semantics enabled, oneof field",
			useProto3Semantics: true,
			proto3Optional:     false,
			oneofIndex:         int32Ptr(0),
			fieldBehaviors:     nil,
			wantRequired:       false,
			wantReadOnly:       false,
			wantWriteOnly:      false,
		},
		{
			name:               "REQUIRED field_behavior annotation",
			useProto3Semantics: false,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     []annotations.FieldBehavior{annotations.FieldBehavior_REQUIRED},
			wantRequired:       true,
			wantReadOnly:       false,
			wantWriteOnly:      false,
		},
		{
			name:               "OUTPUT_ONLY field_behavior annotation",
			useProto3Semantics: false,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     []annotations.FieldBehavior{annotations.FieldBehavior_OUTPUT_ONLY},
			wantRequired:       false,
			wantReadOnly:       true,
			wantWriteOnly:      false,
		},
		{
			name:               "INPUT_ONLY field_behavior annotation",
			useProto3Semantics: false,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     []annotations.FieldBehavior{annotations.FieldBehavior_INPUT_ONLY},
			wantRequired:       false,
			wantReadOnly:       false,
			wantWriteOnly:      true,
		},
		{
			name:               "OPTIONAL overrides proto3 semantics",
			useProto3Semantics: true,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     []annotations.FieldBehavior{annotations.FieldBehavior_OPTIONAL},
			wantRequired:       false,
			wantReadOnly:       false,
			wantWriteOnly:      false,
		},
		{
			name:               "OUTPUT_ONLY with REQUIRED",
			useProto3Semantics: false,
			proto3Optional:     false,
			oneofIndex:         nil,
			fieldBehaviors:     []annotations.FieldBehavior{annotations.FieldBehavior_OUTPUT_ONLY, annotations.FieldBehavior_REQUIRED},
			wantRequired:       true,
			wantReadOnly:       true,
			wantWriteOnly:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := descriptor.NewRegistry()
			reg.SetUseProto3FieldSemantics(tt.useProto3Semantics)
			g := &generator{reg: reg}

			parentSchema := &Schema{Required: []string{}}
			fieldSchema := &Schema{Type: "string"}

			// Create field options with field_behavior extension if specified
			var opts *descriptorpb.FieldOptions
			if len(tt.fieldBehaviors) > 0 {
				opts = &descriptorpb.FieldOptions{}
				proto.SetExtension(opts, annotations.E_FieldBehavior, tt.fieldBehaviors)
			}

			field := &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:           stringPtr("test_field"),
					JsonName:       stringPtr("testField"),
					Proto3Optional: &tt.proto3Optional,
					OneofIndex:     tt.oneofIndex,
					Options:        opts,
				},
			}

			g.applyFieldBehaviorToSchema(parentSchema, fieldSchema, field)

			// Check if field is marked as required in parent schema
			// Note: fieldName() returns proto name by default (test_field), not JSON name (testField)
			gotRequired := len(parentSchema.Required) > 0 && parentSchema.Required[0] == "test_field"
			if gotRequired != tt.wantRequired {
				t.Errorf("applyFieldBehaviorToSchema() required=%v (parentSchema.Required=%v), want %v",
					gotRequired, parentSchema.Required, tt.wantRequired)
			}
			if fieldSchema.ReadOnly != tt.wantReadOnly {
				t.Errorf("applyFieldBehaviorToSchema() set fieldSchema.ReadOnly=%v, want %v", fieldSchema.ReadOnly, tt.wantReadOnly)
			}
			if fieldSchema.WriteOnly != tt.wantWriteOnly {
				t.Errorf("applyFieldBehaviorToSchema() set fieldSchema.WriteOnly=%v, want %v", fieldSchema.WriteOnly, tt.wantWriteOnly)
			}
		})
	}
}

func TestFieldBehaviorRequiredIntegration(t *testing.T) {
	// Integration test verifying multiple fields with different behaviors
	// result in correct schema.Required array
	reg := descriptor.NewRegistry()
	g := &generator{reg: reg}

	parentSchema := &Schema{Required: []string{}}

	// Field 1: REQUIRED behavior - should be in Required array
	opts1 := &descriptorpb.FieldOptions{}
	proto.SetExtension(opts1, annotations.E_FieldBehavior, []annotations.FieldBehavior{annotations.FieldBehavior_REQUIRED})
	field1 := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:    stringPtr("required_field"),
			Options: opts1,
		},
	}
	fieldSchema1 := &Schema{Type: "string"}
	g.applyFieldBehaviorToSchema(parentSchema, fieldSchema1, field1)

	// Field 2: OUTPUT_ONLY behavior - should NOT be in Required array
	opts2 := &descriptorpb.FieldOptions{}
	proto.SetExtension(opts2, annotations.E_FieldBehavior, []annotations.FieldBehavior{annotations.FieldBehavior_OUTPUT_ONLY})
	field2 := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:    stringPtr("output_only_field"),
			Options: opts2,
		},
	}
	fieldSchema2 := &Schema{Type: "string"}
	g.applyFieldBehaviorToSchema(parentSchema, fieldSchema2, field2)

	// Field 3: No behavior - should NOT be in Required array
	field3 := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:    stringPtr("plain_field"),
			Options: nil,
		},
	}
	fieldSchema3 := &Schema{Type: "string"}
	g.applyFieldBehaviorToSchema(parentSchema, fieldSchema3, field3)

	// Field 4: REQUIRED + OUTPUT_ONLY - should be in Required array
	opts4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(opts4, annotations.E_FieldBehavior, []annotations.FieldBehavior{
		annotations.FieldBehavior_REQUIRED,
		annotations.FieldBehavior_OUTPUT_ONLY,
	})
	field4 := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:    stringPtr("required_output_field"),
			Options: opts4,
		},
	}
	fieldSchema4 := &Schema{Type: "string"}
	g.applyFieldBehaviorToSchema(parentSchema, fieldSchema4, field4)

	// Verify Required array contains exactly the expected fields
	expectedRequired := []string{"required_field", "required_output_field"}
	if !reflect.DeepEqual(parentSchema.Required, expectedRequired) {
		t.Errorf("parentSchema.Required = %v, want %v", parentSchema.Required, expectedRequired)
	}

	// Verify ReadOnly/WriteOnly are set correctly
	if !fieldSchema2.ReadOnly {
		t.Error("output_only_field should have ReadOnly=true")
	}
	if !fieldSchema4.ReadOnly {
		t.Error("required_output_field should have ReadOnly=true")
	}
}

func TestGetFieldRequiredFromBehavior(t *testing.T) {
	tests := []struct {
		name               string
		useProto3Semantics bool
		proto3Optional     bool
		oneofIndex         *int32
		wantRequired       bool
	}{
		{
			name:               "proto3 semantics disabled",
			useProto3Semantics: false,
			proto3Optional:     false,
			oneofIndex:         nil,
			wantRequired:       false,
		},
		{
			name:               "proto3 semantics enabled, regular field",
			useProto3Semantics: true,
			proto3Optional:     false,
			oneofIndex:         nil,
			wantRequired:       true,
		},
		{
			name:               "proto3 semantics enabled, optional field",
			useProto3Semantics: true,
			proto3Optional:     true,
			oneofIndex:         nil,
			wantRequired:       false,
		},
		{
			name:               "proto3 semantics enabled, oneof field",
			useProto3Semantics: true,
			proto3Optional:     false,
			oneofIndex:         int32Ptr(0),
			wantRequired:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := descriptor.NewRegistry()
			reg.SetUseProto3FieldSemantics(tt.useProto3Semantics)
			g := &generator{reg: reg}

			field := &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:           stringPtr("test_field"),
					Proto3Optional: &tt.proto3Optional,
					OneofIndex:     tt.oneofIndex,
					Options:        nil,
				},
			}

			gotRequired := g.getFieldRequiredFromBehavior(field)
			if gotRequired != tt.wantRequired {
				t.Errorf("getFieldRequiredFromBehavior() returned required=%v, want %v", gotRequired, tt.wantRequired)
			}
		})
	}
}

func TestSchemaReadOnlyWriteOnlyFields(t *testing.T) {
	// Test that Schema can serialize readOnly and writeOnly properly
	schema := &Schema{
		Type:      "string",
		ReadOnly:  true,
		WriteOnly: false,
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if !strings.Contains(string(data), `"readOnly":true`) {
		t.Error("Schema should contain readOnly:true")
	}

	schema2 := &Schema{
		Type:      "string",
		WriteOnly: true,
	}

	data2, err := json.Marshal(schema2)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if !strings.Contains(string(data2), `"writeOnly":true`) {
		t.Error("Schema should contain writeOnly:true")
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}

func TestResolveOpenAPINameWithStrategy(t *testing.T) {
	tests := []struct {
		name         string
		strategy     string
		fqmn         string
		allFQMNs     []string
		expectedName string
	}{
		{
			name:         "fqn strategy",
			strategy:     "fqn",
			fqmn:         ".test.v1.User",
			allFQMNs:     []string{".test.v1.User", ".test.v1.Request"},
			expectedName: "test.v1.User",
		},
		{
			name:         "simple strategy unique name",
			strategy:     "simple",
			fqmn:         ".test.v1.User",
			allFQMNs:     []string{".test.v1.User", ".test.v1.Request"},
			expectedName: "User",
		},
		{
			name:         "simple strategy with collision",
			strategy:     "simple",
			fqmn:         ".pkg1.User",
			allFQMNs:     []string{".pkg1.User", ".pkg2.User"},
			expectedName: "pkg1.User",
		},
		{
			name:         "legacy strategy",
			strategy:     "legacy",
			fqmn:         ".test.v1.User",
			allFQMNs:     []string{".test.v1.User", ".test.v1.Request"},
			expectedName: "v1User",
		},
		{
			name:         "package strategy",
			strategy:     "package",
			fqmn:         ".test.v1.Outer.Inner",
			allFQMNs:     []string{".test.v1.Outer.Inner"},
			expectedName: "Outer.Inner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveFullyQualifiedNameToOpenAPINames(tt.allFQMNs, tt.strategy)
			if result[tt.fqmn] != tt.expectedName {
				t.Errorf("resolveFullyQualifiedNameToOpenAPINames()[%q] = %q, want %q",
					tt.fqmn, result[tt.fqmn], tt.expectedName)
			}
		})
	}
}

func TestSchemaNullableField(t *testing.T) {
	// Test that nullable field is correctly serialized
	schema := &Schema{
		Type:     "string",
		Nullable: true,
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if !strings.Contains(string(data), `"nullable":true`) {
		t.Errorf("Expected nullable:true in output, got: %s", string(data))
	}

	// Test schema without nullable
	schemaNonNullable := &Schema{
		Type: "string",
	}

	data, err = json.Marshal(schemaNonNullable)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if strings.Contains(string(data), `"nullable"`) {
		t.Errorf("Expected no nullable field in output, got: %s", string(data))
	}
}

func TestGenerateFromProtoDescriptor(t *testing.T) {
	tests := []struct {
		name             string
		inputProtoText   string
		wantJSON         string
		registryModifier func(*descriptor.Registry)
	}{
		{
			name:           "simple echo service",
			inputProtoText: "testdata/generator/simple_echo.prototext",
			wantJSON:       "testdata/generator/simple_echo.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetPreserveRPCOrder(false)
			},
		},
		{
			name:           "simple echo service ordered",
			inputProtoText: "testdata/generator/simple_echo.prototext",
			wantJSON:       "testdata/generator/simple_echo_ordered.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetPreserveRPCOrder(true)
			},
		},
		{
			name:           "merged",
			inputProtoText: "testdata/generator/merged.prototext",
			wantJSON:       "testdata/generator/merged.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetAllowMerge(true)
			},
		},
		{
			name:           "merged ordered",
			inputProtoText: "testdata/generator/merged.prototext",
			wantJSON:       "testdata/generator/merged_ordered.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetAllowMerge(true)
				reg.SetPreserveRPCOrder(true)
			},
		},
		{
			name:           "disable default errors",
			inputProtoText: "testdata/generator/simple_echo.prototext",
			wantJSON:       "testdata/generator/disable_default_errors.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetDisableDefaultErrors(true)
			},
		},
		{
			name:           "disable default responses",
			inputProtoText: "testdata/generator/disable_default_responses.prototext",
			wantJSON:       "testdata/generator/disable_default_responses.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetDisableDefaultResponses(true)
				reg.SetDisableDefaultErrors(true)
			},
		},
		{
			name:           "visibility restriction selectors internal",
			inputProtoText: "testdata/generator/visibility.prototext",
			wantJSON:       "testdata/generator/visibility_internal.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{"INTERNAL"})
				reg.SetAllowMerge(true)
			},
		},
		{
			name:           "visibility restriction selectors none",
			inputProtoText: "testdata/generator/visibility.prototext",
			wantJSON:       "testdata/generator/visibility_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
				reg.SetAllowMerge(true)
			},
		},
		{
			name:           "query param field visibility - internal fields should be excluded",
			inputProtoText: "testdata/generator/query_param_visibility.prototext",
			wantJSON:       "testdata/generator/query_param_visibility_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
			},
		},
		{
			name:           "oneof field visibility - internal fields should be excluded from oneof",
			inputProtoText: "testdata/generator/oneof_visibility.prototext",
			wantJSON:       "testdata/generator/oneof_visibility_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
			},
		},
		{
			name:           "oneof all fields internal - all fields filtered by visibility",
			inputProtoText: "testdata/generator/oneof_all_internal.prototext",
			wantJSON:       "testdata/generator/oneof_all_internal_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
			},
		},
		{
			name:           "oneof multiple groups - one group becomes empty after visibility filter",
			inputProtoText: "testdata/generator/oneof_multiple_groups.prototext",
			wantJSON:       "testdata/generator/oneof_multiple_groups_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load prototext file
			b, err := os.ReadFile(tt.inputProtoText)
			if err != nil {
				t.Fatal(err)
			}

			// Unmarshal into CodeGeneratorRequest
			var req pluginpb.CodeGeneratorRequest
			if err := prototext.Unmarshal(b, &req); err != nil {
				t.Fatal(err)
			}

			// Generate OpenAPI spec
			resp := requireGenerate(t, &req, "3.1.0", tt.registryModifier)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}
			got := resp[0].GetContent()

			// Load expected JSON
			wantBytes, err := os.ReadFile(tt.wantJSON)
			if err != nil {
				t.Fatalf("Failed to read expected JSON file: %v", err)
			}
			want := string(wantBytes)

			// Compare generated JSON with expected JSON
			if !jsonEqualOrdered(t, []byte(got), []byte(want)) {
				t.Errorf("Generated JSON does not match expected JSON.\nGot:\n%s\n\nWant:\n%s", got, want)
			}
		})
	}
}

// requireGenerate is a helper function that generates OpenAPI specs from a CodeGeneratorRequest.
func requireGenerate(
	tb testing.TB,
	req *pluginpb.CodeGeneratorRequest,
	openapiVersion string,
	registryModifier func(*descriptor.Registry),
) []*descriptor.ResponseFile {
	tb.Helper()

	reg := descriptor.NewRegistry()
	if registryModifier != nil {
		registryModifier(reg)
	}

	if err := reg.Load(req); err != nil {
		tb.Fatalf("failed to load request: %s", err)
	}

	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			tb.Fatalf("failed to lookup file: %s", err)
		}

		targets = append(targets, f)
	}

	g := New(reg, FormatJSON, openapiVersion)

	resp, err := g.Generate(targets)
	if err != nil {
		tb.Fatalf("failed to generate targets: %s", err)
	}

	return resp
}

// ============================================================================
// Inline Prototext Tests
// ============================================================================

// requireGenerateInline is a helper that generates OpenAPI from inline prototext.
func requireGenerateInline(
	tb testing.TB,
	protoText string,
	registryModifier func(*descriptor.Registry),
) []*descriptor.ResponseFile {
	tb.Helper()

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(protoText), &req); err != nil {
		tb.Fatalf("failed to unmarshal prototext: %s", err)
	}

	return requireGenerate(tb, &req, "3.0.3", registryModifier)
}

// Basic prototext for a simple service with GET/POST endpoints
const basicServiceProtoInline = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "GetUserRequest"
		field: {
			name: "user_id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "userId"
		}
	}
	message_type: {
		name: "User"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
		field: {
			name: "name"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "name"
		}
	}
	message_type: {
		name: "CreateUserRequest"
		field: {
			name: "name"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "name"
		}
	}
	service: {
		name: "UserService"
		method: {
			name: "GetUser"
			input_type: ".test.v1.GetUserRequest"
			output_type: ".test.v1.User"
			options: {
				[google.api.http]: {
					get: "/v1/users/{user_id}"
				}
			}
		}
		method: {
			name: "CreateUser"
			input_type: ".test.v1.CreateUserRequest"
			output_type: ".test.v1.User"
			options: {
				[google.api.http]: {
					post: "/v1/users"
					body: "*"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`

func TestGenerateInline_BasicService(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, basicServiceProtoInline, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Verify basic structure
	assertions := []string{
		`"openapi": "3.0.3"`,
		`"/v1/users/{user_id}"`,
		`"/v1/users"`,
		`"get":`,
		`"post":`,
		`"UserService_GetUser"`,
		`"UserService_CreateUser"`,
	}

	for _, assertion := range assertions {
		if !strings.Contains(content, assertion) {
			t.Errorf("expected output to contain %q", assertion)
		}
	}
}

func TestGenerateInline_NamingStrategies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		strategy       string
		wantSchemaName string
	}{
		{"fqn", "fqn", `"test.v1.User"`},
		{"simple", "simple", `"User"`},
		{"legacy", "legacy", `"v1User"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := requireGenerateInline(t, basicServiceProtoInline, func(reg *descriptor.Registry) {
				reg.SetOpenAPINamingStrategy(tt.strategy)
			})

			if len(resp) != 1 {
				t.Fatalf("expected 1 response file, got %d", len(resp))
			}

			content := resp[0].GetContent()

			if !strings.Contains(content, tt.wantSchemaName) {
				t.Errorf("expected schema name %s with strategy %s", tt.wantSchemaName, tt.strategy)
			}
		})
	}
}

func TestGenerateInline_SimpleOperationIDs(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, basicServiceProtoInline, func(reg *descriptor.Registry) {
		reg.SetSimpleOperationIDs(true)
	})

	content := resp[0].GetContent()

	if !strings.Contains(content, `"operationId": "GetUser"`) {
		t.Error("expected simple operation ID 'GetUser'")
	}

	if strings.Contains(content, `"operationId": "UserService_GetUser"`) {
		t.Error("should not have service prefix in operation ID")
	}
}

// Proto with enum for testing enum options
const enumServiceProtoInline = `
file_to_generate: "test/v1/enum.proto"
proto_file: {
	name: "test/v1/enum.proto"
	package: "test.v1"
	message_type: {
		name: "Task"
		field: {
			name: "status"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_ENUM
			type_name: ".test.v1.TaskStatus"
			json_name: "status"
		}
	}
	message_type: {
		name: "GetTaskRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	enum_type: {
		name: "TaskStatus"
		value: { name: "TASK_STATUS_UNSPECIFIED" number: 0 }
		value: { name: "TASK_STATUS_PENDING" number: 1 }
		value: { name: "TASK_STATUS_COMPLETED" number: 2 }
	}
	service: {
		name: "TaskService"
		method: {
			name: "GetTask"
			input_type: ".test.v1.GetTaskRequest"
			output_type: ".test.v1.Task"
			options: {
				[google.api.http]: {
					get: "/v1/tasks/{id}"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`

func TestGenerateInline_EnumsAsStrings(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, enumServiceProtoInline, nil)
	content := resp[0].GetContent()

	assertions := []string{
		`"TASK_STATUS_UNSPECIFIED"`,
		`"TASK_STATUS_PENDING"`,
		`"TASK_STATUS_COMPLETED"`,
		`"type": "string"`,
	}

	for _, assertion := range assertions {
		if !strings.Contains(content, assertion) {
			t.Errorf("expected output to contain %q", assertion)
		}
	}
}

func TestGenerateInline_EnumsAsInts(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, enumServiceProtoInline, func(reg *descriptor.Registry) {
		reg.SetEnumsAsInts(true)
	})
	content := resp[0].GetContent()

	if !strings.Contains(content, `"type": "integer"`) {
		t.Error("expected enum type to be integer when enums_as_ints enabled")
	}
}

func TestGenerateInline_OmitEnumDefaultValue(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, enumServiceProtoInline, func(reg *descriptor.Registry) {
		reg.SetOmitEnumDefaultValue(true)
	})
	content := resp[0].GetContent()

	if strings.Contains(content, `"TASK_STATUS_UNSPECIFIED"`) {
		t.Error("should not have default enum value when omit_enum_default_value enabled")
	}

	if !strings.Contains(content, `"TASK_STATUS_PENDING"`) {
		t.Error("should have non-default enum values")
	}
}

// Proto with path parameters
const pathParamsProtoInline = `
file_to_generate: "test/v1/path.proto"
proto_file: {
	name: "test/v1/path.proto"
	package: "test.v1"
	message_type: {
		name: "GetResourceRequest"
		field: {
			name: "project_id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "projectId"
		}
		field: {
			name: "resource_id"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "resourceId"
		}
	}
	message_type: {
		name: "Resource"
		field: {
			name: "name"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "name"
		}
	}
	service: {
		name: "ResourceService"
		method: {
			name: "GetResource"
			input_type: ".test.v1.GetResourceRequest"
			output_type: ".test.v1.Resource"
			options: {
				[google.api.http]: {
					get: "/v1/projects/{project_id}/resources/{resource_id}"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`

func TestGenerateInline_PathParameters(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, pathParamsProtoInline, nil)
	content := resp[0].GetContent()

	assertions := []string{
		`/v1/projects/{project_id}/resources/{resource_id}`,
		`"in": "path"`,
		`"required": true`,
		`"name": "project_id"`,
		`"name": "resource_id"`,
	}

	for _, assertion := range assertions {
		if !strings.Contains(content, assertion) {
			t.Errorf("expected output to contain %q", assertion)
		}
	}
}

// Proto with query parameters
const queryParamsProtoInline = `
file_to_generate: "test/v1/query.proto"
proto_file: {
	name: "test/v1/query.proto"
	package: "test.v1"
	message_type: {
		name: "ListUsersRequest"
		field: {
			name: "page_size"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "pageSize"
		}
		field: {
			name: "page_token"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "pageToken"
		}
	}
	message_type: {
		name: "ListUsersResponse"
		field: {
			name: "users"
			number: 1
			label: LABEL_REPEATED
			type: TYPE_STRING
			json_name: "users"
		}
	}
	service: {
		name: "UserService"
		method: {
			name: "ListUsers"
			input_type: ".test.v1.ListUsersRequest"
			output_type: ".test.v1.ListUsersResponse"
			options: {
				[google.api.http]: {
					get: "/v1/users"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`

func TestGenerateInline_QueryParameters(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, queryParamsProtoInline, nil)
	content := resp[0].GetContent()

	assertions := []string{
		`"in": "query"`,
		`"name": "page_size"`,
		`"name": "page_token"`,
	}

	for _, assertion := range assertions {
		if !strings.Contains(content, assertion) {
			t.Errorf("expected output to contain %q", assertion)
		}
	}
}

func TestGenerateInline_QueryParameters_JSONNames(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, queryParamsProtoInline, func(reg *descriptor.Registry) {
		reg.SetUseJSONNamesForFields(true)
	})
	content := resp[0].GetContent()

	assertions := []string{
		`"in": "query"`,
		`"name": "pageSize"`,
		`"name": "pageToken"`,
	}

	for _, assertion := range assertions {
		if !strings.Contains(content, assertion) {
			t.Errorf("expected output to contain %q with JSON names", assertion)
		}
	}
}

// Well-known type schema tests
func TestWellKnownTypeSchema_Timestamp(t *testing.T) {
	t.Parallel()

	schema := wellKnownTypeSchema(".google.protobuf.Timestamp")
	if schema == nil {
		t.Fatal("expected schema for Timestamp")
	}

	if schema.Type != "string" {
		t.Errorf("Timestamp Type = %q, want %q", schema.Type, "string")
	}
	if schema.Format != "date-time" {
		t.Errorf("Timestamp Format = %q, want %q", schema.Format, "date-time")
	}
}

func TestWellKnownTypeSchema_Duration(t *testing.T) {
	t.Parallel()

	schema := wellKnownTypeSchema(".google.protobuf.Duration")
	if schema == nil {
		t.Fatal("expected schema for Duration")
	}

	if schema.Type != "string" {
		t.Errorf("Duration Type = %q, want %q", schema.Type, "string")
	}
}

func TestWellKnownTypeSchema_Wrappers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typeName   string
		expectType string
		expectFmt  string
	}{
		{".google.protobuf.StringValue", "string", ""},
		{".google.protobuf.Int32Value", "integer", "int32"},
		{".google.protobuf.Int64Value", "string", "int64"},
		{".google.protobuf.BoolValue", "boolean", ""},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			schema := wellKnownTypeSchema(tt.typeName)
			if schema == nil {
				t.Fatalf("expected schema for %s", tt.typeName)
			}
			if schema.Type != tt.expectType {
				t.Errorf("Type = %q, want %q", schema.Type, tt.expectType)
			}
			if tt.expectFmt != "" && schema.Format != tt.expectFmt {
				t.Errorf("Format = %q, want %q", schema.Format, tt.expectFmt)
			}
		})
	}
}

func TestWellKnownTypeSchema_StructTypes(t *testing.T) {
	t.Parallel()

	t.Run("Struct", func(t *testing.T) {
		schema := wellKnownTypeSchema(".google.protobuf.Struct")
		if schema == nil {
			t.Fatal("expected schema for Struct")
		}
		if schema.Type != "object" {
			t.Errorf("Type = %q, want %q", schema.Type, "object")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		schema := wellKnownTypeSchema(".google.protobuf.Empty")
		if schema == nil {
			t.Fatal("expected schema for Empty")
		}
		if schema.Type != "object" {
			t.Errorf("Type = %q, want %q", schema.Type, "object")
		}
	})

	t.Run("Any", func(t *testing.T) {
		schema := wellKnownTypeSchema(".google.protobuf.Any")
		if schema == nil {
			t.Fatal("expected schema for Any")
		}
		if schema.Type != "object" {
			t.Errorf("Type = %q, want %q", schema.Type, "object")
		}
		if schema.Properties == nil || schema.Properties["@type"] == nil {
			t.Error("expected @type property")
		}
	})
}
