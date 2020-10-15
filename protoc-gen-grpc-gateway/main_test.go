package main

import (
	"testing"
)

func TestParseFlagsEmptyNoPanic(t *testing.T) {
	parseFlags("")
}

func TestParseFlags(t *testing.T) {
	parseFlags("allow_repeated_fields_in_body=true")
	if *allowRepeatedFieldsInBody != true {
		t.Errorf("flag allow_repeated_fields_in_body was not set correctly, wanted true got %v", *allowRepeatedFieldsInBody)
	}
}

func TestParseFlagsMultiple(t *testing.T) {
	parseFlags("allow_repeated_fields_in_body=true,import_prefix=foo")
	if *allowRepeatedFieldsInBody != true {
		t.Errorf("flag allow_repeated_fields_in_body was not set correctly, wanted 'true' got '%v'", *allowRepeatedFieldsInBody)
	}
	if *importPrefix != "foo" {
		t.Errorf("flag importPrefix was not set correctly, wanted 'foo' got '%v'", *importPrefix)
	}
}
