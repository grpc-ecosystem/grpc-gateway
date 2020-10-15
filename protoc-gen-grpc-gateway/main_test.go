package main

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

func TestParseFlagsEmptyNoPanic(t *testing.T) {
	reg := descriptor.NewRegistry()
	parseFlags(reg, "")
}

func TestParseFlags(t *testing.T) {
	reg := descriptor.NewRegistry()
	parseFlags(reg, "allow_repeated_fields_in_body=true")
	if *allowRepeatedFieldsInBody != true {
		t.Errorf("flag allow_repeated_fields_in_body was not set correctly, wanted true got %v", *allowRepeatedFieldsInBody)
	}
}

func TestParseFlagsMultiple(t *testing.T) {
	reg := descriptor.NewRegistry()
	parseFlags(reg, "allow_repeated_fields_in_body=true,import_prefix=foo")
	if *allowRepeatedFieldsInBody != true {
		t.Errorf("flag allow_repeated_fields_in_body was not set correctly, wanted 'true' got '%v'", *allowRepeatedFieldsInBody)
	}
	if *importPrefix != "foo" {
		t.Errorf("flag importPrefix was not set correctly, wanted 'foo' got '%v'", *importPrefix)
	}
}
