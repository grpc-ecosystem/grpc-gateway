package main

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/go-test/deep"

	"github.com/grpc-ecosystem/grpc-gateway/examples/clients/abe"
	"github.com/grpc-ecosystem/grpc-gateway/examples/clients/echo"
)

func TestClientIntegration(t *testing.T) {
}

func TestEchoClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	config := echo.NewConfiguration()
	config.BasePath = "http://localhost:8080"
	cl := echo.NewAPIClient(config)
	resp, _, err := cl.EchoServiceApi.Echo(context.Background(), "foo")
	if err != nil {
		t.Errorf(`cl.EchoServiceApi.Echo("foo") failed with %v; want success`, err)
	}
	if got, want := resp.Id, "foo"; got != want {
		t.Errorf("resp.Id = %q; want %q", got, want)
	}
}

func TestEchoBodyClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	config := echo.NewConfiguration()
	config.BasePath = "http://localhost:8080"
	cl := echo.NewAPIClient(config)
	req := echo.ExamplepbSimpleMessage{Id: "foo"}
	resp, _, err := cl.EchoServiceApi.EchoBody(context.Background(), req)
	if err != nil {
		t.Errorf("cl.EchoServiceApi.EchoBody(%#v) failed with %v; want success", req, err)
	}
	if got, want := resp.Id, "foo"; got != want {
		t.Errorf("resp.Id = %q; want %q", got, want)
	}
}

func TestAbitOfEverythingClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	config := abe.NewConfiguration()
	config.BasePath = "http://localhost:8080"
	cl := abe.NewAPIClient(config)
	testABEClientCreate(t, cl)
	testABEClientCreateBody(t, cl)
}

func testABEClientCreate(t *testing.T, cl *abe.APIClient) {
	abeZERO := abe.ZERO

	want := &abe.ExamplepbABitOfEverything{
		Nested:                   []abe.ABitOfEverythingNested{},
		FloatValue:               1.5,
		DoubleValue:              2.5,
		Int64Value:               "4294967296",
		Uint64Value:              "9223372036854775807",
		Int32Value:               -2147483648,
		Fixed64Value:             "9223372036854775807",
		Fixed32Value:             4294967295,
		BoolValue:                true,
		StringValue:              "strprefix/foo",
		Uint32Value:              4294967295,
		Sfixed32Value:            2147483647,
		Sfixed64Value:            "-4611686018427387904",
		Sint32Value:              2147483647,
		Sint64Value:              "4611686018427387903",
		NonConventionalNameValue: "camelCase",
		// Fields that aren't explictly sent so they will come back with defaults.
		EnumValue:           &abeZERO,
		RepeatedStringValue: []string{},
		MapValue:            map[string]abe.ExamplepbNumericEnum{},
		MappedStringValue:   map[string]string{},
		MappedNestedValue:   map[string]abe.ABitOfEverythingNested{},
		RepeatedEnumValue:   []abe.ExamplepbNumericEnum{},
	}
	resp, _, err := cl.ABitOfEverythingServiceApi.Create(
		context.Background(),
		want.FloatValue,
		want.DoubleValue,
		want.Int64Value,
		want.Uint64Value,
		want.Int32Value,
		want.Fixed64Value,
		want.Fixed32Value,
		want.BoolValue,
		want.StringValue,
		want.Uint32Value,
		want.Sfixed32Value,
		want.Sfixed64Value,
		want.Sint32Value,
		want.Sint64Value,
		want.NonConventionalNameValue,
	)
	if err != nil {
		t.Errorf("cl.EchoServiceApi.Create(%#v) failed with %v; want success", want, err)
	}
	if resp.Uuid == "" {
		t.Errorf("resp.Uuid is empty; want not empty")
	}
	resp.Uuid = ""
	if diff := deep.Equal(&resp, want); diff != nil {
		t.Errorf("Create: %v", diff)
	}
}

func testABEClientCreateBody(t *testing.T, cl *abe.APIClient) {
	abeOne := abe.ONE
	abeFalse := abe.FALSE
	abeTrue := abe.TRUE
	want := abe.ExamplepbABitOfEverything{
		FloatValue:               1.5,
		DoubleValue:              2.5,
		Int64Value:               "4294967296",
		Uint64Value:              "9223372036854775807",
		Int32Value:               -2147483648,
		Fixed64Value:             "9223372036854775807",
		Fixed32Value:             4294967295,
		BoolValue:                true,
		StringValue:              "strprefix/foo",
		Uint32Value:              4294967295,
		Sfixed32Value:            2147483647,
		Sfixed64Value:            "-4611686018427387904",
		Sint32Value:              2147483647,
		Sint64Value:              "4611686018427387903",
		NonConventionalNameValue: "camelCase",

		Nested: []abe.ABitOfEverythingNested{
			{
				Ok:     &abeFalse,
				Name:   "bar",
				Amount: 10,
			},
			{
				Ok:     &abeTrue,
				Name:   "baz",
				Amount: 20,
			},
		},
		RepeatedStringValue: []string{"a", "b", "c"},
		OneofString:         "x",
		MapValue:            map[string]abe.ExamplepbNumericEnum{
		//TODO: Fix enums
		//"a": &abeOne,
		//"b": &abeZero,
		},
		MappedStringValue: map[string]string{
			"a": "x",
			"b": "y",
		},
		MappedNestedValue: map[string]abe.ABitOfEverythingNested{
			"a": {Ok: &abeFalse, Name: "x", Amount: 1},
			"b": {Ok: &abeTrue, Name: "y", Amount: 2},
		},

		RepeatedEnumValue: []abe.ExamplepbNumericEnum{},
		EnumValue:         &abeOne,
	}
	resp, _, err := cl.ABitOfEverythingServiceApi.CreateBody(context.Background(), want)
	if err != nil {
		t.Errorf("cl.ABitOfEverythingServiceApi.CreateBody(%#v) failed with %v; want success", want, err)
	}
	if resp.Uuid == "" {
		t.Errorf("resp.Uuid is empty; want not empty")
	}
	resp.Uuid = ""
	if diff := deep.Equal(resp, want); diff != nil {
		t.Errorf("CreateBody: %v", diff)
	}
}
