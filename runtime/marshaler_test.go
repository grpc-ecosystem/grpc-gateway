package runtime_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gengo/grpc-gateway/runtime"
)

func TestMarshalerForRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}
	r.Header.Set("Content-Type", "application/x-example")

	reg := runtime.NewMarshalerMIMERegistry()
	mux := runtime.NewServeMux()
	mux.MIMERegistry = reg

	in, out := runtime.MarshalerForRequest(mux, r)
	if _, ok := in.(*runtime.JSONBuiltin); !ok {
		t.Errorf("in = %#v; want a runtime.JSONBuiltin", in)
	}
	if _, ok := out.(*runtime.JSONBuiltin); !ok {
		t.Errorf("out = %#v; want a runtime.JSONBuiltin", in)
	}

	m1 := new(dummyMarshaler)
	reg.AddMarshaler("*", m1, m1)
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, m1; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, m1; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	m2 := new(dummyMarshaler)
	reg.AddInboundMarshaler("*", m2)
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, m2; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, m1; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	m3 := new(dummyMarshaler)
	reg.AddOutboundMarshaler("application/x-example", m3)
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, m2; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, m3; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	m4 := new(dummyMarshaler)
	reg.AddInboundMarshaler("application/x-example", m4)
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, m4; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, m3; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	m5, m6 := new(dummyMarshaler), new(dummyMarshaler)
	reg.AddMarshaler("application/x-example", m5, m6)
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, m5; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, m6; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "application/x-another")
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, m2; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, m1; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}
}

type dummyMarshaler struct{}

func (dummyMarshaler) ContentType() string { return "" }
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

type dummyDecoder struct{}

func (dummyDecoder) Decode(interface{}) error {
	return errors.New("not implemented")
}

type dummyEncoder struct{}

func (dummyEncoder) Encode(interface{}) error {
	return errors.New("not implemented")
}
