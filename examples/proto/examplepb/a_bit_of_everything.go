package examplepb

import (
	"encoding/json"
	"fmt"
	"strings"

	 "github.com/golang/protobuf/jsonpb"
)


func (m MimicObjectHidden) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return []byte(`"` + strings.Join([]string{m.HiddenValueOne, m.HiddenValueTwo, m.HiddenEnum.String()}, "/") + `"`), nil
}

func (m *MimicObjectHidden) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, d []byte) error {
	var s string

	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}

	fields := strings.Split(s, "/")

	if len(fields) != 3 {
		return fmt.Errorf("Unrecognized format for string %q", s)
	}

	m.HiddenValueOne = strings.TrimSpace(fields[0])
	m.HiddenValueTwo = strings.TrimSpace(fields[1])

	if v, ok := MimicObjectHidden_HiddenEnum_value[fields[2]]; !ok {
		return fmt.Errorf("Unrecognized HiddenEnum value %q", fields[2])
	} else {
		m.HiddenEnum = MimicObjectHidden_HiddenEnum(v)
	}

	return nil
}
