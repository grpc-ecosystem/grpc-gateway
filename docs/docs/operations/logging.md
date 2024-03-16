---
layout: default
title: Logging the request body pattern for a request
nav_order: 5
parent: Operations
---

# Logging the request body pattern for a request

If you want to log the request body in the `customErrorHandler` middleware, unfortunately the request body has been consumed in the `customErrorHandler` middleware and can't be read again. To log the request body, you can use one middleware to buffer the request body before it's consumed.

1. `logRequestBody` middleware，which logs the request body when the response status code is not 200.

```go
type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rsp *logResponseWriter) WriteHeader(code int) {
	rsp.statusCode = code
	rsp.ResponseWriter.WriteHeader(code)
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

// logRequestBody logs the request body when the response status code is not 200.
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
```

2. Wrap the `logRequestBody` middleware with the `http.ServeMux`:

```go
    mux := http.NewServeMux()

    s := &http.Server{
        Handler: logRequestBody(mux),
    }
```