package runtime_test

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/grpc-ecosystem/grpc-gateway/examples/proto/examplepb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func TestHTTPBodyMarshal(t *testing.T) {
	msg := examplepb.ABitOfEverything{
		SingleNested:        &examplepb.ABitOfEverything_Nested{},
		RepeatedStringValue: []string{},
		MappedStringValue:   map[string]string{},
		MappedNestedValue:   map[string]*examplepb.ABitOfEverything_Nested{},
		RepeatedEnumValue:   []examplepb.NumericEnum{},
		TimestampValue:      &timestamp.Timestamp{},
		Uuid:                "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
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
		RepeatedEnumAnnotation:   []examplepb.NumericEnum{},
		EnumValueAnnotation:      examplepb.NumericEnum_ONE,
		RepeatedStringAnnotation: []string{},
		RepeatedNestedAnnotation: []*examplepb.ABitOfEverything_Nested{},
		NestedAnnotation:         &examplepb.ABitOfEverything_Nested{},
	}

	for i, spec := range []struct {
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
		m := runtime.HTTPBodyMarshaler{
			defaultMarshaler: &runtime.JSONPb{
				EnumsAsInts:  spec.enumsAsInts,
				EmitDefaults: spec.emitDefaults,
				Indent:       spec.indent,
				OrigName:     spec.origName,
			},
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
			t.Errorf("case %d: got = %v; want %v; spec=%v", i, &got, &want, spec)
		}
		if spec.verifier != nil {
			spec.verifier(string(buf))
		}
	}
}

func TestHTTPBodyMarshalFields(t *testing.T) {
	var d runtime.JSONPb
	d.EnumsAsInts = true // builtin fixtures include an enum, expected to be marshaled as int
	m := runtime.HTTPBodyMarshaler{
		defaultMarshaler: d,
	}
	for _, spec := range builtinFieldFixtures {
		buf, err := m.Marshal(spec.data)
		if err != nil {
			t.Errorf("m.Marshal(%#v) failed with %v; want success", spec.data, err)
		}
		if got, want := string(buf), spec.json; got != want {
			t.Errorf("m.Marshal(%#v) = %q; want %q", spec.data, got, want)
		}
	}

	d.EnumsAsInts = false
	buf, err := m.Marshal(examplepb.NumericEnum_ONE)
	if err != nil {
		t.Errorf("m.Marshal(%#v) failed with %v; want success", examplepb.NumericEnum_ONE, err)
	}
	if got, want := string(buf), `"ONE"`; got != want {
		t.Errorf("m.Marshal(%#v) = %q; want %q", examplepb.NumericEnum_ONE, got, want)
	}
}

func TestHTTPBodyUnmarshal(t *testing.T) {
	var (
		got examplepb.ABitOfEverything
	)
	m := runtime.HTTPBodyMarshaler{
		defaultMarshaler: &runtime.JSONPb{},
	}
	for i, data := range []string{
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
			t.Errorf("case %d: m.Unmarshal(%q, &got) failed with %v; want success", i, data, err)
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
			t.Errorf("case %d: got = %v; want = %v", i, &got, &want)
		}
	}
}

func TestHTTPBodyUnmarshalFields(t *testing.T) {
	m := runtime.HTTPBodyMarshaler{
		defaultMarshaler: &runtime.JSONPb{},
	}
	for _, fixt := range fieldFixtures {
		if fixt.skipUnmarshal {
			continue
		}

		dest := reflect.New(reflect.TypeOf(fixt.data))
		if err := m.Unmarshal([]byte(fixt.json), dest.Interface()); err != nil {
			t.Errorf("m.Unmarshal(%q, %T) failed with %v; want success", fixt.json, dest.Interface(), err)
		}
		if got, want := dest.Elem().Interface(), fixt.data; !reflect.DeepEqual(got, want) {
			t.Errorf("dest = %#v; want %#v; input = %v", got, want, fixt.json)
		}
	}
}

func TestHTTPBodyEncoder(t *testing.T) {
	msg := examplepb.ABitOfEverything{
		SingleNested:        &examplepb.ABitOfEverything_Nested{},
		RepeatedStringValue: []string{},
		MappedStringValue:   map[string]string{},
		MappedNestedValue:   map[string]*examplepb.ABitOfEverything_Nested{},
		RepeatedEnumValue:   []examplepb.NumericEnum{},
		TimestampValue:      &timestamp.Timestamp{},
		Uuid:                "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
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
		RepeatedEnumAnnotation:   []examplepb.NumericEnum{},
		EnumValueAnnotation:      examplepb.NumericEnum_ONE,
		RepeatedStringAnnotation: []string{},
		RepeatedNestedAnnotation: []*examplepb.ABitOfEverything_Nested{},
		NestedAnnotation:         &examplepb.ABitOfEverything_Nested{},
	}

	for i, spec := range []struct {
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
		m := runtime.HTTPBodyMarshaler{
			defaultMarshaler: &runtime.JSONPb{
				EnumsAsInts:  spec.enumsAsInts,
				EmitDefaults: spec.emitDefaults,
				Indent:       spec.indent,
				OrigName:     spec.origName,
			},
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
			t.Errorf("case %d: got = %v; want %v; spec=%v", i, &got, &want, spec)
		}
		if spec.verifier != nil {
			spec.verifier(buf.String())
		}
	}
}

func TestHTTPBodyEncoderFields(t *testing.T) {
	var d runtime.JSONPb
	d.EnumsAsInts = true // builtin fixtures include an enum, expected to be marshaled as int
	m := runtime.HTTPBodyMarshaler{
		defaultMarshaler: d,
	}
	for _, fixt := range fieldFixtures {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		if err := enc.Encode(fixt.data); err != nil {
			t.Errorf("enc.Encode(%#v) failed with %v; want success", fixt.data, err)
		}
		if got, want := buf.String(), fixt.json; got != want {
			t.Errorf("enc.Encode(%#v) = %q; want %q", fixt.data, got, want)
		}
	}

	d.EnumsAsInts = true
	buf, err := m.Marshal(examplepb.NumericEnum_ONE)
	if err != nil {
		t.Errorf("m.Marshal(%#v) failed with %v; want success", examplepb.NumericEnum_ONE, err)
	}
	if got, want := string(buf), "1"; got != want {
		t.Errorf("m.Marshal(%#v) = %q; want %q", examplepb.NumericEnum_ONE, got, want)
	}
}

func TestHTTPBodyDecoder(t *testing.T) {
	var (
		got examplepb.ABitOfEverything
	)
	m := runtime.HTTPBodyMarshaler{
		defaultMarshaler: &runtime.JSONPb{},
	}
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

func TestHTTPBodyDecoderFields(t *testing.T) {
	m := runtime.HTTPBodyMarshaler{
		defaultMarshaler: &runtime.JSONPb{},
	}
	for _, fixt := range fieldFixtures {
		if fixt.skipUnmarshal {
			continue
		}

		dest := reflect.New(reflect.TypeOf(fixt.data))
		dec := m.NewDecoder(strings.NewReader(fixt.json))
		if err := dec.Decode(dest.Interface()); err != nil {
			t.Errorf("dec.Decode(%T) failed with %v; want success; input = %q", dest.Interface(), err, fixt.json)
		}
		if got, want := dest.Elem().Interface(), fixt.data; !reflect.DeepEqual(got, want) {
			t.Errorf("dest = %#v; want %#v; input = %v", got, want, fixt.json)
		}
	}
}

func TestHTTPBodyMarshalResponseBodies(t *testing.T) {
	for i, spec := range []struct {
		input        interface{}
		emitDefaults bool
		verifier     func(json string)
	}{
		{
			input: &examplepb.ResponseBodyOut{
				Response: &examplepb.ResponseBodyOut_Response{Data: "abcdef"},
			},
			verifier: func(json string) {
				expected := `{"response":{"data":"abcdef"}}`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			emitDefaults: true,
			input:        &examplepb.ResponseBodyOut{},
			verifier: func(json string) {
				expected := `{"response":null}`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: &examplepb.RepeatedResponseBodyOut_Response{},
			verifier: func(json string) {
				expected := `{}`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			emitDefaults: true,
			input:        &examplepb.RepeatedResponseBodyOut_Response{},
			verifier: func(json string) {
				expected := `{"data":"","type":"UNKNOWN"}`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: ([]*examplepb.RepeatedResponseBodyOut_Response)(nil),
			verifier: func(json string) {
				expected := `null`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			emitDefaults: true,
			input:        ([]*examplepb.RepeatedResponseBodyOut_Response)(nil),
			verifier: func(json string) {
				expected := `[]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: []*examplepb.RepeatedResponseBodyOut_Response{},
			verifier: func(json string) {
				expected := `[]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: []string{"something"},
			verifier: func(json string) {
				expected := `["something"]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: []string{},
			verifier: func(json string) {
				expected := `[]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: ([]string)(nil),
			verifier: func(json string) {
				expected := `null`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			emitDefaults: true,
			input:        ([]string)(nil),
			verifier: func(json string) {
				expected := `[]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			input: []*examplepb.RepeatedResponseBodyOut_Response{
				&examplepb.RepeatedResponseBodyOut_Response{},
				&examplepb.RepeatedResponseBodyOut_Response{
					Data: "abc",
					Type: examplepb.RepeatedResponseBodyOut_Response_A,
				},
			},
			verifier: func(json string) {
				expected := `[{},{"data":"abc","type":"A"}]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
		{
			emitDefaults: true,
			input: []*examplepb.RepeatedResponseBodyOut_Response{
				&examplepb.RepeatedResponseBodyOut_Response{},
				&examplepb.RepeatedResponseBodyOut_Response{
					Data: "abc",
					Type: examplepb.RepeatedResponseBodyOut_Response_B,
				},
			},
			verifier: func(json string) {
				expected := `[{"data":"","type":"UNKNOWN"},{"data":"abc","type":"B"}]`
				if json != expected {
					t.Errorf("json not equal (%q, %q)", json, expected)
				}
			},
		},
	} {

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := runtime.HTTPBodyMarshaler{
				defaultMarshaler: &runtime.JSONPb{
					EmitDefaults: spec.emitDefaults,
				},
			}
			val := spec.input
			buf, err := m.Marshal(val)
			if err != nil {
				t.Errorf("m.Marshal(%v) failed with %v; want success; spec=%v", val, err, spec)
			}
			if spec.verifier != nil {
				spec.verifier(string(buf))
			}
		})
	}
}
