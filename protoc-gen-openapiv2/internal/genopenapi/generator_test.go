package genopenapi_test

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi"
	"gopkg.in/yaml.v3"

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
				options: {}
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
