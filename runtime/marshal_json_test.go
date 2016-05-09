package runtime_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/gengo/grpc-gateway/examples/examplepb"
	"github.com/gengo/grpc-gateway/runtime"
)

func TestJSONBuiltinMarshal(t *testing.T) {
	var m runtime.JSONBuiltin
	msg := examplepb.SimpleMessage{
		Id: "foo",
	}

	buf, err := m.Marshal(&msg)
	if err != nil {
		t.Errorf("m.Marshal(%v) failed with %v; want success", &msg, err)
	}

	var got examplepb.SimpleMessage
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Errorf("json.Unmarshal(%q, &got) failed with %v; want success", buf, err)
	}
	if want := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want %v", &got, &want)
	}
}

func TestJSONBuiltinMarshalPrimitive(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, v := range []interface{}{
		"",
		"foo",
		1,
		0,
		-1,
		-0.0,
		1.5,
	} {
		buf, err := m.Marshal(v)
		if err != nil {
			t.Errorf("m.Marshal(%v) failed with %v; want success", v, err)
		}

		dest := reflect.New(reflect.TypeOf(v))
		if err := json.Unmarshal(buf, dest.Interface()); err != nil {
			t.Errorf("json.Unmarshal(%q, unmarshaled) failed with %v; want success", buf, err)
		}
		if got, want := dest.Elem().Interface(), v; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want %v", &got, &want)
		}
	}
}

func TestJSONBuiltinsnmarshal(t *testing.T) {
	var (
		m   runtime.JSONBuiltin
		got examplepb.SimpleMessage

		data = []byte(`{"id": "foo"}`)
	)
	if err := m.Unmarshal(data, &got); err != nil {
		t.Errorf("m.Unmarshal(%q, &got) failed with %v; want success", data, err)
	}

	want := examplepb.SimpleMessage{
		Id: "foo",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want = %v", &got, &want)
	}
}

func TestJSONBuiltinUnmarshalPrimitive(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, spec := range []struct {
		data string
		want interface{}
	}{
		{data: `""`, want: ""},
		{data: `"foo"`, want: "foo"},
		{data: `1`, want: 1},
		{data: `0`, want: 0},
		{data: `-1`, want: -1},
		{data: `-0.0`, want: -0.0},
		{data: `1.5`, want: 1.5},
	} {
		dest := reflect.New(reflect.TypeOf(spec.want))
		if err := m.Unmarshal([]byte(spec.data), dest.Interface()); err != nil {
			t.Errorf("m.Unmarshal(%q, dest) failed with %v; want success", spec.data, err)
		}

		if got, want := dest.Elem().Interface(), spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want = %v", got, want)
		}
	}
}

func TestJSONBuiltinEncoder(t *testing.T) {
	var m runtime.JSONBuiltin
	msg := examplepb.SimpleMessage{
		Id: "foo",
	}

	var buf bytes.Buffer
	enc := m.NewEncoder(&buf)
	if err := enc.Encode(&msg); err != nil {
		t.Errorf("enc.Encode(%v) failed with %v; want success", &msg, err)
	}

	var got examplepb.SimpleMessage
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Errorf("json.Unmarshal(%q, &got) failed with %v; want success", buf.String(), err)
	}
	if want := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want %v", &got, &want)
	}
}

func TestJSONBuiltinEncoderPrimitive(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, v := range []interface{}{
		"",
		"foo",
		1,
		0,
		-1,
		-0.0,
		1.5,
	} {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		if err := enc.Encode(v); err != nil {
			t.Errorf("enc.Encode(%v) failed with %v; want success", v, err)
		}

		dest := reflect.New(reflect.TypeOf(v))
		if err := json.Unmarshal(buf.Bytes(), dest.Interface()); err != nil {
			t.Errorf("json.Unmarshal(%q, unmarshaled) failed with %v; want success", buf.String(), err)
		}
		if got, want := dest.Elem().Interface(), v; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want %v", &got, &want)
		}
	}
}

func TestJSONBuiltinDecoder(t *testing.T) {
	var (
		m   runtime.JSONBuiltin
		got examplepb.SimpleMessage

		data = `{"id": "foo"}`
	)
	r := strings.NewReader(data)
	dec := m.NewDecoder(r)
	if err := dec.Decode(&got); err != nil {
		t.Errorf("m.Unmarshal(&got) failed with %v; want success", err)
	}

	want := examplepb.SimpleMessage{
		Id: "foo",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want = %v", &got, &want)
	}
}

func TestJSONBuiltinDecoderPrimitive(t *testing.T) {
	var m runtime.JSONBuiltin
	for _, spec := range []struct {
		data string
		want interface{}
	}{
		{data: `""`, want: ""},
		{data: `"foo"`, want: "foo"},
		{data: `1`, want: 1},
		{data: `0`, want: 0},
		{data: `-1`, want: -1},
		{data: `-0.0`, want: -0.0},
		{data: `1.5`, want: 1.5},
	} {
		r := strings.NewReader(spec.data)
		dec := m.NewDecoder(r)
		dest := reflect.New(reflect.TypeOf(spec.want))
		if err := dec.Decode(dest.Interface()); err != nil {
			t.Errorf("dec.Decode(dest) failed with %v; want success; data=%q", err, spec.data)
		}

		if got, want := dest.Elem().Interface(), spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v; want = %v", got, want)
		}
	}
}
