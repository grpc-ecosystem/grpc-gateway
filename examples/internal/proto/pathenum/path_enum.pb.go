// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.0
// source: examples/internal/proto/pathenum/path_enum.proto

package pathenum

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PathEnum int32

const (
	PathEnum_ABC PathEnum = 0
	PathEnum_DEF PathEnum = 1
)

// Enum value maps for PathEnum.
var (
	PathEnum_name = map[int32]string{
		0: "ABC",
		1: "DEF",
	}
	PathEnum_value = map[string]int32{
		"ABC": 0,
		"DEF": 1,
	}
)

func (x PathEnum) Enum() *PathEnum {
	p := new(PathEnum)
	*p = x
	return p
}

func (x PathEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PathEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_examples_internal_proto_pathenum_path_enum_proto_enumTypes[0].Descriptor()
}

func (PathEnum) Type() protoreflect.EnumType {
	return &file_examples_internal_proto_pathenum_path_enum_proto_enumTypes[0]
}

func (x PathEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PathEnum.Descriptor instead.
func (PathEnum) EnumDescriptor() ([]byte, []int) {
	return file_examples_internal_proto_pathenum_path_enum_proto_rawDescGZIP(), []int{0}
}

type MessagePathEnum_NestedPathEnum int32

const (
	MessagePathEnum_GHI MessagePathEnum_NestedPathEnum = 0
	MessagePathEnum_JKL MessagePathEnum_NestedPathEnum = 1
)

// Enum value maps for MessagePathEnum_NestedPathEnum.
var (
	MessagePathEnum_NestedPathEnum_name = map[int32]string{
		0: "GHI",
		1: "JKL",
	}
	MessagePathEnum_NestedPathEnum_value = map[string]int32{
		"GHI": 0,
		"JKL": 1,
	}
)

func (x MessagePathEnum_NestedPathEnum) Enum() *MessagePathEnum_NestedPathEnum {
	p := new(MessagePathEnum_NestedPathEnum)
	*p = x
	return p
}

func (x MessagePathEnum_NestedPathEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessagePathEnum_NestedPathEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_examples_internal_proto_pathenum_path_enum_proto_enumTypes[1].Descriptor()
}

func (MessagePathEnum_NestedPathEnum) Type() protoreflect.EnumType {
	return &file_examples_internal_proto_pathenum_path_enum_proto_enumTypes[1]
}

func (x MessagePathEnum_NestedPathEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessagePathEnum_NestedPathEnum.Descriptor instead.
func (MessagePathEnum_NestedPathEnum) EnumDescriptor() ([]byte, []int) {
	return file_examples_internal_proto_pathenum_path_enum_proto_rawDescGZIP(), []int{0, 0}
}

type MessagePathEnum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MessagePathEnum) Reset() {
	*x = MessagePathEnum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_internal_proto_pathenum_path_enum_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessagePathEnum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessagePathEnum) ProtoMessage() {}

func (x *MessagePathEnum) ProtoReflect() protoreflect.Message {
	mi := &file_examples_internal_proto_pathenum_path_enum_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessagePathEnum.ProtoReflect.Descriptor instead.
func (*MessagePathEnum) Descriptor() ([]byte, []int) {
	return file_examples_internal_proto_pathenum_path_enum_proto_rawDescGZIP(), []int{0}
}

var File_examples_internal_proto_pathenum_path_enum_proto protoreflect.FileDescriptor

var file_examples_internal_proto_pathenum_path_enum_proto_rawDesc = []byte{
	0x0a, 0x30, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x65, 0x6e,
	0x75, 0x6d, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x27, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2e, 0x70, 0x61, 0x74, 0x68, 0x65, 0x6e, 0x75, 0x6d, 0x22, 0x35, 0x0a, 0x0f, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x74, 0x68, 0x45, 0x6e, 0x75, 0x6d, 0x22, 0x22,
	0x0a, 0x0e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x50, 0x61, 0x74, 0x68, 0x45, 0x6e, 0x75, 0x6d,
	0x12, 0x07, 0x0a, 0x03, 0x47, 0x48, 0x49, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x4a, 0x4b, 0x4c,
	0x10, 0x01, 0x2a, 0x1c, 0x0a, 0x08, 0x50, 0x61, 0x74, 0x68, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x07,
	0x0a, 0x03, 0x41, 0x42, 0x43, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x45, 0x46, 0x10, 0x01,
	0x42, 0x4c, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2d, 0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x76, 0x32, 0x2f, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x65, 0x6e, 0x75, 0x6d, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examples_internal_proto_pathenum_path_enum_proto_rawDescOnce sync.Once
	file_examples_internal_proto_pathenum_path_enum_proto_rawDescData = file_examples_internal_proto_pathenum_path_enum_proto_rawDesc
)

func file_examples_internal_proto_pathenum_path_enum_proto_rawDescGZIP() []byte {
	file_examples_internal_proto_pathenum_path_enum_proto_rawDescOnce.Do(func() {
		file_examples_internal_proto_pathenum_path_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_examples_internal_proto_pathenum_path_enum_proto_rawDescData)
	})
	return file_examples_internal_proto_pathenum_path_enum_proto_rawDescData
}

var file_examples_internal_proto_pathenum_path_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_examples_internal_proto_pathenum_path_enum_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_examples_internal_proto_pathenum_path_enum_proto_goTypes = []interface{}{
	(PathEnum)(0),                       // 0: grpc.gateway.examples.internal.pathenum.PathEnum
	(MessagePathEnum_NestedPathEnum)(0), // 1: grpc.gateway.examples.internal.pathenum.MessagePathEnum.NestedPathEnum
	(*MessagePathEnum)(nil),             // 2: grpc.gateway.examples.internal.pathenum.MessagePathEnum
}
var file_examples_internal_proto_pathenum_path_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_examples_internal_proto_pathenum_path_enum_proto_init() }
func file_examples_internal_proto_pathenum_path_enum_proto_init() {
	if File_examples_internal_proto_pathenum_path_enum_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examples_internal_proto_pathenum_path_enum_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessagePathEnum); i {
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
			RawDescriptor: file_examples_internal_proto_pathenum_path_enum_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_examples_internal_proto_pathenum_path_enum_proto_goTypes,
		DependencyIndexes: file_examples_internal_proto_pathenum_path_enum_proto_depIdxs,
		EnumInfos:         file_examples_internal_proto_pathenum_path_enum_proto_enumTypes,
		MessageInfos:      file_examples_internal_proto_pathenum_path_enum_proto_msgTypes,
	}.Build()
	File_examples_internal_proto_pathenum_path_enum_proto = out.File
	file_examples_internal_proto_pathenum_path_enum_proto_rawDesc = nil
	file_examples_internal_proto_pathenum_path_enum_proto_goTypes = nil
	file_examples_internal_proto_pathenum_path_enum_proto_depIdxs = nil
}
