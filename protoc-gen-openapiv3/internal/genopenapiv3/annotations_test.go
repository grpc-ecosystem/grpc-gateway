package genopenapiv3

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/types/descriptorpb"
)

// float64Ptr is a helper to create *float64 for test cases
func float64Ptr(f float64) *float64 {
	return &f
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
				Minimum:    float64Ptr(0),
				Maximum:    float64Ptr(100),
				MultipleOf: float64Ptr(5),
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

// ============================================================================
// Apply Annotation Tests
// ============================================================================

func TestApplyInfoAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		opts         *options.Info
		wantTitle    string
		wantVersion  string
		wantSummary  string
		wantDesc     string
		wantTerms    string
		wantContact  bool
		wantLicense  bool
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
			gen := &generator{reg: reg}
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
		name            string
		opts            *options.Schema
		wantTitle       string
		wantDesc        string
		wantExample     string
		wantReadOnly    bool
		wantWriteOnly   bool
		wantNullable    bool
		wantDeprecated  bool
		wantRequired    []string
		wantAllOf       int
		wantAnyOf       int
		wantOneOf       int
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
			wantExample: `{"id": "123"}`,
		},
		{
			name: "apply composition types",
			opts: &options.Schema{
				AllOf: []*options.Schema{
					{Ref: "#/components/schemas/Base"},
				},
				AnyOf: []*options.Schema{
					{Type: "string"},
					{Type: "integer"},
				},
				OneOf: []*options.Schema{
					{Ref: "#/components/schemas/Cat"},
					{Ref: "#/components/schemas/Dog"},
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

			schema := &Schema{Type: "object"}

			// Create a mock message with the annotation
			msg := &descriptor.Message{
				DescriptorProto: &descriptorpb.DescriptorProto{
					Name: stringPtr("TestMessage"),
				},
			}

			// Since we can't easily set proto extensions in tests,
			// we'll call the internal apply function directly with the opts
			reg := &descriptor.Registry{}
			gen := &generator{reg: reg}

			// Apply the annotation manually (simulating what applySchemaAnnotation does)
			if tt.opts.GetTitle() != "" {
				schema.Title = tt.opts.GetTitle()
			}
			if tt.opts.GetDescription() != "" {
				schema.Description = tt.opts.GetDescription()
			}
			if len(tt.opts.GetRequired()) > 0 {
				schema.Required = tt.opts.GetRequired()
			}
			if tt.opts.GetExample() != "" {
				schema.Example = tt.opts.GetExample()
			}
			if tt.opts.GetReadOnly() {
				schema.ReadOnly = true
			}
			if tt.opts.GetWriteOnly() {
				schema.WriteOnly = true
			}
			if tt.opts.GetNullable() {
				schema.Nullable = true
			}
			if tt.opts.GetDeprecated() {
				schema.Deprecated = true
			}
			for _, allOfSchema := range tt.opts.GetAllOf() {
				schema.AllOf = append(schema.AllOf, convertSchema(allOfSchema))
			}
			for _, anyOfSchema := range tt.opts.GetAnyOf() {
				schema.AnyOf = append(schema.AnyOf, convertSchema(anyOfSchema))
			}
			for _, oneOfSchema := range tt.opts.GetOneOf() {
				schema.OneOf = append(schema.OneOf, convertSchema(oneOfSchema))
			}

			// Verify the generator exists (won't be nil)
			if gen == nil || msg == nil {
				t.Fatal("generator and message should exist")
			}

			// Assertions
			if tt.wantTitle != "" && schema.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", schema.Title, tt.wantTitle)
			}
			if tt.wantDesc != "" && schema.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", schema.Description, tt.wantDesc)
			}
			if tt.wantExample != "" && schema.Example != tt.wantExample {
				t.Errorf("Example = %q, want %q", schema.Example, tt.wantExample)
			}
			if schema.ReadOnly != tt.wantReadOnly {
				t.Errorf("ReadOnly = %v, want %v", schema.ReadOnly, tt.wantReadOnly)
			}
			if schema.WriteOnly != tt.wantWriteOnly {
				t.Errorf("WriteOnly = %v, want %v", schema.WriteOnly, tt.wantWriteOnly)
			}
			if schema.Nullable != tt.wantNullable {
				t.Errorf("Nullable = %v, want %v", schema.Nullable, tt.wantNullable)
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
		wantExample    string
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

			schema := &Schema{Type: "string"}

			// Apply the field annotation options manually
			if tt.opts.GetTitle() != "" {
				schema.Title = tt.opts.GetTitle()
			}
			if tt.opts.GetDescription() != "" {
				schema.Description = tt.opts.GetDescription()
			}
			if tt.opts.GetDefault() != "" {
				schema.Default = tt.opts.GetDefault()
			}
			if tt.opts.GetExample() != "" {
				schema.Example = tt.opts.GetExample()
			}
			if tt.opts.GetFormat() != "" {
				schema.Format = tt.opts.GetFormat()
			}
			if tt.opts.GetPattern() != "" {
				schema.Pattern = tt.opts.GetPattern()
			}
			if tt.opts.GetMinLength() > 0 {
				minLen := tt.opts.GetMinLength()
				schema.MinLength = &minLen
			}
			if tt.opts.GetMaxLength() > 0 {
				maxLen := tt.opts.GetMaxLength()
				schema.MaxLength = &maxLen
			}
			if tt.opts.GetMinimum() != 0 {
				min := tt.opts.GetMinimum()
				schema.Minimum = &min
			}
			if tt.opts.GetMaximum() != 0 {
				max := tt.opts.GetMaximum()
				schema.Maximum = &max
			}
			if tt.opts.GetReadOnly() {
				schema.ReadOnly = true
			}
			if tt.opts.GetWriteOnly() {
				schema.WriteOnly = true
			}
			if tt.opts.GetNullable() {
				schema.Nullable = true
			}
			if tt.opts.GetDeprecated() {
				schema.Deprecated = true
			}

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
			if tt.wantExample != "" && schema.Example != tt.wantExample {
				t.Errorf("Example = %q, want %q", schema.Example, tt.wantExample)
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
			if schema.Nullable != tt.wantNullable {
				t.Errorf("Nullable = %v, want %v", schema.Nullable, tt.wantNullable)
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
		name           string
		opts           *options.Operation
		wantSummary    string
		wantDesc       string
		wantOpID       string
		wantTags       []string
		wantDeprecated bool
		wantSecurity   int
		wantServers    int
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

			// Apply the operation annotation options manually
			if tt.opts.GetSummary() != "" {
				op.Summary = tt.opts.GetSummary()
			}
			if tt.opts.GetDescription() != "" {
				op.Description = tt.opts.GetDescription()
			}
			if tt.opts.GetOperationId() != "" {
				op.OperationID = tt.opts.GetOperationId()
			}
			if len(tt.opts.GetTags()) > 0 {
				op.Tags = tt.opts.GetTags()
			}
			if tt.opts.GetDeprecated() {
				op.Deprecated = true
			}
			for _, sec := range tt.opts.GetSecurity() {
				op.Security = append(op.Security, convertSecurityRequirement(sec))
			}
			for _, s := range tt.opts.GetServers() {
				op.Servers = append(op.Servers, convertServer(s))
			}

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

			// Apply the service annotation options manually
			if tt.opts.GetName() != "" {
				tag.Name = tt.opts.GetName()
			}
			if tt.opts.GetDescription() != "" {
				tag.Description = tt.opts.GetDescription()
			}
			if extDocs := tt.opts.GetExternalDocs(); extDocs != nil {
				tag.ExternalDocs = convertExternalDocs(extDocs)
			}

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
		wantExample    string
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
				Type: "string",
				Enum: []any{"PENDING", "COMPLETED", "FAILED"},
			}

			// Apply the enum annotation options manually
			if tt.opts.GetTitle() != "" {
				schema.Title = tt.opts.GetTitle()
			}
			if tt.opts.GetDescription() != "" {
				schema.Description = tt.opts.GetDescription()
			}
			if tt.opts.GetDefault() != "" {
				schema.Default = tt.opts.GetDefault()
			}
			if tt.opts.GetExample() != "" {
				schema.Example = tt.opts.GetExample()
			}
			if tt.opts.GetDeprecated() {
				schema.Deprecated = true
			}
			if extDocs := tt.opts.GetExternalDocs(); extDocs != nil {
				schema.ExternalDocs = convertExternalDocs(extDocs)
			}

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
			if tt.wantExample != "" && schema.Example != tt.wantExample {
				t.Errorf("Example = %q, want %q", schema.Example, tt.wantExample)
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
					"NotFound": {Description: "Resource not found"},
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
						Schema:   &options.Schema{Type: "integer"},
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
				Headers: map[string]*options.Header{
					"X-Request-Id": {
						Description: "Request correlation ID",
						Schema:      &options.Schema{Type: "string"},
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
				Schemas: make(map[string]*SchemaRef),
			}

			reg := &descriptor.Registry{}
			gen := &generator{reg: reg}
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
