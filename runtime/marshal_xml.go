package runtime

import (
	"bytes"
	"encoding/xml"
	"io"
	"reflect"
)

// XMLBuiltin is a Marshaler which marshals/unmarshals into/from XML
// with the standard "encoding/xml" package of Golang.
type XMLBuiltin struct{}

// ContentType always Returns "application/xml".
func (*XMLBuiltin) ContentType() string {
	return "application/xml"
}

// Marshal marshals "v" into XML
func (j *XMLBuiltin) Marshal(v interface{}) ([]byte, error) {
	var rootname string

	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		rootname = t.Elem().Name()
	} else {
		rootname = t.Name()
	}

	data, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}

	data = bytes.Replace(data, []byte(rootname), []byte("xml"), -1)
	return data, nil
}

// Unmarshal unmarshals XML data into "v".
func (j *XMLBuiltin) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads XML stream from "r".
func (j *XMLBuiltin) NewDecoder(r io.Reader) Decoder {
	return xml.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes XML stream into "w".
func (j *XMLBuiltin) NewEncoder(w io.Writer) Encoder {
	return xml.NewEncoder(w)
}
