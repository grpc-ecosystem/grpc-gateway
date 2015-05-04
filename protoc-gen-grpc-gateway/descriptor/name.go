package descriptor

import (
	"regexp"
	"strings"
)

var (
	upperPattern = regexp.MustCompile("[A-Z]")
)

func toCamel(str string) string {
	var components []string
	for _, c := range strings.Split(str, "_") {
		components = append(components, strings.Title(strings.ToLower(c)))
	}
	return strings.Join(components, "")
}
