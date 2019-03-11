package runtime

import (
	"google.golang.org/genproto/googleapis/api/httpbody"
	"io"
)

// SetHTTPBodyMarshaler overwrite  the default marshaler with the HTTPBodyMarshaler
func SetHTTPBodyMarshaler(serveMux *ServeMux) {
	serveMux.marshalers.mimeMap[MIMEWildcard] = &HTTPBodyMarshaler{
		DefaultMarshaler: &JSONPb{OrigName: true},
	}
}

// HTTPBodyMarshaler is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is
// the actual message used as the response. If not, then this will
// simply fallback to the JSONPb marshaler.
type HTTPBodyMarshaler struct {
	DefaultMarshaler Marshaler
}

// ContentType in case v is a google.api.HttpBody message it returns
// its specified content type otherwise fall back to the JSONPb marshaler.
func (h *HTTPBodyMarshaler) ContentType(v interface{}) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.DefaultMarshaler.ContentType(v)
}

// Marshal marshals "v" by returning the body bytes if v is a
// google.api.HttpBody message, otherwise it falls back to the JSONPb marshaler.
func (h *HTTPBodyMarshaler) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.Data, nil
	}
	return h.DefaultMarshaler.Marshal(v)
}

// Unmarshal unmarshals JSON data into "v".
// google.api.HttpBody messages are not supported for request messages.
func (h *HTTPBodyMarshaler) Unmarshal(data []byte, v interface{}) error {
	return h.DefaultMarshaler.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (h *HTTPBodyMarshaler) NewDecoder(r io.Reader) Decoder {
	return h.DefaultMarshaler.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (h *HTTPBodyMarshaler) NewEncoder(w io.Writer) Encoder {
	return h.DefaultMarshaler.NewEncoder(w)
}
