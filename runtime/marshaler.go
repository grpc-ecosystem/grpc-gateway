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

const MIMEWILDCARD = "*"

var inboundDefaultMarshaler = &JSONBuiltin{}
var outboundDefaultMarshaler = &JSONBuiltin{}
var contentTypeHeader = http.CanonicalHeaderKey("Content-Type")

// Get the inbound/outbound marshalers for this request. Checks the registry on the ServeMux for
// the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If that isn't set, uses the ServerMux's InboundMarshaler/OutboundMarshaler.
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry. Otherwise, follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
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
			//Nil mimeMap, no need to bother checking for MIMEWILDCARD
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
		if m, ok := mux.MIMERegistry.mimeMap[MIMEWILDCARD]; ok {
			if inbound == nil && m.inbound != nil {
				inbound = m.inbound
			}
			if outbound == nil && m.outbound != nil {
				outbound = m.outbound
			}
		}
	}

	//Haven't gotten anywhere with any of the headers or MIMEWILDCARD
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

type MarshalerMIMERegistry struct {
	mimeMap map[string]*mimeMarshaler
}

type mimeMarshaler struct {
	inbound  Marshaler
	outbound Marshaler
}

// Add an inbound and outbund marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Inbound is the marshaler that is used when marshaling inbound requests from the client.
// Outbound is the marshaler that is used when marshaling outbound responses to the client.
func (m *MarshalerMIMERegistry) AddMarshaler(mime string, inbound Marshaler, outbound Marshaler) {

	if len(mime) == 0 {
		panic("Mime can't be an empty string")
	}

	m.mimeMap[mime] = &mimeMarshaler{
		inbound:  inbound,
		outbound: outbound,
	}

}

// Add an inbound marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
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

// Add an outbund marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
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

// Create a new MIME marshaler registry. Allows for a mapping of case-sensitive
// Content-Type MIME type string to runtime.Marshaler interfaces.
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

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	NewDecoder(r io.Reader) Decoder
	NewEncoder(w io.Writer) Encoder
	ContentType() string
}

type Decoder interface {
	Decode(v interface{}) error
}

type Encoder interface {
	Encode(v interface{}) error
}

type JSONBuiltin struct{}

func (*JSONBuiltin) ContentType() string {
	return "application/json"
}

func (j *JSONBuiltin) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSONBuiltin) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j *JSONBuiltin) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func (j *JSONBuiltin) NewEncoder(w io.Writer) Encoder {
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

func (*JSONPb) ContentType() string {
	return "application/json"
}

func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	m := &jsonpb.Marshaler{
		EnumsAsInts:  j.EnumsAsInts,
		EmitDefaults: j.EmitDefaults,
		Indent:       j.Indent,
		OrigName:     j.OrigName,
	}
	if p, ok := v.(proto.Message); ok {
		var buf bytes.Buffer
		if err := m.Marshal(&buf, p); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	} else {
		return nil, errors.New("Interface is not proto.Message")
	}

}

func (j *JSONPb) Unmarshal(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	if p, ok := v.(proto.Message); ok {
		return jsonpb.Unmarshal(r, p)
	} else {
		return errors.New("Interface is not proto.Message")
	}
}

func (j *JSONPb) NewDecoder(r io.Reader) Decoder {
	return &JSONPbDecoder{decoder: json.NewDecoder(r)}
}

func (j *JSONPb) NewEncoder(w io.Writer) Encoder {
	m := &jsonpb.Marshaler{
		EnumsAsInts:  j.EnumsAsInts,
		EmitDefaults: j.EmitDefaults,
		Indent:       j.Indent,
		OrigName:     j.OrigName,
	}

	return &JSONPbEncoder{
		marshal: m,
		writer:  w,
	}
}

type JSONPbDecoder struct {
	decoder *json.Decoder
}

func (j *JSONPbDecoder) Decode(v interface{}) error {
	if p, ok := v.(proto.Message); ok {
		return jsonpb.UnmarshalNext(j.decoder, p)
	} else {
		return errors.New("Interface is not proto.Message")
	}
}

type JSONPbEncoder struct {
	marshal *jsonpb.Marshaler
	writer  io.Writer
}

func (j *JSONPbEncoder) Encode(v interface{}) error {
	if p, ok := v.(proto.Message); ok {
		return j.marshal.Marshal(j.writer, p)
	} else {
		return errors.New("Interface is not proto.Message")
	}
}
