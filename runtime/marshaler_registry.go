package runtime

import (
	"errors"
	"net/http"
)

// MIMEWildcard is the fallback MIME type used for requests which do not match
// a registered MIME type.
const MIMEWildcard = "*"

var (
	acceptHeader      = http.CanonicalHeaderKey("Accept")
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")

	defaultMarshaler = &JSONPb{OrigName: true}
)

// MarshalerForRequest returns the outbound marshalers for this request.
// It checks the registry on the ServeMux for the MIME type set by the Accept header.
// If it isn't set (or the request Accept is empty), it checks for the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/OutboundMarshaler.
func MarshalerForRequest(mux *ServeMux, r *http.Request) (outbound Marshaler) {
	for _, acceptVal := range r.Header[acceptHeader] {
		if m, ok := mux.marshalers.mimeMap[acceptVal]; ok {
			outbound = m
			break
		}
	}

	if outbound == nil {
		for _, contentTypeVal := range r.Header[contentTypeHeader] {
			if m, ok := mux.marshalers.mimeMap[contentTypeVal]; ok {
				outbound = m
				break
			}
		}
	}
	if outbound == nil {
		outbound = mux.marshalers.mimeMap[MIMEWildcard]
	}

	return outbound
}

// UnmarshalerForRequest returns the inbound marshalers for this request.
// It checks the registry on the ServeMux for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
func UnmarshalerForRequest(mux *ServeMux, r *http.Request) (inbound Unmarshaler) {
	for _, contentTypeVal := range r.Header[contentTypeHeader] {
		if m, ok := mux.unmarshalers.mimeMap[contentTypeVal]; ok {
			inbound = m
			break
		}
	}

	if inbound == nil {
		inbound = mux.unmarshalers.mimeMap[MIMEWildcard]
	}

	return inbound
}

// marshalerRegistry is a mapping from MIME types to Marshalers.
type marshalerRegistry struct {
	mimeMap map[string]Marshaler
}

// unmarshalerRegistry is a mapping from MIME types to Unmarshalers.
type unmarshalerRegistry struct {
	mimeMap map[string]Unmarshaler
}

// add adds a marshaler for a case-sensitive MIME type string ("*" to match any
// MIME type).
func (m marshalerRegistry) add(mime string, marshaler Marshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}

	m.mimeMap[mime] = marshaler

	return nil
}

// add adds a marshaler for a case-sensitive MIME type string ("*" to match any
// MIME type).
func (m unmarshalerRegistry) add(mime string, unmarshaler Unmarshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}

	m.mimeMap[mime] = unmarshaler

	return nil
}

// makeMarshalerMIMERegistry returns a new registry of marshalers.
// It allows for a mapping of case-sensitive Content-Type MIME type string to runtime.Marshaler interfaces.
//
// For example, you could allow the client to specify the use of the runtime.JSONPb marshaler
// with a "application/jsonpb" Content-Type and the use of the runtime.JSONBuiltin marshaler
// with a "application/json" Content-Type.
// "*" can be used to match any Content-Type.
// This can be attached to a ServerMux with the marshaler option.
func makeMarshalerMIMERegistry() marshalerRegistry {
	return marshalerRegistry{
		mimeMap: map[string]Marshaler{
			MIMEWildcard: defaultMarshaler,
		},
	}
}

func makeUnmarshalerMIMERegistry() unmarshalerRegistry {
	return unmarshalerRegistry{
		mimeMap: map[string]Unmarshaler{
			MIMEWildcard: defaultMarshaler,
		},
	}
}

// WithMarshalerOption returns a ServeMuxOption which associates outbound
// Marshalers to a MIME type in mux.
func WithMarshalerOption(mime string, marshaler Marshaler) ServeMuxOption {
	return func(mux *ServeMux) {
		if err := mux.marshalers.add(mime, marshaler); err != nil {
			panic(err)
		}
	}
}

// WithUnmarshalerOption returns a ServeMuxOption which associates inbound
// Marshalers to a MIME type in mux.
func WithUnmarshalerOption(mime string, unmarshaler Unmarshaler) ServeMuxOption {
	return func(mux *ServeMux) {
		if err := mux.unmarshalers.add(mime, unmarshaler); err != nil {
			panic(err)
		}
	}
}
