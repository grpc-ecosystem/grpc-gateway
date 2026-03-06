package genopenapiv3

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/encoding/prototext"
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

func TestGenerateOneOfSchemas(t *testing.T) {
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

	oneOfSchemas := g.generateOneOfSchemas(parentSchema, groups)

	// Should have 2 oneOf options (one per field in the oneof)
	if len(oneOfSchemas) != 2 {
		t.Fatalf("Expected 2 oneOf schemas, got %d", len(oneOfSchemas))
	}

	// First option should be for stringValue
	opt1 := oneOfSchemas[0].Value
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

	// Second option should be for intValue
	opt2 := oneOfSchemas[1].Value
	if _, ok := opt2.Properties["intValue"]; !ok {
		t.Error("Option 2 should have 'intValue' property")
	}
	if len(opt2.Required) != 1 || opt2.Required[0] != "intValue" {
		t.Errorf("Option 2 Required = %v, want [intValue]", opt2.Required)
	}
}

func TestGenerateOneOfSchemasMultipleGroups(t *testing.T) {
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

	oneOfSchemas := g.generateOneOfSchemas(parentSchema, groups)

	// Should have 4 oneOf options (2 per group)
	if len(oneOfSchemas) != 4 {
		t.Fatalf("Expected 4 oneOf schemas, got %d", len(oneOfSchemas))
	}

	// Check that titles identify the oneof group
	titles := make(map[string]bool)
	for _, schema := range oneOfSchemas {
		titles[schema.Value.Title] = true
	}

	expectedTitles := []string{"event.createEvent", "event.updateEvent", "result.error", "result.success"}
	for _, title := range expectedTitles {
		if !titles[title] {
			t.Errorf("Missing expected oneOf option with title %q", title)
		}
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

			schema := &Schema{Type: "string"}
			field := &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:           stringPtr("test_field"),
					Proto3Optional: &tt.proto3Optional,
					OneofIndex:     tt.oneofIndex,
					Options:        nil,
				},
			}

			gotRequired := g.applyFieldBehaviorToSchema(schema, field)
			if gotRequired != tt.wantRequired {
				t.Errorf("applyFieldBehaviorToSchema() returned required=%v, want %v", gotRequired, tt.wantRequired)
			}
		})
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
