// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.0
// source: examples/internal/helloworld/helloworld.proto

package helloworld

import (
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	StrVal    *wrappers.StringValue `protobuf:"bytes,2,opt,name=strVal,proto3" json:"strVal,omitempty"`
	FloatVal  *wrappers.FloatValue  `protobuf:"bytes,3,opt,name=floatVal,proto3" json:"floatVal,omitempty"`
	DoubleVal *wrappers.DoubleValue `protobuf:"bytes,4,opt,name=doubleVal,proto3" json:"doubleVal,omitempty"`
	BoolVal   *wrappers.BoolValue   `protobuf:"bytes,5,opt,name=boolVal,proto3" json:"boolVal,omitempty"`
	BytesVal  *wrappers.BytesValue  `protobuf:"bytes,6,opt,name=bytesVal,proto3" json:"bytesVal,omitempty"`
	Int32Val  *wrappers.Int32Value  `protobuf:"bytes,7,opt,name=int32Val,proto3" json:"int32Val,omitempty"`
	Uint32Val *wrappers.UInt32Value `protobuf:"bytes,8,opt,name=uint32Val,proto3" json:"uint32Val,omitempty"`
	Int64Val  *wrappers.Int64Value  `protobuf:"bytes,9,opt,name=int64Val,proto3" json:"int64Val,omitempty"`
	Uint64Val *wrappers.UInt64Value `protobuf:"bytes,10,opt,name=uint64Val,proto3" json:"uint64Val,omitempty"`
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_internal_helloworld_helloworld_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_examples_internal_helloworld_helloworld_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_examples_internal_helloworld_helloworld_proto_rawDescGZIP(), []int{0}
}

func (x *HelloRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HelloRequest) GetStrVal() *wrappers.StringValue {
	if x != nil {
		return x.StrVal
	}
	return nil
}

func (x *HelloRequest) GetFloatVal() *wrappers.FloatValue {
	if x != nil {
		return x.FloatVal
	}
	return nil
}

func (x *HelloRequest) GetDoubleVal() *wrappers.DoubleValue {
	if x != nil {
		return x.DoubleVal
	}
	return nil
}

func (x *HelloRequest) GetBoolVal() *wrappers.BoolValue {
	if x != nil {
		return x.BoolVal
	}
	return nil
}

func (x *HelloRequest) GetBytesVal() *wrappers.BytesValue {
	if x != nil {
		return x.BytesVal
	}
	return nil
}

func (x *HelloRequest) GetInt32Val() *wrappers.Int32Value {
	if x != nil {
		return x.Int32Val
	}
	return nil
}

func (x *HelloRequest) GetUint32Val() *wrappers.UInt32Value {
	if x != nil {
		return x.Uint32Val
	}
	return nil
}

func (x *HelloRequest) GetInt64Val() *wrappers.Int64Value {
	if x != nil {
		return x.Int64Val
	}
	return nil
}

func (x *HelloRequest) GetUint64Val() *wrappers.UInt64Value {
	if x != nil {
		return x.Uint64Val
	}
	return nil
}

type HelloReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *HelloReply) Reset() {
	*x = HelloReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examples_internal_helloworld_helloworld_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloReply) ProtoMessage() {}

func (x *HelloReply) ProtoReflect() protoreflect.Message {
	mi := &file_examples_internal_helloworld_helloworld_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloReply.ProtoReflect.Descriptor instead.
func (*HelloReply) Descriptor() ([]byte, []int) {
	return file_examples_internal_helloworld_helloworld_proto_rawDescGZIP(), []int{1}
}

func (x *HelloReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_examples_internal_helloworld_helloworld_proto protoreflect.FileDescriptor

var file_examples_internal_helloworld_helloworld_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2f, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x29, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e,
	0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa6, 0x04, 0x0a, 0x0c, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a,
	0x06, 0x73, 0x74, 0x72, 0x56, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x74, 0x72,
	0x56, 0x61, 0x6c, 0x12, 0x37, 0x0a, 0x08, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x08, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x12, 0x3a, 0x0a, 0x09,
	0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x64,
	0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x12, 0x34, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6c,
	0x56, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x12, 0x37,
	0x0a, 0x08, 0x62, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x12, 0x37, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x33, 0x32,
	0x56, 0x61, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33,
	0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c,
	0x12, 0x3a, 0x0a, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x12, 0x37, 0x0a, 0x08,
	0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x56, 0x61, 0x6c, 0x12, 0x3a, 0x0a, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56,
	0x61, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36,
	0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61,
	0x6c, 0x22, 0x26, 0x0a, 0x0a, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x99, 0x03, 0x0a, 0x07, 0x47, 0x72,
	0x65, 0x65, 0x74, 0x65, 0x72, 0x12, 0x8d, 0x03, 0x0a, 0x08, 0x53, 0x61, 0x79, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x12, 0x37, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x68, 0x65, 0x6c,
	0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x90, 0x02, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x89, 0x02, 0x12, 0x0b, 0x2f, 0x73,
	0x61, 0x79, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x5a, 0x16, 0x12, 0x14, 0x2f, 0x73, 0x61,
	0x79, 0x2f, 0x73, 0x74, 0x72, 0x76, 0x61, 0x6c, 0x2f, 0x7b, 0x73, 0x74, 0x72, 0x56, 0x61, 0x6c,
	0x7d, 0x5a, 0x1a, 0x12, 0x18, 0x2f, 0x73, 0x61, 0x79, 0x2f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x76,
	0x61, 0x6c, 0x2f, 0x7b, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x7d, 0x5a, 0x1c, 0x12,
	0x1a, 0x2f, 0x73, 0x61, 0x79, 0x2f, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x76, 0x61, 0x6c, 0x2f,
	0x7b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x7d, 0x5a, 0x18, 0x12, 0x16, 0x2f,
	0x73, 0x61, 0x79, 0x2f, 0x62, 0x6f, 0x6f, 0x6c, 0x76, 0x61, 0x6c, 0x2f, 0x7b, 0x62, 0x6f, 0x6f,
	0x6c, 0x56, 0x61, 0x6c, 0x7d, 0x5a, 0x1a, 0x12, 0x18, 0x2f, 0x73, 0x61, 0x79, 0x2f, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x76, 0x61, 0x6c, 0x2f, 0x7b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c,
	0x7d, 0x5a, 0x1a, 0x12, 0x18, 0x2f, 0x73, 0x61, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x76,
	0x61, 0x6c, 0x2f, 0x7b, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x7d, 0x5a, 0x1c, 0x12,
	0x1a, 0x2f, 0x73, 0x61, 0x79, 0x2f, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x76, 0x61, 0x6c, 0x2f,
	0x7b, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x7d, 0x5a, 0x1a, 0x12, 0x18, 0x2f,
	0x73, 0x61, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x76, 0x61, 0x6c, 0x2f, 0x7b, 0x69, 0x6e,
	0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x7d, 0x5a, 0x1c, 0x12, 0x1a, 0x2f, 0x73, 0x61, 0x79, 0x2f,
	0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x76, 0x61, 0x6c, 0x2f, 0x7b, 0x75, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x56, 0x61, 0x6c, 0x7d, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x76, 0x32, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examples_internal_helloworld_helloworld_proto_rawDescOnce sync.Once
	file_examples_internal_helloworld_helloworld_proto_rawDescData = file_examples_internal_helloworld_helloworld_proto_rawDesc
)

func file_examples_internal_helloworld_helloworld_proto_rawDescGZIP() []byte {
	file_examples_internal_helloworld_helloworld_proto_rawDescOnce.Do(func() {
		file_examples_internal_helloworld_helloworld_proto_rawDescData = protoimpl.X.CompressGZIP(file_examples_internal_helloworld_helloworld_proto_rawDescData)
	})
	return file_examples_internal_helloworld_helloworld_proto_rawDescData
}

var file_examples_internal_helloworld_helloworld_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_examples_internal_helloworld_helloworld_proto_goTypes = []interface{}{
	(*HelloRequest)(nil),         // 0: grpc.gateway.examples.internal.helloworld.HelloRequest
	(*HelloReply)(nil),           // 1: grpc.gateway.examples.internal.helloworld.HelloReply
	(*wrappers.StringValue)(nil), // 2: google.protobuf.StringValue
	(*wrappers.FloatValue)(nil),  // 3: google.protobuf.FloatValue
	(*wrappers.DoubleValue)(nil), // 4: google.protobuf.DoubleValue
	(*wrappers.BoolValue)(nil),   // 5: google.protobuf.BoolValue
	(*wrappers.BytesValue)(nil),  // 6: google.protobuf.BytesValue
	(*wrappers.Int32Value)(nil),  // 7: google.protobuf.Int32Value
	(*wrappers.UInt32Value)(nil), // 8: google.protobuf.UInt32Value
	(*wrappers.Int64Value)(nil),  // 9: google.protobuf.Int64Value
	(*wrappers.UInt64Value)(nil), // 10: google.protobuf.UInt64Value
}
var file_examples_internal_helloworld_helloworld_proto_depIdxs = []int32{
	2,  // 0: grpc.gateway.examples.internal.helloworld.HelloRequest.strVal:type_name -> google.protobuf.StringValue
	3,  // 1: grpc.gateway.examples.internal.helloworld.HelloRequest.floatVal:type_name -> google.protobuf.FloatValue
	4,  // 2: grpc.gateway.examples.internal.helloworld.HelloRequest.doubleVal:type_name -> google.protobuf.DoubleValue
	5,  // 3: grpc.gateway.examples.internal.helloworld.HelloRequest.boolVal:type_name -> google.protobuf.BoolValue
	6,  // 4: grpc.gateway.examples.internal.helloworld.HelloRequest.bytesVal:type_name -> google.protobuf.BytesValue
	7,  // 5: grpc.gateway.examples.internal.helloworld.HelloRequest.int32Val:type_name -> google.protobuf.Int32Value
	8,  // 6: grpc.gateway.examples.internal.helloworld.HelloRequest.uint32Val:type_name -> google.protobuf.UInt32Value
	9,  // 7: grpc.gateway.examples.internal.helloworld.HelloRequest.int64Val:type_name -> google.protobuf.Int64Value
	10, // 8: grpc.gateway.examples.internal.helloworld.HelloRequest.uint64Val:type_name -> google.protobuf.UInt64Value
	0,  // 9: grpc.gateway.examples.internal.helloworld.Greeter.SayHello:input_type -> grpc.gateway.examples.internal.helloworld.HelloRequest
	1,  // 10: grpc.gateway.examples.internal.helloworld.Greeter.SayHello:output_type -> grpc.gateway.examples.internal.helloworld.HelloReply
	10, // [10:11] is the sub-list for method output_type
	9,  // [9:10] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_examples_internal_helloworld_helloworld_proto_init() }
func file_examples_internal_helloworld_helloworld_proto_init() {
	if File_examples_internal_helloworld_helloworld_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examples_internal_helloworld_helloworld_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloRequest); i {
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
		file_examples_internal_helloworld_helloworld_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloReply); i {
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
			RawDescriptor: file_examples_internal_helloworld_helloworld_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examples_internal_helloworld_helloworld_proto_goTypes,
		DependencyIndexes: file_examples_internal_helloworld_helloworld_proto_depIdxs,
		MessageInfos:      file_examples_internal_helloworld_helloworld_proto_msgTypes,
	}.Build()
	File_examples_internal_helloworld_helloworld_proto = out.File
	file_examples_internal_helloworld_helloworld_proto_rawDesc = nil
	file_examples_internal_helloworld_helloworld_proto_goTypes = nil
	file_examples_internal_helloworld_helloworld_proto_depIdxs = nil
}
