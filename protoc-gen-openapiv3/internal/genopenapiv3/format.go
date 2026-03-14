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

// Encoder is an interface for encoding OpenAPI documents.
type Encoder interface {
	Encode(v any) error
}

// formatHandler contains the encoding and marshaling functions for a format.
type formatHandler struct {
	newEncoder func(io.Writer) Encoder
	marshal    func(any) ([]byte, error)
}

// formatHandlers maps formats to their handlers.
var formatHandlers = map[Format]formatHandler{
	FormatJSON: {
		newEncoder: func(w io.Writer) Encoder {
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			return enc
		},
		marshal: func(v any) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		},
	},
	FormatYAML: {
		newEncoder: func(w io.Writer) Encoder {
			enc := yaml.NewEncoder(w)
			enc.SetIndent(2)
			return enc
		},
		marshal: yaml.Marshal,
	},
}

// Validate checks if the format is valid.
func (f Format) Validate() error {
	if _, ok := formatHandlers[f]; ok {
		return nil
	}
	return fmt.Errorf("invalid output format %q, must be %q or %q", f, FormatJSON, FormatYAML)
}

// NewEncoder creates an encoder for the given format.
func (f Format) NewEncoder(w io.Writer) (Encoder, error) {
	handler, ok := formatHandlers[f]
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", f)
	}
	return handler.newEncoder(w), nil
}

// Marshal serializes the given value to the format.
func (f Format) Marshal(v any) ([]byte, error) {
	handler, ok := formatHandlers[f]
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", f)
	}
	return handler.marshal(v)
}
