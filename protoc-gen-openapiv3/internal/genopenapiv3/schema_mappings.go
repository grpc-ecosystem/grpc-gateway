package genopenapiv3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"google.golang.org/protobuf/types/descriptorpb"
)

var primitiveTypeSchemas = map[descriptorpb.FieldDescriptorProto_Type]*openapi3.Schema{
	descriptorpb.FieldDescriptorProto_TYPE_DOUBLE: {
		Type:   &openapi3.Types{"number"},
		Format: "double",
	},
	descriptorpb.FieldDescriptorProto_TYPE_FLOAT: {
		Type:   &openapi3.Types{"number"},
		Format: "float",
	},
	// 64bit integer types are marshaled as string in the default JSONPb marshaler.
	// This maintains compatibility with JSON's limited number precision.
	descriptorpb.FieldDescriptorProto_TYPE_INT64: {
		Type:   &openapi3.Types{"string"},
		Format: "int64",
	},
	// 64bit integer types are marshaled as string in the default JSONPb marshaler.
	// TODO(yugui) Add an option to declare 64bit integers as int64.
	//
	// NOTE: uint64 is not a standard format in OpenAPI spec.
	// So we cannot expect that uint64 is commonly supported by OpenAPI processors.
	descriptorpb.FieldDescriptorProto_TYPE_UINT64: {
		Type:   &openapi3.Types{"string"},
		Format: "uint64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_INT32: {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	// 64bit types marshaled as string for JSON compatibility
	descriptorpb.FieldDescriptorProto_TYPE_FIXED64: {
		Type:   &openapi3.Types{"string"},
		Format: "uint64",
	},
	// Fixed 32-bit unsigned integer
	descriptorpb.FieldDescriptorProto_TYPE_FIXED32: {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	// NOTE: In OpenAPI v3 specification, format should be empty on boolean type
	descriptorpb.FieldDescriptorProto_TYPE_BOOL: {
		Type: &openapi3.Types{"boolean"},
	},
	// NOTE: In OpenAPI v3 specification, format can be empty on string type
	descriptorpb.FieldDescriptorProto_TYPE_STRING: {
		Type: &openapi3.Types{"string"},
	},
	// Base64 encoded string representation
	descriptorpb.FieldDescriptorProto_TYPE_BYTES: {
		Type:   &openapi3.Types{"string"},
		Format: "byte",
	},
	// 32-bit unsigned integer
	descriptorpb.FieldDescriptorProto_TYPE_UINT32: {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	// 64bit types marshaled as string for JSON compatibility
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: {
		Type:   &openapi3.Types{"string"},
		Format: "int64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_SINT32: {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	// 64bit types marshaled as string for JSON compatibility
	descriptorpb.FieldDescriptorProto_TYPE_SINT64: {
		Type:   &openapi3.Types{"string"},
		Format: "int64",
	},
}

var wktSchemas = map[string]*openapi3.Schema{
	".google.protobuf.FieldMask": {
		Type: &openapi3.Types{"string"},
	},
	".google.protobuf.Timestamp": {
		Type:   &openapi3.Types{"string"},
		Format: "date-time",
	},
	".google.protobuf.Duration": {
		Type: &openapi3.Types{"string"},
	},
	".google.protobuf.StringValue": {
		Type: &openapi3.Types{"string"},
	},
	".google.protobuf.BytesValue": {
		Type:   &openapi3.Types{"string"},
		Format: "byte",
	},
	".google.protobuf.Int32Value": {
		Type:   &openapi3.Types{"integer"},
		Format: "int32",
	},
	".google.protobuf.UInt32Value": {
		Type:   &openapi3.Types{"integer"},
		Format: "int64",
	},
	".google.protobuf.Int64Value": {
		Type:   &openapi3.Types{"string"},
		Format: "int64",
	},
	".google.protobuf.UInt64Value": {
		Type:   &openapi3.Types{"string"},
		Format: "uint64",
	},
	".google.protobuf.FloatValue": {
		Type:   &openapi3.Types{"number"},
		Format: "float",
	},
	".google.protobuf.DoubleValue": {
		Type:   &openapi3.Types{"number"},
		Format: "double",
	},
	".google.protobuf.BoolValue": {
		Type: &openapi3.Types{"boolean"},
	},
	".google.protobuf.Empty": {
		Type: &openapi3.Types{"object"},
	},
	".google.protobuf.Struct": {
		Type: &openapi3.Types{"object"},
	},
	".google.protobuf.Value": {},
	".google.protobuf.ListValue": {
		Type: &openapi3.Types{"array"},
		Items: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"object"},
			},
		},
	},
	".google.protobuf.NullValue": {
		Type: &openapi3.Types{"string"},
	},
}
