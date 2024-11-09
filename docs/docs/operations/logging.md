---
layout: default
title: Logging the request body pattern for a request
nav_order: 5
parent: Operations
---

# Logging the request body pattern for a request

If you want to log the request body of incoming requests, you will need to buffer the body before it reaches the gateway. To log the request body, you can use a middleware `http.Handler` to buffer the request body before it's consumed.

1. Create a `http.Handler` middleware. The `logRequestBody` example middleware logs the request body when the response status code is not 200.

```go
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
func logRequestBody(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := newLogResponseWriter(w)

		// Note that buffering the entire request body could consume a lot of memory.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read body: %v", err), http.StatusBadRequest)
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

2. Wrap the gateway serve mux with the `logRequestBody` middleware:

```go
    mux := runtime.NewServeMux()
    // Register generated gateway handlers

    s := &http.Server{
        Handler: logRequestBody(mux),
    }
```
