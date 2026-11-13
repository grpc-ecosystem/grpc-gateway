---
layout: default
title: ETag and If-None-Match
nav_order: 7
parent: Operations
---

# ETag and If-None-Match

The gateway doesn't add `ETag` headers or handle conditional `GET` requests on its own, but because it's a plain `http.Handler`, you can wrap the `ServeMux` to do it. The handler below buffers each successful `GET` or `HEAD` response, sets an `ETag` computed from the body, and answers `304 Not Modified` with no body when the request's `If-None-Match` header matches it. Error responses and other methods pass through untouched.

```go
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

// etagWriter buffers a successful response so an ETag can be computed over
// the whole body before anything is sent to the client.
type etagWriter struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (w *etagWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
}

func (w *etagWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.buf.Write(b)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (w *etagWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// ifNoneMatch reports whether the If-None-Match header matches etag, using
// the weak comparison that RFC 9110 asks for on this header.
func ifNoneMatch(header, etag string) bool {
	for _, candidate := range strings.Split(header, ",") {
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || strings.TrimPrefix(candidate, "W/") == etag {
			return true
		}
	}
	return false
}

// ETagHandler adds an ETag to successful GET and HEAD responses of the
// wrapped handler, and answers 304 Not Modified when the request's
// If-None-Match header matches it. Responses are buffered to compute the
// ETag, so don't use it in front of server-streaming methods.
func ETagHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			h.ServeHTTP(w, r)
			return
		}

		ew := &etagWriter{ResponseWriter: w}
		h.ServeHTTP(ew, r)
		if ew.status == 0 {
			ew.status = http.StatusOK
		}
		if ew.status != http.StatusOK {
			w.WriteHeader(ew.status)
			_, _ = w.Write(ew.buf.Bytes())
			return
		}

		sum := sha256.Sum256(ew.buf.Bytes())
		etag := `"` + hex.EncodeToString(sum[:16]) + `"`
		w.Header().Set("ETag", etag)
		if ifNoneMatch(r.Header.Get("If-None-Match"), etag) {
			w.Header().Del("Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return
		}
		_, _ = w.Write(ew.buf.Bytes())
	})
}
```

Wrap the mux when you create the server:

```go
mux := runtime.NewServeMux()
// Register your generated gateway handlers on mux here.

s := &http.Server{
	Addr:    ":8081",
	Handler: ETagHandler(mux),
}
```

A few things to keep in mind:

- The response is buffered so the ETag covers the whole body, even if it's written in several chunks. Don't put the handler in front of server-streaming methods, since it would hold every message until the stream ends. Wrap only the routes that return a single message instead, for example with a separate mux.
- The ETag is derived from the response bytes, so the marshaler has to produce the same output for the same message. The default `JSONPb` marshaler does, but if you use a custom one, make sure it doesn't add anything like timestamps.
- The backend still runs the full request: this saves bandwidth, not work. If your service can tell cheaply whether a resource changed, for example from a version field, you can set the ETag from that in a `ForwardResponseOption` instead.

Thanks to @joshgarnett for the original version of this handler in [#4263](https://github.com/grpc-ecosystem/grpc-gateway/issues/4263).
