package genopenapi

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// Annotation lookups. Each returns (nil, false) when the extension is not
// present. A returned annotation's non-empty sub-fields replace defaults the
// generator would otherwise derive from proto comments or proto types.

func fileDocumentAnnotation(file *descriptor.File) (*options.Document, bool) {
	if file.Options == nil || !proto.HasExtension(file.Options, options.E_Openapiv3Document) {
		return nil, false
	}
	d, ok := proto.GetExtension(file.Options, options.E_Openapiv3Document).(*options.Document)
	if !ok || d == nil {
		return nil, false
	}
	return d, true
}

func serviceTagAnnotation(svc *descriptor.Service) (*options.Tag, bool) {
	if svc.Options == nil || !proto.HasExtension(svc.Options, options.E_Openapiv3Tag) {
		return nil, false
	}
	t, ok := proto.GetExtension(svc.Options, options.E_Openapiv3Tag).(*options.Tag)
	if !ok || t == nil {
		return nil, false
	}
	return t, true
}

func methodOperationAnnotation(m *descriptor.Method) (*options.Operation, bool) {
	if m.Options == nil || !proto.HasExtension(m.Options, options.E_Openapiv3Operation) {
		return nil, false
	}
	op, ok := proto.GetExtension(m.Options, options.E_Openapiv3Operation).(*options.Operation)
	if !ok || op == nil {
		return nil, false
	}
	return op, true
}

func messageSchemaAnnotation(msg *descriptor.Message) (*options.Schema, bool) {
	if msg.Options == nil || !proto.HasExtension(msg.Options, options.E_Openapiv3Schema) {
		return nil, false
	}
	s, ok := proto.GetExtension(msg.Options, options.E_Openapiv3Schema).(*options.Schema)
	if !ok || s == nil {
		return nil, false
	}
	return s, true
}

func fieldSchemaAnnotation(field *descriptor.Field) (*options.Schema, bool) {
	if field.Options == nil || !proto.HasExtension(field.Options, options.E_Openapiv3Field) {
		return nil, false
	}
	s, ok := proto.GetExtension(field.Options, options.E_Openapiv3Field).(*options.Schema)
	if !ok || s == nil {
		return nil, false
	}
	return s, true
}

// processExtensions converts a proto options extensions map (as attached to
// openapiv3_document, openapiv3_tag, openapiv3_operation, openapiv3_schema,
// and openapiv3_field annotations) into the sorted, JSON-ready form the
// generator's internal types render inline.
//
// Per the OpenAPI 3.1.0 spec
// (https://spec.openapis.org/oas/v3.1.0#specification-extensions), extension
// keys must start with "x-"; an entry that doesn't is dropped with a log
// line rather than failing generation, consistent with how this generator
// degrades other malformed annotation input (e.g. an unresolvable field
// type) elsewhere. where identifies the annotation site in that log line,
// e.g. "openapiv3_operation".
func processExtensions(where string, exts map[string]*structpb.Value) []extension {
	if len(exts) == 0 {
		return nil
	}
	out := make([]extension, 0, len(exts))
	for k, v := range exts {
		if !strings.HasPrefix(k, "x-") {
			grpclog.Infof("protoc-gen-openapiv3: %s extension key %q does not start with \"x-\"; skipping", where, k)
			continue
		}
		data, err := protojson.Marshal(v)
		if err != nil {
			grpclog.Infof("protoc-gen-openapiv3: %s extension %q: %v; skipping", where, k, err)
			continue
		}
		out = append(out, extension{key: k, value: data})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].key < out[j].key })
	return out
}

// applyDocumentOverride applies file-level Document overrides onto the
// generated OpenAPI document. Non-empty fields replace defaults; empty fields
// leave the current value untouched. Returns an error if the annotation is
// invalid — for example, a License with no name, a Server with no url, a Tag
// with no name, or any ExternalDocs without a url. All four are spec-required
// fields per OpenAPI 3.1.0. Security schemes are validated the same way, and
// document-level security requirements must name a declared scheme.
func applyDocumentOverride(doc *Document, d *options.Document) error {
	if d == nil {
		return nil
	}
	if info := d.GetInfo(); info != nil {
		if v := info.GetTitle(); v != "" {
			doc.Info.Title = v
		}
		if v := info.GetSummary(); v != "" {
			doc.Info.Summary = v
		}
		if v := info.GetDescription(); v != "" {
			doc.Info.Description = v
		}
		if v := info.GetTermsOfService(); v != "" {
			doc.Info.TermsOfService = v
		}
		if v := info.GetVersion(); v != "" {
			doc.Info.Version = v
		}
		if c := info.GetContact(); c != nil {
			doc.Info.Contact = &Contact{
				Name:  c.GetName(),
				URL:   c.GetUrl(),
				Email: c.GetEmail(),
			}
		}
		if l := info.GetLicense(); l != nil {
			if l.GetName() == "" {
				return fmt.Errorf("openapiv3 license: name is required")
			}
			doc.Info.License = &License{
				Name:       l.GetName(),
				Identifier: l.GetIdentifier(),
				URL:        l.GetUrl(),
			}
		}
		doc.Info.Extensions = processExtensions("openapiv3_document.info", info.GetExtensions())
	}
	for i, s := range d.GetServers() {
		if err := validateServer(s); err != nil {
			return fmt.Errorf("openapiv3 document servers[%d]: %w", i, err)
		}
		doc.Servers = append(doc.Servers, &Server{
			URL:         s.GetUrl(),
			Description: s.GetDescription(),
		})
	}
	if ed := d.GetExternalDocs(); ed != nil {
		if err := validateExternalDocs(ed); err != nil {
			return fmt.Errorf("openapiv3 document external_docs: %w", err)
		}
		doc.ExternalDocs = &ExternalDocs{
			Description: ed.GetDescription(),
			URL:         ed.GetUrl(),
		}
	}
	for i, t := range d.GetTags() {
		if t.GetName() == "" {
			return fmt.Errorf("openapiv3 document tags[%d]: name is required", i)
		}
		tag := &Tag{
			Name:        t.GetName(),
			Description: t.GetDescription(),
		}
		if ed := t.GetExternalDocs(); ed != nil {
			if err := validateExternalDocs(ed); err != nil {
				return fmt.Errorf("openapiv3 document tags[%d] external_docs: %w", i, err)
			}
			tag.ExternalDocs = &ExternalDocs{
				Description: ed.GetDescription(),
				URL:         ed.GetUrl(),
			}
		}
		tag.Extensions = processExtensions(fmt.Sprintf("openapiv3_document.tags[%d]", i), t.GetExtensions())
		doc.Tags = append(doc.Tags, tag)
	}
	if c := d.GetComponents(); c != nil {
		schemes, err := convertSecuritySchemes(c.GetSecuritySchemes())
		if err != nil {
			return fmt.Errorf("openapiv3 document components: %w", err)
		}
		doc.Components.SecuritySchemes = schemes
	}
	doc.Security = convertSecurityRequirements(d.GetSecurity())
	if err := checkSecurityRequirements(doc.Security, doc.Components.SecuritySchemes); err != nil {
		return fmt.Errorf("openapiv3 document %w", err)
	}
	doc.Extensions = processExtensions("openapiv3_document", d.GetExtensions())
	return nil
}

// convertSecuritySchemes converts the security schemes of a document-level
// Components annotation. Returns an error if a scheme is missing a field the
// OpenAPI 3.1.0 spec requires for its type.
func convertSecuritySchemes(schemes map[string]*options.SecurityScheme) (map[string]*SecurityScheme, error) {
	if len(schemes) == 0 {
		return nil, nil
	}
	out := make(map[string]*SecurityScheme, len(schemes))
	// Iterate in sorted order so that error messages are deterministic.
	for _, name := range slices.Sorted(maps.Keys(schemes)) {
		s, err := convertSecurityScheme(schemes[name])
		if err != nil {
			return nil, fmt.Errorf("security_schemes[%q]: %w", name, err)
		}
		out[name] = s
	}
	return out, nil
}

// convertSecurityScheme converts a single security scheme. The type and in
// values and the fields required for each type are the ones listed in the
// Security Scheme Object's fixed fields.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#security-scheme-object
func convertSecurityScheme(s *options.SecurityScheme) (*SecurityScheme, error) {
	out := &SecurityScheme{Description: s.GetDescription()}
	// See: https://spec.openapis.org/oas/v3.1.0#securitySchemeType
	switch s.GetType() {
	case options.SecurityScheme_TYPE_API_KEY:
		out.Type = "apiKey"
		if s.GetName() == "" {
			return nil, fmt.Errorf("name is required for apiKey security schemes")
		}
		out.Name = s.GetName()
		// See: https://spec.openapis.org/oas/v3.1.0#securitySchemeIn
		switch s.GetIn() {
		case options.SecurityScheme_IN_QUERY:
			out.In = "query"
		case options.SecurityScheme_IN_HEADER:
			out.In = "header"
		case options.SecurityScheme_IN_COOKIE:
			out.In = "cookie"
		default:
			return nil, fmt.Errorf("in is required for apiKey security schemes")
		}
	case options.SecurityScheme_TYPE_HTTP:
		out.Type = "http"
		if s.GetScheme() == "" {
			return nil, fmt.Errorf("scheme is required for http security schemes")
		}
		out.Scheme = s.GetScheme()
		out.BearerFormat = s.GetBearerFormat()
	case options.SecurityScheme_TYPE_MUTUAL_TLS:
		out.Type = "mutualTLS"
	case options.SecurityScheme_TYPE_OAUTH2:
		out.Type = "oauth2"
		flows, err := convertOAuthFlows(s.GetFlows())
		if err != nil {
			return nil, err
		}
		out.Flows = flows
	case options.SecurityScheme_TYPE_OPEN_ID_CONNECT:
		out.Type = "openIdConnect"
		if s.GetOpenIdConnectUrl() == "" {
			return nil, fmt.Errorf("open_id_connect_url is required for openIdConnect security schemes")
		}
		out.OpenIDConnectURL = s.GetOpenIdConnectUrl()
	default:
		return nil, fmt.Errorf("type is required")
	}
	return out, nil
}

// convertOAuthFlows converts the flows of an oauth2 security scheme. Which
// URLs are required depends on the flow, as listed in the OAuth Flow Object's
// fixed fields.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
func convertOAuthFlows(f *options.OAuthFlows) (*OAuthFlows, error) {
	if f == nil {
		return nil, fmt.Errorf("flows is required for oauth2 security schemes")
	}
	var (
		out OAuthFlows
		err error
	)
	if out.Implicit, err = convertOAuthFlow("implicit", f.GetImplicit(), true, false); err != nil {
		return nil, err
	}
	if out.Password, err = convertOAuthFlow("password", f.GetPassword(), false, true); err != nil {
		return nil, err
	}
	if out.ClientCredentials, err = convertOAuthFlow("client_credentials", f.GetClientCredentials(), false, true); err != nil {
		return nil, err
	}
	if out.AuthorizationCode, err = convertOAuthFlow("authorization_code", f.GetAuthorizationCode(), true, true); err != nil {
		return nil, err
	}
	return &out, nil
}

// convertOAuthFlow converts a single OAuth flow, enforcing the URLs the
// OpenAPI 3.1.0 spec requires for that flow type. A nil flow is not
// supported by the scheme and converts to nil.
//
// Spec: https://spec.openapis.org/oas/v3.1.0#oauth-flow-object
func convertOAuthFlow(name string, f *options.OAuthFlow, needAuthorizationURL, needTokenURL bool) (*OAuthFlow, error) {
	if f == nil {
		return nil, nil
	}
	if needAuthorizationURL && f.GetAuthorizationUrl() == "" {
		return nil, fmt.Errorf("flows.%s: authorization_url is required", name)
	}
	if needTokenURL && f.GetTokenUrl() == "" {
		return nil, fmt.Errorf("flows.%s: token_url is required", name)
	}
	out := &OAuthFlow{
		AuthorizationURL: f.GetAuthorizationUrl(),
		TokenURL:         f.GetTokenUrl(),
		RefreshURL:       f.GetRefreshUrl(),
		// scopes is required by the spec, even when empty.
		Scopes: make(map[string]string, len(f.GetScopes())),
	}
	maps.Copy(out.Scopes, f.GetScopes())
	return out, nil
}

// convertSecurityRequirements converts security requirements from a
// Document or Operation annotation. An empty requirement is kept, since it
// renders as {} and marks security as optional.
func convertSecurityRequirements(reqs []*options.SecurityRequirement) []SecurityRequirement {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]SecurityRequirement, 0, len(reqs))
	for _, r := range reqs {
		req := make(SecurityRequirement, len(r.GetSchemes()))
		for name, scopes := range r.GetSchemes() {
			// The scopes array is required, even when empty.
			req[name] = append([]string{}, scopes.GetScopes()...)
		}
		out = append(out, req)
	}
	return out
}

// checkSecurityRequirements verifies that every scheme named by reqs is
// declared in schemes, as required by the OpenAPI 3.1.0 spec.
func checkSecurityRequirements(reqs []SecurityRequirement, schemes map[string]*SecurityScheme) error {
	for i, req := range reqs {
		for _, name := range slices.Sorted(maps.Keys(req)) {
			if _, ok := schemes[name]; !ok {
				return fmt.Errorf("security[%d]: references undeclared security scheme %q; add it to openapiv3_document.components.security_schemes", i, name)
			}
		}
	}
	return nil
}

// validateServer enforces the OpenAPI 3.1.0 Server Object's required `url`
// field. Empty url is invalid even though `description` alone may look
// useful in proto.
func validateServer(s *options.Server) error {
	if s.GetUrl() == "" {
		return fmt.Errorf("server: url is required")
	}
	return nil
}

// validateExternalDocs enforces the OpenAPI 3.1.0 External Documentation
// Object's required `url` field.
func validateExternalDocs(ed *options.ExternalDocs) error {
	if ed.GetUrl() == "" {
		return fmt.Errorf("external_docs: url is required")
	}
	return nil
}

// applyTagOverride applies a service-level Tag annotation onto the tag
// generated for the service. Non-empty fields replace the defaults (the
// service name and its leading comment). Returns an error if the annotation
// contains an external_docs without a url, which is spec-required per
// OpenAPI 3.1.0.
func applyTagOverride(tag *Tag, t *options.Tag) error {
	if t == nil {
		return nil
	}
	if v := t.GetName(); v != "" {
		tag.Name = v
	}
	if v := t.GetDescription(); v != "" {
		tag.Description = v
	}
	if ed := t.GetExternalDocs(); ed != nil {
		if err := validateExternalDocs(ed); err != nil {
			return err
		}
		tag.ExternalDocs = &ExternalDocs{
			Description: ed.GetDescription(),
			URL:         ed.GetUrl(),
		}
	}
	tag.Extensions = processExtensions("openapiv3_tag", t.GetExtensions())
	return nil
}

// applyOperationOverride applies method-level Operation overrides onto the
// generated operation. Annotation values replace comment-derived summary and
// description; a non-empty tag list replaces the default (the service name);
// servers from the annotation are appended to any defaults; security
// requirements replace the document-level ones. The annotation
// `deprecated` flag is one-way: it can flip deprecation on, but cannot
// clear a flag inherited from the proto cascade.
//
// Returns an error if the annotation contains an external_docs without a url
// or a server without a url, both spec-required per OpenAPI 3.1.0.
func applyOperationOverride(op *Operation, o *options.Operation) error {
	if o == nil {
		return nil
	}
	if tags := o.GetTags(); len(tags) > 0 {
		op.Tags = tags
	}
	if v := o.GetSummary(); v != "" {
		op.Summary = v
	}
	if v := o.GetDescription(); v != "" {
		op.Description = v
	}
	if v := o.GetOperationId(); v != "" {
		op.OperationID = v
	}
	if ed := o.GetExternalDocs(); ed != nil {
		if err := validateExternalDocs(ed); err != nil {
			return fmt.Errorf("external_docs: %w", err)
		}
		op.ExternalDocs = &ExternalDocs{
			Description: ed.GetDescription(),
			URL:         ed.GetUrl(),
		}
	}
	op.Deprecated = op.Deprecated || o.GetDeprecated()
	for i, s := range o.GetServers() {
		if err := validateServer(s); err != nil {
			return fmt.Errorf("servers[%d]: %w", i, err)
		}
		op.Servers = append(op.Servers, &Server{
			URL:         s.GetUrl(),
			Description: s.GetDescription(),
		})
	}
	op.Security = convertSecurityRequirements(o.GetSecurity())
	op.Extensions = processExtensions("openapiv3_operation", o.GetExtensions())
	return nil
}

// applySchemaBodyOverride applies the title and deprecated fields from a
// Schema annotation onto a schema body. Used for both message-level and
// field-level annotations. For $ref-typed fields the caller must first
// ensure an allOf wrapper exists, since neither field can sit alongside
// $ref without one. The annotation `deprecated` flag is one-way: it can
// flip deprecation on, but cannot clear a flag inherited from the proto
// cascade.
func applySchemaBodyOverride(s *Schema, o *options.Schema) {
	if o == nil {
		return
	}
	if v := o.GetTitle(); v != "" {
		s.Title = v
	}
	s.Deprecated = s.Deprecated || o.GetDeprecated()
	if len(o.GetExtensions()) > 0 {
		// Shared by openapiv3_schema (message-level) and openapiv3_field
		// (field-level) annotations; the generic label covers both.
		s.Extensions = processExtensions("openapiv3_schema/openapiv3_field", o.GetExtensions())
	}
}

// applyMessageSchemaOverride applies a message-level Schema annotation onto
// a component schema.
func applyMessageSchemaOverride(s *Schema, o *options.Schema) {
	if o == nil {
		return
	}
	applySchemaBodyOverride(s, o)
	if v := o.GetDescription(); v != "" {
		s.Description = v
	}
}

// annotationNeedsSchemaBody reports whether a field annotation sets anything
// that cannot be expressed as a $ref sibling in OpenAPI 3.1.0. `title`,
// `deprecated`, and `extensions` all require a real schema body;
// `description` can sit as a $ref sibling directly. Used by propertySchema
// to decide when a referenced field needs an allOf wrapper.
func annotationNeedsSchemaBody(o *options.Schema) bool {
	if o == nil {
		return false
	}
	return o.GetTitle() != "" || o.GetDeprecated() || len(o.GetExtensions()) > 0
}
