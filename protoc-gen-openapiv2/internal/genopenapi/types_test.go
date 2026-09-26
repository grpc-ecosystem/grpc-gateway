package genopenapi

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestOpenapiSchemaAdditionalPropertiesEncoding(t *testing.T) {
	tests := []struct {
		name       string
		value      openapiAdditionalProperties
		want       interface{}
		wantExists bool
	}{
		{name: "omitted"},
		{name: "allowed", value: openapiAdditionalPropertiesAllowed(true), want: true, wantExists: true},
		{name: "disallowed", value: openapiAdditionalPropertiesAllowed(false), want: false, wantExists: true},
		{name: "schema", value: &openapiSchemaObject{schemaCore: schemaCore{Type: "string"}}, want: map[string]interface{}{"type": "string"}, wantExists: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := openapiSchemaObject{AdditionalProperties: tt.value}
			for _, format := range []string{"json", "yaml"} {
				t.Run(format, func(t *testing.T) {
					var data []byte
					var err error
					if format == "json" {
						data, err = json.Marshal(schema)
					} else {
						data, err = yaml.Marshal(schema)
					}
					if err != nil {
						t.Fatalf("marshal schema: %v", err)
					}
					var document map[string]interface{}
					if format == "json" {
						err = json.Unmarshal(data, &document)
					} else {
						err = yaml.Unmarshal(data, &document)
					}
					if err != nil {
						t.Fatalf("unmarshal schema: %v", err)
					}
					got, exists := document["additionalProperties"]
					if exists != tt.wantExists || !reflect.DeepEqual(got, tt.want) {
						t.Errorf("additionalProperties = %#v (exists %t), want %#v (exists %t)", got, exists, tt.want, tt.wantExists)
					}
				})
			}
		})
	}
}

func newSpaceReplacer() *strings.Replacer {
	return strings.NewReplacer(" ", "", "\n", "", "\t", "")
}

func TestRawExample(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		In  RawExample
		Exp string
	}{{
		In:  RawExample(`1`),
		Exp: `1`,
	}, {
		In:  RawExample(`"1"`),
		Exp: `"1"`,
	}, {
		In: RawExample(`{"hello":"worldr"}`),
		Exp: `
			hello:
				worldr
		`,
	}}

	sr := newSpaceReplacer()

	for _, tc := range testCases {
		tc := tc

		t.Run(string(tc.In), func(t *testing.T) {
			t.Parallel()

			ex := tc.In

			out, err := yaml.Marshal(ex)
			switch {
			case err != nil:
				t.Fatalf("expect no yaml marshal error, got: %s", err)
			case !json.Valid(tc.In):
				t.Fatalf("json is invalid: %#q", tc.In)
			case sr.Replace(tc.Exp) != sr.Replace(string(out)):
				t.Fatalf("expected: %s, actual: %s", tc.Exp, out)
			}

			out, err = json.Marshal(tc.In)
			switch {
			case err != nil:
				t.Fatalf("expect no json marshal error, got: %s", err)
			case sr.Replace(string(tc.In)) != sr.Replace(string(out)):
				t.Fatalf("expected: %s, actual: %s", tc.In, out)
			}
		})
	}
}

func TestOpenapiSchemaObjectProperties(t *testing.T) {
	t.Parallel()

	v := map[string]interface{}{
		"example": openapiSchemaObjectProperties{{
			Key:   "test1",
			Value: 1,
		}, {
			Key:   "test2",
			Value: 2,
		}},
	}

	t.Run("yaml", func(t *testing.T) {
		t.Parallel()

		const exp = `
			example:
				test1: 1
				test2: 2
			`

		sr := newSpaceReplacer()

		out, err := yaml.Marshal(v)
		switch {
		case err != nil:
			t.Fatalf("expect no marshal error, got: %s", err)
		case sr.Replace(exp) != sr.Replace(string(out)):
			t.Fatalf("expected: %s, actual: %s", exp, out)
		}
	})

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		const exp = `{"example":{"test1":1,"test2":2}}`

		got, err := json.Marshal(v)
		switch {
		case err != nil:
			t.Fatalf("expect no marshal error, got: %s", err)
		case exp != string(got):
			t.Fatalf("expected: %s, actual: %s", exp, got)
		}
	})
}
