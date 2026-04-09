package genopenapi

import (
	"slices"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// fieldNames returns field names in order. Used to compare splitOneofs output
// without dragging full descriptor equality into the assertion.
func fieldNames(fields []*descriptor.Field) []string {
	out := make([]string, len(fields))
	for i, f := range fields {
		out[i] = f.GetName()
	}
	return out
}

// mkField builds a plain regular field.
func mkField(name string) *descriptor.Field {
	return &descriptor.Field{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
		Name: proto.String(name),
	}}
}

// mkOneofField builds a field that belongs to a real (non-synthetic) oneof group.
func mkOneofField(name string, oneofIdx int32) *descriptor.Field {
	return &descriptor.Field{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
		Name:       proto.String(name),
		OneofIndex: proto.Int32(oneofIdx),
	}}
}

// mkProto3OptionalField builds a proto3-optional field. These have an
// OneofIndex pointing at a synthetic single-field oneof that splitOneofs
// must treat as a regular optional field.
func mkProto3OptionalField(name string, oneofIdx int32) *descriptor.Field {
	return &descriptor.Field{FieldDescriptorProto: &descriptorpb.FieldDescriptorProto{
		Name:           proto.String(name),
		OneofIndex:     proto.Int32(oneofIdx),
		Proto3Optional: proto.Bool(true),
	}}
}

// mkMessage builds a message with the given fields and oneof decl names in
// the order given. The first oneof decl name is index 0, etc.
func mkMessage(fields []*descriptor.Field, oneofDecls ...string) *descriptor.Message {
	decls := make([]*descriptorpb.OneofDescriptorProto, len(oneofDecls))
	for i, n := range oneofDecls {
		decls[i] = &descriptorpb.OneofDescriptorProto{Name: proto.String(n)}
	}
	return &descriptor.Message{
		DescriptorProto: &descriptorpb.DescriptorProto{OneofDecl: decls},
		Fields:          fields,
	}
}

func TestSplitOneofs(t *testing.T) {
	t.Parallel()

	type wantGroup struct {
		name   string
		fields []string
	}

	cases := []struct {
		name        string
		msg         *descriptor.Message
		wantRegular []string
		wantGroups  []wantGroup
	}{
		{
			name: "empty message",
			msg:  mkMessage(nil),
		},
		{
			name: "only regular fields",
			msg: mkMessage([]*descriptor.Field{
				mkField("id"),
				mkField("title"),
			}),
			wantRegular: []string{"id", "title"},
		},
		{
			name: "single real oneof group",
			msg: mkMessage([]*descriptor.Field{
				mkField("id"),
				mkOneofField("paperback_isbn", 0),
				mkOneofField("ebook_url", 0),
			}, "format"),
			wantRegular: []string{"id"},
			wantGroups: []wantGroup{
				{name: "format", fields: []string{"paperback_isbn", "ebook_url"}},
			},
		},
		{
			// proto3-optional fields each sit in their own synthetic single-
			// field oneof. splitOneofs must pull them back into `regular`
			// and NOT emit a group for their synthetic oneofs.
			name: "proto3-optional treated as regular",
			msg: mkMessage([]*descriptor.Field{
				mkField("id"),
				mkProto3OptionalField("nickname", 0),
			}, "_nickname"),
			wantRegular: []string{"id", "nickname"},
		},
		{
			// The returned groups must follow OneofDecl order, not the
			// order fields happen to appear in. Here decl order is
			// {format, provenance} while the fields interleave.
			name: "groups ordered by declaration",
			msg: mkMessage([]*descriptor.Field{
				mkOneofField("bought_at", 1),
				mkOneofField("paperback_isbn", 0),
				mkOneofField("borrowed_from", 1),
				mkOneofField("ebook_url", 0),
			}, "format", "provenance"),
			wantGroups: []wantGroup{
				{name: "format", fields: []string{"paperback_isbn", "ebook_url"}},
				{name: "provenance", fields: []string{"bought_at", "borrowed_from"}},
			},
		},
		{
			// A oneof that's declared but has no fields must not appear in
			// the output groups — the loop over OneofDecl checks byIndex
			// before emitting.
			name: "empty oneof declaration is skipped",
			msg: mkMessage([]*descriptor.Field{
				mkField("id"),
				mkOneofField("paperback_isbn", 0),
			}, "format", "provenance"),
			wantRegular: []string{"id"},
			wantGroups: []wantGroup{
				{name: "format", fields: []string{"paperback_isbn"}},
			},
		},
		{
			// Mix of all three kinds: a regular field, a real oneof group,
			// and a proto3-optional field sharing the same message. The
			// proto3-optional field must land in `regular`; only the real
			// group must appear in `groups`.
			name: "mixed regular, real oneof, proto3-optional",
			msg: mkMessage([]*descriptor.Field{
				mkField("id"),
				mkOneofField("paperback_isbn", 0),
				mkProto3OptionalField("nickname", 1),
				mkOneofField("ebook_url", 0),
			}, "format", "_nickname"),
			wantRegular: []string{"id", "nickname"},
			wantGroups: []wantGroup{
				{name: "format", fields: []string{"paperback_isbn", "ebook_url"}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			regular, groups := splitOneofs(tc.msg)

			if got := fieldNames(regular); !slices.Equal(got, tc.wantRegular) {
				t.Errorf("regular fields = %v, want %v", got, tc.wantRegular)
			}
			if len(groups) != len(tc.wantGroups) {
				t.Fatalf("group count = %d, want %d: groups = %+v", len(groups), len(tc.wantGroups), groups)
			}
			for i, g := range groups {
				if g.name != tc.wantGroups[i].name {
					t.Errorf("group[%d] name = %q, want %q", i, g.name, tc.wantGroups[i].name)
				}
				if got := fieldNames(g.fields); !slices.Equal(got, tc.wantGroups[i].fields) {
					t.Errorf("group[%d] fields = %v, want %v", i, got, tc.wantGroups[i].fields)
				}
			}
		})
	}
}

