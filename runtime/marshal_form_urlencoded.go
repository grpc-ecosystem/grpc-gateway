package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
	"io"
	"net/url"
	"reflect"
	"strconv"
)

// FormUrlencoded is a Marshaler which marshals/unmarshals into/from JSON
// with the "google.golang.org/protobuf/encoding/protojson" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a Decoder, so the underlying
// *json.Decoder methods can be used.
type FormUrlencoded JSONPb

// ContentType always returns "application/json".
func (*FormUrlencoded) ContentType(_ interface{}) string {
	return "application/json"
}

// Marshal marshals "v" into JSON.
func (f *FormUrlencoded) Marshal(v interface{}) ([]byte, error) {
	if _, ok := v.(proto.Message); !ok {
		return f.marshalNonProtoField(v)
	}

	var buf bytes.Buffer
	if err := f.marshalTo(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *FormUrlencoded) marshalTo(w io.Writer, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := f.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}
	b, err := f.MarshalOptions.Marshal(p)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

//var (
//	// protoMessageType is stored to prevent constant lookup of the same type at runtime.
//	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
//)

// marshalNonProto marshals a non-message field of a protobuf message.
// This function does not correctly marshal arbitrary data structures into JSON,
// it is only capable of marshaling non-message field values of protobuf,
// i.e. primitive types, enums; pointers to primitives or enums; maps from
// integer/string types to primitives/enums/pointers to messages.
func (f *FormUrlencoded) marshalNonProtoField(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []byte("null"), nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Slice {
		if rv.IsNil() {
			if f.EmitUnpopulated {
				return []byte("[]"), nil
			}
			return []byte("null"), nil
		}

		if rv.Type().Elem().Implements(protoMessageType) {
			var buf bytes.Buffer
			err := buf.WriteByte('[')
			if err != nil {
				return nil, err
			}
			for i := 0; i < rv.Len(); i++ {
				if i != 0 {
					err = buf.WriteByte(',')
					if err != nil {
						return nil, err
					}
				}
				if err = f.marshalTo(&buf, rv.Index(i).Interface().(proto.Message)); err != nil {
					return nil, err
				}
			}
			err = buf.WriteByte(']')
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}

		if rv.Type().Elem().Implements(typeProtoEnum) {
			var buf bytes.Buffer
			err := buf.WriteByte('[')
			if err != nil {
				return nil, err
			}
			for i := 0; i < rv.Len(); i++ {
				if i != 0 {
					err = buf.WriteByte(',')
					if err != nil {
						return nil, err
					}
				}
				if f.UseEnumNumbers {
					_, err = buf.WriteString(strconv.FormatInt(rv.Index(i).Int(), 10))
				} else {
					_, err = buf.WriteString("\"" + rv.Index(i).Interface().(protoEnum).String() + "\"")
				}
				if err != nil {
					return nil, err
				}
			}
			err = buf.WriteByte(']')
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	if rv.Kind() == reflect.Map {
		m := make(map[string]*json.RawMessage)
		for _, k := range rv.MapKeys() {
			buf, err := f.Marshal(rv.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}
			m[fmt.Sprintf("%v", k.Interface())] = (*json.RawMessage)(&buf)
		}
		if f.Indent != "" {
			return json.MarshalIndent(m, "", f.Indent)
		}
		return json.Marshal(m)
	}
	if enum, ok := rv.Interface().(protoEnum); ok && !f.UseEnumNumbers {
		return json.Marshal(enum.String())
	}
	return json.Marshal(rv.Interface())
}

// Unmarshal unmarshals JSON "data" into "v"
func (f *FormUrlencoded) Unmarshal(data []byte, v interface{}) error {
	return unmarshalJSONPb(data, f.UnmarshalOptions, v)
}

// Delimiter for newline encoded JSON streams.
func (f *FormUrlencoded) Delimiter() []byte {
	return []byte("")
}

type FormDecoder struct {
	reader io.Reader
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (f *FormUrlencoded) NewDecoder(r io.Reader) runtime.Decoder {
	return FormDecoder{reader: r}
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (f *FormUrlencoded) NewEncoder(w io.Writer) runtime.Encoder {
	return runtime.EncoderFunc(func(v interface{}) error {
		return f.marshalTo(w, v)
	})
}

func (d FormDecoder) Decode(v interface{}) error {
	formData, err := io.ReadAll(d.reader)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(formData))
	if err != nil {
		return err
	}

	switch val := v.(type) {
	case proto.Message:
		err = PopulateQueryParameters(val, values, &utilities.DoubleArray{})
	default:
		err = decodeSingleFormField(values, v)
	}

	return err
}

func decodeSingleFormField(values url.Values, v interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no form field found")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer type")
	}

	e := rv.Elem()
	for _, arr := range values {
		switch v.(type) {
		case *int8, *int16, *int32, *int64:
			i64, _ := strconv.ParseInt(arr[0], 10, 0)
			e.SetInt(i64)
		case *uint8, *uint16, *uint32, *uint64:
			ui64, _ := strconv.ParseUint(arr[0], 10, 0)
			e.SetUint(ui64)
		case *float32:
			f64, _ := strconv.ParseFloat(arr[0], 32)
			e.SetFloat(f64)
		case *float64:
			f64, _ := strconv.ParseFloat(arr[0], 64)
			e.SetFloat(f64)
		case *bool:
			b, _ := strconv.ParseBool(arr[0])
			e.SetBool(b)
		case *string:
			e.SetString(arr[0])
		case *[]byte:
			e.SetBytes([]byte(arr[0]))
		default:
			return fmt.Errorf("not valid type")
		}
	}

	return nil
}
