package integration_test

import (
	"context"
	"fmt"
	"testing"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	abeclient "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/abe/client"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/abe/client/a_bit_of_everything"
	abemodels "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/abe/models"
	echoclient "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/echo/client"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/echo/client/echo_service"
	echomodels "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/echo/models"
	uaeclient "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/unannotatedecho/client"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/unannotatedecho/client/unannotated_echo_service"
	uaemodels "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/unannotatedecho/models"
	"github.com/rogpeppe/fastuuid"
)

var uuidgen = fastuuid.MustNewGenerator()

// gatewayTransport returns a runtime transport pointed at the integration
// gateway started by main_test.go. Each call yields a fresh transport so
// per-test customization (timeouts, hooks) doesn't leak between tests.
func gatewayTransport() *httptransport.Runtime {
	return httptransport.New("localhost:8088", "/", []string{"http"})
}

// uuid renders a fastuuid value as an RFC 4122-shaped string. The grpc-
// gateway server doesn't enforce the format, but the spec declares
// `format: uuid` and go-openapi's strfmt registry rejects values that
// don't match — matching the canonical layout keeps the params struct
// from refusing to encode the request.
func uuid() strfmt.UUID {
	b := uuidgen.Next()
	return strfmt.UUID(fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]))
}

func TestEchoClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := echoclient.New(gatewayTransport(), strfmt.Default)
	resp, err := cl.EchoService.EchoServiceEcho(
		echo_service.NewEchoServiceEchoParamsWithContext(context.Background()).WithID("foo"),
	)
	if err != nil {
		t.Errorf(`cl.EchoService.EchoServiceEcho("foo") failed with %v; want success`, err)
		return
	}
	if got, want := resp.Payload.ID, "foo"; got != want {
		t.Errorf("resp.Payload.ID = %q; want %q", got, want)
	}
}

func TestEchoBodyClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := echoclient.New(gatewayTransport(), strfmt.Default)
	req := &echomodels.ExamplepbSimpleMessage{ID: "foo"}
	resp, err := cl.EchoService.EchoServiceEchoBody(
		echo_service.NewEchoServiceEchoBodyParamsWithContext(context.Background()).WithBody(req),
	)
	if err != nil {
		t.Errorf("cl.EchoService.EchoServiceEchoBody(%#v) failed with %v; want success", req, err)
		return
	}
	if got, want := resp.Payload.ID, "foo"; got != want {
		t.Errorf("resp.Payload.ID = %q; want %q", got, want)
	}
}

func TestEchoBody2Client(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := echoclient.New(gatewayTransport(), strfmt.Default)
	req := &echomodels.ExamplepbEmbedded{Note: "note"}
	resp, err := cl.EchoService.EchoServiceEchoBody2(
		echo_service.NewEchoServiceEchoBody2ParamsWithContext(context.Background()).
			WithID("foo").
			WithNo(req),
	)
	if err != nil {
		t.Errorf("cl.EchoService.EchoServiceEchoBody2(%#v) failed with %v; want success", req, err)
		return
	}
	if got, want := resp.Payload.ID, "foo"; got != want {
		t.Errorf("resp.Payload.ID = %q; want %q", got, want)
	}
}

func TestAbitOfEverythingClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := abeclient.New(gatewayTransport(), strfmt.Default)
	testABEClientCreate(t, cl)
}

func testABEClientCreate(t *testing.T, cl *abeclient.ABitOfEverything) {
	enumZero := abemodels.ExamplepbNumericEnumZERO
	enumPath := abemodels.PathenumPathEnumABC
	messagePath := abemodels.MessagePathEnumNestedPathEnumJKL

	floatVal := float32(1.5)
	doubleVal := 2.5
	int64Val := "4294967296"
	requestUUID := uuid()
	want := &abemodels.ExamplepbABitOfEverything{
		FloatValue:                               &floatVal,
		DoubleValue:                              &doubleVal,
		Int64Value:                               &int64Val,
		Uint64Value:                              "9223372036854775807",
		Int32Value:                               -2147483648,
		Fixed64Value:                             "9223372036854775807",
		Fixed32Value:                             4294967295,
		BoolValue:                                true,
		StringValue:                              "strprefix/foo",
		Uint32Value:                              4294967295,
		Sfixed32Value:                            2147483647,
		Sfixed64Value:                            "-4611686018427387904",
		Sint32Value:                              2147483647,
		Sint64Value:                              "4611686018427387903",
		NonConventionalNameValue:                 "camelCase",
		EnumValue:                                &enumZero,
		PathEnumValue:                            &enumPath,
		NestedPathEnumValue:                      &messagePath,
		EnumValueAnnotation:                      &enumZero,
		UUID:                                     &requestUUID,
		RequiredFieldBehaviorJSONNameCustom:      strPtr("test"),
		RequiredFieldSchemaJSONNameCustom:        strPtr("test"),
		RequiredStringField1:                     strPtr(""),
		RequiredStringField2:                     strPtr(""),
		RequiredStringViaFieldBehaviorAnnotation: strPtr(""),
	}

	createResp, err := cl.ABitOfEverything.ABitOfEverythingServiceCreate(
		a_bit_of_everything.NewABitOfEverythingServiceCreateParamsWithContext(context.Background()).
			WithFloatValue(*want.FloatValue).
			WithDoubleValue(*want.DoubleValue).
			WithInt64Value(*want.Int64Value).
			WithUint64Value(want.Uint64Value).
			WithInt32Value(want.Int32Value).
			WithFixed64Value(want.Fixed64Value).
			WithFixed32Value(want.Fixed32Value).
			WithBoolValue(want.BoolValue).
			WithStringValue(want.StringValue).
			WithUint32Value(want.Uint32Value).
			WithSfixed32Value(want.Sfixed32Value).
			WithSfixed64Value(want.Sfixed64Value).
			WithSint32Value(want.Sint32Value).
			WithSint64Value(want.Sint64Value).
			WithNonConventionalNameValue(want.NonConventionalNameValue).
			WithEnumValue(string(*want.EnumValue)).
			WithPathEnumValue(string(*want.PathEnumValue)).
			WithNestedPathEnumValue(string(*want.NestedPathEnumValue)).
			WithEnumValueAnnotation(string(*want.EnumValueAnnotation)).
			WithUUID(*want.UUID).
			WithRequiredStringViaFieldBehaviorAnnotation(*want.RequiredStringViaFieldBehaviorAnnotation).
			WithRequiredStringField1(*want.RequiredStringField1).
			WithRequiredStringField2(*want.RequiredStringField2).
			WithRequiredFieldBehaviorJSONNameCustom(*want.RequiredFieldBehaviorJSONNameCustom).
			WithRequiredFieldSchemaJSONNameCustom(*want.RequiredFieldSchemaJSONNameCustom),
		nil,
	)
	if err != nil {
		t.Fatalf("cl.ABitOfEverythingService.ABitOfEverythingServiceCreate(%#v) failed with %v; want success", want, err)
	}
	resp := createResp.Payload
	if resp.UUID == nil || *resp.UUID == "" {
		t.Errorf("resp.UUID is empty; want not empty")
	}
	resp.UUID = nil
	want.UUID = nil

	if resp.FloatValue == nil || *resp.FloatValue != *want.FloatValue {
		t.Error("float")
	}
	if resp.DoubleValue == nil || *resp.DoubleValue != *want.DoubleValue {
		t.Error("double")
	}
	if resp.Int64Value == nil || *resp.Int64Value != *want.Int64Value {
		t.Error("int64")
	}
	if resp.Uint64Value != want.Uint64Value {
		t.Error("uint64")
	}
	if resp.Int32Value != want.Int32Value {
		t.Error("int32")
	}
	if resp.Fixed32Value != want.Fixed32Value {
		t.Error("fixed32")
	}
	if resp.Fixed64Value != want.Fixed64Value {
		t.Error("fixed64")
	}
	if resp.BoolValue != want.BoolValue {
		t.Error("bool")
	}
	if resp.StringValue != want.StringValue {
		t.Error("string")
	}
	if resp.Uint32Value != want.Uint32Value {
		t.Error("uint32")
	}
	if resp.Sfixed32Value != want.Sfixed32Value {
		t.Error("sfixed32")
	}
	if resp.Sfixed64Value != want.Sfixed64Value {
		t.Error("sfixed64")
	}
	if resp.Sint32Value != want.Sint32Value {
		t.Error("sint32")
	}
	if resp.Sint64Value != want.Sint64Value {
		t.Error("sint64")
	}
	if resp.NonConventionalNameValue != want.NonConventionalNameValue {
		t.Error("non-conventional name")
	}
	if resp.EnumValue == nil || *resp.EnumValue != *want.EnumValue {
		t.Error("enum")
	}
	if resp.PathEnumValue == nil || *resp.PathEnumValue != *want.PathEnumValue {
		t.Error("path enum")
	}
	if resp.NestedPathEnumValue == nil || *resp.NestedPathEnumValue != *want.NestedPathEnumValue {
		t.Error("nested path enum")
	}
}

func TestUnannotatedEchoClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := uaeclient.New(gatewayTransport(), strfmt.Default)
	resp, err := cl.UnannotatedEchoService.UnannotatedEchoServiceEcho(
		unannotated_echo_service.NewUnannotatedEchoServiceEchoParamsWithContext(context.Background()).
			WithID("foo").
			WithNum("1"),
		nil,
	)
	if err != nil {
		t.Errorf(`cl.UnannotatedEchoService.UnannotatedEchoServiceEcho("foo", "1") failed with %v; want success`, err)
		return
	}
	if resp.Payload.ID == nil {
		t.Errorf("resp.Payload.ID = nil; want %q", "foo")
		return
	}
	if got, want := *resp.Payload.ID, "foo"; got != want {
		t.Errorf("resp.Payload.ID = %q; want %q", got, want)
	}
}

func TestUnannotatedEchoBodyClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := uaeclient.New(gatewayTransport(), strfmt.Default)
	req := &uaemodels.ExamplepbUnannotatedSimpleMessage{ID: strPtr("foo"), Num: strPtr("1")}
	resp, err := cl.UnannotatedEchoService.UnannotatedEchoServiceEchoBody(
		unannotated_echo_service.NewUnannotatedEchoServiceEchoBodyParamsWithContext(context.Background()).WithBody(req),
		nil,
	)
	if err != nil {
		t.Errorf("cl.UnannotatedEchoService.UnannotatedEchoServiceEchoBody(%#v) failed with %v; want success", req, err)
		return
	}
	if resp.Payload.ID == nil {
		t.Errorf("resp.Payload.ID = nil; want %q", "foo")
		return
	}
	if got, want := *resp.Payload.ID, "foo"; got != want {
		t.Errorf("resp.Payload.ID = %q; want %q", got, want)
	}
}

func strPtr(s string) *string { return &s }
