package runtime

import (
	"regexp"
	"testing"
)

// TestSplitMapKeyMatchesRegexp pins splitMapKey against the regexp it replaced.
//
// The pattern to reproduce is `^(.*)\[(.*)\]$` applied with FindStringSubmatch, and the
// interesting cases are not the obvious ones: the greedy first group picks the last bracket,
// and Go's regexp refuses to match anything containing a newline because "." does not cross
// one. Both are easy to get wrong by hand, so the pattern is kept here and every case is
// checked against it rather than against hand-written expectations.
func TestSplitMapKeyMatchesRegexp(t *testing.T) {
	legacy := regexp.MustCompile(`^(.*)\[(.*)\]$`)

	for _, key := range []string{
		"",
		"a",
		"a]",
		"[",
		"]",
		"[]",
		"a[]",
		"a[b]",
		"a.b[c]",
		"a[b][c]",
		"a[b]c]",
		"a[b.c]",
		"[a]b",
		"a[b",
		"map_value[key]",
		"deeply.nested.map[a.b.c]",
		"a\n[b]",
		"a[b\nc]",
		"\n[]",
		"a[]]",
		"a[[]",
	} {
		prefix, inner, ok := splitMapKey(key)
		match := legacy.FindStringSubmatch(key)
		wantOK := len(match) == 3

		if ok != wantOK {
			t.Errorf("key %q: splitMapKey ok = %v, regexp matched = %v", key, ok, wantOK)
			continue
		}
		if !wantOK {
			continue
		}
		if prefix != match[1] || inner != match[2] {
			t.Errorf("key %q: splitMapKey = (%q, %q), regexp = (%q, %q)",
				key, prefix, inner, match[1], match[2])
		}
		// Parse rebuilds the bracketed part as key[len(prefix):]; it has to equal what the
		// old code built with "[" + match[2] + "]".
		if got, want := key[len(prefix):], "["+match[2]+"]"; got != want {
			t.Errorf("key %q: rebuilt map key = %q, want %q", key, got, want)
		}
	}
}
