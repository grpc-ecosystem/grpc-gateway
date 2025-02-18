package gateway

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

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
