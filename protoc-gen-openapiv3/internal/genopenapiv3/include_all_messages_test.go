package genopenapiv3

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

// TestIncludeAllMessages_UnusedMessageIncluded tests that messages not used in any RPC
// are included in the generated OpenAPI schema when include_all_messages is enabled (default).
func TestIncludeAllMessages_UnusedMessageIncluded(t *testing.T) {
	t.Parallel()

	// Proto with an unused message (ExtendedStatus) that is never referenced by any RPC
	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "GetUserRequest"
		field: {
			name: "user_id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "userId"
		}
	}
	message_type: {
		name: "User"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "ExtendedStatus"
		field: {
			name: "code"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "code"
		}
		field: {
			name: "details"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "details"
		}
	}
	service: {
		name: "UserService"
		method: {
			name: "GetUser"
			input_type: ".test.v1.GetUserRequest"
			output_type: ".test.v1.User"
			options: {
				[google.api.http]: {
					get: "/v1/users/{user_id}"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	resp := requireGenerateInline(t, protoText, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Verify the unused message is included in schemas
	// With legacy naming strategy (default), package "test.v1" becomes "v1" prefix
	if !strings.Contains(content, `"v1ExtendedStatus"`) {
		t.Error("expected v1ExtendedStatus to be included in schemas (unused message should be included when include_all_messages=true)")
	}

	// Also verify the RPC-referenced messages are still there
	if !strings.Contains(content, `"v1GetUserRequest"`) {
		t.Error("expected v1GetUserRequest to be included in schemas")
	}
	if !strings.Contains(content, `"v1User"`) {
		t.Error("expected v1User to be included in schemas")
	}
}

// TestIncludeAllMessages_DisabledExcludesUnusedMessages tests that when include_all_messages
// is disabled, unused messages are NOT included in the generated schema.
func TestIncludeAllMessages_DisabledExcludesUnusedMessages(t *testing.T) {
	t.Parallel()

	// Same proto with an unused message
	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "GetUserRequest"
		field: {
			name: "user_id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "userId"
		}
	}
	message_type: {
		name: "User"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "ExtendedStatus"
		field: {
			name: "code"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "code"
		}
	}
	service: {
		name: "UserService"
		method: {
			name: "GetUser"
			input_type: ".test.v1.GetUserRequest"
			output_type: ".test.v1.User"
			options: {
				[google.api.http]: {
					get: "/v1/users/{user_id}"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	// Disable include_all_messages
	resp := requireGenerateInline(t, protoText, func(reg *descriptor.Registry) {
		reg.SetIncludeAllMessages(false)
	})

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Verify the unused message is NOT included when flag is disabled
	if strings.Contains(content, `"v1ExtendedStatus"`) {
		t.Error("expected v1ExtendedStatus to NOT be included in schemas when include_all_messages=false")
	}

	// RPC-referenced response message should still be there
	// Note: v1GetUserRequest is NOT in schemas because for GET methods,
	// request fields become path/query parameters, not a request body schema
	if !strings.Contains(content, `"v1User"`) {
		t.Errorf("expected v1User to be included in schemas.\nContent: %s", content)
	}
}

// TestIncludeAllMessages_NestedMessages tests that nested messages are included.
func TestIncludeAllMessages_NestedMessages(t *testing.T) {
	t.Parallel()

	// Proto with nested messages that are not directly used in RPCs
	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "Request"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "Response"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "OuterMessage"
		field: {
			name: "name"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "name"
		}
		nested_type: {
			name: "InnerMessage"
			field: {
				name: "value"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "value"
			}
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "DoSomething"
			input_type: ".test.v1.Request"
			output_type: ".test.v1.Response"
			options: {
				[google.api.http]: {
					post: "/v1/do"
					body: "*"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	resp := requireGenerateInline(t, protoText, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Verify the outer message is included (legacy naming: v1OuterMessage)
	if !strings.Contains(content, `"v1OuterMessage"`) {
		t.Error("expected v1OuterMessage to be included in schemas")
	}

	// Verify the nested message is included (legacy naming: OuterMessageInnerMessage)
	if !strings.Contains(content, `"OuterMessageInnerMessage"`) {
		t.Error("expected OuterMessageInnerMessage to be included in schemas")
	}
}

// TestIncludeAllMessages_MessageReferencedByNonRPCMessage tests that messages only
// referenced by other non-RPC messages are included.
func TestIncludeAllMessages_MessageReferencedByNonRPCMessage(t *testing.T) {
	t.Parallel()

	// Proto where ErrorDetail is referenced by ExtendedStatus (unused), but not by any RPC
	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "Request"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "Response"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "ErrorDetail"
		field: {
			name: "reason"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "reason"
		}
	}
	message_type: {
		name: "ExtendedStatus"
		field: {
			name: "code"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "code"
		}
		field: {
			name: "detail"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.ErrorDetail"
			json_name: "detail"
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "DoSomething"
			input_type: ".test.v1.Request"
			output_type: ".test.v1.Response"
			options: {
				[google.api.http]: {
					post: "/v1/do"
					body: "*"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	resp := requireGenerateInline(t, protoText, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Verify ExtendedStatus is included (legacy naming: v1ExtendedStatus)
	if !strings.Contains(content, `"v1ExtendedStatus"`) {
		t.Error("expected v1ExtendedStatus to be included in schemas")
	}

	// Verify ErrorDetail is included (referenced by ExtendedStatus but not by any RPC)
	if !strings.Contains(content, `"v1ErrorDetail"`) {
		t.Error("expected v1ErrorDetail to be included in schemas (referenced by non-RPC message)")
	}
}

// TestIncludeAllMessages_UnusedEnumIncluded tests that enums not used in any RPC
// are included when include_all_messages is enabled.
func TestIncludeAllMessages_UnusedEnumIncluded(t *testing.T) {
	t.Parallel()

	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "Request"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "Response"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	enum_type: {
		name: "UnusedStatus"
		value: {
			name: "UNUSED_STATUS_UNSPECIFIED"
			number: 0
		}
		value: {
			name: "UNUSED_STATUS_ACTIVE"
			number: 1
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "DoSomething"
			input_type: ".test.v1.Request"
			output_type: ".test.v1.Response"
			options: {
				[google.api.http]: {
					post: "/v1/do"
					body: "*"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	resp := requireGenerateInline(t, protoText, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Verify the unused enum is included (legacy naming: v1UnusedStatus)
	if !strings.Contains(content, `"v1UnusedStatus"`) {
		t.Error("expected v1UnusedStatus enum to be included in schemas when include_all_messages=true")
	}
}

// TestIncludeAllMessages_DefaultIsTrue tests that include_all_messages defaults to true.
func TestIncludeAllMessages_DefaultIsTrue(t *testing.T) {
	t.Parallel()

	reg := descriptor.NewRegistry()
	// Should be true by default
	if !reg.GetIncludeAllMessages() {
		t.Error("expected include_all_messages to default to true")
	}
}

// TestIncludeAllMessages_VerifySchemaStructure tests that the schema structure is correct
// for included unused messages.
func TestIncludeAllMessages_VerifySchemaStructure(t *testing.T) {
	t.Parallel()

	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "Request"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "Response"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "UnusedMessage"
		field: {
			name: "stringField"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "stringField"
		}
		field: {
			name: "intField"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "intField"
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "DoSomething"
			input_type: ".test.v1.Request"
			output_type: ".test.v1.Response"
			options: {
				[google.api.http]: {
					post: "/v1/do"
					body: "*"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	resp := requireGenerateInline(t, protoText, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Parse the JSON to verify structure
	var openAPI map[string]any
	if err := json.Unmarshal([]byte(content), &openAPI); err != nil {
		t.Fatalf("failed to parse OpenAPI JSON: %v", err)
	}

	components, ok := openAPI["components"].(map[string]any)
	if !ok {
		t.Fatal("expected components to be present")
	}

	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		t.Fatal("expected schemas to be present")
	}

	// Find the UnusedMessage schema (legacy naming: v1UnusedMessage)
	unusedMessageSchema, ok := schemas["v1UnusedMessage"].(map[string]any)
	if !ok {
		t.Fatal("expected v1UnusedMessage schema to be present")
	}

	// Verify the schema has the correct structure
	if unusedMessageSchema["type"] != "object" {
		t.Errorf("expected v1UnusedMessage type to be 'object', got %v", unusedMessageSchema["type"])
	}

	properties, ok := unusedMessageSchema["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected UnusedMessage to have properties")
	}

	// Verify fields are present
	if _, ok := properties["stringField"]; !ok {
		t.Error("expected stringField property in UnusedMessage")
	}
	if _, ok := properties["intField"]; !ok {
		t.Error("expected intField property in UnusedMessage")
	}
}

// TestIncludeAllMessages_NamingStrategies tests that the feature works with all naming strategies.
func TestIncludeAllMessages_NamingStrategies(t *testing.T) {
	t.Parallel()

	const protoText = `
file_to_generate: "test/v1/service.proto"
proto_file: {
	name: "test/v1/service.proto"
	package: "test.v1"
	message_type: {
		name: "Request"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "Response"
		field: {
			name: "id"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "id"
		}
	}
	message_type: {
		name: "UnusedMessage"
		field: {
			name: "value"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "value"
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "DoSomething"
			input_type: ".test.v1.Request"
			output_type: ".test.v1.Response"
			options: {
				[google.api.http]: {
					post: "/v1/do"
					body: "*"
				}
			}
		}
	}
	options: {
		go_package: "test/v1;testv1"
	}
	syntax: "proto3"
}
`
	tests := []struct {
		name             string
		strategy         string
		wantUnusedSchema string
	}{
		{"legacy", "legacy", `"v1UnusedMessage"`},
		{"simple", "simple", `"UnusedMessage"`},
		{"fqn", "fqn", `"test.v1.UnusedMessage"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := requireGenerateInline(t, protoText, func(reg *descriptor.Registry) {
				reg.SetOpenAPINamingStrategy(tt.strategy)
			})

			if len(resp) != 1 {
				t.Fatalf("expected 1 response file, got %d", len(resp))
			}

			content := resp[0].GetContent()

			// Verify the unused message is included with the correct naming
			if !strings.Contains(content, tt.wantUnusedSchema) {
				t.Errorf("expected schema name %s with strategy %s, content: %s", tt.wantUnusedSchema, tt.strategy, content)
			}
		})
	}
}

