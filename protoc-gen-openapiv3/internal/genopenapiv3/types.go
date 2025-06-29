package genopenapiv3

import (
	"errors"

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


func (f Format) MarshalOpenAPIDoc(doc *openapi3.T) ([]byte, error) {
	switch f {
	case FormatJSON:
		return doc.MarshalJSON()
	case FormatYAML:
		openapiDoc, err := doc.MarshalYAML()
		if err != nil {
			return nil, err
		}
		return yaml.Marshal(openapiDoc)
	default:
		grpclog.Errorf("unsupported format: %s\n", f)
		return nil, ErrUnSupportedFormat
	}
}

type OneOfStrategy string

const (
	OneOfStrategyOneOf = "oneof"
	OneOfStrategyAllOf = "allof"
)
