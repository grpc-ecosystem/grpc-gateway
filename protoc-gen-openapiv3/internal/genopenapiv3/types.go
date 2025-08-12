package genopenapiv3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/yaml.v3"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

var ErrUnSupportedFormat = errors.New("unsupported format provided")

type Encoder interface {
	Encode(v any) error
}

func (f Format) MarshalOpenAPIDoc(doc *openapi3.T) ([]byte, error) {

	openapiDoc, err := doc.MarshalYAML()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	var encoder Encoder

	switch f {
	case FormatJSON:
		encoder = json.NewEncoder(&buf)
		encoder.(*json.Encoder).SetIndent("", "  ")
	case FormatYAML:
		encoder = yaml.NewEncoder(&buf)
		encoder.(*yaml.Encoder).SetIndent(2)
	default:
		grpclog.Errorf("unsupported format: %s\n", f)
		return nil, ErrUnSupportedFormat
	}

	err = encoder.Encode(openapiDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode OpenAPI document: %w", err)
	}

	return buf.Bytes(), nil
}

type OneOfStrategy string

const (
	OneOfStrategyOneOf = "oneof"
	OneOfStrategyAllOf = "allof"
)
