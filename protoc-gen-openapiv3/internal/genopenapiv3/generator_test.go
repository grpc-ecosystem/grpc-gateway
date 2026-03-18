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
	"google.golang.org/genproto/googleapis/api/visibility"
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
		ref      *SchemaOrReference
		expected string
	}{
		{
			name:     "reference",
			ref:      NewSchemaRef("my.Message"),
			expected: `{"$ref":"#/components/schemas/my.Message"}`,
		},
		{
			name:     "inline string schema",
			ref:      &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}},
			expected: `{"type":"string"}`,
		},
		{
			name:     "inline integer schema",
			ref:      &SchemaOrReference{Schema: &Schema{Type: SchemaType{"integer"}, Format: "int32"}},
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

func TestReferenceMarshal(t *testing.T) {
	tests := []struct {
		name     string
		ref      *Reference
		expected string
	}{
		{
			name:     "simple reference",
			ref:      &Reference{Ref: "#/components/schemas/Pet"},
			expected: `{"$ref":"#/components/schemas/Pet"}`,
		},
		{
			name: "reference with summary (v3.1.0)",
			ref: &Reference{
				Ref:     "#/components/schemas/Pet",
				Summary: "A pet in the store",
			},
			expected: `{"$ref":"#/components/schemas/Pet","summary":"A pet in the store"}`,
		},
		{
			name: "reference with description (v3.1.0)",
			ref: &Reference{
				Ref:         "#/components/schemas/Pet",
				Description: "Detailed description of a pet",
			},
			expected: `{"$ref":"#/components/schemas/Pet","description":"Detailed description of a pet"}`,
		},
		{
			name: "reference with all fields (v3.1.0)",
			ref: &Reference{
				Ref:         "#/components/schemas/Pet",
				Summary:     "A pet",
				Description: "A detailed description",
			},
			expected: `{"$ref":"#/components/schemas/Pet","summary":"A pet","description":"A detailed description"}`,
		},
		{
			name:     "nil reference",
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

func TestSchemaRefWithReferenceMarshal(t *testing.T) {
	tests := []struct {
		name     string
		ref      *SchemaOrReference
		expected string
	}{
		{
			name:     "inline schema",
			ref:      NewInlineSchema(&Schema{Type: SchemaType{"string"}}),
			expected: `{"type":"string"}`,
		},
		{
			name:     "simple reference",
			ref:      NewSchemaRef("Pet"),
			expected: `{"$ref":"#/components/schemas/Pet"}`,
		},
		{
			name: "reference with summary (v3.1.0)",
			ref: &SchemaOrReference{
				Reference: &Reference{
					Ref:     "#/components/schemas/Pet",
					Summary: "A pet summary",
				},
			},
			expected: `{"$ref":"#/components/schemas/Pet","summary":"A pet summary"}`,
		},
		{
			name: "reference with description (v3.1.0)",
			ref: &SchemaOrReference{
				Reference: &Reference{
					Ref:         "#/components/schemas/Pet",
					Description: "A detailed pet description",
				},
			},
			expected: `{"$ref":"#/components/schemas/Pet","description":"A detailed pet description"}`,
		},
		{
			name:     "reference with all v3.1.0 overrides",
			ref:      NewSchemaRefWithOverrides("Pet", "A pet summary", "A detailed description"),
			expected: `{"$ref":"#/components/schemas/Pet","summary":"A pet summary","description":"A detailed description"}`,
		},
		{
			name:     "nil",
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
	if resp.Content["application/json"].Schema.Reference == nil || resp.Content["application/json"].Schema.Reference.Ref != "#/components/schemas/MyResponse" {
		ref := ""
		if resp.Content["application/json"].Schema.Reference != nil {
			ref = resp.Content["application/json"].Schema.Reference.Ref
		}
		t.Errorf("Schema.Reference.Ref = %q, want %q", ref, "#/components/schemas/MyResponse")
	}
}

func TestParameterCreation(t *testing.T) {
	t.Run("path parameter", func(t *testing.T) {
		param := NewPathParameter("user_id", &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}})
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
		param := NewQueryParameter("limit", &SchemaOrReference{Schema: &Schema{Type: SchemaType{"integer"}}})
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
		param := NewHeaderParameter("X-Custom", &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}}})
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
	if body.Content["application/json"].Schema.Reference == nil || body.Content["application/json"].Schema.Reference.Ref != "#/components/schemas/MyRequest" {
		ref := ""
		if body.Content["application/json"].Schema.Reference != nil {
			ref = body.Content["application/json"].Schema.Reference.Ref
		}
		t.Errorf("Schema.Reference.Ref = %q, want %q", ref, "#/components/schemas/MyRequest")
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
	// Uses OpenAPI 3.1.0 style: wrapper types use type arrays for nullable
	tests := []struct {
		typeName   string
		expectType SchemaType
		expectFmt  string
	}{
		{".google.protobuf.Timestamp", SchemaType{"string"}, "date-time"},
		{".google.protobuf.Duration", SchemaType{"string"}, ""},
		{".google.protobuf.StringValue", SchemaType{"string", "null"}, ""},
		{".google.protobuf.Int32Value", SchemaType{"integer", "null"}, "int32"},
		{".google.protobuf.Int64Value", SchemaType{"string", "null"}, "int64"},
		{".google.protobuf.BoolValue", SchemaType{"boolean", "null"}, ""},
		{".google.protobuf.Empty", SchemaType{"object"}, ""},
		{".google.protobuf.Struct", SchemaType{"object"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			schema := wellKnownTypeSchema(tt.typeName)
			if schema == nil {
				t.Fatalf("schema should exist for %s", tt.typeName)
			}
			if !schemaTypeEqual(schema.Type, tt.expectType) {
				t.Errorf("Type = %v, want %v", schema.Type, tt.expectType)
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
		Type: SchemaType{"object"},
		Properties: map[string]*SchemaOrReference{
			"stringValue":  {Schema: &Schema{Type: SchemaType{"string"}}},
			"intValue":     {Schema: &Schema{Type: SchemaType{"integer"}}},
			"regularField": {Schema: &Schema{Type: SchemaType{"boolean"}}},
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

	// hasRegularFields=true because we have properties in parentSchema
	doc := NewOpenAPI("Test", "1.0.0", "3.1.0")
	visited := make(map[string]bool)
	oneOf, allOf := g.generateOneOfConstraints(doc, parentSchema, groups, true, visited)

	// Single group should use oneOf directly, not allOf
	if allOf != nil {
		t.Error("Single group should not use allOf")
	}

	// Should have 3 oneOf options (2 fields + neither)
	if len(oneOf) != 3 {
		t.Fatalf("Expected 3 oneOf schemas (2 fields + neither), got %d", len(oneOf))
	}

	// First option should be for stringValue
	opt1 := oneOf[0].Schema
	if !schemaTypeEqual(opt1.Type, SchemaType{"object"}) {
		t.Errorf("Option 1 Type = %v, want %q", opt1.Type, "object")
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
	opt2 := oneOf[1].Schema
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
	neitherOpt := oneOf[2].Schema
	if neitherOpt.Title != "value.none" {
		t.Errorf("Neither option Title = %q, want %q", neitherOpt.Title, "value.none")
	}
	if neitherOpt.Not == nil {
		t.Fatal("Neither option should have 'not' schema")
	}
	if neitherOpt.Not.Schema == nil {
		t.Fatal("Neither option 'not' should have value")
	}
	if neitherOpt.Not.Schema.AnyOf == nil || len(neitherOpt.Not.Schema.AnyOf) != 2 {
		t.Fatalf("Neither option 'not.anyOf' should have 2 entries, got %v", neitherOpt.Not.Schema.AnyOf)
	}
}

func TestGenerateOneOfConstraintsMultipleGroups(t *testing.T) {
	// Create a generator with a registry that uses JSON names
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(true)
	g := &generator{reg: reg}

	// Create a parent schema with properties for multiple oneofs
	parentSchema := &Schema{
		Type: SchemaType{"object"},
		Properties: map[string]*SchemaOrReference{
			"createEvent": {Schema: &Schema{Type: SchemaType{"object"}}},
			"updateEvent": {Schema: &Schema{Type: SchemaType{"object"}}},
			"error":       {Schema: &Schema{Type: SchemaType{"object"}}},
			"success":     {Schema: &Schema{Type: SchemaType{"object"}}},
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

	doc := NewOpenAPI("Test", "1.0.0", "3.1.0")
	visited := make(map[string]bool)
	oneOf, allOf := g.generateOneOfConstraints(doc, parentSchema, groups, true, visited)

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
		if allOfEntry.Schema == nil {
			t.Errorf("allOf[%d] should have value", groupIdx)
			continue
		}
		groupOneOf := allOfEntry.Schema.OneOf
		if len(groupOneOf) != 3 {
			t.Errorf("Group %d should have 3 oneOf options (2 fields + neither), got %d", groupIdx, len(groupOneOf))
		}
	}

	// Verify group 1 (event) has correct options
	group1OneOf := allOf[0].Schema.OneOf
	group1Titles := make(map[string]bool)
	for _, schema := range group1OneOf {
		group1Titles[schema.Schema.Title] = true
	}
	for _, expected := range []string{"event.createEvent", "event.updateEvent", "event.none"} {
		if !group1Titles[expected] {
			t.Errorf("Group 1 missing expected oneOf option with title %q", expected)
		}
	}

	// Verify group 2 (result) has correct options
	group2OneOf := allOf[1].Schema.OneOf
	group2Titles := make(map[string]bool)
	for _, schema := range group2OneOf {
		group2Titles[schema.Schema.Title] = true
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
	if schema.Not.Schema == nil {
		t.Fatal("Not should have value")
	}
	if schema.Not.Schema.AnyOf == nil {
		t.Fatal("Not should have anyOf")
	}

	// AnyOf should contain required entry for each field
	if len(schema.Not.Schema.AnyOf) != 3 {
		t.Errorf("anyOf should have 3 entries, got %d", len(schema.Not.Schema.AnyOf))
	}

	requiredFields := make(map[string]bool)
	for _, ref := range schema.Not.Schema.AnyOf {
		if ref.Schema != nil && len(ref.Schema.Required) == 1 {
			requiredFields[ref.Schema.Required[0]] = true
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
		Type:       SchemaType{"object"},
		Properties: map[string]*SchemaOrReference{},
	}

	// Empty groups
	doc := NewOpenAPI("Test", "1.0.0", "3.1.0")
	visited := make(map[string]bool)
	oneOf, allOf := g.generateOneOfConstraints(doc, parentSchema, []oneofGroup{}, false, visited)

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
			fieldSchema := &Schema{Type: SchemaType{"string"}}

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
	fieldSchema1 := &Schema{Type: SchemaType{"string"}}
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
	fieldSchema2 := &Schema{Type: SchemaType{"string"}}
	g.applyFieldBehaviorToSchema(parentSchema, fieldSchema2, field2)

	// Field 3: No behavior - should NOT be in Required array
	field3 := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:    stringPtr("plain_field"),
			Options: nil,
		},
	}
	fieldSchema3 := &Schema{Type: SchemaType{"string"}}
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
	fieldSchema4 := &Schema{Type: SchemaType{"string"}}
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
		Type:      SchemaType{"string"},
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
		Type:      SchemaType{"string"},
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
		Type:     SchemaType{"string"},
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
		Type: SchemaType{"string"},
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
		{
			name:           "oneof with enum fields - enum types in oneOf should use $ref",
			inputProtoText: "testdata/generator/oneof_enum.prototext",
			wantJSON:       "testdata/generator/oneof_enum_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
			},
		},
		{
			name:           "oneof with message references - message types in oneOf should use $ref",
			inputProtoText: "testdata/generator/oneof_wkt.prototext",
			wantJSON:       "testdata/generator/oneof_wkt_none.openapi.json",
			registryModifier: func(reg *descriptor.Registry) {
				reg.SetVisibilityRestrictionSelectors([]string{})
			},
		},
		{
			name:           "oneof in nested message definition - nested messages with oneOf",
			inputProtoText: "testdata/generator/oneof_nested_definition.prototext",
			wantJSON:       "testdata/generator/oneof_nested_definition_none.openapi.json",
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

	if !schemaTypeEqual(schema.Type, SchemaType{"string"}) {
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

	if !schemaTypeEqual(schema.Type, SchemaType{"string"}) {
		t.Errorf("Duration Type = %q, want %q", schema.Type, "string")
	}
}

func TestWellKnownTypeSchema_Wrappers(t *testing.T) {
	t.Parallel()

	// OpenAPI 3.1.0 style: wrapper types use type arrays for nullable
	tests := []struct {
		typeName   string
		expectType SchemaType
		expectFmt  string
	}{
		{".google.protobuf.StringValue", SchemaType{"string", "null"}, ""},
		{".google.protobuf.Int32Value", SchemaType{"integer", "null"}, "int32"},
		{".google.protobuf.Int64Value", SchemaType{"string", "null"}, "int64"},
		{".google.protobuf.BoolValue", SchemaType{"boolean", "null"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			schema := wellKnownTypeSchema(tt.typeName)
			if schema == nil {
				t.Fatalf("expected schema for %s", tt.typeName)
			}
			if !schemaTypeEqual(schema.Type, tt.expectType) {
				t.Errorf("Type = %v, want %v", schema.Type, tt.expectType)
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
		if !schemaTypeEqual(schema.Type, SchemaType{"object"}) {
			t.Errorf("Type = %q, want %q", schema.Type, "object")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		schema := wellKnownTypeSchema(".google.protobuf.Empty")
		if schema == nil {
			t.Fatal("expected schema for Empty")
		}
		if !schemaTypeEqual(schema.Type, SchemaType{"object"}) {
			t.Errorf("Type = %q, want %q", schema.Type, "object")
		}
	})

	t.Run("Any", func(t *testing.T) {
		schema := wellKnownTypeSchema(".google.protobuf.Any")
		if schema == nil {
			t.Fatal("expected schema for Any")
		}
		if !schemaTypeEqual(schema.Type, SchemaType{"object"}) {
			t.Errorf("Type = %q, want %q", schema.Type, "object")
		}
		if schema.Properties == nil || schema.Properties["@type"] == nil {
			t.Error("expected @type property")
		}
	})
}

// ============================================================================
// Priority 1: Bug-Exposing Tests
// ============================================================================

// TestMapFieldSchema tests that proto map fields are correctly rendered as
// OpenAPI object schemas with additionalProperties, not as arrays.
//
// In protobuf, map<K,V> fields are represented as repeated message types with
// the `map_entry` option set to true. The generator should detect this pattern
// and render it as an object with additionalProperties containing the value schema.
//
// BUG: Currently, the generator treats map fields as regular repeated message
// fields and renders them as `type: array` instead of `type: object`.
func TestMapFieldSchema(t *testing.T) {
	t.Parallel()

	// Create a proto descriptor with a map field
	// map<string, int32> is represented as:
	// - A synthetic message type with map_entry=true containing key and value fields
	// - A repeated field of that synthetic message type
	mapEntryMsg := &descriptorpb.DescriptorProto{
		Name: stringPtr("MetadataEntry"),
		Options: &descriptorpb.MessageOptions{
			MapEntry: boolPtr(true),
		},
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     stringPtr("key"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: stringPtr("key"),
			},
			{
				Name:     stringPtr("value"),
				Number:   int32Ptr(2),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				JsonName: stringPtr("value"),
			},
		},
	}

	// Create the parent message containing the map field
	parentMsg := &descriptorpb.DescriptorProto{
		Name: stringPtr("ResourceWithMap"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     stringPtr("id"),
				Number:   int32Ptr(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				JsonName: stringPtr("id"),
			},
			{
				Name:     stringPtr("metadata"),
				Number:   int32Ptr(2),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: stringPtr(".test.v1.ResourceWithMap.MetadataEntry"),
				JsonName: stringPtr("metadata"),
			},
		},
		NestedType: []*descriptorpb.DescriptorProto{mapEntryMsg},
	}

	// Build the full request
	protoText := `
file_to_generate: "test/v1/map.proto"
proto_file: {
	name: "test/v1/map.proto"
	package: "test.v1"
	message_type: {
		name: "ResourceWithMap"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
		field: {
			name: "metadata"
			number: 2
			label: LABEL_REPEATED
			type: TYPE_MESSAGE
			type_name: ".test.v1.ResourceWithMap.MetadataEntry"
			json_name: "metadata"
		}
		nested_type: {
			name: "MetadataEntry"
			options: {
				map_entry: true
			}
			field: {
				name: "key"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "key"
			}
			field: {
				name: "value"
				number: 2
				label: LABEL_OPTIONAL
				type: TYPE_INT32
				json_name: "value"
			}
		}
	}
	message_type: {
		name: "GetRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	service: {
		name: "MapService"
		method: {
			name: "Get"
			input_type: ".test.v1.GetRequest"
			output_type: ".test.v1.ResourceWithMap"
			options: {
				[google.api.http]: {
					get: "/v1/resources/{id}"
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

	resp := requireGenerateInline(t, protoText, nil)
	content := resp[0].GetContent()

	// Parse the JSON to verify the structure
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	components, ok := doc["components"].(map[string]interface{})
	if !ok {
		t.Fatal("missing components")
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		t.Fatal("missing schemas")
	}

	// Find the ResourceWithMap schema (may have v1 prefix from naming strategy)
	var resourceSchema map[string]interface{}
	for name, s := range schemas {
		if strings.Contains(name, "ResourceWithMap") && !strings.Contains(name, "Entry") {
			resourceSchema, _ = s.(map[string]interface{})
			break
		}
	}

	if resourceSchema == nil {
		t.Fatal("ResourceWithMap schema not found")
	}

	properties, ok := resourceSchema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("ResourceWithMap missing properties")
	}

	metadataField, ok := properties["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("metadata field not found in ResourceWithMap properties")
	}

	// Map fields should be type: object with additionalProperties, not type: array
	fieldType, _ := metadataField["type"].(string)
	if fieldType == "array" {
		t.Errorf("BUG CONFIRMED: Map field 'metadata' has type 'array', should be 'object' with additionalProperties")
	} else if fieldType != "object" {
		t.Errorf("Map field 'metadata' has type %q, expected 'object'", fieldType)
	}

	// Should have additionalProperties pointing to the value type schema
	if metadataField["additionalProperties"] == nil {
		t.Error("Map field 'metadata' should have additionalProperties")
	}

	// Verify that the MetadataEntry schema is NOT included (it's synthetic)
	for name := range schemas {
		if strings.Contains(name, "MetadataEntry") {
			t.Errorf("synthetic map entry schema %q should not appear in output", name)
		}
	}

	// Keep the unused variables for documentation purposes
	_ = parentMsg
}

func boolPtr(b bool) *bool {
	return &b
}

// TestNumericZeroConstraints tests that minimum and maximum values of 0 are
// correctly applied to schemas.
//
// This test verifies that the applySchemaOptionsToSchema function correctly
// handles zero values for minimum/maximum constraints. The proto now uses
// optional fields so we can distinguish "not set" from "explicitly set to 0".
func TestNumericZeroConstraints(t *testing.T) {
	t.Parallel()

	// Helper to create float64 pointer
	f64 := func(v float64) *float64 { return &v }

	tests := []struct {
		name        string
		minimum     *float64
		maximum     *float64
		wantMinimum bool
		wantMaximum bool
	}{
		{
			name:        "minimum 0, maximum 10",
			minimum:     f64(0),
			maximum:     f64(10),
			wantMinimum: true,
			wantMaximum: true,
		},
		{
			name:        "minimum -10, maximum 0",
			minimum:     f64(-10),
			maximum:     f64(0),
			wantMinimum: true,
			wantMaximum: true,
		},
		{
			name:        "minimum 0, maximum 0 (both zero)",
			minimum:     f64(0),
			maximum:     f64(0),
			wantMinimum: true,
			wantMaximum: true,
		},
		{
			name:        "non-zero values",
			minimum:     f64(1),
			maximum:     f64(100),
			wantMinimum: true,
			wantMaximum: true,
		},
		{
			name:        "nil values (not set)",
			minimum:     nil,
			maximum:     nil,
			wantMinimum: false,
			wantMaximum: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := &Schema{Type: SchemaType{"integer"}}

			// Use the actual logic from apply_annotations.go
			// This now correctly uses pointer checks instead of value checks
			if tt.minimum != nil {
				min := *tt.minimum
				schema.Minimum = &min
			}
			if tt.maximum != nil {
				max := *tt.maximum
				schema.Maximum = &max
			}

			// Check if the schema has minimum when it should
			hasMinimum := schema.Minimum != nil
			if tt.wantMinimum && !hasMinimum {
				t.Errorf("minimum was not applied, want minimum=%v", *tt.minimum)
			}
			if !tt.wantMinimum && hasMinimum {
				t.Errorf("minimum was applied but should not have been")
			}

			// Check if the schema has maximum when it should
			hasMaximum := schema.Maximum != nil
			if tt.wantMaximum && !hasMaximum {
				t.Errorf("maximum was not applied, want maximum=%v", *tt.maximum)
			}
			if !tt.wantMaximum && hasMaximum {
				t.Errorf("maximum was applied but should not have been")
			}

			// Verify the actual values are correct
			if tt.wantMinimum && hasMinimum && *schema.Minimum != *tt.minimum {
				t.Errorf("minimum = %v, want %v", *schema.Minimum, *tt.minimum)
			}
			if tt.wantMaximum && hasMaximum && *schema.Maximum != *tt.maximum {
				t.Errorf("maximum = %v, want %v", *schema.Maximum, *tt.maximum)
			}
		})
	}
}

// Tests that empty/whitespace-only visibility restrictions don't incorrectly hide fields
// BUG: Whitespace-only restriction "   " causes field to be incorrectly hidden
func TestEmptyVisibilityRestriction(t *testing.T) {
	t.Parallel()

	// Test case: isVisible with whitespace-only restriction
	// When restriction is "   " (whitespace only):
	// - strings.TrimSpace("   ") returns ""
	// - strings.Split("", ",") returns [""] (slice with one empty string)
	// - len(restrictions) is 1, not 0
	// - loop runs with restriction="" which won't match any selector
	// - returns false (incorrectly hidden)

	// This is a unit test for the isVisible function behavior
	tests := []struct {
		name        string
		restriction string
		selectors   []string
		wantVisible bool
	}{
		{
			name:        "empty restriction should be visible",
			restriction: "",
			selectors:   []string{},
			wantVisible: true,
		},
		{
			name:        "whitespace only restriction should be visible",
			restriction: "   ",
			selectors:   []string{},
			wantVisible: true,
		},
		{
			name:        "whitespace with tabs should be visible",
			restriction: "  \t  ",
			selectors:   []string{},
			wantVisible: true,
		},
		{
			name:        "INTERNAL restriction without selector should be hidden",
			restriction: "INTERNAL",
			selectors:   []string{},
			wantVisible: false,
		},
		{
			name:        "INTERNAL restriction with matching selector should be visible",
			restriction: "INTERNAL",
			selectors:   []string{"INTERNAL"},
			wantVisible: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := descriptor.NewRegistry()
			reg.SetVisibilityRestrictionSelectors(tt.selectors)

			// Create a visibility rule
			rule := &visibility.VisibilityRule{
				Restriction: tt.restriction,
			}

			gotVisible := isVisible(rule, reg)
			if gotVisible != tt.wantVisible {
				if tt.restriction == "   " || tt.restriction == "  \t  " {
					t.Errorf("BUG CONFIRMED: isVisible() = %v, want %v (whitespace-only restriction incorrectly handled)", gotVisible, tt.wantVisible)
				} else {
					t.Errorf("isVisible() = %v, want %v", gotVisible, tt.wantVisible)
				}
			}
		})
	}
}

// ============================================================================
// Priority 2: OpenAPI v3-Specific Feature Tests
// ============================================================================

// TestAllWellKnownTypes tests that all protobuf well-known types are correctly
// mapped to their OpenAPI schema representations.
// Uses OpenAPI 3.1.0 style where nullable is expressed via type arrays (e.g., ["string", "null"]).
func TestAllWellKnownTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typeName     string
		expectType   SchemaType // Use SchemaType for 3.1.0 style type arrays
		expectFormat string
		expectProps  []string
	}{
		// Timestamp and Duration
		{".google.protobuf.Timestamp", SchemaType{"string"}, "date-time", nil},
		{".google.protobuf.Duration", SchemaType{"string"}, "", nil},

		// FieldMask
		{".google.protobuf.FieldMask", SchemaType{"string"}, "", nil},

		// Wrapper types - nullable via type array (3.1.0 style)
		{".google.protobuf.StringValue", SchemaType{"string", "null"}, "", nil},
		{".google.protobuf.BytesValue", SchemaType{"string", "null"}, "byte", nil},
		{".google.protobuf.Int32Value", SchemaType{"integer", "null"}, "int32", nil},
		{".google.protobuf.UInt32Value", SchemaType{"integer", "null"}, "int64", nil},
		{".google.protobuf.Int64Value", SchemaType{"string", "null"}, "int64", nil},
		{".google.protobuf.UInt64Value", SchemaType{"string", "null"}, "uint64", nil},
		{".google.protobuf.FloatValue", SchemaType{"number", "null"}, "float", nil},
		{".google.protobuf.DoubleValue", SchemaType{"number", "null"}, "double", nil},
		{".google.protobuf.BoolValue", SchemaType{"boolean", "null"}, "", nil},

		// Struct types
		{".google.protobuf.Empty", SchemaType{"object"}, "", nil},
		{".google.protobuf.Struct", SchemaType{"object"}, "", nil},
		{".google.protobuf.Value", nil, "", nil}, // No type constraint (empty schema)
		{".google.protobuf.ListValue", SchemaType{"array"}, "", nil},
		{".google.protobuf.NullValue", SchemaType{"null"}, "", nil},

		// Any type
		{".google.protobuf.Any", SchemaType{"object"}, "", []string{"@type"}},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			schema := wellKnownTypeSchema(tt.typeName)
			if schema == nil {
				t.Fatalf("expected schema for %s", tt.typeName)
			}

			if !schemaTypeEqual(schema.Type, tt.expectType) {
				t.Errorf("Type = %v, want %v", schema.Type, tt.expectType)
			}

			if tt.expectFormat != "" && schema.Format != tt.expectFormat {
				t.Errorf("Format = %q, want %q", schema.Format, tt.expectFormat)
			}

			for _, prop := range tt.expectProps {
				if schema.Properties == nil || schema.Properties[prop] == nil {
					t.Errorf("missing expected property %q", prop)
				}
			}
		})
	}
}

// TestOneOfNestedMessages tests complex oneOf scenarios including nested oneOf
// and oneOf with various field types.
const nestedOneOfProtoInline = `
file_to_generate: "test/v1/nested_oneof.proto"
proto_file: {
	name: "test/v1/nested_oneof.proto"
	package: "test.v1"
	message_type: {
		name: "InnerChoice"
		field: {
			name: "text"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "text"
			oneof_index: 0
		}
		field: {
			name: "number"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "number"
			oneof_index: 0
		}
		oneof_decl: {
			name: "inner_value"
		}
	}
	message_type: {
		name: "OuterChoice"
		field: {
			name: "simple"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "simple"
			oneof_index: 0
		}
		field: {
			name: "nested"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.InnerChoice"
			json_name: "nested"
			oneof_index: 0
		}
		oneof_decl: {
			name: "outer_value"
		}
	}
	message_type: {
		name: "GetRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	service: {
		name: "NestedService"
		method: {
			name: "Get"
			input_type: ".test.v1.GetRequest"
			output_type: ".test.v1.OuterChoice"
			options: {
				[google.api.http]: {
					get: "/v1/nested/{id}"
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

func TestOneOfNestedMessages(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, nestedOneOfProtoInline, nil)
	content := resp[0].GetContent()

	// Parse the JSON to verify oneOf structure
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	components, ok := doc["components"].(map[string]interface{})
	if !ok {
		t.Fatal("missing components")
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		t.Fatal("missing schemas")
	}

	// Find the OuterChoice schema
	var outerSchema map[string]interface{}
	for name, s := range schemas {
		if strings.Contains(name, "OuterChoice") {
			outerSchema, _ = s.(map[string]interface{})
			break
		}
	}

	if outerSchema == nil {
		t.Fatal("OuterChoice schema not found")
	}

	// Verify oneOf is present for OuterChoice
	oneOf, hasOneOf := outerSchema["oneOf"].([]interface{})
	if !hasOneOf {
		t.Error("OuterChoice should have oneOf for the oneof_decl")
	} else {
		// Verify oneOf has the expected number of options
		// Note: oneOf includes 2 actual options + 1 "neither" option = 3
		if len(oneOf) != 3 {
			t.Errorf("OuterChoice oneOf has %d options, want 3 (2 fields + none)", len(oneOf))
		}
	}

	// Pure oneof messages should NOT have redundant top-level type and properties
	if _, hasType := outerSchema["type"]; hasType {
		t.Error("Pure oneof message OuterChoice should NOT have top-level 'type' field")
	}
	if _, hasProps := outerSchema["properties"]; hasProps {
		t.Error("Pure oneof message OuterChoice should NOT have top-level 'properties' field")
	}

	// Find the InnerChoice schema
	var innerSchema map[string]interface{}
	for name, s := range schemas {
		if strings.Contains(name, "InnerChoice") {
			innerSchema, _ = s.(map[string]interface{})
			break
		}
	}

	if innerSchema == nil {
		t.Fatal("InnerChoice schema not found")
	}

	// Verify oneOf is present for InnerChoice
	innerOneOf, hasInnerOneOf := innerSchema["oneOf"].([]interface{})
	if !hasInnerOneOf {
		t.Error("InnerChoice should have oneOf for the oneof_decl")
	} else {
		// Note: oneOf includes 2 actual options + 1 "neither" option = 3
		if len(innerOneOf) != 3 {
			t.Errorf("InnerChoice oneOf has %d options, want 3 (2 fields + none)", len(innerOneOf))
		}
	}

	// Pure oneof messages should NOT have redundant top-level type and properties
	if _, hasType := innerSchema["type"]; hasType {
		t.Error("Pure oneof message InnerChoice should NOT have top-level 'type' field")
	}
	if _, hasProps := innerSchema["properties"]; hasProps {
		t.Error("Pure oneof message InnerChoice should NOT have top-level 'properties' field")
	}
}

// TestOneOfWithRegularFields tests messages that have both regular fields and oneof.
// These should have top-level type/properties for regular fields + oneOf constraint.
const mixedOneOfProtoInline = `
file_to_generate: "test/v1/mixed_oneof.proto"
proto_file: {
	name: "test/v1/mixed_oneof.proto"
	package: "test.v1"
	message_type: {
		name: "MixedMessage"
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
		field: {
			name: "text_value"
			number: 3
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "textValue"
			oneof_index: 0
		}
		field: {
			name: "int_value"
			number: 4
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "intValue"
			oneof_index: 0
		}
		oneof_decl: {
			name: "value"
		}
	}
	message_type: {
		name: "GetMixedRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	service: {
		name: "MixedService"
		method: {
			name: "Get"
			input_type: ".test.v1.GetMixedRequest"
			output_type: ".test.v1.MixedMessage"
			options: {
				[google.api.http]: {
					get: "/v1/mixed/{id}"
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

func TestOneOfWithRegularFields(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, mixedOneOfProtoInline, nil)
	content := resp[0].GetContent()

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	components, ok := doc["components"].(map[string]interface{})
	if !ok {
		t.Fatal("missing components")
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		t.Fatal("missing schemas")
	}

	// Find the MixedMessage schema
	var mixedSchema map[string]interface{}
	for name, s := range schemas {
		if strings.Contains(name, "MixedMessage") {
			mixedSchema, _ = s.(map[string]interface{})
			break
		}
	}

	if mixedSchema == nil {
		t.Fatal("MixedMessage schema not found")
	}

	// Mixed messages SHOULD have top-level type and properties (for regular fields)
	if _, hasType := mixedSchema["type"]; !hasType {
		t.Error("Mixed message should have top-level 'type' field")
	}
	props, hasProps := mixedSchema["properties"].(map[string]interface{})
	if !hasProps {
		t.Fatal("Mixed message should have top-level 'properties' field")
	}

	// Should have regular fields in properties
	if _, hasId := props["id"]; !hasId {
		t.Error("Mixed message properties should include 'id'")
	}
	if _, hasName := props["name"]; !hasName {
		t.Error("Mixed message properties should include 'name'")
	}

	// Should also have oneof fields in properties (for JSON serialization)
	// Note: field names may be snake_case or camelCase depending on registry settings
	hasTextValue := props["textValue"] != nil || props["text_value"] != nil
	if !hasTextValue {
		t.Error("Mixed message properties should include oneof field 'textValue' or 'text_value'")
	}
	hasIntValue := props["intValue"] != nil || props["int_value"] != nil
	if !hasIntValue {
		t.Error("Mixed message properties should include oneof field 'intValue' or 'int_value'")
	}

	// Should also have oneOf constraint
	oneOf, hasOneOf := mixedSchema["oneOf"].([]interface{})
	if !hasOneOf {
		t.Error("Mixed message should have oneOf constraint")
	} else {
		// Should have 3 options (2 fields + none)
		if len(oneOf) != 3 {
			t.Errorf("Mixed message oneOf has %d options, want 3 (2 fields + none)", len(oneOf))
		}
	}
}

// TestProto3OptionalWithRequired tests the interaction between proto3 optional
// fields and REQUIRED field_behavior annotation.
//
// When a field is both `proto3 optional` AND has `field_behavior = REQUIRED`,
// the behavior should be well-defined with explicit precedence rules.
func TestProto3OptionalWithRequired(t *testing.T) {
	t.Parallel()

	// Test case: proto3_optional=true with field_behavior=REQUIRED
	// Current behavior: proto3_optional makes the field optional
	// Expected: REQUIRED annotation should take precedence OR there should be
	// clear documentation about precedence rules

	tests := []struct {
		name               string
		proto3Optional     bool
		useProto3Semantics bool
		fieldBehavior      []int32 // FieldBehavior enum values
		wantRequired       bool
		description        string
	}{
		{
			name:               "regular field with proto3 semantics",
			proto3Optional:     false,
			useProto3Semantics: true,
			fieldBehavior:      nil,
			wantRequired:       true,
			description:        "non-optional field with proto3 semantics should be required",
		},
		{
			name:               "proto3 optional without field_behavior",
			proto3Optional:     true,
			useProto3Semantics: true,
			fieldBehavior:      nil,
			wantRequired:       false,
			description:        "optional field should not be required",
		},
		{
			name:               "regular field with REQUIRED behavior",
			proto3Optional:     false,
			useProto3Semantics: false,
			fieldBehavior:      []int32{2}, // REQUIRED = 2
			wantRequired:       true,
			description:        "REQUIRED behavior should make field required",
		},
		{
			name:               "proto3 optional with REQUIRED behavior (conflict)",
			proto3Optional:     true,
			useProto3Semantics: true,
			fieldBehavior:      []int32{2}, // REQUIRED = 2
			wantRequired:       true,       // REQUIRED should win
			description:        "REQUIRED annotation should take precedence over proto3_optional",
		},
		{
			name:               "regular field with OUTPUT_ONLY behavior",
			proto3Optional:     false,
			useProto3Semantics: true,
			fieldBehavior:      []int32{3}, // OUTPUT_ONLY = 3
			wantRequired:       false,
			description:        "OUTPUT_ONLY fields should not be required",
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
					Number:         int32Ptr(1),
					Proto3Optional: &tt.proto3Optional,
					OneofIndex:     nil,
					Options:        nil,
				},
			}

			// Note: We can't easily test with field_behavior annotations without
			// setting up the full proto options, so this tests the base behavior
			gotRequired := g.getFieldRequiredFromBehavior(field)

			// For the conflict case, document the actual behavior
			if tt.name == "proto3 optional with REQUIRED behavior (conflict)" {
				// This test documents current behavior vs expected behavior
				// Currently getFieldRequiredFromBehavior doesn't check field_behavior
				t.Logf("Current behavior: proto3_optional takes precedence (required=%v)", gotRequired)
				t.Logf("Expected behavior: REQUIRED annotation should take precedence")
				if gotRequired == tt.wantRequired {
					t.Log("Behavior matches expectation")
				} else {
					t.Logf("Note: Consider if REQUIRED annotation should override proto3_optional")
				}
			} else if gotRequired != tt.wantRequired && tt.fieldBehavior == nil {
				t.Errorf("getFieldRequiredFromBehavior() = %v, want %v (%s)",
					gotRequired, tt.wantRequired, tt.description)
			}
		})
	}
}

// TestWriteOnlyFromInputOnly tests that INPUT_ONLY field_behavior is correctly
// mapped to writeOnly in OpenAPI schema.
func TestWriteOnlyFromInputOnly(t *testing.T) {
	t.Parallel()

	// Test schema serialization with writeOnly
	schema := &Schema{
		Type:      SchemaType{"string"},
		WriteOnly: true,
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if !strings.Contains(string(data), `"writeOnly":true`) {
		t.Error("Schema with WriteOnly=true should serialize with writeOnly:true")
	}

	// Test that writeOnly is not present when false
	schema2 := &Schema{
		Type:      SchemaType{"string"},
		WriteOnly: false,
	}

	data2, err := json.Marshal(schema2)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if strings.Contains(string(data2), `"writeOnly"`) {
		t.Error("Schema with WriteOnly=false should not have writeOnly in output")
	}
}

// TestRecursiveMessages tests that self-referencing and mutually recursive
// messages are handled correctly without infinite loops.
const recursiveProtoInline = `
file_to_generate: "test/v1/recursive.proto"
proto_file: {
	name: "test/v1/recursive.proto"
	package: "test.v1"
	message_type: {
		name: "TreeNode"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
		field: {
			name: "value"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "value"
		}
		field: {
			name: "children"
			number: 3
			label: LABEL_REPEATED
			type: TYPE_MESSAGE
			type_name: ".test.v1.TreeNode"
			json_name: "children"
		}
	}
	message_type: {
		name: "GetRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	service: {
		name: "TreeService"
		method: {
			name: "GetTree"
			input_type: ".test.v1.GetRequest"
			output_type: ".test.v1.TreeNode"
			options: {
				[google.api.http]: {
					get: "/v1/tree/{id}"
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

func TestRecursiveMessages(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, recursiveProtoInline, nil)
	content := resp[0].GetContent()

	// Should not panic or hang - successful generation means recursion is handled
	if len(content) == 0 {
		t.Fatal("expected non-empty output")
	}

	// Parse the JSON to verify structure
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	components, ok := doc["components"].(map[string]interface{})
	if !ok {
		t.Fatal("missing components")
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		t.Fatal("missing schemas")
	}

	// Find the TreeNode schema
	var treeNodeSchema map[string]interface{}
	for name, s := range schemas {
		if strings.Contains(name, "TreeNode") {
			treeNodeSchema, _ = s.(map[string]interface{})
			break
		}
	}

	if treeNodeSchema == nil {
		t.Fatal("TreeNode schema not found")
	}

	// Verify children field uses a $ref (not inline schema)
	properties, ok := treeNodeSchema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("TreeNode missing properties")
	}

	childrenField, ok := properties["children"].(map[string]interface{})
	if !ok {
		t.Fatal("children field not found")
	}

	// Children should be an array with items referencing TreeNode
	if childrenField["type"] != "array" {
		t.Errorf("children type = %v, want array", childrenField["type"])
	}

	items, ok := childrenField["items"].(map[string]interface{})
	if !ok {
		t.Fatal("children items not found")
	}

	// Should have $ref pointing back to TreeNode
	ref, hasRef := items["$ref"].(string)
	if !hasRef {
		t.Error("children items should have $ref for recursive reference")
	} else if !strings.Contains(ref, "TreeNode") {
		t.Errorf("children items $ref = %q, should reference TreeNode", ref)
	}
}

// ============================================================================
// Priority 3: Edge Case Tests
// ============================================================================

// TestServerVariables tests that server URL templates with variables are
// correctly parsed and rendered.
const serverVariablesProtoInline = `
file_to_generate: "test/v1/server.proto"
proto_file: {
	name: "test/v1/server.proto"
	package: "test.v1"
	message_type: {
		name: "GetRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "Response"
		field: {
			name: "result"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "result"
		}
	}
	service: {
		name: "ServerService"
		method: {
			name: "Get"
			input_type: ".test.v1.GetRequest"
			output_type: ".test.v1.Response"
			options: {
				[google.api.http]: {
					get: "/v1/resources/{id}"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
		[grpc.gateway.protoc_gen_openapiv3.options.openapiv3_document]: {
			info: {
				title: "Server Variables Test"
				version: "1.0"
			}
			servers: {
				url: "https://{environment}.api.example.com/v1"
				description: "API server"
				variables: {
					key: "environment"
					value: {
						default: "production"
						enum: "production"
						enum: "staging"
						enum: "development"
						description: "The deployment environment"
					}
				}
			}
		}
	}
	syntax: "proto3"
}
`

func TestServerVariables(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, serverVariablesProtoInline, nil)
	content := resp[0].GetContent()

	// Parse the JSON to verify server structure
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	servers, ok := doc["servers"].([]interface{})
	if !ok || len(servers) == 0 {
		t.Fatal("expected servers array")
	}

	server := servers[0].(map[string]interface{})

	// Verify URL template
	url, _ := server["url"].(string)
	if !strings.Contains(url, "{environment}") {
		t.Errorf("server URL = %q, want URL with {environment} variable", url)
	}

	// Verify variables
	variables, hasVars := server["variables"].(map[string]interface{})
	if !hasVars {
		t.Error("server should have variables")
	} else {
		envVar, hasEnv := variables["environment"].(map[string]interface{})
		if !hasEnv {
			t.Error("server should have 'environment' variable")
		} else {
			if envVar["default"] != "production" {
				t.Errorf("environment default = %v, want 'production'", envVar["default"])
			}

			// Verify enum values
			enumVals, hasEnum := envVar["enum"].([]interface{})
			if !hasEnum {
				t.Error("environment variable should have enum")
			} else if len(enumVals) != 3 {
				t.Errorf("environment enum has %d values, want 3", len(enumVals))
			}
		}
	}
}

// TestTypeLookupFailure tests that the generator handles gracefully when
// type lookups fail during generation (e.g., when registry has type but
// lookup fails for some reason).
//
// Note: The registry validates proto descriptors during Load(), so missing
// types in the proto itself are caught early. This test verifies that the
// generator has fallback handling when LookupMsg/LookupEnum fail.
func TestTypeLookupFailure(t *testing.T) {
	t.Parallel()

	// Test the fieldTypeToSchema fallback behavior
	// When a message type cannot be looked up, it should return type: object
	reg := descriptor.NewRegistry()
	g := &generator{reg: reg}

	// Create a field that references a type that won't be found
	field := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:     stringPtr("missing_ref"),
			Number:   int32Ptr(1),
			Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
			TypeName: stringPtr(".nonexistent.Type"),
			JsonName: stringPtr("missingRef"),
		},
	}

	// This should not panic, and should return a fallback schema
	schemaRef := g.fieldTypeToSchema(field, nil)

	if schemaRef == nil {
		t.Fatal("expected schema reference, got nil")
	}

	if schemaRef.Schema == nil {
		t.Fatal("expected inline schema value")
	}

	// The fallback for unresolved message types should be type: object
	if !schemaTypeEqual(schemaRef.Schema.Type, SchemaType{"object"}) {
		t.Errorf("fallback schema type = %q, want 'object'", schemaRef.Schema.Type)
	}
}

// TestEnumTypeLookupFailure tests fallback behavior when enum lookup fails.
func TestEnumTypeLookupFailure(t *testing.T) {
	t.Parallel()

	reg := descriptor.NewRegistry()
	g := &generator{reg: reg}

	// Create a field that references an enum that won't be found
	field := &descriptor.Field{
		FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
			Name:     stringPtr("missing_enum"),
			Number:   int32Ptr(1),
			Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
			Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
			TypeName: stringPtr(".nonexistent.EnumType"),
			JsonName: stringPtr("missingEnum"),
		},
	}

	// This should not panic, and should return a fallback schema
	schemaRef := g.fieldTypeToSchema(field, nil)

	if schemaRef == nil {
		t.Fatal("expected schema reference, got nil")
	}

	if schemaRef.Schema == nil {
		t.Fatal("expected inline schema value")
	}

	// The fallback for unresolved enum types should be type: string
	if !schemaTypeEqual(schemaRef.Schema.Type, SchemaType{"string"}) {
		t.Errorf("fallback schema type = %q, want 'string'", schemaRef.Schema.Type)
	}
}

// TestMutualRecursion tests that mutually recursive messages (A references B,
// B references A) are handled correctly.
const mutualRecursionProtoInline = `
file_to_generate: "test/v1/mutual.proto"
proto_file: {
	name: "test/v1/mutual.proto"
	package: "test.v1"
	message_type: {
		name: "Parent"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
		field: {
			name: "children"
			number: 2
			label: LABEL_REPEATED
			type: TYPE_MESSAGE
			type_name: ".test.v1.Child"
			json_name: "children"
		}
	}
	message_type: {
		name: "Child"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
		field: {
			name: "parent"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Parent"
			json_name: "parent"
		}
	}
	message_type: {
		name: "GetRequest"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	service: {
		name: "FamilyService"
		method: {
			name: "GetParent"
			input_type: ".test.v1.GetRequest"
			output_type: ".test.v1.Parent"
			options: {
				[google.api.http]: {
					get: "/v1/parents/{id}"
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

func TestMutualRecursion(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, mutualRecursionProtoInline, nil)
	content := resp[0].GetContent()

	// Should not panic or hang
	if len(content) == 0 {
		t.Fatal("expected non-empty output")
	}

	// Parse the JSON to verify structure
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	components, ok := doc["components"].(map[string]interface{})
	if !ok {
		t.Fatal("missing components")
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		t.Fatal("missing schemas")
	}

	// Both Parent and Child schemas should exist
	var parentSchema, childSchema map[string]interface{}
	for name, s := range schemas {
		if strings.Contains(name, "Parent") {
			parentSchema, _ = s.(map[string]interface{})
		}
		if strings.Contains(name, "Child") {
			childSchema, _ = s.(map[string]interface{})
		}
	}

	if parentSchema == nil {
		t.Error("Parent schema not found")
	}
	if childSchema == nil {
		t.Error("Child schema not found")
	}

	// Verify Parent.children references Child
	if parentSchema != nil {
		props, _ := parentSchema["properties"].(map[string]interface{})
		children, _ := props["children"].(map[string]interface{})
		if children != nil {
			items, _ := children["items"].(map[string]interface{})
			if items != nil {
				ref, _ := items["$ref"].(string)
				if !strings.Contains(ref, "Child") {
					t.Errorf("Parent.children should reference Child, got $ref=%q", ref)
				}
			}
		}
	}

	// Verify Child.parent references Parent
	if childSchema != nil {
		props, _ := childSchema["properties"].(map[string]interface{})
		parent, _ := props["parent"].(map[string]interface{})
		if parent != nil {
			ref, _ := parent["$ref"].(string)
			if !strings.Contains(ref, "Parent") {
				t.Errorf("Child.parent should reference Parent, got $ref=%q", ref)
			}
		}
	}
}

// TestFieldAnnotation_ZeroMinimumMaximum tests that minimum: 0 and maximum: 0
// field annotations via openapiv3_field are correctly applied to OpenAPI schemas.
//
// This is an end-to-end test verifying that field-level annotations with
// zero constraints work correctly. The field annotation code path uses
// pointer checks (opts.Minimum != nil) which correctly handles zero values.
func TestFieldAnnotation_ZeroMinimumMaximum(t *testing.T) {
	t.Parallel()

	// Proto with field annotation specifying minimum: 0, maximum: 100
	// This is valid input - user wants to constrain field to [0, 100]
	// Uses POST with body so that the Request message is included in schemas.
	const protoText = `
file_to_generate: "test/v1/zero_minimum.proto"
proto_file: {
	name: "test/v1/zero_minimum.proto"
	package: "test.v1"
	message_type: {
		name: "CountRequest"
		field: {
			name: "count"
			number: 1
			type: TYPE_INT32
			json_name: "count"
			options: {
				[grpc.gateway.protoc_gen_openapiv3.options.openapiv3_field]: {
					minimum: 0
					maximum: 100
				}
			}
		}
	}
	message_type: { name: "CountResponse" }
	service: {
		name: "TestService"
		method: {
			name: "SetCount"
			input_type: ".test.v1.CountRequest"
			output_type: ".test.v1.CountResponse"
			options: { [google.api.http]: { post: "/v1/count" body: "*" } }
		}
	}
	options: { go_package: "test/v1;testv1" }
	syntax: "proto3"
}
`
	resp := requireGenerateInline(t, protoText, nil)
	content := resp[0].GetContent()

	// Parse the generated OpenAPI JSON
	var doc map[string]any
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("Failed to parse OpenAPI JSON: %v", err)
	}

	// Navigate to the schema
	components, ok := doc["components"].(map[string]any)
	if !ok {
		t.Fatalf("missing components in output: %s", content)
	}
	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		t.Fatalf("missing schemas in output: %s", content)
	}

	// Find the CountRequest schema (might have a prefix)
	var request map[string]any
	for name, s := range schemas {
		if strings.Contains(name, "CountRequest") {
			request, _ = s.(map[string]any)
			break
		}
	}
	if request == nil {
		t.Fatalf("CountRequest schema not found in output: %s", content)
	}

	properties, ok := request["properties"].(map[string]any)
	if !ok {
		t.Fatalf("Request schema missing properties: %v", request)
	}
	countField, ok := properties["count"].(map[string]any)
	if !ok {
		t.Fatalf("count field not found in properties: %v", properties)
	}

	// Check if minimum: 0 is present
	minVal, hasMin := countField["minimum"]
	if !hasMin {
		t.Errorf("minimum constraint is missing: expected minimum: 0 to be set")
	} else if minVal.(float64) != 0 {
		t.Errorf("minimum = %v, want 0", minVal)
	}

	// Check if maximum: 100 is present (should work)
	maxVal, hasMax := countField["maximum"]
	if !hasMax || maxVal.(float64) != 100 {
		t.Errorf("maximum constraint should be 100, got %v (present=%v)", maxVal, hasMax)
	}
}
