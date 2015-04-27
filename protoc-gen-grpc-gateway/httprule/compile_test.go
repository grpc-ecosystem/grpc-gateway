package httprule

import (
	"reflect"
	"testing"

	"github.com/gengo/grpc-gateway/internal"
)

const (
	operandFiller = 0
)

func TestCompile(t *testing.T) {
	for _, spec := range []struct {
		segs []segment
		verb string

		ops    []int
		pool   []string
		fields []string
	}{
		{},
		{
			segs: []segment{
				wildcard{},
			},
			ops: []int{int(internal.OpPush), operandFiller},
		},
		{
			segs: []segment{
				deepWildcard{},
			},
			ops: []int{int(internal.OpPushM), operandFiller},
		},
		{
			segs: []segment{
				literal("v1"),
			},
			ops:  []int{int(internal.OpLitPush), 0},
			pool: []string{"v1"},
		},
		{
			segs: []segment{
				literal("v1"),
			},
			verb: "LOCK",
			ops:  []int{int(internal.OpLitPush), 0},
			pool: []string{"v1"},
		},
		{
			segs: []segment{
				variable{
					path: "name.nested",
					segments: []segment{
						wildcard{},
					},
				},
			},
			ops: []int{
				int(internal.OpPush), operandFiller,
				int(internal.OpConcatN), 1,
				int(internal.OpCapture), 0,
			},
			pool:   []string{"name.nested"},
			fields: []string{"name.nested"},
		},
		{
			segs: []segment{
				literal("obj"),
				variable{
					path: "name.nested",
					segments: []segment{
						literal("a"),
						wildcard{},
						literal("b"),
					},
				},
				variable{
					path: "obj",
					segments: []segment{
						deepWildcard{},
					},
				},
			},
			ops: []int{
				int(internal.OpLitPush), 0,
				int(internal.OpLitPush), 1,
				int(internal.OpPush), operandFiller,
				int(internal.OpLitPush), 2,
				int(internal.OpConcatN), 3,
				int(internal.OpCapture), 3,
				int(internal.OpPushM), operandFiller,
				int(internal.OpConcatN), 1,
				int(internal.OpCapture), 0,
			},
			pool:   []string{"obj", "a", "b", "name.nested"},
			fields: []string{"name.nested", "obj"},
		},
	} {
		tmpl := template{
			segments: spec.segs,
			verb:     spec.verb,
		}
		compiled := tmpl.Compile()
		if got, want := compiled.Version, opcodeVersion; got != want {
			t.Errorf("tmpl.Compile().Version = %d; want %d; segs=%#v, verb=%q", got, want, spec.segs, spec.verb)
		}
		if got, want := compiled.OpCodes, spec.ops; !reflect.DeepEqual(got, want) {
			t.Errorf("tmpl.Compile().OpCodes = %v; want %v; segs=%#v, verb=%q", got, want, spec.segs, spec.verb)
		}
		if got, want := compiled.Pool, spec.pool; !reflect.DeepEqual(got, want) {
			t.Errorf("tmpl.Compile().Pool = %q; want %q; segs=%#v, verb=%q", got, want, spec.segs, spec.verb)
		}
		if got, want := compiled.Verb, spec.verb; got != want {
			t.Errorf("tmpl.Compile().Verb = %q; want %q; segs=%#v, verb=%q", got, want, spec.segs, spec.verb)
		}
		if got, want := compiled.Fields, spec.fields; !reflect.DeepEqual(got, want) {
			t.Errorf("tmpl.Compile().Fields = %q; want %q; segs=%#v, verb=%q", got, want, spec.segs, spec.verb)
		}
	}
}
