package runtime

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

// SetHTTPBodyMarshaler overwrite the default marshaler with the HTTPBodyMarshaler
func SetHTTPBodyMarshaler(serveMux *ServeMux) {
	serveMux.marshalers.mimeMap[MIMEWildcard] = &HTTPBodyMarshaler{
		Marshaler: &JSONPb{OrigName: true},
	}
}

// HTTPBodyMarshaler is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is the actual
// message used as the response. Stream google.api.HttpBody message is also
// supported. If not, then this will simply fallback to the Marshaler specified
// as its default Marshaler.
type HTTPBodyMarshaler struct {
	Marshaler
}

// ContentType implementation to keep backwards compatability with marshal interface
func (h *HTTPBodyMarshaler) ContentType() string {
	return h.ContentTypeFromMessage(nil)
}

// ContentTypeFromMessage in case v is a google.api.HttpBody message it returns
// its specified content type otherwise fall back to the default Marshaler.
func (h *HTTPBodyMarshaler) ContentTypeFromMessage(v interface{}) string {
	if httpBody := tryHttpBody(v); httpBody != nil {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType()
}

// Marshal marshals "v" by returning the body bytes if v is a
// google.api.HttpBody message, otherwise it falls back to the default Marshaler.
func (h *HTTPBodyMarshaler) Marshal(v interface{}) ([]byte, error) {
	if httpBody := tryHttpBody(v); httpBody != nil {
		return httpBody.GetData(), nil
	}
	return h.Marshaler.Marshal(v)
}

// Delimiter for encoded multi-part streams.
func (h *HTTPBodyMarshaler) Delimiter() []byte {
	return []byte("")
}

func tryHttpBody(v interface{}) *httpbody.HttpBody {
	rv := reflect.ValueOf(v)
	// The handler wraps streamed chunks in a map.
	// If we're sending an HTTP body as a chunk, we need to unpack it.
	if rv.Kind() == reflect.Map && rv.Type().ConvertibleTo(reflect.TypeOf(map[string]proto.Message{})) {
		m := v.(map[string]proto.Message)
		if r, ok := m["result"]; ok && proto.MessageName(r) == "google.api.HttpBody" {
			return r.(*httpbody.HttpBody)
		}
		return nil
	}
	if rv.Kind() != reflect.Ptr {
		return nil
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		if rv.Type().ConvertibleTo(typeProtoMessage) {
			pb := rv.Interface().(proto.Message)
			if proto.MessageName(pb) == "google.api.HttpBody" {
				return v.(*httpbody.HttpBody)
			}
		}
		rv = rv.Elem()
	}
	return nil

}
