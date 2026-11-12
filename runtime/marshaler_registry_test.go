package runtime_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func TestMarshalerForRequest(t *testing.T) {
	ctx := context.Background()
	r, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}

	mux := runtime.NewServeMux()

	r.Header.Set("Accept", "application/x-out")
	r.Header.Set("Content-Type", "application/x-in")
	in, out := runtime.MarshalerForRequest(mux, r)
	if _, ok := in.(*runtime.HTTPBodyMarshaler); !ok {
		t.Errorf("in = %#v; want a runtime.HTTPBodyMarshaler", in)
	}
	if _, ok := out.(*runtime.HTTPBodyMarshaler); !ok {
		t.Errorf("out = %#v; want a runtime.HTTPBodyMarshaler", in)
	}

	marshalers := []dummyMarshaler{0, 1, 2}
	specs := []struct {
		opt runtime.ServeMuxOption

		wantIn  runtime.Marshaler
		wantOut runtime.Marshaler
	}{
		// The option with wildcard overwrites the default configuration
		{
			opt:     runtime.WithMarshalerOption(runtime.MIMEWildcard, &marshalers[0]),
			wantIn:  &marshalers[0],
			wantOut: &marshalers[0],
		},
		// You can specify a marshaler for a specific MIME type.
		// The output marshaler follows the input one unless specified.
		{
			opt:     runtime.WithMarshalerOption("application/x-in", &marshalers[1]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[1],
		},
		// You can also separately specify an output marshaler
		{
			opt:     runtime.WithMarshalerOption("application/x-out", &marshalers[2]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[2],
		},
	}
	for i, spec := range specs {
		var opts []runtime.ServeMuxOption
		for _, s := range specs[:i+1] {
			opts = append(opts, s.opt)
		}
		mux = runtime.NewServeMux(opts...)

		in, out = runtime.MarshalerForRequest(mux, r)
		if got, want := in, spec.wantIn; got != want {
			t.Errorf("in = %#v; want %#v", got, want)
		}
		if got, want := out, spec.wantOut; got != want {
			t.Errorf("out = %#v; want %#v", got, want)
		}
	}

	r.Header.Set("Content-Type", "application/x-in; charset=UTF-8")
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, &marshalers[1]; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, &marshalers[2]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "application/x-another")
	r.Header.Set("Accept", "application/x-another")
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, &marshalers[0]; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, &marshalers[0]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}
}

type dummyMarshaler int

func (dummyMarshaler) ContentType(_ interface{}) string { return "" }
func (dummyMarshaler) Marshal(interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (dummyMarshaler) Unmarshal([]byte, interface{}) error {
	return errors.New("not implemented")
}

func (dummyMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return dummyEncoder{}
}

func (m dummyMarshaler) GoString() string {
	return fmt.Sprintf("dummyMarshaler(%d)", m)
}

type dummyDecoder struct{}

func (dummyDecoder) Decode(interface{}) error {
	return errors.New("not implemented")
}

type dummyEncoder struct{}

func (dummyEncoder) Encode(interface{}) error {
	return errors.New("not implemented")
}

// BenchmarkMarshalerForRequest covers the marshaller lookup every generated handler performs
// once per request, and the error paths perform again.
//
// The Content-Type is the bare media type clients send, which is also exactly how marshallers
// are registered, so this is the case worth keeping cheap.
func BenchmarkMarshalerForRequest(b *testing.B) {
	// A marshaller registered under the concrete media type, as the customisation examples in
	// the documentation do. The default registry holds only MIMEWildcard, in which case the
	// lookups below miss and the Content-Type still has to be parsed.
	mux := runtime.NewServeMux(runtime.WithMarshalerOption("application/json", &runtime.JSONPb{}))
	req, err := http.NewRequest(http.MethodPost, "https://example.com/v1/echo", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inbound, outbound := runtime.MarshalerForRequest(mux, req)
		if inbound == nil || outbound == nil {
			b.Fatal("got a nil marshaler")
		}
	}
}

// BenchmarkMarshalerForRequestWithParams keeps the path that still has to parse measured: a
// Content-Type carrying parameters is not a registry key, so it falls through to
// mime.ParseMediaType.
func BenchmarkMarshalerForRequestWithParams(b *testing.B) {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption("application/json", &runtime.JSONPb{}))
	req, err := http.NewRequest(http.MethodPost, "https://example.com/v1/echo", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inbound, outbound := runtime.MarshalerForRequest(mux, req)
		if inbound == nil || outbound == nil {
			b.Fatal("got a nil marshaler")
		}
	}
}
