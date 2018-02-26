package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
)

// This Marshal is to parse application/x-www-form-urlencoded and return json
// You should add `runtime.WithMarshalerOption("application/x-www-form-urlencoded", &runtime.FORMPb{}),`
// before MIMEWildcard

// FORMPb is a Marshaler which marshals into JSON
// with the "github.com/golang/protobuf/jsonpb".
// It supports fully functionality of protobuf unlike JSONBuiltin.
type FORMPb jsonpb.Marshaler

// NewDecoder returns a Decoder which reads form data stream from "r".
func (j *FORMPb) NewDecoder(r io.Reader) Decoder {
	return DecoderFunc(func(v interface{}) error { return decodeFORMPb(r, v) })
}

func decodeFORMPb(d io.Reader, v interface{}) error {
	msg, ok := v.(proto.Message)

	if !ok {
		return fmt.Errorf("not proto message")
	}

	formData, err := ioutil.ReadAll(d)

	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(formData))

	if err != nil {
		return err
	}

	filter := &utilities.DoubleArray{}

	err = PopulateQueryParameters(msg, values, filter)

	if err != nil {
		return err
	}

	return nil
}

// Unmarshal not realized for the moment
func (j *FORMPb) Unmarshal(data []byte, v interface{}) error {
	return nil
}

// ContentType always returns "application/json".
func (*FORMPb) ContentType() string {
	return "application/json"
}

// Marshal marshals "v" into JSON
// Currently it can marshal only proto.Message.
// TODO(yugui) Support fields of primitive types in a message.
func (j *FORMPb) Marshal(v interface{}) ([]byte, error) {
	if _, ok := v.(proto.Message); !ok {
		return j.marshalNonProtoField(v)
	}

	var buf bytes.Buffer
	if err := j.marshalTo(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *FORMPb) marshalTo(w io.Writer, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := j.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}
	return (*jsonpb.Marshaler)(j).Marshal(w, p)
}

func (j *FORMPb) marshalNonProtoField(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []byte("null"), nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Map {
		m := make(map[string]*json.RawMessage)
		for _, k := range rv.MapKeys() {
			buf, err := j.Marshal(rv.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}
			m[fmt.Sprintf("%v", k.Interface())] = (*json.RawMessage)(&buf)
		}
		if j.Indent != "" {
			return json.MarshalIndent(m, "", j.Indent)
		}
		return json.Marshal(m)
	}
	if enum, ok := rv.Interface().(protoEnum); ok && !j.EnumsAsInts {
		return json.Marshal(enum.String())
	}
	return json.Marshal(rv.Interface())
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *FORMPb) NewEncoder(w io.Writer) Encoder {
	return EncoderFunc(func(v interface{}) error { return j.marshalTo(w, v) })
}
