package convert

import (
	"strconv"
)

func String(val string) (string, error) {
	return val, nil
}

func Bool(val string) (bool, error) {
	return strconv.ParseBool(val)
}

func Float64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func Float32(val string) (float32, error) {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

func Int32(val string) (int32, error) {
	i, err := strconv.ParseInt(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func Uint64(val string) (uint64, error) {
	return strconv.ParseUint(val, 0, 64)
}

func Uint32(val string) (uint32, error) {
	i, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}
