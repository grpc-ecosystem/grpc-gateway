package genopenapi

import (
	"slices"
	"testing"
)

func TestConvertPathTemplate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		in         string
		wantPath   string
		wantParams []pathParam
	}{
		{
			name:       "plain",
			in:         "/v1/books/{id}",
			wantPath:   "/v1/books/{id}",
			wantParams: []pathParam{{openAPIName: "id", fieldName: "id"}},
		},
		{
			name:       "dotted",
			in:         "/v1/books/{book.id}",
			wantPath:   "/v1/books/{book.id}",
			wantParams: []pathParam{{openAPIName: "book.id", fieldName: "book.id"}},
		},
		{
			name:       "single_constraint",
			in:         "/v1/{name=shelves/*}",
			wantPath:   "/v1/shelves/{name}",
			wantParams: []pathParam{{openAPIName: "name", fieldName: "name"}},
		},
		{
			name:     "multi_segment_constraint",
			in:       "/v1/{name=shelves/*/books/*}",
			wantPath: "/v1/shelves/{name}/books/{name_1}",
			wantParams: []pathParam{
				{openAPIName: "name", fieldName: "name"},
				{openAPIName: "name_1", fieldName: "name"},
			},
		},
		{
			name:     "multiple_params",
			in:       "/v1/{shelf}/books/{book}",
			wantPath: "/v1/{shelf}/books/{book}",
			wantParams: []pathParam{
				{openAPIName: "shelf", fieldName: "shelf"},
				{openAPIName: "book", fieldName: "book"},
			},
		},
		{
			name:     "mixed_constraint_and_plain",
			in:       "/v1/{shelf=shelves/*}/books/{book}",
			wantPath: "/v1/shelves/{shelf}/books/{book}",
			wantParams: []pathParam{
				{openAPIName: "shelf", fieldName: "shelf"},
				{openAPIName: "book", fieldName: "book"},
			},
		},
		{
			name:       "double_wildcard",
			in:         "/v1/{name=files/**}",
			wantPath:   "/v1/files/{name}",
			wantParams: []pathParam{{openAPIName: "name", fieldName: "name"}},
		},
		{
			name:       "no_params",
			in:         "/v1/healthz",
			wantPath:   "/v1/healthz",
			wantParams: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotPath, gotParams := convertPathTemplate(tc.in)
			if gotPath != tc.wantPath {
				t.Errorf("convertPathTemplate(%q) path = %q, want %q", tc.in, gotPath, tc.wantPath)
			}
			if !slices.Equal(gotParams, tc.wantParams) {
				t.Errorf("convertPathTemplate(%q) params = %+v, want %+v", tc.in, gotParams, tc.wantParams)
			}
		})
	}
}
