package runtime

import (
	"errors"
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
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
func MarshalerForRequest(mux *ServeMux, r *http.Request) (inbound Marshaler, outbound Marshaler) {
	headerVals := append(append([]string(nil), r.Header[contentTypeHeader]...), "*")

	for _, val := range headerVals {
		m := mux.marshalers.lookup(val)
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
		inbound = inboundDefaultMarshaler
	}
	if outbound == nil {
		outbound = outboundDefaultMarshaler
	}
	return inbound, outbound

}

// marshalerRegistry keeps a mapping from MIME types to mimeMarshalers.
type marshalerRegistry map[string]*mimeMarshaler

type mimeMarshaler struct {
	inbound  Marshaler
	outbound Marshaler
}

// addMarshaler adds an inbound and outbund marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Inbound is the marshaler that is used when marshaling inbound requests from the client.
// Outbound is the marshaler that is used when marshaling outbound responses to the client.
func (r *marshalerRegistry) add(mime string, inbound, outbound Marshaler) error {
	if mime == "" {
		return errors.New("empty MIME type")
	}
	(*r)[mime] = &mimeMarshaler{
		inbound:  inbound,
		outbound: outbound,
	}
	return nil
}

// addInboundMarshaler adds an inbound marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Inbound is the marshaler that is used when marshaling inbound requests from the client.
func (r *marshalerRegistry) addInbound(mime string, inbound Marshaler) error {
	if mime == "" {
		return errors.New("empty MIME type")
	}
	if entry := (*r)[mime]; entry != nil {
		entry.inbound = inbound
		return nil
	}
	(*r)[mime] = &mimeMarshaler{inbound: inbound}
	return nil
}

// addOutBound adds an outbund marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
// Outbound is the marshaler that is used when marshaling outbound responses to the client.
func (r *marshalerRegistry) addOutbound(mime string, outbound Marshaler) error {
	mime = http.CanonicalHeaderKey(mime)
	if mime == "" {
		return errors.New("empty MIME type")
	}
	if entry := (*r)[mime]; entry != nil {
		entry.outbound = outbound
		return nil
	}
	(*r)[mime] = &mimeMarshaler{outbound: outbound}
	return nil

}

func (r *marshalerRegistry) lookup(mime string) *mimeMarshaler {
	if r == nil {
		return nil
	}
	return (*r)[mime]
}

// WithMarshalerOption returns a ServeMuxOption which associates inbound and outbound
// Marshalers to a MIME type in mux.
func WithMarshalerOption(mime string, in, out Marshaler) ServeMuxOption {
	return func(mux *ServeMux) {
		if err := mux.marshalers.add(mime, in, out); err != nil {
			panic(err)
		}
	}
}

// WithInboundMarshalerOption returns a ServeMuxOption which associates an inbound
// Marshaler to a MIME type in mux.
func WithInboundMarshalerOption(mime string, in Marshaler) ServeMuxOption {
	return func(mux *ServeMux) {
		if err := mux.marshalers.addInbound(mime, in); err != nil {
			panic(err)
		}
	}
}

// WithOutboundMarshalerOption returns a ServeMuxOption which associates an outbound
// Marshaler to a MIME type in mux.
func WithOutboundMarshalerOption(mime string, out Marshaler) ServeMuxOption {
	return func(mux *ServeMux) {
		if err := mux.marshalers.addOutbound(mime, out); err != nil {
			panic(err)
		}
	}
}
