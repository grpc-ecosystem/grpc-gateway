package descriptor

import (
	"fmt"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/casing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/httprule"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// IsWellKnownType returns true if the provided fully qualified type name is considered 'well-known'.
func IsWellKnownType(typeName string) bool {
	_, ok := wellKnownTypeConv[typeName]
	return ok
}

// GoPackage represents a golang package
type GoPackage struct {
	// Path is the package path to the package.
	Path string
	// Name is the package name of the package
	Name string
	// Alias is an alias of the package unique within the current invokation of grpc-gateway generator.
	Alias string
}

// Standard returns whether the import is a golang standard package.
func (p GoPackage) Standard() bool {
	return !strings.Contains(p.Path, ".")
}

// String returns a string representation of this package in the form of import line in golang.
func (p GoPackage) String() string {
	if p.Alias == "" {
		return fmt.Sprintf("%q", p.Path)
	}
	return fmt.Sprintf("%s %q", p.Alias, p.Path)
}

// File wraps descriptorpb.FileDescriptorProto for richer features.
type File struct {
	*descriptorpb.FileDescriptorProto
	// GoPkg is the go package of the go file generated from this file.
	GoPkg GoPackage
	// Messages is the list of messages defined in this file.
	Messages []*Message
	// Enums is the list of enums defined in this file.
	Enums []*Enum
	// Services is the list of services defined in this file.
	Services []*Service
}

// Pkg returns package name or alias if it's present
func (f *File) Pkg() string {
	pkg := f.GoPkg.Name
	if alias := f.GoPkg.Alias; alias != "" {
		pkg = alias
	}
	return pkg
}

// proto2 determines if the syntax of the file is proto2.
func (f *File) proto2() bool {
	return f.Syntax == nil || f.GetSyntax() == "proto2"
}

// ResponseFile wraps pluginpb.CodeGeneratorResponse_File.
type ResponseFile struct {
	*pluginpb.CodeGeneratorResponse_File

	// GoPkg is the Go package of the generated file.
	GoPkg GoPackage
}

// Message describes a protocol buffer message types
type Message struct {
	// File is the file where the message is defined
	File *File
	// Outers is a list of outer messages if this message is a nested type.
	Outers []string
	*descriptorpb.DescriptorProto
	Fields []*Field

	// Index is proto path index of this message in File.
	Index int

	ForcePrefixedName bool
}

// FQMN returns a fully qualified message name of this message.
func (m *Message) FQMN() string {
	components := []string{""}
	if m.File.Package != nil {
		components = append(components, m.File.GetPackage())
	}
	components = append(components, m.Outers...)
	components = append(components, m.GetName())
	return strings.Join(components, ".")
}

// GoType returns a go type name for the message type.
// It prefixes the type name with the package alias if
// its belonging package is not "currentPackage".
func (m *Message) GoType(currentPackage string) string {
	var components []string
	components = append(components, m.Outers...)
	components = append(components, m.GetName())

	name := strings.Join(components, "_")
	if !m.ForcePrefixedName && m.File.GoPkg.Path == currentPackage {
		return name
	}
	return fmt.Sprintf("%s.%s", m.File.Pkg(), name)
}

// Enum describes a protocol buffer enum types
type Enum struct {
	// File is the file where the enum is defined
	File *File
	// Outers is a list of outer messages if this enum is a nested type.
	Outers []string
	*descriptorpb.EnumDescriptorProto

	Index int

	ForcePrefixedName bool
}

// FQEN returns a fully qualified enum name of this enum.
func (e *Enum) FQEN() string {
	components := []string{""}
	if e.File.Package != nil {
		components = append(components, e.File.GetPackage())
	}
	components = append(components, e.Outers...)
	components = append(components, e.GetName())
	return strings.Join(components, ".")
}

// GoType returns a go type name for the enum type.
// It prefixes the type name with the package alias if
// its belonging package is not "currentPackage".
func (e *Enum) GoType(currentPackage string) string {
	var components []string
	components = append(components, e.Outers...)
	components = append(components, e.GetName())

	name := strings.Join(components, "_")
	if !e.ForcePrefixedName && e.File.GoPkg.Path == currentPackage {
		return name
	}
	return fmt.Sprintf("%s.%s", e.File.Pkg(), name)
}

// Service wraps descriptorpb.ServiceDescriptorProto for richer features.
type Service struct {
	// File is the file where this service is defined.
	File *File
	*descriptorpb.ServiceDescriptorProto
	// Methods is the list of methods defined in this service.
	Methods []*Method

	ForcePrefixedName bool
}

// FQSN returns the fully qualified service name of this service.
func (s *Service) FQSN() string {
	components := []string{""}
	if s.File.Package != nil {
		components = append(components, s.File.GetPackage())
	}
	components = append(components, s.GetName())
	return strings.Join(components, ".")
}

// InstanceName returns object name of the service with package prefix if needed
func (s *Service) InstanceName() string {
	if !s.ForcePrefixedName {
		return s.GetName()
	}
	return fmt.Sprintf("%s.%s", s.File.Pkg(), s.GetName())
}

// ClientConstructorName returns name of the Client constructor with package prefix if needed
func (s *Service) ClientConstructorName() string {
	constructor := "New" + s.GetName() + "Client"
	if !s.ForcePrefixedName {
		return constructor
	}
	return fmt.Sprintf("%s.%s", s.File.Pkg(), constructor)
}

// Method wraps descriptorpb.MethodDescriptorProto for richer features.
type Method struct {
	// Service is the service which this method belongs to.
	Service *Service
	*descriptorpb.MethodDescriptorProto

	// RequestType is the message type of requests to this method.
	RequestType *Message
	// ResponseType is the message type of responses from this method.
	ResponseType *Message
	Bindings     []*Binding
}

// FQMN returns a fully qualified rpc method name of this method.
func (m *Method) FQMN() string {
	components := []string{}
	components = append(components, m.Service.FQSN())
	components = append(components, m.GetName())
	return strings.Join(components, ".")
}

// Binding describes how an HTTP endpoint is bound to a gRPC method.
type Binding struct {
	// Method is the method which the endpoint is bound to.
	Method *Method
	// Index is a zero-origin index of the binding in the target method
	Index int
	// PathTmpl is path template where this method is mapped to.
	PathTmpl httprule.Template
	// HTTPMethod is the HTTP method which this method is mapped to.
	HTTPMethod string
	// PathParams is the list of parameters provided in HTTP request paths.
	PathParams []Parameter
	// Body describes parameters provided in HTTP request body.
	Body *Body
	// ResponseBody describes field in response struct to marshal in HTTP response body.
	ResponseBody *Body
}

// ExplicitParams returns a list of explicitly bound parameters of "b",
// i.e. a union of field path for body and field paths for path parameters.
func (b *Binding) ExplicitParams() []string {
	var result []string
	if b.Body != nil {
		result = append(result, b.Body.FieldPath.String())
	}
	for _, p := range b.PathParams {
		result = append(result, p.FieldPath.String())
	}
	return result
}

// Field wraps descriptorpb.FieldDescriptorProto for richer features.
type Field struct {
	// Message is the message type which this field belongs to.
	Message *Message
	// FieldMessage is the message type of the field.
	FieldMessage *Message
	*descriptorpb.FieldDescriptorProto

	ForcePrefixedName bool
}

// FQFN returns a fully qualified field name of this field.
func (f *Field) FQFN() string {
	return strings.Join([]string{f.Message.FQMN(), f.GetName()}, ".")
}

// Parameter is a parameter provided in http requests
type Parameter struct {
	// FieldPath is a path to a proto field which this parameter is mapped to.
	FieldPath
	// Target is the proto field which this parameter is mapped to.
	Target *Field
	// Method is the method which this parameter is used for.
	Method *Method
}

// ConvertFuncExpr returns a go expression of a converter function.
// The converter function converts a string into a value for the parameter.
func (p Parameter) ConvertFuncExpr() (string, error) {
	tbl := proto3ConvertFuncs
	if !p.IsProto2() && p.IsRepeated() {
		tbl = proto3RepeatedConvertFuncs
	} else if p.IsProto2() && !p.IsRepeated() {
		tbl = proto2ConvertFuncs
	} else if p.IsProto2() && p.IsRepeated() {
		tbl = proto2RepeatedConvertFuncs
	}
	typ := p.Target.GetType()
	conv, ok := tbl[typ]
	if !ok {
		conv, ok = wellKnownTypeConv[p.Target.GetTypeName()]
	}
	if !ok {
		return "", fmt.Errorf("unsupported field type %s of parameter %s in %s.%s", typ, p.FieldPath, p.Method.Service.GetName(), p.Method.GetName())
	}
	return conv, nil
}

// IsEnum returns true if the field is an enum type, otherwise false is returned.
func (p Parameter) IsEnum() bool {
	return p.Target.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM
}

// IsRepeated returns true if the field is repeated, otherwise false is returned.
func (p Parameter) IsRepeated() bool {
	return p.Target.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}

// IsProto2 returns true if the field is proto2, otherwise false is returned.
func (p Parameter) IsProto2() bool {
	return p.Target.Message.File.proto2()
}

// Body describes a http (request|response) body to be sent to the (method|client).
// This is used in body and response_body options in google.api.HttpRule
type Body struct {
	// FieldPath is a path to a proto field which the (request|response) body is mapped to.
	// The (request|response) body is mapped to the (request|response) type itself if FieldPath is empty.
	FieldPath FieldPath
}

// AssignableExpr returns an assignable expression in Go to be used to initialize method request object.
// It starts with "msgExpr", which is the go expression of the method request object.
func (b Body) AssignableExpr(msgExpr string) string {
	return b.FieldPath.AssignableExpr(msgExpr)
}

// FieldPath is a path to a field from a request message.
type FieldPath []FieldPathComponent

// String returns a string representation of the field path.
func (p FieldPath) String() string {
	var components []string
	for _, c := range p {
		components = append(components, c.Name)
	}
	return strings.Join(components, ".")
}

// IsNestedProto3 indicates whether the FieldPath is a nested Proto3 path.
func (p FieldPath) IsNestedProto3() bool {
	if len(p) > 1 && !p[0].Target.Message.File.proto2() {
		return true
	}
	return false
}

// AssignableExpr is an assignable expression in Go to be used to assign a value to the target field.
// It starts with "msgExpr", which is the go expression of the method request object.
func (p FieldPath) AssignableExpr(msgExpr string) string {
	l := len(p)
	if l == 0 {
		return msgExpr
	}

	var preparations []string
	components := msgExpr
	for i, c := range p {
		// Check if it is a oneOf field.
		if c.Target.OneofIndex != nil {
			index := c.Target.OneofIndex
			msg := c.Target.Message
			oneOfName := casing.Camel(msg.GetOneofDecl()[*index].GetName())
			oneofFieldName := msg.GetName() + "_" + c.AssignableExpr()

			if c.Target.ForcePrefixedName {
				oneofFieldName = msg.File.Pkg() + "." + oneofFieldName
			}

			components = components + "." + oneOfName
			s := `if %s == nil {
				%s =&%s{}
			} else if _, ok := %s.(*%s); !ok {
				return nil, metadata, status.Errorf(codes.InvalidArgument, "expect type: *%s, but: %%t\n",%s)
			}`

			preparations = append(preparations, fmt.Sprintf(s, components, components, oneofFieldName, components, oneofFieldName, oneofFieldName, components))
			components = components + ".(*" + oneofFieldName + ")"
		}

		if i == l-1 {
			components = components + "." + c.AssignableExpr()
			continue
		}
		components = components + "." + c.ValueExpr()
	}

	preparations = append(preparations, components)
	return strings.Join(preparations, "\n")
}

// FieldPathComponent is a path component in FieldPath
type FieldPathComponent struct {
	// Name is a name of the proto field which this component corresponds to.
	// TODO(yugui) is this necessary?
	Name string
	// Target is the proto field which this component corresponds to.
	Target *Field
}

// AssignableExpr returns an assignable expression in go for this field.
func (c FieldPathComponent) AssignableExpr() string {
	return casing.Camel(c.Name)
}

// ValueExpr returns an expression in go for this field.
func (c FieldPathComponent) ValueExpr() string {
	if c.Target.Message.File.proto2() {
		return fmt.Sprintf("Get%s()", casing.Camel(c.Name))
	}
	return casing.Camel(c.Name)
}

var (
	proto3ConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.Bool",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.String",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:    "runtime.Bytes",
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.Enum",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64",
	}

	proto3RepeatedConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolSlice",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.StringSlice",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:    "runtime.BytesSlice",
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumSlice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64Slice",
	}

	proto2ConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64P",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32P",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64P",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64P",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32P",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64P",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32P",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolP",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.StringP",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(yugui) Handle bytes
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32P",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumP",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32P",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64P",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32P",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64P",
	}

	proto2RepeatedConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolSlice",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.StringSlice",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(maros7) Handle bytes
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumSlice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64Slice",
	}

	wellKnownTypeConv = map[string]string{
		".google.protobuf.Timestamp":   "runtime.Timestamp",
		".google.protobuf.Duration":    "runtime.Duration",
		".google.protobuf.StringValue": "runtime.StringValue",
		".google.protobuf.FloatValue":  "runtime.FloatValue",
		".google.protobuf.DoubleValue": "runtime.DoubleValue",
		".google.protobuf.BoolValue":   "runtime.BoolValue",
		".google.protobuf.BytesValue":  "runtime.BytesValue",
		".google.protobuf.Int32Value":  "runtime.Int32Value",
		".google.protobuf.UInt32Value": "runtime.UInt32Value",
		".google.protobuf.Int64Value":  "runtime.Int64Value",
		".google.protobuf.UInt64Value": "runtime.UInt64Value",
	}
)
