package runtime

import (
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"io"
)

var (
	jsonPbMarshaller = &JSONPb{OrigName: true}
)

// SetHTTPBodyMarshaler overwrite  the default marshaler with the HTTPBodyMarshaler
func SetHTTPBodyMarshaler(serveMux *ServeMux) {
	serveMux.marshalers.mimeMap[MIMEWildcard] = &HTTPBodyMarshaler{OrigName: true}
}

// HTTPBodyMarshaler is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is
// the actual message used as the response. If not, then this will
// simply fallback to the JSONPb marshaler.
type HTTPBodyMarshaler jsonpb.Marshaler

// ContentType in case v is a google.api.HttpBody message it returns
// its specified content type otherwise fall back to the JSONPb marshaler.
func (*HTTPBodyMarshaler) ContentType(v interface{}) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return jsonPbMarshaller.ContentType(v)
}

// Marshal marshals "v" by returning the body bytes if v is a
// google.api.HttpBody message, otherwise it falls back to the JSONPb marshaler.
func (*HTTPBodyMarshaler) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.Data, nil
	}
	return jsonPbMarshaller.Marshal(v)
}

// Unmarshal unmarshals JSON data into "v".
// google.api.HttpBody messages are not supported for request messages.
func (*HTTPBodyMarshaler) Unmarshal(data []byte, v interface{}) error {
	return jsonPbMarshaller.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (*HTTPBodyMarshaler) NewDecoder(r io.Reader) Decoder {
	return jsonPbMarshaller.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (*HTTPBodyMarshaler) NewEncoder(w io.Writer) Encoder {
	return jsonPbMarshaller.NewEncoder(w)
}
