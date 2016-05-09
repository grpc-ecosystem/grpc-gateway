package runtime

import (
	"io"
	"net/http"
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
	headerVals := append(append([]string(nil), r.Header[contentTypeHeader]...), "*")

	for _, val := range headerVals {
		m := mux.MIMERegistry.lookup(val)
		if m != nil {
			if inbound == nil {
				inbound = m.inbound
			}
			if outbound == nil {
				outbound = m.outbound
			}
		}
		if inbound != nil && outbound != nil {
			// Got them both, return
			return inbound, outbound
		}
	}

	if inbound == nil {
		inbound = mux.InboundMarshaler
	}
	if inbound == nil {
		inbound = inboundDefaultMarshaler
	}

	if outbound == nil {
		outbound = mux.OutboundMarshaler
	}
	if outbound == nil {
		outbound = outboundDefaultMarshaler
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

func (m *MarshalerMIMERegistry) lookup(mime string) *mimeMarshaler {
	if m == nil {
		return nil
	}
	return m.mimeMap[mime]
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
