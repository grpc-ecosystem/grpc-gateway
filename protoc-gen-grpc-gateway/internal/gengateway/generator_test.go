package gengateway

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func newExampleFileDescriptor() *descriptor.File {
	return newExampleFileDescriptorWithGoPkg(
		&descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
	)
}

func newExampleFileDescriptorWithGoPkg(gp *descriptor.GoPackage) *descriptor.File {
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
		GoPkg:    *gp,
		Messages: []*descriptor.Message{msg},
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

func TestGenerateServiceWithoutBindings(t *testing.T) {
	file := newExampleFileDescriptor()
	g := &generator{}
	got, err := g.generate(crossLinkFixture(file))
	if err != nil {
		t.Errorf("generate(%#v) failed with %v; want success", file, err)
		return
	}
	if notwanted := `"github.com/golang/protobuf/ptypes/empty"`; strings.Contains(got, notwanted) {
		t.Errorf("generate(%#v) = %s; does not want to contain %s", file, got, notwanted)
	}
}

func TestGenerateOutputPath(t *testing.T) {
	cases := []struct {
		file       *descriptor.File
		pathType   pathType
		modulePath string

		// the path that function Generate should output
		expectedPath string
		// the path after protogen has remove the module prefix
		expectedFinalPath string
		expectedError     bool
	}{
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/to/example",
					Name: "example_pb",
				},
			),
			expectedPath:      "example.com/path/to/example",
			expectedFinalPath: "example.com/path/to/example",
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example",
					Name: "example_pb",
				},
			),
			expectedPath:      "example",
			expectedFinalPath: "example",
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/to/example",
					Name: "example_pb",
				},
			),
			pathType: pathTypeSourceRelative,

			expectedPath:      ".",
			expectedFinalPath: ".",
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example",
					Name: "example_pb",
				},
			),
			pathType: pathTypeSourceRelative,

			expectedPath:      ".",
			expectedFinalPath: ".",
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/root",
					Name: "example_pb",
				},
			),
			modulePath: "example.com/path/root",

			expectedPath:      "example.com/path/root",
			expectedFinalPath: ".",
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/to/example",
					Name: "example_pb",
				},
			),
			modulePath: "example.com/path/to",

			expectedPath:      "example.com/path/to/example",
			expectedFinalPath: "example",
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/to/example/with/many/nested/paths",
					Name: "example_pb",
				},
			),
			modulePath: "example.com/path/to",

			expectedPath:      "example.com/path/to/example/with/many/nested/paths",
			expectedFinalPath: "example/with/many/nested/paths",
		},

		// Error cases
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/root",
					Name: "example_pb",
				},
			),
			modulePath:    "example.com/path/root",
			pathType:      pathTypeSourceRelative, // Not allowed
			expectedError: true,
		},
		{
			file: newExampleFileDescriptorWithGoPkg(
				&descriptor.GoPackage{
					Path: "example.com/path/rootextra",
					Name: "example_pb",
				},
			),
			modulePath:    "example.com/path/root", // Not a prefix of path
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			g := &generator{
				pathType: c.pathType,
			}

			file := c.file
			gots, err := g.Generate([]*descriptor.File{crossLinkFixture(file)})

			// We don't expect an error during generation
			if err != nil {
				t.Errorf("Generate(%#v) failed with %v; wants success", file, err)
				return
			}

			if len(gots) != 1 {
				t.Errorf("Generate(%#v) failed; expects one result, got %d", file, len(gots))
				return
			}

			got := gots[0]
			if got.Name == nil {
				t.Errorf("Generate(%#v) failed; expects non-nil Name(%v)", file, got.Name)
				return
			}

			gotPath := filepath.Dir(*got.Name)
			if gotPath != c.expectedPath && !c.expectedError {
				t.Errorf("Generate(%#v) failed; got path: %s expected path: %s", file, gotPath, c.expectedPath)
				return
			}

			// We now use codegen to verify how it optionally removes the module prefix

			reqParam := ""
			if c.modulePath != "" {
				reqParam = "module=" + c.modulePath
			}
			req := &pluginpb.CodeGeneratorRequest{Parameter: &reqParam}
			plugin, err := protogen.Options{}.New(req)
			if err != nil {
				t.Errorf("Unexpected error during plugin creation: %v", err)
			}

			genFile := plugin.NewGeneratedFile(got.GetName(), protogen.GoImportPath(got.GoPkg.Path))
			_, _ = genFile.Write([]byte(got.GetContent()))
			resp := plugin.Response()

			if !c.expectedError && resp.GetError() != "" {
				t.Errorf("Unexpected error in protogen response: %v", resp.GetError())
				return
			}

			if c.expectedError {
				if resp.GetError() == "" {
					t.Error("Expected an non-null error in protogen response")
				}
				return
			}

			finalName := resp.File[0].GetName()
			gotPath = filepath.Dir(finalName)
			if gotPath != c.expectedFinalPath {
				t.Errorf("After protogen, got path: %s expected path: %s", gotPath, c.expectedFinalPath)
				return
			}
		})
	}
}
