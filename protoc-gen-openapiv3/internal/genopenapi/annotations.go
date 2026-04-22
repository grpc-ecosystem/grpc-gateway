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
// invalid (e.g. License sets both `identifier` and `url`, which the OpenAPI
// 3.1.0 spec declares mutually exclusive).
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
			if l.GetIdentifier() != "" && l.GetUrl() != "" {
				return fmt.Errorf("openapiv3 license: identifier and url are mutually exclusive, set only one")
			}
			doc.Info.License = &License{
				Name:       l.GetName(),
				Identifier: l.GetIdentifier(),
				URL:        l.GetUrl(),
			}
		}
	}
	for _, s := range d.GetServers() {
		doc.Servers = append(doc.Servers, &Server{
			URL:         s.GetUrl(),
			Description: s.GetDescription(),
		})
	}
	if ed := d.GetExternalDocs(); ed != nil {
		doc.ExternalDocs = &ExternalDocs{
			Description: ed.GetDescription(),
			URL:         ed.GetUrl(),
		}
	}
	for _, t := range d.GetTags() {
		tag := &Tag{
			Name:        t.GetName(),
			Description: t.GetDescription(),
		}
		if ed := t.GetExternalDocs(); ed != nil {
			tag.ExternalDocs = &ExternalDocs{
				Description: ed.GetDescription(),
				URL:         ed.GetUrl(),
			}
		}
		doc.Tags = append(doc.Tags, tag)
	}
	return nil
}

// applyOperationOverride applies method-level Operation overrides onto the
// generated operation. Annotation values replace comment-derived summary and
// description; a non-empty tag list replaces the default (the service name);
// the deprecated flag is ORed with the proto cascade, so proto-level
// deprecation cannot be revoked by annotation.
func applyOperationOverride(op *Operation, o *options.Operation) {
	if o == nil {
		return
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
		op.ExternalDocs = &ExternalDocs{
			Description: ed.GetDescription(),
			URL:         ed.GetUrl(),
		}
	}
	if o.GetDeprecated() {
		op.Deprecated = true
	}
}

// applySchemaTitle applies the title field of a Schema annotation onto a
// schema body. Used for both message-level and field-level annotations, and
// (for fields) only when an allOf wrapper is present — a bare $ref cannot
// carry a title sibling in OpenAPI 3.1.0.
func applySchemaTitle(s *Schema, o *options.Schema) {
	if o == nil {
		return
	}
	if v := o.GetTitle(); v != "" {
		s.Title = v
	}
}

// applyMessageSchemaOverride applies a message-level Schema annotation onto
// a component schema.
func applyMessageSchemaOverride(s *Schema, o *options.Schema) {
	if o == nil {
		return
	}
	applySchemaTitle(s, o)
	if v := o.GetDescription(); v != "" {
		s.Description = v
	}
}

// annotationNeedsSchemaBody reports whether a field annotation sets anything
// that cannot be expressed as a $ref sibling in OpenAPI 3.1.0. Currently
// that's just `title`: `description` can sit as a sibling directly. Used by
// propertySchema to decide when a referenced field needs an allOf wrapper.
func annotationNeedsSchemaBody(o *options.Schema) bool {
	if o == nil {
		return false
	}
	return o.GetTitle() != ""
}
