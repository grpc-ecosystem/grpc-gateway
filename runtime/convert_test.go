package runtime_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func TestConvertTimestamp(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		output  *timestamp.Timestamp
		wanterr error
	}{
		{
			name:  "a valid RFC3339 timestamp",
			input: `"2016-05-10T10:19:13.123Z"`,
			output: &timestamp.Timestamp{
				Seconds: 1462875553,
				Nanos:   123000000,
			},
			wanterr: nil,
		},
		{
			name:    "invalid timestamp",
			input:   `"05-10-2016T10:19:13.123Z"`,
			output:  nil,
			wanterr: fmt.Errorf(`bad Timestamp: parsing time "05-10-2016T10:19:13.123Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-2016T10:19:13.123Z" as "2006"`),
		},
		{
			name:   "JSON number",
			input:  "123",
			output: nil,
			wanterr: &json.UnmarshalTypeError{
				Value:  "number",
				Type:   reflect.TypeOf("123"),
				Offset: 3,
			},
		},
		{
			name:   "JSON bool",
			input:  "true",
			output: nil,
			wanterr: &json.UnmarshalTypeError{
				Value:  "bool",
				Type:   reflect.TypeOf("123"),
				Offset: 4,
			},
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			ts, err := runtime.Timestamp(spec.input)
			if spec.wanterr != nil {
				if !reflect.DeepEqual(err, spec.wanterr) {
					t.Errorf("got unexpected error\n%#v\nexpected\n%#v", err, spec.wanterr)
				}
				return
			}
			if !proto.Equal(ts, spec.output) {
				t.Errorf(
					"when testing %s; got\n%#v\nexpected\n%#v",
					spec.name,
					ts,
					spec.output,
				)
			}
		})
	}
}

func TestConvertDuration(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		output  *duration.Duration
		wanterr error
	}{
		{
			name:  "a valid duration",
			input: `"123.456s"`,
			output: &duration.Duration{
				Seconds: 123,
				Nanos:   456000000,
			},
			wanterr: nil,
		},
		{
			name:    "invalid duration",
			input:   `"123years"`,
			output:  nil,
			wanterr: fmt.Errorf(`bad Duration: time: unknown unit years in duration 123years`),
		},
		{
			name:   "JSON number",
			input:  "123",
			output: nil,
			wanterr: &json.UnmarshalTypeError{
				Value:  "number",
				Type:   reflect.TypeOf("123"),
				Offset: 3,
			},
		},
		{
			name:   "JSON bool",
			input:  "true",
			output: nil,
			wanterr: &json.UnmarshalTypeError{
				Value:  "bool",
				Type:   reflect.TypeOf("123"),
				Offset: 4,
			},
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			ts, err := runtime.Duration(spec.input)
			if spec.wanterr != nil {
				if !reflect.DeepEqual(err, spec.wanterr) {
					t.Errorf("got unexpected error\n%#v\nexpected\n%#v", err, spec.wanterr)
				}
				return
			}
			if !proto.Equal(ts, spec.output) {
				t.Errorf(
					"when testing %s; got\n%#v\nexpected\n%#v",
					spec.name,
					ts,
					spec.output,
				)
			}
		})
	}
}
