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
// OpenAPI 3.1.0 specifics handled by transformation:
// - nullable via type arrays: ["string", "null"]
// - examples as array (JSON Schema style)
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
// It transforms the canonical model to 3.1.0-specific output types, then serializes.
func (a *Adapter) Adapt(doc *model.Document, format Format) ([]byte, error) {
	if doc == nil {
		return nil, nil
	}

	// Transform canonical model to 3.1.0 output format
	out := TransformDocument(doc)

	switch format {
	case FormatYAML:
		return yaml.Marshal(out)
	default:
		return json.MarshalIndent(out, "", "  ")
	}
}
