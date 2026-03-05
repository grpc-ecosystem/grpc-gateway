package genopenapiv3

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
)

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
				Schema: &options.Schema{
					Type: "string",
				},
			},
		},
	}

	result := convertResponse(input)

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
		wantRef  bool
		expected *Schema
	}{
		{
			name: "reference schema",
			input: &options.Schema{
				Ref: "#/components/schemas/User",
			},
			wantRef:  true,
			expected: nil,
		},
		{
			name: "inline string schema",
			input: &options.Schema{
				Type:      "string",
				MinLength: 1,
				MaxLength: 100,
				Pattern:   "^[a-z]+$",
			},
			wantRef: false,
			expected: &Schema{
				Type:    "string",
				Pattern: "^[a-z]+$",
			},
		},
		{
			name: "schema with validation",
			input: &options.Schema{
				Type:       "integer",
				Minimum:    0,
				Maximum:    100,
				MultipleOf: 5,
			},
			wantRef: false,
			expected: &Schema{
				Type: "integer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertSchema(tt.input)
			if tt.wantRef {
				if result.Ref != tt.input.Ref {
					t.Errorf("Ref = %q, want %q", result.Ref, tt.input.Ref)
				}
			} else {
				if result.Value == nil {
					t.Fatal("Value should not be nil for inline schema")
				}
				if result.Value.Type != tt.expected.Type {
					t.Errorf("Type = %q, want %q", result.Value.Type, tt.expected.Type)
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
		Schema: &options.Schema{
			Type:   "string",
			Format: "uuid",
		},
	}

	result := convertParameter(input)

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
				Schema: &options.Schema{
					Type: "object",
				},
			},
		},
	}

	result := convertRequestBody(input)

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
		Schema: &options.Schema{
			Type: "string",
		},
	}

	result := convertHeader(input)

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
			OneOf: []*options.Schema{
				{Type: "string"},
				{Type: "integer"},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Value.OneOf) != 2 {
			t.Errorf("OneOf count = %d, want %d", len(result.Value.OneOf), 2)
		}
		if result.Value.OneOf[0].Value.Type != "string" {
			t.Errorf("OneOf[0].Type = %q, want %q", result.Value.OneOf[0].Value.Type, "string")
		}
		if result.Value.OneOf[1].Value.Type != "integer" {
			t.Errorf("OneOf[1].Type = %q, want %q", result.Value.OneOf[1].Value.Type, "integer")
		}
	})

	t.Run("anyOf schema", func(t *testing.T) {
		input := &options.Schema{
			AnyOf: []*options.Schema{
				{Type: "string"},
				{Type: "number"},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Value.AnyOf) != 2 {
			t.Errorf("AnyOf count = %d, want %d", len(result.Value.AnyOf), 2)
		}
	})

	t.Run("allOf schema", func(t *testing.T) {
		input := &options.Schema{
			AllOf: []*options.Schema{
				{Ref: "#/components/schemas/Base"},
				{Type: "object", Required: []string{"extra_field"}},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Value.AllOf) != 2 {
			t.Errorf("AllOf count = %d, want %d", len(result.Value.AllOf), 2)
		}
		if result.Value.AllOf[0].Ref != "#/components/schemas/Base" {
			t.Errorf("AllOf[0].Ref = %q, want %q", result.Value.AllOf[0].Ref, "#/components/schemas/Base")
		}
	})

	t.Run("not schema", func(t *testing.T) {
		input := &options.Schema{
			Type: "string",
			Not:  &options.Schema{Enum: []string{"forbidden"}},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if result.Value.Not == nil {
			t.Fatal("Not should not be nil")
		}
		if len(result.Value.Not.Value.Enum) != 1 {
			t.Errorf("Not.Enum count = %d, want %d", len(result.Value.Not.Value.Enum), 1)
		}
	})

	t.Run("discriminator", func(t *testing.T) {
		input := &options.Schema{
			OneOf: []*options.Schema{
				{Ref: "#/components/schemas/Cat"},
				{Ref: "#/components/schemas/Dog"},
			},
			Discriminator: &options.Discriminator{
				PropertyName: "petType",
				Mapping: map[string]string{
					"cat": "#/components/schemas/Cat",
					"dog": "#/components/schemas/Dog",
				},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if result.Value.Discriminator == nil {
			t.Fatal("Discriminator should not be nil")
		}
		if result.Value.Discriminator.PropertyName != "petType" {
			t.Errorf("Discriminator.PropertyName = %q, want %q", result.Value.Discriminator.PropertyName, "petType")
		}
		if len(result.Value.Discriminator.Mapping) != 2 {
			t.Errorf("Discriminator.Mapping count = %d, want %d", len(result.Value.Discriminator.Mapping), 2)
		}
	})

	t.Run("items for array", func(t *testing.T) {
		input := &options.Schema{
			Type:  "array",
			Items: &options.Schema{Type: "string"},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if result.Value.Items == nil {
			t.Fatal("Items should not be nil")
		}
		if result.Value.Items.Value.Type != "string" {
			t.Errorf("Items.Type = %q, want %q", result.Value.Items.Value.Type, "string")
		}
	})

	t.Run("properties for object", func(t *testing.T) {
		input := &options.Schema{
			Type: "object",
			Properties: map[string]*options.Schema{
				"name": {Type: "string"},
				"age":  {Type: "integer"},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if len(result.Value.Properties) != 2 {
			t.Errorf("Properties count = %d, want %d", len(result.Value.Properties), 2)
		}
		if result.Value.Properties["name"].Value.Type != "string" {
			t.Errorf("Properties[name].Type = %q, want %q", result.Value.Properties["name"].Value.Type, "string")
		}
		if result.Value.Properties["age"].Value.Type != "integer" {
			t.Errorf("Properties[age].Type = %q, want %q", result.Value.Properties["age"].Value.Type, "integer")
		}
	})

	t.Run("additionalProperties with schema", func(t *testing.T) {
		input := &options.Schema{
			Type: "object",
			AdditionalProperties: &options.AdditionalPropertiesItem{
				Kind: &options.AdditionalPropertiesItem_Schema{
					Schema: &options.Schema{Type: "string"},
				},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if result.Value.AdditionalProperties == nil {
			t.Fatal("AdditionalProperties should not be nil")
		}
		if result.Value.AdditionalProperties.Value.Type != "string" {
			t.Errorf("AdditionalProperties.Type = %q, want %q", result.Value.AdditionalProperties.Value.Type, "string")
		}
	})

	t.Run("additionalProperties allows", func(t *testing.T) {
		input := &options.Schema{
			Type: "object",
			AdditionalProperties: &options.AdditionalPropertiesItem{
				Kind: &options.AdditionalPropertiesItem_Allows{Allows: true},
			},
		}

		result := convertSchema(input)

		if result.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if result.Value.AdditionalProperties == nil {
			t.Fatal("AdditionalProperties should not be nil when allows=true")
		}
	})
}
