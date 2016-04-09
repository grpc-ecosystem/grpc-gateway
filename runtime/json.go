package runtime

import (
  "errors"
  "io"
  "bytes"
  "encoding/json"
  
  "github.com/golang/protobuf/jsonpb"
  "github.com/golang/protobuf/proto"
)

type JSONAdapter interface {
  Marshal(v interface{}) ([]byte, error)
  Unmarshal(data []byte, v interface{}) error
  NewDecoder(r io.Reader) JSONDecoderAdapter
  NewEncoder(w io.Writer) JSONEncoderAdapter
}

type JSONDecoderAdapter interface {
  Decode(v interface{}) error
}

type JSONEncoderAdapter interface {
  Encode(v interface{}) error
}

type JSONBuiltin struct { }

func (j *JSONBuiltin) Marshal(v interface{}) ([]byte, error) {
  return json.Marshal(v)
}

func (j *JSONBuiltin) Unmarshal(data []byte, v interface{}) error {
  return json.Unmarshal(data,v)
}

func (j *JSONBuiltin) NewDecoder(r io.Reader) JSONDecoderAdapter {
  return json.NewDecoder(r)
}

func (j *JSONBuiltin) NewEncoder(w io.Writer) JSONEncoderAdapter {
  return json.NewEncoder(w)
}


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

func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
  m := &jsonpb.Marshaler{
    EnumsAsInts:j.EnumsAsInts,
    EmitDefaults:j.EmitDefaults,
    Indent:j.Indent,
    OrigName:j.OrigName,
  }
  if p, ok := v.(proto.Message); ok {
    var buf bytes.Buffer
  	if err := m.Marshal(&buf, p); err != nil {
  		return nil, err
  	}
  	return buf.Bytes(), nil
  } else {
    _ = v.(proto.Message)
    return nil,errors.New("Interface is not proto.Message")
  }
  
}

func (j *JSONPb) Unmarshal(data []byte, v interface{}) error {
  r := bytes.NewReader(data)
  if p, ok := v.(proto.Message); ok {
    return jsonpb.Unmarshal(r,p)
  } else {
    _ = v.(proto.Message)
    return errors.New("Interface is not proto.Message")
  }
}

func (j *JSONPb) NewDecoder(r io.Reader) JSONDecoderAdapter {
  return &JSONPbDecoder{decoder:json.NewDecoder(r)}
}

func (j *JSONPb) NewEncoder(w io.Writer) JSONEncoderAdapter {
  m := &jsonpb.Marshaler{
    EnumsAsInts:j.EnumsAsInts,
    EmitDefaults:j.EmitDefaults,
    Indent:j.Indent,
    OrigName:j.OrigName,
  }
  
  
  return &JSONPbEncoder{
    marshal:m,
    writer:w,
  }
}



type JSONPbDecoder struct {
  decoder *json.Decoder
}

func (j *JSONPbDecoder) Decode(v interface{}) error {
  if p, ok := v.(proto.Message); ok {
    return jsonpb.UnmarshalNext(j.decoder,p)
  } else {
    _ = v.(proto.Message)
    return errors.New("Interface is not proto.Message")
  }
}

type JSONPbEncoder struct {
  marshal *jsonpb.Marshaler
  writer io.Writer
}

func (j *JSONPbEncoder) Encode(v interface{}) error {
  if p, ok := v.(proto.Message); ok {
    return j.marshal.Marshal(j.writer,p)
  } else {
    _ = v.(proto.Message)
    return errors.New("Interface is not proto.Message")
  }
}