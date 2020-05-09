package integration_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/pathenum"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/sub"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestEcho(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	for _, apiPrefix := range []string{"v1", "v2"} {
		t.Run(apiPrefix, func(t *testing.T) {
			testEcho(t, 8088, apiPrefix, "application/json")
			testEchoOneof(t, 8088, apiPrefix, "application/json")
			testEchoOneof1(t, 8088, apiPrefix, "application/json")
			testEchoOneof2(t, 8088, apiPrefix, "application/json")
			testEchoBody(t, 8088, apiPrefix)
		})
	}
}

func TestForwardResponseOption(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	port := 7079
	go func() {
		if err := runGateway(
			ctx,
			fmt.Sprintf(":%d", port),
			runtime.WithForwardResponseOption(
				func(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
					w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1.1+json")
					return nil
				},
			),
		); err != nil {
			t.Errorf("runGateway() failed with %v; want success", err)
			return
		}
	}()
	if err := waitForGateway(ctx, uint16(port)); err != nil {
		t.Errorf("waitForGateway(ctx, %d) failed with %v; want success", port, err)
	}
	testEcho(t, port, "v1", "application/vnd.docker.plugins.v1.1+json")
}

func testEcho(t *testing.T, port int, apiPrefix string, contentType string) {
	apiURL := fmt.Sprintf("http://localhost:%d/%s/example/echo/myid", port, apiPrefix)
	resp, err := http.Post(apiURL, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.SimpleMessage
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got, want := msg.Id, "myid"; got != want {
		t.Errorf("msg.Id = %q; want %q", got, want)
	}

	if value := resp.Header.Get("Content-Type"); value != contentType {
		t.Errorf("Content-Type was %s, wanted %s", value, contentType)
	}
}

func testEchoOneof(t *testing.T, port int, apiPrefix string, contentType string) {
	apiURL := fmt.Sprintf("http://localhost:%d/%s/example/echo/myid/10/golang", port, apiPrefix)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.SimpleMessage
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got, want := msg.GetLang(), "golang"; got != want {
		t.Errorf("msg.GetLang() = %q; want %q", got, want)
	}

	if value := resp.Header.Get("Content-Type"); value != contentType {
		t.Errorf("Content-Type was %s, wanted %s", value, contentType)
	}
}

func testEchoOneof1(t *testing.T, port int, apiPrefix string, contentType string) {
	apiURL := fmt.Sprintf("http://localhost:%d/%s/example/echo1/myid/10/golang", port, apiPrefix)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.SimpleMessage
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got, want := msg.GetStatus().GetNote(), "golang"; got != want {
		t.Errorf("msg.GetStatus().GetNote() = %q; want %q", got, want)
	}

	if value := resp.Header.Get("Content-Type"); value != contentType {
		t.Errorf("Content-Type was %s, wanted %s", value, contentType)
	}
}

func testEchoOneof2(t *testing.T, port int, apiPrefix string, contentType string) {
	apiURL := fmt.Sprintf("http://localhost:%d/%s/example/echo2/golang", port, apiPrefix)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.SimpleMessage
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got, want := msg.GetNo().GetNote(), "golang"; got != want {
		t.Errorf("msg.GetNo().GetNote() = %q; want %q", got, want)
	}

	if value := resp.Header.Get("Content-Type"); value != contentType {
		t.Errorf("Content-Type was %s, wanted %s", value, contentType)
	}
}

func testEchoBody(t *testing.T, port int, apiPrefix string) {
	sent := examplepb.SimpleMessage{Id: "example"}
	var m jsonpb.Marshaler
	payload, err := m.MarshalToString(&sent)
	if err != nil {
		t.Fatalf("m.MarshalToString(%#v) failed with %v; want success", payload, err)
	}

	apiURL := fmt.Sprintf("http://localhost:%d/%s/example/echo_body", port, apiPrefix)
	resp, err := http.Post(apiURL, "", strings.NewReader(payload))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var received examplepb.SimpleMessage
	if err := jsonpb.UnmarshalString(string(buf), &received); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if diff := cmp.Diff(received, sent, protocmp.Transform()); diff != "" {
		t.Errorf(diff)
	}

	if got, want := resp.Header.Get("Grpc-Metadata-Foo"), "foo1"; got != want {
		t.Errorf("Grpc-Metadata-Foo was %q, wanted %q", got, want)
	}
	if got, want := resp.Header.Get("Grpc-Metadata-Bar"), "bar1"; got != want {
		t.Errorf("Grpc-Metadata-Bar was %q, wanted %q", got, want)
	}

	if got, want := resp.Trailer.Get("Grpc-Trailer-Foo"), "foo2"; got != want {
		t.Errorf("Grpc-Trailer-Foo was %q, wanted %q", got, want)
	}
	if got, want := resp.Trailer.Get("Grpc-Trailer-Bar"), "bar2"; got != want {
		t.Errorf("Grpc-Trailer-Bar was %q, wanted %q", got, want)
	}
}

func TestABE(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	testABECreate(t, 8088)
	testABECreateBody(t, 8088)
	testABEBulkCreate(t, 8088)
	testABEBulkCreateWithError(t, 8088)
	testABELookup(t, 8088)
	testABELookupNotFound(t, 8088)
	testABEList(t, 8088)
	testABEDownload(t, 8088)
	testABEBulkEcho(t, 8088)
	testABEBulkEchoZeroLength(t, 8088)
	testAdditionalBindings(t, 8088)
	testABERepeated(t, 8088)
}

func testABECreate(t *testing.T, port int) {
	want := examplepb.ABitOfEverything{
		FloatValue:               1.5,
		DoubleValue:              2.5,
		Int64Value:               4294967296,
		Uint64Value:              9223372036854775807,
		Int32Value:               -2147483648,
		Fixed64Value:             9223372036854775807,
		Fixed32Value:             4294967295,
		BoolValue:                true,
		StringValue:              "strprefix/foo",
		Uint32Value:              4294967295,
		Sfixed32Value:            2147483647,
		Sfixed64Value:            -4611686018427387904,
		Sint32Value:              2147483647,
		Sint64Value:              4611686018427387903,
		NonConventionalNameValue: "camelCase",
		EnumValue:                examplepb.NumericEnum_ZERO,
		PathEnumValue:            pathenum.PathEnum_DEF,
		NestedPathEnumValue:      pathenum.MessagePathEnum_JKL,
		EnumValueAnnotation:      examplepb.NumericEnum_ONE,
	}
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/%f/%f/%d/separator/%d/%d/%d/%d/%v/%s/%d/%d/%d/%d/%d/%s/%s/%s/%s/%s", port, want.FloatValue, want.DoubleValue, want.Int64Value, want.Uint64Value, want.Int32Value, want.Fixed64Value, want.Fixed32Value, want.BoolValue, want.StringValue, want.Uint32Value, want.Sfixed32Value, want.Sfixed64Value, want.Sint32Value, want.Sint64Value, want.NonConventionalNameValue, want.EnumValue, want.PathEnumValue, want.NestedPathEnumValue, want.EnumValueAnnotation)

	resp, err := http.Post(apiURL, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.ABitOfEverything
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if msg.Uuid == "" {
		t.Error("msg.Uuid is empty; want not empty")
	}
	msg.Uuid = ""
	if diff := cmp.Diff(msg, want, protocmp.Transform()); diff != "" {
		t.Errorf(diff)
	}
}

func testABECreateBody(t *testing.T, port int) {
	want := examplepb.ABitOfEverything{
		FloatValue:               1.5,
		DoubleValue:              2.5,
		Int64Value:               4294967296,
		Uint64Value:              9223372036854775807,
		Int32Value:               -2147483648,
		Fixed64Value:             9223372036854775807,
		Fixed32Value:             4294967295,
		BoolValue:                true,
		StringValue:              "strprefix/foo",
		Uint32Value:              4294967295,
		Sfixed32Value:            2147483647,
		Sfixed64Value:            -4611686018427387904,
		Sint32Value:              2147483647,
		Sint64Value:              4611686018427387903,
		NonConventionalNameValue: "camelCase",
		EnumValue:                examplepb.NumericEnum_ONE,
		PathEnumValue:            pathenum.PathEnum_ABC,
		NestedPathEnumValue:      pathenum.MessagePathEnum_GHI,

		Nested: []*examplepb.ABitOfEverything_Nested{
			{
				Name:   "bar",
				Amount: 10,
			},
			{
				Name:   "baz",
				Amount: 20,
			},
		},
		RepeatedStringValue: []string{"a", "b", "c"},
		OneofValue: &examplepb.ABitOfEverything_OneofString{
			OneofString: "x",
		},
		MapValue: map[string]examplepb.NumericEnum{
			"a": examplepb.NumericEnum_ONE,
			"b": examplepb.NumericEnum_ZERO,
		},
		MappedStringValue: map[string]string{
			"a": "x",
			"b": "y",
		},
		MappedNestedValue: map[string]*examplepb.ABitOfEverything_Nested{
			"a": {Name: "x", Amount: 1},
			"b": {Name: "y", Amount: 2},
		},
		RepeatedEnumAnnotation: []examplepb.NumericEnum{
			examplepb.NumericEnum_ONE,
			examplepb.NumericEnum_ZERO,
		},
		EnumValueAnnotation: examplepb.NumericEnum_ONE,
		RepeatedStringAnnotation: []string{
			"a",
			"b",
		},
		RepeatedNestedAnnotation: []*examplepb.ABitOfEverything_Nested{
			{
				Name:   "hoge",
				Amount: 10,
			},
			{
				Name:   "fuga",
				Amount: 20,
			},
		},
		NestedAnnotation: &examplepb.ABitOfEverything_Nested{
			Name:   "hoge",
			Amount: 10,
		},
	}
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything", port)
	var m jsonpb.Marshaler
	payload, err := m.MarshalToString(&want)
	if err != nil {
		t.Fatalf("m.MarshalToString(%#v) failed with %v; want success", want, err)
	}

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(payload))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.ABitOfEverything
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if msg.Uuid == "" {
		t.Error("msg.Uuid is empty; want not empty")
	}
	msg.Uuid = ""
	if diff := cmp.Diff(msg, want, protocmp.Transform()); diff != "" {
		t.Errorf(diff)
	}
}

func testABEBulkCreate(t *testing.T, port int) {
	count := 0
	r, w := io.Pipe()
	go func(w io.WriteCloser) {
		defer func() {
			if cerr := w.Close(); cerr != nil {
				t.Errorf("w.Close() failed with %v; want success", cerr)
			}
		}()
		for _, val := range []string{
			"foo", "bar", "baz", "qux", "quux",
		} {
			want := examplepb.ABitOfEverything{
				FloatValue:               1.5,
				DoubleValue:              2.5,
				Int64Value:               4294967296,
				Uint64Value:              9223372036854775807,
				Int32Value:               -2147483648,
				Fixed64Value:             9223372036854775807,
				Fixed32Value:             4294967295,
				BoolValue:                true,
				StringValue:              fmt.Sprintf("strprefix/%s", val),
				Uint32Value:              4294967295,
				Sfixed32Value:            2147483647,
				Sfixed64Value:            -4611686018427387904,
				Sint32Value:              2147483647,
				Sint64Value:              4611686018427387903,
				NonConventionalNameValue: "camelCase",
				EnumValue:                examplepb.NumericEnum_ONE,
				PathEnumValue:            pathenum.PathEnum_ABC,
				NestedPathEnumValue:      pathenum.MessagePathEnum_GHI,

				Nested: []*examplepb.ABitOfEverything_Nested{
					{
						Name:   "hoge",
						Amount: 10,
					},
					{
						Name:   "fuga",
						Amount: 20,
					},
				},
				RepeatedEnumAnnotation: []examplepb.NumericEnum{
					examplepb.NumericEnum_ONE,
					examplepb.NumericEnum_ZERO,
				},
				EnumValueAnnotation: examplepb.NumericEnum_ONE,
				RepeatedStringAnnotation: []string{
					"a",
					"b",
				},
				RepeatedNestedAnnotation: []*examplepb.ABitOfEverything_Nested{
					{
						Name:   "hoge",
						Amount: 10,
					},
					{
						Name:   "fuga",
						Amount: 20,
					},
				},
				NestedAnnotation: &examplepb.ABitOfEverything_Nested{
					Name:   "hoge",
					Amount: 10,
				},
			}
			var m jsonpb.Marshaler
			if err := m.Marshal(w, &want); err != nil {
				t.Fatalf("m.Marshal(%#v, w) failed with %v; want success", want, err)
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				t.Errorf("w.Write(%q) failed with %v; want success", "\n", err)
				return
			}
			count++
		}
	}(w)
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/bulk", port)
	resp, err := http.Post(apiURL, "application/json", r)
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg empty.Empty
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}

	if got, want := resp.Header.Get("Grpc-Metadata-Count"), fmt.Sprintf("%d", count); got != want {
		t.Errorf("Grpc-Metadata-Count was %q, wanted %q", got, want)
	}

	if got, want := resp.Trailer.Get("Grpc-Trailer-Foo"), "foo2"; got != want {
		t.Errorf("Grpc-Trailer-Foo was %q, wanted %q", got, want)
	}
	if got, want := resp.Trailer.Get("Grpc-Trailer-Bar"), "bar2"; got != want {
		t.Errorf("Grpc-Trailer-Bar was %q, wanted %q", got, want)
	}
}

func testABEBulkCreateWithError(t *testing.T, port int) {
	count := 0
	r, w := io.Pipe()
	go func(w io.WriteCloser) {
		defer func() {
			if cerr := w.Close(); cerr != nil {
				t.Errorf("w.Close() failed with %v; want success", cerr)
			}
		}()
		for _, val := range []string{
			"foo", "bar", "baz", "qux", "quux",
		} {
			time.Sleep(1 * time.Millisecond)

			want := examplepb.ABitOfEverything{
				StringValue: fmt.Sprintf("strprefix/%s", val),
			}
			var m jsonpb.Marshaler
			if err := m.Marshal(w, &want); err != nil {
				t.Fatalf("m.Marshal(%#v, w) failed with %v; want success", want, err)
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				t.Errorf("w.Write(%q) failed with %v; want success", "\n", err)
				return
			}
			count++
		}
	}(w)

	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/bulk", port)
	request, err := http.NewRequest("POST", apiURL, r)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, nil) failed with %v; want success", "POST", apiURL, err)
	}
	request.Header.Add("Grpc-Metadata-error", "some error")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg spb.Status
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Fatalf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
	}
}

func testABELookup(t *testing.T, port int) {
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything", port)
	cresp, err := http.Post(apiURL, "application/json", strings.NewReader(`
		{"bool_value": true, "string_value": "strprefix/example"}
	`))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer cresp.Body.Close()
	buf, err := ioutil.ReadAll(cresp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(cresp.Body) failed with %v; want success", err)
		return
	}
	if got, want := cresp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
		return
	}

	var want examplepb.ABitOfEverything
	if err := jsonpb.UnmarshalString(string(buf), &want); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &want) failed with %v; want success", buf, err)
		return
	}

	apiURL = fmt.Sprintf("%s/%s", apiURL, want.Uuid)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	var msg examplepb.ABitOfEverything
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if diff := cmp.Diff(msg, want, protocmp.Transform()); diff != "" {
		t.Errorf(diff)
	}

	if got, want := resp.Header.Get("Grpc-Metadata-Uuid"), want.Uuid; got != want {
		t.Errorf("Grpc-Metadata-Uuid was %s, wanted %s", got, want)
	}
}

// TestABEPatch demonstrates partially updating a resource.
// First, we'll create an ABE resource with known values for string_value and int32_value
// Then, issue a PATCH request updating only the string_value
// Then, GET the resource and verify that string_value is changed, but int32_value isn't
func TestABEPatch(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	port := 8088

	// create a record with a known string_value and int32_value
	uuid := postABE(t, port, examplepb.ABitOfEverything{StringValue: "strprefix/bar", Int32Value: 32})

	// issue PATCH request, only updating string_value
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("http://localhost:%d/v2/example/a_bit_of_everything/%s", port, uuid),
		strings.NewReader(`{"string_value": "strprefix/foo"}`),
	)
	if err != nil {
		t.Fatalf("http.NewRequest(PATCH) failed with %v; want success", err)
	}
	patchResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to issue PATCH request: %v", err)
	}
	if got, want := patchResp.StatusCode, http.StatusOK; got != want {
		if body, err := ioutil.ReadAll(patchResp.Body); err != nil {
			t.Errorf("patchResp body couldn't be read: %v", err)
		} else {
			t.Errorf("patchResp.StatusCode= %d; want %d resp: %v", got, want, string(body))
		}
	}

	// issue GET request, verifying that string_value is changed and int32_value is not
	getRespBody := getABE(t, port, uuid)
	if got, want := getRespBody.StringValue, "strprefix/foo"; got != want {
		t.Errorf("string_value= %q; want %q", got, want)
	}
	if got, want := getRespBody.Int32Value, int32(32); got != want {
		t.Errorf("int_32_value= %d; want %d", got, want)
	}
}

// TestABEPatchBody demonstrates the ability to specify an update mask within the request body.
// This binding does not use an automatically generated update_mask.
func TestABEPatchBody(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	port := 8088

	for _, tc := range []struct {
		name          string
		originalValue examplepb.ABitOfEverything
		input         examplepb.UpdateV2Request
		want          examplepb.ABitOfEverything
	}{
		{
			name: "with fieldmask provided",
			originalValue: examplepb.ABitOfEverything{
				Int32Value:   42,
				StringValue:  "rabbit",
				SingleNested: &examplepb.ABitOfEverything_Nested{Name: "some value that will get overwritten", Amount: 345},
			},
			input: examplepb.UpdateV2Request{Abe: &examplepb.ABitOfEverything{
				StringValue:  "some value that won't get updated because it's not in the field mask",
				SingleNested: &examplepb.ABitOfEverything_Nested{Amount: 456},
			}, UpdateMask: &field_mask.FieldMask{Paths: []string{"single_nested"}}},
			want: examplepb.ABitOfEverything{
				Int32Value:   42,
				StringValue:  "rabbit",
				SingleNested: &examplepb.ABitOfEverything_Nested{Amount: 456},
			},
		},
		{
			// N.B. This case passes the empty field mask to the UpdateV2 method so falls back to PUT semantics as per the implementation.
			name: "with empty fieldmask",
			originalValue: examplepb.ABitOfEverything{
				Int32Value:   42,
				StringValue:  "some value that will get overwritten",
				SingleNested: &examplepb.ABitOfEverything_Nested{Name: "value that will get empty", Amount: 345},
			},
			input: examplepb.UpdateV2Request{Abe: &examplepb.ABitOfEverything{
				StringValue:  "some updated value because the fieldMask is nil",
				SingleNested: &examplepb.ABitOfEverything_Nested{Amount: 456},
			}, UpdateMask: &field_mask.FieldMask{}},
			want: examplepb.ABitOfEverything{
				StringValue:  "some updated value because the fieldMask is nil",
				SingleNested: &examplepb.ABitOfEverything_Nested{Amount: 456},
			},
		},
		{
			// N.B. This case passes the nil field mask to the UpdateV2 method so falls back to PUT semantics as per the implementation.
			name: "with nil fieldmask",
			originalValue: examplepb.ABitOfEverything{
				Int32Value:   42,
				StringValue:  "some value that will get overwritten",
				SingleNested: &examplepb.ABitOfEverything_Nested{Name: "value that will get empty", Amount: 123},
			},
			input: examplepb.UpdateV2Request{Abe: &examplepb.ABitOfEverything{
				StringValue:  "some updated value because the fieldMask is nil",
				SingleNested: &examplepb.ABitOfEverything_Nested{Amount: 657},
			}, UpdateMask: nil},
			want: examplepb.ABitOfEverything{
				StringValue:  "some updated value because the fieldMask is nil",
				SingleNested: &examplepb.ABitOfEverything_Nested{Amount: 657},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			originalABE := tc.originalValue
			uuid := postABE(t, port, originalABE)

			patchBody := tc.input
			patchReq, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("http://localhost:%d/v2a/example/a_bit_of_everything/%s", port, uuid),
				strings.NewReader(mustMarshal(t, patchBody)),
			)
			if err != nil {
				t.Fatalf("http.NewRequest(PATCH) failed with %v; want success", err)
			}
			patchResp, err := http.DefaultClient.Do(patchReq)
			if err != nil {
				t.Fatalf("failed to issue PATCH request: %v", err)
			}
			if got, want := patchResp.StatusCode, http.StatusOK; got != want {
				if body, err := ioutil.ReadAll(patchResp.Body); err != nil {
					t.Errorf("patchResp body couldn't be read: %v", err)
				} else {
					t.Errorf("patchResp.StatusCode= %d; want %d resp: %v", got, want, string(body))
				}
			}

			want, got := tc.want, getABE(t, port, uuid)
			got.Uuid = "" // empty out uuid so we don't need to worry about it in comparisons
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

// mustMarshal marshals the given object into a json string, calling t.Fatal if an error occurs. Useful in testing to
// inline marshalling whenever you don't expect the marshalling to return an error
func mustMarshal(t *testing.T, i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		t.Fatalf("failed to marshal %#v: %v", i, err)
	}

	return string(b)
}

// postABE conveniently creates a new ABE record for ease in testing
func postABE(t *testing.T, port int, abe examplepb.ABitOfEverything) (uuid string) {
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything", port)
	postResp, err := http.Post(apiURL, "application/json", strings.NewReader(mustMarshal(t, abe)))
	if err != nil {
		t.Fatalf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		t.Fatalf("postResp body couldn't be read: %v", err)
	}
	var f struct {
		UUID string `json:"uuid"`
	}
	if err := json.Unmarshal(body, &f); err != nil {
		t.Fatalf("postResp body couldn't be unmarshalled: %v. body: %s", err, string(body))
	}
	if f.UUID == "" {
		t.Fatalf("want uuid from postResp, but got none. body: %s", string(body))
	}
	return f.UUID
}

// getABE conveniently fetches an ABE record for ease in testing
func getABE(t *testing.T, port int, uuid string) examplepb.ABitOfEverything {
	gURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/%s", port, uuid)
	getResp, err := http.Get(gURL)
	if err != nil {
		t.Fatalf("http.Get(%s) failed with %v; want success", gURL, err)
	}
	defer getResp.Body.Close()

	if got, want := getResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("getResp.StatusCode= %d, want %d. resp: %v", got, want, getResp)
	}
	var getRespBody examplepb.ABitOfEverything
	body, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		t.Fatalf("getResp body couldn't be read: %v", err)
	}
	if err := json.Unmarshal(body, &getRespBody); err != nil {
		t.Fatalf("getResp body couldn't be unmarshalled: %v body: %s", err, string(body))
	}

	return getRespBody
}

func testABELookupNotFound(t *testing.T, port int) {
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything", port)
	uuid := "not_exist"
	apiURL = fmt.Sprintf("%s/%s", apiURL, uuid)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
		return
	}

	var msg spb.Status
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}

	if got, want := msg.Code, int32(codes.NotFound); got != want {
		t.Errorf("msg.Code = %d; want %d", got, want)
		return
	}

	if got, want := msg.Message, "not found"; got != want {
		t.Errorf("msg.Message = %s; want %s", got, want)
		return
	}

	if got, want := resp.Header.Get("Grpc-Metadata-Uuid"), uuid; got != want {
		t.Errorf("Grpc-Metadata-Uuid was %s, wanted %s", got, want)
	}
	if got, want := resp.Trailer.Get("Grpc-Trailer-Foo"), "foo2"; got != want {
		t.Errorf("Grpc-Trailer-Foo was %q, wanted %q", got, want)
	}
	if got, want := resp.Trailer.Get("Grpc-Trailer-Bar"), "bar2"; got != want {
		t.Errorf("Grpc-Trailer-Bar was %q, wanted %q", got, want)
	}
}

func testABEList(t *testing.T, port int) {
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything", port)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var i int
	for i = 0; ; i++ {
		var item struct {
			Result json.RawMessage        `json:"result"`
			Error  map[string]interface{} `json:"error"`
		}
		err := dec.Decode(&item)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("dec.Decode(&item) failed with %v; want success; i = %d", err, i)
		}
		if len(item.Error) != 0 {
			t.Errorf("item.Error = %#v; want empty; i = %d", item.Error, i)
			continue
		}
		var msg examplepb.ABitOfEverything
		if err := jsonpb.UnmarshalString(string(item.Result), &msg); err != nil {
			t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", item.Result, err)
		}
	}
	if i <= 0 {
		t.Errorf("i == %d; want > 0", i)
	}

	value := resp.Header.Get("Grpc-Metadata-Count")
	if value == "" {
		t.Errorf("Grpc-Metadata-Count should not be empty")
	}

	count, err := strconv.Atoi(value)
	if err != nil {
		t.Errorf("failed to Atoi %q: %v", value, err)
	}

	if count <= 0 {
		t.Errorf("count == %d; want > 0", count)
	}
}

func testABEDownload(t *testing.T, port int) {
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/download", port)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()

	wantHeader := "text/html"
	if value := resp.Header.Get("Content-Type"); value != wantHeader {
		t.Fatalf("testABEDownload() Content-Type failed: got %s, want %s", value, wantHeader)
	}

	body, err := readAll(resp.Body)
	if err != nil {
		t.Fatalf("readAll(resp.Body) failed with %v; want success", err)
	}

	want := []string{"Hello 1", "Hello 2"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("testABEDownload() failed: got %v, want %v", body, want)
	}
}

func testABEBulkEcho(t *testing.T, port int) {
	reqr, reqw := io.Pipe()
	var wg sync.WaitGroup
	var want []*sub.StringMessage
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer reqw.Close()
		var m jsonpb.Marshaler
		for i := 0; i < 1000; i++ {
			msg := sub.StringMessage{Value: proto.String(fmt.Sprintf("message %d", i))}
			buf, err := m.MarshalToString(&msg)
			if err != nil {
				t.Errorf("m.Marshal(%v) failed with %v; want success", &msg, err)
				return
			}
			if _, err := fmt.Fprintln(reqw, buf); err != nil {
				t.Errorf("fmt.Fprintln(reqw, %q) failed with %v; want success", buf, err)
				return
			}
			want = append(want, &msg)
		}
	}()

	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/echo", port)
	req, err := http.NewRequest("POST", apiURL, reqr)
	if err != nil {
		t.Errorf("http.NewRequest(%q, %q, reqr) failed with %v; want success", "POST", apiURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "chunked")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http.Post(%q, %q, req) failed with %v; want success", apiURL, "application/json", err)
		return
	}
	defer resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
	}

	var got []*sub.StringMessage
	wg.Add(1)
	go func() {
		defer wg.Done()

		dec := json.NewDecoder(resp.Body)
		for i := 0; ; i++ {
			var item struct {
				Result json.RawMessage        `json:"result"`
				Error  map[string]interface{} `json:"error"`
			}
			err := dec.Decode(&item)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("dec.Decode(&item) failed with %v; want success; i = %d", err, i)
			}
			if len(item.Error) != 0 {
				t.Errorf("item.Error = %#v; want empty; i = %d", item.Error, i)
				continue
			}
			var msg sub.StringMessage
			if err := jsonpb.UnmarshalString(string(item.Result), &msg); err != nil {
				t.Errorf("jsonpb.UnmarshalString(%q, &msg) failed with %v; want success", item.Result, err)
			}
			got = append(got, &msg)
		}
	}()

	wg.Wait()
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf(diff)
	}
}

func testABEBulkEchoZeroLength(t *testing.T, port int) {
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/echo", port)
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(nil))
	if err != nil {
		t.Errorf("http.NewRequest(%q, %q, bytes.NewReader(nil)) failed with %v; want success", "POST", apiURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "chunked")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http.Post(%q, %q, req) failed with %v; want success", apiURL, "application/json", err)
		return
	}
	defer resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
	}

	dec := json.NewDecoder(resp.Body)
	var item struct {
		Result json.RawMessage        `json:"result"`
		Error  map[string]interface{} `json:"error"`
	}
	if err := dec.Decode(&item); err == nil {
		t.Errorf("dec.Decode(&item) succeeded; want io.EOF; item = %#v", item)
	} else if err != io.EOF {
		t.Errorf("dec.Decode(&item) failed with %v; want success", err)
		return
	}
}

func testAdditionalBindings(t *testing.T, port int) {
	for i, f := range []func() *http.Response{
		func() *http.Response {
			apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/echo/hello", port)
			resp, err := http.Get(apiURL)
			if err != nil {
				t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
				return nil
			}
			return resp
		},
		func() *http.Response {
			apiURL := fmt.Sprintf("http://localhost:%d/v2/example/echo", port)
			resp, err := http.Post(apiURL, "application/json", strings.NewReader(`"hello"`))
			if err != nil {
				t.Errorf("http.Post(%q, %q, %q) failed with %v; want success", apiURL, "application/json", `"hello"`, err)
				return nil
			}
			return resp
		},
		func() *http.Response {
			r, w := io.Pipe()
			go func() {
				defer w.Close()
				w.Write([]byte(`"hello"`))
			}()
			apiURL := fmt.Sprintf("http://localhost:%d/v2/example/echo", port)
			resp, err := http.Post(apiURL, "application/json", r)
			if err != nil {
				t.Errorf("http.Post(%q, %q, %q) failed with %v; want success", apiURL, "application/json", `"hello"`, err)
				return nil
			}
			return resp
		},
		func() *http.Response {
			apiURL := fmt.Sprintf("http://localhost:%d/v2/example/echo?value=hello", port)
			resp, err := http.Get(apiURL)
			if err != nil {
				t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
				return nil
			}
			return resp
		},
	} {
		resp := f()
		if resp == nil {
			continue
		}

		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success; i=%d", err, i)
			return
		}
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("resp.StatusCode = %d; want %d; i=%d", got, want, i)
			t.Logf("%s", buf)
		}

		var msg sub.StringMessage
		if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
			t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success; %d", buf, err, i)
			return
		}
		if got, want := msg.GetValue(), "hello"; got != want {
			t.Errorf("msg.GetValue() = %q; want %q", got, want)
		}
	}
}

func testABERepeated(t *testing.T, port int) {
	f := func(v reflect.Value) string {
		var f func(v reflect.Value, idx int) string
		s := make([]string, v.Len())
		switch v.Index(0).Kind() {
		case reflect.Slice:
			f = func(v reflect.Value, idx int) string {
				t := v.Index(idx).Type().Elem().Kind()
				if t == reflect.Uint8 {
					return base64.URLEncoding.EncodeToString(v.Index(idx).Interface().([]byte))
				}
				// Could handle more elegantly
				panic("unknown slice of type: " + t.String())
			}
		default:
			f = func(v reflect.Value, idx int) string {
				return fmt.Sprintf("%v", v.Index(idx).Interface())
			}
		}
		for i := 0; i < v.Len(); i++ {
			s[i] = f(v, i)
		}
		return strings.Join(s, ",")
	}
	want := examplepb.ABitOfEverythingRepeated{
		PathRepeatedFloatValue: []float32{
			1.5,
			-1.5,
		},
		PathRepeatedDoubleValue: []float64{
			2.5,
			-2.5,
		},
		PathRepeatedInt64Value: []int64{
			4294967296,
			-4294967296,
		},
		PathRepeatedUint64Value: []uint64{
			0,
			9223372036854775807,
		},
		PathRepeatedInt32Value: []int32{
			2147483647,
			-2147483648,
		},
		PathRepeatedFixed64Value: []uint64{
			0,
			9223372036854775807,
		},
		PathRepeatedFixed32Value: []uint32{
			0,
			4294967295,
		},
		PathRepeatedBoolValue: []bool{
			true,
			false,
		},
		PathRepeatedStringValue: []string{
			"foo",
			"bar",
		},
		PathRepeatedBytesValue: [][]byte{
			[]byte{0x00},
			[]byte{0xFF},
		},
		PathRepeatedUint32Value: []uint32{
			0,
			4294967295,
		},
		PathRepeatedEnumValue: []examplepb.NumericEnum{
			examplepb.NumericEnum_ZERO,
			examplepb.NumericEnum_ONE,
		},
		PathRepeatedSfixed32Value: []int32{
			2147483647,
			-2147483648,
		},
		PathRepeatedSfixed64Value: []int64{
			4294967296,
			-4294967296,
		},
		PathRepeatedSint32Value: []int32{
			2147483647,
			-2147483648,
		},
		PathRepeatedSint64Value: []int64{
			4611686018427387903,
			-4611686018427387904,
		},
	}
	apiURL := fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything_repeated/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s/%s", port, f(reflect.ValueOf(want.PathRepeatedFloatValue)), f(reflect.ValueOf(want.PathRepeatedDoubleValue)), f(reflect.ValueOf(want.PathRepeatedInt64Value)), f(reflect.ValueOf(want.PathRepeatedUint64Value)), f(reflect.ValueOf(want.PathRepeatedInt32Value)), f(reflect.ValueOf(want.PathRepeatedFixed64Value)), f(reflect.ValueOf(want.PathRepeatedFixed32Value)), f(reflect.ValueOf(want.PathRepeatedBoolValue)), f(reflect.ValueOf(want.PathRepeatedStringValue)), f(reflect.ValueOf(want.PathRepeatedBytesValue)), f(reflect.ValueOf(want.PathRepeatedUint32Value)), f(reflect.ValueOf(want.PathRepeatedEnumValue)), f(reflect.ValueOf(want.PathRepeatedSfixed32Value)), f(reflect.ValueOf(want.PathRepeatedSfixed64Value)), f(reflect.ValueOf(want.PathRepeatedSint32Value)), f(reflect.ValueOf(want.PathRepeatedSint64Value)))

	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg examplepb.ABitOfEverythingRepeated
	if err := jsonpb.UnmarshalString(string(buf), &msg); err != nil {
		t.Errorf("jsonpb.UnmarshalString(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if diff := cmp.Diff(msg, want, protocmp.Transform()); diff != "" {
		t.Errorf(diff)
	}
}

func TestTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	apiURL := "http://localhost:8088/v2/example/timeout"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		t.Errorf(`http.NewRequest("GET", %q, nil) failed with %v; want success`, apiURL, err)
		return
	}
	req.Header.Set("Grpc-Timeout", "10m")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("http.DefaultClient.Do(%#v) failed with %v; want success", req, err)
		return
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusGatewayTimeout; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
	}
}

func TestPostWithEmptyBody(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	apiURL := "http://localhost:8088/v2/example/postwithemptybody/name"
	rep, err := http.Post(apiURL, "application/json", nil)

	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}

	if rep.StatusCode != http.StatusOK {
		t.Errorf("http.Post(%q) response code is %d; want %d", apiURL,
			rep.StatusCode, http.StatusOK)
		return
	}
}

func TestUnknownPath(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	apiURL := "http://localhost:8088"
	resp, err := http.Post(apiURL, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	apiURL := "http://localhost:8088/v1/example/echo/myid"
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusMethodNotAllowed; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}
}

func TestInvalidArgument(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	apiURL := "http://localhost:8088/v1/example/echo/myid/not_int64"
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}
}

func TestResponseBody(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	testResponseBody(t, 8088)
	testResponseBodies(t, 8088)
	testResponseStrings(t, 8088)
}

func testResponseBody(t *testing.T, port int) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
	}{{
		name:       "unary case",
		url:        "http://localhost:%d/responsebody/foo",
		wantStatus: http.StatusOK,
		wantBody:   `{"data":"foo"}`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL := fmt.Sprintf(tt.url, port)
			resp, err := http.Get(apiURL)
			if err != nil {
				t.Fatalf("http.Get(%q) failed with %v; want success", apiURL, err)
			}

			defer resp.Body.Close()
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
			}

			if got, want := resp.StatusCode, tt.wantStatus; got != want {
				t.Errorf("resp.StatusCode = %d; want %d", got, want)
				t.Logf("%s", buf)
			}

			if got, want := string(buf), tt.wantBody; got != want {
				t.Errorf("response = %q; want %q", got, want)
			}
		})
	}
}

func TestResponseBodyStream(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   []string
	}{{
		name:       "stream case",
		url:        "http://localhost:%d/responsebody/stream/foo",
		wantStatus: http.StatusOK,
		wantBody:   []string{`{"result":{"data":"first foo"}}`, `{"result":{"data":"second foo"}}`},
	}}

	port := 8088
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL := fmt.Sprintf(tt.url, port)
			resp, err := http.Get(apiURL)
			if err != nil {
				t.Fatalf("http.Get(%q) failed with %v; want success", apiURL, err)
			}

			defer resp.Body.Close()
			body, err := readAll(resp.Body)
			if err != nil {
				t.Fatalf("readAll(resp.Body) failed with %v; want success", err)
			}

			if got, want := resp.StatusCode, tt.wantStatus; got != want {
				t.Errorf("resp.StatusCode = %d; want %d", got, want)
			}

			if !reflect.DeepEqual(tt.wantBody, body) {
				t.Errorf("response = %v; want %v", body, tt.wantBody)
			}
		})
	}
}

func readAll(body io.ReadCloser) ([]string, error) {
	var b []string
	reader := bufio.NewReader(body)
	for {
		l, err := reader.ReadBytes('\n')
		switch {
		case err == io.EOF:
			return b, nil
		case err != nil:
			return nil, err
		}

		b = append(b, string(bytes.TrimSpace(l)))
	}
}

func testResponseBodies(t *testing.T, port int) {
	apiURL := fmt.Sprintf("http://localhost:%d/responsebodies/foo", port)
	resp, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	if got, want := string(buf), `[{"data":"foo"}]`; got != want {
		t.Errorf("response = %q; want %q", got, want)
	}
}

func testResponseStrings(t *testing.T, port int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port = 8087
	// Run Secondary server with different marshalling
	ch := make(chan error)
	go func() {
		if err := runGateway(ctx, fmt.Sprintf(":%d", port), runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EnumsAsInts: false, EmitDefaults: true})); err != nil {
			ch <- fmt.Errorf("cannot run gateway service: %v", err)
		}
	}()

	if err := waitForGateway(ctx, uint16(port)); err != nil {
		t.Fatalf("waitForGateway(ctx, %d) failed with %v; want success", port, err)
	}

	for i, spec := range []struct {
		endpoint     string
		expectedCode int
		expectedBody string
	}{
		{
			endpoint:     fmt.Sprintf("http://localhost:%d/responsestrings/foo", port),
			expectedCode: http.StatusOK,
			expectedBody: `["hello","foo"]`,
		},
		{
			endpoint:     fmt.Sprintf("http://localhost:%d/responsestrings/empty", port),
			expectedCode: http.StatusOK,
			expectedBody: `[]`,
		},
		{
			endpoint:     fmt.Sprintf("http://localhost:%d/responsebodies/foo", port),
			expectedCode: http.StatusOK,
			expectedBody: `[{"data":"foo","type":"UNKNOWN"}]`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			apiURL := spec.endpoint
			resp, err := http.Get(apiURL)
			if err != nil {
				t.Errorf("http.Get(%q) failed with %v; want success", apiURL, err)
				return
			}
			defer resp.Body.Close()
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
				return
			}

			if got, want := resp.StatusCode, spec.expectedCode; got != want {
				t.Errorf("resp.StatusCode = %d; want %d", got, want)
				t.Logf("%s", buf)
			}

			if got, want := string(buf), spec.expectedBody; got != want {
				t.Errorf("response = %q; want %q", got, want)
			}
		})
	}

}

func TestRequestQueryParams(t *testing.T) {
	testRequestQueryParams(t, 8088)
}

func TestRequestQueryParamsInProcessGateway(t *testing.T) {
	testRequestQueryParams(t, 8089)
}

func testRequestQueryParams(t *testing.T, port int) {
	if testing.Short() {
		t.Skip()
		return
	}

	formValues := url.Values{}
	formValues.Set("string_value", "hello-world")
	formValues.Add("repeated_string_value", "demo1")
	formValues.Add("repeated_string_value", "demo2")

	testCases := []struct {
		name           string
		httpMethod     string
		contentType    string
		apiURL         string
		wantContent    string
		requestContent io.Reader
	}{
		{
			name:        "get url query values",
			httpMethod:  "GET",
			contentType: "application/json",
			apiURL:      fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/params/get/foo?double_value=%v&bool_value=%v", port, 1234.56, true),
			wantContent: `{"single_nested":{"name":"foo"},"double_value":1234.56,"bool_value":true}`,
		},
		{
			name:        "get nested enum url parameter",
			httpMethod:  "GET",
			contentType: "application/json",
			// If nested_enum.OK were FALSE, the content of single_nested would be {} due to how 0 values are serialized
			apiURL:      fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/params/get/nested_enum/TRUE", port),
			wantContent: `{"single_nested":{"ok":"TRUE"}}`,
		},
		{
			name:           "post url query values",
			httpMethod:     "POST",
			contentType:    "application/json",
			apiURL:         fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/params/post/hello-world?double_value=%v&bool_value=%v", port, 1234.56, true),
			wantContent:    `{"single_nested":{"name":"foo","amount":100},"double_value":1234.56,"bool_value":true,"string_value":"hello-world"}`,
			requestContent: strings.NewReader(`{"name":"foo","amount":100}`),
		},
		{
			name:           "post form and url query values",
			httpMethod:     "POST",
			contentType:    "application/x-www-form-urlencoded",
			apiURL:         fmt.Sprintf("http://localhost:%d/v1/example/a_bit_of_everything/params/get/foo?double_value=%v&bool_value=%v", port, 1234.56, true),
			wantContent:    `{"single_nested":{"name":"foo"},"double_value":1234.56,"bool_value":true,"string_value":"hello-world","repeated_string_value":["demo1","demo2"]}`,
			requestContent: strings.NewReader(formValues.Encode()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.httpMethod, tc.apiURL, tc.requestContent)
			if err != nil {
				t.Errorf("http.method (%q) http.url (%q) failed with %v; want success", tc.httpMethod, tc.apiURL, err)
				return
			}

			req.Header.Add("Content-Type", tc.contentType)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("http.method (%q) http.url (%q) failed with %v; want success", tc.httpMethod, tc.apiURL, err)
				return
			}
			defer resp.Body.Close()

			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
				return
			}

			if gotCode, wantCode := resp.StatusCode, http.StatusOK; gotCode != wantCode {
				t.Errorf("resp.StatusCode = %d; want %d", gotCode, wantCode)
				t.Logf("%s", buf)
			}

			gotContent := string(buf)
			if gotContent != tc.wantContent {
				t.Errorf("http.method (%q) http.url (%q) response = %q; want %q", tc.httpMethod, tc.apiURL, gotContent, tc.wantContent)
			}
		})
	}
}

func TestNonStandardNames(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := runGateway(
			ctx,
			":8081",
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
		); err != nil {
			t.Errorf("runGateway() failed with %v; want success", err)
			return
		}
	}()
	go func() {
		if err := runGateway(
			ctx,
			":8082",
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: false, EmitDefaults: true}),
		); err != nil {
			t.Errorf("runGateway() failed with %v; want success", err)
			return
		}
	}()

	if err := waitForGateway(ctx, 8081); err != nil {
		t.Errorf("waitForGateway(ctx, 8081) failed with %v; want success", err)
	}
	if err := waitForGateway(ctx, 8082); err != nil {
		t.Errorf("waitForGateway(ctx, 8082) failed with %v; want success", err)
	}

	for _, tc := range []struct {
		name     string
		port     int
		method   string
		jsonBody string
	}{
		{
			"Test standard update method",
			8081,
			"update",
			`{"id":"foo","Num":"1","line_num":"42","langIdent":"English","STATUS":"good","en_GB":"1","no":"yes","thing":{"subThing":{"sub_value":"hi"}}}`,
		},
		{
			"Test update method using json_names in message",
			8081,
			"update_with_json_names",
			// N.B. json_names have no effect if not using OrigName: false
			`{"id":"foo","Num":"1","line_num":"42","langIdent":"English","STATUS":"good","en_GB":"1","no":"yes","thing":{"subThing":{"sub_value":"hi"}}}`,
		},
		{
			"Test standard update method with OrigName: false marshaller option",
			8082,
			"update",
			`{"id":"foo","Num":"1","lineNum":"42","langIdent":"English","STATUS":"good","enGB":"1","no":"yes","thing":{"subThing":{"subValue":"hi"}}}`,
		},
		{
			"Test update method using json_names in message with OrigName: false marshaller option",
			8082,
			"update_with_json_names",
			`{"ID":"foo","Num":"1","LineNum":"42","langIdent":"English","status":"good","En_GB":"1","yes":"no","Thingy":{"SubThing":{"sub_Value":"hi"}}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testNonStandardNames(t, tc.port, tc.method, tc.jsonBody)
		})
	}
}

func testNonStandardNames(t *testing.T, port int, method string, jsonBody string) {
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("http://localhost:%d/v1/example/non_standard/%s", port, method),
		strings.NewReader(jsonBody),
	)
	if err != nil {
		t.Fatalf("http.NewRequest(PATCH) failed with %v; want success", err)
	}
	patchResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to issue PATCH request: %v", err)
	}

	body, err := ioutil.ReadAll(patchResp.Body)
	if err != nil {
		t.Errorf("patchResp body couldn't be read: %v", err)
	}

	if got, want := patchResp.StatusCode, http.StatusOK; got != want {
		t.Errorf("patchResp.StatusCode= %d; want %d resp: %v", got, want, string(body))
	}

	if got, want := string(body), jsonBody; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}
