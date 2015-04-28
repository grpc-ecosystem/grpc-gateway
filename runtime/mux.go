package runtime

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
)

// A HandlerFunc handles a specific pair of path pattern and HTTP method.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)

// ServeMux is a request multiplexer for grpc-gateway.
// It matches http requests to patterns and invokes the corresponding handler.
type ServeMux struct {
	// handlers maps HTTP method to a list of handlers.
	handlers map[string][]handler
}

// NewServeMux returns a new MuxHandler whose internal mapping is empty.
func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string][]handler),
	}
}

// Handle associates "h" to the pair of HTTP method and path pattern.
func (s *ServeMux) Handle(meth string, pat Pattern, h HandlerFunc) {
	s.handlers[meth] = append(s.handlers[meth], handler{pat: pat, h: h})
}

// ServeHTTP dispatches the request to the first handler whose pattern matches to r.Method and r.Path.
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	components := strings.Split(path[1:], "/")
	l := len(components)
	var verb string
	if idx := strings.LastIndex(components[l-1], ":"); idx == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if idx > 0 {
		c := components[l-1]
		verb, components[l-1] = c[:idx], c[idx+1:]
	}

	for _, h := range s.handlers[r.Method] {
		pathParams, err := h.pat.Match(components, verb)
		if err != nil {
			glog.V(3).Infof("path mismatch: %q to %q", path, h.pat)
			continue
		}
		h.h(w, r, pathParams)
		return
	}

	// lookup other methods to determine if it is MethodNotAllowed
	for m, handlers := range s.handlers {
		if m == r.Method {
			continue
		}
		for _, h := range handlers {
			if _, err := h.pat.Match(components, verb); err == nil {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

type handler struct {
	pat Pattern
	h   HandlerFunc
}
