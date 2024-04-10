// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: examples/internal/proto/examplepb/enum_with_single_value.proto

package examplepb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// EnumWithSingleValue is an enum with a single value. Since it has just a single value
// the `enum` field should be omitted in the generated OpenAPI spec for the type when
// the omit_enum_default_value option is set to true.
type EnumWithSingleValue int32

const (
	EnumWithSingleValue_ONLY_VALUE EnumWithSingleValue = 0
)

// Enum value maps for EnumWithSingleValue.
var (
	EnumWithSingleValue_name = map[int32]string{
		0: "ONLY_VALUE",
	}
	EnumWithSingleValue_value = map[string]int32{
		"ONLY_VALUE": 0,
	}
)

func (x EnumWithSingleValue) Enum() *EnumWithSingleValue {
	p := new(EnumWithSingleValue)
	*p = x
	return p
}

func (x EnumWithSingleValue) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EnumWithSingleValue) Descriptor() protoreflect.EnumDescriptor {
	return file_examples_internal_proto_examplepb_enum_with_single_value_proto_enumTypes[0].Descriptor()
}

func (EnumWithSingleValue) Type() protoreflect.EnumType {
	return &file_examples_internal_proto_examplepb_enum_with_single_value_proto_enumTypes[0]
}

func (x EnumWithSingleValue) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EnumWithSingleValue.Descriptor instead.
func (EnumWithSingleValue) EnumDescriptor() ([]byte, []int) {
	return file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescGZIP(), []int{0}
}

type EnumWithSingleValueEchoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value EnumWithSingleValue `protobuf:"varint,1,opt,name=value,proto3,enum=grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValue" json:"value,omitempty"`
}

func (x *EnumWithSingleValueEchoRequest) Reset() {
	*x = EnumWithSingleValueEchoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumWithSingleValueEchoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumWithSingleValueEchoRequest) ProtoMessage() {}

func (x *EnumWithSingleValueEchoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumWithSingleValueEchoRequest.ProtoReflect.Descriptor instead.
func (*EnumWithSingleValueEchoRequest) Descriptor() ([]byte, []int) {
	return file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescGZIP(), []int{0}
}

func (x *EnumWithSingleValueEchoRequest) GetValue() EnumWithSingleValue {
	if x != nil {
		return x.Value
	}
	return EnumWithSingleValue_ONLY_VALUE
}

type EnumWithSingleValueEchoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EnumWithSingleValueEchoResponse) Reset() {
	*x = EnumWithSingleValueEchoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumWithSingleValueEchoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumWithSingleValueEchoResponse) ProtoMessage() {}

func (x *EnumWithSingleValueEchoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumWithSingleValueEchoResponse.ProtoReflect.Descriptor instead.
func (*EnumWithSingleValueEchoResponse) Descriptor() ([]byte, []int) {
	return file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescGZIP(), []int{1}
}

var File_examples_internal_proto_examplepb_enum_with_single_value_proto protoreflect.FileDescriptor

var file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDesc = []byte{
	0x0a, 0x3e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x70, 0x62, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x73, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61,
	0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b,
	0x0a, 0x1e, 0x45, 0x6e, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x59, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x43, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62,
	0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a, 0x1f, 0x45,
	0x6e, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x25,
	0x0a, 0x13, 0x45, 0x6e, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4f, 0x4e, 0x4c, 0x59, 0x5f, 0x56, 0x41,
	0x4c, 0x55, 0x45, 0x10, 0x00, 0x32, 0xfa, 0x01, 0x0a, 0x1a, 0x45, 0x6e, 0x75, 0x6d, 0x57, 0x69,
	0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0xdb, 0x01, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x4e, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x45,
	0x6e, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x4f, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x45,
	0x6e, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c, 0x3a, 0x01, 0x2a, 0x22, 0x27, 0x2f, 0x76, 0x31, 0x2f, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x2d, 0x77, 0x69, 0x74, 0x68,
	0x2d, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x2d, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2f, 0x65, 0x63,
	0x68, 0x6f, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x76, 0x32, 0x2f,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescOnce sync.Once
	file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescData = file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDesc
)

func file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescGZIP() []byte {
	file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescOnce.Do(func() {
		file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescData = protoimpl.X.CompressGZIP(file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescData)
	})
	return file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDescData
}

var file_examples_internal_proto_examplepb_enum_with_single_value_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_examples_internal_proto_examplepb_enum_with_single_value_proto_goTypes = []interface{}{
	(EnumWithSingleValue)(0),                // 0: grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValue
	(*EnumWithSingleValueEchoRequest)(nil),  // 1: grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueEchoRequest
	(*EnumWithSingleValueEchoResponse)(nil), // 2: grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueEchoResponse
}
var file_examples_internal_proto_examplepb_enum_with_single_value_proto_depIdxs = []int32{
	0, // 0: grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueEchoRequest.value:type_name -> grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValue
	1, // 1: grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueService.Echo:input_type -> grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueEchoRequest
	2, // 2: grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueService.Echo:output_type -> grpc.gateway.examples.internal.proto.examplepb.EnumWithSingleValueEchoResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_examples_internal_proto_examplepb_enum_with_single_value_proto_init() }
func file_examples_internal_proto_examplepb_enum_with_single_value_proto_init() {
	if File_examples_internal_proto_examplepb_enum_with_single_value_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumWithSingleValueEchoRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumWithSingleValueEchoResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examples_internal_proto_examplepb_enum_with_single_value_proto_goTypes,
		DependencyIndexes: file_examples_internal_proto_examplepb_enum_with_single_value_proto_depIdxs,
		EnumInfos:         file_examples_internal_proto_examplepb_enum_with_single_value_proto_enumTypes,
		MessageInfos:      file_examples_internal_proto_examplepb_enum_with_single_value_proto_msgTypes,
	}.Build()
	File_examples_internal_proto_examplepb_enum_with_single_value_proto = out.File
	file_examples_internal_proto_examplepb_enum_with_single_value_proto_rawDesc = nil
	file_examples_internal_proto_examplepb_enum_with_single_value_proto_goTypes = nil
	file_examples_internal_proto_examplepb_enum_with_single_value_proto_depIdxs = nil
}
