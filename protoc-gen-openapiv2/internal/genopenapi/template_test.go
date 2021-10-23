package genopenapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor/openapiconfig"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/httprule"
	openapi_options "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/pluginpb"
)

var marshaler = &runtime.JSONPb{}

func crossLinkFixture(f *descriptor.File) *descriptor.File {
	for _, m := range f.Messages {
		m.File = f
	}
	for _, svc := range f.Services {
		svc.File = f
		for _, m := range svc.Methods {
			m.Service = svc
			for _, b := range m.Bindings {
				b.Method = m
				for _, param := range b.PathParams {
					param.Method = m
				}
			}
		}
	}
	return f
}

func reqFromFile(f *descriptor.File) *pluginpb.CodeGeneratorRequest {
	return &pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			f.FileDescriptorProto,
		},
		FileToGenerate: []string{f.GetName()},
	}
}

func TestMessageToQueryParametersWithEnumAsInt(t *testing.T) {
	type test struct {
		MsgDescs []*descriptorpb.DescriptorProto
		Message  string
		Params   []openapiParameterObject
	}

	tests := []test{
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("a"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(1),
						},
						{
							Name:   proto.String("b"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum(),
							Number: proto.Int32(2),
						},
						{
							Name:   proto.String("c"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
							Number: proto.Int32(3),
						},
					},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "a",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "b",
					In:       "query",
					Required: false,
					Type:     "number",
					Format:   "double",
				},
				{
					Name:             "c",
					In:               "query",
					Required:         false,
					Type:             "array",
					CollectionFormat: "multi",
				},
			},
		},
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("nested"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.Nested"),
							Number:   proto.Int32(1),
						},
					},
				},
				{
					Name: proto.String("Nested"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("a"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(1),
						},
						{
							Name:     proto.String("deep"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.Nested.DeepNested"),
							Number:   proto.Int32(2),
						},
					},
					NestedType: []*descriptorpb.DescriptorProto{{
						Name: proto.String("DeepNested"),
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:   proto.String("b"),
								Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
								Number: proto.Int32(1),
							},
							{
								Name:     proto.String("c"),
								Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
								TypeName: proto.String(".example.Nested.DeepNested.DeepEnum"),
								Number:   proto.Int32(2),
							},
						},
						EnumType: []*descriptorpb.EnumDescriptorProto{
							{
								Name: proto.String("DeepEnum"),
								Value: []*descriptorpb.EnumValueDescriptorProto{
									{Name: proto.String("FALSE"), Number: proto.Int32(0)},
									{Name: proto.String("TRUE"), Number: proto.Int32(1)},
								},
							},
						},
					}},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "nested.a",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "nested.deep.b",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "nested.deep.c",
					In:       "query",
					Required: false,
					Type:     "integer",
					Enum:     []string{"0", "1"},
					Default:  "0",
				},
			},
		},
	}

	for _, test := range tests {
		reg := descriptor.NewRegistry()
		reg.SetEnumsAsInts(true)
		msgs := []*descriptor.Message{}
		for _, msgdesc := range test.MsgDescs {
			msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
		}
		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				Dependency:     []string{},
				MessageType:    test.MsgDescs,
				Service:        []*descriptorpb.ServiceDescriptorProto{},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: msgs,
		}
		err := reg.Load(&pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
		})
		if err != nil {
			t.Fatalf("failed to load code generator request: %v", err)
		}

		message, err := reg.LookupMsg("", ".example."+test.Message)
		if err != nil {
			t.Fatalf("failed to lookup message: %s", err)
		}
		params, err := messageToQueryParameters(message, reg, []descriptor.Parameter{}, nil)
		if err != nil {
			t.Fatalf("failed to convert message to query parameters: %s", err)
		}
		// avoid checking Items for array types
		for i := range params {
			params[i].Items = nil
		}
		if !reflect.DeepEqual(params, test.Params) {
			t.Errorf("expected %v, got %v", test.Params, params)
		}
	}
}

func TestMessageToQueryParameters(t *testing.T) {
	type test struct {
		MsgDescs []*descriptorpb.DescriptorProto
		Message  string
		Params   []openapiParameterObject
	}

	tests := []test{
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("a"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(1),
						},
						{
							Name:   proto.String("b"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum(),
							Number: proto.Int32(2),
						},
						{
							Name:   proto.String("c"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
							Number: proto.Int32(3),
						},
					},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "a",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "b",
					In:       "query",
					Required: false,
					Type:     "number",
					Format:   "double",
				},
				{
					Name:             "c",
					In:               "query",
					Required:         false,
					Type:             "array",
					CollectionFormat: "multi",
				},
			},
		},
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("nested"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.Nested"),
							Number:   proto.Int32(1),
						},
					},
				},
				{
					Name: proto.String("Nested"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("a"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(1),
						},
						{
							Name:     proto.String("deep"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.Nested.DeepNested"),
							Number:   proto.Int32(2),
						},
					},
					NestedType: []*descriptorpb.DescriptorProto{{
						Name: proto.String("DeepNested"),
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:   proto.String("b"),
								Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
								Number: proto.Int32(1),
							},
							{
								Name:     proto.String("c"),
								Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
								TypeName: proto.String(".example.Nested.DeepNested.DeepEnum"),
								Number:   proto.Int32(2),
							},
						},
						EnumType: []*descriptorpb.EnumDescriptorProto{
							{
								Name: proto.String("DeepEnum"),
								Value: []*descriptorpb.EnumValueDescriptorProto{
									{Name: proto.String("FALSE"), Number: proto.Int32(0)},
									{Name: proto.String("TRUE"), Number: proto.Int32(1)},
								},
							},
						},
					}},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "nested.a",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "nested.deep.b",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "nested.deep.c",
					In:       "query",
					Required: false,
					Type:     "string",
					Enum:     []string{"FALSE", "TRUE"},
					Default:  "FALSE",
				},
			},
		},
	}

	for _, test := range tests {
		reg := descriptor.NewRegistry()
		msgs := []*descriptor.Message{}
		for _, msgdesc := range test.MsgDescs {
			msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
		}
		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				Dependency:     []string{},
				MessageType:    test.MsgDescs,
				Service:        []*descriptorpb.ServiceDescriptorProto{},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: msgs,
		}
		err := reg.Load(&pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
		})
		if err != nil {
			t.Fatalf("failed to load code generator request: %v", err)
		}

		message, err := reg.LookupMsg("", ".example."+test.Message)
		if err != nil {
			t.Fatalf("failed to lookup message: %s", err)
		}
		params, err := messageToQueryParameters(message, reg, []descriptor.Parameter{}, nil)
		if err != nil {
			t.Fatalf("failed to convert message to query parameters: %s", err)
		}
		// avoid checking Items for array types
		for i := range params {
			params[i].Items = nil
		}
		if !reflect.DeepEqual(params, test.Params) {
			t.Errorf("expected %v, got %v", test.Params, params)
		}
	}
}

// TestMessagetoQueryParametersNoRecursive, is a check that cyclical references between messages
//  are not falsely detected given previous known edge-cases.
func TestMessageToQueryParametersNoRecursive(t *testing.T) {
	type test struct {
		MsgDescs []*descriptorpb.DescriptorProto
		Message  string
	}

	tests := []test{
		// First test:
		// Here is a message that has two of another message adjacent to one another in a nested message.
		// There is no loop but this was previouly falsely flagged as a cycle.
		// Example proto:
		// message NonRecursiveMessage {
		//      string field = 1;
		// }
		// message BaseMessage {
		//      NonRecursiveMessage first = 1;
		//      NonRecursiveMessage second = 2;
		// }
		// message QueryMessage {
		//      BaseMessage first = 1;
		//      string second = 2;
		// }
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("QueryMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("first"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.BaseMessage"),
							Number:   proto.Int32(1),
						},
						{
							Name:   proto.String("second"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(2),
						},
					},
				},
				{
					Name: proto.String("BaseMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("first"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.NonRecursiveMessage"),
							Number:   proto.Int32(1),
						},
						{
							Name:     proto.String("second"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.NonRecursiveMessage"),
							Number:   proto.Int32(2),
						},
					},
				},
				// Note there is no recursive nature to this message
				{
					Name: proto.String("NonRecursiveMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name: proto.String("field"),
							//Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(1),
						},
					},
				},
			},
			Message: "QueryMessage",
		},
	}

	for _, test := range tests {
		reg := descriptor.NewRegistry()
		msgs := []*descriptor.Message{}
		for _, msgdesc := range test.MsgDescs {
			msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
		}
		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				Dependency:     []string{},
				MessageType:    test.MsgDescs,
				Service:        []*descriptorpb.ServiceDescriptorProto{},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: msgs,
		}
		err := reg.Load(&pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
		})
		if err != nil {
			t.Fatalf("failed to load code generator request: %v", err)
		}

		message, err := reg.LookupMsg("", ".example."+test.Message)
		if err != nil {
			t.Fatalf("failed to lookup message: %s", err)
		}

		_, err = messageToQueryParameters(message, reg, []descriptor.Parameter{}, nil)
		if err != nil {
			t.Fatalf("No recursion error should be thrown: %s", err)
		}
	}
}

// TestMessagetoQueryParametersRecursive, is a check that cyclical references between messages
//  are handled gracefully. The goal is to insure that attempts to add messages with cyclical
//  references to query-parameters returns an error message.
func TestMessageToQueryParametersRecursive(t *testing.T) {
	type test struct {
		MsgDescs []*descriptorpb.DescriptorProto
		Message  string
	}

	tests := []test{
		// First test:
		// Here we test that a message that references it self through a field will return an error.
		// Example proto:
		// message DirectRecursiveMessage {
		//      DirectRecursiveMessage nested = 1;
		// }
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("DirectRecursiveMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("nested"),
							Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.DirectRecursiveMessage"),
							Number:   proto.Int32(1),
						},
					},
				},
			},
			Message: "DirectRecursiveMessage",
		},
		// Second test:
		// Here we test that a cycle through multiple messages is detected and that an error is returned.
		// Sample:
		// message Root { NodeMessage nested = 1; }
		// message NodeMessage { CycleMessage nested = 1; }
		// message CycleMessage { Root nested = 1; }
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("RootMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("nested"),
							Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.NodeMessage"),
							Number:   proto.Int32(1),
						},
					},
				},
				{
					Name: proto.String("NodeMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("nested"),
							Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.CycleMessage"),
							Number:   proto.Int32(1),
						},
					},
				},
				{
					Name: proto.String("CycleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("nested"),
							Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.RootMessage"),
							Number:   proto.Int32(1),
						},
					},
				},
			},
			Message: "RootMessage",
		},
	}

	for _, test := range tests {
		reg := descriptor.NewRegistry()
		msgs := []*descriptor.Message{}
		for _, msgdesc := range test.MsgDescs {
			msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
		}
		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				Dependency:     []string{},
				MessageType:    test.MsgDescs,
				Service:        []*descriptorpb.ServiceDescriptorProto{},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: msgs,
		}
		err := reg.Load(&pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
		})
		if err != nil {
			t.Fatalf("failed to load code generator request: %v", err)
		}

		message, err := reg.LookupMsg("", ".example."+test.Message)
		if err != nil {
			t.Fatalf("failed to lookup message: %s", err)
		}
		_, err = messageToQueryParameters(message, reg, []descriptor.Parameter{}, nil)
		if err == nil {
			t.Fatalf("It should not be allowed to have recursive query parameters")
		}
	}
}

func TestMessageToQueryParametersWithJsonName(t *testing.T) {
	type test struct {
		MsgDescs []*descriptorpb.DescriptorProto
		Message  string
		Params   []openapiParameterObject
	}

	tests := []test{
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("test_field_a"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number:   proto.Int32(1),
							JsonName: proto.String("testFieldA"),
						},
					},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "testFieldA",
					In:       "query",
					Required: false,
					Type:     "string",
				},
			},
		},
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("SubMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("test_field_a"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number:   proto.Int32(1),
							JsonName: proto.String("testFieldA"),
						},
					},
				},
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("sub_message"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".example.SubMessage"),
							Number:   proto.Int32(1),
							JsonName: proto.String("subMessage"),
						},
					},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "subMessage.testFieldA",
					In:       "query",
					Required: false,
					Type:     "string",
				},
			},
		},
	}

	for _, test := range tests {
		reg := descriptor.NewRegistry()
		reg.SetUseJSONNamesForFields(true)
		msgs := []*descriptor.Message{}
		for _, msgdesc := range test.MsgDescs {
			msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
		}
		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				Dependency:     []string{},
				MessageType:    test.MsgDescs,
				Service:        []*descriptorpb.ServiceDescriptorProto{},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: msgs,
		}
		err := reg.Load(&pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
		})
		if err != nil {
			t.Fatalf("failed to load code generator request: %v", err)
		}

		message, err := reg.LookupMsg("", ".example."+test.Message)
		if err != nil {
			t.Fatalf("failed to lookup message: %s", err)
		}
		params, err := messageToQueryParameters(message, reg, []descriptor.Parameter{}, nil)
		if err != nil {
			t.Fatalf("failed to convert message to query parameters: %s", err)
		}
		if !reflect.DeepEqual(params, test.Params) {
			t.Errorf("expected %v, got %v", test.Params, params)
		}
	}
}

func TestMessageToQueryParametersWellKnownTypes(t *testing.T) {
	type test struct {
		MsgDescs          []*descriptorpb.DescriptorProto
		WellKnownMsgDescs []*descriptorpb.DescriptorProto
		Message           string
		Params            []openapiParameterObject
	}

	tests := []test{
		{
			MsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("ExampleMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("a_field_mask"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".google.protobuf.FieldMask"),
							Number:   proto.Int32(1),
						},
						{
							Name:     proto.String("a_timestamp"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".google.protobuf.Timestamp"),
							Number:   proto.Int32(2),
						},
					},
				},
			},
			WellKnownMsgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("FieldMask"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("paths"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
							Number: proto.Int32(1),
						},
					},
				},
				{
					Name: proto.String("Timestamp"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("seconds"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum(),
							Number: proto.Int32(1),
						},
						{
							Name:   proto.String("nanos"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
							Number: proto.Int32(2),
						},
					},
				},
			},
			Message: "ExampleMessage",
			Params: []openapiParameterObject{
				{
					Name:     "a_field_mask",
					In:       "query",
					Required: false,
					Type:     "string",
				},
				{
					Name:     "a_timestamp",
					In:       "query",
					Required: false,
					Type:     "string",
					Format:   "date-time",
				},
			},
		},
	}

	for _, test := range tests {
		reg := descriptor.NewRegistry()
		reg.SetEnumsAsInts(true)
		err := reg.Load(&pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{
				{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
					Name:           proto.String("google/well_known.proto"),
					Package:        proto.String("google.protobuf"),
					Dependency:     []string{},
					MessageType:    test.WellKnownMsgDescs,
					Service:        []*descriptorpb.ServiceDescriptorProto{},
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("google/well_known"),
					},
				},
				{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
					Name:           proto.String("acme/example.proto"),
					Package:        proto.String("example"),
					Dependency:     []string{"google/well_known.proto"},
					MessageType:    test.MsgDescs,
					Service:        []*descriptorpb.ServiceDescriptorProto{},
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("acme/example"),
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("failed to load CodeGeneratorRequest: %v", err)
		}

		message, err := reg.LookupMsg("", ".example."+test.Message)
		if err != nil {
			t.Fatalf("failed to lookup message: %s", err)
		}
		params, err := messageToQueryParameters(message, reg, []descriptor.Parameter{}, nil)
		if err != nil {
			t.Fatalf("failed to convert message to query parameters: %s", err)
		}
		if !reflect.DeepEqual(params, test.Params) {
			t.Errorf("expected %v, got %v", test.Params, params)
		}
	}
}

func TestApplyTemplateSimple(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			MessageType:    []*descriptorpb.DescriptorProto{msgdesc},
			Service:        []*descriptorpb.ServiceDescriptorProto{svc},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
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
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/echo", // TODO(achew22): Figure out what this should really be
								},
							},
						},
					},
				},
			},
		},
	}
	reg := descriptor.NewRegistry()
	if err := AddErrorDefs(reg); err != nil {
		t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
		return
	}
	fileCL := crossLinkFixture(&file)
	err := reg.Load(reqFromFile(fileCL))
	if err != nil {
		t.Errorf("reg.Load(%#v) failed with %v; want success", file, err)
		return
	}
	result, err := applyTemplate(param{File: fileCL, reg: reg})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want, is, name := "2.0", result.Swagger, "Swagger"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := "", result.BasePath, "BasePath"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := ([]string)(nil), result.Schemes, "Schemes"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"application/json"}, result.Consumes, "Consumes"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"application/json"}, result.Produces, "Produces"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", fmt.Sprint(result))
	}
}

func TestApplyTemplateMultiService(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}

	// Create two services that have the same method name. We will test that the
	// operation IDs are different
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	svc2 := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("OtherService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}

	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			MessageType:    []*descriptorpb.DescriptorProto{msgdesc},
			Service:        []*descriptorpb.ServiceDescriptorProto{svc},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
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
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/echo",
								},
							},
						},
					},
				},
			},
			{
				ServiceDescriptorProto: svc2,
				Methods: []*descriptor.Method{
					{
						MethodDescriptorProto: meth,
						RequestType:           msg,
						ResponseType:          msg,
						Bindings: []*descriptor.Binding{
							{
								HTTPMethod: "GET",
								Body:       &descriptor.Body{FieldPath: nil},
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/ping",
								},
							},
						},
					},
				},
			},
		},
	}
	reg := descriptor.NewRegistry()
	if err := AddErrorDefs(reg); err != nil {
		t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
		return
	}
	fileCL := crossLinkFixture(&file)
	err := reg.Load(reqFromFile(fileCL))
	if err != nil {
		t.Errorf("reg.Load(%#v) failed with %v; want success", file, err)
		return
	}
	result, err := applyTemplate(param{File: fileCL, reg: reg})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}

	// Check that the two services have unique operation IDs even though they
	// have the same method name.
	if want, is := "ExampleService_Example", result.Paths["/v1/echo"].Get.OperationID; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).Paths[0].Get.OperationID = %s want to be %s", file, is, want)
	}
	if want, is := "OtherService_Example", result.Paths["/v1/ping"].Get.OperationID; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(%#v).Paths[0].Get.OperationID = %s want to be %s", file, is, want)
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", fmt.Sprint(result))
	}
}

func TestApplyTemplateOverrideOperationID(t *testing.T) {
	newFile := func() *descriptor.File {
		msgdesc := &descriptorpb.DescriptorProto{
			Name: proto.String("ExampleMessage"),
		}
		meth := &descriptorpb.MethodDescriptorProto{
			Name:       proto.String("Example"),
			InputType:  proto.String("ExampleMessage"),
			OutputType: proto.String("ExampleMessage"),
			Options:    &descriptorpb.MethodOptions{},
		}
		svc := &descriptorpb.ServiceDescriptorProto{
			Name:   proto.String("ExampleService"),
			Method: []*descriptorpb.MethodDescriptorProto{meth},
		}
		msg := &descriptor.Message{
			DescriptorProto: msgdesc,
		}
		return &descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				MessageType:    []*descriptorpb.DescriptorProto{msgdesc},
				Service:        []*descriptorpb.ServiceDescriptorProto{svc},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
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
									PathTmpl: httprule.Template{
										Version:  1,
										OpCodes:  []int{0, 0},
										Template: "/v1/echo", // TODO(achew22): Figure out what this should really be
									},
								},
							},
						},
					},
				},
			},
		}
	}

	verifyTemplateFromReq := func(t *testing.T, reg *descriptor.Registry, file *descriptor.File, opts *openapiconfig.OpenAPIOptions) {
		if err := AddErrorDefs(reg); err != nil {
			t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
			return
		}
		fileCL := crossLinkFixture(file)
		err := reg.Load(reqFromFile(fileCL))
		if err != nil {
			t.Errorf("reg.Load(%#v) failed with %v; want success", *file, err)
			return
		}
		if opts != nil {
			if err := reg.RegisterOpenAPIOptions(opts); err != nil {
				t.Fatalf("failed to register OpenAPI options: %s", err)
			}
		}
		result, err := applyTemplate(param{File: fileCL, reg: reg})
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", *file, err)
			return
		}
		if want, is := "MyExample", result.Paths["/v1/echo"].Get.OperationID; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).Paths[0].Get.OperationID = %s want to be %s", *file, is, want)
		}

		// If there was a failure, print out the input and the json result for debugging.
		if t.Failed() {
			t.Errorf("had: %s", *file)
			t.Errorf("got: %s", fmt.Sprint(result))
		}
	}

	openapiOperation := openapi_options.Operation{
		OperationId: "MyExample",
	}

	t.Run("verify override via method option", func(t *testing.T) {
		file := newFile()
		proto.SetExtension(proto.Message(file.Services[0].Methods[0].MethodDescriptorProto.Options),
			openapi_options.E_Openapiv2Operation, &openapiOperation)

		reg := descriptor.NewRegistry()
		verifyTemplateFromReq(t, reg, file, nil)
	})

	t.Run("verify override options annotations", func(t *testing.T) {
		file := newFile()
		reg := descriptor.NewRegistry()
		opts := &openapiconfig.OpenAPIOptions{
			Method: []*openapiconfig.OpenAPIMethodOption{
				{
					Method: "example.ExampleService.Example",
					Option: &openapiOperation,
				},
			},
		}
		verifyTemplateFromReq(t, reg, file, opts)
	})
}

func TestApplyTemplateExtensions(t *testing.T) {
	newFile := func() *descriptor.File {
		msgdesc := &descriptorpb.DescriptorProto{
			Name: proto.String("ExampleMessage"),
		}
		meth := &descriptorpb.MethodDescriptorProto{
			Name:       proto.String("Example"),
			InputType:  proto.String("ExampleMessage"),
			OutputType: proto.String("ExampleMessage"),
			Options:    &descriptorpb.MethodOptions{},
		}
		svc := &descriptorpb.ServiceDescriptorProto{
			Name:   proto.String("ExampleService"),
			Method: []*descriptorpb.MethodDescriptorProto{meth},
		}
		msg := &descriptor.Message{
			DescriptorProto: msgdesc,
		}
		return &descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				MessageType:    []*descriptorpb.DescriptorProto{msgdesc},
				Service:        []*descriptorpb.ServiceDescriptorProto{svc},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
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
									PathTmpl: httprule.Template{
										Version:  1,
										OpCodes:  []int{0, 0},
										Template: "/v1/echo", // TODO(achew22): Figure out what this should really be
									},
								},
							},
						},
					},
				},
			},
		}
	}
	swagger := openapi_options.Swagger{
		Info: &openapi_options.Info{
			Title: "test",
			Extensions: map[string]*structpb.Value{
				"x-info-extension": {Kind: &structpb.Value_StringValue{StringValue: "bar"}},
			},
		},
		Extensions: map[string]*structpb.Value{
			"x-foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}},
			"x-bar": {Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{
				Values: []*structpb.Value{{Kind: &structpb.Value_StringValue{StringValue: "baz"}}},
			}}},
		},
		SecurityDefinitions: &openapi_options.SecurityDefinitions{
			Security: map[string]*openapi_options.SecurityScheme{
				"somescheme": {
					Extensions: map[string]*structpb.Value{
						"x-security-baz": {Kind: &structpb.Value_BoolValue{BoolValue: true}},
					},
				},
			},
		},
	}
	openapiOperation := openapi_options.Operation{
		Responses: map[string]*openapi_options.Response{
			"200": {
				Extensions: map[string]*structpb.Value{
					"x-resp-id": {Kind: &structpb.Value_StringValue{StringValue: "resp1000"}},
				},
			},
		},
		Extensions: map[string]*structpb.Value{
			"x-op-foo": {Kind: &structpb.Value_StringValue{StringValue: "baz"}},
		},
	}
	verifyTemplateExtensions := func(t *testing.T, reg *descriptor.Registry, file *descriptor.File,
		opts *openapiconfig.OpenAPIOptions) {
		if err := AddErrorDefs(reg); err != nil {
			t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
			return
		}
		fileCL := crossLinkFixture(file)
		err := reg.Load(reqFromFile(fileCL))
		if err != nil {
			t.Errorf("reg.Load(%#v) failed with %v; want success", file, err)
			return
		}
		if opts != nil {
			if err := reg.RegisterOpenAPIOptions(opts); err != nil {
				t.Fatalf("failed to register OpenAPI annotations: %s", err)
			}
		}
		result, err := applyTemplate(param{File: fileCL, reg: reg})
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}
		if want, is, name := "2.0", result.Swagger, "Swagger"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}
		if got, want := len(result.extensions), 2; got != want {
			t.Fatalf("len(applyTemplate(%#v).Extensions) = %d want to be %d", file, got, want)
		}
		if got, want := result.extensions[0].key, "x-bar"; got != want {
			t.Errorf("applyTemplate(%#v).Extensions[0].key = %s want to be %s", file, got, want)
		}
		if got, want := result.extensions[1].key, "x-foo"; got != want {
			t.Errorf("applyTemplate(%#v).Extensions[1].key = %s want to be %s", file, got, want)
		}
		{
			var got []string
			err = marshaler.Unmarshal(result.extensions[0].value, &got)
			if err != nil {
				t.Fatalf("marshaler.Unmarshal failed: %v", err)
			}
			want := []string{"baz"}
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf(diff)
			}
		}
		{
			var got string
			err = marshaler.Unmarshal(result.extensions[1].value, &got)
			if err != nil {
				t.Fatalf("marshaler.Unmarshal failed: %v", err)
			}
			want := "bar"
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf(diff)
			}
		}

		var scheme openapiSecuritySchemeObject
		for _, v := range result.SecurityDefinitions {
			scheme = v
		}
		if want, is, name := []extension{
			{key: "x-security-baz", value: json.RawMessage("true")},
		}, scheme.extensions, "SecurityScheme.Extensions"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}

		if want, is, name := []extension{
			{key: "x-info-extension", value: json.RawMessage("\"bar\"")},
		}, result.Info.extensions, "Info.Extensions"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}

		var operation *openapiOperationObject
		var response openapiResponseObject
		for _, v := range result.Paths {
			operation = v.Get
			response = v.Get.Responses["200"]
		}
		if want, is, name := []extension{
			{key: "x-op-foo", value: json.RawMessage("\"baz\"")},
		}, operation.extensions, "operation.Extensions"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}
		if want, is, name := []extension{
			{key: "x-resp-id", value: json.RawMessage("\"resp1000\"")},
		}, response.extensions, "response.Extensions"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}
	}
	t.Run("verify template options set via proto options", func(t *testing.T) {
		file := newFile()
		proto.SetExtension(proto.Message(file.FileDescriptorProto.Options), openapi_options.E_Openapiv2Swagger, &swagger)
		proto.SetExtension(proto.Message(file.Services[0].Methods[0].Options), openapi_options.E_Openapiv2Operation, &openapiOperation)
		reg := descriptor.NewRegistry()
		verifyTemplateExtensions(t, reg, file, nil)
	})
	t.Run("verify template options set via annotations", func(t *testing.T) {
		file := newFile()
		opts := &openapiconfig.OpenAPIOptions{
			File: []*openapiconfig.OpenAPIFileOption{
				{
					File:   "example.proto",
					Option: &swagger,
				},
			},
			Method: []*openapiconfig.OpenAPIMethodOption{
				{
					Method: "example.ExampleService.Example",
					Option: &openapiOperation,
				},
			},
		}
		reg := descriptor.NewRegistry()
		verifyTemplateExtensions(t, reg, file, opts)
	})
}

func TestApplyTemplateHeaders(t *testing.T) {
	newFile := func() *descriptor.File {
		msgdesc := &descriptorpb.DescriptorProto{
			Name: proto.String("ExampleMessage"),
		}
		meth := &descriptorpb.MethodDescriptorProto{
			Name:       proto.String("Example"),
			InputType:  proto.String("ExampleMessage"),
			OutputType: proto.String("ExampleMessage"),
			Options:    &descriptorpb.MethodOptions{},
		}
		svc := &descriptorpb.ServiceDescriptorProto{
			Name:   proto.String("ExampleService"),
			Method: []*descriptorpb.MethodDescriptorProto{meth},
		}
		msg := &descriptor.Message{
			DescriptorProto: msgdesc,
		}
		return &descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("example.proto"),
				Package:        proto.String("example"),
				MessageType:    []*descriptorpb.DescriptorProto{msgdesc},
				Service:        []*descriptorpb.ServiceDescriptorProto{svc},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
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
									PathTmpl: httprule.Template{
										Version:  1,
										OpCodes:  []int{0, 0},
										Template: "/v1/echo", // TODO(achew22): Figure out what this should really be
									},
								},
							},
						},
					},
				},
			},
		}
	}
	openapiOperation := openapi_options.Operation{
		Responses: map[string]*openapi_options.Response{
			"200": &openapi_options.Response{
				Description: "Testing Headers",
				Headers: map[string]*openapi_options.Header{
					"string": {
						Description: "string header description",
						Type:        "string",
						Format:      "uuid",
						Pattern:     "",
					},
					"boolean": {
						Description: "boolean header description",
						Type:        "boolean",
						Default:     "true",
						Pattern:     "^true|false$",
					},
					"integer": {
						Description: "integer header description",
						Type:        "integer",
						Default:     "0",
						Pattern:     "^[0-9]$",
					},
					"number": {
						Description: "number header description",
						Type:        "number",
						Default:     "1.2",
						Pattern:     "^[-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?$",
					},
				},
			},
		},
	}
	verifyTemplateHeaders := func(t *testing.T, reg *descriptor.Registry, file *descriptor.File,
		opts *openapiconfig.OpenAPIOptions) {
		if err := AddErrorDefs(reg); err != nil {
			t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
			return
		}
		fileCL := crossLinkFixture(file)
		err := reg.Load(reqFromFile(fileCL))
		if err != nil {
			t.Errorf("reg.Load(%#v) failed with %v; want success", file, err)
			return
		}
		if opts != nil {
			if err := reg.RegisterOpenAPIOptions(opts); err != nil {
				t.Fatalf("failed to register OpenAPI annotations: %s", err)
			}
		}
		result, err := applyTemplate(param{File: fileCL, reg: reg})
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}
		if want, is, name := "2.0", result.Swagger, "Swagger"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}

		var response openapiResponseObject
		for _, v := range result.Paths {
			response = v.Get.Responses["200"]
		}
		if want, is, name := []openapiHeadersObject{
			{
				"String": openapiHeaderObject{
					Description: "string header description",
					Type:        "string",
					Format:      "uuid",
					Pattern:     "",
				},
				"Boolean": openapiHeaderObject{
					Description: "boolean header description",
					Type:        "boolean",
					Default:     json.RawMessage("true"),
					Pattern:     "^true|false$",
				},
				"Integer": openapiHeaderObject{
					Description: "integer header description",
					Type:        "integer",
					Default:     json.RawMessage("0"),
					Pattern:     "^[0-9]$",
				},
				"Number": openapiHeaderObject{
					Description: "number header description",
					Type:        "number",
					Default:     json.RawMessage("1.2"),
					Pattern:     "^[-+]?[0-9]*\\.?[0-9]+([eE][-+]?[0-9]+)?$",
				},
			},
		}[0], response.Headers, "response.Headers"; !reflect.DeepEqual(is, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, is, want)
		}

	}
	t.Run("verify template options set via proto options", func(t *testing.T) {
		file := newFile()
		proto.SetExtension(proto.Message(file.Services[0].Methods[0].Options), openapi_options.E_Openapiv2Operation, &openapiOperation)
		reg := descriptor.NewRegistry()
		verifyTemplateHeaders(t, reg, file, nil)
	})
}

func TestValidateHeaderType(t *testing.T) {
	type test struct {
		Type          string
		Format        string
		expectedError error
	}
	tests := []test{
		{
			"string",
			"date-time",
			nil,
		},
		{
			"boolean",
			"",
			nil,
		},
		{
			"integer",
			"uint",
			nil,
		},
		{
			"integer",
			"uint8",
			nil,
		},
		{
			"integer",
			"uint16",
			nil,
		},
		{
			"integer",
			"uint32",
			nil,
		},
		{
			"integer",
			"uint64",
			nil,
		},
		{
			"integer",
			"int",
			nil,
		},
		{
			"integer",
			"int8",
			nil,
		},
		{
			"integer",
			"int16",
			nil,
		},
		{
			"integer",
			"int32",
			nil,
		},
		{
			"integer",
			"int64",
			nil,
		},
		{
			"integer",
			"float64",
			errors.New("the provided format \"float64\" is not a valid extension of the type \"integer\""),
		},
		{
			"integer",
			"uuid",
			errors.New("the provided format \"uuid\" is not a valid extension of the type \"integer\""),
		},
		{
			"number",
			"uint",
			nil,
		},
		{
			"number",
			"uint8",
			nil,
		},
		{
			"number",
			"uint16",
			nil,
		},
		{
			"number",
			"uint32",
			nil,
		},
		{
			"number",
			"uint64",
			nil,
		},
		{
			"number",
			"int",
			nil,
		},
		{
			"number",
			"int8",
			nil,
		},
		{
			"number",
			"int16",
			nil,
		},
		{
			"number",
			"int32",
			nil,
		},
		{
			"number",
			"int64",
			nil,
		},
		{
			"number",
			"float",
			nil,
		},
		{
			"number",
			"float32",
			nil,
		},
		{
			"number",
			"float64",
			nil,
		},
		{
			"number",
			"complex64",
			nil,
		},
		{
			"number",
			"complex128",
			nil,
		},
		{
			"number",
			"double",
			nil,
		},
		{
			"number",
			"byte",
			nil,
		},
		{
			"number",
			"rune",
			nil,
		},
		{
			"number",
			"uintptr",
			nil,
		},
		{
			"number",
			"date",
			errors.New("the provided format \"date\" is not a valid extension of the type \"number\""),
		},
		{
			"array",
			"",
			errors.New("the provided header type \"array\" is not supported"),
		},
		{
			"foo",
			"",
			errors.New("the provided header type \"foo\" is not supported"),
		},
	}
	for _, v := range tests {
		err := validateHeaderTypeAndFormat(v.Type, v.Format)

		if v.expectedError == nil {
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
		} else {
			if err == nil {
				t.Fatal("expected header error not returned")
			}
			if err.Error() != v.expectedError.Error() {
				t.Errorf("expected error malformed, expected %q, got %q", v.expectedError.Error(), err.Error())
			}
		}
	}

}

func TestValidateDefaultValueType(t *testing.T) {
	type test struct {
		Type          string
		Value         string
		Format        string
		expectedError error
	}
	tests := []test{
		{
			"string",
			`"string"`,
			"",
			nil,
		},
		{
			"string",
			"\"2012-11-01T22:08:41+00:00\"",
			"date-time",
			nil,
		},
		{
			"string",
			"\"2012-11-01\"",
			"date",
			nil,
		},
		{
			"string",
			"0",
			"",
			errors.New("the provided default value \"0\" does not match provider type \"string\", or is not properly quoted with escaped quotations"),
		},
		{
			"string",
			"false",
			"",
			errors.New("the provided default value \"false\" does not match provider type \"string\", or is not properly quoted with escaped quotations"),
		},
		{
			"boolean",
			"true",
			"",
			nil,
		},
		{
			"boolean",
			"0",
			"",
			errors.New("the provided default value \"0\" does not match provider type \"boolean\""),
		},
		{
			"boolean",
			`"string"`,
			"",
			errors.New("the provided default value \"\\\"string\\\"\" does not match provider type \"boolean\""),
		},
		{
			"number",
			"1.2",
			"",
			nil,
		},
		{
			"number",
			"123",
			"",
			nil,
		},
		{
			"number",
			"nan",
			"",
			errors.New("the provided number \"nan\" is not a valid JSON number"),
		},
		{
			"number",
			"NaN",
			"",
			errors.New("the provided number \"NaN\" is not a valid JSON number"),
		},
		{
			"number",
			"-459.67",
			"",
			nil,
		},
		{
			"number",
			"inf",
			"",
			errors.New("the provided number \"inf\" is not a valid JSON number"),
		},
		{
			"number",
			"infinity",
			"",
			errors.New("the provided number \"infinity\" is not a valid JSON number"),
		},
		{
			"number",
			"Inf",
			"",
			errors.New("the provided number \"Inf\" is not a valid JSON number"),
		},
		{
			"number",
			"Infinity",
			"",
			errors.New("the provided number \"Infinity\" is not a valid JSON number"),
		},
		{
			"number",
			"false",
			"",
			errors.New("the provided default value \"false\" does not match provider type \"number\""),
		},
		{
			"number",
			`"string"`,
			"",
			errors.New("the provided default value \"\\\"string\\\"\" does not match provider type \"number\""),
		},
		{
			"integer",
			"2",
			"",
			nil,
		},
		{
			"integer",
			fmt.Sprint(math.MaxInt32),
			"int32",
			nil,
		},
		{
			"integer",
			fmt.Sprint(math.MaxInt32 + 1),
			"int32",
			errors.New("the provided default value \"2147483648\" does not match provided format \"int32\""),
		},
		{
			"integer",
			fmt.Sprint(math.MaxInt64),
			"int64",
			nil,
		},
		{
			"integer",
			"9223372036854775808",
			"int64",
			errors.New("the provided default value \"9223372036854775808\" does not match provided format \"int64\""),
		},
		{
			"integer",
			"18446744073709551615",
			"uint64",
			nil,
		},
		{
			"integer",
			"false",
			"",
			errors.New("the provided default value \"false\" does not match provided type \"integer\""),
		},
		{
			"integer",
			"1.2",
			"",
			errors.New("the provided default value \"1.2\" does not match provided type \"integer\""),
		},
		{
			"integer",
			`"string"`,
			"",
			errors.New("the provided default value \"\\\"string\\\"\" does not match provided type \"integer\""),
		},
	}
	for _, v := range tests {
		err := validateDefaultValueTypeAndFormat(v.Type, v.Value, v.Format)

		if v.expectedError == nil {
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
		} else {
			if err == nil {
				t.Error("expected update error not returned")
			}
			if err.Error() != v.expectedError.Error() {
				t.Errorf("expected error malformed, expected %q, got %q", v.expectedError.Error(), err.Error())
			}
		}
	}

}

func TestApplyTemplateRequestWithoutClientStreaming(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &descriptorpb.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(false),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}

	meth.ServerStreaming = proto.Bool(false)

	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	nested := &descriptor.Message{
		DescriptorProto: nesteddesc,
	}

	nestedField := &descriptor.Field{
		Message:              msg,
		FieldDescriptorProto: msg.GetField()[0],
	}
	intField := &descriptor.Field{
		Message:              nested,
		FieldDescriptorProto: nested.GetField()[0],
	}
	boolField := &descriptor.Field{
		Message:              nested,
		FieldDescriptorProto: nested.GetField()[1],
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			MessageType:    []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
			Service:        []*descriptorpb.ServiceDescriptorProto{svc},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor.Message{msg, nested},
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
								HTTPMethod: "POST",
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/echo", // TODO(achew): Figure out what this hsould really be
								},
								PathParams: []descriptor.Parameter{
									{
										FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
											{
												Name:   "nested",
												Target: nestedField,
											},
											{
												Name:   "int32",
												Target: intField,
											},
										}),
										Target: intField,
									},
								},
								Body: &descriptor.Body{
									FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
										{
											Name:   "nested",
											Target: nestedField,
										},
										{
											Name:   "bool",
											Target: boolField,
										},
									}),
								},
							},
						},
					},
				},
			},
		},
	}
	reg := descriptor.NewRegistry()
	if err := AddErrorDefs(reg); err != nil {
		t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
		return
	}
	err := reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
	})
	if err != nil {
		t.Fatalf("failed to load code generator request: %v", err)
	}
	result, err := applyTemplate(param{File: crossLinkFixture(&file), reg: reg})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want, got := "2.0", result.Swagger; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).Swagger = %s want to be %s", file, got, want)
	}
	if want, got := "", result.BasePath; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).BasePath = %s want to be %s", file, got, want)
	}
	if want, got := ([]string)(nil), result.Schemes; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).Schemes = %s want to be %s", file, got, want)
	}
	if want, got := []string{"application/json"}, result.Consumes; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).Consumes = %s want to be %s", file, got, want)
	}
	if want, got := []string{"application/json"}, result.Produces; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).Produces = %s want to be %s", file, got, want)
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", fmt.Sprint(result))
	}
}

func TestApplyTemplateRequestWithClientStreaming(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &descriptorpb.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(true),
		ServerStreaming: proto.Bool(true),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}

	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	nested := &descriptor.Message{
		DescriptorProto: nesteddesc,
	}

	nestedField := &descriptor.Field{
		Message:              msg,
		FieldDescriptorProto: msg.GetField()[0],
	}
	intField := &descriptor.Field{
		Message:              nested,
		FieldDescriptorProto: nested.GetField()[0],
	}
	boolField := &descriptor.Field{
		Message:              nested,
		FieldDescriptorProto: nested.GetField()[1],
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			MessageType:    []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
			Service:        []*descriptorpb.ServiceDescriptorProto{svc},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor.Message{msg, nested},
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
								HTTPMethod: "POST",
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/echo", // TODO(achew): Figure out what this hsould really be
								},
								PathParams: []descriptor.Parameter{
									{
										FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
											{
												Name:   "nested",
												Target: nestedField,
											},
											{
												Name:   "int32",
												Target: intField,
											},
										}),
										Target: intField,
									},
								},
								Body: &descriptor.Body{
									FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
										{
											Name:   "nested",
											Target: nestedField,
										},
										{
											Name:   "bool",
											Target: boolField,
										},
									}),
								},
							},
						},
					},
				},
			},
		},
	}
	reg := descriptor.NewRegistry()
	if err := AddErrorDefs(reg); err != nil {
		t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
		return
	}
	err := reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
	})
	if err != nil {
		t.Fatalf("failed to load code generator request: %v", err)
	}
	result, err := applyTemplate(param{File: crossLinkFixture(&file), reg: reg})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}

	// Only ExampleMessage must be present, not NestedMessage
	if want, got, name := 3, len(result.Definitions), "len(Definitions)"; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).%s = %d want to be %d", file, name, got, want)
	}
	if _, ok := result.Paths["/v1/echo"].Post.Responses["200"]; !ok {
		t.Errorf("applyTemplate(%#v).%s = expected 200 response to be defined", file, `result.Paths["/v1/echo"].Post.Responses["200"]`)
	} else {
		if want, got, name := "A successful response.(streaming responses)", result.Paths["/v1/echo"].Post.Responses["200"].Description, `result.Paths["/v1/echo"].Post.Responses["200"].Description`; !reflect.DeepEqual(got, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
		}
		streamExampleExampleMessage := result.Paths["/v1/echo"].Post.Responses["200"].Schema
		if want, got, name := "object", streamExampleExampleMessage.Type, `result.Paths["/v1/echo"].Post.Responses["200"].Schema.Type`; !reflect.DeepEqual(got, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
		}
		if want, got, name := "Stream result of exampleExampleMessage", streamExampleExampleMessage.Title, `result.Paths["/v1/echo"].Post.Responses["200"].Schema.Title`; !reflect.DeepEqual(got, want) {
			t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
		}
		streamExampleExampleMessageProperties := *(streamExampleExampleMessage.Properties)
		if want, got, name := 2, len(streamExampleExampleMessageProperties), `len(StreamDefinitions["exampleExampleMessage"].Properties)`; !reflect.DeepEqual(got, want) {
			t.Errorf("applyTemplate(%#v).%s = %d want to be %d", file, name, got, want)
		} else {
			resultProperty := streamExampleExampleMessageProperties[0]
			if want, got, name := "result", resultProperty.Key, `(*(StreamDefinitions["exampleExampleMessage"].Properties))[0].Key`; !reflect.DeepEqual(got, want) {
				t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
			}
			result := resultProperty.Value.(openapiSchemaObject)
			if want, got, name := "#/definitions/exampleExampleMessage", result.Ref, `((*(StreamDefinitions["exampleExampleMessage"].Properties))[0].Value.(openapiSchemaObject)).Ref`; !reflect.DeepEqual(got, want) {
				t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
			}
			errorProperty := streamExampleExampleMessageProperties[1]
			if want, got, name := "error", errorProperty.Key, `(*(StreamDefinitions["exampleExampleMessage"].Properties))[0].Key`; !reflect.DeepEqual(got, want) {
				t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
			}
			err := errorProperty.Value.(openapiSchemaObject)
			if want, got, name := "#/definitions/rpcStatus", err.Ref, `((*(StreamDefinitions["exampleExampleMessage"].Properties))[0].Value.(openapiSchemaObject)).Ref`; !reflect.DeepEqual(got, want) {
				t.Errorf("applyTemplate(%#v).%s = %s want to be %s", file, name, got, want)
			}
		}
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", fmt.Sprint(result))
	}
}

func TestApplyTemplateRequestWithUnusedReferences(t *testing.T) {
	reqdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("string"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Number: proto.Int32(1),
			},
		},
	}
	respdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("EmptyMessage"),
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Example"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("EmptyMessage"),
		ClientStreaming: proto.Bool(false),
		ServerStreaming: proto.Bool(false),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}

	req := &descriptor.Message{
		DescriptorProto: reqdesc,
	}
	resp := &descriptor.Message{
		DescriptorProto: respdesc,
	}
	stringField := &descriptor.Field{
		Message:              req,
		FieldDescriptorProto: req.GetField()[0],
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			MessageType:    []*descriptorpb.DescriptorProto{reqdesc, respdesc},
			Service:        []*descriptorpb.ServiceDescriptorProto{svc},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor.Message{req, resp},
		Services: []*descriptor.Service{
			{
				ServiceDescriptorProto: svc,
				Methods: []*descriptor.Method{
					{
						MethodDescriptorProto: meth,
						RequestType:           req,
						ResponseType:          resp,
						Bindings: []*descriptor.Binding{
							{
								HTTPMethod: "GET",
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/example",
								},
							},
							{
								HTTPMethod: "POST",
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/example/{string}",
								},
								PathParams: []descriptor.Parameter{
									{
										FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
											{
												Name:   "string",
												Target: stringField,
											},
										}),
										Target: stringField,
									},
								},
								Body: &descriptor.Body{
									FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
										{
											Name:   "string",
											Target: stringField,
										},
									}),
								},
							},
						},
					},
				},
			},
		},
	}

	reg := descriptor.NewRegistry()
	if err := AddErrorDefs(reg); err != nil {
		t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
		return
	}
	err := reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
	})
	if err != nil {
		t.Fatalf("failed to load code generator request: %v", err)
	}
	result, err := applyTemplate(param{File: crossLinkFixture(&file), reg: reg})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}

	// Only EmptyMessage must be present, not ExampleMessage (plus error status)
	if want, got, name := 3, len(result.Definitions), "len(Definitions)"; !reflect.DeepEqual(got, want) {
		t.Errorf("applyTemplate(%#v).%s = %d want to be %d", file, name, got, want)
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", fmt.Sprint(result))
	}
}

func TestApplyTemplateRequestWithBodyQueryParameters(t *testing.T) {
	bookDesc := &descriptorpb.DescriptorProto{
		Name: proto.String("Book"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("name"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_REQUIRED.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("id"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_REQUIRED.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	createDesc := &descriptorpb.DescriptorProto{
		Name: proto.String("CreateBookRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("parent"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_REQUIRED.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("book"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_REQUIRED.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Number: proto.Int32(2),
			},
			{
				Name:   proto.String("book_id"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_REQUIRED.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Number: proto.Int32(3),
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("CreateBook"),
		InputType:  proto.String("CreateBookRequest"),
		OutputType: proto.String("Book"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("BookService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}

	bookMsg := &descriptor.Message{
		DescriptorProto: bookDesc,
	}
	createMsg := &descriptor.Message{
		DescriptorProto: createDesc,
	}

	parentField := &descriptor.Field{
		Message:              createMsg,
		FieldDescriptorProto: createMsg.GetField()[0],
	}
	bookField := &descriptor.Field{
		Message:              createMsg,
		FieldMessage:         bookMsg,
		FieldDescriptorProto: createMsg.GetField()[1],
	}
	bookIDField := &descriptor.Field{
		Message:              createMsg,
		FieldDescriptorProto: createMsg.GetField()[2],
	}

	createMsg.Fields = []*descriptor.Field{parentField, bookField, bookIDField}

	newFile := func() descriptor.File {
		return descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
				Name:           proto.String("book.proto"),
				MessageType:    []*descriptorpb.DescriptorProto{bookDesc, createDesc},
				Service:        []*descriptorpb.ServiceDescriptorProto{svc},
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
				},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/book.pb",
				Name: "book_pb",
			},
			Messages: []*descriptor.Message{bookMsg, createMsg},
			Services: []*descriptor.Service{
				{
					ServiceDescriptorProto: svc,
					Methods: []*descriptor.Method{
						{
							MethodDescriptorProto: meth,
							RequestType:           createMsg,
							ResponseType:          bookMsg,
							Bindings: []*descriptor.Binding{
								{
									HTTPMethod: "POST",
									PathTmpl: httprule.Template{
										Version:  1,
										OpCodes:  []int{0, 0},
										Template: "/v1/{parent=publishers/*}/books",
									},
									PathParams: []descriptor.Parameter{
										{
											FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
												{
													Name:   "parent",
													Target: parentField,
												},
											}),
											Target: parentField,
										},
									},
									Body: &descriptor.Body{
										FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{
											{
												Name:   "book",
												Target: bookField,
											},
										}),
									},
								},
							},
						},
					},
				},
			},
		}
	}
	type args struct {
		file descriptor.File
	}
	type paramOut struct {
		Name     string
		In       string
		Required bool
	}
	tests := []struct {
		name string
		args args
		want []paramOut
	}{
		{
			name: "book_in_body",
			args: args{file: newFile()},
			want: []paramOut{
				{"parent", "path", true},
				{"body", "body", true},
				{"book_id", "query", false},
			},
		},
		{
			name: "book_in_query",
			args: args{file: func() descriptor.File {
				f := newFile()
				f.Services[0].Methods[0].Bindings[0].Body = nil
				return f
			}()},
			want: []paramOut{
				{"parent", "path", true},
				{"book", "query", false},
				{"book_id", "query", false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := descriptor.NewRegistry()
			if err := AddErrorDefs(reg); err != nil {
				t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
				return
			}
			err := reg.Load(&pluginpb.CodeGeneratorRequest{ProtoFile: []*descriptorpb.FileDescriptorProto{tt.args.file.FileDescriptorProto}})
			if err != nil {
				t.Errorf("Registry.Load() failed with %v; want success", err)
				return
			}
			result, err := applyTemplate(param{File: crossLinkFixture(&tt.args.file), reg: reg})
			if err != nil {
				t.Errorf("applyTemplate(%#v) failed with %v; want success", tt.args.file, err)
				return
			}

			if _, ok := result.Paths["/v1/{parent=publishers/*}/books"].Post.Responses["200"]; !ok {
				t.Errorf("applyTemplate(%#v).%s = expected 200 response to be defined", tt.args.file, `result.Paths["/v1/{parent=publishers/*}/books"].Post.Responses["200"]`)
			} else {

				if want, got, name := 3, len(result.Paths["/v1/{parent=publishers/*}/books"].Post.Parameters), `len(result.Paths["/v1/{parent=publishers/*}/books"].Post.Parameters)`; !reflect.DeepEqual(got, want) {
					t.Errorf("applyTemplate(%#v).%s = %d want to be %d", tt.args.file, name, got, want)
				}

				for i, want := range tt.want {
					p := result.Paths["/v1/{parent=publishers/*}/books"].Post.Parameters[i]
					if got, name := (paramOut{p.Name, p.In, p.Required}), `result.Paths["/v1/{parent=publishers/*}/books"].Post.Parameters[0]`; !reflect.DeepEqual(got, want) {
						t.Errorf("applyTemplate(%#v).%s = %v want to be %v", tt.args.file, name, got, want)
					}
				}

			}

			// If there was a failure, print out the input and the json result for debugging.
			if t.Failed() {
				t.Errorf("had: %s", tt.args.file)
				t.Errorf("got: %s", fmt.Sprint(result))
			}
		})
	}

}

// TestApplyTemplateProtobufAny tests that the protobufAny definition is correctly rendered with the @type field and
// allowing additional properties.
func TestApplyTemplateProtobufAny(t *testing.T) {
	// checkProtobufAnyFormat verifies the only property should be @type and additional properties are allowed
	checkProtobufAnyFormat := func(t *testing.T, protobufAny openapiSchemaObject) {
		anyPropsJSON, err := protobufAny.Properties.MarshalJSON()
		if err != nil {
			t.Errorf("protobufAny.Properties.MarshalJSON(), got error = %v", err)
		}
		var anyPropsMap map[string]interface{}
		if err := json.Unmarshal(anyPropsJSON, &anyPropsMap); err != nil {
			t.Errorf("json.Unmarshal(), got error = %v", err)
		}

		// @type should exist
		if _, ok := anyPropsMap["@type"]; !ok {
			t.Errorf("protobufAny.Properties missing key, \"@type\". got = %#v", anyPropsMap)
		}

		// and @type should be the only property
		if len(anyPropsMap) > 1 {
			t.Errorf("len(protobufAny.Properties) = %v, want = %v", len(anyPropsMap), 1)
		}

		// protobufAny should have additionalProperties allowed
		if protobufAny.AdditionalProperties == nil {
			t.Errorf("protobufAny.AdditionalProperties = nil, want not-nil")
		}
	}

	type args struct {
		regConfig      func(registry *descriptor.Registry)
		msgContainsAny bool
	}
	tests := []struct {
		name               string
		args               args
		wantNumDefinitions int
	}{
		{
			// our proto schema doesn't directly use protobufAny, but it is implicitly used by rpcStatus being
			// automatically rendered
			name: "default_protobufAny_from_rpcStatus",
			args: args{
				msgContainsAny: false,
			},
			wantNumDefinitions: 4,
		},
		{
			// we have a protobufAny in a message, it should contain a ref inside the custom message
			name: "protobufAny_referenced_in_message",
			args: args{
				msgContainsAny: true,
			},
			wantNumDefinitions: 4,
		},
		{
			// we have a protobufAny in a message but with automatic rendering of rpcStatus disabled
			name: "protobufAny_referenced_in_message_with_default_errors_disabled",
			args: args{
				msgContainsAny: true,
				regConfig: func(reg *descriptor.Registry) {
					reg.SetDisableDefaultErrors(true)
				},
			},
			wantNumDefinitions: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqdesc := &descriptorpb.DescriptorProto{
				Name: proto.String("ExampleMessage"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("name"),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						Number: proto.Int32(1),
					},
				},
			}
			respdesc := &descriptorpb.DescriptorProto{
				Name: proto.String("EmptyMessage"),
			}
			meth := &descriptorpb.MethodDescriptorProto{
				Name:            proto.String("Example"),
				InputType:       proto.String("ExampleMessage"),
				OutputType:      proto.String("EmptyMessage"),
				ClientStreaming: proto.Bool(false),
				ServerStreaming: proto.Bool(false),
			}
			svc := &descriptorpb.ServiceDescriptorProto{
				Name:   proto.String("ExampleService"),
				Method: []*descriptorpb.MethodDescriptorProto{meth},
			}

			req := &descriptor.Message{
				DescriptorProto: reqdesc,
			}
			resp := &descriptor.Message{
				DescriptorProto: respdesc,
			}
			file := descriptor.File{
				FileDescriptorProto: &descriptorpb.FileDescriptorProto{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
					Name:           proto.String("example.proto"),
					Package:        proto.String("example"),
					MessageType:    []*descriptorpb.DescriptorProto{reqdesc, respdesc},
					Service:        []*descriptorpb.ServiceDescriptorProto{svc},
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
					},
				},
				GoPkg: descriptor.GoPackage{
					Path: "example.com/path/to/example/example.pb",
					Name: "example_pb",
				},
				Messages: []*descriptor.Message{req, resp},
				Services: []*descriptor.Service{
					{
						ServiceDescriptorProto: svc,
						Methods: []*descriptor.Method{
							{
								MethodDescriptorProto: meth,
								RequestType:           req,
								ResponseType:          resp,
							},
						},
					},
				},
			}

			reg := descriptor.NewRegistry()
			reg.SetGenerateUnboundMethods(true)

			if tt.args.regConfig != nil {
				tt.args.regConfig(reg)
			}

			if err := AddErrorDefs(reg); err != nil {
				t.Errorf("AddErrorDefs(%#v) failed with %v; want success", reg, err)
				return
			}

			protoFiles := []*descriptorpb.FileDescriptorProto{
				file.FileDescriptorProto,
			}

			if tt.args.msgContainsAny {
				// add an Any field to the request message
				reqdesc.Field = append(reqdesc.Field, &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("any_value"),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					TypeName: proto.String(".google.protobuf.Any"),
					Number:   proto.Int32(2),
				})

				// update the dependencies to import it
				file.Dependency = append(file.Dependency, "google/protobuf/any.proto")

				anyDescriptorProto := protodesc.ToFileDescriptorProto((&anypb.Any{}).ProtoReflect().Descriptor().ParentFile())
				anyDescriptorProto.SourceCodeInfo = &descriptorpb.SourceCodeInfo{}

				// prepend the anyDescriptorProto to the protoFiles slice so that the dependency can be resolved
				protoFiles = append(append(make([]*descriptorpb.FileDescriptorProto, 0, len(protoFiles)+1), anyDescriptorProto), protoFiles[0:]...)
			}

			err := reg.Load(&pluginpb.CodeGeneratorRequest{
				ProtoFile:      protoFiles,
				FileToGenerate: []string{file.GetName()},
			})
			if err != nil {
				t.Fatalf("failed to load code generator request: %v", err)
			}

			target, err := reg.LookupFile(file.GetName())
			if err != nil {
				t.Fatalf("failed to lookup file from reg: %v", err)
			}
			result, err := applyTemplate(param{File: crossLinkFixture(target), reg: reg})
			if err != nil {
				t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
				return
			}

			if want, got, name := tt.wantNumDefinitions, len(result.Definitions), "len(Definitions)"; !reflect.DeepEqual(got, want) {
				t.Errorf("applyTemplate(%#v).%s = %d want to be %d", file, name, got, want)
			}

			protobufAny, ok := result.Definitions["protobufAny"]
			if !ok {
				t.Error("expecting Definitions to contain protobufAny")
			}

			checkProtobufAnyFormat(t, protobufAny)

			// If there was a failure, print out the input and the json result for debugging.
			if t.Failed() {
				t.Errorf("had: %s", file)
				resultJSON, _ := json.Marshal(result)
				t.Errorf("got: %s", resultJSON)
			}
		})
	}
}

func generateFieldsForJSONReservedName() []*descriptor.Field {
	fields := make([]*descriptor.Field, 0)
	fieldName := string("json_name")
	fieldJSONName := string("jsonNAME")
	fieldDescriptor := descriptorpb.FieldDescriptorProto{Name: &fieldName, JsonName: &fieldJSONName}
	field := &descriptor.Field{FieldDescriptorProto: &fieldDescriptor}
	return append(fields, field)
}

func generateMsgsForJSONReservedName() []*descriptor.Message {
	result := make([]*descriptor.Message, 0)
	// The first message, its field is field_abc and its type is NewType
	// NewType field_abc
	fieldName := "field_abc"
	fieldJSONName := "fieldAbc"
	messageName1 := "message1"
	messageType := "pkg.a.NewType"
	pfd := descriptorpb.FieldDescriptorProto{Name: &fieldName, JsonName: &fieldJSONName, TypeName: &messageType}
	result = append(result,
		&descriptor.Message{
			DescriptorProto: &descriptorpb.DescriptorProto{
				Name: &messageName1, Field: []*descriptorpb.FieldDescriptorProto{&pfd},
			},
		})
	// The second message, its name is NewName, its type is string
	// message NewType {
	//    string field_newName [json_name = RESERVEDJSONNAME]
	// }
	messageName := "NewType"
	field := "field_newName"
	fieldJSONName2 := "RESERVEDJSONNAME"
	pfd2 := descriptorpb.FieldDescriptorProto{Name: &field, JsonName: &fieldJSONName2}
	result = append(result, &descriptor.Message{
		DescriptorProto: &descriptorpb.DescriptorProto{
			Name: &messageName, Field: []*descriptorpb.FieldDescriptorProto{&pfd2},
		},
	})
	return result
}

func TestTemplateWithJsonCamelCase(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/test/{test_id}", "/test/{testId}"},
		{"/test1/{test1_id}/test2/{test2_id}", "/test1/{test1Id}/test2/{test2Id}"},
		{"/test1/{test1_id}/{test2_id}", "/test1/{test1Id}/{test2Id}"},
		{"/test1/test2/{test1_id}/{test2_id}", "/test1/test2/{test1Id}/{test2Id}"},
		{"/test1/{test1_id1_id2}", "/test1/{test1Id1Id2}"},
		{"/test1/{test1_id1_id2}/test2/{test2_id3_id4}", "/test1/{test1Id1Id2}/test2/{test2Id3Id4}"},
		{"/test1/test2/{test1_id1_id2}/{test2_id3_id4}", "/test1/test2/{test1Id1Id2}/{test2Id3Id4}"},
		{"test/{a}", "test/{a}"},
		{"test/{ab}", "test/{ab}"},
		{"test/{a_a}", "test/{aA}"},
		{"test/{ab_c}", "test/{abC}"},
		{"test/{json_name}", "test/{jsonNAME}"},
		{"test/{field_abc.field_newName}", "test/{fieldAbc.RESERVEDJSONNAME}"},
	}
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(true)
	for _, data := range tests {
		actual := templateToOpenAPIPath(data.input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		if data.expected != actual {
			t.Errorf("Expected templateToOpenAPIPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
}

func TestTemplateWithoutJsonCamelCase(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/test/{test_id}", "/test/{test_id}"},
		{"/test1/{test1_id}/test2/{test2_id}", "/test1/{test1_id}/test2/{test2_id}"},
		{"/test1/{test1_id}/{test2_id}", "/test1/{test1_id}/{test2_id}"},
		{"/test1/test2/{test1_id}/{test2_id}", "/test1/test2/{test1_id}/{test2_id}"},
		{"/test1/{test1_id1_id2}", "/test1/{test1_id1_id2}"},
		{"/test1/{test1_id1_id2}/test2/{test2_id3_id4}", "/test1/{test1_id1_id2}/test2/{test2_id3_id4}"},
		{"/test1/test2/{test1_id1_id2}/{test2_id3_id4}", "/test1/test2/{test1_id1_id2}/{test2_id3_id4}"},
		{"test/{a}", "test/{a}"},
		{"test/{ab}", "test/{ab}"},
		{"test/{a_a}", "test/{a_a}"},
		{"test/{json_name}", "test/{json_name}"},
		{"test/{field_abc.field_newName}", "test/{field_abc.field_newName}"},
	}
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(false)
	for _, data := range tests {
		actual := templateToOpenAPIPath(data.input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		if data.expected != actual {
			t.Errorf("Expected templateToOpenAPIPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
}

func TestTemplateToOpenAPIPath(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/test", "/test"},
		{"/{test}", "/{test}"},
		{"/{test=prefix/*}", "/{test}"},
		{"/{test=prefix/that/has/multiple/parts/to/it/*}", "/{test}"},
		{"/{test1}/{test2}", "/{test1}/{test2}"},
		{"/{test1}/{test2}/", "/{test1}/{test2}/"},
		{"/{name=prefix/*}", "/{name=prefix/*}"},
		{"/{name=prefix1/*/prefix2/*}", "/{name=prefix1/*/prefix2/*}"},
		{"/{user.name=prefix/*}", "/{user.name=prefix/*}"},
		{"/{user.name=prefix1/*/prefix2/*}", "/{user.name=prefix1/*/prefix2/*}"},
		{"/{parent=prefix/*}/children", "/{parent=prefix/*}/children"},
		{"/{name=prefix/*}:customMethod", "/{name=prefix/*}:customMethod"},
		{"/{name=prefix1/*/prefix2/*}:customMethod", "/{name=prefix1/*/prefix2/*}:customMethod"},
		{"/{user.name=prefix/*}:customMethod", "/{user.name=prefix/*}:customMethod"},
		{"/{user.name=prefix1/*/prefix2/*}:customMethod", "/{user.name=prefix1/*/prefix2/*}:customMethod"},
		{"/{parent=prefix/*}/children:customMethod", "/{parent=prefix/*}/children:customMethod"},
	}
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(false)
	for _, data := range tests {
		actual := templateToOpenAPIPath(data.input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		if data.expected != actual {
			t.Errorf("Expected templateToOpenAPIPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
	reg.SetUseJSONNamesForFields(true)
	for _, data := range tests {
		actual := templateToOpenAPIPath(data.input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		if data.expected != actual {
			t.Errorf("Expected templateToOpenAPIPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
}

func BenchmarkTemplateToOpenAPIPath(b *testing.B) {
	const input = "/{user.name=prefix1/*/prefix2/*}:customMethod"

	b.Run("with JSON names", func(b *testing.B) {
		reg := descriptor.NewRegistry()
		reg.SetUseJSONNamesForFields(false)

		for i := 0; i < b.N; i++ {
			_ = templateToOpenAPIPath(input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		}
	})

	b.Run("without JSON names", func(b *testing.B) {
		reg := descriptor.NewRegistry()
		reg.SetUseJSONNamesForFields(true)

		for i := 0; i < b.N; i++ {
			_ = templateToOpenAPIPath(input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		}
	})
}

func TestResolveFullyQualifiedNameToOpenAPIName(t *testing.T) {
	var tests = []struct {
		input          string
		output         string
		listOfFQMNs    []string
		namingStrategy string
	}{
		{
			".a.b.C",
			"C",
			[]string{
				".a.b.C",
			},
			"legacy",
		},
		{
			".a.b.C",
			"C",
			[]string{
				".a.b.C",
			},
			"simple",
		},
		{
			".a.b.C",
			"abC",
			[]string{
				".a.C",
				".a.b.C",
			},
			"legacy",
		},
		{
			".a.b.C",
			"b.C",
			[]string{
				".a.C",
				".a.b.C",
			},
			"simple",
		},
		{
			".a.b.C",
			"abC",
			[]string{
				".C",
				".a.C",
				".a.b.C",
			},
			"legacy",
		},
		{
			".a.b.C",
			"b.C",
			[]string{
				".C",
				".a.C",
				".a.b.C",
			},
			"simple",
		},
		{
			".a.b.C",
			"a.b.C",
			[]string{
				".C",
				".a.C",
				".a.b.C",
			},
			"fqn",
		},
	}

	for _, data := range tests {
		names := resolveFullyQualifiedNameToOpenAPINames(data.listOfFQMNs, data.namingStrategy)
		output := names[data.input]
		if output != data.output {
			t.Errorf("Expected fullyQualifiedNameToOpenAPIName(%v, %s) to be %s but got %s",
				data.input, data.namingStrategy, data.output, output)
		}
	}
}

func TestFQMNtoOpenAPIName(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/test", "/test"},
		{"/{test}", "/{test}"},
		{"/{test=prefix/*}", "/{test}"},
		{"/{test=prefix/that/has/multiple/parts/to/it/*}", "/{test}"},
		{"/{test1}/{test2}", "/{test1}/{test2}"},
		{"/{test1}/{test2}/", "/{test1}/{test2}/"},
	}
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(false)
	for _, data := range tests {
		actual := templateToOpenAPIPath(data.input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		if data.expected != actual {
			t.Errorf("Expected templateToOpenAPIPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
	reg.SetUseJSONNamesForFields(true)
	for _, data := range tests {
		actual := templateToOpenAPIPath(data.input, reg, generateFieldsForJSONReservedName(), generateMsgsForJSONReservedName())
		if data.expected != actual {
			t.Errorf("Expected templateToOpenAPIPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
}

func TestSchemaOfField(t *testing.T) {
	type test struct {
		field          *descriptor.Field
		refs           refMap
		expected       openapiSchemaObject
		openAPIOptions *openapiconfig.OpenAPIOptions
	}

	jsonSchema := &openapi_options.JSONSchema{
		Title:       "field title",
		Description: "field description",
	}

	var fieldOptions = new(descriptorpb.FieldOptions)
	proto.SetExtension(fieldOptions, openapi_options.E_Openapiv2Field, jsonSchema)

	var requiredField = []annotations.FieldBehavior{annotations.FieldBehavior_REQUIRED}
	var requiredFieldOptions = new(descriptorpb.FieldOptions)
	proto.SetExtension(requiredFieldOptions, annotations.E_FieldBehavior, requiredField)

	var outputOnlyField = []annotations.FieldBehavior{annotations.FieldBehavior_OUTPUT_ONLY}
	var outputOnlyOptions = new(descriptorpb.FieldOptions)
	proto.SetExtension(outputOnlyOptions, annotations.E_FieldBehavior, outputOnlyField)

	tests := []test{
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name: proto.String("primitive_field"),
					Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:  proto.String("repeated_primitive_field"),
					Type:  descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					Label: descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "array",
					Items: &openapiItemsObject{
						Type: "string",
					},
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.FieldMask"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.Timestamp"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "string",
					Format: "date-time",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.Duration"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.StringValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("repeated_wrapped_field"),
					TypeName: proto.String(".google.protobuf.StringValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "array",
					Items: &openapiItemsObject{
						Type: "string",
					},
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.BytesValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "string",
					Format: "byte",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.Int32Value"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "integer",
					Format: "int32",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.UInt32Value"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "integer",
					Format: "int64",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.Int64Value"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "string",
					Format: "int64",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.UInt64Value"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "string",
					Format: "uint64",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.FloatValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "number",
					Format: "float",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.DoubleValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "number",
					Format: "double",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.BoolValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "boolean",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.Struct"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "object",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.Value"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "object",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.ListValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "array",
					Items: (*openapiItemsObject)(&schemaCore{
						Type: "object",
					}),
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("wrapped_field"),
					TypeName: proto.String(".google.protobuf.NullValue"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("message_field"),
					TypeName: proto.String(".example.Message"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				},
			},
			refs: refMap{".example.Message": struct{}{}},
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Ref: "#/definitions/exampleMessage",
				},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("map_field"),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					TypeName: proto.String(".example.Message.MapFieldEntry"),
					Options:  fieldOptions,
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "object",
				},
				AdditionalProperties: &openapiSchemaObject{
					schemaCore: schemaCore{Type: "string"},
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:    proto.String("array_field"),
					Label:   descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
					Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					Options: fieldOptions,
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:  "array",
					Items: (*openapiItemsObject)(&schemaCore{Type: "string"}),
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:    proto.String("primitive_field"),
					Label:   descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:    descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
					Options: fieldOptions,
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "integer",
					Format: "int32",
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("message_field"),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					TypeName: proto.String(".example.Empty"),
					Options:  fieldOptions,
				},
			},
			refs: refMap{".example.Empty": struct{}{}},
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Ref: "#/definitions/exampleEmpty",
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("map_field"), // should be called map_field_option but it's not valid map field name
					Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					TypeName: proto.String(".example.Message.MapFieldEntry"),
				},
			},
			openAPIOptions: &openapiconfig.OpenAPIOptions{
				Field: []*openapiconfig.OpenAPIFieldOption{
					{
						Field:  "example.Message.map_field",
						Option: jsonSchema,
					},
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "object",
				},
				AdditionalProperties: &openapiSchemaObject{
					schemaCore: schemaCore{Type: "string"},
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:  proto.String("array_field_option"),
					Label: descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
					Type:  descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				},
			},
			openAPIOptions: &openapiconfig.OpenAPIOptions{
				Field: []*openapiconfig.OpenAPIFieldOption{
					{
						Field:  "example.Message.array_field_option",
						Option: jsonSchema,
					},
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:  "array",
					Items: (*openapiItemsObject)(&schemaCore{Type: "string"}),
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:  proto.String("primitive_field_option"),
					Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:  descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				},
			},
			openAPIOptions: &openapiconfig.OpenAPIOptions{
				Field: []*openapiconfig.OpenAPIFieldOption{
					{
						Field:  "example.Message.primitive_field_option",
						Option: jsonSchema,
					},
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type:   "integer",
					Format: "int32",
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("message_field_option"),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					TypeName: proto.String(".example.Empty"),
				},
			},
			openAPIOptions: &openapiconfig.OpenAPIOptions{
				Field: []*openapiconfig.OpenAPIFieldOption{
					{
						Field:  "example.Message.message_field_option",
						Option: jsonSchema,
					},
				},
			},
			refs: refMap{".example.Empty": struct{}{}},
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Ref: "#/definitions/exampleEmpty",
				},
				Title:       "field title",
				Description: "field description",
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:    proto.String("required_via_field_behavior_field"),
					Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					Options: requiredFieldOptions,
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
				Required: []string{"required_via_field_behavior_field"},
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:    proto.String("readonly_via_field_behavior_field"),
					Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					Options: outputOnlyOptions,
				},
			},
			refs: make(refMap),
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Type: "string",
				},
				ReadOnly: true,
			},
		},
		{
			field: &descriptor.Field{
				FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("message_field"),
					TypeName: proto.String(".example.Message"),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					Options:  requiredFieldOptions,
				},
			},
			refs: refMap{".example.Message": struct{}{}},
			expected: openapiSchemaObject{
				schemaCore: schemaCore{
					Ref: "#/definitions/exampleMessage",
				},
			},
		},
	}
	for _, test := range tests {
		reg := descriptor.NewRegistry()
		req := &pluginpb.CodeGeneratorRequest{
			ProtoFile: []*descriptorpb.FileDescriptorProto{
				{
					Name:    proto.String("third_party/google.proto"),
					Package: proto.String("google.protobuf"),
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("third_party/google"),
					},
					MessageType: []*descriptorpb.DescriptorProto{
						protodesc.ToDescriptorProto((&structpb.Struct{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&structpb.Value{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&structpb.ListValue{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&field_mask.FieldMask{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&timestamppb.Timestamp{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&durationpb.Duration{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.StringValue{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.BytesValue{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.Int32Value{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.UInt32Value{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.Int64Value{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.UInt64Value{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.FloatValue{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.DoubleValue{}).ProtoReflect().Descriptor()),
						protodesc.ToDescriptorProto((&wrapperspb.BoolValue{}).ProtoReflect().Descriptor()),
					},
					EnumType: []*descriptorpb.EnumDescriptorProto{
						protodesc.ToEnumDescriptorProto(structpb.NullValue(0).Descriptor()),
					},
				},
				{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
					Name:           proto.String("example.proto"),
					Package:        proto.String("example"),
					Dependency:     []string{"third_party/google.proto"},
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
					},
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: proto.String("Message"),
							Field: []*descriptorpb.FieldDescriptorProto{
								{
									Name:   proto.String("value"),
									Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
									Number: proto.Int32(1),
								},
								func() *descriptorpb.FieldDescriptorProto {
									fd := test.field.FieldDescriptorProto
									fd.Number = proto.Int32(2)
									return fd
								}(),
							},
							NestedType: []*descriptorpb.DescriptorProto{
								{
									Name:    proto.String("MapFieldEntry"),
									Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(true)},
									Field: []*descriptorpb.FieldDescriptorProto{
										{
											Name:   proto.String("key"),
											Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
											Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
											Number: proto.Int32(1),
										},
										{
											Name:   proto.String("value"),
											Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
											Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
											Number: proto.Int32(2),
										},
									},
								},
							},
						},
						{
							Name: proto.String("Empty"),
						},
					},
					EnumType: []*descriptorpb.EnumDescriptorProto{
						{
							Name: proto.String("MessageType"),
							Value: []*descriptorpb.EnumValueDescriptorProto{
								{
									Name:   proto.String("MESSAGE_TYPE_1"),
									Number: proto.Int32(0),
								},
							},
						},
					},
					Service: []*descriptorpb.ServiceDescriptorProto{},
				},
			},
		}
		err := reg.Load(req)
		if err != nil {
			t.Errorf("failed to reg.Load(req): %v", err)
		}

		// set field's parent message pointer to message so field can resolve its FQFN
		test.field.Message = &descriptor.Message{
			DescriptorProto: req.ProtoFile[1].MessageType[0],
			File: &descriptor.File{
				FileDescriptorProto: req.ProtoFile[1],
			},
		}

		if test.openAPIOptions != nil {
			if err := reg.RegisterOpenAPIOptions(test.openAPIOptions); err != nil {
				t.Fatalf("failed to register OpenAPI options: %s", err)
			}
		}

		refs := make(refMap)
		actual := schemaOfField(test.field, reg, refs)
		expectedSchemaObject := test.expected
		if e, a := expectedSchemaObject, actual; !reflect.DeepEqual(a, e) {
			t.Errorf("Expected schemaOfField(%v) = \n%#+v, actual: \n%#+v", test.field, e, a)
		}
		if !reflect.DeepEqual(refs, test.refs) {
			t.Errorf("Expected schemaOfField(%v) to add refs %v, not %v", test.field, test.refs, refs)
		}
	}
}

func TestRenderMessagesAsDefinition(t *testing.T) {
	jsonSchema := &openapi_options.JSONSchema{
		Title:       "field title",
		Description: "field description",
		Required:    []string{"aRequiredField"},
	}

	var requiredField = new(descriptorpb.FieldOptions)
	proto.SetExtension(requiredField, openapi_options.E_Openapiv2Field, jsonSchema)

	var fieldBehaviorRequired = []annotations.FieldBehavior{annotations.FieldBehavior_REQUIRED}
	var requiredFieldOptions = new(descriptorpb.FieldOptions)
	proto.SetExtension(requiredFieldOptions, annotations.E_FieldBehavior, fieldBehaviorRequired)

	var fieldBehaviorOutputOnlyField = []annotations.FieldBehavior{annotations.FieldBehavior_OUTPUT_ONLY}
	var fieldBehaviorOutputOnlyOptions = new(descriptorpb.FieldOptions)
	proto.SetExtension(fieldBehaviorOutputOnlyOptions, annotations.E_FieldBehavior, fieldBehaviorOutputOnlyField)

	tests := []struct {
		descr          string
		msgDescs       []*descriptorpb.DescriptorProto
		schema         map[string]openapi_options.Schema // per-message schema to add
		defs           openapiDefinitionsObject
		openAPIOptions *openapiconfig.OpenAPIOptions
		excludedFields []*descriptor.Field
	}{
		{
			descr: "no OpenAPI options",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{},
			defs: map[string]openapiSchemaObject{
				"Message": {schemaCore: schemaCore{Type: "object"}},
			},
		},
		{
			descr: "example option",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					Example: `{"foo":"bar"}`,
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {schemaCore: schemaCore{
					Type:    "object",
					Example: json.RawMessage(`{"foo":"bar"}`),
				}},
			},
		},
		{
			descr: "example option with something non-json",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					Example: `XXXX anything goes XXXX`,
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {schemaCore: schemaCore{
					Type:    "object",
					Example: json.RawMessage(`XXXX anything goes XXXX`),
				}},
			},
		},
		{
			descr: "external docs option",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					ExternalDocs: &openapi_options.ExternalDocumentation{
						Description: "glorious docs",
						Url:         "https://nada",
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					ExternalDocs: &openapiExternalDocumentationObject{
						Description: "glorious docs",
						URL:         "https://nada",
					},
				},
			},
		},
		{
			descr: "JSONSchema options",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					JsonSchema: &openapi_options.JSONSchema{
						Title:            "title",
						Description:      "desc",
						MultipleOf:       100,
						Maximum:          101,
						ExclusiveMaximum: true,
						Minimum:          1,
						ExclusiveMinimum: true,
						MaxLength:        10,
						MinLength:        3,
						Pattern:          "[a-z]+",
						MaxItems:         20,
						MinItems:         2,
						UniqueItems:      true,
						MaxProperties:    33,
						MinProperties:    22,
						Required:         []string{"req"},
						ReadOnly:         true,
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:            "title",
					Description:      "desc",
					MultipleOf:       100,
					Maximum:          101,
					ExclusiveMaximum: true,
					Minimum:          1,
					ExclusiveMinimum: true,
					MaxLength:        10,
					MinLength:        3,
					Pattern:          "[a-z]+",
					MaxItems:         20,
					MinItems:         2,
					UniqueItems:      true,
					MaxProperties:    33,
					MinProperties:    22,
					Required:         []string{"req"},
					ReadOnly:         true,
				},
			},
		},
		{
			descr: "JSONSchema options from registry",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			openAPIOptions: &openapiconfig.OpenAPIOptions{
				Message: []*openapiconfig.OpenAPIMessageOption{
					{
						Message: "example.Message",
						Option: &openapi_options.Schema{
							JsonSchema: &openapi_options.JSONSchema{
								Title:            "title",
								Description:      "desc",
								MultipleOf:       100,
								Maximum:          101,
								ExclusiveMaximum: true,
								Minimum:          1,
								ExclusiveMinimum: true,
								MaxLength:        10,
								MinLength:        3,
								Pattern:          "[a-z]+",
								MaxItems:         20,
								MinItems:         2,
								UniqueItems:      true,
								MaxProperties:    33,
								MinProperties:    22,
								Required:         []string{"req"},
								ReadOnly:         true,
							},
						},
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:            "title",
					Description:      "desc",
					MultipleOf:       100,
					Maximum:          101,
					ExclusiveMaximum: true,
					Minimum:          1,
					ExclusiveMinimum: true,
					MaxLength:        10,
					MinLength:        3,
					Pattern:          "[a-z]+",
					MaxItems:         20,
					MinItems:         2,
					UniqueItems:      true,
					MaxProperties:    33,
					MinProperties:    22,
					Required:         []string{"req"},
					ReadOnly:         true,
				},
			},
		},
		{
			descr: "JSONSchema with required properties",
			msgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("Message"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:    proto.String("aRequiredField"),
							Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number:  proto.Int32(1),
							Options: requiredField,
						},
					},
				},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					JsonSchema: &openapi_options.JSONSchema{
						Title:       "title",
						Description: "desc",
						Required:    []string{"req"},
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:       "title",
					Description: "desc",
					Required:    []string{"req", "aRequiredField"},
					Properties: &openapiSchemaObjectProperties{
						{
							Key: "aRequiredField",
							Value: openapiSchemaObject{
								schemaCore: schemaCore{
									Type: "string",
								},
								Description: "field description",
								Title:       "field title",
								Required:    []string{"aRequiredField"},
							},
						},
					},
				},
			},
		},
		{
			descr: "JSONSchema with excluded fields",
			msgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("Message"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:    proto.String("aRequiredField"),
							Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number:  proto.Int32(1),
							Options: requiredField,
						},
						{
							Name:   proto.String("anExcludedField"),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number: proto.Int32(2),
						},
					},
				},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					JsonSchema: &openapi_options.JSONSchema{
						Title:       "title",
						Description: "desc",
						Required:    []string{"req"},
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:       "title",
					Description: "desc",
					Required:    []string{"req", "aRequiredField"},
					Properties: &openapiSchemaObjectProperties{
						{
							Key: "aRequiredField",
							Value: openapiSchemaObject{
								schemaCore: schemaCore{
									Type: "string",
								},
								Description: "field description",
								Title:       "field title",
								Required:    []string{"aRequiredField"},
							},
						},
					},
				},
			},
			excludedFields: []*descriptor.Field{
				{
					FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
						Name: strPtr("anExcludedField"),
					},
				},
			},
		},
		{
			descr: "JSONSchema with required properties via field_behavior",
			msgDescs: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("Message"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:    proto.String("aRequiredField"),
							Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number:  proto.Int32(1),
							Options: requiredFieldOptions,
						},
						{
							Name:    proto.String("aOutputOnlyField"),
							Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Number:  proto.Int32(2),
							Options: fieldBehaviorOutputOnlyOptions,
						},
					},
				},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					JsonSchema: &openapi_options.JSONSchema{
						Title:       "title",
						Description: "desc",
						Required:    []string{"req"},
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:       "title",
					Description: "desc",
					Required:    []string{"req", "aRequiredField"},
					Properties: &openapiSchemaObjectProperties{
						{
							Key: "aRequiredField",
							Value: openapiSchemaObject{
								schemaCore: schemaCore{
									Type: "string",
								},
								Required: []string{"aRequiredField"},
							},
						},
						{
							Key: "aOutputOnlyField",
							Value: openapiSchemaObject{
								schemaCore: schemaCore{
									Type: "string",
								},
								ReadOnly: true,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.descr, func(t *testing.T) {

			msgs := []*descriptor.Message{}
			for _, msgdesc := range test.msgDescs {
				msgdesc.Options = &descriptorpb.MessageOptions{}
				msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
			}

			reg := descriptor.NewRegistry()
			file := descriptor.File{
				FileDescriptorProto: &descriptorpb.FileDescriptorProto{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
					Name:           proto.String("example.proto"),
					Package:        proto.String("example"),
					Dependency:     []string{},
					MessageType:    test.msgDescs,
					EnumType:       []*descriptorpb.EnumDescriptorProto{},
					Service:        []*descriptorpb.ServiceDescriptorProto{},
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
					},
				},
				Messages: msgs,
			}
			err := reg.Load(&pluginpb.CodeGeneratorRequest{
				ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
			})
			if err != nil {
				t.Fatalf("failed to load code generator request: %v", err)
			}

			msgMap := map[string]*descriptor.Message{}
			for _, d := range test.msgDescs {
				name := d.GetName()
				msg, err := reg.LookupMsg("example", name)
				if err != nil {
					t.Fatalf("lookup message %v: %v", name, err)
				}
				msgMap[msg.FQMN()] = msg

				if schema, ok := test.schema[name]; ok {
					proto.SetExtension(d.Options, openapi_options.E_Openapiv2Schema, &schema)
				}
			}

			if test.openAPIOptions != nil {
				if err := reg.RegisterOpenAPIOptions(test.openAPIOptions); err != nil {
					t.Fatalf("failed to register OpenAPI options: %s", err)
				}
			}

			refs := make(refMap)
			actual := make(openapiDefinitionsObject)
			renderMessagesAsDefinition(msgMap, actual, reg, refs, test.excludedFields)

			if !reflect.DeepEqual(actual, test.defs) {
				t.Errorf("Expected renderMessagesAsDefinition() to add defs %+v, not %+v", test.defs, actual)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestUpdateOpenAPIDataFromComments(t *testing.T) {

	tests := []struct {
		descr                 string
		openapiSwaggerObject  interface{}
		comments              string
		expectedError         error
		expectedOpenAPIObject interface{}
		useGoTemplate         bool
	}{
		{
			descr:                 "empty comments",
			openapiSwaggerObject:  nil,
			expectedOpenAPIObject: nil,
			comments:              "",
			expectedError:         nil,
		},
		{
			descr:                "set field to read only",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				ReadOnly:    true,
				Description: "... Output only. ...",
			},
			comments:      "... Output only. ...",
			expectedError: nil,
		},
		{
			descr:                "set title",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Title: "Comment with no trailing dot",
			},
			comments:      "Comment with no trailing dot",
			expectedError: nil,
		},
		{
			descr:                "set description",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Description: "Comment with trailing dot.",
			},
			comments:      "Comment with trailing dot.",
			expectedError: nil,
		},
		{
			descr: "use info object",
			openapiSwaggerObject: &openapiSwaggerObject{
				Info: openapiInfoObject{},
			},
			expectedOpenAPIObject: &openapiSwaggerObject{
				Info: openapiInfoObject{
					Description: "Comment with trailing dot.",
				},
			},
			comments:      "Comment with trailing dot.",
			expectedError: nil,
		},
		{
			descr:                "multi line comment with title",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Title:       "First line",
				Description: "Second line",
			},
			comments:      "First line\n\nSecond line",
			expectedError: nil,
		},
		{
			descr:                "multi line comment no title",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Description: "First line.\n\nSecond line",
			},
			comments:      "First line.\n\nSecond line",
			expectedError: nil,
		},
		{
			descr:                "multi line comment with summary with dot",
			openapiSwaggerObject: &openapiOperationObject{},
			expectedOpenAPIObject: &openapiOperationObject{
				Summary:     "First line.",
				Description: "Second line",
			},
			comments:      "First line.\n\nSecond line",
			expectedError: nil,
		},
		{
			descr:                "multi line comment with summary no dot",
			openapiSwaggerObject: &openapiOperationObject{},
			expectedOpenAPIObject: &openapiOperationObject{
				Summary:     "First line",
				Description: "Second line",
			},
			comments:      "First line\n\nSecond line",
			expectedError: nil,
		},
		{
			descr:                 "multi line comment with summary no dot",
			openapiSwaggerObject:  &schemaCore{},
			expectedOpenAPIObject: &schemaCore{},
			comments:              "Any comment",
			expectedError:         errors.New("no description nor summary property"),
		},
		{
			descr:                "without use_go_template",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Title:       "First line",
				Description: "{{import \"documentation.md\"}}",
			},
			comments:      "First line\n\n{{import \"documentation.md\"}}",
			expectedError: nil,
		},
		{
			descr:                "error with use_go_template",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Title:       "First line",
				Description: "open noneexistingfile.txt: no such file or directory",
			},
			comments:      "First line\n\n{{import \"noneexistingfile.txt\"}}",
			expectedError: nil,
			useGoTemplate: true,
		},
		{
			descr:                "template with use_go_template",
			openapiSwaggerObject: &openapiSchemaObject{},
			expectedOpenAPIObject: &openapiSchemaObject{
				Title:       "Template",
				Description: `Description "which means nothing"`,
			},
			comments:      "Template\n\nDescription {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
			expectedError: nil,
			useGoTemplate: true,
		},
	}

	for _, test := range tests {
		t.Run(test.descr, func(t *testing.T) {
			reg := descriptor.NewRegistry()
			if test.useGoTemplate {
				reg.SetUseGoTemplate(true)
			}
			err := updateOpenAPIDataFromComments(reg, test.openapiSwaggerObject, nil, test.comments, false)
			if test.expectedError == nil {
				if err != nil {
					t.Errorf("unexpected error '%v'", err)
				}
				if !reflect.DeepEqual(test.openapiSwaggerObject, test.expectedOpenAPIObject) {
					t.Errorf("openapiSwaggerObject was not updated correctly, expected '%+v', got '%+v'", test.expectedOpenAPIObject, test.openapiSwaggerObject)
				}
			} else {
				if err == nil {
					t.Error("expected update error not returned")
				}
				if !reflect.DeepEqual(test.openapiSwaggerObject, test.expectedOpenAPIObject) {
					t.Errorf("openapiSwaggerObject was not updated correctly, expected '%+v', got '%+v'", test.expectedOpenAPIObject, test.openapiSwaggerObject)
				}
				if err.Error() != test.expectedError.Error() {
					t.Errorf("expected error malformed, expected %q, got %q", test.expectedError.Error(), err.Error())
				}
			}
		})
	}
}

func TestMessageOptionsWithGoTemplate(t *testing.T) {
	tests := []struct {
		descr          string
		msgDescs       []*descriptorpb.DescriptorProto
		schema         map[string]openapi_options.Schema // per-message schema to add
		defs           openapiDefinitionsObject
		openAPIOptions *openapiconfig.OpenAPIOptions
		useGoTemplate  bool
	}{
		{
			descr: "external docs option",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					JsonSchema: &openapi_options.JSONSchema{
						Title:       "{{.Name}}",
						Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
					},
					ExternalDocs: &openapi_options.ExternalDocumentation{
						Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:       "Message",
					Description: `Description "which means nothing"`,
					ExternalDocs: &openapiExternalDocumentationObject{
						Description: `Description "which means nothing"`,
					},
				},
			},
			useGoTemplate: true,
		},
		{
			descr: "external docs option",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			schema: map[string]openapi_options.Schema{
				"Message": {
					JsonSchema: &openapi_options.JSONSchema{
						Title:       "{{.Name}}",
						Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
					},
					ExternalDocs: &openapi_options.ExternalDocumentation{
						Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:       "{{.Name}}",
					Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
					ExternalDocs: &openapiExternalDocumentationObject{
						Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
					},
				},
			},
			useGoTemplate: false,
		},
		{
			descr: "registered OpenAPIOption",
			msgDescs: []*descriptorpb.DescriptorProto{
				{Name: proto.String("Message")},
			},
			openAPIOptions: &openapiconfig.OpenAPIOptions{
				Message: []*openapiconfig.OpenAPIMessageOption{
					{
						Message: "example.Message",
						Option: &openapi_options.Schema{
							JsonSchema: &openapi_options.JSONSchema{
								Title:       "{{.Name}}",
								Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
							},
							ExternalDocs: &openapi_options.ExternalDocumentation{
								Description: "Description {{with \"which means nothing\"}}{{printf \"%q\" .}}{{end}}",
							},
						},
					},
				},
			},
			defs: map[string]openapiSchemaObject{
				"Message": {
					schemaCore: schemaCore{
						Type: "object",
					},
					Title:       "Message",
					Description: `Description "which means nothing"`,
					ExternalDocs: &openapiExternalDocumentationObject{
						Description: `Description "which means nothing"`,
					},
				},
			},
			useGoTemplate: true,
		},
	}

	for _, test := range tests {
		t.Run(test.descr, func(t *testing.T) {

			msgs := []*descriptor.Message{}
			for _, msgdesc := range test.msgDescs {
				msgdesc.Options = &descriptorpb.MessageOptions{}
				msgs = append(msgs, &descriptor.Message{DescriptorProto: msgdesc})
			}

			reg := descriptor.NewRegistry()
			reg.SetUseGoTemplate(test.useGoTemplate)
			file := descriptor.File{
				FileDescriptorProto: &descriptorpb.FileDescriptorProto{
					SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
					Name:           proto.String("example.proto"),
					Package:        proto.String("example"),
					Dependency:     []string{},
					MessageType:    test.msgDescs,
					EnumType:       []*descriptorpb.EnumDescriptorProto{},
					Service:        []*descriptorpb.ServiceDescriptorProto{},
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
					},
				},
				Messages: msgs,
			}
			err := reg.Load(&pluginpb.CodeGeneratorRequest{
				ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto},
			})
			if err != nil {
				t.Fatalf("failed to load code generator request: %v", err)
			}

			msgMap := map[string]*descriptor.Message{}
			for _, d := range test.msgDescs {
				name := d.GetName()
				msg, err := reg.LookupMsg("example", name)
				if err != nil {
					t.Fatalf("lookup message %v: %v", name, err)
				}
				msgMap[msg.FQMN()] = msg

				if schema, ok := test.schema[name]; ok {
					proto.SetExtension(d.Options, openapi_options.E_Openapiv2Schema, &schema)
				}
			}

			if test.openAPIOptions != nil {
				if err := reg.RegisterOpenAPIOptions(test.openAPIOptions); err != nil {
					t.Fatalf("failed to register OpenAPI options: %s", err)
				}
			}

			refs := make(refMap)
			actual := make(openapiDefinitionsObject)
			renderMessagesAsDefinition(msgMap, actual, reg, refs, nil)

			if !reflect.DeepEqual(actual, test.defs) {
				t.Errorf("Expected renderMessagesAsDefinition() to add defs %+v, not %+v", test.defs, actual)
			}
		})
	}
}

func TestTemplateWithoutErrorDefinition(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name:  proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Echo"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}

	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}

	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			MessageType:    []*descriptorpb.DescriptorProto{msgdesc},
			Service:        []*descriptorpb.ServiceDescriptorProto{svc},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
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
								HTTPMethod: "POST",
								PathTmpl: httprule.Template{
									Version:  1,
									OpCodes:  []int{0, 0},
									Template: "/v1/echo",
								},
								Body: &descriptor.Body{
									FieldPath: descriptor.FieldPath([]descriptor.FieldPathComponent{}),
								},
							},
						},
					},
				},
			},
		},
	}
	reg := descriptor.NewRegistry()
	err := reg.Load(&pluginpb.CodeGeneratorRequest{ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto}})
	if err != nil {
		t.Errorf("failed to reg.Load(): %v", err)
		return
	}
	result, err := applyTemplate(param{File: crossLinkFixture(&file), reg: reg})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}

	defRsp, ok := result.Paths["/v1/echo"].Post.Responses["default"]
	if !ok {
		return
	}

	ref := defRsp.Schema.schemaCore.Ref
	refName := strings.TrimPrefix(ref, "#/definitions/")
	if refName == "" {
		t.Fatal("created default Error response with empty reflink")
	}

	if _, ok := result.Definitions[refName]; !ok {
		t.Errorf("default Error response with reflink '%v', but its definition was not found", refName)
	}
}

func Test_getReservedJsonName(t *testing.T) {
	type args struct {
		fieldName                     string
		messageNameToFieldsToJSONName map[string]map[string]string
		fieldNameToType               map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test case 1: single dot use case",
			args{
				fieldName: "abc.a_1",
				messageNameToFieldsToJSONName: map[string]map[string]string{
					"Msg": {
						"a_1": "a1JSONNAME",
						"b_1": "b1JSONNAME",
					},
				},
				fieldNameToType: map[string]string{
					"abc": "pkg1.test.Msg",
					"bcd": "pkg1.test.Msg",
				},
			},
			"a1JSONNAME",
		},
		{
			"test case 2: single dot use case with no existing field",
			args{
				fieldName: "abc.d_1",
				messageNameToFieldsToJSONName: map[string]map[string]string{
					"Msg": {
						"a_1": "a1JSONNAME",
						"b_1": "b1JSONNAME",
					},
				},
				fieldNameToType: map[string]string{
					"abc": "pkg1.test.Msg",
					"bcd": "pkg1.test.Msg",
				},
			},
			"",
		},
		{
			"test case 3: double dot use case",
			args{
				fieldName: "pkg.abc.a_1",
				messageNameToFieldsToJSONName: map[string]map[string]string{
					"Msg": {
						"a_1": "a1JSONNAME",
						"b_1": "b1JSONNAME",
					},
				},
				fieldNameToType: map[string]string{
					"abc": "pkg1.test.Msg",
					"bcd": "pkg1.test.Msg",
				},
			},
			"a1JSONNAME",
		},
		{
			"test case 4: double dot use case with a not existed field",
			args{
				fieldName: "pkg.abc.c_1",
				messageNameToFieldsToJSONName: map[string]map[string]string{
					"Msg": {
						"a_1": "a1JSONNAME",
						"b_1": "b1JSONNAME",
					},
				},
				fieldNameToType: map[string]string{
					"abc": "pkg1.test.Msg",
					"bcd": "pkg1.test.Msg",
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getReservedJSONName(tt.args.fieldName, tt.args.messageNameToFieldsToJSONName, tt.args.fieldNameToType); got != tt.want {
				t.Errorf("getReservedJSONName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIncompleteSecurityRequirement(t *testing.T) {
	swagger := openapi_options.Swagger{
		Security: []*openapi_options.SecurityRequirement{
			{
				SecurityRequirement: map[string]*openapi_options.SecurityRequirement_SecurityRequirementValue{
					"key": nil,
				},
			},
		},
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			SourceCodeInfo: &descriptorpb.SourceCodeInfo{},
			Name:           proto.String("example.proto"),
			Package:        proto.String("example"),
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("github.com/grpc-ecosystem/grpc-gateway/runtime/internal/examplepb;example"),
			},
		},
	}
	proto.SetExtension(proto.Message(file.FileDescriptorProto.Options), openapi_options.E_Openapiv2Swagger, &swagger)
	reg := descriptor.NewRegistry()
	err := reg.Load(&pluginpb.CodeGeneratorRequest{ProtoFile: []*descriptorpb.FileDescriptorProto{file.FileDescriptorProto}})
	if err != nil {
		t.Errorf("failed to reg.Load(): %v", err)
		return
	}
	_, err = applyTemplate(param{File: crossLinkFixture(&file), reg: reg})
	if err == nil {
		t.Errorf("applyTemplate(%#v) did not error as expected", file)
		return
	}
}
