package runtime_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func BenchmarkPopulatePathParameters(b *testing.B) {
	timeT := time.Date(2016, time.December, 15, 12, 23, 32, 49, time.UTC)
	timeStr := timeT.Format(time.RFC3339Nano)

	durationT := 13 * time.Hour
	durationStr := durationT.String()

	fieldmaskStr := "float_value,double_value"

	msg := &examplepb.Proto3Message{}
	values := map[string]string{
		"float_value":            "1.5",
		"double_value":           "2.5",
		"int64_value":            "-1",
		"int32_value":            "-2",
		"uint64_value":           "3",
		"uint32_value":           "4",
		"bool_value":             "true",
		"string_value":           "str",
		"bytes_value":            "Ynl0ZXM=",
		"repeated_value":         "a,b,c",
		"repeated_message":       "1,2,3",
		"enum_value":             "1",
		"repeated_enum":          "1,2,0",
		"timestamp_value":        timeStr,
		"duration_value":         durationStr,
		"fieldmask_value":        fieldmaskStr,
		"optional_string_value":  "optional-str",
		"wrapper_float_value":    "1.5",
		"wrapper_double_value":   "2.5",
		"wrapper_int64_value":    "-1",
		"wrapper_int32_value":    "-2",
		"wrapper_u_int64_value":  "3",
		"wrapper_u_int32_value":  "4",
		"wrapper_bool_value":     "true",
		"wrapper_string_value":   "str",
		"wrapper_bytes_value":    "Ynl0ZXM=",
	}

	for i := 0; i < b.N; i++ {
		_ = runtime.PopulatePathParameters(msg, values)
	}
}

func TestPopulatePathParameters(t *testing.T) {
	timeT := time.Date(2016, time.December, 15, 12, 23, 32, 49, time.UTC)
	timeStr := timeT.Format(time.RFC3339Nano)
	timePb := timestamppb.New(timeT)

	durationT := 13 * time.Hour
	durationStr := durationT.String()
	durationPb := durationpb.New(durationT)

	fieldmaskStr := "float_value,double_value"
	fieldmaskPb := &field_mask.FieldMask{Paths: []string{"float_value", "double_value"}}

	for i, spec := range []struct {
		values  map[string]string
		want    proto.Message
		wanterr error
	}{
		{
			values: map[string]string{
				"float_value":            "1.5",
				"double_value":           "2.5",
				"int64_value":            "-1",
				"int32_value":            "-2",
				"uint64_value":           "3",
				"uint32_value":           "4",
				"bool_value":             "true",
				"string_value":           "str",
				"bytes_value":            "YWJjMTIzIT8kKiYoKSctPUB-",
				"repeated_value":         "a,b,c",
				"repeated_message":       "1,2,3",
				"enum_value":             "1",
				"repeated_enum":          "1,2,0",
				"timestamp_value":        timeStr,
				"duration_value":         durationStr,
				"fieldmask_value":        fieldmaskStr,
				"wrapper_float_value":    "1.5",
				"wrapper_double_value":   "2.5",
				"wrapper_int64_value":    "-1",
				"wrapper_int32_value":    "-2",
				"wrapper_u_int64_value":  "3",
				"wrapper_u_int32_value":  "4",
				"wrapper_bool_value":     "true",
				"wrapper_string_value":   "str",
				"wrapper_bytes_value":    "YWJjMTIzIT8kKiYoKSctPUB-",
			},
			want: &examplepb.Proto3Message{
				FloatValue:         1.5,
				DoubleValue:        2.5,
				Int64Value:         -1,
				Int32Value:         -2,
				Uint64Value:        3,
				Uint32Value:        4,
				BoolValue:          true,
				StringValue:        "str",
				BytesValue:         []byte("abc123!?$*&()'-=@~"),
				RepeatedValue:      []string{"a", "b", "c"},
				RepeatedMessage:    []*wrapperspb.UInt64Value{{Value: 1}, {Value: 2}, {Value: 3}},
				EnumValue:          examplepb.EnumValue_Y,
				RepeatedEnum:       []examplepb.EnumValue{examplepb.EnumValue_Y, examplepb.EnumValue_Z, examplepb.EnumValue_X},
				TimestampValue:     timePb,
				DurationValue:      durationPb,
				FieldmaskValue:     fieldmaskPb,
				WrapperFloatValue:  &wrapperspb.FloatValue{Value: 1.5},
				WrapperDoubleValue: &wrapperspb.DoubleValue{Value: 2.5},
				WrapperInt64Value:  &wrapperspb.Int64Value{Value: -1},
				WrapperInt32Value:  &wrapperspb.Int32Value{Value: -2},
				WrapperUInt64Value: &wrapperspb.UInt64Value{Value: 3},
				WrapperUInt32Value: &wrapperspb.UInt32Value{Value: 4},
				WrapperBoolValue:   &wrapperspb.BoolValue{Value: true},
				WrapperStringValue: &wrapperspb.StringValue{Value: "str"},
				WrapperBytesValue:  &wrapperspb.BytesValue{Value: []byte("abc123!?$*&()'-=@~")},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			msg := spec.want.ProtoReflect().New().Interface()
			err := runtime.PopulatePathParameters(msg, spec.values)
			if spec.wanterr != nil {
				if err == nil || err.Error() != spec.wanterr.Error() {
					t.Errorf("runtime.PopulatePathParameters(msg, %v) failed with %q; want error %q", spec.values, err, spec.wanterr)
				}
				return
			}

			if err != nil {
				t.Errorf("runtime.PopulatePathParameters(msg, %v) failed with %v; want success", spec.values, err)
				return
			}
			if diff := cmp.Diff(spec.want, msg, protocmp.Transform()); diff != "" {
				t.Errorf("runtime.PopulatePathParameters(msg, %v): %s", spec.values, diff)
			}
		})
	}
}

func TestPopulatePathParametersWithInvalidNestedParameters(t *testing.T) {
	for _, spec := range []struct {
		msg    proto.Message
		values map[string]string
	}{
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"float_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"double_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"int64_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"int32_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"uint64_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"uint32_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"bool_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"string_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"repeated_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"enum_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"enum_value.nested": "test",
			},
		},
		{
			msg: &examplepb.Proto3Message{},
			values: map[string]string{
				"repeated_enum.nested": "test",
			},
		},
	} {
		spec.msg = spec.msg.ProtoReflect().New().Interface()
		err := runtime.PopulatePathParameters(spec.msg, spec.values)
		if err == nil {
			t.Errorf("runtime.PopulatePathParameters(msg, %v) did not fail; want error", spec.values)
		}
	}
}
