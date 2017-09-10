package runtime

import (
	"io"
	"reflect"

	"github.com/golang/protobuf/proto"
	hb "google.golang.org/genproto/googleapis/api/httpbody"
)

var (
	backupMarshaler = &JSONPb{OrigName: true}
)

// HttpBodyMarshaler is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is
// the actual message used as the response. If not, then this will
// simply fallback to the JSONPb marshaler.
type HttpBodyMarshaler struct{}

// ContentType returns the type specified in the google.api.HttpBody
// proto if "v" is a google.api.HttpBody proto, otherwise returns
// "application/json".
func (*HttpBodyMarshaler) ContentType(v interface{}) string {
	if h := tryHttpBody(v); h != nil {
		return h.GetContentType()
	}
	return "application/json"
}

// Marshal marshals "v" by returning the body bytes if v is a
// google.api.HttpBody message, or it marshals to JSON.
func (*HttpBodyMarshaler) Marshal(v interface{}) ([]byte, error) {
	if h := tryHttpBody(v); h != nil {
		return h.GetData(), nil
	}
	return backupMarshaler.Marshal(v)
}

// Unmarshal unmarshals JSON data into "v".
// google.api.HttpBody messages are not supported on the request.
func (*HttpBodyMarshaler) Unmarshal(data []byte, v interface{}) error {
	return backupMarshaler.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (*HttpBodyMarshaler) NewDecoder(r io.Reader) Decoder {
	return backupMarshaler.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (*HttpBodyMarshaler) NewEncoder(w io.Writer) Encoder {
	return backupMarshaler.NewEncoder(w)
}

func tryHttpBody(v interface{}) *hb.HttpBody {
	rv := reflect.ValueOf(v)
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
				return v.(*hb.HttpBody)
			}
		}
		rv = rv.Elem()
	}
	return nil
}
