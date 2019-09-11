package runtime

import (
	"bytes"
	"fmt"
	"testing"

	"google.golang.org/genproto/protobuf/field_mask"
)

func fieldMasksEqual(fm1, fm2 *field_mask.FieldMask) bool {
	if fm1 == nil && fm2 == nil {
		return true
	}
	if fm1 == nil || fm2 == nil {
		return false
	}
	if len(fm1.GetPaths()) != len(fm2.GetPaths()) {
		return false
	}

	paths := make(map[string]bool)
	for _, path := range fm1.GetPaths() {
		paths[path] = true
	}
	for _, path := range fm2.GetPaths() {
		if _, ok := paths[path]; !ok {
			return false
		}
	}

	return true
}

func newFieldMask(paths ...string) *field_mask.FieldMask {
	return &field_mask.FieldMask{Paths: paths}
}

func fieldMaskString(fm *field_mask.FieldMask) string {
	if fm == nil {
		return ""
	}
	return fmt.Sprintf("%v", fm.GetPaths())
}

func TestFieldMaskFromRequestBody(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input       string
		expected    *field_mask.FieldMask
		expectedErr error
	}{
		{name: "empty", expected: newFieldMask()},
		{name: "simple", input: `{"foo":1, "bar":"baz"}`, expected: newFieldMask("foo", "bar")},
		{name: "nested", input: `{"foo": {"bar":1, "baz": 2}, "qux": 3}`, expected: newFieldMask("foo.bar", "foo.baz", "qux")},
		{name: "canonical", input: `{"f": {"b": {"d": 1, "x": 2}, "c": 1}}`, expected: newFieldMask("f.b.d", "f.b.x", "f.c")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := FieldMaskFromRequestBody(bytes.NewReader([]byte(tc.input)))
			if !fieldMasksEqual(actual, tc.expected) {
				t.Errorf("want %v; got %v", fieldMaskString(tc.expected), fieldMaskString(actual))
			}
			if err != tc.expectedErr {
				t.Errorf("want %v; got %v", tc.expectedErr, err)
			}
		})
	}
}
