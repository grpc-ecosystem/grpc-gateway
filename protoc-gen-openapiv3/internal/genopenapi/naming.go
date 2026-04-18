package genopenapi

import "strings"

// schemaName returns the OpenAPI component name for a fully-qualified proto
// name. We use a single naming strategy: strip the leading dot.
//
// Examples:
//
//	".foo.bar.Baz"      -> "foo.bar.Baz"
//	".foo.bar.Baz.Qux"  -> "foo.bar.Baz.Qux" (nested)
func schemaName(fqn string) string {
	return strings.TrimPrefix(fqn, ".")
}
