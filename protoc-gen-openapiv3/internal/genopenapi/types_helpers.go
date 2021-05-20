package genopenapi

import (
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

func marshalRef(value string, otherwise interface{}) ([]byte, error) {
	if len(value) > 0 {
		return json.Marshal(&refProps{
			Ref: value,
		})
	}
	return json.Marshal(otherwise)
}

func unmarshalRef(data []byte, destRef *string, destOtherwise interface{}) error {
	refProps := &refProps{}
	if err := json.Unmarshal(data, refProps); err == nil {
		ref := refProps.Ref
		if len(ref) > 0 {
			*destRef = ref
			return nil
		}
	}
	return json.Unmarshal(data, destOtherwise)
}

type refProps struct {
	Ref string `json:"$ref,omitempty"`
}

func (s *SchemaRef) setRefFromFQN(ref string, reg *descriptor.Registry) error {
	name, ok := fullyQualifiedNameToOpenAPIName(ref, reg)
	if !ok {
		return fmt.Errorf("setRefFromFQN: can't resolve OpenAPI name from '%v'", ref)
	}
	s.Ref = fmt.Sprintf("#/definitions/%s", name)
	return nil
}
