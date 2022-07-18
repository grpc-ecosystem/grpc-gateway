package genopenapi_test

import (
	"strings"
	"testing"

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

	resp := requireGenerate(t, req, genopenapi.FormatYAML)
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

			resp := requireGenerate(t, &req, format)
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

func requireGenerate(
	tb testing.TB,
	req *pluginpb.CodeGeneratorRequest,
	format genopenapi.Format,
) []*descriptor.ResponseFile {
	tb.Helper()

	reg := descriptor.NewRegistry()

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
	case len(resp) != len(targets):
		tb.Fatalf("invalid count, expected: %d, actual: %d", len(targets), len(resp))
	}

	return resp
}
