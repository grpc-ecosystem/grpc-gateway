package openapi31

import (
	"encoding/json"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
)

func TestTransformSchema_Nullable(t *testing.T) {
	tests := []struct {
		name     string
		input    *model.Schema
		wantType any // string or []string for nullable
	}{
		{
			name:     "non-nullable string",
			input:    &model.Schema{Type: "string"},
			wantType: "string",
		},
		{
			name:     "nullable string outputs type array",
			input:    &model.Schema{Type: "string", IsNullable: true},
			wantType: []string{"string", "null"},
		},
		{
			name:     "nullable integer outputs type array",
			input:    &model.Schema{Type: "integer", IsNullable: true},
			wantType: []string{"integer", "null"},
		},
		{
			name:     "non-nullable object",
			input:    &model.Schema{Type: "object"},
			wantType: "object",
		},
		{
			name:     "nullable object outputs type array",
			input:    &model.Schema{Type: "object", IsNullable: true},
			wantType: []string{"object", "null"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTransformer().transformSchema(tt.input)
			if got == nil {
				t.Fatal("transformSchema returned nil")
			}

			// Marshal to JSON to verify output format
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}

			// Parse back to check type field
			var result map[string]any
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("json.Unmarshal failed: %v", err)
			}

			typeVal, ok := result["type"]
			if !ok && tt.wantType != nil {
				t.Errorf("expected type field, got none")
				return
			}

			switch want := tt.wantType.(type) {
			case string:
				if typeVal != want {
					t.Errorf("type = %v, want %v", typeVal, want)
				}
			case []string:
				// JSON unmarshals arrays as []interface{}
				arr, ok := typeVal.([]any)
				if !ok {
					t.Errorf("expected type array, got %T: %v", typeVal, typeVal)
					return
				}
				if len(arr) != len(want) {
					t.Errorf("type array length = %d, want %d", len(arr), len(want))
					return
				}
				for i, v := range want {
					if arr[i] != v {
						t.Errorf("type[%d] = %v, want %v", i, arr[i], v)
					}
				}
			}
		})
	}
}

func TestTransformSchemaOrRef_RefSiblings(t *testing.T) {
	tests := []struct {
		name            string
		input           *model.SchemaOrRef
		wantRef         string
		wantSummary     string
		wantDescription string
	}{
		{
			name: "ref without siblings",
			input: &model.SchemaOrRef{
				Ref: "#/components/schemas/User",
			},
			wantRef: "#/components/schemas/User",
		},
		{
			name: "ref with summary (3.1.0 feature)",
			input: &model.SchemaOrRef{
				Ref:     "#/components/schemas/User",
				Summary: "A user object",
			},
			wantRef:     "#/components/schemas/User",
			wantSummary: "A user object",
		},
		{
			name: "ref with description (3.1.0 feature)",
			input: &model.SchemaOrRef{
				Ref:         "#/components/schemas/User",
				Description: "Represents a user in the system",
			},
			wantRef:         "#/components/schemas/User",
			wantDescription: "Represents a user in the system",
		},
		{
			name: "ref with both summary and description",
			input: &model.SchemaOrRef{
				Ref:         "#/components/schemas/User",
				Summary:     "A user object",
				Description: "Represents a user in the system",
			},
			wantRef:         "#/components/schemas/User",
			wantSummary:     "A user object",
			wantDescription: "Represents a user in the system",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTransformer().transformSchemaOrRef(tt.input)
			if got == nil {
				t.Fatal("transformSchemaOrRef returned nil")
			}

			// Marshal to JSON
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}

			var result map[string]any
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("json.Unmarshal failed: %v", err)
			}

			// Check $ref
			if ref, ok := result["$ref"].(string); !ok || ref != tt.wantRef {
				t.Errorf("$ref = %v, want %v", result["$ref"], tt.wantRef)
			}

			// Check summary (3.1.0 allows $ref siblings)
			if tt.wantSummary != "" {
				if summary, ok := result["summary"].(string); !ok || summary != tt.wantSummary {
					t.Errorf("summary = %v, want %v", result["summary"], tt.wantSummary)
				}
			}

			// Check description (3.1.0 allows $ref siblings)
			if tt.wantDescription != "" {
				if desc, ok := result["description"].(string); !ok || desc != tt.wantDescription {
					t.Errorf("description = %v, want %v", result["description"], tt.wantDescription)
				}
			}
		})
	}
}

func TestTransformDocument_SetsVersion(t *testing.T) {
	doc := model.NewDocument("ignored")
	doc.Info.Title = "Test API"
	doc.Info.Version = "1.0.0"

	got := TransformDocument(doc)
	if got == nil {
		t.Fatal("TransformDocument returned nil")
	}

	if got.OpenAPI != "3.1.0" {
		t.Errorf("OpenAPI = %v, want 3.1.0", got.OpenAPI)
	}
}

func TestTransformDocument_PreservesInfo(t *testing.T) {
	doc := model.NewDocument("3.1.0")
	doc.Info.Title = "My API"
	doc.Info.Version = "2.0.0"
	doc.Info.Description = "Test description"
	doc.Info.Summary = "Test summary"

	got := TransformDocument(doc)
	if got == nil {
		t.Fatal("TransformDocument returned nil")
	}

	if got.Info.Title != "My API" {
		t.Errorf("Info.Title = %v, want My API", got.Info.Title)
	}
	if got.Info.Version != "2.0.0" {
		t.Errorf("Info.Version = %v, want 2.0.0", got.Info.Version)
	}
	if got.Info.Description != "Test description" {
		t.Errorf("Info.Description = %v, want Test description", got.Info.Description)
	}
	if got.Info.Summary != "Test summary" {
		t.Errorf("Info.Summary = %v, want Test summary", got.Info.Summary)
	}
}

func TestTransformSchema_PreservesAllFields(t *testing.T) {
	minLen := uint64(1)
	maxLen := uint64(100)
	minimum := float64(0)
	maximum := float64(1000)

	input := &model.Schema{
		Type:        "string",
		Format:      "email",
		Title:       "Email Address",
		Description: "User's email",
		MinLength:   &minLen,
		MaxLength:   &maxLen,
		Minimum:     &minimum,
		Maximum:     &maximum,
		Pattern:     "^[a-z]+@[a-z]+\\.[a-z]+$",
		Enum:        []any{"a@b.com", "c@d.com"},
		Default:     "default@example.com",
		Deprecated:  true,
		ReadOnly:    true,
	}

	got := NewTransformer().transformSchema(input)
	if got == nil {
		t.Fatal("transformSchema returned nil")
	}

	// Marshal and check fields are preserved
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	checks := map[string]any{
		"type":        "string",
		"format":      "email",
		"title":       "Email Address",
		"description": "User's email",
		"minLength":   float64(1),
		"maxLength":   float64(100),
		"minimum":     float64(0),
		"maximum":     float64(1000),
		"pattern":     "^[a-z]+@[a-z]+\\.[a-z]+$",
		"default":     "default@example.com",
		"deprecated":  true,
		"readOnly":    true,
	}

	for field, want := range checks {
		got := result[field]
		if got != want {
			t.Errorf("%s = %v (%T), want %v (%T)", field, got, got, want, want)
		}
	}
}
