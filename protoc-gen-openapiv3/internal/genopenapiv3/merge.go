package genopenapiv3

import (
	"fmt"
	"maps"

	"github.com/getkin/kin-openapi/openapi3"
)

func MergeOpenAPISpecs(specs ...*openapi3.T) (*openapi3.T, error) {
	if len(specs) == 0 {
		return nil, fmt.Errorf("at least one OpenAPI spec is required")
	}

	merged := &openapi3.T{
		OpenAPI: specs[0].OpenAPI,
		Info:    specs[0].Info,
		Servers: make(openapi3.Servers, 0),
		Paths:   new(openapi3.Paths),
		Components: &openapi3.Components{
			Schemas:         make(openapi3.Schemas),
			Parameters:      make(openapi3.ParametersMap),
			Headers:         make(openapi3.Headers),
			RequestBodies:   make(openapi3.RequestBodies),
			Responses:       make(openapi3.ResponseBodies),
			SecuritySchemes: make(openapi3.SecuritySchemes),
			Examples:        make(openapi3.Examples),
			Links:           make(openapi3.Links),
			Callbacks:       make(openapi3.Callbacks),
		},
		Security: make(openapi3.SecurityRequirements, 0),
		Tags:     make(openapi3.Tags, 0),
	}

	serverMap := make(map[string]bool)
	tagMap := make(map[string]bool)

	for i, spec := range specs {
		if spec == nil {
			continue
		}

		if spec.Servers != nil {
			for _, server := range spec.Servers {
				serverKey := server.URL
				if !serverMap[serverKey] {
					merged.Servers = append(merged.Servers, server)
					serverMap[serverKey] = true
				}
			}
		}

		if spec.Paths != nil {
			for path, pathItem := range spec.Paths.Map() {
				if existingPathItem, exists := merged.Paths.Map()[path]; exists {
					mergedPathItem, err := mergePathItems(existingPathItem, pathItem)
					if err != nil {
						return nil, fmt.Errorf("error merging path %s from spec %d: %v", path, i, err)
					}

					merged.Paths.Set(path, mergedPathItem)
				} else {

					merged.Paths.Set(path, pathItem)
				}
			}
		}

		if spec.Components != nil {
			if err := mergeComponents(merged.Components, spec.Components, i); err != nil {
				return nil, fmt.Errorf("error merging components from spec %d: %v", i, err)
			}
		}

		if spec.Tags != nil {
			for _, tag := range spec.Tags {
				if !tagMap[tag.Name] {
					merged.Tags = append(merged.Tags, tag)
					tagMap[tag.Name] = true
				}
			}
		}

		if spec.Security != nil {
			merged.Security = append(merged.Security, spec.Security...)
		}
	}

	return merged, nil
}

// mergePathItems merges two path items, combining their operations
func mergePathItems(existing, new *openapi3.PathItem) (*openapi3.PathItem, error) {
	merged := &openapi3.PathItem{
		Ref:         existing.Ref,
		Summary:     existing.Summary,
		Description: existing.Description,
		Servers:     existing.Servers,
		Parameters:  existing.Parameters,
	}

	// Copy existing operations
	if existing.Get != nil {
		merged.Get = existing.Get
	}
	if existing.Put != nil {
		merged.Put = existing.Put
	}
	if existing.Post != nil {
		merged.Post = existing.Post
	}
	if existing.Delete != nil {
		merged.Delete = existing.Delete
	}
	if existing.Options != nil {
		merged.Options = existing.Options
	}
	if existing.Head != nil {
		merged.Head = existing.Head
	}
	if existing.Patch != nil {
		merged.Patch = existing.Patch
	}
	if existing.Trace != nil {
		merged.Trace = existing.Trace
	}

	// Add new operations (will overwrite if same method exists)
	if new.Get != nil {
		if merged.Get != nil {
			return nil, fmt.Errorf("conflicting GET operation")
		}
		merged.Get = new.Get
	}
	if new.Put != nil {
		if merged.Put != nil {
			return nil, fmt.Errorf("conflicting PUT operation")
		}
		merged.Put = new.Put
	}
	if new.Post != nil {
		if merged.Post != nil {
			return nil, fmt.Errorf("conflicting POST operation")
		}
		merged.Post = new.Post
	}
	if new.Delete != nil {
		if merged.Delete != nil {
			return nil, fmt.Errorf("conflicting DELETE operation")
		}
		merged.Delete = new.Delete
	}
	if new.Options != nil {
		if merged.Options != nil {
			return nil, fmt.Errorf("conflicting OPTIONS operation")
		}
		merged.Options = new.Options
	}
	if new.Head != nil {
		if merged.Head != nil {
			return nil, fmt.Errorf("conflicting HEAD operation")
		}
		merged.Head = new.Head
	}
	if new.Patch != nil {
		if merged.Patch != nil {
			return nil, fmt.Errorf("conflicting PATCH operation")
		}
		merged.Patch = new.Patch
	}
	if new.Trace != nil {
		if merged.Trace != nil {
			return nil, fmt.Errorf("conflicting TRACE operation")
		}
		merged.Trace = new.Trace
	}

	// Merge parameters
	if new.Parameters != nil {
		merged.Parameters = append(merged.Parameters, new.Parameters...)
	}

	// Merge servers
	if new.Servers != nil {
		merged.Servers = append(merged.Servers, new.Servers...)
	}

	return merged, nil
}

// mergeComponents merges components from source into target
func mergeComponents(target, source *openapi3.Components, specIndex int) error {
	// Merge schemas
	if source.Schemas != nil {
		maps.Copy(target.Schemas, source.Schemas)
	}

	// Merge parameters
	if source.Parameters != nil {
		for name, param := range source.Parameters {
			if _, exists := target.Parameters[name]; exists {
				return fmt.Errorf("parameter %s already exists (from spec %d)", name, specIndex)
			}
			target.Parameters[name] = param
		}
	}

	// Merge headers
	if source.Headers != nil {
		for name, header := range source.Headers {
			if _, exists := target.Headers[name]; exists {
				return fmt.Errorf("header %s already exists (from spec %d)", name, specIndex)
			}
			target.Headers[name] = header
		}
	}

	// Merge request bodies
	if source.RequestBodies != nil {
		for name, requestBody := range source.RequestBodies {
			if _, exists := target.RequestBodies[name]; exists {
				return fmt.Errorf("request body %s already exists (from spec %d)", name, specIndex)
			}
			target.RequestBodies[name] = requestBody
		}
	}

	// Merge responses
	if source.Responses != nil {
		for name, response := range source.Responses {
			if _, exists := target.Responses[name]; exists {
				return fmt.Errorf("response %s already exists (from spec %d)", name, specIndex)
			}
			target.Responses[name] = response
		}
	}

	// Merge security schemes
	if source.SecuritySchemes != nil {
		for name, securityScheme := range source.SecuritySchemes {
			if _, exists := target.SecuritySchemes[name]; exists {
				return fmt.Errorf("security scheme %s already exists (from spec %d)", name, specIndex)
			}
			target.SecuritySchemes[name] = securityScheme
		}
	}

	// Merge examples
	if source.Examples != nil {
		for name, example := range source.Examples {
			if _, exists := target.Examples[name]; exists {
				return fmt.Errorf("example %s already exists (from spec %d)", name, specIndex)
			}
			target.Examples[name] = example
		}
	}

	// Merge links
	if source.Links != nil {
		for name, link := range source.Links {
			if _, exists := target.Links[name]; exists {
				return fmt.Errorf("link %s already exists (from spec %d)", name, specIndex)
			}
			target.Links[name] = link
		}
	}

	// Merge callbacks
	if source.Callbacks != nil {
		for name, callback := range source.Callbacks {
			if _, exists := target.Callbacks[name]; exists {
				return fmt.Errorf("callback %s already exists (from spec %d)", name, specIndex)
			}
			target.Callbacks[name] = callback
		}
	}

	return nil
}
