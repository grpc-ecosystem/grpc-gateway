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
	// applyDocumentOverride must run before this point so its tags are
	// included; the post-iteration orphan check below depends on this set
	// being complete.
	seenTags := make(map[string]bool, len(doc.Tags))
	for _, t := range doc.Tags {
		seenTags[t.Name] = true
	}

	// Tracks operationIds across all bindings of all services in this file
	// to enforce the OpenAPI 3.1.0 uniqueness rule. Generator-derived ids
	// (`<Service>_<Method>`) are naturally unique; collisions here mean
	// two annotations (or an annotation colliding with a default) chose
	// the same id.
	seenOpIDs := make(map[string]string)

	// Collected as we go so we can validate against the final tag set
	// after service iteration completes (services may add their default
	// tag after their operations have already declared tags).
	var opTagRefs []opTagRef

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
				op, err := buildOperation(b, svc, method, binding, i, pathParams)
				if err != nil {
					return nil, false, err
				}
				if prev, dup := seenOpIDs[op.OperationID]; dup {
					return nil, false, fmt.Errorf("openapiv3: duplicate operationId %q (used by %s and %s.%s); set a unique openapiv3_operation.operation_id",
						op.OperationID, prev, svc.GetName(), method.GetName())
				}
				seenOpIDs[op.OperationID] = svc.GetName() + "." + method.GetName()
				for _, t := range op.Tags {
					opTagRefs = append(opTagRefs, opTagRef{operationID: op.OperationID, tag: t})
				}

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

	// Validate that every tag referenced by an operation is declared
	// somewhere — either by a service-derived default or by a
	// document-level openapiv3_document.tags annotation. An orphan tag
	// would render as a name with no metadata in the document, which is
	// almost always a mistake (a typo, or a forgotten doc.tags entry).
	for _, ref := range opTagRefs {
		if !seenTags[ref.tag] {
			return nil, false, fmt.Errorf("openapiv3: operation %q references undeclared tag %q; add it to openapiv3_document.tags",
				ref.operationID, ref.tag)
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

// opTagRef remembers a tag reference from an operation, deferred for
// validation after service iteration completes.
type opTagRef struct {
	operationID string
	tag         string
}
