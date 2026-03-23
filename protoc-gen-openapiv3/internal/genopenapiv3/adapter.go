package genopenapiv3

import (
	"fmt"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/v31"
)

// Adapter converts a canonical Document to version-specific output.
// Each OpenAPI version has its own adapter implementation.
type Adapter interface {
	// Version returns the OpenAPI version this adapter produces (e.g., "3.1.0").
	Version() string

	// Adapt converts a canonical Document to JSON/YAML bytes.
	Adapt(doc *model.Document, format Format) ([]byte, error)
}

// v31AdapterWrapper wraps v31.Adapter to implement the Adapter interface.
type v31AdapterWrapper struct {
	adapter *v31.Adapter
}

func (w *v31AdapterWrapper) Version() string {
	return w.adapter.Version()
}

func (w *v31AdapterWrapper) Adapt(doc *model.Document, format Format) ([]byte, error) {
	// Convert Format to v31.Format
	var v31Format v31.Format
	if format == FormatYAML {
		v31Format = v31.FormatYAML
	} else {
		v31Format = v31.FormatJSON
	}
	return w.adapter.Adapt(doc, v31Format)
}

// AdapterRegistry holds available adapters keyed by version.
type AdapterRegistry struct {
	adapters map[string]Adapter
}

// NewAdapterRegistry creates a registry with default adapters.
func NewAdapterRegistry() *AdapterRegistry {
	r := &AdapterRegistry{
		adapters: make(map[string]Adapter),
	}
	// Register built-in adapters
	r.Register(&v31AdapterWrapper{adapter: v31.New()})
	// r.Register(&v30AdapterWrapper{adapter: v30.New()}) // Future: 3.0.x support
	return r
}

// Register adds an adapter to the registry.
func (r *AdapterRegistry) Register(a Adapter) {
	r.adapters[a.Version()] = a
}

// Get returns the adapter for the given version.
func (r *AdapterRegistry) Get(version string) (Adapter, error) {
	// Exact match first
	if a, ok := r.adapters[version]; ok {
		return a, nil
	}

	// Try prefix match for minor versions (e.g., "3.1" matches "3.1.0")
	for v, a := range r.adapters {
		if strings.HasPrefix(v, version) || strings.HasPrefix(version, strings.TrimSuffix(v, ".0")) {
			return a, nil
		}
	}

	return nil, fmt.Errorf("no adapter for OpenAPI version %s", version)
}

// SupportedVersions returns list of supported versions.
func (r *AdapterRegistry) SupportedVersions() []string {
	versions := make([]string, 0, len(r.adapters))
	for v := range r.adapters {
		versions = append(versions, v)
	}
	return versions
}
