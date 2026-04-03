package genopenapiv3

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/pluginpb"
)

// validateRefs extracts all $ref values and checks if corresponding schemas exist.
// Returns a list of schema names that are referenced but don't exist.
func validateRefs(t *testing.T, jsonContent string) []string {
	t.Helper()

	var doc map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &doc); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Extract all schema names from components/schemas
	existingSchemas := make(map[string]bool)
	if components, ok := doc["components"].(map[string]any); ok {
		if schemas, ok := components["schemas"].(map[string]any); ok {
			for name := range schemas {
				existingSchemas[name] = true
			}
		}
	}

	// Find all $refs in the document
	refPattern := regexp.MustCompile(`"\$ref"\s*:\s*"#/components/schemas/([^"]+)"`)
	matches := refPattern.FindAllStringSubmatch(jsonContent, -1)

	var missing []string
	seen := make(map[string]bool)
	for _, match := range matches {
		schemaName := match[1]
		if seen[schemaName] {
			continue
		}
		seen[schemaName] = true

		if !existingSchemas[schemaName] {
			missing = append(missing, schemaName)
		}
	}

	return missing
}

// TestRefValidation_FindsMissingSchemas verifies that our validation function
// actually catches missing schemas.
func TestRefValidation_FindsMissingSchemas(t *testing.T) {
	badJSON := `{
		"components": {
			"schemas": {
				"ExistingSchema": {"type": "object"}
			}
		},
		"paths": {
			"/test": {
				"get": {
					"responses": {
						"200": {
							"content": {
								"application/json": {
									"schema": {"$ref": "#/components/schemas/MissingSchema"}
								}
							}
						}
					}
				}
			}
		}
	}`

	missing := validateRefs(t, badJSON)
	if len(missing) != 1 || missing[0] != "MissingSchema" {
		t.Errorf("expected to find MissingSchema, got: %v", missing)
	}
}

// TestRefValidation_PassesForValidJSON verifies validation passes for valid JSON.
func TestRefValidation_PassesForValidJSON(t *testing.T) {
	goodJSON := `{
		"components": {
			"schemas": {
				"MySchema": {"type": "object"},
				"google.rpc.Status": {"type": "object"}
			}
		},
		"paths": {
			"/test": {
				"get": {
					"responses": {
						"200": {
							"content": {
								"application/json": {
									"schema": {"$ref": "#/components/schemas/MySchema"}
								}
							}
						},
						"default": {
							"content": {
								"application/json": {
									"schema": {"$ref": "#/components/schemas/google.rpc.Status"}
								}
							}
						}
					}
				}
			}
		}
	}`

	missing := validateRefs(t, goodJSON)
	if len(missing) != 0 {
		t.Errorf("expected no missing schemas, got: %v", missing)
	}
}

// TestAllRefsHaveSchemas_BasicService validates refs for a basic service.
func TestAllRefsHaveSchemas_BasicService(t *testing.T) {
	prototextInput := `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "google/protobuf/timestamp.proto"
	package: "google.protobuf"
	message_type: {
		name: "Timestamp"
		field: { name: "seconds" number: 1 label: LABEL_OPTIONAL type: TYPE_INT64 json_name: "seconds" }
		field: { name: "nanos" number: 2 label: LABEL_OPTIONAL type: TYPE_INT32 json_name: "nanos" }
	}
	options: { go_package: "google.golang.org/protobuf/types/known/timestamppb" }
	syntax: "proto3"
}
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	dependency: "google/protobuf/timestamp.proto"
	message_type: {
		name: "Nested"
		field: { name: "value" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "value" }
		field: { name: "timestamp" number: 2 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".google.protobuf.Timestamp" json_name: "timestamp" }
	}
	message_type: {
		name: "Request"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
		field: { name: "nested" number: 2 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".test.v1.Nested" json_name: "nested" }
	}
	message_type: {
		name: "Response"
		field: { name: "result" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "result" }
	}
	service: {
		name: "TestService"
		method: {
			name: "Test"
			input_type: ".test.v1.Request"
			output_type: ".test.v1.Response"
			options: { [google.api.http]: { get: "/test" } }
		}
	}
	options: { go_package: "test/v1;testv1" }
	syntax: "proto3"
}
`
	content := generateFromPrototext(t, prototextInput, nil)
	missing := validateRefs(t, content)
	if len(missing) > 0 {
		t.Errorf("found references to missing schemas: %v", missing)
		t.Logf("Generated content:\n%s", content)
	}
}

// TestAllRefsHaveSchemas_WKTAsRefs validates refs with wkt_as_refs=true.
func TestAllRefsHaveSchemas_WKTAsRefs(t *testing.T) {
	prototextInput := `
file_to_generate: "test/v1/wkt.proto"
proto_file: {
	name: "google/protobuf/timestamp.proto"
	package: "google.protobuf"
	message_type: {
		name: "Timestamp"
		field: { name: "seconds" number: 1 label: LABEL_OPTIONAL type: TYPE_INT64 json_name: "seconds" }
		field: { name: "nanos" number: 2 label: LABEL_OPTIONAL type: TYPE_INT32 json_name: "nanos" }
	}
	options: { go_package: "google.golang.org/protobuf/types/known/timestamppb" }
	syntax: "proto3"
}
proto_file: {
	name: "test/v1/wkt.proto"
	package: "test.v1"
	dependency: "google/protobuf/timestamp.proto"
	message_type: {
		name: "Event"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
		field: { name: "created_at" number: 2 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".google.protobuf.Timestamp" json_name: "createdAt" }
	}
	message_type: {
		name: "GetEventRequest"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
	}
	service: {
		name: "EventService"
		method: {
			name: "GetEvent"
			input_type: ".test.v1.GetEventRequest"
			output_type: ".test.v1.Event"
			options: { [google.api.http]: { get: "/events/{id}" } }
		}
	}
	options: { go_package: "test/v1;testv1" }
	syntax: "proto3"
}
`
	content := generateFromPrototext(t, prototextInput, func(reg *descriptor.Registry) {
		reg.SetWKTAsRefs(true)
	})
	missing := validateRefs(t, content)
	if len(missing) > 0 {
		t.Errorf("found references to missing schemas: %v", missing)
		t.Logf("Generated content:\n%s", content)
	}
}

// TestAllRefsHaveSchemas_WKTInlined validates refs with wkt_as_refs=false (default).
func TestAllRefsHaveSchemas_WKTInlined(t *testing.T) {
	prototextInput := `
file_to_generate: "test/v1/wkt.proto"
proto_file: {
	name: "google/protobuf/timestamp.proto"
	package: "google.protobuf"
	message_type: {
		name: "Timestamp"
		field: { name: "seconds" number: 1 label: LABEL_OPTIONAL type: TYPE_INT64 json_name: "seconds" }
		field: { name: "nanos" number: 2 label: LABEL_OPTIONAL type: TYPE_INT32 json_name: "nanos" }
	}
	options: { go_package: "google.golang.org/protobuf/types/known/timestamppb" }
	syntax: "proto3"
}
proto_file: {
	name: "test/v1/wkt.proto"
	package: "test.v1"
	dependency: "google/protobuf/timestamp.proto"
	message_type: {
		name: "Event"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
		field: { name: "created_at" number: 2 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".google.protobuf.Timestamp" json_name: "createdAt" }
	}
	message_type: {
		name: "GetEventRequest"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
	}
	service: {
		name: "EventService"
		method: {
			name: "GetEvent"
			input_type: ".test.v1.GetEventRequest"
			output_type: ".test.v1.Event"
			options: { [google.api.http]: { get: "/events/{id}" } }
		}
	}
	options: { go_package: "test/v1;testv1" }
	syntax: "proto3"
}
`
	content := generateFromPrototext(t, prototextInput, nil)
	missing := validateRefs(t, content)
	if len(missing) > 0 {
		t.Errorf("found references to missing schemas: %v", missing)
		t.Logf("Generated content:\n%s", content)
	}

	// Verify no WKT schema is generated when inlining
	if strings.Contains(content, `"protobufTimestamp"`) {
		t.Error("expected no protobufTimestamp schema when wkt_as_refs=false")
	}
}

// TestAllRefsHaveSchemas_EnumReference validates refs for enums.
func TestAllRefsHaveSchemas_EnumReference(t *testing.T) {
	prototextInput := `
file_to_generate: "test/v1/enum.proto"
proto_file: {
	name: "test/v1/enum.proto"
	package: "test.v1"
	enum_type: {
		name: "Status"
		value: { name: "UNKNOWN" number: 0 }
		value: { name: "ACTIVE" number: 1 }
		value: { name: "INACTIVE" number: 2 }
	}
	message_type: {
		name: "Item"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
		field: { name: "status" number: 2 label: LABEL_OPTIONAL type: TYPE_ENUM type_name: ".test.v1.Status" json_name: "status" }
	}
	message_type: {
		name: "GetItemRequest"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
	}
	service: {
		name: "ItemService"
		method: {
			name: "GetItem"
			input_type: ".test.v1.GetItemRequest"
			output_type: ".test.v1.Item"
			options: { [google.api.http]: { get: "/items/{id}" } }
		}
	}
	options: { go_package: "test/v1;testv1" }
	syntax: "proto3"
}
`
	content := generateFromPrototext(t, prototextInput, nil)
	missing := validateRefs(t, content)
	if len(missing) > 0 {
		t.Errorf("found references to missing schemas: %v", missing)
		t.Logf("Generated content:\n%s", content)
	}
}

// TestAllRefsHaveSchemas_DeeplyNested validates refs for deeply nested messages.
func TestAllRefsHaveSchemas_DeeplyNested(t *testing.T) {
	prototextInput := `
file_to_generate: "test/v1/nested.proto"
proto_file: {
	name: "test/v1/nested.proto"
	package: "test.v1"
	message_type: {
		name: "Level3"
		field: { name: "value" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "value" }
	}
	message_type: {
		name: "Level2"
		field: { name: "level3" number: 1 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".test.v1.Level3" json_name: "level3" }
	}
	message_type: {
		name: "Level1"
		field: { name: "level2" number: 1 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".test.v1.Level2" json_name: "level2" }
	}
	message_type: {
		name: "GetRequest"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
	}
	service: {
		name: "NestedService"
		method: {
			name: "Get"
			input_type: ".test.v1.GetRequest"
			output_type: ".test.v1.Level1"
			options: { [google.api.http]: { get: "/nested/{id}" } }
		}
	}
	options: { go_package: "test/v1;testv1" }
	syntax: "proto3"
}
`
	content := generateFromPrototext(t, prototextInput, nil)
	missing := validateRefs(t, content)
	if len(missing) > 0 {
		t.Errorf("found references to missing schemas: %v", missing)
		t.Logf("Generated content:\n%s", content)
	}
}

// TestAllRefsHaveSchemas_ExternalPackage validates refs from external packages.
func TestAllRefsHaveSchemas_ExternalPackage(t *testing.T) {
	prototextInput := `
file_to_generate: "myapp/v1/service.proto"
proto_file: {
	name: "common/v1/types.proto"
	package: "common.v1"
	message_type: {
		name: "Metadata"
		field: { name: "created_by" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "createdBy" }
		field: { name: "version" number: 2 label: LABEL_OPTIONAL type: TYPE_INT32 json_name: "version" }
	}
	options: { go_package: "common/v1;commonv1" }
	syntax: "proto3"
}
proto_file: {
	name: "myapp/v1/service.proto"
	package: "myapp.v1"
	dependency: "common/v1/types.proto"
	message_type: {
		name: "Resource"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
		field: { name: "metadata" number: 2 label: LABEL_OPTIONAL type: TYPE_MESSAGE type_name: ".common.v1.Metadata" json_name: "metadata" }
	}
	message_type: {
		name: "GetResourceRequest"
		field: { name: "id" number: 1 label: LABEL_OPTIONAL type: TYPE_STRING json_name: "id" }
	}
	service: {
		name: "ResourceService"
		method: {
			name: "GetResource"
			input_type: ".myapp.v1.GetResourceRequest"
			output_type: ".myapp.v1.Resource"
			options: { [google.api.http]: { get: "/resources/{id}" } }
		}
	}
	options: { go_package: "myapp/v1;myappv1" }
	syntax: "proto3"
}
`
	content := generateFromPrototext(t, prototextInput, nil)
	missing := validateRefs(t, content)
	if len(missing) > 0 {
		t.Errorf("found references to missing schemas: %v", missing)
		t.Logf("Generated content:\n%s", content)
	}

	// Verify Metadata from external package is generated
	if !strings.Contains(content, "Metadata") {
		t.Error("expected Metadata schema from external package to be generated")
	}
}

// generateFromPrototext is a helper to generate OpenAPI from prototext input.
func generateFromPrototext(t *testing.T, prototextInput string, configure func(*descriptor.Registry)) string {
	t.Helper()

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(prototextInput), &req); err != nil {
		t.Fatalf("failed to unmarshal prototext: %v", err)
	}

	reg := descriptor.NewRegistry()
	if err := reg.Load(&req); err != nil {
		t.Fatalf("failed to load registry: %v", err)
	}

	if configure != nil {
		configure(reg)
	}

	targets := make([]*descriptor.File, 0)
	for _, name := range req.FileToGenerate {
		f, err := reg.LookupFile(name)
		if err != nil {
			t.Fatalf("failed to lookup file %s: %v", name, err)
		}
		targets = append(targets, f)
	}

	gen, err := New(reg, FormatJSON, "3.1.0")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	resp, err := gen.Generate(targets)
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	return resp[0].GetContent()
}
