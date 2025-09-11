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
					err := mergePathItems(existingPathItem, pathItem)
					if err != nil {
						return nil, fmt.Errorf("error merging path %s from spec %d: %v", path, i, err)
					}

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

func mergePathItems(existing, new *openapi3.PathItem) error {
	// TODO: error and log warn when new one overrides existing one
	if new.Get != nil {
		existing.Get = new.Get
	}
	if new.Put != nil {
		existing.Put = new.Put
	}
	if new.Post != nil {
		existing.Post = new.Post
	}
	if new.Delete != nil {
		existing.Delete = new.Delete
	}
	if new.Options != nil {
		existing.Options = new.Options
	}
	if new.Head != nil {
		existing.Head = new.Head
	}
	if new.Patch != nil {
		existing.Patch = new.Patch
	}
	if new.Trace != nil {
		existing.Trace = new.Trace
	}

	if new.Parameters != nil {
		existing.Parameters = mergeParameters(existing.Parameters, new.Parameters)
	}

	if new.Servers != nil {
		existing.Servers = append(existing.Servers, new.Servers...)
	}

	return nil
}

func mergeComponents(target, source *openapi3.Components, specIndex int) error {
	if source.Schemas != nil {
		maps.Copy(target.Schemas, source.Schemas)
	}

	if source.Parameters != nil {
		for name, param := range source.Parameters {
			if _, exists := target.Parameters[name]; exists {
				return fmt.Errorf("parameter %s already exists (from spec %d)", name, specIndex)
			}
			target.Parameters[name] = param
		}
	}

	if source.Headers != nil {
		for name, header := range source.Headers {
			if _, exists := target.Headers[name]; exists {
				return fmt.Errorf("header %s already exists (from spec %d)", name, specIndex)
			}
			target.Headers[name] = header
		}
	}

	if source.RequestBodies != nil {
		for name, requestBody := range source.RequestBodies {
			if _, exists := target.RequestBodies[name]; exists {
				return fmt.Errorf("request body %s already exists (from spec %d)", name, specIndex)
			}
			target.RequestBodies[name] = requestBody
		}
	}

	if source.Responses != nil {
		for name, response := range source.Responses {
			if _, exists := target.Responses[name]; exists {
				return fmt.Errorf("response %s already exists (from spec %d)", name, specIndex)
			}
			target.Responses[name] = response
		}
	}

	if source.SecuritySchemes != nil {
		for name, securityScheme := range source.SecuritySchemes {
			if _, exists := target.SecuritySchemes[name]; exists {
				return fmt.Errorf("security scheme %s already exists (from spec %d)", name, specIndex)
			}
			target.SecuritySchemes[name] = securityScheme
		}
	}

	if source.Examples != nil {
		for name, example := range source.Examples {
			if _, exists := target.Examples[name]; exists {
				return fmt.Errorf("example %s already exists (from spec %d)", name, specIndex)
			}
			target.Examples[name] = example
		}
	}

	if source.Links != nil {
		for name, link := range source.Links {
			if _, exists := target.Links[name]; exists {
				return fmt.Errorf("link %s already exists (from spec %d)", name, specIndex)
			}
			target.Links[name] = link
		}
	}

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

func mergeParameters(existing, new openapi3.Parameters) openapi3.Parameters {
	if len(existing) == 0 {
		return new
	}
	if len(new) == 0 {
		return existing
	}

	paramMap := make(map[string]*openapi3.ParameterRef)
	var result openapi3.Parameters

	for _, param := range existing {
		if param != nil && param.Value != nil {
			key := getParameterKey(param.Value)
			paramMap[key] = param
			result = append(result, param)
		}
	}

	for _, param := range new {
		if param != nil && param.Value != nil {
			key := getParameterKey(param.Value)
			if existingParam, exists := paramMap[key]; exists {
				for i, resultParam := range result {
					if resultParam == existingParam {
						result[i] = param
						break
					}
				}
				paramMap[key] = param
			} else {
				// Add new parameter
				paramMap[key] = param
				result = append(result, param)
			}
		}
	}

	return result
}

func getParameterKey(param *openapi3.Parameter) string {
	return param.Name + ":" + param.In
}
