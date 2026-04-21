package genopenapi

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// namedDoc pairs a Document with the output file base name (no extension).
type namedDoc struct {
	name string
	doc  *Document
}

// Generate produces one OpenAPI 3.1.0 JSON document per input proto file,
// or a single merged document when reg.IsAllowMerge() is true.
// Files without HTTP-bound services are skipped (no output file is emitted).
func Generate(reg *descriptor.Registry, files []*descriptor.File) ([]*pluginpb.CodeGeneratorResponse_File, error) {
	var docs []namedDoc
	for _, file := range files {
		doc, ok := generateFile(reg, file)
		if !ok {
			// No HTTP annotations found, skip file.
			continue
		}
		base := strings.TrimSuffix(file.GetName(), path.Ext(file.GetName()))
		docs = append(docs, namedDoc{name: base, doc: doc})
	}

	if reg.IsAllowMerge() {
		merged := mergeDocuments(docs, reg.GetMergeFileName())
		if merged == nil {
			return nil, nil
		}
		body, err := json.MarshalIndent(merged.doc, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal merged: %w", err)
		}
		return []*pluginpb.CodeGeneratorResponse_File{{
			Name:    proto.String(merged.name + ".openapi.json"),
			Content: proto.String(string(body)),
		}}, nil
	}

	var out []*pluginpb.CodeGeneratorResponse_File
	for _, fd := range docs {
		body, err := json.MarshalIndent(fd.doc, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %w", fd.name, err)
		}
		out = append(out, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(fd.name + ".openapi.json"),
			Content: proto.String(string(body)),
		})
	}
	return out, nil
}

// mergeDocuments combines all per-file documents into a single OpenAPI
// document named mergeFileName. The merged document inherits Info from the
// first document. Returns nil when docs is empty.
func mergeDocuments(docs []namedDoc, mergeFileName string) *namedDoc {
	if len(docs) == 0 {
		return nil
	}
	merged := &Document{
		OpenAPI: "3.1.0",
		Info:    docs[0].doc.Info,
		Paths:   NewPaths(),
		Components: &Components{
			Schemas:         make(map[string]*SchemaOrRef),
			SecuritySchemes: make(map[string]*SecurityScheme),
		},
	}
	seenTags := make(map[string]bool)
	for _, fd := range docs {
		// Merge paths in insertion order.
		if fd.doc.Paths != nil {
			for _, urlPath := range fd.doc.Paths.order {
				merged.Paths.Set(urlPath, fd.doc.Paths.items[urlPath])
			}
		}
		// Merge component schemas and security schemes.
		if fd.doc.Components != nil {
			for k, v := range fd.doc.Components.Schemas {
				merged.Components.Schemas[k] = v
			}
			for k, v := range fd.doc.Components.SecuritySchemes {
				merged.Components.SecuritySchemes[k] = v
			}
		}
		// Merge tags, deduplicating by name.
		for _, tag := range fd.doc.Tags {
			if !seenTags[tag.Name] {
				seenTags[tag.Name] = true
				merged.Tags = append(merged.Tags, tag)
			}
		}
		// Merge top-level security requirements and servers.
		merged.Security = append(merged.Security, fd.doc.Security...)
		merged.Servers = append(merged.Servers, fd.doc.Servers...)
	}
	if merged.Components.Empty() {
		merged.Components = nil
	}
	return &namedDoc{name: mergeFileName, doc: merged}
}

// generateFile builds a Document for a single proto file. The boolean return
// is false when the file has no HTTP-bound operations to emit.
func generateFile(reg *descriptor.Registry, file *descriptor.File) (*Document, bool) {
	name := file.GetName()
	title := strings.TrimSuffix(path.Base(name), path.Ext(name))
	doc := NewDocument(title, "1.0.0")
	b := newSchemaBuilder(reg, doc)

	seenTags := make(map[string]bool)

	for _, svc := range file.Services {
		if !isVisible(serviceVisibility(svc), reg) {
			// Service is hidden by visibility rules, skip.
			continue
		}

		tag := svc.GetName()
		// Track whether any methods of this service are visible,
		// to avoid emitting a tag for a service with only hidden methods.
		hasVisibleMethod := false

		for _, method := range svc.Methods {
			if !isVisible(methodVisibility(method), reg) {
				// Method is hidden by visibility rules, skip.
				continue
			}
			for i, binding := range method.Bindings {
				urlPath, pathParams := convertPathTemplate(binding.PathTmpl.Template)
				op := buildOperation(b, svc, method, binding, i, pathParams)

				item, ok := doc.Paths.Get(urlPath)
				if !ok {
					item = &PathItem{}
					doc.Paths.Set(urlPath, item)
				}
				item.SetOperation(binding.HTTPMethod, op)
			}
			// At least one method is visible, so the service's tag
			// should be emitted.
			hasVisibleMethod = true
		}

		// Emit the service's tag if it has any visible
		// methods and the tag hasn't already been emitted.
		if hasVisibleMethod && !seenTags[tag] {
			seenTags[tag] = true
			doc.Tags = append(doc.Tags, &Tag{
				Name:        tag,
				Description: serviceComments(svc),
			})
		}
	}

	if doc.Paths.Len() == 0 {
		return nil, false
	}
	if doc.Components.Empty() {
		doc.Components = nil
	}
	return doc, true
}
