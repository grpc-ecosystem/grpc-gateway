package genopenapiv3

import (
	"encoding/json"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

// Format represents the output format of the OpenAPI specification.
type Format string

const (
	// FormatJSON outputs JSON format.
	FormatJSON Format = "json"
	// FormatYAML outputs YAML format.
	FormatYAML Format = "yaml"
)

// Validate checks if the format is valid.
func (f Format) Validate() error {
	switch f {
	case FormatJSON, FormatYAML:
		return nil
	default:
		return fmt.Errorf("invalid output format %q, must be %q or %q", f, FormatJSON, FormatYAML)
	}
}

// Encoder is an interface for encoding OpenAPI documents.
type Encoder interface {
	Encode(v interface{}) error
}

// NewEncoder creates an encoder for the given format.
func (f Format) NewEncoder(w io.Writer) (Encoder, error) {
	switch f {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc, nil
	case FormatYAML:
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		return enc, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", f)
	}
}

// Marshal serializes the given value to the format.
func (f Format) Marshal(v interface{}) ([]byte, error) {
	switch f {
	case FormatJSON:
		return json.MarshalIndent(v, "", "  ")
	case FormatYAML:
		return yaml.Marshal(v)
	default:
		return nil, fmt.Errorf("unknown format: %s", f)
	}
}
