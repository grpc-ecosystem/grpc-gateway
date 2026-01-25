package genopenapi_test

import (
	"bytes"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi"
	"go.yaml.in/yaml/v3"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerate_YAML(t *testing.T) {
	t.Parallel()

	req := &pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name:    proto.String("file.proto"),
			Package: proto.String("example"),
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("goexample/v1;goexample"),
			},
		}},
		FileToGenerate: []string{
			"file.proto",
		},
	}

	resp := requireGenerate(t, req, genopenapi.FormatYAML, false, false)
	if len(resp) != 1 {
		t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
	}

	var p map[string]interface{}
	err := yaml.Unmarshal([]byte(resp[0].GetContent()), &p)
	if err != nil {
		t.Fatalf("failed to unmarshall yaml: %s", err)
	}
}

func TestGenerateExtension(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
				options: {
					[grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field]: {
						description: "This is bar"
						extensions: {
							key: "x-go-default"
							value: {
								string_value: "0.5s"
							}
						}
					}
				}
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/v1/test"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format

		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, false, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			if !strings.Contains(content, "x-go-default") {
				t.Fatal("x-go-default not found in content message")
			}
		})
	}
}

func TestGenerateYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputProtoText string
		wantYAML       string
	}{
		{
			// It tests https://github.com/grpc-ecosystem/grpc-gateway/issues/3557.
			name:           "path item object",
			inputProtoText: "testdata/generator/path_item_object.prototext",
			wantYAML:       "testdata/generator/path_item_object.swagger.yaml",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := os.ReadFile(tt.inputProtoText)
			if err != nil {
				t.Fatal(err)
			}
			var req pluginpb.CodeGeneratorRequest
			if err := prototext.Unmarshal(b, &req); err != nil {
				t.Fatal(err)
			}

			resp := requireGenerate(t, &req, genopenapi.FormatYAML, false, true)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}
			got := resp[0].GetContent()

			want, err := os.ReadFile(tt.wantYAML)
			if err != nil {
				t.Fatal(err)
			}
			diff := cmp.Diff(string(want), got)
			if diff != "" {
				t.Fatalf("content not match\n%s", diff)
			}
		})
	}
}

func requireGenerate(
	tb testing.TB,
	req *pluginpb.CodeGeneratorRequest,
	format genopenapi.Format,
	preserveRPCOrder bool,
	allowMerge bool,
) []*descriptor.ResponseFile {
	tb.Helper()

	reg := descriptor.NewRegistry()
	reg.SetPreserveRPCOrder(preserveRPCOrder)
	reg.SetAllowMerge(allowMerge)

	if err := reg.Load(req); err != nil {
		tb.Fatalf("failed to load request: %s", err)
	}

	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			tb.Fatalf("failed to lookup file: %s", err)
		}

		targets = append(targets, f)
	}

	g := genopenapi.New(reg, format)

	resp, err := g.Generate(targets)
	switch {
	case err != nil:
		tb.Fatalf("failed to generate targets: %s", err)
	case len(resp) != len(targets) && !allowMerge:
		tb.Fatalf("invalid count, expected: %d, actual: %d", len(targets), len(resp))
	}

	return resp
}

func TestGeneratedYAMLIndent(t *testing.T) {
	// It tests https://github.com/grpc-ecosystem/grpc-gateway/issues/2745.
	const in = `
	file_to_generate: "exampleproto/v1/exampleproto.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/exampleproto.proto"
		package: "repro"
		message_type: {
			name: "RollupRequest"
			field: {
				name: "type"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_ENUM
				type_name: ".repro.RollupType"
				json_name: "type"
			}
		}
		message_type: {
			name: "RollupResponse"
		}
		enum_type: {
			name: "RollupType"
			value: {
				name: "UNKNOWN_ROLLUP"
				number: 0
			}
			value: {
				name: "APPLE"
				number: 1
			}
			value: {
				name: "BANANA"
				number: 2
			}
			value: {
				name: "CARROT"
				number: 3
			}
		}
		service: {
			name: "Repro"
			method: {
				name: "GetRollup"
				input_type: ".repro.RollupRequest"
				output_type: ".repro.RollupResponse"
				options: {
					[google.api.http]: {
						get: "/rollup"
					}
				}
			}
		}
		options: {
			go_package: "repro/foobar"
		}
		source_code_info: {
			location: {
				path: 5
				path: 0
				path: 2
				path: 1
				span: 24
				span: 4
				span: 14
				leading_comments: " Apples are good\n"
			}
			location: {
				path: 5
				path: 0
				path: 2
				path: 3
				span: 28
				span: 4
				span: 15
				leading_comments: " Carrots are mediocre\n"
			}
		}
		syntax: "proto3"
	}
	`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	resp := requireGenerate(t, &req, genopenapi.FormatYAML, false, false)
	if len(resp) != 1 {
		t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
	}

	content := resp[0].GetContent()

	err := yaml.Unmarshal([]byte(content), map[string]interface{}{})
	if err != nil {
		t.Log(content)
		t.Fatalf("got invalid yaml: %s", err)
	}
}

func TestGenerateRPCOrderPreserved(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, true, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/b/first", "/a/second", "/c/third"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}

}

func TestGenerateRPCOrderNotPreserved(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, false, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)
			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/a/second", "/b/first", "/c/third"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}

}

func TestGenerateRPCOrderPreservedMultipleServices(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestServiceOne"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/d/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/e/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
					}
				}
			}
		}
		service: {
			name: "TestServiceTwo"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/g/third"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, true, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/d/first", "/e/second", "/c/third", "/b/first", "/a/second", "/g/third"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOrderNotPreservedMultipleServices(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestServiceOne"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/d/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/e/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
					}
				}
			}
		}
		service: {
			name: "TestServiceTwo"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/g/third"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, false, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/d/first", "/e/second", "/c/third", "/b/first", "/a/second", "/g/third"}
			sort.Strings(expectedPaths)

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOrderPreservedMergeFiles(t *testing.T) {
	t.Parallel()

	const in1 = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/cpath"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/bpath"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/apath"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	const in2 = `
	file_to_generate: "exampleproto/v2/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v2/example.proto"
		package: "example.v2"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/f/fpath"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/e/epath"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/d/dpath"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v2;exampleproto"
		}
	}`

	var req1, req2 pluginpb.CodeGeneratorRequest

	if err := prototext.Unmarshal([]byte(in1), &req1); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}
	if err := prototext.Unmarshal([]byte(in2), &req2); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	req1.ProtoFile = append(req1.ProtoFile, req2.ProtoFile...)
	req1.FileToGenerate = append(req1.FileToGenerate, req2.FileToGenerate...)
	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req1, format, true, true)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/c/cpath", "/b/bpath", "/a/apath", "/f/fpath", "/e/epath", "/d/dpath"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOrderNotPreservedMergeFiles(t *testing.T) {
	t.Parallel()

	const in1 = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/cpath"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/bpath"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/apath"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	const in2 = `
	file_to_generate: "exampleproto/v2/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v2/example.proto"
		package: "example.v2"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/f/fpath"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/e/epath"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/d/dpath"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v2;exampleproto"
		}
	}`

	var req1, req2 pluginpb.CodeGeneratorRequest

	if err := prototext.Unmarshal([]byte(in1), &req1); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}
	if err := prototext.Unmarshal([]byte(in2), &req2); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	req1.ProtoFile = append(req1.ProtoFile, req2.ProtoFile...)
	req1.FileToGenerate = append(req1.FileToGenerate, req2.FileToGenerate...)
	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req1, format, false, true)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/c/cpath", "/b/bpath", "/a/apath", "/f/fpath", "/e/epath", "/d/dpath"}
			sort.Strings(expectedPaths)

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOrderPreservedAdditionalBindings(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
						additional_bindings {
							get: "/a/additional"
						}
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
						additional_bindings {
							get: "/z/zAdditional"
						}
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
						additional_bindings {
							get: "/b/bAdditional"
						}
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, true, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/b/first", "/a/additional", "/a/second", "/z/zAdditional", "/c/third", "/b/bAdditional"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOneOfFieldBodyAdditionalBindings(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			oneof_decl: {
				name: "foo"
			}
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
				oneof_index: 0
			}
			field: {
				name: "baz"
				number: 2
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
				oneof_index: 0
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						post: "/b/foo"
						body: "*"
						additional_bindings {
							post: "/b/foo/bar"
							body: "bar"
						}
						additional_bindings {
							post: "/b/foo/baz"
							body: "baz"
						}
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, true, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/b/foo", "/b/foo/bar", "/b/foo/baz"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}

			// The input message only contains oneof fields, so no other fields should be mapped to the query.
			if strings.Contains(content, "query") {
				t.Fatalf("Found query in content, expected not to find any")
			}
		})
	}
}

func TestGenerateRPCOrderNotPreservedAdditionalBindings(t *testing.T) {
	t.Parallel()

	const in = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
						additional_bindings {
							get: "/a/additional"
						}
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
						additional_bindings {
							get: "/z/zAdditional"
						}
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
						additional_bindings {
							get: "/b/bAdditional"
						}
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req, format, false, false)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/b/first", "/a/additional", "/a/second", "/z/zAdditional", "/c/third", "/b/bAdditional"}
			sort.Strings(expectedPaths)

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOrderPreservedMergeFilesAdditionalBindingsMultipleServices(t *testing.T) {
	t.Parallel()

	const in1 = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestServiceOne"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/d/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/e/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
					}
				}
			}
		}
		service: {
			name: "TestServiceTwo"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/g/third"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	const in2 = `
	file_to_generate: "exampleproto/v2/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v2/example.proto"
		package: "example.v2"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/b/bpath"
						additional_bindings {
							get: "/a/additional"
						}
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/a/apath"
						additional_bindings {
							get: "/z/zAdditional"
						}
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/c/cpath"
						additional_bindings {
							get: "/b/bAdditional"
						}
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v2;exampleproto"
		}
	}`

	var req1, req2 pluginpb.CodeGeneratorRequest

	if err := prototext.Unmarshal([]byte(in1), &req1); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}
	if err := prototext.Unmarshal([]byte(in2), &req2); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	req1.ProtoFile = append(req1.ProtoFile, req2.ProtoFile...)
	req1.FileToGenerate = append(req1.FileToGenerate, req2.FileToGenerate...)
	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req1, format, true, true)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/d/first", "/e/second", "/c/third",
				"/b/first", "/a/second", "/g/third", "/b/bpath", "/a/additional",
				"/a/apath", "/z/zAdditional", "/c/cpath", "/b/bAdditional"}

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

func TestGenerateRPCOrderNotPreservedMergeFilesAdditionalBindingsMultipleServices(t *testing.T) {
	t.Parallel()

	const in1 = `
	file_to_generate: "exampleproto/v1/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v1/example.proto"
		package: "example.v1"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestServiceOne"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/d/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/e/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/c/third"
					}
				}
			}
		}
		service: {
			name: "TestServiceTwo"
			method: {
				name: "Test1"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/b/first"
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/a/second"
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v1.Foo"
				output_type: ".example.v1.Foo"
				options: {
					[google.api.http]: {
						get: "/g/third"
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v1;exampleproto"
		}
	}`

	const in2 = `
	file_to_generate: "exampleproto/v2/example.proto"
	parameter: "output_format=yaml,allow_delete_body=true"
	proto_file: {
		name: "exampleproto/v2/example.proto"
		package: "example.v2"
		message_type: {
			name: "Foo"
			field: {
				name: "bar"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "bar"
			}
		}
		service: {
			name: "TestService"
			method: {
				name: "Test1"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/b/bpath"
						additional_bindings {
							get: "/a/additional"
						}
					}
				}
			}
			method: {
				name: "Test2"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/a/apath"
						additional_bindings {
							get: "/z/zAdditional"
						}
					}
				}
			}
			method: {
				name: "Test3"
				input_type: ".example.v2.Foo"
				output_type: ".example.v2.Foo"
				options: {
					[google.api.http]: {
						get: "/c/cpath"
						additional_bindings {
							get: "/b/bAdditional"
						}
					}
				}
			}
		}
		options: {
			go_package: "exampleproto/v2;exampleproto"
		}
	}`

	var req1, req2 pluginpb.CodeGeneratorRequest

	if err := prototext.Unmarshal([]byte(in1), &req1); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}
	if err := prototext.Unmarshal([]byte(in2), &req2); err != nil {
		t.Fatalf("failed to marshall yaml: %s", err)
	}

	req1.ProtoFile = append(req1.ProtoFile, req2.ProtoFile...)
	req1.FileToGenerate = append(req1.FileToGenerate, req2.FileToGenerate...)
	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			resp := requireGenerate(t, &req1, format, false, true)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()

			t.Log(content)

			contentsSlice := strings.Fields(content)
			expectedPaths := []string{"/d/first", "/e/second", "/c/third",
				"/b/first", "/a/second", "/g/third", "/b/bpath", "/a/additional",
				"/a/apath", "/z/zAdditional", "/c/cpath", "/b/bAdditional"}
			sort.Strings(expectedPaths)

			foundPaths := []string{}
			for _, contentValue := range contentsSlice {
				findExpectedPaths(&foundPaths, expectedPaths, contentValue)
			}

			if allPresent := reflect.DeepEqual(foundPaths, expectedPaths); !allPresent {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, expectedPaths)
			}
		})
	}
}

// Tries to find expected paths from a provided substring and store them in the foundPaths
// slice.
func findExpectedPaths(foundPaths *[]string, expectedPaths []string, potentialPath string) {
	seenPaths := map[string]struct{}{}

	// foundPaths may not be empty when this function is called multiple times,
	// so we add them to seenPaths map to avoid duplicates.
	for _, path := range *foundPaths {
		seenPaths[path] = struct{}{}
	}

	for _, path := range expectedPaths {
		_, pathAlreadySeen := seenPaths[path]
		if strings.Contains(potentialPath, path) && !pathAlreadySeen {
			*foundPaths = append(*foundPaths, path)
			seenPaths[path] = struct{}{}
		}
	}
}

func TestFindExpectedPaths(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		testName           string
		requiredPaths      []string
		potentialPath      string
		expectedPathsFound []string
	}{
		{
			testName:           "One potential path present",
			requiredPaths:      []string{"/d/first", "/e/second", "/c/third", "/b/first"},
			potentialPath:      "[{\"path: \"/d/first\"",
			expectedPathsFound: []string{"/d/first"},
		},
		{
			testName:           "No potential Paths present",
			requiredPaths:      []string{"/d/first", "/e/second", "/c/third", "/b/first"},
			potentialPath:      "[{\"path: \"/z/zpath\"",
			expectedPathsFound: []string{},
		},
		{
			testName:           "Multiple potential paths present",
			requiredPaths:      []string{"/d/first", "/e/second", "/c/third", "/b/first", "/d/first"},
			potentialPath:      "[{\"path: \"/d/first\"someData\"/c/third\"someData\"/b/third\"",
			expectedPathsFound: []string{"/d/first", "/c/third"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			foundPaths := []string{}
			findExpectedPaths(&foundPaths, tc.requiredPaths, tc.potentialPath)
			if correctPathsFound := reflect.DeepEqual(foundPaths, tc.expectedPathsFound); !correctPathsFound {
				t.Fatalf("Found paths differed from expected paths. Got: %#v, want %#v", foundPaths, tc.expectedPathsFound)
			}
		})
	}
}

func TestGenerateXGoType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputProtoText string
		wantYAML       string
	}{
		{
			name:           "x-go-type extension",
			inputProtoText: "testdata/generator/x_go_type.prototext",
			wantYAML:       "testdata/generator/x_go_type.swagger.yaml",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := os.ReadFile(tt.inputProtoText)
			if err != nil {
				t.Fatal(err)
			}
			var req pluginpb.CodeGeneratorRequest
			if err := prototext.Unmarshal(b, &req); err != nil {
				t.Fatal(err)
			}

			reg := descriptor.NewRegistry()
			reg.SetGenerateXGoType(true)
			if err := reg.Load(&req); err != nil {
				t.Fatalf("failed to load request: %s", err)
			}

			var targets []*descriptor.File
			for _, target := range req.FileToGenerate {
				f, err := reg.LookupFile(target)
				if err != nil {
					t.Fatalf("failed to lookup file: %s", err)
				}
				targets = append(targets, f)
			}

			g := genopenapi.New(reg, genopenapi.FormatYAML)
			resp, err := g.Generate(targets)
			if err != nil {
				t.Fatalf("failed to generate: %s", err)
			}

			if len(resp) != 1 {
				t.Fatalf("expected 1 file, got %d", len(resp))
			}

			want, err := os.ReadFile(tt.wantYAML)
			if err != nil {
				t.Fatal(err)
			}

			var gotMap, wantMap map[string]interface{}
			if err := yaml.Unmarshal([]byte(resp[0].GetContent()), &gotMap); err != nil {
				t.Fatalf("failed to unmarshal generated YAML: %v", err)
			}
			if err := yaml.Unmarshal(want, &wantMap); err != nil {
				t.Fatalf("failed to unmarshal expected YAML: %v", err)
			}

			gotYAML, err := yaml.Marshal(gotMap)
			if err != nil {
				t.Fatalf("failed to marshal got YAML: %v", err)
			}
			wantYAML, err := yaml.Marshal(wantMap)
			if err != nil {
				t.Fatalf("failed to marshal want YAML: %v", err)
			}

			if !bytes.Equal(gotYAML, wantYAML) {
				t.Errorf("YAMLs don't match:\ngot:\n%s\nwant:\n%s", gotYAML, wantYAML)
			}
		})
	}
}

// TestIssue5684_UnusedMethodsNotInOpenAPI tests that methods without HTTP bindings
// do not appear in the OpenAPI definitions.
// See https://github.com/grpc-ecosystem/grpc-gateway/issues/5684
func TestIssue5684_UnusedMethodsNotInOpenAPI(t *testing.T) {
	t.Parallel()

	// Create a proto definition similar to the issue report:
	// - Service with two methods: Add (no HTTP binding) and Show (with HTTP binding)
	// - Only Show should appear in the OpenAPI output
	// - AddRequest and AddResponse should NOT appear in definitions
	const in = `
	file_to_generate: "account/account.proto"
	proto_file: {
		name: "account/account.proto"
		package: "account"

		message_type: {
			name: "Money"
			field: {
				name: "amount"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_INT64
				json_name: "amount"
			}
		}

		message_type: {
			name: "AddRequest"
			field: {
				name: "money"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_MESSAGE
				type_name: ".account.Money"
				json_name: "money"
			}
		}

		message_type: {
			name: "AddResponse"
		}

		message_type: {
			name: "ShowRequest"
		}

		message_type: {
			name: "ShowResponse"
			field: {
				name: "money"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_MESSAGE
				type_name: ".account.Money"
				json_name: "money"
			}
		}

		service: {
			name: "AccountService"
			method: {
				name: "Add"
				input_type: ".account.AddRequest"
				output_type: ".account.AddResponse"
			}
			method: {
				name: "Show"
				input_type: ".account.ShowRequest"
				output_type: ".account.ShowResponse"
				options: {
					[google.api.http]: {
						get: "/v1/account"
					}
				}
			}
		}

		options: {
			go_package: "accounts/pkg/account;account"
		}
	}`

	var req pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(in), &req); err != nil {
		t.Fatalf("failed to unmarshal proto: %v", err)
	}

	resp := requireGenerate(t, &req, genopenapi.FormatYAML, false, false)
	if len(resp) != 1 {
		t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
	}

	var openAPIDoc map[string]interface{}
	if err := yaml.Unmarshal([]byte(resp[0].GetContent()), &openAPIDoc); err != nil {
		t.Fatalf("failed to parse OpenAPI YAML: %v", err)
	}

	definitions, ok := openAPIDoc["definitions"].(map[string]interface{})
	if !ok {
		t.Fatalf("no definitions found in OpenAPI document")
	}

	if _, exists := definitions["accountAddRequest"]; exists {
		t.Error("accountAddRequest found in definitions, but should be excluded (Add method has no HTTP binding)")
	}

	if _, exists := definitions["accountAddResponse"]; exists {
		t.Error("accountAddResponse found in definitions, but should be excluded (Add method has no HTTP binding)")
	}

	if _, exists := definitions["accountShowResponse"]; !exists {
		t.Error("accountShowResponse not found in definitions, but should be included (Show method has HTTP binding)")
	}

	if _, exists := definitions["accountMoney"]; !exists {
		t.Error("accountMoney not found in definitions, but should be included (referenced by ShowResponse)")
	}

	paths, ok := openAPIDoc["paths"].(map[string]interface{})
	if !ok {
		t.Fatalf("no paths found in OpenAPI document")
	}

	if _, exists := paths["/v1/account"]; !exists {
		t.Error("/v1/account path not found, but should be included (Show method)")
	}

	if len(paths) != 1 {
		t.Errorf("expected exactly 1 path, got %d paths", len(paths))
	}
}

// TestGenerateMergeFilesWithBodyAndPathParams tests that OpenAPI generation
// doesn't panic when merging files where a service uses body:"*" with path parameters.
// This reproduces the bug from https://github.com/grpc-ecosystem/grpc-gateway/issues/6274
func TestGenerateMergeFilesWithBodyAndPathParams(t *testing.T) {
	t.Parallel()

	// First proto file: contains only message definitions, with swagger option
	// This file will be the merge target since it has the swagger option
	const messagesProto = `
	proto_file: {
		name: "example/v1/messages.proto"
		package: "example.v1"
		message_type: {
			name: "Item"
			field: {
				name: "id"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "id"
			}
			field: {
				name: "name"
				number: 2
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "name"
			}
		}
		message_type: {
			name: "UpdateItemRequest"
			field: {
				name: "id"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "id"
			}
			field: {
				name: "name"
				number: 2
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "name"
			}
		}
		options: {
			go_package: "example/v1;examplev1"
			[grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger]: {
				info: {
					title: "Test API"
					version: "1.0"
				}
			}
		}
		syntax: "proto3"
	}`

	// Second proto file: contains the service that references messages from first file
	// This file does NOT have the swagger option, so it won't be the merge target
	const serviceProto = `
	proto_file: {
		name: "example/v1/service.proto"
		package: "example.v1"
		dependency: "example/v1/messages.proto"
		service: {
			name: "ItemService"
			method: {
				name: "UpdateItem"
				input_type: ".example.v1.UpdateItemRequest"
				output_type: ".example.v1.Item"
				options: {
					[google.api.http]: {
						put: "/v1/items/{id}"
						body: "*"
					}
				}
			}
		}
		options: {
			go_package: "example/v1;examplev1"
		}
		syntax: "proto3"
	}`

	var msgReq, svcReq pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(messagesProto), &msgReq); err != nil {
		t.Fatalf("failed to unmarshal messages proto: %s", err)
	}
	if err := prototext.Unmarshal([]byte(serviceProto), &svcReq); err != nil {
		t.Fatalf("failed to unmarshal service proto: %s", err)
	}

	// Combine into a single request with both files to generate
	req := &pluginpb.CodeGeneratorRequest{
		ProtoFile:      append(msgReq.ProtoFile, svcReq.ProtoFile...),
		FileToGenerate: []string{"example/v1/messages.proto", "example/v1/service.proto"},
	}

	formats := [...]genopenapi.Format{
		genopenapi.FormatJSON,
		genopenapi.FormatYAML,
	}

	for _, format := range formats {
		format := format
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			// This should not panic - the bug causes panic with
			// "failed to resolve method FQN: '.example.v1.ItemService.UpdateItem'"
			resp := requireGenerate(t, req, format, false, true)
			if len(resp) != 1 {
				t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
			}

			content := resp[0].GetContent()
			t.Log(content)

			// Verify the path exists in output
			if !strings.Contains(content, "/v1/items/{id}") {
				t.Error("expected /v1/items/{id} path in output")
			}

			// Verify the body definition was created (this is what triggers the bug)
			if !strings.Contains(content, "ItemServiceUpdateItemBody") {
				t.Error("expected ItemServiceUpdateItemBody definition in output")
			}
		})
	}
}

// TestGenerateMergeWithServiceNotInTargets tests the scenario where a service
// file is available in the proto file set but not marked for generation.
// This can happen with Edition 2023 files where dependencies are structured differently.
func TestGenerateMergeWithServiceNotInTargets(t *testing.T) {
	t.Parallel()

	// This tests a potential bug where services from non-target files
	// might not have their methods registered in the registry.
	// The test ensures that when generating OpenAPI, only services from
	// target files are processed.

	const messagesProto = `
	proto_file: {
		name: "example/v1/messages.proto"
		package: "example.v1"
		message_type: {
			name: "Item"
			field: {
				name: "id"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "id"
			}
		}
		message_type: {
			name: "UpdateItemRequest"
			field: {
				name: "id"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
				json_name: "id"
			}
		}
		options: {
			go_package: "example/v1;examplev1"
		}
		syntax: "proto3"
	}`

	const serviceProto = `
	proto_file: {
		name: "example/v1/service.proto"
		package: "example.v1"
		dependency: "example/v1/messages.proto"
		service: {
			name: "ItemService"
			method: {
				name: "UpdateItem"
				input_type: ".example.v1.UpdateItemRequest"
				output_type: ".example.v1.Item"
				options: {
					[google.api.http]: {
						put: "/v1/items/{id}"
						body: "*"
					}
				}
			}
		}
		options: {
			go_package: "example/v1;examplev1"
		}
		syntax: "proto3"
	}`

	var msgReq, svcReq pluginpb.CodeGeneratorRequest
	if err := prototext.Unmarshal([]byte(messagesProto), &msgReq); err != nil {
		t.Fatalf("failed to unmarshal messages proto: %s", err)
	}
	if err := prototext.Unmarshal([]byte(serviceProto), &svcReq); err != nil {
		t.Fatalf("failed to unmarshal service proto: %s", err)
	}

	// Include both files in ProtoFile but only service.proto in FileToGenerate
	req := &pluginpb.CodeGeneratorRequest{
		ProtoFile:      append(msgReq.ProtoFile, svcReq.ProtoFile...),
		FileToGenerate: []string{"example/v1/service.proto"},
	}

	resp := requireGenerate(t, req, genopenapi.FormatJSON, false, false)
	if len(resp) != 1 {
		t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
	}

	content := resp[0].GetContent()
	t.Log(content)

	// Verify the path exists in output
	if !strings.Contains(content, "/v1/items/{id}") {
		t.Error("expected /v1/items/{id} path in output")
	}

	// Verify the body definition was created
	if !strings.Contains(content, "ItemServiceUpdateItemBody") {
		t.Error("expected ItemServiceUpdateItemBody definition in output")
	}
}
