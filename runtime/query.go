package runtime

import (
	"net/url"
	"strings"

	"github.com/ajg/form"
	"github.com/gengo/grpc-gateway/internal"
	"github.com/golang/protobuf/proto"
)

func isQueryParam(key string, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(key, f) {
			switch l, m := len(key), len(f); {
			case l == m:
				return false
			case key[m] == '.':
				return false
			}
		}
	}
	return true
}

func convertPath(path string) string {
	var components []string
	for _, c := range strings.Split(path, ".") {
		components = append(components, internal.PascalFromSnake(c))
	}
	return strings.Join(components, ".")
}

// PopulateQueryParameters populates "values" into "msg".
// A value is ignored if its key starts with one of the elements in "filters".
//
// TODO(yugui) Use trie for filters?
func PopulateQueryParameters(msg proto.Message, values url.Values, filters []string) error {
	filtered := make(url.Values)
	for key, values := range values {
		if isQueryParam(key, filters) {
			filtered[convertPath(key)] = values
		}
	}
	return form.DecodeValues(msg, filtered)
}
