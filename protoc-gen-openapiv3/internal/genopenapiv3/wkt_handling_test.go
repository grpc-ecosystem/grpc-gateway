package genopenapiv3

import (
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

// wktTestProto is a proto with nested google.protobuf.Timestamp fields
// to test well-known type handling.
// Includes the google/protobuf/timestamp.proto dependency.
const wktTestProto = `
file_to_generate: "test/v1/wkt.proto"
proto_file: {
	name: "google/protobuf/timestamp.proto"
	package: "google.protobuf"
	message_type: {
		name: "Timestamp"
		field: {
			name: "seconds"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT64
			json_name: "seconds"
		}
		field: {
			name: "nanos"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "nanos"
		}
	}
	options: {
		go_package: "google.golang.org/protobuf/types/known/timestamppb"
	}
	syntax: "proto3"
}
proto_file: {
	name: "test/v1/wkt.proto"
	package: "test.v1"
	dependency: "google/protobuf/timestamp.proto"
	message_type: {
		name: "Nested"
		field: {
			name: "nested_timestamp"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "nestedTimestamp"
		}
	}
	message_type: {
		name: "TestRequest"
		field: {
			name: "timestamp"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "timestamp"
		}
		field: {
			name: "nested"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Nested"
			json_name: "nested"
		}
	}
	message_type: {
		name: "TestResponse"
		field: {
			name: "result"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "result"
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "Test"
			input_type: ".test.v1.TestRequest"
			output_type: ".test.v1.TestResponse"
			options: {
				[google.api.http]: {
					post: "/v1/test"
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

// wktMultiLevelProto tests multi-level nesting with WKTs
const wktMultiLevelProto = `
file_to_generate: "test/v1/multi.proto"
proto_file: {
	name: "google/protobuf/timestamp.proto"
	package: "google.protobuf"
	message_type: {
		name: "Timestamp"
		field: {
			name: "seconds"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT64
			json_name: "seconds"
		}
		field: {
			name: "nanos"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "nanos"
		}
	}
	options: {
		go_package: "google.golang.org/protobuf/types/known/timestamppb"
	}
	syntax: "proto3"
}
proto_file: {
	name: "google/protobuf/duration.proto"
	package: "google.protobuf"
	message_type: {
		name: "Duration"
		field: {
			name: "seconds"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT64
			json_name: "seconds"
		}
		field: {
			name: "nanos"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "nanos"
		}
	}
	options: {
		go_package: "google.golang.org/protobuf/types/known/durationpb"
	}
	syntax: "proto3"
}
proto_file: {
	name: "google/protobuf/wrappers.proto"
	package: "google.protobuf"
	message_type: {
		name: "StringValue"
		field: {
			name: "value"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "value"
		}
	}
	message_type: {
		name: "Int32Value"
		field: {
			name: "value"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "value"
		}
	}
	options: {
		go_package: "google.golang.org/protobuf/types/known/wrapperspb"
	}
	syntax: "proto3"
}
proto_file: {
	name: "test/v1/multi.proto"
	package: "test.v1"
	dependency: "google/protobuf/timestamp.proto"
	dependency: "google/protobuf/duration.proto"
	dependency: "google/protobuf/wrappers.proto"
	message_type: {
		name: "Level3"
		field: {
			name: "timestamp"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "timestamp"
		}
		field: {
			name: "duration"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Duration"
			json_name: "duration"
		}
	}
	message_type: {
		name: "Level2"
		field: {
			name: "level3"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Level3"
			json_name: "level3"
		}
		field: {
			name: "string_value"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.StringValue"
			json_name: "stringValue"
		}
	}
	message_type: {
		name: "Level1"
		field: {
			name: "level2"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Level2"
			json_name: "level2"
		}
		field: {
			name: "int_value"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Int32Value"
			json_name: "intValue"
		}
	}
	message_type: {
		name: "TestRequest"
		field: {
			name: "level1"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Level1"
			json_name: "level1"
		}
	}
	message_type: {
		name: "TestResponse"
		field: {
			name: "result"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "result"
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "Test"
			input_type: ".test.v1.TestRequest"
			output_type: ".test.v1.TestResponse"
			options: {
				[google.api.http]: {
					post: "/v1/test"
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

// TestWKT_DefaultInlinesAndSkipsSchema tests that by default (wkt_as_refs=false),
// well-known types are inlined directly in the schema and no separate WKT schema
// is generated in components/schemas.
func TestWKT_DefaultInlinesAndSkipsSchema(t *testing.T) {
	t.Parallel()

	// Default behavior: WKTs should be inlined, no WKT schema generated
	resp := requireGenerateInline(t, wktTestProto, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Should NOT have any Timestamp schema in components/schemas (with any naming strategy)
	// Legacy naming produces "protobufTimestamp", fqn produces "google.protobuf.Timestamp"
	if strings.Contains(content, `"protobufTimestamp"`) || strings.Contains(content, `"google.protobuf.Timestamp"`) {
		t.Errorf("expected no Timestamp schema when wkt_as_refs=false (default)\nContent: %s", content)
	}

	// Should have inlined date-time format in nested timestamp field
	if !strings.Contains(content, `"format": "date-time"`) {
		t.Errorf("expected inlined date-time format for timestamp fields\nContent: %s", content)
	}
}

// TestWKT_AsRefsGeneratesSchemaAndUsesRef tests that when wkt_as_refs=true,
// well-known types use $ref instead of inlining, and the WKT schema is generated
// in components/schemas.
func TestWKT_AsRefsGeneratesSchemaAndUsesRef(t *testing.T) {
	t.Parallel()

	// With flag: WKTs should use $ref, schema should be generated
	resp := requireGenerateInline(t, wktTestProto, func(reg *descriptor.Registry) {
		reg.SetWKTAsRefs(true)
	})

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Should have Timestamp schema in components/schemas
	// Legacy naming produces "protobufTimestamp"
	if !strings.Contains(content, `"protobufTimestamp"`) {
		t.Errorf("expected protobufTimestamp schema when wkt_as_refs=true\nContent: %s", content)
	}

	// Should have $ref to the WKT schema (with legacy naming)
	if !strings.Contains(content, `"$ref": "#/components/schemas/protobufTimestamp"`) {
		t.Errorf("expected $ref to protobufTimestamp when wkt_as_refs=true\nContent: %s", content)
	}
}

// TestWKT_DefaultIsFalse tests that wkt_as_refs defaults to false.
func TestWKT_DefaultIsFalse(t *testing.T) {
	t.Parallel()

	reg := descriptor.NewRegistry()
	// Should be false by default
	if reg.GetWKTAsRefs() {
		t.Error("expected wkt_as_refs to default to false")
	}
}

// TestWKT_MultipleWKTsDefaultInlines tests that multiple different WKTs
// are all inlined when wkt_as_refs=false.
func TestWKT_MultipleWKTsDefaultInlines(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, wktMultiLevelProto, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Should NOT have any WKT schemas in components/schemas
	wktSchemas := []string{
		`"protobufTimestamp"`,
		`"protobufDuration"`,
		`"protobufStringValue"`,
		`"protobufInt32Value"`,
		`"google.protobuf.Timestamp"`,
		`"google.protobuf.Duration"`,
		`"google.protobuf.StringValue"`,
		`"google.protobuf.Int32Value"`,
	}
	for _, wkt := range wktSchemas {
		if strings.Contains(content, wkt) {
			t.Errorf("expected no %s schema when wkt_as_refs=false (default)\nContent: %s", wkt, content)
		}
	}

	// Should have inlined formats for WKTs
	if !strings.Contains(content, `"format": "date-time"`) {
		t.Errorf("expected inlined date-time format for timestamp fields\nContent: %s", content)
	}
}

// TestWKT_MultipleWKTsAsRefs tests that multiple different WKTs
// all use $ref when wkt_as_refs=true.
func TestWKT_MultipleWKTsAsRefs(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, wktMultiLevelProto, func(reg *descriptor.Registry) {
		reg.SetWKTAsRefs(true)
	})

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Should have WKT schemas in components/schemas
	wktSchemas := []string{
		`"protobufTimestamp"`,
		`"protobufDuration"`,
		`"protobufStringValue"`,
		`"protobufInt32Value"`,
	}
	for _, wkt := range wktSchemas {
		if !strings.Contains(content, wkt) {
			t.Errorf("expected %s schema when wkt_as_refs=true\nContent: %s", wkt, content)
		}
	}

	// Should have $refs to WKT schemas
	wktRefs := []string{
		`"$ref": "#/components/schemas/protobufTimestamp"`,
		`"$ref": "#/components/schemas/protobufDuration"`,
		`"$ref": "#/components/schemas/protobufStringValue"`,
		`"$ref": "#/components/schemas/protobufInt32Value"`,
	}
	for _, ref := range wktRefs {
		if !strings.Contains(content, ref) {
			t.Errorf("expected %s when wkt_as_refs=true\nContent: %s", ref, content)
		}
	}
}

// TestWKT_MultiLevelNestingDefault tests that multi-level nesting works correctly
// with default settings (WKTs inlined at all levels).
func TestWKT_MultiLevelNestingDefault(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, wktMultiLevelProto, nil)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// All level schemas should exist
	levelSchemas := []string{
		`"v1Level1"`,
		`"v1Level2"`,
		`"v1Level3"`,
	}
	for _, schema := range levelSchemas {
		if !strings.Contains(content, schema) {
			t.Errorf("expected %s schema\nContent: %s", schema, content)
		}
	}

	// WKT schemas should NOT exist
	if strings.Contains(content, `"protobufTimestamp"`) {
		t.Errorf("expected no protobufTimestamp schema\nContent: %s", content)
	}

	// Inlined WKT formats should be present in the level schemas
	// Level3 has timestamp (date-time) and duration (string)
	if !strings.Contains(content, `"format": "date-time"`) {
		t.Errorf("expected inlined date-time format\nContent: %s", content)
	}
}

// TestWKT_MultiLevelNestingAsRefs tests that multi-level nesting works correctly
// with wkt_as_refs=true (all WKTs use $ref at all levels).
func TestWKT_MultiLevelNestingAsRefs(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, wktMultiLevelProto, func(reg *descriptor.Registry) {
		reg.SetWKTAsRefs(true)
	})

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// All level schemas should exist
	levelSchemas := []string{
		`"v1Level1"`,
		`"v1Level2"`,
		`"v1Level3"`,
	}
	for _, schema := range levelSchemas {
		if !strings.Contains(content, schema) {
			t.Errorf("expected %s schema\nContent: %s", schema, content)
		}
	}

	// WKT schemas should exist
	wktSchemas := []string{
		`"protobufTimestamp"`,
		`"protobufDuration"`,
	}
	for _, schema := range wktSchemas {
		if !strings.Contains(content, schema) {
			t.Errorf("expected %s schema\nContent: %s", schema, content)
		}
	}
}

// wktDuplicateProto tests that the same WKT used in multiple messages
// results in only one WKT schema being generated.
const wktDuplicateProto = `
file_to_generate: "test/v1/dup.proto"
proto_file: {
	name: "google/protobuf/timestamp.proto"
	package: "google.protobuf"
	message_type: {
		name: "Timestamp"
		field: {
			name: "seconds"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_INT64
			json_name: "seconds"
		}
		field: {
			name: "nanos"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_INT32
			json_name: "nanos"
		}
	}
	options: {
		go_package: "google.golang.org/protobuf/types/known/timestamppb"
	}
	syntax: "proto3"
}
proto_file: {
	name: "test/v1/dup.proto"
	package: "test.v1"
	dependency: "google/protobuf/timestamp.proto"
	message_type: {
		name: "Message1"
		field: {
			name: "created_at"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "createdAt"
		}
		field: {
			name: "updated_at"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "updatedAt"
		}
	}
	message_type: {
		name: "Message2"
		field: {
			name: "timestamp"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "timestamp"
		}
	}
	message_type: {
		name: "Message3"
		field: {
			name: "start_time"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "startTime"
		}
		field: {
			name: "end_time"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".google.protobuf.Timestamp"
			json_name: "endTime"
		}
	}
	message_type: {
		name: "TestRequest"
		field: {
			name: "msg1"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Message1"
			json_name: "msg1"
		}
		field: {
			name: "msg2"
			number: 2
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Message2"
			json_name: "msg2"
		}
		field: {
			name: "msg3"
			number: 3
			label: LABEL_OPTIONAL
			type: TYPE_MESSAGE
			type_name: ".test.v1.Message3"
			json_name: "msg3"
		}
	}
	message_type: {
		name: "TestResponse"
		field: {
			name: "result"
			number: 1
			label: LABEL_OPTIONAL
			type: TYPE_STRING
			json_name: "result"
		}
	}
	service: {
		name: "TestService"
		method: {
			name: "Test"
			input_type: ".test.v1.TestRequest"
			output_type: ".test.v1.TestResponse"
			options: {
				[google.api.http]: {
					post: "/v1/test"
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

// TestWKT_DuplicateWKTGeneratesOnce tests that when the same WKT is used
// in multiple messages, the schema is only generated once in components/schemas.
func TestWKT_DuplicateWKTGeneratesOnce(t *testing.T) {
	t.Parallel()

	resp := requireGenerateInline(t, wktDuplicateProto, func(reg *descriptor.Registry) {
		reg.SetWKTAsRefs(true)
	})

	if len(resp) != 1 {
		t.Fatalf("expected 1 response file, got %d", len(resp))
	}

	content := resp[0].GetContent()

	// Count occurrences of the schema definition
	// The schema key should appear exactly once as a definition in components/schemas
	schemaKey := `"protobufTimestamp": {`
	count := strings.Count(content, schemaKey)
	if count != 1 {
		t.Errorf("expected protobufTimestamp schema to be defined exactly once, found %d occurrences\nContent: %s", count, content)
	}

	// But the $ref should appear multiple times (5 timestamp fields across 3 messages)
	refPattern := `"$ref": "#/components/schemas/protobufTimestamp"`
	refCount := strings.Count(content, refPattern)
	if refCount < 5 {
		t.Errorf("expected at least 5 $ref to protobufTimestamp, found %d\nContent: %s", refCount, content)
	}
}
