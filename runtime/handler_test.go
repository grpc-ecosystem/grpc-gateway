package runtime_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/proto"
	pb "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
)

func TestForwardResponseStream(t *testing.T) {
	var (
		msgs = []proto.Message{
			&pb.SimpleMessage{Id: "One"},
			&pb.SimpleMessage{Id: "Two"},
		}

		ctx = runtime.NewServerMetadataContext(
			context.Background(), runtime.ServerMetadata{},
		)
		mux       = runtime.NewServeMux()
		marshaler = &runtime.JSONPb{}
		req       = httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp      = httptest.NewRecorder()
		count     = 0
		recv      = func() (proto.Message, error) {
			if count >= len(msgs) {
				return nil, io.EOF
			}
			count++
			return msgs[count-1], nil
		}
	)

	runtime.ForwardResponseStream(ctx, mux, marshaler, resp, req, recv)

	w := resp.Result()
	if w.StatusCode != http.StatusOK {
		t.Errorf(" got %d want %d", w.StatusCode, http.StatusOK)
	}
	if h := w.Header.Get("Transfer-Encoding"); h != "chunked" {
		t.Errorf("ForwardResponseStream missing header chunked")
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Failed to read response body with %v", err)
	}
	w.Body.Close()

	var want []byte
	for _, msg := range msgs {
		b, err := marshaler.Marshal(map[string]proto.Message{"result": msg})
		if err != nil {
			t.Errorf("marshaler.Marshal() failed %v", err)
		}
		want = append(want, b...)
	}

	if string(body) != string(want) {
		t.Errorf("ForwardResponseStream() = \"%s\" want \"%s\"", body, want)
	}
}
