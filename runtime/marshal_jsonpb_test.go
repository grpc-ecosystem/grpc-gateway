package runtime_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/gengo/grpc-gateway/examples/examplepb"
	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func TestJSONPbMarshal(t *testing.T) {
	msg := examplepb.ABitOfEverything{
		Uuid: "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
		Nested: []*examplepb.ABitOfEverything_Nested{
			{
				Name:   "foo",
				Amount: 12345,
			},
		},
		Uint64Value: 0xFFFFFFFFFFFFFFFF,
		EnumValue:   examplepb.NumericEnum_ONE,
		OneofValue: &examplepb.ABitOfEverything_OneofString{
			OneofString: "bar",
		},
		MapValue: map[string]examplepb.NumericEnum{
			"a": examplepb.NumericEnum_ONE,
			"b": examplepb.NumericEnum_ZERO,
		},
	}

	for _, spec := range []struct {
		enumsAsInts, emitDefaults bool
		indent                    string
		origName                  bool
		verifier                  func(json string)
	}{
		{
			verifier: func(json string) {
				if strings.ContainsAny(json, " \t\r\n") {
					t.Errorf("strings.ContainsAny(%q, %q) = true; want false", json, " \t\r\n")
				}
				if !strings.Contains(json, "ONE") {
					t.Errorf(`strings.Contains(%q, "ONE") = false; want true`, json)
				}
				if want := "uint64Value"; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
		{
			enumsAsInts: true,
			verifier: func(json string) {
				if strings.Contains(json, "ONE") {
					t.Errorf(`strings.Contains(%q, "ONE") = true; want false`, json)
				}
			},
		},
		{
			emitDefaults: true,
			verifier: func(json string) {
				if want := `"sfixed32Value"`; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
		{
			indent: "\t\t",
			verifier: func(json string) {
				if want := "\t\t\"amount\":"; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
		{
			origName: true,
			verifier: func(json string) {
				if want := "uint64_value"; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
	} {
		m := runtime.JSONPb{
			EnumsAsInts:  spec.enumsAsInts,
			EmitDefaults: spec.emitDefaults,
			Indent:       spec.indent,
			OrigName:     spec.origName,
		}
		buf, err := m.Marshal(&msg)
		if err != nil {
			t.Errorf("m.Marshal(%v) failed with %v; want success; spec=%v", &msg, err, spec)
		}

		var got examplepb.ABitOfEverything
		if err := jsonpb.UnmarshalString(string(buf), &got); err != nil {
			t.Errorf("jsonpb.UnmarshalString(%q, &got) failed with %v; want success; spec=%v", string(buf), err, spec)
		}
		if want := msg; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want %v; spec=%v", &got, &want, spec)
		}
		if spec.verifier != nil {
			spec.verifier(string(buf))
		}
	}
}

func TestJSONPbMarshalNonProto(t *testing.T) {
	var m runtime.JSONPb
	for _, spec := range []struct {
		val  interface{}
		want string
	}{
		{val: int32(1), want: "1"},
		{val: proto.Int32(1), want: "1"},
		{val: int64(1), want: "1"},
		{val: proto.Int64(1), want: "1"},
		{val: uint32(1), want: "1"},
		{val: proto.Uint32(1), want: "1"},
		{val: uint64(1), want: "1"},
		{val: proto.Uint64(1), want: "1"},
		{val: "abc", want: `"abc"`},
		{val: proto.String("abc"), want: `"abc"`},
		{val: float32(1.5), want: "1.5"},
		{val: proto.Float32(1.5), want: "1.5"},
		{val: float64(1.5), want: "1.5"},
		{val: proto.Float64(1.5), want: "1.5"},
		{val: true, want: "true"},
		{val: false, want: "false"},
		{val: (*string)(nil), want: "null"},
		{val: examplepb.NumericEnum_ONE, want: `"ONE"`},
		{
			val:  (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
			want: `"ONE"`,
		},
		{
			val: map[string]int32{
				"foo": 1,
			},
			want: `{"foo":1}`,
		},
		{
			val: map[string]*examplepb.SimpleMessage{
				"foo": {Id: "bar"},
			},
			want: `{"foo":{"id":"bar"}}`,
		},
		{
			val: map[int32]*examplepb.SimpleMessage{
				1: {Id: "foo"},
			},
			want: `{"1":{"id":"foo"}}`,
		},
		{
			val: map[bool]*examplepb.SimpleMessage{
				true: {Id: "foo"},
			},
			want: `{"true":{"id":"foo"}}`,
		},
	} {
		buf, err := m.Marshal(spec.val)
		if err != nil {
			t.Errorf("m.Marshal(%#v) failed with %v; want success", spec.val, err)
		}
		if got, want := string(buf), spec.want; got != want {
			t.Errorf("m.Marshal(%#v) = %q; want %q", spec.val, got, want)
		}
	}

	m.EnumsAsInts = true
	buf, err := m.Marshal(examplepb.NumericEnum_ONE)
	if err != nil {
		t.Errorf("m.Marshal(%#v) failed with %v; want success", examplepb.NumericEnum_ONE, err)
	}
	if got, want := string(buf), "1"; got != want {
		t.Errorf("m.Marshal(%#v) = %q; want %q", examplepb.NumericEnum_ONE, got, want)
	}
}

func TestJSONPbUnmarshal(t *testing.T) {
	var (
		m   runtime.JSONPb
		got examplepb.ABitOfEverything
	)
	for _, data := range []string{
		`{
			"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			"nested": [
				{"name": "foo", "amount": 12345}
			],
			"uint64Value": 18446744073709551615,
			"enumValue": "ONE",
			"oneofString": "bar",
			"mapValue": {
				"a": 1,
				"b": 0
			}
		}`,
		`{
			"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			"nested": [
				{"name": "foo", "amount": 12345}
			],
			"uint64Value": "18446744073709551615",
			"enumValue": "ONE",
			"oneofString": "bar",
			"mapValue": {
				"a": 1,
				"b": 0
			}
		}`,
		`{
			"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			"nested": [
				{"name": "foo", "amount": 12345}
			],
			"uint64Value": 18446744073709551615,
			"enumValue": 1,
			"oneofString": "bar",
			"mapValue": {
				"a": 1,
				"b": 0
			}
		}`,
	} {
		if err := m.Unmarshal([]byte(data), &got); err != nil {
			t.Errorf("m.Unmarshal(%q, &got) failed with %v; want success", data, err)
		}

		want := examplepb.ABitOfEverything{
			Uuid: "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			Nested: []*examplepb.ABitOfEverything_Nested{
				{
					Name:   "foo",
					Amount: 12345,
				},
			},
			Uint64Value: 0xFFFFFFFFFFFFFFFF,
			EnumValue:   examplepb.NumericEnum_ONE,
			OneofValue: &examplepb.ABitOfEverything_OneofString{
				OneofString: "bar",
			},
			MapValue: map[string]examplepb.NumericEnum{
				"a": examplepb.NumericEnum_ONE,
				"b": examplepb.NumericEnum_ZERO,
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want = %v", &got, &want)
		}
	}
}

func TestJSONPbUnmarshalNonProto(t *testing.T) {
	var m runtime.JSONPb
	for _, spec := range []struct {
		input string
		want  interface{}
	}{
		{input: "1", want: int32(1)},
		{input: "1", want: proto.Int32(1)},
		{input: "1", want: int64(1)},
		{input: "1", want: proto.Int64(1)},
		{input: "1", want: uint32(1)},
		{input: "1", want: proto.Uint32(1)},
		{input: "1", want: uint64(1)},
		{input: "1", want: proto.Uint64(1)},
		{input: `"abc"`, want: "abc"},
		{input: `"abc"`, want: proto.String("abc")},
		{input: "1.5", want: float32(1.5)},
		{input: "1.5", want: proto.Float32(1.5)},
		{input: "1.5", want: float64(1.5)},
		{input: "1.5", want: proto.Float64(1.5)},
		{input: "true", want: true},
		{input: "false", want: false},
		{input: "null", want: (*string)(nil)},
		// TODO(yugui) Support symbolic enum
		// {input: `"ONE"`, want: examplepb.NumericEnum_ONE},
		{input: `1`, want: examplepb.NumericEnum_ONE},
		{
			input: `1`,
			want:  (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
		},
		{
			input: `{"foo":{"id":"bar"}}`,
			want: map[string]*examplepb.SimpleMessage{
				"foo": {Id: "bar"},
			},
		},
		{
			input: `{"1":{"id":"foo"}}`,
			want: map[int32]*examplepb.SimpleMessage{
				1: {Id: "foo"},
			},
		},
		{
			input: `{"true":{"id":"foo"}}`,
			want: map[bool]*examplepb.SimpleMessage{
				true: {Id: "foo"},
			},
		},
	} {
		dest := reflect.New(reflect.TypeOf(spec.want))
		if err := m.Unmarshal([]byte(spec.input), dest.Interface()); err != nil {
			t.Errorf("m.Unmarshal(%q, %T) failed with %v; want success", spec.input, dest.Interface(), err)
		}
		if got, want := dest.Elem().Interface(), spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("dest = %#v; want %#v; input = %v", got, want, spec.input)
		}
	}
}

func TestJSONPbEncoder(t *testing.T) {
	msg := examplepb.ABitOfEverything{
		Uuid: "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
		Nested: []*examplepb.ABitOfEverything_Nested{
			{
				Name:   "foo",
				Amount: 12345,
			},
		},
		Uint64Value: 0xFFFFFFFFFFFFFFFF,
		OneofValue: &examplepb.ABitOfEverything_OneofString{
			OneofString: "bar",
		},
		MapValue: map[string]examplepb.NumericEnum{
			"a": examplepb.NumericEnum_ONE,
			"b": examplepb.NumericEnum_ZERO,
		},
	}

	for _, spec := range []struct {
		enumsAsInts, emitDefaults bool
		indent                    string
		origName                  bool
		verifier                  func(json string)
	}{
		{
			verifier: func(json string) {
				if strings.ContainsAny(json, " \t\r\n") {
					t.Errorf("strings.ContainsAny(%q, %q) = true; want false", json, " \t\r\n")
				}
				if strings.Contains(json, "ONE") {
					t.Errorf(`strings.Contains(%q, "ONE") = true; want false`, json)
				}
				if want := "uint64Value"; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
		{
			enumsAsInts: true,
			verifier: func(json string) {
				if strings.Contains(json, "ONE") {
					t.Errorf(`strings.Contains(%q, "ONE") = true; want false`, json)
				}
			},
		},
		{
			emitDefaults: true,
			verifier: func(json string) {
				if want := `"sfixed32Value"`; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
		{
			indent: "\t\t",
			verifier: func(json string) {
				if want := "\t\t\"amount\":"; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
		{
			origName: true,
			verifier: func(json string) {
				if want := "uint64_value"; !strings.Contains(json, want) {
					t.Errorf(`strings.Contains(%q, %q) = false; want true`, json, want)
				}
			},
		},
	} {
		m := runtime.JSONPb{
			EnumsAsInts:  spec.enumsAsInts,
			EmitDefaults: spec.emitDefaults,
			Indent:       spec.indent,
			OrigName:     spec.origName,
		}

		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		if err := enc.Encode(&msg); err != nil {
			t.Errorf("enc.Encode(%v) failed with %v; want success; spec=%v", &msg, err, spec)
		}

		var got examplepb.ABitOfEverything
		if err := jsonpb.UnmarshalString(buf.String(), &got); err != nil {
			t.Errorf("jsonpb.UnmarshalString(%q, &got) failed with %v; want success; spec=%v", buf.String(), err, spec)
		}
		if want := msg; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want %v; spec=%v", &got, &want, spec)
		}
		if spec.verifier != nil {
			spec.verifier(buf.String())
		}
	}
}

func TestJSONPbEncoderNonProto(t *testing.T) {
	var m runtime.JSONPb
	for _, spec := range []struct {
		val  interface{}
		want string
	}{
		{val: int32(1), want: "1"},
		{val: proto.Int32(1), want: "1"},
		{val: int64(1), want: "1"},
		{val: proto.Int64(1), want: "1"},
		{val: uint32(1), want: "1"},
		{val: proto.Uint32(1), want: "1"},
		{val: uint64(1), want: "1"},
		{val: proto.Uint64(1), want: "1"},
		{val: "abc", want: `"abc"`},
		{val: proto.String("abc"), want: `"abc"`},
		{val: float32(1.5), want: "1.5"},
		{val: proto.Float32(1.5), want: "1.5"},
		{val: float64(1.5), want: "1.5"},
		{val: proto.Float64(1.5), want: "1.5"},
		{val: true, want: "true"},
		{val: false, want: "false"},
		{val: (*string)(nil), want: "null"},
		{val: examplepb.NumericEnum_ONE, want: `"ONE"`},
		{
			val:  (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
			want: `"ONE"`,
		},
		{
			val: map[string]int32{
				"foo": 1,
			},
			want: `{"foo":1}`,
		},
		{
			val: map[string]*examplepb.SimpleMessage{
				"foo": {Id: "bar"},
			},
			want: `{"foo":{"id":"bar"}}`,
		},
		{
			val: map[int32]*examplepb.SimpleMessage{
				1: {Id: "foo"},
			},
			want: `{"1":{"id":"foo"}}`,
		},
		{
			val: map[bool]*examplepb.SimpleMessage{
				true: {Id: "foo"},
			},
			want: `{"true":{"id":"foo"}}`,
		},
	} {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		if err := enc.Encode(spec.val); err != nil {
			t.Errorf("enc.Encode(%#v) failed with %v; want success", spec.val, err)
		}
		if got, want := buf.String(), spec.want; got != want {
			t.Errorf("m.Marshal(%#v) = %q; want %q", spec.val, got, want)
		}
	}

	m.EnumsAsInts = true
	buf, err := m.Marshal(examplepb.NumericEnum_ONE)
	if err != nil {
		t.Errorf("m.Marshal(%#v) failed with %v; want success", examplepb.NumericEnum_ONE, err)
	}
	if got, want := string(buf), "1"; got != want {
		t.Errorf("m.Marshal(%#v) = %q; want %q", examplepb.NumericEnum_ONE, got, want)
	}
}

func TestJSONPbDecoder(t *testing.T) {
	var (
		m   runtime.JSONPb
		got examplepb.ABitOfEverything
	)
	for _, data := range []string{
		`{
			"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			"nested": [
				{"name": "foo", "amount": 12345}
			],
			"uint64Value": 18446744073709551615,
			"enumValue": "ONE",
			"oneofString": "bar",
			"mapValue": {
				"a": 1,
				"b": 0
			}
		}`,
		`{
			"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			"nested": [
				{"name": "foo", "amount": 12345}
			],
			"uint64Value": "18446744073709551615",
			"enumValue": "ONE",
			"oneofString": "bar",
			"mapValue": {
				"a": 1,
				"b": 0
			}
		}`,
		`{
			"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			"nested": [
				{"name": "foo", "amount": 12345}
			],
			"uint64Value": 18446744073709551615,
			"enumValue": 1,
			"oneofString": "bar",
			"mapValue": {
				"a": 1,
				"b": 0
			}
		}`,
	} {
		r := strings.NewReader(data)
		dec := m.NewDecoder(r)
		if err := dec.Decode(&got); err != nil {
			t.Errorf("m.Unmarshal(&got) failed with %v; want success; data=%q", err, data)
		}

		want := examplepb.ABitOfEverything{
			Uuid: "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
			Nested: []*examplepb.ABitOfEverything_Nested{
				{
					Name:   "foo",
					Amount: 12345,
				},
			},
			Uint64Value: 0xFFFFFFFFFFFFFFFF,
			EnumValue:   examplepb.NumericEnum_ONE,
			OneofValue: &examplepb.ABitOfEverything_OneofString{
				OneofString: "bar",
			},
			MapValue: map[string]examplepb.NumericEnum{
				"a": examplepb.NumericEnum_ONE,
				"b": examplepb.NumericEnum_ZERO,
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want = %v; data = %v", &got, &want, data)
		}
	}
}

func TestJSONPbDecoderNonProto(t *testing.T) {
	var m runtime.JSONPb
	for _, spec := range []struct {
		input string
		want  interface{}
	}{
		{input: "1", want: int32(1)},
		{input: "1", want: proto.Int32(1)},
		{input: "1", want: int64(1)},
		{input: "1", want: proto.Int64(1)},
		{input: "1", want: uint32(1)},
		{input: "1", want: proto.Uint32(1)},
		{input: "1", want: uint64(1)},
		{input: "1", want: proto.Uint64(1)},
		{input: `"abc"`, want: "abc"},
		{input: `"abc"`, want: proto.String("abc")},
		{input: "1.5", want: float32(1.5)},
		{input: "1.5", want: proto.Float32(1.5)},
		{input: "1.5", want: float64(1.5)},
		{input: "1.5", want: proto.Float64(1.5)},
		{input: "true", want: true},
		{input: "false", want: false},
		{input: "null", want: (*string)(nil)},
		// TODO(yugui) Support symbolic enum
		// {input: `"ONE"`, want: examplepb.NumericEnum_ONE},
		{input: `1`, want: examplepb.NumericEnum_ONE},
		{
			input: `1`,
			want:  (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
		},
		{
			input: `{"foo":{"id":"bar"}}`,
			want: map[string]*examplepb.SimpleMessage{
				"foo": {Id: "bar"},
			},
		},
		{
			input: `{"1":{"id":"foo"}}`,
			want: map[int32]*examplepb.SimpleMessage{
				1: {Id: "foo"},
			},
		},
		{
			input: `{"true":{"id":"foo"}}`,
			want: map[bool]*examplepb.SimpleMessage{
				true: {Id: "foo"},
			},
		},
	} {
		dest := reflect.New(reflect.TypeOf(spec.want))
		dec := m.NewDecoder(strings.NewReader(spec.input))
		if err := dec.Decode(dest.Interface()); err != nil {
			t.Errorf("dec.Decode(%T) failed with %v; want success; input = %q", dest.Interface(), err, spec.input)
		}
		if got, want := dest.Elem().Interface(), spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("dest = %#v; want %#v; input = %v", got, want, spec.input)
		}
	}
}
