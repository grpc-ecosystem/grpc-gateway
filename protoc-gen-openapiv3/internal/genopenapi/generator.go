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

// Generate produces one OpenAPI 3.1.0 JSON document per input proto file.
// Files without HTTP-bound services are skipped (no output file is emitted).
func Generate(reg *descriptor.Registry, files []*descriptor.File) ([]*pluginpb.CodeGeneratorResponse_File, error) {
	var out []*pluginpb.CodeGeneratorResponse_File
	for _, file := range files {
		doc, ok, err := generateFile(reg, file)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file.GetName(), err)
		}
		if !ok {
			// No HTTP annotations found, skip file.
			continue
		}
		body, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %w", file.GetName(), err)
		}
		base := strings.TrimSuffix(file.GetName(), path.Ext(file.GetName()))
		out = append(out, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(base + ".openapi.json"),
			Content: proto.String(string(body)),
		})
	}
	return out, nil
}

// generateFile builds a Document for a single proto file. The boolean return
// is false when the file has no HTTP-bound operations to emit.
func generateFile(reg *descriptor.Registry, file *descriptor.File) (*Document, bool, error) {
	name := file.GetName()
	title := strings.TrimSuffix(path.Base(name), path.Ext(name))
	doc := NewDocument(title, "1.0.0")
	if d, ok := fileDocumentAnnotation(file); ok {
		if err := applyDocumentOverride(doc, d); err != nil {
			return nil, false, err
		}
	}
	b := newSchemaBuilder(reg, doc)

	// Prime seenTags with any document-level annotation tags so a service's
	// default tag does not clobber the annotation-provided description.
	seenTags := make(map[string]bool, len(doc.Tags))
	for _, t := range doc.Tags {
		seenTags[t.Name] = true
	}

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
		return nil, false, nil
	}
	if doc.Components.Empty() {
		doc.Components = nil
	}
	return doc, true, nil
}
