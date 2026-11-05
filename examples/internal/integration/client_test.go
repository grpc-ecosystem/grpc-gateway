package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	want := &echomodels.ExamplepbSimpleMessage{ID: "foo", Num: "0"}
	if diff := cmp.Diff(want, resp.Payload); diff != "" {
		t.Errorf("resp.Payload differs (-want +got):\n%s", diff)
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
	want := &echomodels.ExamplepbSimpleMessage{ID: "foo", Num: "0"}
	if diff := cmp.Diff(want, resp.Payload); diff != "" {
		t.Errorf("resp.Payload differs (-want +got):\n%s", diff)
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

func TestABitOfEverythingClientCreateBody(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	enumOne := abemodels.ExamplepbNumericEnumONE
	enumZero := abemodels.ExamplepbNumericEnumZERO
	enumPath := abemodels.PathenumPathEnumDEF
	messagePath := abemodels.MessagePathEnumNestedPathEnumJKL
	deepFalse := abemodels.NestedDeepEnumFALSE
	floatVal := float32(1.5)
	doubleVal := 2.5
	int64Val := "4294967296"
	timestamp, err := strfmt.ParseDateTime("2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatalf("strfmt.ParseDateTime failed with %v; want success", err)
	}
	want := &abemodels.ExamplepbABitOfEverything{
		FloatValue:               &floatVal,
		DoubleValue:              &doubleVal,
		Int64Value:               &int64Val,
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
		EnumValue:                &enumOne,
		PathEnumValue:            &enumPath,
		NestedPathEnumValue:      &messagePath,
		EnumValueAnnotation:      &enumOne,
		Nested: []*abemodels.ABitOfEverythingNested{
			{Name: "bar", Amount: 10, Ok: &deepFalse},
			{Name: "baz", Amount: 20, Ok: &deepFalse},
		},
		RepeatedStringValue: []string{"a", "b", "c"},
		OneofString:         "x",
		MapValue: map[string]*abemodels.ExamplepbNumericEnum{
			"a": &enumOne,
			"b": &enumZero,
		},
		MappedStringValue: map[string]string{"a": "x", "b": "y"},
		MappedNestedValue: map[string]abemodels.ABitOfEverythingNested{
			"a": {Name: "x", Amount: 1, Ok: &deepFalse},
			"b": {Name: "y", Amount: 2, Ok: &deepFalse},
		},
		TimestampValue:                           timestamp,
		Int64OverrideType:                        "0",
		RequiredStringViaFieldBehaviorAnnotation: strPtr("foo"),
		RequiredStringField1:                     strPtr(""),
		RequiredStringField2:                     strPtr(""),
		RequiredFieldBehaviorJSONNameCustom:      strPtr(""),
		RequiredFieldSchemaJSONNameCustom:        strPtr(""),
	}

	cl := abeclient.New(gatewayTransport(), strfmt.Default)
	got := createABEBody(t, cl, want)
	if got.UUID == nil || *got.UUID == "" {
		t.Errorf("got.UUID is empty; want not empty")
	}
	got.UUID = nil
	if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("CreateBody response differs (-want +got):\n%s", diff)
	}
}

func TestABitOfEverythingClientLookup(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := abeclient.New(gatewayTransport(), strfmt.Default)
	created := createABEBody(t, cl, &abemodels.ExamplepbABitOfEverything{
		BoolValue:   true,
		StringValue: "strprefix/foo",
	})

	resp, err := cl.ABitOfEverything.ABitOfEverythingServiceLookup(
		a_bit_of_everything.NewABitOfEverythingServiceLookupParamsWithContext(context.Background()).
			WithUUID(string(*created.UUID)),
		nil,
	)
	if err != nil {
		t.Fatalf("cl.ABitOfEverything.ABitOfEverythingServiceLookup(%q) failed with %v; want success", *created.UUID, err)
	}
	if diff := cmp.Diff(created, resp.Payload); diff != "" {
		t.Errorf("Lookup response differs (-want +got):\n%s", diff)
	}

	testABEClientLookupNotFound(t, cl, "not_exist")
}

func TestABitOfEverythingClientDelete(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	cl := abeclient.New(gatewayTransport(), strfmt.Default)
	created := createABEBody(t, cl, &abemodels.ExamplepbABitOfEverything{
		BoolValue:   true,
		StringValue: "strprefix/foo",
	})

	resp, err := cl.ABitOfEverything.ABitOfEverythingServiceDelete(
		a_bit_of_everything.NewABitOfEverythingServiceDeleteParamsWithContext(context.Background()).
			WithUUID(string(*created.UUID)),
		nil,
	)
	if err != nil {
		t.Fatalf("cl.ABitOfEverything.ABitOfEverythingServiceDelete(%q) failed with %v; want success", *created.UUID, err)
	}
	if diff := cmp.Diff(map[string]any{}, resp.Payload); diff != "" {
		t.Errorf("Delete response differs (-want +got):\n%s", diff)
	}

	testABEClientLookupNotFound(t, cl, string(*created.UUID))
}

func TestABitOfEverythingClientGetRepeatedQuery(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	enumOne := abemodels.ExamplepbNumericEnumONE
	want := &abemodels.ExamplepbABitOfEverythingRepeated{
		PathRepeatedFloatValue:   []float32{1.5, -1.5},
		PathRepeatedDoubleValue:  []float64{2.5, -2.5},
		PathRepeatedInt64Value:   []string{"4294967296", "-4294967296"},
		PathRepeatedUint64Value:  []string{"0", "9223372036854775807"},
		PathRepeatedInt32Value:   []int32{2147483647, -2147483648},
		PathRepeatedFixed64Value: []string{"0", "9223372036854775807"},
		PathRepeatedFixed32Value: []int64{0, 4294967295},
		PathRepeatedBoolValue:    []bool{true, false},
		PathRepeatedStringValue:  []string{"foo", "bar"},
		// The client sends standard base64 and escapes "/" as %2F, which the
		// gateway's default unescaping mode decodes before routing, so stick to
		// bytes whose encoding has no "/". "+" still checks a non URL-safe
		// encoding.
		PathRepeatedBytesValue:    []strfmt.Base64{{0x00}, {0xfb}},
		PathRepeatedUint32Value:   []int64{4294967295, 0},
		PathRepeatedEnumValue:     []*abemodels.ExamplepbNumericEnum{&enumOne, &enumOne},
		PathRepeatedSfixed32Value: []int32{-2147483648, 2147483647},
		PathRepeatedSfixed64Value: []string{"-4294967296", "4294967296"},
		PathRepeatedSint32Value:   []int32{2147483646, -2147483647},
		PathRepeatedSint64Value:   []string{"4611686018427387903", "-4611686018427387904"},
	}

	cl := abeclient.New(gatewayTransport(), strfmt.Default)
	resp, err := cl.ABitOfEverything.ABitOfEverythingServiceGetRepeatedQuery(
		a_bit_of_everything.NewABitOfEverythingServiceGetRepeatedQueryParamsWithContext(context.Background()).
			WithPathRepeatedFloatValue(want.PathRepeatedFloatValue).
			WithPathRepeatedDoubleValue(want.PathRepeatedDoubleValue).
			WithPathRepeatedInt64Value(want.PathRepeatedInt64Value).
			WithPathRepeatedUint64Value(want.PathRepeatedUint64Value).
			WithPathRepeatedInt32Value(want.PathRepeatedInt32Value).
			WithPathRepeatedFixed64Value(want.PathRepeatedFixed64Value).
			WithPathRepeatedFixed32Value(want.PathRepeatedFixed32Value).
			WithPathRepeatedBoolValue(want.PathRepeatedBoolValue).
			WithPathRepeatedStringValue(want.PathRepeatedStringValue).
			WithPathRepeatedBytesValue(want.PathRepeatedBytesValue).
			WithPathRepeatedUint32Value(want.PathRepeatedUint32Value).
			WithPathRepeatedEnumValue([]string{string(enumOne), string(enumOne)}).
			WithPathRepeatedSfixed32Value(want.PathRepeatedSfixed32Value).
			WithPathRepeatedSfixed64Value(want.PathRepeatedSfixed64Value).
			WithPathRepeatedSint32Value(want.PathRepeatedSint32Value).
			WithPathRepeatedSint64Value(want.PathRepeatedSint64Value),
		nil,
	)
	if err != nil {
		t.Fatalf("cl.ABitOfEverything.ABitOfEverythingServiceGetRepeatedQuery failed with %v; want success", err)
	}
	if diff := cmp.Diff(want, resp.Payload); diff != "" {
		t.Errorf("GetRepeatedQuery response differs (-want +got):\n%s", diff)
	}
}

// createABEBody creates msg through the CreateBody endpoint and returns the
// stored message.
func createABEBody(t *testing.T, cl *abeclient.ABitOfEverything, msg *abemodels.ExamplepbABitOfEverything) *abemodels.ExamplepbABitOfEverything {
	t.Helper()
	resp, err := cl.ABitOfEverything.ABitOfEverythingServiceCreateBody(
		a_bit_of_everything.NewABitOfEverythingServiceCreateBodyParamsWithContext(context.Background()).WithBody(msg),
		nil,
	)
	if err != nil {
		t.Fatalf("cl.ABitOfEverything.ABitOfEverythingServiceCreateBody(%#v) failed with %v; want success", msg, err)
	}
	return resp.Payload
}

func testABEClientLookupNotFound(t *testing.T, cl *abeclient.ABitOfEverything, uuid string) {
	t.Helper()
	var code int
	_, err := cl.ABitOfEverything.ABitOfEverythingServiceLookup(
		a_bit_of_everything.NewABitOfEverythingServiceLookupParamsWithContext(context.Background()).WithUUID(uuid),
		nil,
		recordStatusCode(&code),
	)
	if err == nil {
		t.Errorf("cl.ABitOfEverything.ABitOfEverythingServiceLookup(%q) succeeded; want failure", uuid)
	}
	if code != http.StatusNotFound {
		t.Errorf("cl.ABitOfEverything.ABitOfEverythingServiceLookup(%q) returned status %d; want %d", uuid, code, http.StatusNotFound)
	}
}

// recordStatusCode stores the HTTP status code of the response in code. The
// ABE spec declares a string body for 404 responses, so the generated client
// fails to decode the gateway's error body and the status code is what we
// can check.
func recordStatusCode(code *int) a_bit_of_everything.ClientOption {
	return func(op *runtime.ClientOperation) {
		reader := op.Reader
		op.Reader = runtime.ClientResponseReaderFunc(func(resp runtime.ClientResponse, consumer runtime.Consumer) (any, error) {
			*code = resp.Code()
			return reader.ReadResponse(resp, consumer)
		})
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
