package genopenapi

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/proto"
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

// applyDocumentOverride applies file-level Document overrides onto the
// generated OpenAPI document. Non-empty fields replace defaults; empty fields
// leave the current value untouched. Returns an error if the annotation is
// invalid — for example, a License with no name, a Server with no url, a Tag
// with no name, or any ExternalDocs without a url. All four are spec-required
// fields per OpenAPI 3.1.0.
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
		doc.Tags = append(doc.Tags, tag)
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

// applyOperationOverride applies method-level Operation overrides onto the
// generated operation. Annotation values replace comment-derived summary and
// description; a non-empty tag list replaces the default (the service name);
// servers from the annotation are appended to any defaults. The annotation
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
// that cannot be expressed as a $ref sibling in OpenAPI 3.1.0. `title` and
// `deprecated` both require a real schema body; `description` can sit as a
// $ref sibling directly. Used by propertySchema to decide when a referenced
// field needs an allOf wrapper.
func annotationNeedsSchemaBody(o *options.Schema) bool {
	if o == nil {
		return false
	}
	return o.GetTitle() != "" || o.GetDeprecated()
}
