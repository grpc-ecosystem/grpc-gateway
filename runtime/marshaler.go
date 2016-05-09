package runtime

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const mimeWildcard = "*"

var (
	inboundDefaultMarshaler  = new(JSONBuiltin)
	outboundDefaultMarshaler = new(JSONBuiltin)

	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")
)

// MarshalerForRequest returns the inbound/outbound marshalers for this request.
// It checks the registry on the ServeMux for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If that isn't set, uses the ServerMux's InboundMarshaler/OutboundMarshaler.
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
func MarshalerForRequest(mux *ServeMux, r *http.Request) (inbound Marshaler, outbound Marshaler) {
	inbound = nil
	outbound = nil
	headerVals := r.Header[contentTypeHeader]

	for _, val := range headerVals {
		if mux.MIMERegistry != nil && mux.MIMERegistry.mimeMap != nil {
			if m, ok := mux.MIMERegistry.mimeMap[val]; ok {
				if inbound == nil && m.inbound != nil {
					inbound = m.inbound
				}
				if outbound == nil && m.outbound != nil {
					outbound = m.outbound
				}
			}
		} else {
			//Nil mimeMap, no need to bother checking for mimeWildcard
			if mux.InboundMarshaler == nil {
				//Its nil, use our default
				inbound = inboundDefaultMarshaler
			} else {
				inbound = mux.InboundMarshaler
			}

			if mux.OutboundMarshaler == nil {
				//Its nil, use our default
				outbound = outboundDefaultMarshaler
			} else {
				outbound = mux.OutboundMarshaler
			}
		}

		if inbound != nil && outbound != nil {
			//Got them both, return
			return inbound, outbound
		}
	}

	if mux.MIMERegistry != nil && mux.MIMERegistry.mimeMap != nil {
		if m, ok := mux.MIMERegistry.mimeMap[mimeWildcard]; ok {
			if inbound == nil && m.inbound != nil {
				inbound = m.inbound
			}
			if outbound == nil && m.outbound != nil {
				outbound = m.outbound
			}
		}
	}

	//Haven't gotten anywhere with any of the headers or mimeWildcard
	//Try to use the mux, otherwise use our default
	if inbound == nil {
		if mux.InboundMarshaler == nil {
			//Its nil, use our default
			inbound = inboundDefaultMarshaler
		} else {
			inbound = mux.InboundMarshaler
		}
	}

	if outbound == nil {
		if mux.OutboundMarshaler == nil {
			//Its nil, use our default
			outbound = outboundDefaultMarshaler
		} else {
			outbound = mux.OutboundMarshaler
		}
	}

	return inbound, outbound

}

// MarshalerMIMERegistry keeps a mapping from MIME types to mimeMarshalers.
type MarshalerMIMERegistry struct {
	mimeMap map[string]*mimeMarshaler
}

type mimeMarshaler struct {
	inbound  Marshaler
	outbound Marshaler
}

// AddMarshaler adds an inbound and outbund marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Inbound is the marshaler that is used when marshaling inbound requests from the client.
// Outbound is the marshaler that is used when marshaling outbound responses to the client.
func (m *MarshalerMIMERegistry) AddMarshaler(mime string, inbound, outbound Marshaler) {

	if len(mime) == 0 {
		panic("Mime can't be an empty string")
	}

	m.mimeMap[mime] = &mimeMarshaler{
		inbound:  inbound,
		outbound: outbound,
	}

}

// AddInboundMarshaler adds an inbound marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Inbound is the marshaler that is used when marshaling inbound requests from the client.
func (m *MarshalerMIMERegistry) AddInboundMarshaler(mime string, inbound Marshaler) {

	if len(mime) == 0 {
		panic("Mime can't be an empty string")
	}

	if _, ok := m.mimeMap[mime]; ok {
		//Already have this mime, just change inbound
		m.mimeMap[mime].inbound = inbound
	} else {
		m.mimeMap[mime] = &mimeMarshaler{
			inbound:  inbound,
			outbound: nil,
		}
	}

}

// AddOutboundMarshaler adds an outbund marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Outbound is the marshaler that is used when marshaling outbound responses to the client.
func (m *MarshalerMIMERegistry) AddOutboundMarshaler(mime string, outbound Marshaler) {
	mime = http.CanonicalHeaderKey(mime)
	if len(mime) == 0 {
		panic("Mime can't be an empty string")
	}

	if _, ok := m.mimeMap[mime]; ok {
		//Already have this mime, just change outbound
		m.mimeMap[mime].outbound = outbound
	} else {
		m.mimeMap[mime] = &mimeMarshaler{
			inbound:  nil,
			outbound: outbound,
		}
	}

}

// NewMarshalerMIMERegistry returns a new registry of marshalers.
// It allows for a mapping of case-sensitive Content-Type MIME type string to runtime.Marshaler interfaces.
//
// For example, you could allow the client to specify the use of the runtime.JSONPb marshaler
// with a "applicaton/jsonpb" Content-Type and the use of the runtime.JSONBuiltin marshaler
// with a "application/json" Content-Type.
// "*" can be used to match any Content-Type.
// This can be attached to a ServerMux with the MIMERegistry option.
func NewMarshalerMIMERegistry() *MarshalerMIMERegistry {
	return &MarshalerMIMERegistry{
		mimeMap: make(map[string]*mimeMarshaler),
	}
}

// Marshaler defines a conversion between byte sequence and gRPC payloads / fields.
type Marshaler interface {
	// Marshal marshals "v" into byte sequence.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal unmarshals "data" into "v".
	// "v" must be a pointer value.
	Unmarshal(data []byte, v interface{}) error
	// NewDecoder returns a Decoder which reads byte sequence from "r".
	NewDecoder(r io.Reader) Decoder
	// NewEncoder returns an Encoder which writes bytes sequence into "w".
	NewEncoder(w io.Writer) Encoder
	// ContentType returns the Content-Type which this marshaler is responsible for.
	ContentType() string
}

// Decoder decodes a byte sequence
type Decoder interface {
	Decode(v interface{}) error
}

// Encoder encodes gRPC payloads / fields into byte sequence.
type Encoder interface {
	Encode(v interface{}) error
}

// JSONBuiltin is a Marshaler which marshals/unmarshals into/from JSON
// with the standard "encoding/json" package of Golang.
// Although it is generally faster for simple proto messages than JSONPb,
// it does not support advanced features of protobuf, e.g. map, oneof, ....
type JSONBuiltin struct{}

// ContentType always Returns "application/json".
func (*JSONBuiltin) ContentType() string {
	return "application/json"
}

// Marshal marshals "v" into JSON
func (j *JSONBuiltin) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals JSON data into "v".
func (j *JSONBuiltin) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSONBuiltin) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONBuiltin) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

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
