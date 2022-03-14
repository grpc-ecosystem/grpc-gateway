package genopenapi_test

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi"
	"gopkg.in/yaml.v2"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerate_YAML(t *testing.T) {
	t.Parallel()

	reg := descriptor.NewRegistry()
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

	if err := reg.Load(req); err != nil {
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
	switch {
	case err != nil:
		t.Fatalf("failed to generate targets: %s", err)
	case len(resp) != 1:
		t.Fatalf("invalid count, expected: 1, actual: %d", len(resp))
	}

	var p map[string]interface{}
	err = yaml.Unmarshal([]byte(resp[0].GetContent()), &p)
	if err != nil {
		t.Fatalf("failed to unmarshall yaml: %s", err)
	}
}
