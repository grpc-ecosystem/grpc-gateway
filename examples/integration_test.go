package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	gw "github.com/gengo/grpc-gateway/examples/examplepb"
	server "github.com/gengo/grpc-gateway/examples/server"
	sub "github.com/gengo/grpc-gateway/examples/sub"
	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/protobuf/proto"
)

type aBitOfEverything struct {
	Result gw.ABitOfEverything          `json:"result,omitempty"`
	Error  *runtime.ResponseStreamError `json:"error,omitempty"`
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	go func() {
		if err := server.Run(); err != nil {
			t.Errorf("server.Run() failed with %v; want success", err)
			return
		}
	}()
	go func() {
		if err := Run(
			":8080",
			runtime.WithForwardResponseOption(
				func(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
					if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
						for k, vs := range md.HeaderMD {
							for i := range vs {
								w.Header().Add(fmt.Sprintf("Grpc-Header-%s", k), vs[i])
							}
						}
						for k, vs := range md.TrailerMD {
							for i := range vs {
								w.Header().Add(fmt.Sprintf("Grpc-Trailer-%s", k), vs[i])
							}
						}
					}
					return nil
				},
			),
		); err != nil {
			t.Errorf("gw.Run() failed with %v; want success", err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	testEcho(t, 8080, "application/json")
	testEchoBody(t)
	testABECreate(t)
	testABECreateBody(t)
	testABEBulkCreate(t)
	testABELookup(t)
	testABEList(t)
	testABEListError(t)
	testAdditionalBindings(t)

	go func() {
		if err := Run(
			":8081",
			runtime.WithForwardResponseOption(
				func(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
					w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1.1+json")
					return nil
				},
			),
		); err != nil {
			t.Errorf("gw.Run() failed with %v; want success", err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	testEcho(t, 8081, "application/vnd.docker.plugins.v1.1+json")
}

func testEcho(t *testing.T, port int, contentType string) {
	url := fmt.Sprintf("http://localhost:%d/v1/example/echo/myid", port)
	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("iotuil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg gw.SimpleMessage
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got, want := msg.Id, "myid"; got != want {
		t.Errorf("msg.Id = %q; want %q", got, want)
	}

	if value := resp.Header.Get("Content-Type"); value != contentType {
		t.Errorf("Content-Type was %s, wanted %s", value, contentType)
	}
}

func testEchoBody(t *testing.T) {
	sent := gw.SimpleMessage{Id: "example"}
	buf, err := json.Marshal(sent)
	if err != nil {
		t.Fatalf("json.Marshal(%#v) failed with %v; want success", sent, err)
	}

	url := "http://localhost:8080/v1/example/echo_body"
	resp, err := http.Post(url, "", bytes.NewReader(buf))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("iotuil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var received gw.SimpleMessage
	if err := json.Unmarshal(buf, &received); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got, want := received, sent; !reflect.DeepEqual(got, want) {
		t.Errorf("msg.Id = %q; want %q", got, want)
	}

	if value := resp.Header.Get("Grpc-Header-foo"); value != "foo1" {
		t.Errorf("Grpc-Header-foo was %s, wanted %s", value, "foo1")
	}

	if value := resp.Header.Get("Grpc-Header-bar"); value != "bar1" {
		t.Errorf("Grpc-Header-bar was %s, wanted %s", value, "bar1")
	}

	if value := resp.Header.Get("Grpc-Trailer-foo"); value != "foo2" {
		t.Errorf("Grpc-Trailer-foo was %s, wanted %s", value, "foo2")
	}

	if value := resp.Header.Get("Grpc-Trailer-bar"); value != "bar2" {
		t.Errorf("Grpc-Trailer-bar was %s, wanted %s", value, "bar2")
	}
}

func testABECreate(t *testing.T) {
	want := gw.ABitOfEverything{
		FloatValue:    1.5,
		DoubleValue:   2.5,
		Int64Value:    4294967296,
		Uint64Value:   9223372036854775807,
		Int32Value:    -2147483648,
		Fixed64Value:  9223372036854775807,
		Fixed32Value:  4294967295,
		BoolValue:     true,
		StringValue:   "strprefix/foo",
		Uint32Value:   4294967295,
		Sfixed32Value: 2147483647,
		Sfixed64Value: -4611686018427387904,
		Sint32Value:   2147483647,
		Sint64Value:   4611686018427387903,
	}
	url := fmt.Sprintf("http://localhost:8080/v1/example/a_bit_of_everything/%f/%f/%d/separator/%d/%d/%d/%d/%v/%s/%d/%d/%d/%d/%d", want.FloatValue, want.DoubleValue, want.Int64Value, want.Uint64Value, want.Int32Value, want.Fixed64Value, want.Fixed32Value, want.BoolValue, want.StringValue, want.Uint32Value, want.Sfixed32Value, want.Sfixed64Value, want.Sint32Value, want.Sint64Value)

	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("iotuil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg gw.ABitOfEverything
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if msg.Uuid == "" {
		t.Error("msg.Uuid is empty; want not empty")
	}
	msg.Uuid = ""
	if got := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("msg= %v; want %v", &got, &want)
	}
}

func testABECreateBody(t *testing.T) {
	want := gw.ABitOfEverything{
		FloatValue:    1.5,
		DoubleValue:   2.5,
		Int64Value:    4294967296,
		Uint64Value:   9223372036854775807,
		Int32Value:    -2147483648,
		Fixed64Value:  9223372036854775807,
		Fixed32Value:  4294967295,
		BoolValue:     true,
		StringValue:   "strprefix/foo",
		Uint32Value:   4294967295,
		Sfixed32Value: 2147483647,
		Sfixed64Value: -4611686018427387904,
		Sint32Value:   2147483647,
		Sint64Value:   4611686018427387903,

		Nested: []*gw.ABitOfEverything_Nested{
			{
				Name:   "bar",
				Amount: 10,
			},
			{
				Name:   "baz",
				Amount: 20,
			},
		},
	}
	url := "http://localhost:8080/v1/example/a_bit_of_everything"
	buf, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json.Marshal(%#v) failed with %v; want success", want, err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("iotuil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg gw.ABitOfEverything
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if msg.Uuid == "" {
		t.Error("msg.Uuid is empty; want not empty")
	}
	msg.Uuid = ""
	if got := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("msg= %v; want %v", &got, &want)
	}
}

func testABEBulkCreate(t *testing.T) {
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
			want := gw.ABitOfEverything{
				FloatValue:    1.5,
				DoubleValue:   2.5,
				Int64Value:    4294967296,
				Uint64Value:   9223372036854775807,
				Int32Value:    -2147483648,
				Fixed64Value:  9223372036854775807,
				Fixed32Value:  4294967295,
				BoolValue:     true,
				StringValue:   fmt.Sprintf("strprefix/%s", val),
				Uint32Value:   4294967295,
				Sfixed32Value: 2147483647,
				Sfixed64Value: -4611686018427387904,
				Sint32Value:   2147483647,
				Sint64Value:   4611686018427387903,

				Nested: []*gw.ABitOfEverything_Nested{
					{
						Name:   "hoge",
						Amount: 10,
					},
					{
						Name:   "fuga",
						Amount: 20,
					},
				},
			}
			buf, err := json.Marshal(want)
			if err != nil {
				t.Fatalf("json.Marshal(%#v) failed with %v; want success", want, err)
			}
			if _, err := w.Write(buf); err != nil {
				t.Errorf("w.Write(%s) failed with %v; want success", buf, err)
				return
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				t.Errorf("w.Write(%s) failed with %v; want success", buf, err)
				return
			}
		}
	}(w)
	url := "http://localhost:8080/v1/example/a_bit_of_everything/bulk"
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("iotuil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
	}

	var msg gw.EmptyMessage
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}

	if value := resp.Header.Get("Grpc-Trailer-foo"); value != "foo2" {
		t.Errorf("Grpc-Trailer-foo was %q, wanted %q", value, "foo2")
	}

	if value := resp.Header.Get("Grpc-Trailer-bar"); value != "bar2" {
		t.Errorf("Grpc-Trailer-bar was %q, wanted %q", value, "bar2")
	}
}

func testABELookup(t *testing.T) {
	url := "http://localhost:8080/v1/example/a_bit_of_everything"
	cresp, err := http.Post(url, "application/json", strings.NewReader(`
		{"bool_value": true, "string_value": "strprefix/example"}
	`))
	if err != nil {
		t.Errorf("http.Post(%q) failed with %v; want success", url, err)
		return
	}
	defer cresp.Body.Close()
	buf, err := ioutil.ReadAll(cresp.Body)
	if err != nil {
		t.Errorf("iotuil.ReadAll(cresp.Body) failed with %v; want success", err)
		return
	}
	if got, want := cresp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
		t.Logf("%s", buf)
		return
	}

	var want gw.ABitOfEverything
	if err := json.Unmarshal(buf, &want); err != nil {
		t.Errorf("json.Unmarshal(%s, &want) failed with %v; want success", buf, err)
		return
	}

	url = fmt.Sprintf("%s/%s", url, want.Uuid)
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(resp.Body) failed with %v; want success", err)
		return
	}

	var msg gw.ABitOfEverything
	if err := json.Unmarshal(buf, &msg); err != nil {
		t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success", buf, err)
		return
	}
	if got := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("msg= %v; want %v", &got, &want)
	}
}

func testABEList(t *testing.T) {
	url := "http://localhost:8080/v1/example/a_bit_of_everything"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("http.Get(%q) failed with %v; want success", url, err)
		return
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var i int
	for i = 0; ; i++ {
		var msg aBitOfEverything
		err := dec.Decode(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("dec.Decode(&msg) failed with %v; want success; i = %d", err, i)
		}
	}
	if i <= 0 {
		t.Errorf("i == %d; want > 0", i)
	}
}

func testABEListError(t *testing.T) {
	url := "http://localhost:8080/v1/example/a_bit_of_everything"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("http.NewReuest(%q) failed with %v; want success", url, err)
		return
	}
	req.Header.Set("Grpc-Metadata-error", "foo")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("client.Do failed with %v; want success", err)
		return
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var i int
	var lastMsg aBitOfEverything
	for i = 0; ; i++ {
		var msg aBitOfEverything
		err := dec.Decode(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("dec.Decode(&msg) failed with %v; want success; i = %d", err, i)
		}
		lastMsg = msg
	}
	if i <= 0 {
		t.Errorf("i == %d; want > 0", i)
	}

	if got, want := lastMsg.Error.HTTPCode, http.StatusBadRequest; got != want {
		t.Errorf("lastMsg.Error.HTTPCode = %d; want %d", got, want)
		return
	}

	md := lastMsg.Error.Trailer
	v, ok := md["foo"]
	if !ok || len(v) == 0 {
		t.Errorf("Trailer doesn't contain %q", "foo")
	}
	if !ok && v[0] != "foo1" {
		t.Errorf("Trailer %q = %q; want %q", "foo", v[0], "foo2")
	}

	v, ok = md["bar"]
	if !ok || len(v) == 0 {
		t.Errorf("Trailer doesn't contain %q", "bar")
	}
	if !ok && v[0] != "bar1" {
		t.Errorf("Trailer %q = %q; want %q", "bar", v[0], "bar2")
	}
}

func testAdditionalBindings(t *testing.T) {
	for i, f := range []func() *http.Response{
		func() *http.Response {
			url := "http://localhost:8080/v1/example/a_bit_of_everything/echo/hello"
			resp, err := http.Get(url)
			if err != nil {
				t.Errorf("http.Get(%q) failed with %v; want success", url, err)
				return nil
			}
			return resp
		},
		func() *http.Response {
			url := "http://localhost:8080/v2/example/echo"
			resp, err := http.Post(url, "application/json", strings.NewReader(`"hello"`))
			if err != nil {
				t.Errorf("http.Post(%q, %q, %q) failed with %v; want success", url, "application/json", `"hello"`, err)
				return nil
			}
			return resp
		},
		func() *http.Response {
			url := "http://localhost:8080/v2/example/echo?value=hello"
			resp, err := http.Get(url)
			if err != nil {
				t.Errorf("http.Get(%q) failed with %v; want success", url, err)
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
			t.Errorf("iotuil.ReadAll(resp.Body) failed with %v; want success; i=%d", err, i)
			return
		}
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("resp.StatusCode = %d; want %d; i=%d", got, want, i)
			t.Logf("%s", buf)
		}

		var msg sub.StringMessage
		if err := json.Unmarshal(buf, &msg); err != nil {
			t.Errorf("json.Unmarshal(%s, &msg) failed with %v; want success; %i", buf, err, i)
			return
		}
		if got, want := msg.GetValue(), "hello"; got != want {
			t.Errorf("msg.GetValue() = %q; want %q", got, want)
		}
	}
}
