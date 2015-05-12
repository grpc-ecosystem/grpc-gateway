package internal

import (
	"strings"
)

// PascalFromSnake converts an identifier in snake_case into PascalCase.
func PascalFromSnake(str string) string {
	var components []string
	for _, c := range strings.Split(str, "_") {
		components = append(components, strings.Title(strings.ToLower(c)))
	}
	return strings.Join(components, "")
}
