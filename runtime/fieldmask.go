package runtime

import (
	"encoding/json"
	"io"
	"strings"

	"google.golang.org/genproto/protobuf/field_mask"
)

// FieldMaskFromRequestBody creates a FieldMask printing all complete paths from the JSON body.
func FieldMaskFromRequestBody(r io.Reader) (*field_mask.FieldMask, error) {
	fm := &field_mask.FieldMask{}
	var root interface{}
	if err := json.NewDecoder(r).Decode(&root); err != nil {
		if err == io.EOF {
			return fm, nil
		}
		return nil, err
	}

	queue := []fieldMaskPathItem{{node: root}}
	for len(queue) > 0 {
		// dequeue an item
		item := queue[0]
		queue = queue[1:]

		if m, ok := item.node.(map[string]interface{}); ok {
			// if the item is an object, then enqueue all of its children
			for k, v := range m {
				queue = append(queue, fieldMaskPathItem{path: append(item.path, k), node: v})
			}
		} else if len(item.path) > 0 {
			// otherwise, it's a leaf node so print its path
			fm.Paths = append(fm.Paths, strings.Join(item.path, "."))
		}
	}

	return fm, nil
}

// fieldMaskPathItem stores a in-progress deconstruction of a path for a fieldmask
type fieldMaskPathItem struct {
	// the list of prior fields leading up to node
	path []string

	// a generic decoded json object the current item to inspect for further path extraction
	node interface{}
}
