package runtime

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var currentPathParser PathParameterParser = &defaultPathParser{}

// PathParameterParser defines interface for all path parameter parsers
type PathParameterParser interface {
	Parse(msg proto.Message, pathParams map[string]string) error
}

// PopulatePathParameters parses path parameters
// into "msg" using current path parser
func PopulatePathParameters(msg proto.Message, pathParams map[string]string) error {
	return currentPathParser.Parse(msg, pathParams)
}

type defaultPathParser struct{}

// Parse populates "values" into "msg".
// This is mainly similar to the query parser except support for maps and proto messages
// is removed. Repeated fields (comma separated) are supported, but any unsupported field types
// will just be skipped to let the query parser handle.
func (*defaultPathParser) Parse(msg proto.Message, values map[string]string) error {
	for key, value := range values {
		fieldPath := strings.Split(key, ".")
		if err := populateFieldValueFromPath(msg.ProtoReflect(), fieldPath, value); err != nil {
			return err
		}
	}
	return nil
}

// PopulateFieldFromPathParam sets a value in a nested Protobuf structure.
func PopulateFieldFromPathParam(msg proto.Message, fieldPathString string, value string) error {
	fieldPath := strings.Split(fieldPathString, ".")
	return populateFieldValueFromPath(msg.ProtoReflect(), fieldPath, value)
}

func populateFieldValueFromPath(msgValue protoreflect.Message, fieldPath []string, value string) error {
	if len(fieldPath) < 1 {
		return errors.New("no field path")
	}
	if value == "" {
		return errors.New("no value provided")
	}

	var fieldDescriptor protoreflect.FieldDescriptor
	for i, fieldName := range fieldPath {
		fields := msgValue.Descriptor().Fields()

		// Get field by name
		fieldDescriptor = fields.ByName(protoreflect.Name(fieldName))
		if fieldDescriptor == nil {
			fieldDescriptor = fields.ByJSONName(fieldName)
			if fieldDescriptor == nil {
				// We're not returning an error here because this could just be
				// an extra query parameter that isn't part of the request.
				//grpclog.Infof("field not found in %q: %q", msgValue.Descriptor().FullName(), strings.Join(fieldPath, "."))
				return nil
			}
		}

		// If this is the last element, we're done
		if i == len(fieldPath)-1 {
			break
		}

		// Only singular message fields are allowed
		if fieldDescriptor.Message() == nil || fieldDescriptor.Cardinality() == protoreflect.Repeated {
			return fmt.Errorf("invalid path: %q is not a message", fieldName)
		}

		// Get the nested message
		msgValue = msgValue.Mutable(fieldDescriptor).Message()
	}

	// Check if oneof already set
	if of := fieldDescriptor.ContainingOneof(); of != nil {
		if f := msgValue.WhichOneof(of); f != nil {
			return fmt.Errorf("field already set for oneof %q", of.FullName().Name())
		}
	}

	switch {
	case fieldDescriptor.IsList():
		return populateRepeatedPathField(fieldDescriptor, msgValue.Mutable(fieldDescriptor).List(), value)
	case fieldDescriptor.IsMap():
		// Don't continue further if it is a proto message. We shouldn't be handling these.
		return nil
		//case fieldDescriptor.Kind() == protoreflect.MessageKind:
		//	// Don't continue further if it is a proto message. We shouldn't be handling these.
		//	return nil
		//case fieldDescriptor.Kind() == protoreflect.GroupKind:
		//	return nil
	}


	return populatePathField(fieldDescriptor, msgValue, value)
}

func populatePathField(fieldDescriptor protoreflect.FieldDescriptor, msgValue protoreflect.Message, value string) error {
	v, err := parsePathField(fieldDescriptor, value)
	if err != nil {
		return fmt.Errorf("parsing field %q: %w", fieldDescriptor.FullName().Name(), err)
	}

	msgValue.Set(fieldDescriptor, v)
	return nil
}

func populateRepeatedPathField(fieldDescriptor protoreflect.FieldDescriptor, list protoreflect.List, value string) error {
	for _, value := range strings.Split(value, ",") {
		v, err := parsePathField(fieldDescriptor, value)
		if err != nil {
			return fmt.Errorf("parsing list %q: %w", fieldDescriptor.FullName().Name(), err)
		}
		list.Append(v)
	}

	return nil
}

func parsePathField(fieldDescriptor protoreflect.FieldDescriptor, value string) (protoreflect.Value, error) {
	switch fieldDescriptor.Kind() {
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBool(v), nil
	case protoreflect.EnumKind:
		enum, err := protoregistry.GlobalTypes.FindEnumByName(fieldDescriptor.Enum().FullName())
		switch {
		case errors.Is(err, protoregistry.NotFound):
			return protoreflect.Value{}, fmt.Errorf("enum %q is not registered", fieldDescriptor.Enum().FullName())
		case err != nil:
			return protoreflect.Value{}, fmt.Errorf("failed to look up enum: %w", err)
		}
		// Look for enum by name
		v := enum.Descriptor().Values().ByName(protoreflect.Name(value))
		if v == nil {
			i, err := strconv.Atoi(value)
			if err != nil {
				return protoreflect.Value{}, fmt.Errorf("%q is not a valid value", value)
			}
			// Look for enum by number
			v = enum.Descriptor().Values().ByNumber(protoreflect.EnumNumber(i))
			if v == nil {
				return protoreflect.Value{}, fmt.Errorf("%q is not a valid value", value)
			}
		}
		return protoreflect.ValueOfEnum(v.Number()), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt64(v), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint64(v), nil
	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil
	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat64(v), nil
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(value), nil
	case protoreflect.BytesKind:
		v, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBytes(v), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return parsePathMessage(fieldDescriptor.Message(), value)
	default:
		panic(fmt.Sprintf("unknown field kind: %v", fieldDescriptor.Kind()))
	}
}

func parsePathMessage(msgDescriptor protoreflect.MessageDescriptor, value string) (protoreflect.Value, error) {
	var msg proto.Message
	switch msgDescriptor.FullName() {
	case "google.protobuf.Timestamp":
		if value == "null" {
			break
		}
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = timestamppb.New(t)
	case "google.protobuf.Duration":
		if value == "null" {
			break
		}
		d, err := time.ParseDuration(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = durationpb.New(d)
	case "google.protobuf.DoubleValue":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.DoubleValue{Value: v}
	case "google.protobuf.FloatValue":
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.FloatValue{Value: float32(v)}
	case "google.protobuf.Int64Value":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.Int64Value{Value: v}
	case "google.protobuf.Int32Value":
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.Int32Value{Value: int32(v)}
	case "google.protobuf.UInt64Value":
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.UInt64Value{Value: v}
	case "google.protobuf.UInt32Value":
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.UInt32Value{Value: uint32(v)}
	case "google.protobuf.BoolValue":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.BoolValue{Value: v}
	case "google.protobuf.StringValue":
		msg = &wrapperspb.StringValue{Value: value}
	case "google.protobuf.BytesValue":
		v, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrapperspb.BytesValue{Value: v}
	case "google.protobuf.FieldMask":
		fm := &field_mask.FieldMask{}
		fm.Paths = append(fm.Paths, strings.Split(value, ",")...)
		msg = fm
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported message type: %q", string(msgDescriptor.FullName()))
	}

	return protoreflect.ValueOfMessage(msg.ProtoReflect()), nil
}
