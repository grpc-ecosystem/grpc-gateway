package genopenapi

import (
	"fmt"
	"regexp"
	"strings"
)

// pathParamPattern matches a single proto path parameter expression with an
// optional `=constraint` clause, e.g. {field} or {field=lit/*/lit/*}.
var pathParamPattern = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_.]*)(?:=([^}]*))?\}`)

// pathParam describes one OpenAPI path parameter synthesized from the
// expansion of a proto URL template.
//
// fieldName is the proto request field this parameter ultimately populates
// (e.g. "name"). openAPIName is the unique name as it appears in the OpenAPI
// URL — usually equal to fieldName, but with a "_<n>" suffix when one proto
// field expands into multiple wildcards.
type pathParam struct {
	openAPIName string
	fieldName   string
}

// convertPathTemplate translates a proto HTTP rule URL template to an OpenAPI
// URL template plus the list of synthetic OpenAPI path parameters in URL
// order.
//
// The proto rule {name=shelves/*} matches URLs of the form /shelves/<id> and
// sets the proto field `name` to the full matched substring (e.g.
// "shelves/abc"). The literal portion is part of the URL the client must
// send. This expansion emits the literal verbatim and synthesizes one OpenAPI parameter
// per wildcard:
//
//	/v1/{name}                   → /v1/{name}
//	/v1/{name=shelves/*}         → /v1/shelves/{name}
//	/v1/{name=shelves/*/books/*} → /v1/shelves/{name}/books/{name_1}
//	/v1/{name=files/**}          → /v1/files/{name}
//
// Synthetic OpenAPI parameter names do not need to match any single proto
// field — they exist only so the generated URL is valid OpenAPI 3.1. The
// grpc-gateway runtime reconstructs the original proto value from the full
// matched substring.
func convertPathTemplate(template string) (string, []pathParam) {
	var params []pathParam
	out := pathParamPattern.ReplaceAllStringFunc(template, func(match string) string {
		sub := pathParamPattern.FindStringSubmatch(match)
		field, constraint := sub[1], sub[2]
		if constraint == "" {
			params = append(params, pathParam{openAPIName: field, fieldName: field})
			return "{" + field + "}"
		}
		var b strings.Builder
		wildIdx := 0
		for i, segment := range strings.Split(constraint, "/") {
			if i > 0 {
				b.WriteByte('/')
			}
			if segment != "*" && segment != "**" {
				b.WriteString(segment)
				continue
			}
			name := field
			if wildIdx > 0 {
				name = fmt.Sprintf("%s_%d", field, wildIdx)
			}
			b.WriteString("{" + name + "}")
			params = append(params, pathParam{openAPIName: name, fieldName: field})
			wildIdx++
		}
		return b.String()
	})
	return out, params
}
