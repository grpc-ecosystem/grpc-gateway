package main

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

// TODO(asnyder): These tests are covering a function parseFlags() which became unused after https://github.com/grpc-ecosystem/grpc-gateway/pull/1756. The function and tests may be deleted.

func TestParseFlagsEmptyNoPanic(t *testing.T) {
	reg := descriptor.NewRegistry()
	parseFlags(reg, "")
}

func TestParseFlags(t *testing.T) {
	reg := descriptor.NewRegistry()
	parseFlags(reg, "generate_unbound_methods=true")
	if *generateUnboundMethods != true {
		t.Errorf("flag generate_unbound_methods was not set correctly, wanted true got %v", *generateUnboundMethods)
	}
}

func TestParseFlagsMultiple(t *testing.T) {
	reg := descriptor.NewRegistry()
	parseFlags(reg, "generate_unbound_methods=true,repeated_path_param_separator=csv")
	if *generateUnboundMethods != true {
		t.Errorf("flag generate_unbound_methods was not set correctly, wanted 'true' got '%v'", *generateUnboundMethods)
	}
	if *repeatedPathParamSeparator != "csv" {
		t.Errorf("flag importPrefix was not set correctly, wanted 'csv' got '%v'", *repeatedPathParamSeparator)
	}
}
