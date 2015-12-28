package genswagger

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/httprule"
	"github.com/golang/protobuf/proto"
	protodescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

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

func TestApplyTemplateSimple(t *testing.T) {
	msgdesc := &protodescriptor.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	meth := &protodescriptor.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	svc := &protodescriptor.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*protodescriptor.MethodDescriptorProto{meth},
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	file := descriptor.File{
		FileDescriptorProto: &protodescriptor.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*protodescriptor.DescriptorProto{msgdesc},
			Service:     []*protodescriptor.ServiceDescriptorProto{svc},
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
									Template: "/v1/echo", // TODO(achew): Figure out what this hsould really be
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := applyTemplate(param{File: crossLinkFixture(&file)})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	got := new(swaggerObject)
	err = json.Unmarshal([]byte(result), got)
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want, is, name := "2.0", got.Swagger, "Swagger"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := "", got.BasePath, "BasePath"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"http", "https"}, got.Schemes, "Schemes"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"application/json"}, got.Consumes, "Consumes"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"application/json"}, got.Produces, "Produces"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", result)
	}
}

func TestApplyTemplateRequestWithoutClientStreaming(t *testing.T) {
	t.Skip()
	msgdesc := &protodescriptor.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*protodescriptor.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    protodescriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     protodescriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &protodescriptor.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*protodescriptor.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  protodescriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   protodescriptor.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  protodescriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   protodescriptor.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &protodescriptor.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(false),
	}
	svc := &protodescriptor.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*protodescriptor.MethodDescriptorProto{meth},
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
		FileDescriptorProto: &protodescriptor.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			MessageType: []*protodescriptor.DescriptorProto{msgdesc, nesteddesc},
			Service:     []*protodescriptor.ServiceDescriptorProto{svc},
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
	result, err := applyTemplate(param{File: crossLinkFixture(&file)})
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	got := new(swaggerObject)
	err = json.Unmarshal([]byte(result), got)
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want, is, name := "2.0", got.Swagger, "Swagger"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := "", got.BasePath, "BasePath"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"http", "https"}, got.Schemes, "Schemes"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"application/json"}, got.Consumes, "Consumes"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := []string{"application/json"}, got.Produces, "Produces"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", file, name, is, want)
	}
	if want, is, name := "Generated for ExampleService.Echo - ", got.Paths["/v1/echo"].Post.Summary, "Paths[/v1/echo].Post.Summary"; !reflect.DeepEqual(is, want) {
		t.Errorf("applyTemplate(file).%s = %s want to be %s", name, is, want)
	}

	// If there was a failure, print out the input and the json result for debugging.
	if t.Failed() {
		t.Errorf("had: %s", file)
		t.Errorf("got: %s", result)
	}
}

func TestApplyTemplateRequestWithClientStreaming(t *testing.T) {
	t.Skip()
	msgdesc := &protodescriptor.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*protodescriptor.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    protodescriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     protodescriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &protodescriptor.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*protodescriptor.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  protodescriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   protodescriptor.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  protodescriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   protodescriptor.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &protodescriptor.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(true),
		ServerStreaming: proto.Bool(true),
	}
	svc := &protodescriptor.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*protodescriptor.MethodDescriptorProto{meth},
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
		FileDescriptorProto: &protodescriptor.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			MessageType: []*protodescriptor.DescriptorProto{msgdesc, nesteddesc},
			Service:     []*protodescriptor.ServiceDescriptorProto{svc},
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
	_, err := applyTemplate(param{File: crossLinkFixture(&file)})
	if err == nil {
		t.Errorf("applyTemplate(%#v) should have failed cause swagger doesn't support streaming", file)
		return
	}
}

func TestTemplateToSwaggerPath(t *testing.T) {
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

	for _, data := range tests {
		actual := templateToSwaggerPath(data.input)
		if data.expected != actual {
			t.Errorf("Expected templateToSwaggerPath(%v) = %v, actual: %v", data.input, data.expected, actual)
		}
	}
}

func TestTemplateToSwaggerPathThatShouldBlowUp(t *testing.T) {
	var tests = []struct {
		input string
	}{
		{"/}}"},
	}

	for _, data := range tests {
		// Wrap in a closure so that the panic catcher is destroyed every time
		func() {
			var err interface{}
			defer func() {
				// recover from panic if one occured. Set err to nil otherwise.
				err = recover()
			}()

			_ = templateToSwaggerPath(data.input)
			if err == nil {
				t.Errorf("Expected templateToSwaggerPath(%v) to blow up", data.input)
			}
		}()
	}
}
