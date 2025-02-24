package runtime

import (
	"fmt"
	"io"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
)

type UrlEncodedDecoder struct {
	r io.Reader
}

func NewUrlEncodedDecoder(r io.Reader) Decoder {
	return &UrlEncodedDecoder{r: r}
}

func (u *UrlEncodedDecoder) Decode(v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("not proto message")
	}

	formData, err := io.ReadAll(u.r)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(formData))
	if err != nil {
		return err
	}

	filter := &utilities.DoubleArray{}

	err = PopulateQueryParameters(msg, values, filter)
	if err != nil {
		return err
	}

	return nil
}

type UrlEncodeMarshal struct {
	Marshaler
}

// ContentType means the content type of the response
func (u *UrlEncodeMarshal) ContentType(_ interface{}) string {
	return "application/json"
}

func (u *UrlEncodeMarshal) Marshal(v interface{}) ([]byte, error) {
	return u.Marshaler.Marshal(v)
}

// NewDecoder indicates how to decode the request
func (u UrlEncodeMarshal) NewDecoder(r io.Reader) Decoder {
	return NewUrlEncodedDecoder(r)
}
