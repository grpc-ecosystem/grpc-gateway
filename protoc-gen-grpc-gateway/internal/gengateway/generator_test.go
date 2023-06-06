package gengateway

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func newExampleFileDescriptorWithGoPkg(gp *descriptor.GoPackage, filenamePrefix string) *descriptor.File {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	msg1 := &descriptor.Message{
		DescriptorProto: msgdesc,
		File: &descriptor.File{
			GoPkg: descriptor.GoPackage{
				Path: "github.com/golang/protobuf/ptypes/empty",
				Name: "emptypb",
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	meth1 := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ExampleWithoutBindings"),
		InputType:  proto.String("empty.Empty"),
		OutputType: proto.String("empty.Empty"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth, meth1},
	}
	return &descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg:                   *gp,
		GeneratedFilenamePrefix: filenamePrefix,
		Messages:                []*descriptor.Message{msg},
		Services: []*descriptor.Service{
			{
				ServiceDescriptorProto: svc,
				Methods: []*descriptor.Method{
					{
						MethodDescriptorProto: meth,
						RequestType:           msg,
						ResponseType:          msg,
						Bindings: []*descriptor.Binding{
							{
								HTTPMethod: "GET",
								Body:       &descriptor.Body{FieldPath: nil},
							},
						},
					},
					{
						MethodDescriptorProto: meth1,
						RequestType:           msg1,
						ResponseType:          msg1,
					},
				},
			},
		},
	}
}

func TestGenerator_Generate(t *testing.T) {
	g := new(generator)
	g.reg = descriptor.NewRegistry()
	result, err := g.Generate([]*descriptor.File{
		crossLinkFixture(newExampleFileDescriptorWithGoPkg(&descriptor.GoPackage{
			Path: "example.com/path/to/example",
			Name: "example_pb",
		}, "path/to/example")),
	})
	if err != nil {
		t.Fatalf("failed to generate stubs: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected to generate one file, got: %d", len(result))
	}
	expectedName := "path/to/example.pb.gw.go"
	gotName := result[0].GetName()
	if gotName != expectedName {
		t.Fatalf("invalid name %q, expected %q", gotName, expectedName)
	}
}

func TestGenerator_GenerateSeparatePackage(t *testing.T) {
	reg := descriptor.NewRegistry()
	reg.SetSeparatePackage(true)
	reg.SetStandalone(true)
	g := New(reg, true, "Handler", true, true, true)
	targets := []*descriptor.File{
		crossLinkFixture(newExampleFileDescriptorWithGoPkg(&descriptor.GoPackage{
			Path:  "example.com/path/to/example",
			Name:  "example" + "gateway",
			Alias: "extexample",
		}, "path/to/example")),
	}
	// Set ForcePrefixedName (usually set when standalone=true).
	for _, f := range targets {
		for _, msg := range f.Messages {
			msg.ForcePrefixedName = true
			for _, field := range msg.Fields {
				field.ForcePrefixedName = true
			}
		}
		for _, enum := range f.Enums {
			enum.ForcePrefixedName = true
		}
		imports := make([]string, 0)
		for _, svc := range f.Services {
			packageName := strings.TrimSuffix(svc.File.GoPkg.Name, "gateway") + "grpc"
			svc.ForcePrefixedName = true
			// replicates behavior in internal/descriptor/services.go (loadServices)
			svc.GRPCFile = &descriptor.File{
				GoPkg: descriptor.GoPackage{
					Path: strings.Replace(
						filepath.Join(svc.File.GoPkg.Path, packageName),
						"protocolbuffers/go",
						"grpc/go",
						1,
					),
					Name: strings.Replace(packageName, "protocolbuffers/go", "grpc/go", 1),
				},
			}
			imports = append(reg.GetAdditionalImports(), svc.GRPCFile.GoPkg.Path)
		}
		reg.SetAdditionalImports(imports)
	}
	result, err := g.Generate(targets)
	if err != nil {
		t.Fatalf("failed to generate stubs: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected to generate one file, got: %d", len(result))
	}
	expectedName := "path/to/example/examplegateway/example.pb.gw.go"
	gotName := result[0].GetName()
	if gotName != expectedName {
		t.Fatalf("invalid name %q, expected %q", gotName, expectedName)
	}
	if result[0].GoPkg.Path != "example.com/path/to/example/examplegateway" {
		t.Fatalf("invalid path %q, expected %q", result[0].GoPkg.Path, "example.com/path/to/example/examplegateway")
	}
	if result[0].GoPkg.Name != "examplegateway" {
		t.Fatalf("invalid name %q, expected %q", result[0].GoPkg.Name, "examplegateway")
	}
	// Require the two dependencies to be declared as imported packages
	for _, expectedImport := range []string{
		`extexample "example.com/path/to/example"`,
		`"example.com/path/to/example/examplegrpc"`,
	} {
		if !strings.Contains(result[0].GetContent(), expectedImport) {
			t.Fatalf("expected to find import %q in the generated file", expectedImport)
		}
	}
}
