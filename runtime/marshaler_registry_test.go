package runtime_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func TestMarshalerForRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}
	r.Header.Set("Accept", "application/x-out")
	r.Header.Set("Content-Type", "application/x-in")

	mux := runtime.NewServeMux()

	out := runtime.MarshalerForRequest(mux, r)
	in := runtime.UnmarshalerForRequest(mux, r)
	if _, ok := in.(*runtime.JSONPb); !ok {
		t.Errorf("in = %#v; want a runtime.JSONPb", in)
	}
	if _, ok := out.(*runtime.JSONPb); !ok {
		t.Errorf("out = %#v; want a runtime.JSONPb", in)
	}

	var marshalers [3]dummyMarshaler
	specs := []struct {
		optIn  runtime.ServeMuxOption
		optOut runtime.ServeMuxOption

		wantIn  runtime.Unmarshaler
		wantOut runtime.Marshaler
	}{
		{
			optIn:   runtime.WithMarshalerOption(runtime.MIMEWildcard, &marshalers[0]),
			optOut:  runtime.WithUnmarshalerOption(runtime.MIMEWildcard, &marshalers[0]),
			wantIn:  &marshalers[0],
			wantOut: &marshalers[0],
		},
		{
			optIn:   runtime.WithMarshalerOption("application/x-in", &marshalers[1]),
			optOut:  runtime.WithUnmarshalerOption("application/x-in", &marshalers[1]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[0],
		},
		{
			optIn:   runtime.WithMarshalerOption("application/x-out", &marshalers[2]),
			optOut:  runtime.WithUnmarshalerOption("application/x-out", &marshalers[2]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[2],
		},
	}
	for i, spec := range specs {
		var opts []runtime.ServeMuxOption
		for _, s := range specs[:i+1] {
			opts = append(opts, s.optIn, s.optOut)
		}
		mux = runtime.NewServeMux(opts...)

		out = runtime.MarshalerForRequest(mux, r)
		in = runtime.UnmarshalerForRequest(mux, r)
		if got, want := in, spec.wantIn; got != want {
			t.Errorf("in = %#v; want %#v", got, want)
		}
		if got, want := out, spec.wantOut; got != want {
			t.Errorf("out = %#v; want %#v", got, want)
		}
	}

	r.Header.Set("Content-Type", "application/x-another")
	out = runtime.MarshalerForRequest(mux, r)
	in = runtime.UnmarshalerForRequest(mux, r)
	if got, want := in, &marshalers[1]; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, &marshalers[0]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
	}
}

type dummyMarshaler struct{}

func (dummyMarshaler) ContentType() string { return "" }
func (dummyMarshaler) Marshal(context.Context, interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (dummyMarshaler) Unmarshal(context.Context, []byte, interface{}) error {
	return errors.New("not implemented")
}

func (dummyMarshaler) NewDecoder(ctx context.Context, r io.Reader) runtime.Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler) NewEncoder(ctx context.Context, w io.Writer) runtime.Encoder {
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
