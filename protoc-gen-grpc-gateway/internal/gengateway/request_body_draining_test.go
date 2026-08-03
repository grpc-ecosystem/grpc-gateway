package gengateway

import (
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const generatedRequestBodyDrain = "io.Copy(io.Discard, req.Body)"

func newRequestBodyDrainingFixture(
	t *testing.T,
	scenario string,
) (*descriptor.File, bool) {
	t.Helper()

	switch scenario {
	case "body_mapping":
		file := newExampleFileDescriptorWithGoPkg(
			&descriptor.GoPackage{
				Path: "example.com/path/to/example",
				Name: "example_pb",
			},
			"path/to/example",
		)
		return crossLinkFixture(file), false

	case "no_body_mapping":
		file := newExampleFileDescriptorWithGoPkg(
			&descriptor.GoPackage{
				Path: "example.com/path/to/example",
				Name: "example_pb",
			},
			"path/to/example",
		)
		file.Services[0].Methods[0].Bindings[0].Body = nil
		return crossLinkFixture(file), false

	case "patch_field_mask":
		updateMaskDesc := &descriptorpb.FieldDescriptorProto{
			Name:     proto.String("UpdateMask"),
			Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
			Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
			TypeName: proto.String(".google.protobuf.FieldMask"),
			Number:   proto.Int32(1),
		}

		messageDesc := &descriptorpb.DescriptorProto{
			Name: proto.String("ExampleMessage"),
			Field: []*descriptorpb.FieldDescriptorProto{
				updateMaskDesc,
			},
		}

		methodDesc := &descriptorpb.MethodDescriptorProto{
			Name:       proto.String("Example"),
			InputType:  proto.String("ExampleMessage"),
			OutputType: proto.String("ExampleMessage"),
		}

		serviceDesc := &descriptorpb.ServiceDescriptorProto{
			Name: proto.String("ExampleService"),
			Method: []*descriptorpb.MethodDescriptorProto{
				methodDesc,
			},
		}

		message := &descriptor.Message{
			DescriptorProto: messageDesc,
		}

		updateMaskField := &descriptor.Field{
			Message:              message,
			FieldDescriptorProto: updateMaskDesc,
		}
		message.Fields = append(message.Fields, updateMaskField)

		file := &descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				Name:    proto.String("example.proto"),
				Package: proto.String("example"),
				MessageType: []*descriptorpb.DescriptorProto{
					messageDesc,
				},
				Service: []*descriptorpb.ServiceDescriptorProto{
					serviceDesc,
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: []*descriptor.Message{
				message,
			},
			Services: []*descriptor.Service{
				{
					ServiceDescriptorProto: serviceDesc,
					Methods: []*descriptor.Method{
						{
							MethodDescriptorProto: methodDesc,
							RequestType:           message,
							ResponseType:          message,
							Bindings: []*descriptor.Binding{
								{
									HTTPMethod: "PATCH",
									Body: &descriptor.Body{
										FieldPath: descriptor.FieldPath{
											{
												Name:   "abe",
												Target: updateMaskField,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		return crossLinkFixture(file), true

	default:
		t.Fatalf("unknown scenario %q", scenario)
		return nil, false
	}
}

func TestDisableRequestBodyDraining(t *testing.T) {
	scenarios := []string{
		"body_mapping",
		"patch_field_mask",
		"no_body_mapping",
	}

	for _, scenario := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			for _, disableRequestBodyDraining := range []bool{
				false,
				true,
			} {
				file, allowPatchFeature :=
					newRequestBodyDrainingFixture(t, scenario)

				generator := New(
					descriptor.NewRegistry(),
					true,
					"Handler",
					allowPatchFeature,
					false,
					false,
					disableRequestBodyDraining,
				)

				result, err := generator.Generate(
					[]*descriptor.File{file},
				)
				if err != nil {
					t.Fatalf(
						"Generate() failed: %v",
						err,
					)
				}

				if len(result) != 1 {
					t.Fatalf(
						"Generate() returned %d files, want 1",
						len(result),
					)
				}

				hasDrain := strings.Contains(
					result[0].GetContent(),
					generatedRequestBodyDrain,
				)

				if disableRequestBodyDraining && hasDrain {
					t.Errorf(
						"generated code contains request-body drain " +
							"when disabling is enabled",
					)
				}

				if !disableRequestBodyDraining && !hasDrain {
					t.Errorf(
						"generated code does not contain the default " +
							"request-body drain",
					)
				}
			}
		})
	}
}
