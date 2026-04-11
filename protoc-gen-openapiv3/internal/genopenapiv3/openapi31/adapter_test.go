package openapi31

import (
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
)

func TestDocumentFieldOrdering_JSON(t *testing.T) {
	adapter := New()
	doc := model.NewDocument("3.1.0")
	doc.Info.Title = "Test API"
	doc.Info.Version = "1.0.0"
	doc.Components.Schemas["TestSchema"] = model.NewInlineSchema(&model.Schema{Type: "string"})
	doc.Tags = []*model.Tag{{Name: "test"}}

	jsonOutput, err := adapter.Adapt(doc, FormatJSON)
	if err != nil {
		t.Fatalf("Adapt failed: %v", err)
	}

	output := string(jsonOutput)

	// Verify field order: openapi before info before components before tags
	openapiIdx := strings.Index(output, `"openapi"`)
	infoIdx := strings.Index(output, `"info"`)
	componentsIdx := strings.Index(output, `"components"`)
	tagsIdx := strings.Index(output, `"tags"`)

	if openapiIdx == -1 {
		t.Error("openapi field should exist in output")
	}
	if infoIdx == -1 {
		t.Error("info field should exist in output")
	}
	if componentsIdx == -1 {
		t.Error("components field should exist in output")
	}
	if tagsIdx == -1 {
		t.Error("tags field should exist in output")
	}

	if openapiIdx >= infoIdx {
		t.Errorf("openapi (at %d) should appear before info (at %d)", openapiIdx, infoIdx)
	}
	if infoIdx >= componentsIdx {
		t.Errorf("info (at %d) should appear before components (at %d)", infoIdx, componentsIdx)
	}
	if componentsIdx >= tagsIdx {
		t.Errorf("components (at %d) should appear before tags (at %d)", componentsIdx, tagsIdx)
	}
}

func TestDocumentFieldOrdering_YAML(t *testing.T) {
	adapter := New()
	doc := model.NewDocument("3.1.0")
	doc.Info.Title = "Test API"
	doc.Info.Version = "1.0.0"
	doc.Components.Schemas["TestSchema"] = model.NewInlineSchema(&model.Schema{Type: "string"})
	doc.Tags = []*model.Tag{{Name: "test"}}

	yamlOutput, err := adapter.Adapt(doc, FormatYAML)
	if err != nil {
		t.Fatalf("Adapt failed: %v", err)
	}

	output := string(yamlOutput)

	// Verify field order in YAML
	openapiIdx := strings.Index(output, "openapi:")
	infoIdx := strings.Index(output, "info:")
	componentsIdx := strings.Index(output, "components:")
	tagsIdx := strings.Index(output, "tags:")

	if openapiIdx == -1 {
		t.Error("openapi field should exist in output")
	}
	if infoIdx == -1 {
		t.Error("info field should exist in output")
	}
	if componentsIdx == -1 {
		t.Error("components field should exist in output")
	}
	if tagsIdx == -1 {
		t.Error("tags field should exist in output")
	}

	if openapiIdx >= infoIdx {
		t.Errorf("openapi (at %d) should appear before info (at %d)", openapiIdx, infoIdx)
	}
	if infoIdx >= componentsIdx {
		t.Errorf("info (at %d) should appear before components (at %d)", infoIdx, componentsIdx)
	}
	if componentsIdx >= tagsIdx {
		t.Errorf("components (at %d) should appear before tags (at %d)", componentsIdx, tagsIdx)
	}
}

func TestOutputDeterminism(t *testing.T) {
	adapter := New()
	doc := createComplexDocument()

	// Run multiple times - output should be identical
	var jsonOutputs []string
	var yamlOutputs []string
	for i := 0; i < 10; i++ {
		jsonOutput, err := adapter.Adapt(doc, FormatJSON)
		if err != nil {
			t.Fatalf("Adapt JSON failed: %v", err)
		}
		jsonOutputs = append(jsonOutputs, string(jsonOutput))

		yamlOutput, err := adapter.Adapt(doc, FormatYAML)
		if err != nil {
			t.Fatalf("Adapt YAML failed: %v", err)
		}
		yamlOutputs = append(yamlOutputs, string(yamlOutput))
	}

	for i := 1; i < len(jsonOutputs); i++ {
		if jsonOutputs[0] != jsonOutputs[i] {
			t.Errorf("JSON output should be deterministic across runs, run 0 vs run %d differ", i)
		}
	}

	for i := 1; i < len(yamlOutputs); i++ {
		if yamlOutputs[0] != yamlOutputs[i] {
			t.Errorf("YAML output should be deterministic across runs, run 0 vs run %d differ", i)
		}
	}
}

func createComplexDocument() *model.Document {
	doc := model.NewDocument("3.1.0")
	doc.Info.Title = "Complex API"
	doc.Info.Version = "2.0.0"
	doc.Info.Description = "A complex API for testing"

	// Add multiple schemas
	doc.Components.Schemas["User"] = model.NewInlineSchema(&model.Schema{
		Type: "object",
		Properties: map[string]*model.SchemaOrRef{
			"id":   model.NewInlineSchema(&model.Schema{Type: "integer"}),
			"name": model.NewInlineSchema(&model.Schema{Type: "string"}),
		},
	})
	doc.Components.Schemas["Error"] = model.NewInlineSchema(&model.Schema{
		Type: "object",
		Properties: map[string]*model.SchemaOrRef{
			"code":    model.NewInlineSchema(&model.Schema{Type: "integer"}),
			"message": model.NewInlineSchema(&model.Schema{Type: "string"}),
		},
	})

	// Add tags
	doc.Tags = []*model.Tag{
		{Name: "users", Description: "User operations"},
		{Name: "admin", Description: "Admin operations"},
	}

	// Add servers
	doc.Servers = []*model.Server{
		{URL: "https://api.example.com", Description: "Production"},
		{URL: "https://staging.example.com", Description: "Staging"},
	}

	return doc
}
