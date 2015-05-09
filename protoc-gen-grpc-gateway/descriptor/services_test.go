package descriptor

import (
	"reflect"
	"testing"

	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/httprule"
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func compilePath(t *testing.T, path string) httprule.Template {
	parsed, err := httprule.Parse(path)
	if err != nil {
		t.Fatalf("httprule.Parse(%q) failed with %v; want success", path, err)
	}
	return parsed.Compile()
}

func testExtractServices(t *testing.T, input []*descriptor.FileDescriptorProto, target string, wantSvcs []*Service) {
	reg := NewRegistry()
	for _, file := range input {
		reg.loadFile(file)
	}
	err := reg.loadServices(target)
	if err != nil {
		t.Errorf("loadServices(%q) failed with %v; want success; files=%v", target, err, input)
	}

	file := reg.files[target]
	svcs := file.Services
	var i int
	for i = 0; i < len(svcs) && i < len(wantSvcs); i++ {
		svc, wantSvc := svcs[i], wantSvcs[i]
		if got, want := svc.ServiceDescriptorProto, wantSvc.ServiceDescriptorProto; !proto.Equal(got, want) {
			t.Errorf("svcs[%d].ServiceDescriptorProto = %v; want %v; input = %v", i, got, want, input)
			continue
		}
		var j int
		for j = 0; j < len(svc.Methods) && j < len(wantSvc.Methods); j++ {
			meth, wantMeth := svc.Methods[j], wantSvc.Methods[j]
			if got, want := meth.MethodDescriptorProto, wantMeth.MethodDescriptorProto; !proto.Equal(got, want) {
				t.Errorf("svcs[%d].Methods[%d].MethodDescriptorProto = %v; want %v; input = %v", i, j, got, want, input)
				continue
			}
			if got, want := meth.PathTmpl, wantMeth.PathTmpl; !reflect.DeepEqual(got, want) {
				t.Errorf("svcs[%d].Methods[%d].PathTmpl = %#v; want %#v; input = %v", i, j, got, want, input)
			}
			if got, want := meth.HTTPMethod, wantMeth.HTTPMethod; got != want {
				t.Errorf("svcs[%d].Methods[%d].HTTPMethod = %q; want %q; input = %v", i, j, got, want, input)
			}
			if got, want := meth.RequestType, wantMeth.RequestType; got.FQMN() != want.FQMN() {
				t.Errorf("svcs[%d].Methods[%d].RequestType = %s; want %s; input = %v", i, j, got.FQMN(), want.FQMN(), input)
			}
			if got, want := meth.ResponseType, wantMeth.ResponseType; got.FQMN() != want.FQMN() {
				t.Errorf("svcs[%d].Methods[%d].ResponseType = %s; want %s; input = %v", i, j, got.FQMN(), want.FQMN(), input)
			}

			var k int
			for k = 0; k < len(meth.PathParams) && k < len(wantMeth.PathParams); k++ {
				param, wantParam := meth.PathParams[k], wantMeth.PathParams[k]
				if got, want := param.FieldPath.String(), wantParam.FieldPath.String(); got != want {
					t.Errorf("svcs[%d].Methods[%d].PathParams[%d].FieldPath.String() = %q; want %q; input = %v", i, j, k, got, want, input)
					continue
				}
				for l := 0; l < len(param.FieldPath) && l < len(wantParam.FieldPath); l++ {
					field, wantField := param.FieldPath[l].Target, wantParam.FieldPath[l].Target
					if got, want := field.FieldDescriptorProto, wantField.FieldDescriptorProto; !proto.Equal(got, want) {
						t.Errorf("svcs[%d].Methods[%d].PathParams[%d].FieldPath[%d].Target.FieldDescriptorProto = %v; want %v; input = %v", i, j, k, l, got, want, input)
					}
				}
			}
			for ; k < len(meth.PathParams); k++ {
				got := meth.PathParams[k].FieldPath.String()
				t.Errorf("svcs[%d].Methods[%d].PathParams[%d] = %q; want it to be missing; input = %v", i, j, k, got, input)
			}
			for ; k < len(wantMeth.PathParams); k++ {
				want := wantMeth.PathParams[k].FieldPath.String()
				t.Errorf("svcs[%d].Methods[%d].PathParams[%d] missing; want %q; input = %v", i, j, k, want, input)
			}

			if got, want := (meth.Body != nil), (wantMeth.Body != nil); got != want {
				if got {
					t.Errorf("svcs[%d].Methods[%d].Body = %q; want it to be missing; input = %v", i, j, meth.Body.FieldPath.String(), input)
				} else {
					t.Errorf("svcs[%d].Methods[%d].Body missing; want %q; input = %v", i, j, wantMeth.Body.FieldPath.String(), input)
				}
			}
		}
		for ; j < len(svc.Methods); j++ {
			got := svc.Methods[j].MethodDescriptorProto
			t.Errorf("svcs[%d].Methods[%d] = %v; want it to be missing; input = %v", i, j, got, input)
		}
		for ; j < len(wantSvc.Methods); j++ {
			want := wantSvc.Methods[j].MethodDescriptorProto
			t.Errorf("svcs[%d].Methods[%d] missing; want %v; input = %v", i, j, want, input)
		}
	}
	for ; i < len(svcs); i++ {
		got := svcs[i].ServiceDescriptorProto
		t.Errorf("svcs[%d] = %v; want it to be missing; input = %v", i, got, input)
	}
	for ; i < len(wantSvcs); i++ {
		want := wantSvcs[i].ServiceDescriptorProto
		t.Errorf("svcs[%d] missing; want %v; input = %v", i, want, input)
	}
}

func crossLinkFixture(f *File) *File {
	for _, m := range f.Messages {
		m.File = f
		for _, f := range m.Fields {
			f.Message = m
		}
	}
	for _, svc := range f.Services {
		svc.File = f
		for _, m := range svc.Methods {
			m.Service = svc
			for _, param := range m.PathParams {
				param.Method = m
			}
			for _, param := range m.QueryParams {
				param.Method = m
			}
		}
	}
	return f
}

func TestExtractServicesSimple(t *testing.T) {
	src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "StringMessage"
			field <
				name: "string"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
			>
		>
		service <
			name: "ExampleService"
			method <
				name: "Echo"
				input_type: "StringMessage"
				output_type: "StringMessage"
				options <
					[google.api.http] <
						post: "/v1/example/echo"
						body: "*"
					>
				>
			>
		>
	`
	var fd descriptor.FileDescriptorProto
	if err := proto.UnmarshalText(src, &fd); err != nil {
		t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
	}
	msg := &Message{
		DescriptorProto: fd.MessageType[0],
		Fields: []*Field{
			{
				FieldDescriptorProto: fd.MessageType[0].Field[0],
			},
		},
	}
	file := &File{
		FileDescriptorProto: &fd,
		GoPkg: GoPackage{
			Path: "path/to/example.pb",
			Name: "example_pb",
		},
		Messages: []*Message{msg},
		Services: []*Service{
			{
				ServiceDescriptorProto: fd.Service[0],
				Methods: []*Method{
					{
						MethodDescriptorProto: fd.Service[0].Method[0],
						PathTmpl:              compilePath(t, "/v1/example/echo"),
						HTTPMethod:            "POST",
						RequestType:           msg,
						ResponseType:          msg,
						Body:                  &Body{FieldPath: nil},
					},
				},
			},
		},
	}

	crossLinkFixture(file)
	testExtractServices(t, []*descriptor.FileDescriptorProto{&fd}, "path/to/example.proto", file.Services)
}

func TestExtractServicesCrossPackage(t *testing.T) {
	srcs := []string{
		`
			name: "path/to/example.proto",
			package: "example"
			message_type <
				name: "StringMessage"
				field <
					name: "string"
					number: 1
					label: LABEL_OPTIONAL
					type: TYPE_STRING
				>
			>
			service <
				name: "ExampleService"
				method <
					name: "ToString"
					input_type: ".another.example.BoolMessage"
					output_type: "StringMessage"
					options <
						[google.api.http] <
							post: "/v1/example/to_s"
							body: "*"
						>
					>
				>
			>
		`, `
			name: "path/to/another/example.proto",
			package: "another.example"
			message_type <
				name: "BoolMessage"
				field <
					name: "bool"
					number: 1
					label: LABEL_OPTIONAL
					type: TYPE_BOOL
				>
			>
		`,
	}
	var fds []*descriptor.FileDescriptorProto
	for _, src := range srcs {
		var fd descriptor.FileDescriptorProto
		if err := proto.UnmarshalText(src, &fd); err != nil {
			t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
		}
		fds = append(fds, &fd)
	}
	stringMsg := &Message{
		DescriptorProto: fds[0].MessageType[0],
		Fields: []*Field{
			{
				FieldDescriptorProto: fds[0].MessageType[0].Field[0],
			},
		},
	}
	boolMsg := &Message{
		DescriptorProto: fds[1].MessageType[0],
		Fields: []*Field{
			{
				FieldDescriptorProto: fds[1].MessageType[0].Field[0],
			},
		},
	}
	files := []*File{
		{
			FileDescriptorProto: fds[0],
			GoPkg: GoPackage{
				Path: "path/to/example.pb",
				Name: "example_pb",
			},
			Messages: []*Message{stringMsg},
			Services: []*Service{
				{
					ServiceDescriptorProto: fds[0].Service[0],
					Methods: []*Method{
						{
							MethodDescriptorProto: fds[0].Service[0].Method[0],
							PathTmpl:              compilePath(t, "/v1/example/to_s"),
							HTTPMethod:            "POST",
							RequestType:           boolMsg,
							ResponseType:          stringMsg,
							Body:                  &Body{FieldPath: nil},
						},
					},
				},
			},
		},
		{
			FileDescriptorProto: fds[1],
			GoPkg: GoPackage{
				Path: "path/to/another/example.pb",
				Name: "example_pb",
			},
			Messages: []*Message{boolMsg},
		},
	}

	for _, file := range files {
		crossLinkFixture(file)
	}
	testExtractServices(t, fds, "path/to/example.proto", files[0].Services)
}

func TestExtractServicesWithBodyPath(t *testing.T) {
	src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "OuterMessage"
			nested_type <
				name: "StringMessage"
				field <
					name: "string"
					number: 1
					label: LABEL_OPTIONAL
					type: TYPE_STRING
				>
			>
			field <
				name: "nested"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_MESSAGE
				type_name: "StringMessage"
			>
		>
		service <
			name: "ExampleService"
			method <
				name: "Echo"
				input_type: "OuterMessage"
				output_type: "OuterMessage"
				options <
					[google.api.http] <
						post: "/v1/example/echo"
						body: "nested"
					>
				>
			>
		>
	`
	var fd descriptor.FileDescriptorProto
	if err := proto.UnmarshalText(src, &fd); err != nil {
		t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
	}
	msg := &Message{
		DescriptorProto: fd.MessageType[0],
		Fields: []*Field{
			{
				FieldDescriptorProto: fd.MessageType[0].Field[0],
			},
		},
	}
	file := &File{
		FileDescriptorProto: &fd,
		GoPkg: GoPackage{
			Path: "path/to/example.pb",
			Name: "example_pb",
		},
		Messages: []*Message{msg},
		Services: []*Service{
			{
				ServiceDescriptorProto: fd.Service[0],
				Methods: []*Method{
					{
						MethodDescriptorProto: fd.Service[0].Method[0],
						PathTmpl:              compilePath(t, "/v1/example/echo"),
						HTTPMethod:            "POST",
						RequestType:           msg,
						ResponseType:          msg,
						Body: &Body{
							FieldPath: FieldPath{
								{
									Name:   "nested",
									Target: msg.Fields[0],
								},
							},
						},
					},
				},
			},
		},
	}

	crossLinkFixture(file)
	testExtractServices(t, []*descriptor.FileDescriptorProto{&fd}, "path/to/example.proto", file.Services)
}

func TestExtractServicesWithPathParam(t *testing.T) {
	src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "StringMessage"
			field <
				name: "string"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
			>
		>
		service <
			name: "ExampleService"
			method <
				name: "Echo"
				input_type: "StringMessage"
				output_type: "StringMessage"
				options <
					[google.api.http] <
						get: "/v1/example/echo/{string=*}"
					>
				>
			>
		>
	`
	var fd descriptor.FileDescriptorProto
	if err := proto.UnmarshalText(src, &fd); err != nil {
		t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
	}
	msg := &Message{
		DescriptorProto: fd.MessageType[0],
		Fields: []*Field{
			{
				FieldDescriptorProto: fd.MessageType[0].Field[0],
			},
		},
	}
	file := &File{
		FileDescriptorProto: &fd,
		GoPkg: GoPackage{
			Path: "path/to/example.pb",
			Name: "example_pb",
		},
		Messages: []*Message{msg},
		Services: []*Service{
			{
				ServiceDescriptorProto: fd.Service[0],
				Methods: []*Method{
					{
						MethodDescriptorProto: fd.Service[0].Method[0],
						PathTmpl:              compilePath(t, "/v1/example/echo/{string=*}"),
						HTTPMethod:            "GET",
						RequestType:           msg,
						ResponseType:          msg,
						PathParams: []Parameter{
							{
								FieldPath: FieldPath{
									{
										Name:   "string",
										Target: msg.Fields[0],
									},
								},
								Target: msg.Fields[0],
							},
						},
					},
				},
			},
		},
	}

	crossLinkFixture(file)
	testExtractServices(t, []*descriptor.FileDescriptorProto{&fd}, "path/to/example.proto", file.Services)
}

func TestExtractServicesWithError(t *testing.T) {
	for _, spec := range []struct {
		target string
		srcs   []string
	}{
		{
			target: "path/to/example.proto",
			srcs: []string{
				// message not found
				`
					name: "path/to/example.proto",
					package: "example"
					service <
						name: "ExampleService"
						method <
							name: "Echo"
							input_type: "StringMessage"
							output_type: "StringMessage"
							options <
								[google.api.http] <
									post: "/v1/example/echo"
									body: "*"
								>
							>
						>
					>
				`,
			},
		},
		// body field path not resolved
		{
			target: "path/to/example.proto",
			srcs: []string{`
						name: "path/to/example.proto",
						package: "example"
						message_type <
							name: "StringMessage"
							field <
								name: "string"
								number: 1
								label: LABEL_OPTIONAL
								type: TYPE_STRING
							>
						>
						service <
							name: "ExampleService"
							method <
								name: "Echo"
								input_type: "StringMessage"
								output_type: "StringMessage"
								options <
									[google.api.http] <
										post: "/v1/example/echo"
										body: "bool"
									>
								>
							>
						>`,
			},
		},
		// param field path not resolved
		{
			target: "path/to/example.proto",
			srcs: []string{
				`
					name: "path/to/example.proto",
					package: "example"
					message_type <
						name: "StringMessage"
						field <
							name: "string"
							number: 1
							label: LABEL_OPTIONAL
							type: TYPE_STRING
						>
					>
					service <
						name: "ExampleService"
						method <
							name: "Echo"
							input_type: "StringMessage"
							output_type: "StringMessage"
							options <
								[google.api.http] <
									post: "/v1/example/echo/{bool=*}"
								>
							>
						>
					>
				`,
			},
		},
		// non aggregate type on field path
		{
			target: "path/to/example.proto",
			srcs: []string{
				`
					name: "path/to/example.proto",
					package: "example"
					message_type <
						name: "OuterMessage"
						field <
							name: "mid"
							number: 1
							label: LABEL_OPTIONAL
							type: TYPE_STRING
						>
						field <
							name: "bool"
							number: 2
							label: LABEL_OPTIONAL
							type: TYPE_BOOL
						>
					>
					service <
						name: "ExampleService"
						method <
							name: "Echo"
							input_type: "OuterMessage"
							output_type: "OuterMessage"
							options <
								[google.api.http] <
									post: "/v1/example/echo/{mid.bool=*}"
								>
							>
						>
					>
				`,
			},
		},
		// path param in client streaming
		{
			target: "path/to/example.proto",
			srcs: []string{
				`
					name: "path/to/example.proto",
					package: "example"
					message_type <
						name: "StringMessage"
						field <
							name: "string"
							number: 1
							label: LABEL_OPTIONAL
							type: TYPE_STRING
						>
					>
					service <
						name: "ExampleService"
						method <
							name: "Echo"
							input_type: "StringMessage"
							output_type: "StringMessage"
							options <
								[google.api.http] <
									post: "/v1/example/echo/{bool=*}"
								>
							>
							client_streaming: true
						>
					>
				`,
			},
		},
		// body for GET
		{
			target: "path/to/example.proto",
			srcs: []string{
				`
					name: "path/to/example.proto",
					package: "example"
					message_type <
						name: "StringMessage"
						field <
							name: "string"
							number: 1
							label: LABEL_OPTIONAL
							type: TYPE_STRING
						>
					>
					service <
						name: "ExampleService"
						method <
							name: "Echo"
							input_type: "StringMessage"
							output_type: "StringMessage"
							options <
								[google.api.http] <
									get: "/v1/example/echo"
									body: "string"
								>
							>
						>
					>
				`,
			},
		},
		// body for DELETE
		{
			target: "path/to/example.proto",
			srcs: []string{
				`
					name: "path/to/example.proto",
					package: "example"
					message_type <
						name: "StringMessage"
						field <
							name: "string"
							number: 1
							label: LABEL_OPTIONAL
							type: TYPE_STRING
						>
					>
					service <
						name: "ExampleService"
						method <
							name: "RemoveResource"
							input_type: "StringMessage"
							output_type: "StringMessage"
							options <
								[google.api.http] <
									delete: "/v1/example/resource"
									body: "string"
								>
							>
						>
					>
				`,
			},
		},
		// no pattern specified
		{
			target: "path/to/example.proto",
			srcs: []string{
				`
					name: "path/to/example.proto",
					package: "example"
					service <
						name: "ExampleService"
						method <
							name: "RemoveResource"
							input_type: "StringMessage"
							output_type: "StringMessage"
							options <
								[google.api.http] <
									body: "string"
								>
							>
						>
					>
				`,
			},
		},
		// unsupported path parameter type
		{
			target: "path/to/example.proto",
			srcs: []string{`
					name: "path/to/example.proto",
					package: "example"
					message_type <
						name: "OuterMessage"
						nested_type <
							name: "StringMessage"
							field <
								name: "value"
								number: 1
								label: LABEL_OPTIONAL
								type: TYPE_STRING
							>
						>
						field <
							name: "string"
							number: 1
							label: LABEL_OPTIONAL
							type: TYPE_MESSAGE
							type_name: "StringMessage"
						>
					>
					service <
						name: "ExampleService"
						method <
							name: "Echo"
							input_type: "OuterMessage"
							output_type: "OuterMessage"
							options <
								[google.api.http] <
									get: "/v1/example/echo/{string=*}"
								>
							>
						>
					>
				`,
			},
		},
	} {
		reg := NewRegistry()

		var fds []*descriptor.FileDescriptorProto
		for _, src := range spec.srcs {
			var fd descriptor.FileDescriptorProto
			if err := proto.UnmarshalText(src, &fd); err != nil {
				t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
			}
			reg.loadFile(&fd)
			fds = append(fds, &fd)
		}
		err := reg.loadServices(spec.target)
		if err == nil {
			t.Errorf("loadServices(%q) succeeded; want an error; files=%v", spec.target, spec.srcs)
		}
		t.Log(err)
	}
}

func TestResolveFieldPath(t *testing.T) {
	for _, spec := range []struct {
		src     string
		path    string
		wantErr bool
	}{
		{
			src: `
				name: 'example.proto'
				package: 'example'
				message_type <
					name: 'ExampleMessage'
					field <
						name: 'string'
						type: TYPE_STRING
						label: LABEL_OPTIONAL
						number: 1
					>
				>
			`,
			path:    "string",
			wantErr: false,
		},
		// no such field
		{
			src: `
				name: 'example.proto'
				package: 'example'
				message_type <
					name: 'ExampleMessage'
					field <
						name: 'string'
						type: TYPE_STRING
						label: LABEL_OPTIONAL
						number: 1
					>
				>
			`,
			path:    "something_else",
			wantErr: true,
		},
		// repeated field
		{
			src: `
				name: 'example.proto'
				package: 'example'
				message_type <
					name: 'ExampleMessage'
					field <
						name: 'string'
						type: TYPE_STRING
						label: LABEL_REPEATED
						number: 1
					>
				>
			`,
			path:    "string",
			wantErr: true,
		},
		// nested field
		{
			src: `
				name: 'example.proto'
				package: 'example'
				message_type <
					name: 'ExampleMessage'
					field <
						name: 'nested'
						type: TYPE_MESSAGE
						type_name: 'AnotherMessage'
						label: LABEL_OPTIONAL
						number: 1
					>
					field <
						name: 'terminal'
						type: TYPE_BOOL
						label: LABEL_OPTIONAL
						number: 2
					>
				>
				message_type <
					name: 'AnotherMessage'
					field <
						name: 'nested2'
						type: TYPE_MESSAGE
						type_name: 'ExampleMessage'
						label: LABEL_OPTIONAL
						number: 1
					>
				>
			`,
			path:    "nested.nested2.nested.nested2.nested.nested2.terminal",
			wantErr: false,
		},
		// non aggregate field on the path
		{
			src: `
				name: 'example.proto'
				package: 'example'
				message_type <
					name: 'ExampleMessage'
					field <
						name: 'nested'
						type: TYPE_MESSAGE
						type_name: 'AnotherMessage'
						label: LABEL_OPTIONAL
						number: 1
					>
					field <
						name: 'terminal'
						type: TYPE_BOOL
						label: LABEL_OPTIONAL
						number: 2
					>
				>
				message_type <
					name: 'AnotherMessage'
					field <
						name: 'nested2'
						type: TYPE_MESSAGE
						type_name: 'ExampleMessage'
						label: LABEL_OPTIONAL
						number: 1
					>
				>
			`,
			path:    "nested.terminal.nested2",
			wantErr: true,
		},
		// repeated field
		{
			src: `
				name: 'example.proto'
				package: 'example'
				message_type <
					name: 'ExampleMessage'
					field <
						name: 'nested'
						type: TYPE_MESSAGE
						type_name: 'AnotherMessage'
						label: LABEL_OPTIONAL
						number: 1
					>
					field <
						name: 'terminal'
						type: TYPE_BOOL
						label: LABEL_OPTIONAL
						number: 2
					>
				>
				message_type <
					name: 'AnotherMessage'
					field <
						name: 'nested2'
						type: TYPE_MESSAGE
						type_name: 'ExampleMessage'
						label: LABEL_REPEATED
						number: 1
					>
				>
			`,
			path:    "nested.nested2.terminal",
			wantErr: true,
		},
	} {
		var file descriptor.FileDescriptorProto
		if err := proto.UnmarshalText(spec.src, &file); err != nil {
			t.Fatalf("proto.Unmarshal(%s) failed with %v; want success", spec.src, err)
		}
		reg := NewRegistry()
		reg.loadFile(&file)
		f, err := reg.LookupFile(file.GetName())
		if err != nil {
			t.Fatalf("reg.LookupFile(%q) failed with %v; want success; on file=%s", file.GetName(), err, spec.src)
		}
		_, err = reg.resolveFiledPath(f.Messages[0], spec.path)
		if got, want := err != nil, spec.wantErr; got != want {
			if want {
				t.Errorf("reg.resolveFiledPath(%q, %q) succeeded; want an error", f.Messages[0].GetName(), spec.path)
				continue
			}
			t.Errorf("reg.resolveFiledPath(%q, %q) failed with %v; want success", f.Messages[0].GetName(), spec.path, err)
		}
	}
}
