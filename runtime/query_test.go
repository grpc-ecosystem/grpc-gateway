package runtime_test

import (
	"errors"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	wrapperspb "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/utilities"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func BenchmarkPopulateQueryParameters(b *testing.B) {
	timeT := time.Date(2016, time.December, 15, 12, 23, 32, 49, time.UTC)
	timeStr := timeT.Format(time.RFC3339Nano)

	durationT := 13 * time.Hour
	durationStr := durationT.String()

	fieldmaskStr := "float_value,double_value"

	msg := &examplepb.Proto3Message{}
	values := url.Values{
		"float_value":            {"1.5"},
		"double_value":           {"2.5"},
		"int64_value":            {"-1"},
		"int32_value":            {"-2"},
		"uint64_value":           {"3"},
		"uint32_value":           {"4"},
		"bool_value":             {"true"},
		"string_value":           {"str"},
		"bytes_value":            {"Ynl0ZXM="},
		"repeated_value":         {"a", "b", "c"},
		"enum_value":             {"1"},
		"repeated_enum":          {"1", "2", "0"},
		"timestamp_value":        {timeStr},
		"duration_value":         {durationStr},
		"fieldmask_value":        {fieldmaskStr},
		"wrapper_float_value":    {"1.5"},
		"wrapper_double_value":   {"2.5"},
		"wrapper_int64_value":    {"-1"},
		"wrapper_int32_value":    {"-2"},
		"wrapper_u_int64_value":  {"3"},
		"wrapper_u_int32_value":  {"4"},
		"wrapper_bool_value":     {"true"},
		"wrapper_string_value":   {"str"},
		"wrapper_bytes_value":    {"Ynl0ZXM="},
		"map_value[key]":         {"value"},
		"map_value[second]":      {"bar"},
		"map_value[third]":       {"zzz"},
		"map_value[fourth]":      {""},
		`map_value[~!@#$%^&*()]`: {"value"},
		"map_value2[key]":        {"-2"},
		"map_value3[-2]":         {"value"},
		"map_value4[key]":        {"-1"},
		"map_value5[-1]":         {"value"},
		"map_value6[key]":        {"3"},
		"map_value7[3]":          {"value"},
		"map_value8[key]":        {"4"},
		"map_value9[4]":          {"value"},
		"map_value10[key]":       {"1.5"},
		"map_value11[1.5]":       {"value"},
		"map_value12[key]":       {"2.5"},
		"map_value13[2.5]":       {"value"},
		"map_value14[key]":       {"true"},
		"map_value15[true]":      {"value"},
	}
	filter := utilities.NewDoubleArray([][]string{
		{"bool_value"}, {"repeated_value"},
	})

	for i := 0; i < b.N; i++ {
		_ = runtime.PopulateQueryParameters(msg, values, filter)
	}
}

func TestPopulateParameters(t *testing.T) {
	timeT := time.Date(2016, time.December, 15, 12, 23, 32, 49, time.UTC)
	timeStr := timeT.Format(time.RFC3339Nano)
	timePb, err := ptypes.TimestampProto(timeT)
	if err != nil {
		t.Fatalf("Couldn't setup timestamp in Protobuf format: %v", err)
	}

	durationT := 13 * time.Hour
	durationStr := durationT.String()
	durationPb := ptypes.DurationProto(durationT)

	fieldmaskStr := "float_value,double_value"
	fieldmaskPb := &field_mask.FieldMask{Paths: []string{"float_value", "double_value"}}

	for i, spec := range []struct {
		values  url.Values
		filter  *utilities.DoubleArray
		want    proto.Message
		wanterr error
	}{
		{
			values: url.Values{
				"float_value":            {"1.5"},
				"double_value":           {"2.5"},
				"int64_value":            {"-1"},
				"int32_value":            {"-2"},
				"uint64_value":           {"3"},
				"uint32_value":           {"4"},
				"bool_value":             {"true"},
				"string_value":           {"str"},
				"bytes_value":            {"Ynl0ZXM="},
				"repeated_value":         {"a", "b", "c"},
				"repeated_message":       {"1", "2", "3"},
				"enum_value":             {"1"},
				"repeated_enum":          {"1", "2", "0"},
				"timestamp_value":        {timeStr},
				"duration_value":         {durationStr},
				"fieldmask_value":        {fieldmaskStr},
				"wrapper_float_value":    {"1.5"},
				"wrapper_double_value":   {"2.5"},
				"wrapper_int64_value":    {"-1"},
				"wrapper_int32_value":    {"-2"},
				"wrapper_u_int64_value":  {"3"},
				"wrapper_u_int32_value":  {"4"},
				"wrapper_bool_value":     {"true"},
				"wrapper_string_value":   {"str"},
				"wrapper_bytes_value":    {"Ynl0ZXM="},
				"map_value[key]":         {"value"},
				"map_value[second]":      {"bar"},
				"map_value[third]":       {"zzz"},
				"map_value[fourth]":      {""},
				`map_value[~!@#$%^&*()]`: {"value"},
				"map_value2[key]":        {"-2"},
				"map_value3[-2]":         {"value"},
				"map_value4[key]":        {"-1"},
				"map_value5[-1]":         {"value"},
				"map_value6[key]":        {"3"},
				"map_value7[3]":          {"value"},
				"map_value8[key]":        {"4"},
				"map_value9[4]":          {"value"},
				"map_value10[key]":       {"1.5"},
				"map_value11[1.5]":       {"value"},
				"map_value12[key]":       {"2.5"},
				"map_value13[2.5]":       {"value"},
				"map_value14[key]":       {"true"},
				"map_value15[true]":      {"value"},
				"map_value16[key]":       {"2"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				FloatValue:         1.5,
				DoubleValue:        2.5,
				Int64Value:         -1,
				Int32Value:         -2,
				Uint64Value:        3,
				Uint32Value:        4,
				BoolValue:          true,
				StringValue:        "str",
				BytesValue:         []byte("bytes"),
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
				WrapperBytesValue:  &wrapperspb.BytesValue{Value: []byte("bytes")},
				MapValue: map[string]string{
					"key":         "value",
					"second":      "bar",
					"third":       "zzz",
					"fourth":      "",
					`~!@#$%^&*()`: "value",
				},
				MapValue2:  map[string]int32{"key": -2},
				MapValue3:  map[int32]string{-2: "value"},
				MapValue4:  map[string]int64{"key": -1},
				MapValue5:  map[int64]string{-1: "value"},
				MapValue6:  map[string]uint32{"key": 3},
				MapValue7:  map[uint32]string{3: "value"},
				MapValue8:  map[string]uint64{"key": 4},
				MapValue9:  map[uint64]string{4: "value"},
				MapValue10: map[string]float32{"key": 1.5},
				MapValue12: map[string]float64{"key": 2.5},
				MapValue14: map[string]bool{"key": true},
				MapValue15: map[bool]string{true: "value"},
				MapValue16: map[string]*wrapperspb.UInt64Value{"key": {Value: 2}},
			},
		},
		{
			values: url.Values{
				"floatValue":         {"1.5"},
				"doubleValue":        {"2.5"},
				"int64Value":         {"-1"},
				"int32Value":         {"-2"},
				"uint64Value":        {"3"},
				"uint32Value":        {"4"},
				"boolValue":          {"true"},
				"stringValue":        {"str"},
				"bytesValue":         {"Ynl0ZXM="},
				"repeatedValue":      {"a", "b", "c"},
				"enumValue":          {"1"},
				"repeatedEnum":       {"1", "2", "0"},
				"timestampValue":     {timeStr},
				"durationValue":      {durationStr},
				"fieldmaskValue":     {fieldmaskStr},
				"wrapperFloatValue":  {"1.5"},
				"wrapperDoubleValue": {"2.5"},
				"wrapperInt64Value":  {"-1"},
				"wrapperInt32Value":  {"-2"},
				"wrapperUInt64Value": {"3"},
				"wrapperUInt32Value": {"4"},
				"wrapperBoolValue":   {"true"},
				"wrapperStringValue": {"str"},
				"wrapperBytesValue":  {"Ynl0ZXM="},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				FloatValue:         1.5,
				DoubleValue:        2.5,
				Int64Value:         -1,
				Int32Value:         -2,
				Uint64Value:        3,
				Uint32Value:        4,
				BoolValue:          true,
				StringValue:        "str",
				BytesValue:         []byte("bytes"),
				RepeatedValue:      []string{"a", "b", "c"},
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
				WrapperBytesValue:  &wrapperspb.BytesValue{Value: []byte("bytes")},
			},
		},
		{
			values: url.Values{
				"enum_value":    {"Z"},
				"repeated_enum": {"X", "2", "0"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				EnumValue:    examplepb.EnumValue_Z,
				RepeatedEnum: []examplepb.EnumValue{examplepb.EnumValue_X, examplepb.EnumValue_Z, examplepb.EnumValue_X},
			},
		},
		{
			values: url.Values{
				"float_value":    {"1.5"},
				"double_value":   {"2.5"},
				"int64_value":    {"-1"},
				"int32_value":    {"-2"},
				"uint64_value":   {"3"},
				"uint32_value":   {"4"},
				"bool_value":     {"true"},
				"string_value":   {"str"},
				"repeated_value": {"a", "b", "c"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto2Message{
				FloatValue:    proto.Float32(1.5),
				DoubleValue:   proto.Float64(2.5),
				Int64Value:    proto.Int64(-1),
				Int32Value:    proto.Int32(-2),
				Uint64Value:   proto.Uint64(3),
				Uint32Value:   proto.Uint32(4),
				BoolValue:     proto.Bool(true),
				StringValue:   proto.String("str"),
				RepeatedValue: []string{"a", "b", "c"},
			},
		},
		{
			values: url.Values{
				"floatValue":    {"1.5"},
				"doubleValue":   {"2.5"},
				"int64Value":    {"-1"},
				"int32Value":    {"-2"},
				"uint64Value":   {"3"},
				"uint32Value":   {"4"},
				"boolValue":     {"true"},
				"stringValue":   {"str"},
				"repeatedValue": {"a", "b", "c"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto2Message{
				FloatValue:    proto.Float32(1.5),
				DoubleValue:   proto.Float64(2.5),
				Int64Value:    proto.Int64(-1),
				Int32Value:    proto.Int32(-2),
				Uint64Value:   proto.Uint64(3),
				Uint32Value:   proto.Uint32(4),
				BoolValue:     proto.Bool(true),
				StringValue:   proto.String("str"),
				RepeatedValue: []string{"a", "b", "c"},
			},
		},
		{
			values: url.Values{
				"nested.nested.nested.repeated_value": {"a", "b", "c"},
				"nested.nested.nested.string_value":   {"s"},
				"nested.nested.string_value":          {"t"},
				"nested.string_value":                 {"u"},
				"nested.nested.map_value[first]":      {"foo"},
				"nested.nested.map_value[second]":     {"bar"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				Nested: &examplepb.Proto3Message{
					Nested: &examplepb.Proto3Message{
						MapValue: map[string]string{
							"first":  "foo",
							"second": "bar",
						},
						Nested: &examplepb.Proto3Message{
							RepeatedValue: []string{"a", "b", "c"},
							StringValue:   "s",
						},
						StringValue: "t",
					},
					StringValue: "u",
				},
			},
		},
		{
			values: url.Values{
				"oneof_string_value": {"foobar"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				OneofValue: &examplepb.Proto3Message_OneofStringValue{
					OneofStringValue: "foobar",
				},
			},
		},
		{
			values: url.Values{
				"oneofStringValue": {"foobar"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				OneofValue: &examplepb.Proto3Message_OneofStringValue{
					OneofStringValue: "foobar",
				},
			},
		},
		{
			values: url.Values{
				"oneof_bool_value": {"true"},
			},
			filter: utilities.NewDoubleArray(nil),
			want: &examplepb.Proto3Message{
				OneofValue: &examplepb.Proto3Message_OneofBoolValue{
					OneofBoolValue: true,
				},
			},
		},
		{
			// Don't allow setting a oneof more than once
			values: url.Values{
				"oneof_bool_value":   {"true"},
				"oneof_string_value": {"foobar"},
			},
			filter:  utilities.NewDoubleArray(nil),
			want:    &examplepb.Proto3Message{},
			wanterr: errors.New("field already set for oneof \"oneof_value\""),
		},
		{
			// Error when there are too many values
			values: url.Values{
				"uint64_value": {"1", "2"},
			},
			filter:  utilities.NewDoubleArray(nil),
			want:    &examplepb.Proto3Message{},
			wanterr: errors.New("too many values for field \"uint64_value\": 1, 2"),
		},
		{
			// Error when dereferencing a list of messages
			values: url.Values{
				"repeated_message.value": {"1"},
			},
			filter:  utilities.NewDoubleArray(nil),
			want:    &examplepb.Proto3Message{},
			wanterr: errors.New("invalid path: \"repeated_message\" is not a message"),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			msg := spec.want.ProtoReflect().New().Interface()
			err := runtime.PopulateQueryParameters(msg, spec.values, spec.filter)
			if spec.wanterr != nil {
				if err == nil || err.Error() != spec.wanterr.Error() {
					t.Errorf("runtime.PopulateQueryParameters(msg, %v, %v) failed with %q; want error %q", spec.values, spec.filter, err, spec.wanterr)
				}
				return
			}

			if err != nil {
				t.Errorf("runtime.PopulateQueryParameters(msg, %v, %v) failed with %v; want success", spec.values, spec.filter, err)
				return
			}
			if diff := cmp.Diff(spec.want, msg, protocmp.Transform()); diff != "" {
				t.Errorf("runtime.PopulateQueryParameters(msg, %v, %v): %s", spec.values, spec.filter, diff)
			}
		})
	}
}

func TestPopulateParametersWithFilters(t *testing.T) {
	for _, spec := range []struct {
		values url.Values
		filter *utilities.DoubleArray
		want   proto.Message
	}{
		{
			values: url.Values{
				"bool_value":     {"true"},
				"string_value":   {"str"},
				"repeated_value": {"a", "b", "c"},
			},
			filter: utilities.NewDoubleArray([][]string{
				{"bool_value"}, {"repeated_value"},
			}),
			want: &examplepb.Proto3Message{
				StringValue: "str",
			},
		},
		{
			values: url.Values{
				"nested.nested.bool_value":   {"true"},
				"nested.nested.string_value": {"str"},
				"nested.string_value":        {"str"},
				"string_value":               {"str"},
			},
			filter: utilities.NewDoubleArray([][]string{
				{"nested"},
			}),
			want: &examplepb.Proto3Message{
				StringValue: "str",
			},
		},
		{
			values: url.Values{
				"nested.nested.bool_value":   {"true"},
				"nested.nested.string_value": {"str"},
				"nested.string_value":        {"str"},
				"string_value":               {"str"},
			},
			filter: utilities.NewDoubleArray([][]string{
				{"nested", "nested"},
			}),
			want: &examplepb.Proto3Message{
				Nested: &examplepb.Proto3Message{
					StringValue: "str",
				},
				StringValue: "str",
			},
		},
		{
			values: url.Values{
				"nested.nested.bool_value":   {"true"},
				"nested.nested.string_value": {"str"},
				"nested.string_value":        {"str"},
				"string_value":               {"str"},
			},
			filter: utilities.NewDoubleArray([][]string{
				{"nested", "nested", "string_value"},
			}),
			want: &examplepb.Proto3Message{
				Nested: &examplepb.Proto3Message{
					StringValue: "str",
					Nested: &examplepb.Proto3Message{
						BoolValue: true,
					},
				},
				StringValue: "str",
			},
		},
	} {
		msg := spec.want.ProtoReflect().New().Interface()
		err := runtime.PopulateQueryParameters(msg, spec.values, spec.filter)
		if err != nil {
			t.Errorf("runtime.PoplateQueryParameters(msg, %v, %v) failed with %v; want success", spec.values, spec.filter, err)
			continue
		}
		if got, want := msg, spec.want; !proto.Equal(got, want) {
			t.Errorf("runtime.PopulateQueryParameters(msg, %v, %v = %v; want %v", spec.values, spec.filter, got, want)
		}
	}
}

func TestPopulateQueryParametersWithInvalidNestedParameters(t *testing.T) {
	for _, spec := range []struct {
		msg    proto.Message
		values url.Values
		filter *utilities.DoubleArray
	}{
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"float_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"double_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"int64_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"int32_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"uint64_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"uint32_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"bool_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"string_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"repeated_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"enum_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"enum_value.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
		{
			msg: &examplepb.Proto3Message{},
			values: url.Values{
				"repeated_enum.nested": {"test"},
			},
			filter: utilities.NewDoubleArray(nil),
		},
	} {
		spec.msg = spec.msg.ProtoReflect().New().Interface()
		err := runtime.PopulateQueryParameters(spec.msg, spec.values, spec.filter)
		if err == nil {
			t.Errorf("runtime.PopulateQueryParameters(msg, %v, %v) did not fail; want error", spec.values, spec.filter)
		}
	}
}
