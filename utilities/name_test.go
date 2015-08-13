package utilities_test

import (
	"testing"

	"github.com/gengo/grpc-gateway/utilities"
)

func TestPascalToSnake(t *testing.T) {
	for _, spec := range []struct {
		input, want string
	}{
		{input: "value", want: "Value"},
		{input: "prefixed_value", want: "PrefixedValue"},
		{input: "foo_id", want: "FooId"},
	} {
		got := utilities.PascalFromSnake(spec.input)
		if got != spec.want {
			t.Errorf("utilities.PascalFromSnake(%q) = %q; want %q", spec.input, got, spec.want)
		}
	}
}
