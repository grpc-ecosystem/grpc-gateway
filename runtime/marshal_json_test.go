package runtime_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestJSONBuiltinMarshal(t *testing.T) {
	var m runtime.JSONBuiltin
	msg := &examplepb.SimpleMessage{
		Id: "foo",
	}

	buf, err := m.Marshal(msg)
	if err != nil {
		t.Errorf("m.Marshal(%v) failed with %v; want success", msg, err)
	}

	got := new(examplepb.SimpleMessage)
	if err := json.Unmarshal(buf, got); err != nil {
		t.Errorf("json.Unmarshal(%q, got) failed with %v; want success", buf, err)
	}
	if diff := cmp.Diff(got, msg, protocmp.Transform()); diff != "" {
		t.Error(diff)
	}
}

func TestJSONBuiltinMarshalField(t *testing.T) {
	var (
		m   runtime.JSONBuiltin
		buf []byte
		err error
	)

	for _, fixt := range builtinFieldFixtures {
		if len(fixt.indent) == 0 {
			buf, err = m.Marshal(fixt.data)
			if err != nil {
				t.Errorf("m.Marshal(%v) failed with %v; want success", fixt.data, err)
			}
		} else {
			buf, err = m.MarshalIndent(fixt.data, "", fixt.indent)
			if err != nil {
				t.Errorf("m.MarshalIndent(%v, \"\", \"%s\") failed with %v; want success", fixt.data, fixt.indent, err)
			}
		}

		if got, want := string(buf), fixt.json; got != want {
			t.Errorf("got = %q; want %q; data = %#v", got, want, fixt.data)
		}
	}
}

func TestJSONBuiltinMarshalFieldKnownErrors(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, fixt := range builtinKnownErrors {
		buf, err := m.Marshal(fixt.data)
		if err != nil {
			t.Errorf("m.Marshal(%v) failed with %v; want success", fixt.data, err)
		}
		if got, want := string(buf), fixt.json; got == want {
			t.Errorf("surprisingly got = %q; as want %q; data = %#v", got, want, fixt.data)
		}
	}
}

func TestJSONBuiltinsnmarshal(t *testing.T) {
	var (
		m   runtime.JSONBuiltin
		got = new(examplepb.SimpleMessage)

		data = []byte(`{"id": "foo"}`)
	)
	if err := m.Unmarshal(data, got); err != nil {
		t.Errorf("m.Unmarshal(%q, got) failed with %v; want success", data, err)
	}

	want := &examplepb.SimpleMessage{
		Id: "foo",
	}
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Error(diff)
	}
}

func TestJSONBuiltinUnmarshalField(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, fixt := range builtinFieldFixtures {
		dest := alloc(reflect.TypeOf(fixt.data))
		if err := m.Unmarshal([]byte(fixt.json), dest.Interface()); err != nil {
			t.Errorf("m.Unmarshal(%q, dest) failed with %v; want success", fixt.json, err)
		}

		got, want := dest.Elem().Interface(), fixt.data
		if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
			t.Error(diff)
		}
	}
}

func alloc(t reflect.Type) reflect.Value {
	if t == nil {
		return reflect.ValueOf(new(interface{}))
	}
	return reflect.New(t)
}

func TestJSONBuiltinUnmarshalFieldKnownErrors(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, fixt := range builtinKnownErrors {
		dest := reflect.New(reflect.TypeOf(fixt.data))
		if err := m.Unmarshal([]byte(fixt.json), dest.Interface()); err == nil {
			t.Errorf("m.Unmarshal(%q, dest) succeeded; want an error", fixt.json)
		}
	}
}

func TestJSONBuiltinEncoder(t *testing.T) {
	var m runtime.JSONBuiltin
	msg := &examplepb.SimpleMessage{
		Id: "foo",
	}

	var buf bytes.Buffer
	enc := m.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		t.Errorf("enc.Encode(%v) failed with %v; want success", msg, err)
	}

	got := new(examplepb.SimpleMessage)
	if err := json.Unmarshal(buf.Bytes(), got); err != nil {
		t.Errorf("json.Unmarshal(%q, got) failed with %v; want success", buf.String(), err)
	}
	if diff := cmp.Diff(got, msg, protocmp.Transform()); diff != "" {
		t.Error(diff)
	}
}

func TestJSONBuiltinEncoderFields(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, fixt := range builtinFieldFixtures {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)

		if fixt.indent != "" {
			if e, ok := enc.(*json.Encoder); ok {
				e.SetIndent("", fixt.indent)
			} else {
				// By default, JSONBuiltin.NewEncoder returns *json.Encoder as runtime.Encoder.
				// Otherwise it's better to fail the tests than skip fixtures with non empty indent
				t.Errorf("enc is not *json.Encoder, unable to set indentation settings. " +
					"This failure prevents testing the correctness of indentation in JSON output.")
			}
		}

		if err := enc.Encode(fixt.data); err != nil {
			t.Errorf("enc.Encode(%#v) failed with %v; want success", fixt.data, err)
		}

		if got, want := buf.String(), fixt.json+"\n"; got != want {
			t.Errorf("got = %q; want %q; data = %#v", got, want, fixt.data)
		}
	}
}

func TestJSONBuiltinDecoder(t *testing.T) {
	var (
		m   runtime.JSONBuiltin
		got = new(examplepb.SimpleMessage)

		data = `{"id": "foo"}`
	)
	r := strings.NewReader(data)
	dec := m.NewDecoder(r)
	if err := dec.Decode(got); err != nil {
		t.Errorf("m.Unmarshal(got) failed with %v; want success", err)
	}

	want := &examplepb.SimpleMessage{
		Id: "foo",
	}
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("got = %v; want = %v", got, want)
	}
}

func TestJSONBuiltinDecoderFields(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, fixt := range builtinFieldFixtures {
		r := strings.NewReader(fixt.json)
		dec := m.NewDecoder(r)
		dest := alloc(reflect.TypeOf(fixt.data))
		if err := dec.Decode(dest.Interface()); err != nil {
			t.Errorf("dec.Decode(dest) failed with %v; want success; data = %q", err, fixt.json)
		}

		got, want := dest.Elem().Interface(), fixt.data
		if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
			t.Error(diff)
		}
	}
}

var (
	defaultIndent        = "  "
	builtinFieldFixtures = []struct {
		data   interface{}
		indent string
		json   string
	}{
		{data: "", json: `""`},
		{data: "", indent: defaultIndent, json: `""`},
		{data: proto.String(""), json: `""`},
		{data: proto.String(""), indent: defaultIndent, json: `""`},
		{data: "foo", json: `"foo"`},
		{data: "foo", indent: defaultIndent, json: `"foo"`},
		{data: []byte("foo"), json: `"Zm9v"`},
		{data: []byte("foo"), indent: defaultIndent, json: `"Zm9v"`},
		{data: []byte{}, json: `""`},
		{data: []byte{}, indent: defaultIndent, json: `""`},
		{data: proto.String("foo"), json: `"foo"`},
		{data: proto.String("foo"), indent: defaultIndent, json: `"foo"`},
		{data: int32(-1), json: "-1"},
		{data: int32(-1), indent: defaultIndent, json: "-1"},
		{data: proto.Int32(-1), json: "-1"},
		{data: proto.Int32(-1), indent: defaultIndent, json: "-1"},
		{data: int64(-1), json: "-1"},
		{data: int64(-1), indent: defaultIndent, json: "-1"},
		{data: proto.Int64(-1), json: "-1"},
		{data: proto.Int64(-1), indent: defaultIndent, json: "-1"},
		{data: uint32(123), json: "123"},
		{data: uint32(123), indent: defaultIndent, json: "123"},
		{data: proto.Uint32(123), json: "123"},
		{data: proto.Uint32(123), indent: defaultIndent, json: "123"},
		{data: uint64(123), json: "123"},
		{data: uint64(123), indent: defaultIndent, json: "123"},
		{data: proto.Uint64(123), json: "123"},
		{data: proto.Uint64(123), indent: defaultIndent, json: "123"},
		{data: float32(-1.5), json: "-1.5"},
		{data: float32(-1.5), indent: defaultIndent, json: "-1.5"},
		{data: proto.Float32(-1.5), json: "-1.5"},
		{data: proto.Float32(-1.5), indent: defaultIndent, json: "-1.5"},
		{data: float64(-1.5), json: "-1.5"},
		{data: float64(-1.5), indent: defaultIndent, json: "-1.5"},
		{data: proto.Float64(-1.5), json: "-1.5"},
		{data: proto.Float64(-1.5), indent: defaultIndent, json: "-1.5"},
		{data: true, json: "true"},
		{data: true, indent: defaultIndent, json: "true"},
		{data: proto.Bool(true), json: "true"},
		{data: proto.Bool(true), indent: defaultIndent, json: "true"},
		{data: (*string)(nil), json: "null"},
		{data: (*string)(nil), indent: defaultIndent, json: "null"},
		{data: new(emptypb.Empty), json: "{}"},
		{data: new(emptypb.Empty), indent: defaultIndent, json: "{}"},
		{data: examplepb.NumericEnum_ONE, json: "1"},
		{data: examplepb.NumericEnum_ONE, indent: defaultIndent, json: "1"},
		{data: nil, json: "null"},
		{data: nil, indent: defaultIndent, json: "null"},
		{data: (*string)(nil), json: "null"},
		{data: (*string)(nil), indent: defaultIndent, json: "null"},
		{data: []interface{}{nil, "foo", -1.0, 1.234, true}, json: `[null,"foo",-1,1.234,true]`},
		{data: []interface{}{nil, "foo", -1.0, 1.234, true}, indent: defaultIndent, json: "[\n  null,\n  \"foo\",\n  -1,\n  1.234,\n  true\n]"},
		{
			data: map[string]interface{}{"bar": nil, "baz": -1.0, "fiz": 1.234, "foo": true},
			json: `{"bar":null,"baz":-1,"fiz":1.234,"foo":true}`,
		},
		{
			data:   map[string]interface{}{"bar": nil, "baz": -1.0, "fiz": 1.234, "foo": true},
			indent: defaultIndent,
			json:   "{\n  \"bar\": null,\n  \"baz\": -1,\n  \"fiz\": 1.234,\n  \"foo\": true\n}",
		},
		{
			data: (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
			json: "1",
		},
		{
			data:   (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
			indent: defaultIndent,
			json:   "1",
		},
		{data: map[string]int{"FOO": 0, "BAR": -1}, json: "{\"BAR\":-1,\"FOO\":0}"},
		{data: map[string]int{"FOO": 0, "BAR": -1}, indent: defaultIndent, json: "{\n  \"BAR\": -1,\n  \"FOO\": 0\n}"},
		{data: struct {
			A string
			B int
			C map[string]int
		}{A: "Go", B: 3, C: map[string]int{"FOO": 0, "BAR": -1}},
			json: "{\"A\":\"Go\",\"B\":3,\"C\":{\"BAR\":-1,\"FOO\":0}}"},
		{data: struct {
			A string
			B int
			C map[string]int
		}{A: "Go", B: 3, C: map[string]int{"FOO": 0, "BAR": -1}}, indent: defaultIndent,
			json: "{\n  \"A\": \"Go\",\n  \"B\": 3,\n  \"C\": {\n    \"BAR\": -1,\n    \"FOO\": 0\n  }\n}"},
	}
	builtinKnownErrors = []struct {
		data interface{}
		json string
	}{
		{data: examplepb.NumericEnum_ONE, json: "ONE"},
		{
			data: (*examplepb.NumericEnum)(proto.Int32(int32(examplepb.NumericEnum_ONE))),
			json: "ONE",
		},
		{
			data: &examplepb.ABitOfEverything_OneofString{OneofString: "abc"},
			json: `"abc"`,
		},
		{
			data: &timestamppb.Timestamp{
				Seconds: 1462875553,
				Nanos:   123000000,
			},
			json: `"2016-05-10T10:19:13.123Z"`,
		},
		{
			data: wrapperspb.Int32(123),
			json: "123",
		},
	}
)
