package openapi31

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
	"go.yaml.in/yaml/v3"
)

// Format represents output format.
type Format int

const (
	FormatJSON Format = iota
	FormatYAML
)

// Adapter converts canonical model to OpenAPI 3.1.0 output.
//
// OpenAPI 3.1.0 specifics:
// - nullable via type arrays: ["string", "null"] (handled by Schema.MarshalJSON/YAML)
// - examples as map, not singular example
// - $ref can have sibling properties (summary, description)
// - JSON Schema draft 2020-12 alignment
type Adapter struct{}

// New creates a new OpenAPI 3.1.0 adapter.
func New() *Adapter {
	return &Adapter{}
}

// Version returns "3.1.0".
func (a *Adapter) Version() string {
	return "3.1.0"
}

// Adapt converts a canonical Document to OpenAPI 3.1.0 JSON or YAML.
// The canonical model handles its own marshaling with proper field ordering.
func (a *Adapter) Adapt(doc *model.Document, format Format) ([]byte, error) {
	if doc == nil {
		return nil, nil
	}

	// Ensure OpenAPI version is set for 3.1.0
	doc.OpenAPIVersion = "3.1.0"

	switch format {
	case FormatYAML:
		return yaml.Marshal(doc)
	default:
		return json.MarshalIndent(doc, "", "  ")
	}
}
