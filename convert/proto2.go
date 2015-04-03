package convert

import (
	"github.com/golang/protobuf/proto"
)

func StringP(val string) (*string, error) {
	return proto.String(val), nil
}

func BoolP(val string) (*bool, error) {
	b, err := Bool(val)
	if err != nil {
		return nil, err
	}
	return proto.Bool(b), nil
}

func Float64P(val string) (*float64, error) {
	f, err := Float64(val)
	if err != nil {
		return nil, err
	}
	return proto.Float64(f), nil
}

func Float32P(val string) (*float32, error) {
	f, err := Float32(val)
	if err != nil {
		return nil, err
	}
	return proto.Float32(f), nil
}

func Int64P(val string) (*int64, error) {
	i, err := Int64(val)
	if err != nil {
		return nil, err
	}
	return proto.Int64(i), nil
}

func Int32P(val string) (*int32, error) {
	i, err := Int32(val)
	if err != nil {
		return nil, err
	}
	return proto.Int32(i), err
}

func Uint64P(val string) (*uint64, error) {
	i, err := Uint64(val)
	if err != nil {
		return nil, err
	}
	return proto.Uint64(i), err
}

func Uint32P(val string) (*uint32, error) {
	i, err := Uint32(val)
	if err != nil {
		return nil, err
	}
	return proto.Uint32(i), err
}
