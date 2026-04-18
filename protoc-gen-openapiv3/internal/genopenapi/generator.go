package genopenapi

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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
		doc, ok := generateFile(reg, file)
		if !ok {
			// No HTTP annotations found, skip file.
			continue
		}
		body, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %w", file.GetName(), err)
		}
		base := strings.TrimSuffix(file.GetName(), filepath.Ext(file.GetName()))
		out = append(out, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(base + ".openapi.json"),
			Content: proto.String(string(body)),
		})
	}
	return out, nil
}

// generateFile builds a Document for a single proto file. The boolean return
// is false when the file has no HTTP-bound operations to emit.
func generateFile(reg *descriptor.Registry, file *descriptor.File) (*Document, bool) {
	name := file.GetName()
	title := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	doc := NewDocument(title, "1.0.0")
	b := newSchemaBuilder(reg, doc)

	seenTags := make(map[string]bool)

	for _, svc := range file.Services {
		tag := svc.GetName()
		if !seenTags[tag] {
			seenTags[tag] = true
			doc.Tags = append(doc.Tags, &Tag{
				Name:        tag,
				Description: serviceComments(svc),
			})
		}

		for _, method := range svc.Methods {
			for i, binding := range method.Bindings {
				path, pathParams := convertPathTemplate(binding.PathTmpl.Template)
				op := buildOperation(b, svc, method, binding, i, pathParams)

				item, ok := doc.Paths.Get(path)
				if !ok {
					item = &PathItem{}
					doc.Paths.Set(path, item)
				}
				item.SetOperation(binding.HTTPMethod, op)
			}
		}
	}

	if doc.Paths.Len() == 0 {
		return nil, false
	}
	return doc, true
}
