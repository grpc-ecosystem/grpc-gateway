package gengateway

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
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

	// Test with Open Struct API (default)
	t.Run("OpenStructAPI", func(t *testing.T) {
		testGeneratorGenerate(t, false)
	})

	// Test with Opaque API
	t.Run("OpaqueAPI", func(t *testing.T) {
		testGeneratorGenerate(t, true)
	})
}

func testGeneratorGenerate(t *testing.T, useOpaqueAPI bool) {
	g := new(generator)
	g.reg = descriptor.NewRegistry()
	g.useOpaqueAPI = useOpaqueAPI
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

func TestAddBodyFieldImportsOpaqueOnly(t *testing.T) {
	reg, file := buildBodyImportTestFile(t)
	svc := file.Services[0]
	m := svc.Methods[0]

	bookField := m.RequestType.Fields[0]
	bookMessage, err := reg.LookupMsg("", ".example.sub.CreateBook")
	if err != nil {
		t.Fatalf("lookup book message: %v", err)
	}
	bookField.FieldMessage = bookMessage

	bodyPath := descriptor.FieldPath{
		{
			Name:   bookField.GetName(),
			Target: bookField,
		},
	}
	m.Bindings = []*descriptor.Binding{
		{
			HTTPMethod: "POST",
			Body: &descriptor.Body{
				FieldPath: bodyPath,
			},
		},
	}

	g := &generator{
		reg:          reg,
		useOpaqueAPI: false,
	}

	if got := g.addBodyFieldImports(file, m, map[string]bool{}); len(got) != 0 {
		t.Fatalf("expected no imports when opaque API disabled, got %v", got)
	}

	g.useOpaqueAPI = true
	imports := g.addBodyFieldImports(file, m, map[string]bool{})
	if len(imports) != 1 {
		t.Fatalf("expected 1 import when opaque API enabled, got %d", len(imports))
	}
	if imports[0].Path != bookMessage.File.GoPkg.Path {
		t.Fatalf("import path mismatch: got %q want %q", imports[0].Path, bookMessage.File.GoPkg.Path)
	}
}

func buildBodyImportTestFile(t *testing.T) (*descriptor.Registry, *descriptor.File) {
	t.Helper()

	subFile := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("sub.proto"),
		Package: proto.String("example.sub"),
		Syntax:  proto.String("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/sub;subpb"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("CreateBook"),
			},
		},
	}

	mainFile := &descriptorpb.FileDescriptorProto{
		Name:       proto.String("svc.proto"),
		Package:    proto.String("example.svc"),
		Syntax:     proto.String("proto3"),
		Dependency: []string{"sub.proto"},
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/svc;svcpb"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("CreateBookRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("book"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: proto.String(".example.sub.CreateBook"),
					},
				},
			},
			{
				Name: proto.String("CreateBookResponse"),
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("LibraryService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("CreateBook"),
						InputType:  proto.String(".example.svc.CreateBookRequest"),
						OutputType: proto.String(".example.svc.CreateBookResponse"),
					},
				},
			},
		},
	}

	req := &pluginpb.CodeGeneratorRequest{
		ProtoFile:      []*descriptorpb.FileDescriptorProto{subFile, mainFile},
		FileToGenerate: []string{"svc.proto"},
		CompilerVersion: &pluginpb.Version{
			Major: proto.Int32(3),
			Minor: proto.Int32(21),
		},
	}

	reg := descriptor.NewRegistry()
	if err := reg.Load(req); err != nil {
		t.Fatalf("registry load failed: %v", err)
	}

	file, err := reg.LookupFile("svc.proto")
	if err != nil {
		t.Fatalf("lookup svc file: %v", err)
	}
	return reg, crossLinkFixture(file)
}
