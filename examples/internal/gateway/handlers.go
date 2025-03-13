package gateway

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/emptypb"
)

// openAPIServer returns OpenAPI specification files located under "/openapiv2/"
func openAPIServer(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
			grpclog.Errorf("Not Found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		grpclog.Infof("Serving %s", r.URL.Path)
		p := strings.TrimPrefix(r.URL.Path, "/openapiv2/")
		p = path.Join(dir, p)
		http.ServeFile(w, r, p)
	}
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	grpclog.Infof("Preflight request for %s", r.URL.Path)
}

// healthzServer returns a simple health handler which returns ok.
func healthzServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s == connectivity.Idle {
			// Invoke method to move connection from Idle to Ready
			conn.Invoke(r.Context(), "/grpc.health.v1.Health/Check", &emptypb.Empty{}, &emptypb.Empty{})
		}
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rsp *logResponseWriter) WriteHeader(code int) {
	rsp.statusCode = code
	rsp.ResponseWriter.WriteHeader(code)
}

// Unwrap returns the original http.ResponseWriter. This is necessary
// to expose Flush() and Push() on the underlying response writer.
func (rsp *logResponseWriter) Unwrap() http.ResponseWriter {
	return rsp.ResponseWriter
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

// logRequestBody logs the request body when the response status code is not 200.
// This addresses the issue of being unable to retrieve the request body in the customErrorHandler middleware.
func logRequestBody(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := newLogResponseWriter(w)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("grpc server read request body err %+v", err), http.StatusBadRequest)
			return
		}
		clonedR := r.Clone(r.Context())
		clonedR.Body = io.NopCloser(bytes.NewReader(body))

		h.ServeHTTP(lw, clonedR)

		if lw.statusCode != http.StatusOK {
			grpclog.Errorf("http error %+v request body %+v", lw.statusCode, string(body))
		}
	})
}

// websocketGateway transparently upgrades WebSocket connections and forwards
// their traffic to the given handler, each WebSocket message corresponding to
// a line in the HTTP request/response body.
func websocketGateway(h http.Handler) http.Handler {
	upgrader := websocket.Upgrader{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "websocket" {
			h.ServeHTTP(w, r)
			return
		}

		// Upgrade connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
			return
		}

		// WebSocket connection / HTTP request pipes for passing data around
		rReq, wReq := io.Pipe()
		rResp, wResp := io.Pipe()

		// Close request pipe when connection is closed
		conn.SetCloseHandler(func(code int, text string) error {
			return wReq.Close()
		})

		// Read from conn, write to wReq.
		// Close wReq when conn is closed.
		go func() {
			defer wReq.Close()
			for {
				t, p, err := conn.ReadMessage()
				if err != nil {
					return
				}
				switch t {
				case websocket.TextMessage:
					_, err = wReq.Write(p)
				case websocket.BinaryMessage:
					_, err = wReq.Write(p)
				case websocket.CloseMessage:
					return
				default:
					continue
				}
				if err != nil {
					return
				}
			}
		}()

		// Read from rResp, write to conn.
		// Close conn when rResp is closed.
		go func() {
			defer conn.Close()
			s := bufio.NewScanner(rResp)
			for s.Scan() {
				err := conn.WriteMessage(websocket.TextMessage, s.Bytes())
				if err != nil {
					break
				}
			}
			if err != nil {
				conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()),
				)
			} else if s.Err() != nil {
				conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseGoingAway, err.Error()),
				)
			} else {
				conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				)
			}
		}()

		// Replace request body with one that reads from rReq
		r, err = http.NewRequestWithContext(r.Context(), http.MethodPost, r.URL.String(), eofPipeReader{rReq})
		if err != nil {
			return
		}

		// Replace response writer with one that writes to wResp
		w = &responseForwarder{wResp, w.Header()}

		// Execute HTTP request with wrapped request and response
		h.ServeHTTP(w, r)

		// Close wResp to signal end of response
		wResp.Close()
	})
}

// Implementation of http.ResponseWriter interface that redirects all response
// data to embedded *io.PipeWriter.
type responseForwarder struct {
	*io.PipeWriter
	h http.Header
}

func (rf *responseForwarder) Header() http.Header {
	return rf.h
}

func (rf *responseForwarder) WriteHeader(int) {
}

func (rf *responseForwarder) Flush() {
}

// Wrapper around *io.PipeReader that returns io.EOF when underlying
// *io.PipeReader is closed.
type eofPipeReader struct {
	*io.PipeReader
}

func (r eofPipeReader) Read(p []byte) (int, error) {
	n, err := r.PipeReader.Read(p)
	if err == io.ErrClosedPipe {
		err = io.EOF
	}
	return n, err
}
