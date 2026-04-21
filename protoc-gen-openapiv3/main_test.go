package main

import (
	"flag"
	"testing"
)

func TestParseReqParam_DisableDefaultErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		param     string
		wantErr   bool
		wantValue bool
	}{
		{
			name:      "bare flag name implies true",
			param:     "disable_default_errors",
			wantValue: true,
		},
		{
			name:      "explicit true",
			param:     "disable_default_errors=true",
			wantValue: true,
		},
		{
			name:      "explicit false",
			param:     "disable_default_errors=false",
			wantValue: false,
		},
		{
			name:    "unknown flag returns error",
			param:   "unknown_flag",
			wantErr: true,
		},
		{
			name:  "empty param is a no-op",
			param: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := flag.NewFlagSet("test", flag.ContinueOnError)
			got := f.Bool("disable_default_errors", false, "")

			err := parseReqParam(tc.param, f)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if *got != tc.wantValue {
				t.Errorf("disable_default_errors = %v, want %v", *got, tc.wantValue)
			}
		})
	}
}
