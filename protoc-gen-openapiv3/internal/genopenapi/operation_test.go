package genopenapi

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func strField(name string, num int32) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:   proto.String(name),
		Number: proto.Int32(num),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
	}
}

// TestFieldConsumedByPathParams covers the #2624 case: a body field whose only
// sub-field is bound to a nested path parameter (id.value) must be reported as
// consumed, while a sibling scalar body field must not.
func TestFieldConsumedByPathParams(t *testing.T) {
	fdp := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Package: proto.String("test"),
		Syntax:  proto.String("proto3"),
		Options: &descriptorpb.FileOptions{GoPackage: proto.String("github.com/example/test")},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name:  proto.String("EntityId"),
				Field: []*descriptorpb.FieldDescriptorProto{strField("value", 1)},
			},
			{
				Name: proto.String("UpdateEntityRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("id"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: proto.String(".test.EntityId"),
					},
					strField("other_field", 2),
				},
			},
		},
	}

	reg := descriptor.NewRegistry()
	if err := reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile:      []*descriptorpb.FileDescriptorProto{fdp},
		FileToGenerate: []string{"test.proto"},
	}); err != nil {
		t.Fatalf("registry load: %v", err)
	}

	req, err := reg.LookupMsg("", ".test.UpdateEntityRequest")
	if err != nil {
		t.Fatalf("lookup request message: %v", err)
	}
	entityID, err := reg.LookupMsg("", ".test.EntityId")
	if err != nil {
		t.Fatalf("lookup EntityId message: %v", err)
	}

	var idField, otherField *descriptor.Field
	for _, f := range req.Fields {
		switch f.GetName() {
		case "id":
			idField = f
		case "other_field":
			otherField = f
		}
	}
	if idField == nil || otherField == nil {
		t.Fatalf("did not resolve request fields: id=%v other=%v", idField, otherField)
	}
	valueField := entityID.Fields[0]

	// Path parameter id.value: FieldPath [id, value], leaf Target = value.
	params := []descriptor.Parameter{{
		FieldPath: descriptor.FieldPath{
			{Name: "id", Target: idField},
			{Name: "value", Target: valueField},
		},
		Target: valueField,
	}}

	b := &schemaBuilder{reg: reg}
	if !b.fieldConsumedByPathParams(idField, params) {
		t.Error("id's only sub-field is the id.value path parameter, so it should be consumed")
	}
	if b.fieldConsumedByPathParams(otherField, params) {
		t.Error("other_field is a real body field and should not be consumed")
	}
}
