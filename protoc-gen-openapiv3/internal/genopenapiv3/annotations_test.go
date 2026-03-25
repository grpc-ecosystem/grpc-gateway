package genopenapiv3

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// testGenerator creates a generator instance for testing convert methods.
func testGenerator() *generator {
	registry := NewAdapterRegistry()
	adapter, _ := registry.Get("3.1.0")
	return &generator{
		openapiVersion: "3.1.0",
		adapter:        adapter,
	}
}

// testGeneratorWithReg creates a generator instance with a registry for testing.
func testGeneratorWithReg(reg *descriptor.Registry) *generator {
	registry := NewAdapterRegistry()
	adapter, _ := registry.Get("3.1.0")
	return &generator{
		reg:            reg,
		openapiVersion: "3.1.0",
		adapter:        adapter,
	}
}

// float64Ptr is a helper to create *float64 for test cases
func float64Ptr(f float64) *float64 {
	return &f
}

// uint64Ptr is a helper to create *uint64 for test cases
func uint64Ptr(u uint64) *uint64 {
	return &u
}

// ptrEqual compares two pointers for equality (both nil or same dereferenced value).
func ptrEqual[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// ptrVal returns the dereferenced value or zero value if nil.
func ptrVal[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// schemaTypeEqual compares two SchemaType slices for equality.
func schemaTypeEqual(a, b SchemaType) bool {
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

// isNullable checks if a schema is nullable by checking the Nullable field.
// The adapter handles converting this to type arrays for OpenAPI 3.1.0 output.
func isNullable(schema *Schema) bool {
	if schema == nil {
		return false
	}
	return schema.Nullable
}

// getExampleValue extracts the example value from an examples map.
// Returns nil if the map is empty or the "example" key doesn't exist.
func getExampleValue(examples map[string]*ExampleRef) any {
	if examples == nil {
		return nil
	}
	if ex, ok := examples["example"]; ok && ex != nil && ex.Value != nil {
		return ex.Value.Value
	}
	return nil
}

// findParamByName finds a parameter by name in a slice of ParameterRef.
// Returns nil if not found.
func findParamByName(params []*ParameterRef, name string) *ParameterRef {
	for _, p := range params {
		if p.Value != nil && p.Value.Name == name {
			return p
		}
	}
	return nil
}

// findParamByRef finds a parameter by $ref in a slice of ParameterRef.
// Returns nil if not found.
func findParamByRef(params []*ParameterRef, ref string) *ParameterRef {
	for _, p := range params {
		if p.Ref == ref {
			return p
		}
	}
	return nil
}

// TestParseExampleValue tests the parseExampleValue function that converts
// string examples from proto annotations into properly typed JSON values.
func TestParseExampleValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "empty string returns nil",
			input:    "",
			expected: nil,
		},
		{
			name:     "JSON object",
			input:    `{"id": 123, "name": "John"}`,
			expected: map[string]any{"id": float64(123), "name": "John"},
		},
		{
			name:     "JSON array",
			input:    `[1, 2, 3]`,
			expected: []any{float64(1), float64(2), float64(3)},
		},
		{
			name:     "JSON number",
			input:    `42`,
			expected: float64(42),
		},
		{
			name:     "JSON float",
			input:    `3.14`,
			expected: float64(3.14),
		},
		{
			name:     "JSON boolean true",
			input:    `true`,
			expected: true,
		},
		{
			name:     "JSON boolean false",
			input:    `false`,
			expected: false,
		},
		{
			name:     "JSON null",
			input:    `null`,
			expected: nil,
		},
		{
			name:     "JSON string",
			input:    `"hello world"`,
			expected: "hello world",
		},
		{
			name:     "plain string (not valid JSON) stays as string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "nested JSON object",
			input:    `{"user": {"id": 1, "tags": ["admin", "user"]}}`,
			expected: map[string]any{"user": map[string]any{"id": float64(1), "tags": []any{"admin", "user"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExampleValue(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseExampleValue(%q) = %v (%T), want %v (%T)",
					tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// TestParseExampleValueJSONOutput verifies that parsed examples serialize correctly to JSON.
func TestParseExampleValueJSONOutput(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedJSON string
	}{
		{
			name:         "object example",
			input:        `{"id": 123}`,
			expectedJSON: `{"id":123}`,
		},
		{
			name:         "array example",
			input:        `[1, 2, 3]`,
			expectedJSON: `[1,2,3]`,
		},
		{
			name:         "number example",
			input:        `42`,
			expectedJSON: `42`,
		},
		{
			name:         "boolean example",
			input:        `true`,
			expectedJSON: `true`,
		},
		{
			name:         "string example from JSON",
			input:        `"hello"`,
			expectedJSON: `"hello"`,
		},
		{
			name:         "plain string stays as string",
			input:        `hello`,
			expectedJSON: `"hello"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExampleValue(tt.input)
			jsonBytes, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}
			if string(jsonBytes) != tt.expectedJSON {
				t.Errorf("JSON output = %s, want %s", string(jsonBytes), tt.expectedJSON)
			}
		})
	}
}

// TestConvertExampleOrReference tests conversion of ExampleOrReference proto to ExampleRef.
func TestConvertExampleOrReference(t *testing.T) {
	tests := []struct {
		name    string
		input   *options.ExampleOrReference
		wantRef string
		wantVal *Example
	}{
		{
			name:    "nil input",
			input:   nil,
			wantRef: "",
			wantVal: nil,
		},
		{
			name: "inline example",
			input: &options.ExampleOrReference{
				Oneof: &options.ExampleOrReference_Example{
					Example: &options.Example{
						Summary:     "Example summary",
						Description: "Example description",
						Value:       `{"id": 123}`,
					},
				},
			},
			wantRef: "",
			wantVal: &Example{
				Summary:     "Example summary",
				Description: "Example description",
				Value:       map[string]any{"id": float64(123)},
			},
		},
		{
			name: "reference to component example",
			input: &options.ExampleOrReference{
				Oneof: &options.ExampleOrReference_Reference{
					Reference: &options.Reference{
						Ref: "#/components/examples/MyExample",
					},
				},
			},
			wantRef: "#/components/examples/MyExample",
			wantVal: nil,
		},
		{
			name: "example with external value",
			input: &options.ExampleOrReference{
				Oneof: &options.ExampleOrReference_Example{
					Example: &options.Example{
						Summary:       "External example",
						ExternalValue: "https://example.com/example.json",
					},
				},
			},
			wantRef: "",
			wantVal: &Example{
				Summary:       "External example",
				ExternalValue: "https://example.com/example.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertExampleOrReference(tt.input)

			if tt.input == nil {
				if got != nil {
					t.Errorf("Expected nil, got %+v", got)
				}
				return
			}

			if got.Ref != tt.wantRef {
				t.Errorf("Ref = %q, want %q", got.Ref, tt.wantRef)
			}

			if tt.wantVal == nil {
				if got.Value != nil {
					t.Errorf("Expected Value to be nil, got %+v", got.Value)
				}
			} else {
				if got.Value == nil {
					t.Errorf("Expected Value to be non-nil")
					return
				}
				if got.Value.Summary != tt.wantVal.Summary {
					t.Errorf("Summary = %q, want %q", got.Value.Summary, tt.wantVal.Summary)
				}
				if got.Value.Description != tt.wantVal.Description {
					t.Errorf("Description = %q, want %q", got.Value.Description, tt.wantVal.Description)
				}
				if got.Value.ExternalValue != tt.wantVal.ExternalValue {
					t.Errorf("ExternalValue = %q, want %q", got.Value.ExternalValue, tt.wantVal.ExternalValue)
				}
				if !reflect.DeepEqual(got.Value.Value, tt.wantVal.Value) {
					t.Errorf("Value = %v, want %v", got.Value.Value, tt.wantVal.Value)
				}
			}
		})
	}
}

// TestConvertExamplesMap tests conversion of examples map with mixed inline and reference entries.
func TestConvertExamplesMap(t *testing.T) {
	input := map[string]*options.ExampleOrReference{
		"inline": {
			Oneof: &options.ExampleOrReference_Example{
				Example: &options.Example{
					Summary: "Inline example",
					Value:   `"test"`,
				},
			},
		},
		"referenced": {
			Oneof: &options.ExampleOrReference_Reference{
				Reference: &options.Reference{
					Ref: "#/components/examples/SharedExample",
				},
			},
		},
	}

	got := convertExamplesMap(input)

	if len(got) != 2 {
		t.Fatalf("Expected 2 examples, got %d", len(got))
	}

	// Check inline example
	if inline, ok := got["inline"]; !ok {
		t.Error("Missing 'inline' example")
	} else {
		if inline.Ref != "" {
			t.Errorf("inline.Ref = %q, want empty", inline.Ref)
		}
		if inline.Value == nil || inline.Value.Summary != "Inline example" {
			t.Errorf("inline.Value.Summary = %v, want 'Inline example'", inline.Value)
		}
	}

	// Check referenced example
	if ref, ok := got["referenced"]; !ok {
		t.Error("Missing 'referenced' example")
	} else {
		if ref.Ref != "#/components/examples/SharedExample" {
			t.Errorf("referenced.Ref = %q, want '#/components/examples/SharedExample'", ref.Ref)
		}
		if ref.Value != nil {
			t.Errorf("referenced.Value should be nil for reference, got %+v", ref.Value)
		}
	}
}

// schemaEqual compares two Schema structs for equality on commonly used fields.
func schemaEqual(t *testing.T, got, want *Schema) {
	t.Helper()
	if !schemaTypeEqual(got.Type, want.Type) {
		t.Errorf("Type = %v, want %v", got.Type, want.Type)
	}
	if got.Format != want.Format {
		t.Errorf("Format = %q, want %q", got.Format, want.Format)
	}
	if got.Pattern != want.Pattern {
		t.Errorf("Pattern = %q, want %q", got.Pattern, want.Pattern)
	}
	if !ptrEqual(got.MinLength, want.MinLength) {
		t.Errorf("MinLength = %v, want %v", ptrVal(got.MinLength), ptrVal(want.MinLength))
	}
	if !ptrEqual(got.MaxLength, want.MaxLength) {
		t.Errorf("MaxLength = %v, want %v", ptrVal(got.MaxLength), ptrVal(want.MaxLength))
	}
	if !ptrEqual(got.Minimum, want.Minimum) {
		t.Errorf("Minimum = %v, want %v", ptrVal(got.Minimum), ptrVal(want.Minimum))
	}
	if !ptrEqual(got.Maximum, want.Maximum) {
		t.Errorf("Maximum = %v, want %v", ptrVal(got.Maximum), ptrVal(want.Maximum))
	}
	if !ptrEqual(got.MultipleOf, want.MultipleOf) {
		t.Errorf("MultipleOf = %v, want %v", ptrVal(got.MultipleOf), ptrVal(want.MultipleOf))
	}
}

// referenceEqual compares two Reference structs for equality.
func referenceEqual(t *testing.T, got, want *Reference) {
	t.Helper()
	if got.Ref != want.Ref {
		t.Errorf("Ref = %q, want %q", got.Ref, want.Ref)
	}
	if got.Summary != want.Summary {
		t.Errorf("Summary = %q, want %q", got.Summary, want.Summary)
	}
	if got.Description != want.Description {
		t.Errorf("Description = %q, want %q", got.Description, want.Description)
	}
}

// schemaRefEqual compares two SchemaRef structs for equality.
// It handles both reference and inline schema cases.
func schemaRefEqual(t *testing.T, got, want *SchemaOrReference) {
	t.Helper()
	if got == nil && want == nil {
		return
	}
	if got == nil {
		t.Error("got nil, want non-nil SchemaRef")
		return
	}
	if want == nil {
		t.Error("got non-nil SchemaRef, want nil")
		return
	}
	if want.Reference != nil {
		if got.Reference == nil {
			t.Error("Reference should not be nil")
			return
		}
		referenceEqual(t, got.Reference, want.Reference)
		if got.Schema != nil {
			t.Error("Value should be nil for reference")
		}
		return
	}
	if want.Schema != nil {
		if got.Schema == nil {
			t.Error("Value should not be nil for inline schema")
			return
		}
		schemaEqual(t, got.Schema, want.Schema)
	}
}

func TestConvertServer(t *testing.T) {
	tests := []struct {
		name     string
		input    *options.Server
		expected *Server
	}{
		{
			name: "basic server",
			input: &options.Server{
				Url:         "https://api.example.com",
				Description: "Production server",
			},
			expected: &Server{
				URL:         "https://api.example.com",
				Description: "Production server",
			},
		},
		{
			name: "server with variables",
			input: &options.Server{
				Url: "https://{environment}.api.example.com",
				Variables: map[string]*options.ServerVariable{
					"environment": {
						Default:     "prod",
						Enum:        []string{"prod", "staging", "dev"},
						Description: "Server environment",
					},
				},
			},
			expected: &Server{
				URL: "https://{environment}.api.example.com",
				Variables: map[string]*ServerVariable{
					"environment": {
						Default:     "prod",
						Enum:        []string{"prod", "staging", "dev"},
						Description: "Server environment",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertServer(tt.input)
			if result.URL != tt.expected.URL {
				t.Errorf("URL = %q, want %q", result.URL, tt.expected.URL)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Description = %q, want %q", result.Description, tt.expected.Description)
			}
			if tt.expected.Variables != nil {
				if result.Variables == nil {
					t.Error("Variables should not be nil")
				} else {
					for key, expected := range tt.expected.Variables {
						if got, ok := result.Variables[key]; !ok {
							t.Errorf("Missing variable %q", key)
						} else {
							if got.Default != expected.Default {
								t.Errorf("Variable %q Default = %q, want %q", key, got.Default, expected.Default)
							}
						}
					}
				}
			}
		})
	}
}

func TestConvertTag(t *testing.T) {
	tests := []struct {
		name     string
		input    *options.Tag
		expected *Tag
	}{
		{
			name: "basic tag",
			input: &options.Tag{
				Name:        "Users",
				Description: "User management operations",
			},
			expected: &Tag{
				Name:        "Users",
				Description: "User management operations",
			},
		},
		{
			name: "tag with external docs",
			input: &options.Tag{
				Name:        "Auth",
				Description: "Authentication operations",
				ExternalDocs: &options.ExternalDocumentation{
					Description: "More info",
					Url:         "https://docs.example.com/auth",
				},
			},
			expected: &Tag{
				Name:        "Auth",
				Description: "Authentication operations",
				ExternalDocs: &ExternalDocumentation{
					Description: "More info",
					URL:         "https://docs.example.com/auth",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTag(tt.input)
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %q, want %q", result.Name, tt.expected.Name)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Description = %q, want %q", result.Description, tt.expected.Description)
			}
			if tt.expected.ExternalDocs != nil {
				if result.ExternalDocs == nil {
					t.Error("ExternalDocs should not be nil")
				} else {
					if result.ExternalDocs.URL != tt.expected.ExternalDocs.URL {
						t.Errorf("ExternalDocs.URL = %q, want %q", result.ExternalDocs.URL, tt.expected.ExternalDocs.URL)
					}
				}
			}
		})
	}
}

func TestConvertExternalDocs(t *testing.T) {
	input := &options.ExternalDocumentation{
		Description: "Additional documentation",
		Url:         "https://docs.example.com",
	}

	result := convertExternalDocs(input)

	if result.Description != "Additional documentation" {
		t.Errorf("Description = %q, want %q", result.Description, "Additional documentation")
	}
	if result.URL != "https://docs.example.com" {
		t.Errorf("URL = %q, want %q", result.URL, "https://docs.example.com")
	}
}

func TestConvertSecurityRequirement(t *testing.T) {
	input := &options.SecurityRequirement{
		SecurityRequirement: map[string]*options.SecurityRequirement_SecurityRequirementValue{
			"oauth2": {
				Scope: []string{"read:users", "write:users"},
			},
			"apiKey": {
				Scope: []string{},
			},
		},
	}

	result := convertSecurityRequirement(input)

	if scopes, ok := result["oauth2"]; !ok {
		t.Error("Missing oauth2 requirement")
	} else if len(scopes) != 2 {
		t.Errorf("oauth2 scopes count = %d, want %d", len(scopes), 2)
	}

	if scopes, ok := result["apiKey"]; !ok {
		t.Error("Missing apiKey requirement")
	} else if len(scopes) != 0 {
		t.Errorf("apiKey scopes count = %d, want %d", len(scopes), 0)
	}
}

func TestConvertSecurityScheme(t *testing.T) {
	tests := []struct {
		name     string
		input    *options.SecurityScheme
		expected *SecurityScheme
	}{
		{
			name: "api key scheme",
			input: &options.SecurityScheme{
				Type:        options.SecurityScheme_TYPE_API_KEY,
				Name:        "X-API-Key",
				In:          options.SecurityScheme_IN_HEADER,
				Description: "API key authentication",
			},
			expected: &SecurityScheme{
				Type:        "apiKey",
				Name:        "X-API-Key",
				In:          "header",
				Description: "API key authentication",
			},
		},
		{
			name: "http bearer scheme",
			input: &options.SecurityScheme{
				Type:         options.SecurityScheme_TYPE_HTTP,
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT authentication",
			},
			expected: &SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT authentication",
			},
		},
		{
			name: "oauth2 scheme",
			input: &options.SecurityScheme{
				Type: options.SecurityScheme_TYPE_OAUTH2,
				Flows: &options.OAuthFlows{
					AuthorizationCode: &options.OAuthFlow{
						AuthorizationUrl: "https://auth.example.com/authorize",
						TokenUrl:         "https://auth.example.com/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				},
			},
			expected: &SecurityScheme{
				Type: "oauth2",
				Flows: &OAuthFlows{
					AuthorizationCode: &OAuthFlow{
						AuthorizationURL: "https://auth.example.com/authorize",
						TokenURL:         "https://auth.example.com/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertSecurityScheme(tt.input)
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %q, want %q", result.Type, tt.expected.Type)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name = %q, want %q", result.Name, tt.expected.Name)
			}
			if result.In != tt.expected.In {
				t.Errorf("In = %q, want %q", result.In, tt.expected.In)
			}
			if result.Scheme != tt.expected.Scheme {
				t.Errorf("Scheme = %q, want %q", result.Scheme, tt.expected.Scheme)
			}
			if tt.expected.Flows != nil {
				if result.Flows == nil {
					t.Error("Flows should not be nil")
				}
			}
		})
	}
}

func TestConvertResponse(t *testing.T) {
	input := &options.Response{
		Description: "User not found",
		Headers: map[string]*options.Header{
			"X-Request-Id": {
				Description: "Request ID",
				Schema: &options.SchemaOrReference{
					Oneof: &options.SchemaOrReference_Value{
						Value: &options.Schema{
							Type: []string{"string"},
						},
					},
				},
			},
		},
	}

	result := testGenerator().convertResponse(input)

	if result.Description != "User not found" {
		t.Errorf("Description = %q, want %q", result.Description, "User not found")
	}
	if result.Headers == nil {
		t.Fatal("Headers should not be nil")
	}
	if _, ok := result.Headers["X-Request-Id"]; !ok {
		t.Error("Missing X-Request-Id header")
	}
}

func TestConvertSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *options.Schema
		expected *Schema
	}{
		{
			name: "inline string schema",
			input: &options.Schema{
				Type:      []string{"string"},
				MinLength: 1,
				MaxLength: 100,
				Pattern:   "^[a-z]+$",
			},
			expected: &Schema{
				Type:      SchemaType{"string"},
				MinLength: uint64Ptr(1),
				MaxLength: uint64Ptr(100),
				Pattern:   "^[a-z]+$",
			},
		},
		{
			name: "schema with validation",
			input: &options.Schema{
				Type:       []string{"integer"},
				Minimum:    float64Ptr(0),
				Maximum:    float64Ptr(100),
				MultipleOf: float64Ptr(5),
			},
			expected: &Schema{
				Type:       SchemaType{"integer"},
				Minimum:    float64Ptr(0),
				Maximum:    float64Ptr(100),
				MultipleOf: float64Ptr(5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testGenerator().convertSchema(tt.input)
			if result.Schema == nil {
				t.Fatal("Value should not be nil for inline schema")
			}
			schemaEqual(t, result.Schema, tt.expected)
		})
	}
}

func TestConvertSchemaOrReference(t *testing.T) {
	tests := []struct {
		name       string
		input      *options.SchemaOrReference
		wantNil    bool
		wantRef    *Reference
		wantSchema *Schema
	}{
		{
			name: "reference",
			input: &options.SchemaOrReference{
				Oneof: &options.SchemaOrReference_Reference{
					Reference: &options.Reference{
						Ref:         "#/components/schemas/User",
						Summary:     "User reference",
						Description: "Reference to User schema",
					},
				},
			},
			wantRef: &Reference{
				Ref:         "#/components/schemas/User",
				Summary:     "User reference",
				Description: "Reference to User schema",
			},
		},
		{
			name: "inline string schema",
			input: &options.SchemaOrReference{
				Oneof: &options.SchemaOrReference_Value{
					Value: &options.Schema{
						Type:    []string{"string"},
						Pattern: "^[a-z]+$",
					},
				},
			},
			wantSchema: &Schema{
				Type:    SchemaType{"string"},
				Pattern: "^[a-z]+$",
			},
		},
		{
			name: "inline integer schema with format",
			input: &options.SchemaOrReference{
				Oneof: &options.SchemaOrReference_Value{
					Value: &options.Schema{
						Type:    []string{"integer"},
						Format:  "int64",
						Minimum: float64Ptr(0),
						Maximum: float64Ptr(100),
					},
				},
			},
			wantSchema: &Schema{
				Type:    SchemaType{"integer"},
				Format:  "int64",
				Minimum: float64Ptr(0),
				Maximum: float64Ptr(100),
			},
		},
		{
			name:    "nil input",
			input:   nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testGenerator().convertSchemaOrReference(tt.input)

			if tt.wantNil {
				if result != nil {
					t.Error("Expected nil result for nil input")
				}
				return
			}

			if result == nil {
				t.Fatal("Unexpected nil result")
			}

			if tt.wantRef != nil {
				if result.Reference == nil {
					t.Fatal("Reference should not be nil")
				}
				referenceEqual(t, result.Reference, tt.wantRef)
				if result.Schema != nil {
					t.Error("Value should be nil for reference")
				}
				return
			}

			if tt.wantSchema != nil {
				if result.Schema == nil {
					t.Fatal("Value should not be nil for inline schema")
				}
				schemaEqual(t, result.Schema, tt.wantSchema)
			}
		})
	}
}

// TestConvertSchema_ZeroMinimumMaximum verifies that explicitly setting
// minimum: 0 or maximum: 0 in a Schema annotation is correctly preserved.
func TestConvertSchema_ZeroMinimumMaximum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       *options.Schema
		wantMinimum *float64
		wantMaximum *float64
	}{
		{
			name: "minimum 0, maximum 100",
			input: &options.Schema{
				Type:    []string{"integer"},
				Minimum: float64Ptr(0),
				Maximum: float64Ptr(100),
			},
			wantMinimum: float64Ptr(0),
			wantMaximum: float64Ptr(100),
		},
		{
			name: "minimum -10, maximum 0",
			input: &options.Schema{
				Type:    []string{"integer"},
				Minimum: float64Ptr(-10),
				Maximum: float64Ptr(0),
			},
			wantMinimum: float64Ptr(-10),
			wantMaximum: float64Ptr(0),
		},
		{
			name: "both zero",
			input: &options.Schema{
				Type:    []string{"integer"},
				Minimum: float64Ptr(0),
				Maximum: float64Ptr(0),
			},
			wantMinimum: float64Ptr(0),
			wantMaximum: float64Ptr(0),
		},
		{
			name: "non-zero values",
			input: &options.Schema{
				Type:    []string{"integer"},
				Minimum: float64Ptr(1),
				Maximum: float64Ptr(100),
			},
			wantMinimum: float64Ptr(1),
			wantMaximum: float64Ptr(100),
		},
		{
			name: "only minimum set to 0",
			input: &options.Schema{
				Type:    []string{"integer"},
				Minimum: float64Ptr(0),
			},
			wantMinimum: float64Ptr(0),
			wantMaximum: nil,
		},
		{
			name: "only maximum set to 0",
			input: &options.Schema{
				Type:    []string{"integer"},
				Maximum: float64Ptr(0),
			},
			wantMinimum: nil,
			wantMaximum: float64Ptr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testGenerator().convertSchema(tt.input)

			if result.Schema == nil {
				t.Fatal("Value should not be nil for inline schema")
			}

			// Check minimum
			if tt.wantMinimum != nil {
				if result.Schema.Minimum == nil {
					t.Errorf("minimum constraint is missing: expected minimum: %v to be set",
						*tt.wantMinimum)
				} else if *result.Schema.Minimum != *tt.wantMinimum {
					t.Errorf("minimum = %v, want %v", *result.Schema.Minimum, *tt.wantMinimum)
				}
			} else {
				if result.Schema.Minimum != nil {
					t.Errorf("minimum should be nil, got %v", *result.Schema.Minimum)
				}
			}

			// Check maximum
			if tt.wantMaximum != nil {
				if result.Schema.Maximum == nil {
					t.Errorf("maximum constraint is missing: expected maximum: %v to be set",
						*tt.wantMaximum)
				} else if *result.Schema.Maximum != *tt.wantMaximum {
					t.Errorf("maximum = %v, want %v", *result.Schema.Maximum, *tt.wantMaximum)
				}
			} else {
				if result.Schema.Maximum != nil {
					t.Errorf("maximum should be nil, got %v", *result.Schema.Maximum)
				}
			}
		})
	}
}

func TestConvertParameter(t *testing.T) {
	input := &options.Parameter{
		Name:        "user_id",
		In:          "path",
		Description: "User identifier",
		Required:    true,
		Schema: &options.SchemaOrReference{
			Oneof: &options.SchemaOrReference_Value{
				Value: &options.Schema{
					Type:   []string{"string"},
					Format: "uuid",
				},
			},
		},
	}

	result := testGenerator().convertParameter(input)

	if result.Name != "user_id" {
		t.Errorf("Name = %q, want %q", result.Name, "user_id")
	}
	if result.In != "path" {
		t.Errorf("In = %q, want %q", result.In, "path")
	}
	if result.Description != "User identifier" {
		t.Errorf("Description = %q, want %q", result.Description, "User identifier")
	}
	if !result.Required {
		t.Error("Required should be true")
	}
	if result.Schema == nil {
		t.Fatal("Schema should not be nil")
	}
}

func TestConvertRequestBody(t *testing.T) {
	input := &options.RequestBody{
		Description: "User data",
		Required:    true,
		Content: map[string]*options.MediaType{
			"application/json": {
				Schema: &options.SchemaOrReference{
					Oneof: &options.SchemaOrReference_Value{
						Value: &options.Schema{
							Type: []string{"object"},
						},
					},
				},
			},
		},
	}

	result := testGenerator().convertRequestBody(input)

	if result.Description != "User data" {
		t.Errorf("Description = %q, want %q", result.Description, "User data")
	}
	if !result.Required {
		t.Error("Required should be true")
	}
	if result.Content == nil {
		t.Fatal("Content should not be nil")
	}
	if _, ok := result.Content["application/json"]; !ok {
		t.Error("Missing application/json content")
	}
}

func TestConvertHeader(t *testing.T) {
	input := &options.Header{
		Description: "Request correlation ID",
		Required:    true,
		Schema: &options.SchemaOrReference{
			Oneof: &options.SchemaOrReference_Value{
				Value: &options.Schema{
					Type: []string{"string"},
				},
			},
		},
	}

	result := testGenerator().convertHeader(input)

	if result.Description != "Request correlation ID" {
		t.Errorf("Description = %q, want %q", result.Description, "Request correlation ID")
	}
	if !result.Required {
		t.Error("Required should be true")
	}
	if result.Schema == nil {
		t.Fatal("Schema should not be nil")
	}
}

func TestConvertSchemaComposition(t *testing.T) {
	t.Run("oneOf schema", func(t *testing.T) {
		input := &options.Schema{
			OneOf: []*options.SchemaOrReference{
				{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}, Format: "email"}}},
				{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"integer"}, Format: "int64"}}},
			},
		}
		wantOneOf := []*SchemaOrReference{
			{Schema: &Schema{Type: SchemaType{"string"}, Format: "email"}},
			{Schema: &Schema{Type: SchemaType{"integer"}, Format: "int64"}},
		}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Schema.OneOf) != len(wantOneOf) {
			t.Fatalf("OneOf count = %d, want %d", len(result.Schema.OneOf), len(wantOneOf))
		}
		for i, want := range wantOneOf {
			schemaRefEqual(t, result.Schema.OneOf[i], want)
		}
	})

	t.Run("anyOf schema", func(t *testing.T) {
		input := &options.Schema{
			AnyOf: []*options.SchemaOrReference{
				{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}, Pattern: "^[a-z]+$"}}},
				{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"number"}, Format: "double"}}},
			},
		}
		wantAnyOf := []*SchemaOrReference{
			{Schema: &Schema{Type: SchemaType{"string"}, Pattern: "^[a-z]+$"}},
			{Schema: &Schema{Type: SchemaType{"number"}, Format: "double"}},
		}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Schema.AnyOf) != len(wantAnyOf) {
			t.Fatalf("AnyOf count = %d, want %d", len(result.Schema.AnyOf), len(wantAnyOf))
		}
		for i, want := range wantAnyOf {
			schemaRefEqual(t, result.Schema.AnyOf[i], want)
		}
	})

	t.Run("allOf schema with reference", func(t *testing.T) {
		input := &options.Schema{
			AllOf: []*options.SchemaOrReference{
				{Oneof: &options.SchemaOrReference_Reference{Reference: &options.Reference{Ref: "#/components/schemas/Base"}}},
				{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"object"}}}},
			},
		}
		wantAllOf := []*SchemaOrReference{
			{Reference: &Reference{Ref: "#/components/schemas/Base"}},
			{Schema: &Schema{Type: SchemaType{"object"}}},
		}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Schema.AllOf) != len(wantAllOf) {
			t.Fatalf("AllOf count = %d, want %d", len(result.Schema.AllOf), len(wantAllOf))
		}
		for i, want := range wantAllOf {
			schemaRefEqual(t, result.Schema.AllOf[i], want)
		}
	})

	t.Run("not schema", func(t *testing.T) {
		input := &options.Schema{
			Type: []string{"string"},
			Not:  &options.SchemaOrReference{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}, Pattern: "^forbidden$"}}},
		}
		wantNot := &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}, Pattern: "^forbidden$"}}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		schemaRefEqual(t, result.Schema.Not, wantNot)
	})

	t.Run("discriminator", func(t *testing.T) {
		input := &options.Schema{
			OneOf: []*options.SchemaOrReference{
				{Oneof: &options.SchemaOrReference_Reference{Reference: &options.Reference{Ref: "#/components/schemas/Cat"}}},
				{Oneof: &options.SchemaOrReference_Reference{Reference: &options.Reference{Ref: "#/components/schemas/Dog"}}},
			},
			Discriminator: &options.Discriminator{
				PropertyName: "petType",
				Mapping: map[string]string{
					"cat": "#/components/schemas/Cat",
					"dog": "#/components/schemas/Dog",
				},
			},
		}
		wantOneOf := []*SchemaOrReference{
			{Reference: &Reference{Ref: "#/components/schemas/Cat"}},
			{Reference: &Reference{Ref: "#/components/schemas/Dog"}},
		}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Schema.OneOf) != len(wantOneOf) {
			t.Fatalf("OneOf count = %d, want %d", len(result.Schema.OneOf), len(wantOneOf))
		}
		for i, want := range wantOneOf {
			schemaRefEqual(t, result.Schema.OneOf[i], want)
		}
		if result.Schema.Discriminator == nil {
			t.Fatal("Discriminator should not be nil")
		}
		if result.Schema.Discriminator.PropertyName != "petType" {
			t.Errorf("Discriminator.PropertyName = %q, want %q", result.Schema.Discriminator.PropertyName, "petType")
		}
		if len(result.Schema.Discriminator.Mapping) != 2 {
			t.Errorf("Discriminator.Mapping count = %d, want %d", len(result.Schema.Discriminator.Mapping), 2)
		}
	})

	t.Run("items for array", func(t *testing.T) {
		input := &options.Schema{
			Type:  []string{"array"},
			Items: &options.SchemaOrReference{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}, Format: "uuid"}}},
		}
		wantItems := &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}, Format: "uuid"}}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		schemaRefEqual(t, result.Schema.Items, wantItems)
	})

	t.Run("properties for object", func(t *testing.T) {
		input := &options.Schema{
			Type: []string{"object"},
			Properties: []*options.NamedSchemaOrReference{
				{Name: "name", Value: &options.SchemaOrReference{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}, MinLength: 1}}}},
				{Name: "age", Value: &options.SchemaOrReference{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"integer"}, Format: "int32"}}}},
			},
		}
		wantProperties := map[string]*SchemaOrReference{
			"name": {Schema: &Schema{Type: SchemaType{"string"}, MinLength: uint64Ptr(1)}},
			"age":  {Schema: &Schema{Type: SchemaType{"integer"}, Format: "int32"}},
		}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Schema.Properties) != len(wantProperties) {
			t.Fatalf("Properties count = %d, want %d", len(result.Schema.Properties), len(wantProperties))
		}
		for name, want := range wantProperties {
			got, ok := result.Schema.Properties[name]
			if !ok {
				t.Errorf("missing property %q", name)
				continue
			}
			schemaRefEqual(t, got, want)
		}
	})

	t.Run("additionalProperties with schema", func(t *testing.T) {
		input := &options.Schema{
			Type: []string{"object"},
			AdditionalProperties: &options.AdditionalPropertiesItem{
				Kind: &options.AdditionalPropertiesItem_SchemaOrReference{
					SchemaOrReference: &options.SchemaOrReference{
						Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}, Format: "uri"}},
					},
				},
			},
		}
		wantAdditionalProperties := &SchemaOrReference{Schema: &Schema{Type: SchemaType{"string"}, Format: "uri"}}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		schemaRefEqual(t, result.Schema.AdditionalProperties, wantAdditionalProperties)
	})

	t.Run("additionalProperties allows", func(t *testing.T) {
		input := &options.Schema{
			Type: []string{"object"},
			AdditionalProperties: &options.AdditionalPropertiesItem{
				Kind: &options.AdditionalPropertiesItem_Allows{Allows: true},
			},
		}

		result := testGenerator().convertSchema(input)

		if result.Schema == nil {
			t.Fatal("Value should not be nil")
		}
		if result.Schema.AdditionalProperties == nil {
			t.Fatal("AdditionalProperties should not be nil when allows=true")
		}
	})
}

// ============================================================================
// Apply Annotation Tests
// ============================================================================

func TestApplyInfoAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        *options.Info
		wantTitle   string
		wantVersion string
		wantSummary string
		wantDesc    string
		wantTerms   string
		wantContact bool
		wantLicense bool
	}{
		{
			name: "apply all info fields",
			opts: &options.Info{
				Title:          "My API",
				Summary:        "A brief summary",
				Description:    "Full description",
				TermsOfService: "https://example.com/tos",
				Version:        "2.0.0",
				Contact: &options.Contact{
					Name:  "Support",
					Url:   "https://support.example.com",
					Email: "support@example.com",
				},
				License: &options.License{
					Name:       "MIT",
					Identifier: "MIT",
					Url:        "https://opensource.org/licenses/MIT",
				},
			},
			wantTitle:   "My API",
			wantVersion: "2.0.0",
			wantSummary: "A brief summary",
			wantDesc:    "Full description",
			wantTerms:   "https://example.com/tos",
			wantContact: true,
			wantLicense: true,
		},
		{
			name: "partial info update",
			opts: &options.Info{
				Title: "Updated Title",
			},
			wantTitle:   "Updated Title",
			wantVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := &Info{
				Title:   "Original Title",
				Version: "1.0.0",
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applyInfoAnnotation(info, tt.opts)

			if tt.wantTitle != "" && info.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", info.Title, tt.wantTitle)
			}
			if tt.wantVersion != "" && info.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", info.Version, tt.wantVersion)
			}
			if tt.wantSummary != "" && info.Summary != tt.wantSummary {
				t.Errorf("Summary = %q, want %q", info.Summary, tt.wantSummary)
			}
			if tt.wantDesc != "" && info.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", info.Description, tt.wantDesc)
			}
			if tt.wantTerms != "" && info.TermsOfService != tt.wantTerms {
				t.Errorf("TermsOfService = %q, want %q", info.TermsOfService, tt.wantTerms)
			}
			if tt.wantContact && info.Contact == nil {
				t.Error("Contact should not be nil")
			}
			if tt.wantLicense && info.License == nil {
				t.Error("License should not be nil")
			}
		})
	}
}

func TestApplySchemaAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		opts           *options.Schema
		wantTitle      string
		wantDesc       string
		wantExample    any
		wantReadOnly   bool
		wantWriteOnly  bool
		wantNullable   bool
		wantDeprecated bool
		wantRequired   []string
		wantAllOf      int
		wantAnyOf      int
		wantOneOf      int
	}{
		{
			name: "apply title and description",
			opts: &options.Schema{
				Title:       "User Schema",
				Description: "A user object",
			},
			wantTitle: "User Schema",
			wantDesc:  "A user object",
		},
		{
			name: "apply read/write only",
			opts: &options.Schema{
				ReadOnly:  true,
				WriteOnly: false,
			},
			wantReadOnly:  true,
			wantWriteOnly: false,
		},
		{
			name: "apply nullable and deprecated",
			opts: &options.Schema{
				Nullable:   true,
				Deprecated: true,
			},
			wantNullable:   true,
			wantDeprecated: true,
		},
		{
			name: "apply required fields",
			opts: &options.Schema{
				Required: []string{"id", "name"},
			},
			wantRequired: []string{"id", "name"},
		},
		{
			name: "apply example",
			opts: &options.Schema{
				Example: `{"id": "123"}`,
			},
			// Example is now parsed as JSON, so we expect a map, not a string
			wantExample: map[string]any{"id": "123"},
		},
		{
			name: "apply composition types",
			opts: &options.Schema{
				AllOf: []*options.SchemaOrReference{
					{Oneof: &options.SchemaOrReference_Reference{Reference: &options.Reference{Ref: "#/components/schemas/Base"}}},
				},
				AnyOf: []*options.SchemaOrReference{
					{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"string"}}}},
					{Oneof: &options.SchemaOrReference_Value{Value: &options.Schema{Type: []string{"integer"}}}},
				},
				OneOf: []*options.SchemaOrReference{
					{Oneof: &options.SchemaOrReference_Reference{Reference: &options.Reference{Ref: "#/components/schemas/Cat"}}},
					{Oneof: &options.SchemaOrReference_Reference{Reference: &options.Reference{Ref: "#/components/schemas/Dog"}}},
				},
			},
			wantAllOf: 1,
			wantAnyOf: 2,
			wantOneOf: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			schema := &Schema{Type: SchemaType{"object"}}

			// Create message options with the OpenAPI annotation
			msgOpts := &descriptorpb.MessageOptions{}
			proto.SetExtension(msgOpts, options.E_Openapiv3Schema, tt.opts)

			// Create a mock message with the annotation
			msg := &descriptor.Message{
				DescriptorProto: &descriptorpb.DescriptorProto{
					Name:    stringPtr("TestMessage"),
					Options: msgOpts,
				},
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applySchemaAnnotation(schema, msg)

			// Assertions
			if tt.wantTitle != "" && schema.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", schema.Title, tt.wantTitle)
			}
			if tt.wantDesc != "" && schema.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", schema.Description, tt.wantDesc)
			}
			if tt.wantExample != nil {
				gotExample := getExampleValue(schema.Examples)
				if !reflect.DeepEqual(gotExample, tt.wantExample) {
					t.Errorf("Example = %v, want %v", gotExample, tt.wantExample)
				}
			}
			if schema.ReadOnly != tt.wantReadOnly {
				t.Errorf("ReadOnly = %v, want %v", schema.ReadOnly, tt.wantReadOnly)
			}
			if schema.WriteOnly != tt.wantWriteOnly {
				t.Errorf("WriteOnly = %v, want %v", schema.WriteOnly, tt.wantWriteOnly)
			}
			// Check nullable via type array (OpenAPI 3.1.0 style)
			if isNullable(schema) != tt.wantNullable {
				t.Errorf("Nullable (via type array) = %v, want %v", isNullable(schema), tt.wantNullable)
			}
			if schema.Deprecated != tt.wantDeprecated {
				t.Errorf("Deprecated = %v, want %v", schema.Deprecated, tt.wantDeprecated)
			}
			if len(tt.wantRequired) > 0 && len(schema.Required) != len(tt.wantRequired) {
				t.Errorf("Required count = %d, want %d", len(schema.Required), len(tt.wantRequired))
			}
			if tt.wantAllOf > 0 && len(schema.AllOf) != tt.wantAllOf {
				t.Errorf("AllOf count = %d, want %d", len(schema.AllOf), tt.wantAllOf)
			}
			if tt.wantAnyOf > 0 && len(schema.AnyOf) != tt.wantAnyOf {
				t.Errorf("AnyOf count = %d, want %d", len(schema.AnyOf), tt.wantAnyOf)
			}
			if tt.wantOneOf > 0 && len(schema.OneOf) != tt.wantOneOf {
				t.Errorf("OneOf count = %d, want %d", len(schema.OneOf), tt.wantOneOf)
			}
		})
	}
}

func TestApplyFieldAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		opts           *options.Schema
		wantTitle      string
		wantDesc       string
		wantDefault    string
		wantExample    any
		wantFormat     string
		wantPattern    string
		wantMinLength  uint64
		wantMaxLength  uint64
		wantMinimum    float64
		wantMaximum    float64
		wantReadOnly   bool
		wantWriteOnly  bool
		wantNullable   bool
		wantDeprecated bool
	}{
		{
			name: "string field with validation",
			opts: &options.Schema{
				Title:       "Email",
				Description: "User email address",
				Format:      "email",
				Pattern:     "^[\\w-\\.]+@[\\w-]+\\.[a-z]{2,}$",
				MinLength:   5,
				MaxLength:   100,
			},
			wantTitle:     "Email",
			wantDesc:      "User email address",
			wantFormat:    "email",
			wantPattern:   "^[\\w-\\.]+@[\\w-]+\\.[a-z]{2,}$",
			wantMinLength: 5,
			wantMaxLength: 100,
		},
		{
			name: "numeric field with constraints",
			opts: &options.Schema{
				Minimum:    float64Ptr(0),
				Maximum:    float64Ptr(100),
				MultipleOf: float64Ptr(5),
			},
			wantMinimum: 0,
			wantMaximum: 100,
		},
		{
			name: "field with default and example",
			opts: &options.Schema{
				Default: "active",
				Example: "pending",
			},
			wantDefault: "active",
			wantExample: "pending",
		},
		{
			name: "read-only field",
			opts: &options.Schema{
				ReadOnly: true,
			},
			wantReadOnly: true,
		},
		{
			name: "write-only field",
			opts: &options.Schema{
				WriteOnly: true,
			},
			wantWriteOnly: true,
		},
		{
			name: "nullable deprecated field",
			opts: &options.Schema{
				Nullable:   true,
				Deprecated: true,
			},
			wantNullable:   true,
			wantDeprecated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			schema := &Schema{Type: SchemaType{"string"}}

			// Create field options with the OpenAPI annotation
			fieldOpts := &descriptorpb.FieldOptions{}
			proto.SetExtension(fieldOpts, options.E_Openapiv3Field, tt.opts)

			// Create a mock field with the annotation
			field := &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:    stringPtr("test_field"),
					Options: fieldOpts,
				},
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applyFieldAnnotation(schema, field)

			// Assertions
			if tt.wantTitle != "" && schema.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", schema.Title, tt.wantTitle)
			}
			if tt.wantDesc != "" && schema.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", schema.Description, tt.wantDesc)
			}
			if tt.wantDefault != "" && schema.Default != tt.wantDefault {
				t.Errorf("Default = %q, want %q", schema.Default, tt.wantDefault)
			}
			if tt.wantExample != nil {
				gotExample := getExampleValue(schema.Examples)
				if !reflect.DeepEqual(gotExample, tt.wantExample) {
					t.Errorf("Example = %v, want %v", gotExample, tt.wantExample)
				}
			}
			if tt.wantFormat != "" && schema.Format != tt.wantFormat {
				t.Errorf("Format = %q, want %q", schema.Format, tt.wantFormat)
			}
			if tt.wantPattern != "" && schema.Pattern != tt.wantPattern {
				t.Errorf("Pattern = %q, want %q", schema.Pattern, tt.wantPattern)
			}
			if tt.wantMinLength > 0 && (schema.MinLength == nil || *schema.MinLength != tt.wantMinLength) {
				t.Errorf("MinLength = %v, want %v", schema.MinLength, tt.wantMinLength)
			}
			if tt.wantMaxLength > 0 && (schema.MaxLength == nil || *schema.MaxLength != tt.wantMaxLength) {
				t.Errorf("MaxLength = %v, want %v", schema.MaxLength, tt.wantMaxLength)
			}
			if schema.ReadOnly != tt.wantReadOnly {
				t.Errorf("ReadOnly = %v, want %v", schema.ReadOnly, tt.wantReadOnly)
			}
			if schema.WriteOnly != tt.wantWriteOnly {
				t.Errorf("WriteOnly = %v, want %v", schema.WriteOnly, tt.wantWriteOnly)
			}
			// Check nullable via type array (OpenAPI 3.1.0 style)
			if isNullable(schema) != tt.wantNullable {
				t.Errorf("Nullable (via type array) = %v, want %v", isNullable(schema), tt.wantNullable)
			}
			if schema.Deprecated != tt.wantDeprecated {
				t.Errorf("Deprecated = %v, want %v", schema.Deprecated, tt.wantDeprecated)
			}
		})
	}
}

func TestApplyOperationAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		opts             *options.Operation
		wantSummary      string
		wantDesc         string
		wantOpID         string
		wantTags         []string
		wantDeprecated   bool
		wantSecurity     int
		wantServers      int
		wantHeaderParams int
		wantCookieParams int
		verifyParams     func(t *testing.T, params []*ParameterRef)
	}{
		{
			name: "override summary and description",
			opts: &options.Operation{
				Summary:     "Get a user",
				Description: "Retrieves a user by ID",
			},
			wantSummary: "Get a user",
			wantDesc:    "Retrieves a user by ID",
		},
		{
			name: "override operation ID",
			opts: &options.Operation{
				OperationId: "getUserById",
			},
			wantOpID: "getUserById",
		},
		{
			name: "override tags",
			opts: &options.Operation{
				Tags: []string{"Users", "Admin"},
			},
			wantTags: []string{"Users", "Admin"},
		},
		{
			name: "mark deprecated",
			opts: &options.Operation{
				Deprecated: true,
			},
			wantDeprecated: true,
		},
		{
			name: "add security requirements",
			opts: &options.Operation{
				Security: []*options.SecurityRequirement{
					{
						SecurityRequirement: map[string]*options.SecurityRequirement_SecurityRequirementValue{
							"oauth2": {Scope: []string{"read:users"}},
						},
					},
				},
			},
			wantSecurity: 1,
		},
		{
			name: "add servers",
			opts: &options.Operation{
				Servers: []*options.Server{
					{Url: "https://api.example.com"},
				},
			},
			wantServers: 1,
		},
		{
			name: "add header parameters",
			opts: &options.Operation{
				Parameters: &options.OperationParameters{
					Headers: []*options.HeaderParameterOrReference{
						{
							Oneof: &options.HeaderParameterOrReference_Header{
								Header: &options.HeaderParameter{
									Name:        "X-Request-ID",
									Description: "Request tracking ID",
									Required:    true,
								},
							},
						},
						{
							Oneof: &options.HeaderParameterOrReference_Header{
								Header: &options.HeaderParameter{
									Name:        "X-Api-Version",
									Description: "API version header",
								},
							},
						},
					},
				},
			},
			wantHeaderParams: 2,
			verifyParams: func(t *testing.T, params []*ParameterRef) {
				t.Helper()
				// Find X-Request-ID header
				p1 := findParamByName(params, "X-Request-ID")
				if p1 == nil {
					t.Fatal("X-Request-ID header not found")
				}
				if p1.Value.In != "header" {
					t.Errorf("X-Request-ID in = %q, want header", p1.Value.In)
				}
				if !p1.Value.Required {
					t.Error("X-Request-ID should be required")
				}
				if p1.Value.Description != "Request tracking ID" {
					t.Errorf("X-Request-ID description = %q, want %q", p1.Value.Description, "Request tracking ID")
				}
				// Find X-Api-Version header
				p2 := findParamByName(params, "X-Api-Version")
				if p2 == nil {
					t.Fatal("X-Api-Version header not found")
				}
				if p2.Value.In != "header" {
					t.Errorf("X-Api-Version in = %q, want header", p2.Value.In)
				}
				if p2.Value.Description != "API version header" {
					t.Errorf("X-Api-Version description = %q, want %q", p2.Value.Description, "API version header")
				}
			},
		},
		{
			name: "add cookie parameters",
			opts: &options.Operation{
				Parameters: &options.OperationParameters{
					Cookies: []*options.CookieParameterOrReference{
						{
							Oneof: &options.CookieParameterOrReference_Cookie{
								Cookie: &options.CookieParameter{
									Name:        "session_id",
									Description: "Session identifier",
									Required:    true,
								},
							},
						},
					},
				},
			},
			wantCookieParams: 1,
			verifyParams: func(t *testing.T, params []*ParameterRef) {
				t.Helper()
				p := findParamByName(params, "session_id")
				if p == nil {
					t.Fatal("session_id cookie not found")
				}
				if p.Value.In != "cookie" {
					t.Errorf("session_id in = %q, want cookie", p.Value.In)
				}
				if !p.Value.Required {
					t.Error("session_id should be required")
				}
				if p.Value.Description != "Session identifier" {
					t.Errorf("session_id description = %q, want %q", p.Value.Description, "Session identifier")
				}
			},
		},
		{
			name: "add header parameter with reference",
			opts: &options.Operation{
				Parameters: &options.OperationParameters{
					Headers: []*options.HeaderParameterOrReference{
						{
							Oneof: &options.HeaderParameterOrReference_Reference{
								Reference: &options.Reference{
									Ref: "#/components/parameters/AuthToken",
								},
							},
						},
					},
				},
			},
			wantHeaderParams: 1,
			verifyParams: func(t *testing.T, params []*ParameterRef) {
				t.Helper()
				p := findParamByRef(params, "#/components/parameters/AuthToken")
				if p == nil {
					t.Fatal("parameter with ref #/components/parameters/AuthToken not found")
				}
			},
		},
		{
			name: "add mixed header and cookie parameters",
			opts: &options.Operation{
				Parameters: &options.OperationParameters{
					Headers: []*options.HeaderParameterOrReference{
						{
							Oneof: &options.HeaderParameterOrReference_Header{
								Header: &options.HeaderParameter{
									Name:        "X-Request-ID",
									Description: "Request ID",
								},
							},
						},
					},
					Cookies: []*options.CookieParameterOrReference{
						{
							Oneof: &options.CookieParameterOrReference_Cookie{
								Cookie: &options.CookieParameter{
									Name:        "session",
									Description: "Session cookie",
								},
							},
						},
					},
				},
			},
			wantHeaderParams: 1,
			wantCookieParams: 1,
			verifyParams: func(t *testing.T, params []*ParameterRef) {
				t.Helper()
				// Find header by name
				header := findParamByName(params, "X-Request-ID")
				if header == nil {
					t.Fatal("X-Request-ID header not found")
				}
				if header.Value.In != "header" {
					t.Errorf("X-Request-ID in = %q, want header", header.Value.In)
				}
				if header.Value.Description != "Request ID" {
					t.Errorf("X-Request-ID description = %q, want %q", header.Value.Description, "Request ID")
				}
				// Find cookie by name
				cookie := findParamByName(params, "session")
				if cookie == nil {
					t.Fatal("session cookie not found")
				}
				if cookie.Value.In != "cookie" {
					t.Errorf("session in = %q, want cookie", cookie.Value.In)
				}
				if cookie.Value.Description != "Session cookie" {
					t.Errorf("session description = %q, want %q", cookie.Value.Description, "Session cookie")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			op := &Operation{
				Summary:     "Original summary",
				Description: "Original description",
				OperationID: "originalOpId",
				Tags:        []string{"Original"},
			}

			// Create method options with the OpenAPI annotation
			methodOpts := &descriptorpb.MethodOptions{}
			proto.SetExtension(methodOpts, options.E_Openapiv3Operation, tt.opts)

			// Create a mock method with the annotation
			method := &descriptor.Method{
				MethodDescriptorProto: &descriptorpb.MethodDescriptorProto{
					Name:    stringPtr("TestMethod"),
					Options: methodOpts,
				},
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applyOperationAnnotation(op, method)

			// Assertions
			if tt.wantSummary != "" && op.Summary != tt.wantSummary {
				t.Errorf("Summary = %q, want %q", op.Summary, tt.wantSummary)
			}
			if tt.wantDesc != "" && op.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", op.Description, tt.wantDesc)
			}
			if tt.wantOpID != "" && op.OperationID != tt.wantOpID {
				t.Errorf("OperationID = %q, want %q", op.OperationID, tt.wantOpID)
			}
			if len(tt.wantTags) > 0 && len(op.Tags) != len(tt.wantTags) {
				t.Errorf("Tags count = %d, want %d", len(op.Tags), len(tt.wantTags))
			}
			if op.Deprecated != tt.wantDeprecated {
				t.Errorf("Deprecated = %v, want %v", op.Deprecated, tt.wantDeprecated)
			}
			if tt.wantSecurity > 0 && len(op.Security) != tt.wantSecurity {
				t.Errorf("Security count = %d, want %d", len(op.Security), tt.wantSecurity)
			}
			if tt.wantServers > 0 && len(op.Servers) != tt.wantServers {
				t.Errorf("Servers count = %d, want %d", len(op.Servers), tt.wantServers)
			}
			// Verify header params count
			if tt.wantHeaderParams > 0 || tt.wantCookieParams > 0 {
				expectedTotal := tt.wantHeaderParams + tt.wantCookieParams
				if len(op.Parameters) != expectedTotal {
					t.Errorf("Parameters count = %d, want %d", len(op.Parameters), expectedTotal)
				}
			}
			// Run custom param verification if provided
			if tt.verifyParams != nil {
				tt.verifyParams(t, op.Parameters)
			}
		})
	}
}

func TestApplyServiceAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        *options.Tag
		wantName    string
		wantDesc    string
		wantExtDocs bool
	}{
		{
			name: "override name and description",
			opts: &options.Tag{
				Name:        "Users API",
				Description: "User management operations",
			},
			wantName: "Users API",
			wantDesc: "User management operations",
		},
		{
			name: "add external docs",
			opts: &options.Tag{
				ExternalDocs: &options.ExternalDocumentation{
					Description: "See more",
					Url:         "https://docs.example.com/users",
				},
			},
			wantExtDocs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := &Tag{
				Name:        "OriginalService",
				Description: "Original description",
			}

			// Create service options with the OpenAPI annotation
			svcOpts := &descriptorpb.ServiceOptions{}
			proto.SetExtension(svcOpts, options.E_Openapiv3Tag, tt.opts)

			// Create a mock service with the annotation
			svc := &descriptor.Service{
				ServiceDescriptorProto: &descriptorpb.ServiceDescriptorProto{
					Name:    stringPtr("TestService"),
					Options: svcOpts,
				},
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applyServiceAnnotation(tag, svc)

			// Assertions
			if tt.wantName != "" && tag.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tag.Name, tt.wantName)
			}
			if tt.wantDesc != "" && tag.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", tag.Description, tt.wantDesc)
			}
			if tt.wantExtDocs && tag.ExternalDocs == nil {
				t.Error("ExternalDocs should not be nil")
			}
		})
	}
}

func TestApplyEnumAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		opts           *options.EnumSchema
		wantTitle      string
		wantDesc       string
		wantDefault    string
		wantExample    any
		wantDeprecated bool
		wantExtDocs    bool
	}{
		{
			name: "apply title and description",
			opts: &options.EnumSchema{
				Title:       "Task Status",
				Description: "The status of a task",
			},
			wantTitle: "Task Status",
			wantDesc:  "The status of a task",
		},
		{
			name: "apply default and example",
			opts: &options.EnumSchema{
				Default: "PENDING",
				Example: "COMPLETED",
			},
			wantDefault: "PENDING",
			wantExample: "COMPLETED",
		},
		{
			name: "mark deprecated",
			opts: &options.EnumSchema{
				Deprecated: true,
			},
			wantDeprecated: true,
		},
		{
			name: "add external docs",
			opts: &options.EnumSchema{
				ExternalDocs: &options.ExternalDocumentation{
					Url: "https://docs.example.com/status",
				},
			},
			wantExtDocs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			schema := &Schema{
				Type: SchemaType{"string"},
				Enum: []any{"PENDING", "COMPLETED", "FAILED"},
			}

			// Create enum options with the OpenAPI annotation
			enumOpts := &descriptorpb.EnumOptions{}
			proto.SetExtension(enumOpts, options.E_Openapiv3Enum, tt.opts)

			// Create a mock enum with the annotation
			enum := &descriptor.Enum{
				EnumDescriptorProto: &descriptorpb.EnumDescriptorProto{
					Name:    stringPtr("TestEnum"),
					Options: enumOpts,
				},
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applyEnumAnnotation(schema, enum)

			// Assertions
			if tt.wantTitle != "" && schema.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", schema.Title, tt.wantTitle)
			}
			if tt.wantDesc != "" && schema.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", schema.Description, tt.wantDesc)
			}
			if tt.wantDefault != "" && schema.Default != tt.wantDefault {
				t.Errorf("Default = %q, want %q", schema.Default, tt.wantDefault)
			}
			if tt.wantExample != nil {
				gotExample := getExampleValue(schema.Examples)
				if !reflect.DeepEqual(gotExample, tt.wantExample) {
					t.Errorf("Example = %v, want %v", gotExample, tt.wantExample)
				}
			}
			if schema.Deprecated != tt.wantDeprecated {
				t.Errorf("Deprecated = %v, want %v", schema.Deprecated, tt.wantDeprecated)
			}
			if tt.wantExtDocs && schema.ExternalDocs == nil {
				t.Error("ExternalDocs should not be nil")
			}
		})
	}
}

func TestApplyComponentsAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		opts                *options.Components
		wantSecuritySchemes int
		wantResponses       int
		wantParameters      int
		wantRequestBodies   int
		wantHeaders         int
	}{
		{
			name: "add security schemes",
			opts: &options.Components{
				SecuritySchemes: map[string]*options.SecurityScheme{
					"bearerAuth": {
						Type:         options.SecurityScheme_TYPE_HTTP,
						Scheme:       "bearer",
						BearerFormat: "JWT",
					},
					"apiKey": {
						Type: options.SecurityScheme_TYPE_API_KEY,
						Name: "X-API-Key",
						In:   options.SecurityScheme_IN_HEADER,
					},
				},
			},
			wantSecuritySchemes: 2,
		},
		{
			name: "add responses",
			opts: &options.Components{
				Responses: map[string]*options.Response{
					"NotFound":   {Description: "Resource not found"},
					"BadRequest": {Description: "Invalid request"},
				},
			},
			wantResponses: 2,
		},
		{
			name: "add parameters",
			opts: &options.Components{
				Parameters: map[string]*options.Parameter{
					"PageSize": {
						Name:     "page_size",
						In:       "query",
						Required: false,
						Schema: &options.SchemaOrReference{
							Oneof: &options.SchemaOrReference_Value{
								Value: &options.Schema{Type: []string{"integer"}},
							},
						},
					},
				},
			},
			wantParameters: 1,
		},
		{
			name: "add request bodies",
			opts: &options.Components{
				RequestBodies: map[string]*options.RequestBody{
					"UserInput": {
						Description: "User data",
						Required:    true,
					},
				},
			},
			wantRequestBodies: 1,
		},
		{
			name: "add headers",
			opts: &options.Components{
				Headers: map[string]*options.HeaderOrReference{
					"X-Request-Id": {
						Oneof: &options.HeaderOrReference_Header{
							Header: &options.Header{
								Description: "Request correlation ID",
								Schema: &options.SchemaOrReference{
									Oneof: &options.SchemaOrReference_Value{
										Value: &options.Schema{Type: []string{"string"}},
									},
								},
							},
						},
					},
				},
			},
			wantHeaders: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			comp := &Components{
				Schemas: make(map[string]*SchemaOrReference),
			}

			reg := &descriptor.Registry{}
			gen := testGeneratorWithReg(reg)
			gen.applyComponentsAnnotation(comp, tt.opts)

			if tt.wantSecuritySchemes > 0 && len(comp.SecuritySchemes) != tt.wantSecuritySchemes {
				t.Errorf("SecuritySchemes count = %d, want %d", len(comp.SecuritySchemes), tt.wantSecuritySchemes)
			}
			if tt.wantResponses > 0 && len(comp.Responses) != tt.wantResponses {
				t.Errorf("Responses count = %d, want %d", len(comp.Responses), tt.wantResponses)
			}
			if tt.wantParameters > 0 && len(comp.Parameters) != tt.wantParameters {
				t.Errorf("Parameters count = %d, want %d", len(comp.Parameters), tt.wantParameters)
			}
			if tt.wantRequestBodies > 0 && len(comp.RequestBodies) != tt.wantRequestBodies {
				t.Errorf("RequestBodies count = %d, want %d", len(comp.RequestBodies), tt.wantRequestBodies)
			}
			if tt.wantHeaders > 0 && len(comp.Headers) != tt.wantHeaders {
				t.Errorf("Headers count = %d, want %d", len(comp.Headers), tt.wantHeaders)
			}
		})
	}
}

// stringPtr is defined in generator_test.go
