package runtime

import (
	"context"
	"io"
)

// Marshaler defines an outbound conversion between a gRPC payload and a byte sequence.
type Marshaler interface {
	// Marshal marshals "v" into a byte sequence.
	Marshal(ctx context.Context, v interface{}) ([]byte, error)
	// NewEncoder returns an Encoder which writes byte sequences into "w".
	NewEncoder(ctx context.Context, w io.Writer) Encoder
	// ContentType returns the Content-Type which this marshaler is responsible for.
	ContentType() string
}

// Unmarshaler defines an inbound conversion between a byte sequence and a gRPC payload.
type Unmarshaler interface {
	// Unmarshal unmarshals "data" into "v".
	// "v" must be a pointer value.
	Unmarshal(ctx context.Context, data []byte, v interface{}) error
	// NewDecoder returns a Decoder which reads byte sequences from "r".
	NewDecoder(ctx context.Context, r io.Reader) Decoder
	// ContentType returns the Content-Type which this marshaler is responsible for.
	ContentType() string
}

// Marshalers that implement contentTypeMarshaler will have their ContentTypeFromMessage method called
// to set the Content-Type header on the response
type contentTypeMarshaler interface {
	// ContentTypeFromMessage returns the Content-Type this marshaler produces from the provided message
	ContentTypeFromMessage(v interface{}) string
}

// Decoder decodes a byte sequence
type Decoder interface {
	Decode(v interface{}) error
}

// Encoder encodes gRPC payloads / fields into byte sequence.
type Encoder interface {
	Encode(v interface{}) error
}

// DecoderFunc adapts an decoder function into Decoder.
type DecoderFunc func(v interface{}) error

// Decode delegates invocations to the underlying function itself.
func (f DecoderFunc) Decode(v interface{}) error { return f(v) }

// EncoderFunc adapts an encoder function into Encoder
type EncoderFunc func(v interface{}) error

// Encode delegates invocations to the underlying function itself.
func (f EncoderFunc) Encode(v interface{}) error { return f(v) }

// Delimited defines the streaming delimiter.
type Delimited interface {
	// Delimiter returns the record seperator for the stream.
	Delimiter() []byte
}
