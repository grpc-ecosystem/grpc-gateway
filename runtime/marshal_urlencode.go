package runtime

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
)

type UrlEncodeMarshal struct {
	Marshaler
}

// ContentType means the content type of the response
func (u UrlEncodeMarshal) ContentType(_ interface{}) string {
	return "application/json"
}

func (u UrlEncodeMarshal) Marshal(v interface{}) ([]byte, error) {
	// can marshal the response in proto message format
	j := JSONPb{}
	return j.Marshal(v)
}

// NewDecoder indicates how to decode the request
func (u UrlEncodeMarshal) NewDecoder(r io.Reader) Decoder {
	return DecoderFunc(func(p interface{}) error {
		msg, ok := p.(proto.Message)
		if !ok {
			return fmt.Errorf("not proto message")
		}

		formData, err := ioutil.ReadAll(r)
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
	})
}
