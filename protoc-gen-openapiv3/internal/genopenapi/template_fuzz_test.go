//go:build go1.18
// +build go1.18

package genopenapi

import (
	"regexp"
	"testing"
)

var replaceInternalCommentsRegex = regexp.MustCompile(`(?s)(\r?\n)?[ \t]*(\(--)((.*?--\))|.*$)?`)

func FuzzRemoveInternalComments(f *testing.F) {
	f.Add("Text\n\n(-- Comment --)\n\nMore Text\n")
	f.Add("Text\n\n(-- Multi\nLine\n\nComment --)\n\nMore Text\n")
	f.Add("(-- Starting with comment --)\n\nMore Text\n")
	f.Add("\n\n(-- Starting with new line and comment --)\n\nMore Text\n")
	f.Add("Ending with\n\n(-- Comment --)")
	f.Fuzz(func(t *testing.T, s string) {
		s1 := removeInternalComments(s)
		s2 := replaceInternalCommentsRegex.ReplaceAllString(s, "")
		if s1 != s2 {
			t.Errorf("Unexpected comment removal difference: our function produced %q but regex produced %q on %q", s1, s2, s)
		}
	})
}
