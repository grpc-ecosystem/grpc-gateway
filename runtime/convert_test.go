package runtime_test

import (
	"reflect"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConvertTimestamp(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		output  *timestamppb.Timestamp
		wanterr bool
	}{
		{
			name:  "a valid RFC3339 timestamp",
			input: `"2016-05-10T10:19:13.123Z"`,
			output: &timestamppb.Timestamp{
				Seconds: 1462875553,
				Nanos:   123000000,
			},
			wanterr: false,
		},
		{
			name:  "a valid RFC3339 timestamp without double quotation",
			input: "2016-05-10T10:19:13.123Z",
			output: &timestamppb.Timestamp{
				Seconds: 1462875553,
				Nanos:   123000000,
			},
			wanterr: false,
		},
		{
			name:    "invalid timestamp",
			input:   `"05-10-2016T10:19:13.123Z"`,
			output:  nil,
			wanterr: true,
		},
		{
			name:    "JSON number",
			input:   "123",
			output:  nil,
			wanterr: true,
		},
		{
			name:    "JSON bool",
			input:   "true",
			output:  nil,
			wanterr: true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			ts, err := runtime.Timestamp(spec.input)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !proto.Equal(ts, spec.output):
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
		output  *durationpb.Duration
		wanterr bool
	}{
		{
			name:  "a valid duration",
			input: `"123.456s"`,
			output: &durationpb.Duration{
				Seconds: 123,
				Nanos:   456000000,
			},
			wanterr: false,
		},
		{
			name:  "a valid duration without double quotation",
			input: "123.456s",
			output: &durationpb.Duration{
				Seconds: 123,
				Nanos:   456000000,
			},
			wanterr: false,
		},
		{
			name:    "invalid duration",
			input:   `"123years"`,
			output:  nil,
			wanterr: true,
		},
		{
			name:    "JSON number",
			input:   "123",
			output:  nil,
			wanterr: true,
		},
		{
			name:    "JSON bool",
			input:   "true",
			output:  nil,
			wanterr: true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			ts, err := runtime.Duration(spec.input)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !proto.Equal(ts, spec.output):
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

func TestConvertStringSlice(t *testing.T) {
	specs := []struct {
		name   string
		input  string
		sep    string
		output []string
	}{
		{name: "multiple values", input: "a,b,c", sep: ",", output: []string{"a", "b", "c"}},
		{name: "single value", input: "a", sep: ",", output: []string{"a"}},
		{name: "custom separator", input: "a|b", sep: "|", output: []string{"a", "b"}},
		{name: "empty string yields one empty element", input: "", sep: ",", output: []string{""}},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.StringSlice(spec.input, spec.sep)
			if err != nil {
				t.Errorf("got unexpected error\n%#v", err)
			}
			if !reflect.DeepEqual(got, spec.output) {
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertBoolSlice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []bool
		wanterr bool
	}{
		{name: "valid values", input: "true,false,1,0", sep: ",", output: []bool{true, false, true, false}},
		{name: "single value", input: "true", sep: ",", output: []bool{true}},
		{name: "invalid element", input: "true,notabool", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.BoolSlice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertFloat64Slice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []float64
		wanterr bool
	}{
		{name: "valid values", input: "1.5,2,-3.25", sep: ",", output: []float64{1.5, 2, -3.25}},
		{name: "single value", input: "1.5", sep: ",", output: []float64{1.5}},
		{name: "invalid element", input: "1.5,notafloat", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.Float64Slice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertFloat32Slice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []float32
		wanterr bool
	}{
		{name: "valid values", input: "1.5,2,-3.25", sep: ",", output: []float32{1.5, 2, -3.25}},
		{name: "single value", input: "1.5", sep: ",", output: []float32{1.5}},
		{name: "invalid element", input: "1.5,notafloat", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.Float32Slice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertInt64Slice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []int64
		wanterr bool
	}{
		{name: "valid values", input: "1,-2,3", sep: ",", output: []int64{1, -2, 3}},
		{name: "single value", input: "42", sep: ",", output: []int64{42}},
		{name: "invalid element", input: "1,notanint", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.Int64Slice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertInt32Slice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []int32
		wanterr bool
	}{
		{name: "valid values", input: "1,-2,3", sep: ",", output: []int32{1, -2, 3}},
		{name: "single value", input: "42", sep: ",", output: []int32{42}},
		{name: "invalid element", input: "1,notanint", sep: ",", wanterr: true},
		{name: "overflows int32", input: "2147483648", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.Int32Slice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertUint64Slice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []uint64
		wanterr bool
	}{
		{name: "valid values", input: "1,2,3", sep: ",", output: []uint64{1, 2, 3}},
		{name: "single value", input: "42", sep: ",", output: []uint64{42}},
		{name: "invalid element", input: "1,-2", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.Uint64Slice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertUint32Slice(t *testing.T) {
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []uint32
		wanterr bool
	}{
		{name: "valid values", input: "1,2,3", sep: ",", output: []uint32{1, 2, 3}},
		{name: "single value", input: "42", sep: ",", output: []uint32{42}},
		{name: "invalid element", input: "1,-2", sep: ",", wanterr: true},
		{name: "overflows uint32", input: "4294967296", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.Uint32Slice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertBytesSlice(t *testing.T) {
	// "Zm9v" and "YmFy" are base64 for "foo" and "bar".
	specs := []struct {
		name    string
		input   string
		sep     string
		output  [][]byte
		wanterr bool
	}{
		{name: "valid values", input: "Zm9v,YmFy", sep: ",", output: [][]byte{[]byte("foo"), []byte("bar")}},
		{name: "single value", input: "Zm9v", sep: ",", output: [][]byte{[]byte("foo")}},
		{name: "invalid base64 element", input: "Zm9v,!!!!", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.BytesSlice(spec.input, spec.sep)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}

func TestConvertEnumSlice(t *testing.T) {
	enumValMap := map[string]int32{"A": 0, "B": 1, "C": 2}
	specs := []struct {
		name    string
		input   string
		sep     string
		output  []int32
		wanterr bool
	}{
		{name: "by name", input: "A,C", sep: ",", output: []int32{0, 2}},
		{name: "by numeric value", input: "0,2", sep: ",", output: []int32{0, 2}},
		{name: "mixed name and value", input: "A,1", sep: ",", output: []int32{0, 1}},
		{name: "unknown name", input: "A,Z", sep: ",", wanterr: true},
		{name: "numeric value out of range", input: "0,9", sep: ",", wanterr: true},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got, err := runtime.EnumSlice(spec.input, spec.sep, enumValMap)
			switch {
			case err != nil && !spec.wanterr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && spec.wanterr:
				t.Errorf("did not error when expected")
			case !spec.wanterr && !reflect.DeepEqual(got, spec.output):
				t.Errorf("got\n%#v\nexpected\n%#v", got, spec.output)
			}
		})
	}
}
