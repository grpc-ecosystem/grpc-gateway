package runtime_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gengo/grpc-gateway/examples/examplepb"
	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type stubResponseWriter struct {
	header  http.Header
	results bytes.Buffer
}

func (w *stubResponseWriter) Header() http.Header { return w.header }
func (w *stubResponseWriter) Write(p []byte) (int, error) {
	return w.results.Write(p)
}
func (w *stubResponseWriter) WriteHeader(int) {}

func TestForwardResponseMessage(t *testing.T) {
	writer := stubResponseWriter{header: http.Header{}}
	m1 := examplepb.ABitOfEverything{Uuid: "foo"}
	m1.Nested = append(m1.Nested, &examplepb.ABitOfEverything_Nested{Name: "bar"})

	runtime.ForwardResponseMessage(nil, &writer, nil, &m1)

	got := writer.results.String()
	want := "{\"uuid\":\"foo\",\"nested\":[{\"name\":\"bar\"}]}"
	if got != want {
		t.Errorf("Got %s (expected %s)", got, want)
	} else {
		glog.Infof("Got %s", got)
	}
	m2 := examplepb.ABitOfEverything{}
	err := json.Unmarshal(writer.results.Bytes(), &m2)
	if err != nil {
		t.Errorf("Invalid JSON: %v", err)
	}
	if !proto.Equal(&m2, &m1) {
		t.Errorf("Got %v (expected %v)", m2, m1)
	}
}

func TestForwardResponseMessage_WithResponseIndent(t *testing.T) {
	writer := stubResponseWriter{header: http.Header{}}
	m1 := examplepb.ABitOfEverything{Uuid: "foo"}
	m1.Nested = append(m1.Nested, &examplepb.ABitOfEverything_Nested{Name: "bar"})

	runtime.ResponseIndent = " "
	defer func() {
		runtime.ResponseIndent = ""
	}()
	runtime.ForwardResponseMessage(nil, &writer, nil, &m1)

	got := writer.results.String()
	want := "{\n \"uuid\": \"foo\",\n \"nested\": [\n  {\n   \"name\": \"bar\"\n  }\n ]\n}"
	if got != want {
		t.Errorf("Got:\n%s\nExpected:\n%s)", got, want)
	} else {
		glog.Infof("Got:\n%s", got)
	}
	m2 := examplepb.ABitOfEverything{}
	err := json.Unmarshal(writer.results.Bytes(), &m2)
	if err != nil {
		t.Errorf("Invalid JSON: %v", err)
	}
	if !proto.Equal(&m2, &m1) {
		t.Errorf("Got %v (expected %v)", m2, m1)
	}
}
