package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

type responseStreamChunk struct {
	Result proto.Message `json:"result,omitempty"`
	Error  string        `json:"error,omitempty"`
}

// ForwardResponseStream forwards the stream from gRPC server to REST client.
func ForwardResponseStream(w http.ResponseWriter, recv func() (proto.Message, error)) {
	f, ok := w.(http.Flusher)
	if !ok {
		glog.Errorf("Flush not supported in %T", w)
		http.Error(w, "unexpected type of web server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	f.Flush()
	for {
		resp, err := recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			buf, merr := json.Marshal(responseStreamChunk{Error: err.Error()})
			if merr != nil {
				glog.Errorf("Failed to marshal an error: %v", merr)
				return
			}
			if _, werr := fmt.Fprintf(w, "%s\n", buf); werr != nil {
				glog.Errorf("Failed to notify error to client: %v", werr)
				return
			}
			return
		}
		buf, err := json.Marshal(responseStreamChunk{Result: resp})
		if err != nil {
			glog.Errorf("Failed to marshal response chunk: %v", err)
			return
		}
		if _, err = fmt.Fprintf(w, "%s\n", buf); err != nil {
			glog.Errorf("Failed to send response chunk: %v", err)
			return
		}
	}
}

// ForwardResponseMessage forwards the message from gRPC server to REST client.
func ForwardResponseMessage(ctx context.Context, w http.ResponseWriter, resp proto.Message) {
	buf, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("Marshal error: %v", err)
		HTTPError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(buf); err != nil {
		glog.Errorf("Failed to write response: %v", err)
	}
}
