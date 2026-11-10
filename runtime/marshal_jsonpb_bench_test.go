package runtime_test

import (
	"bytes"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// benchMarshalCases covers the two shapes that matter for transcoding overhead: a small
// response, where any fixed per-call cost dominates, and a message with many fields, where
// protojson's own reflection dominates.
func benchMarshalCases() []struct {
	name string
	msg  proto.Message
} {
	return []struct {
		name string
		msg  proto.Message
	}{
		{
			name: "small",
			msg:  &examplepb.SimpleMessage{Id: "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7"},
		},
		{
			name: "many_fields",
			msg: &examplepb.ABitOfEverything{
				Uuid:                "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
				SingleNested:        &examplepb.ABitOfEverything_Nested{Name: "bar", Amount: 10},
				Nested:              []*examplepb.ABitOfEverything_Nested{{Name: "bar", Amount: 10}, {Name: "baz", Amount: 20}},
				FloatValue:          1.5,
				DoubleValue:         2.5,
				Int64Value:          4294967296,
				Uint64Value:         9223372036854775807,
				Int32Value:          -2147483648,
				BoolValue:           true,
				StringValue:         "strprefix/foo",
				Uint32Value:         4294967295,
				EnumValue:           examplepb.NumericEnum_ONE,
				RepeatedStringValue: []string{"a", "b", "c"},
				MappedStringValue:   map[string]string{"one": "1", "two": "2"},
			},
		},
	}
}

// BenchmarkJSONPbMarshal measures JSONPb.Marshal against the protojson call it wraps.
//
// The gap between the two arms is the marshaler's own overhead on top of transcoding: every
// unary response goes through JSONPb.Marshal (see ForwardResponseMessage), so it is paid once
// per request.
func BenchmarkJSONPbMarshal(b *testing.B) {
	for _, c := range benchMarshalCases() {
		msg := c.msg

		b.Run(c.name+"/JSONPb", func(b *testing.B) {
			m := &runtime.JSONPb{}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := m.Marshal(msg); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(c.name+"/protojson", func(b *testing.B) {
			var opts protojson.MarshalOptions
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := opts.Marshal(msg); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkJSONPbDecode measures the decoder returned by JSONPb.NewDecoder against a direct
// protojson.Unmarshal of the same bytes.
//
// Both arms allocate a fresh destination message, so the gap is what the decoder adds: the
// *json.Decoder, plus a scan of the whole value by encoding/json and a copy of it into a
// json.RawMessage before protojson parses those bytes a second time.
func BenchmarkJSONPbDecode(b *testing.B) {
	for _, c := range benchMarshalCases() {
		msg := c.msg
		data, err := protojson.Marshal(msg)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(c.name+"/NewDecoder", func(b *testing.B) {
			m := &runtime.JSONPb{}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				dst := msg.ProtoReflect().New().Interface()
				if err := m.NewDecoder(bytes.NewReader(data)).Decode(dst); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(c.name+"/protojson", func(b *testing.B) {
			var opts protojson.UnmarshalOptions
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				dst := msg.ProtoReflect().New().Interface()
				if err := opts.Unmarshal(data, dst); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
