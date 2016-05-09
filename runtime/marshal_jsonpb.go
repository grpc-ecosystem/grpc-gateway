package runtime

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// JSONPb is a Marshaler which marshals/unmarshals into/from JSON
// with the "github.com/golang/protobuf/jsonpb".
// It supports fully functionality of protobuf unlike JSONBuiltin.
type JSONPb struct {
	// Whether to render enum values as integers, as opposed to string values.
	EnumsAsInts bool

	// Whether to render fields with zero values.
	EmitDefaults bool

	// A string to indent each level by. The presence of this field will
	// also cause a space to appear between the field separator and
	// value, and for newlines to be appear between fields and array
	// elements.
	Indent string

	// Whether to use the original (.proto) name for fields.
	OrigName bool
}

// ContentType always returns "application/json".
func (*JSONPb) ContentType() string {
	return "application/json"
}

// Marshal marshals "v" into JSON
// Currently it can marshal only proto.Message.
// TODO(yugui) Support fields of primitive types in a message.
func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	m := &jsonpb.Marshaler{
		EnumsAsInts:  j.EnumsAsInts,
		EmitDefaults: j.EmitDefaults,
		Indent:       j.Indent,
		OrigName:     j.OrigName,
	}
	p, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("interface is not proto.Message")
	}

	var buf bytes.Buffer
	if err := m.Marshal(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

// Unmarshal unmarshals JSON "data" into "v"
// Currently it can marshal only proto.Message.
// TODO(yugui) Support fields of primitive types in a message.
func (j *JSONPb) Unmarshal(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	p, ok := v.(proto.Message)
	if !ok {
		return errors.New("interface is not proto.Message")
	}
	return jsonpb.Unmarshal(r, p)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSONPb) NewDecoder(r io.Reader) Decoder {
	return &jsonPbDecoder{decoder: json.NewDecoder(r)}
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPb) NewEncoder(w io.Writer) Encoder {
	m := &jsonpb.Marshaler{
		EnumsAsInts:  j.EnumsAsInts,
		EmitDefaults: j.EmitDefaults,
		Indent:       j.Indent,
		OrigName:     j.OrigName,
	}

	return &jsonPbEncoder{
		marshal: m,
		writer:  w,
	}
}

type jsonPbDecoder struct {
	decoder *json.Decoder
}

func (j *jsonPbDecoder) Decode(v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		return errors.New("interface is not proto.Message")
	}

	return jsonpb.UnmarshalNext(j.decoder, p)
}

type jsonPbEncoder struct {
	marshal *jsonpb.Marshaler
	writer  io.Writer
}

func (j *jsonPbEncoder) Encode(v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		return errors.New("interface is not proto.Message")
	}
	return j.marshal.Marshal(j.writer, p)
}
