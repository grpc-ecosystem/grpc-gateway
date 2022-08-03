/*
 * A Bit of Everything
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0
 * Contact: none@example.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package abe

import (
	"time"
)

// Intentionally complicated message type to cover many features of Protobuf.
type ExamplepbABitOfEverything struct {
	SingleNested *ABitOfEverythingNested `json:"singleNested,omitempty"`
	Uuid string `json:"uuid"`
	Nested []ABitOfEverythingNested `json:"nested,omitempty"`
	// Float value field
	FloatValue float32 `json:"floatValue"`
	DoubleValue float64 `json:"doubleValue"`
	Int64Value string `json:"int64Value"`
	Uint64Value string `json:"uint64Value,omitempty"`
	Int32Value int32 `json:"int32Value,omitempty"`
	Fixed64Value string `json:"fixed64Value,omitempty"`
	Fixed32Value int64 `json:"fixed32Value,omitempty"`
	BoolValue bool `json:"boolValue,omitempty"`
	StringValue string `json:"stringValue,omitempty"`
	BytesValue string `json:"bytesValue,omitempty"`
	Uint32Value int64 `json:"uint32Value,omitempty"`
	EnumValue *ExamplepbNumericEnum `json:"enumValue,omitempty"`
	PathEnumValue *PathenumPathEnum `json:"pathEnumValue,omitempty"`
	NestedPathEnumValue *MessagePathEnumNestedPathEnum `json:"nestedPathEnumValue,omitempty"`
	Sfixed32Value int32 `json:"sfixed32Value,omitempty"`
	Sfixed64Value string `json:"sfixed64Value,omitempty"`
	Sint32Value int32 `json:"sint32Value,omitempty"`
	Sint64Value string `json:"sint64Value,omitempty"`
	RepeatedStringValue []string `json:"repeatedStringValue,omitempty"`
	OneofEmpty *interface{} `json:"oneofEmpty,omitempty"`
	OneofString string `json:"oneofString,omitempty"`
	MapValue map[string]ExamplepbNumericEnum `json:"mapValue,omitempty"`
	MappedStringValue map[string]string `json:"mappedStringValue,omitempty"`
	MappedNestedValue map[string]ABitOfEverythingNested `json:"mappedNestedValue,omitempty"`
	NonConventionalNameValue string `json:"nonConventionalNameValue,omitempty"`
	TimestampValue time.Time `json:"timestampValue,omitempty"`
	RepeatedEnumValue []ExamplepbNumericEnum `json:"repeatedEnumValue,omitempty"`
	// Repeated numeric enum description.
	RepeatedEnumAnnotation []ExamplepbNumericEnum `json:"repeatedEnumAnnotation,omitempty"`
	// Numeric enum description.
	EnumValueAnnotation *ExamplepbNumericEnum `json:"enumValueAnnotation,omitempty"`
	// Repeated string description.
	RepeatedStringAnnotation []string `json:"repeatedStringAnnotation,omitempty"`
	// Repeated nested object description.
	RepeatedNestedAnnotation []ABitOfEverythingNested `json:"repeatedNestedAnnotation,omitempty"`
	// Nested object description.
	NestedAnnotation *ABitOfEverythingNested `json:"nestedAnnotation,omitempty"`
	Int64OverrideType int64 `json:"int64OverrideType,omitempty"`
	RequiredStringViaFieldBehaviorAnnotation string `json:"requiredStringViaFieldBehaviorAnnotation"`
	OutputOnlyStringViaFieldBehaviorAnnotation string `json:"outputOnlyStringViaFieldBehaviorAnnotation,omitempty"`
	OptionalStringValue string `json:"optionalStringValue,omitempty"`
	// Only digits are allowed.
	ProductId []string `json:"productId,omitempty"`
	FieldOne string `json:"FieldOne,omitempty"`
	FieldTwo string `json:"FieldTwo"`
	FieldThree string `json:"FieldThree"`
}
