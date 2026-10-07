package genopenapi

import (
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
)

func TestConvertSecuritySchemeErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		scheme  *options.SecurityScheme
		wantErr string
	}{
		{
			name:    "missing type",
			scheme:  &options.SecurityScheme{Description: "No type."},
			wantErr: "type is required",
		},
		{
			name: "apiKey missing name",
			scheme: &options.SecurityScheme{
				Type: options.SecurityScheme_TYPE_API_KEY,
				In:   options.SecurityScheme_IN_HEADER,
			},
			wantErr: "name is required",
		},
		{
			name: "apiKey missing in",
			scheme: &options.SecurityScheme{
				Type: options.SecurityScheme_TYPE_API_KEY,
				Name: "X-API-Key",
			},
			wantErr: "in is required",
		},
		{
			name:    "http missing scheme",
			scheme:  &options.SecurityScheme{Type: options.SecurityScheme_TYPE_HTTP},
			wantErr: "scheme is required",
		},
		{
			name:    "oauth2 missing flows",
			scheme:  &options.SecurityScheme{Type: options.SecurityScheme_TYPE_OAUTH2},
			wantErr: "flows is required",
		},
		{
			name: "implicit flow missing authorization_url",
			scheme: &options.SecurityScheme{
				Type: options.SecurityScheme_TYPE_OAUTH2,
				Flows: &options.OAuthFlows{
					Implicit: &options.OAuthFlow{TokenUrl: "https://auth.example.com/token"},
				},
			},
			wantErr: "flows.implicit: authorization_url is required",
		},
		{
			name: "password flow missing token_url",
			scheme: &options.SecurityScheme{
				Type: options.SecurityScheme_TYPE_OAUTH2,
				Flows: &options.OAuthFlows{
					Password: &options.OAuthFlow{},
				},
			},
			wantErr: "flows.password: token_url is required",
		},
		{
			name: "authorization_code flow missing token_url",
			scheme: &options.SecurityScheme{
				Type: options.SecurityScheme_TYPE_OAUTH2,
				Flows: &options.OAuthFlows{
					AuthorizationCode: &options.OAuthFlow{AuthorizationUrl: "https://auth.example.com/authorize"},
				},
			},
			wantErr: "flows.authorization_code: token_url is required",
		},
		{
			name:    "openIdConnect missing url",
			scheme:  &options.SecurityScheme{Type: options.SecurityScheme_TYPE_OPEN_ID_CONNECT},
			wantErr: "open_id_connect_url is required",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := convertSecuritySchemes(map[string]*options.SecurityScheme{"scheme": tc.scheme})
			if err == nil {
				t.Fatalf("want error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error: want substring %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}

func TestApplyDocumentOverrideUndeclaredSecurityScheme(t *testing.T) {
	t.Parallel()

	doc := NewDocument("test", "1.0.0")
	err := applyDocumentOverride(doc, &options.Document{
		Components: &options.Components{
			SecuritySchemes: map[string]*options.SecurityScheme{
				"bearer": {Type: options.SecurityScheme_TYPE_HTTP, Scheme: "bearer"},
			},
		},
		Security: []*options.SecurityRequirement{
			{Schemes: map[string]*options.SecurityRequirement_Scopes{"bearer": {}}},
			{Schemes: map[string]*options.SecurityRequirement_Scopes{"basic": {}}},
		},
	})
	want := `openapiv3 document security[1]: references undeclared security scheme "basic"`
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("error: want substring %q, got %v", want, err)
	}
}
