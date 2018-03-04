package runtime_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func TestMarshalerForRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}
	r.Header.Set("Accept", "application/x-out")
	r.Header.Set("Content-Type", "application/x-in")

	mux := runtime.NewServeMux()

	in, out := runtime.MarshalerForRequest(mux, r)
	if _, ok := in.(*runtime.JSONPb); !ok {
		t.Errorf("in = %#v; want a runtime.JSONPb", in)
	}
	if _, ok := out.(*runtime.JSONPb); !ok {
		t.Errorf("out = %#v; want a runtime.JSONPb", in)
	}

	var marshalers [3]dummyMarshaler
	specs := []struct {
		opt runtime.ServeMuxOption

		wantIn  runtime.Marshaler
		wantOut runtime.Marshaler
	}{
		{
			opt:     runtime.WithMarshalerOption(runtime.MIMEWildcard, &marshalers[0]),
			wantIn:  &marshalers[0],
			wantOut: &marshalers[0],
		},
		{
			opt:     runtime.WithMarshalerOption("application/x-in", &marshalers[1]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[0],
		},
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

	r.Header.Set("Content-Type", "application/x-another")
	in, out = runtime.MarshalerForRequest(mux, r)
	if got, want := in, &marshalers[1]; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, &marshalers[0]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}
}

func TestMarshalerForRequest2(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}

	marshalers := [4]*dummyMarshaler2{{0}, {1}, {2}, {3}}
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("multipart/form-data; boundary=*", marshalers[0]),
		runtime.WithMarshalerOption("multipart/*", marshalers[1]),
		runtime.WithMarshalerOption("multipart/form-data", marshalers[2]),
		runtime.WithMarshalerOption("*/form-data", marshalers[3]),
	)

	r.Header.Set("Content-Type", "multipart/form-data")
	_, out := runtime.MarshalerForRequest(mux, r)
	if got, want := out, marshalers[2]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "multipart/form-data; boundary=foobar")
	_, out = runtime.MarshalerForRequest(mux, r)
	if got, want := out, marshalers[0]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "multipart/form-data; boundary=*")
	_, out = runtime.MarshalerForRequest(mux, r)
	if got, want := out, marshalers[0]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "multipart/octet-stream")
	_, out = runtime.MarshalerForRequest(mux, r)
	if got, want := out, marshalers[1]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "application/form-data")
	_, out = runtime.MarshalerForRequest(mux, r)
	if got, want := out, marshalers[3]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}

	r.Header.Set("Content-Type", "application/json")
	_, out = runtime.MarshalerForRequest(mux, r)
	if _, ok := out.(*runtime.JSONPb); !ok {
		t.Errorf("in = %#v; want a runtime.JSONPb", out)
	}

	r.Header.Del("Content-Type")
	_, out = runtime.MarshalerForRequest(mux, r)
	if _, ok := out.(*runtime.JSONPb); !ok {
		t.Errorf("in = %#v; want a runtime.JSONPb", out)
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

type dummyMarshaler2 struct {
	label int
}

func (dummyMarshaler2) ContentType() string { return "" }
func (dummyMarshaler2) Marshal(interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (dummyMarshaler2) Unmarshal([]byte, interface{}) error {
	return errors.New("not implemented")
}

func (dummyMarshaler2) NewDecoder(r io.Reader) runtime.Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler2) NewEncoder(w io.Writer) runtime.Encoder {
	return dummyEncoder{}
}
