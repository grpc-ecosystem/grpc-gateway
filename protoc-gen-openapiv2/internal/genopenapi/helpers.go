//go:build go1.12
// +build go1.12

package genopenapi

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func fieldName(k string) string {
	return strings.ReplaceAll(cases.Title(language.AmericanEnglish).String(k), "-", "_")
}
