package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func getFieldByName(fields protoreflect.FieldDescriptors, name string) protoreflect.FieldDescriptor {
	fd := fields.ByName(protoreflect.Name(name))
	if fd != nil {
		return fd
	}

	return fields.ByJSONName(name)
}

// FieldMaskFromRequestBody creates a FieldMask printing all complete paths from the JSON body.
func FieldMaskFromRequestBody(r io.Reader, msg proto.Message) (*field_mask.FieldMask, error) {
	fm := &field_mask.FieldMask{}
	var root interface{}

	if err := json.NewDecoder(r).Decode(&root); err != nil {
		if err == io.EOF {
			return fm, nil
		}
		return nil, err
	}

	queue := []fieldMaskPathItem{{node: root, msg: msg.ProtoReflect()}}
	var repeatedChild *fieldMaskPathItem
	for len(queue) > 0 {
		// dequeue an item
		item := queue[0]
		queue = queue[1:]

		m, ok := item.node.(map[string]interface{})
		switch {
		case ok:
			// if the item is an object, then enqueue all of its children
			for k, v := range m {
				if item.msg == nil {
					return nil, fmt.Errorf("JSON structure did not match request type")
				}

				fd := getFieldByName(item.msg.Descriptor().Fields(), k)
				if fd == nil {
					return nil, fmt.Errorf("could not find field %q in %q", k, item.msg.Descriptor().FullName())
				}
				child := fieldMaskPathItem{
					path: append(item.path, string(fd.FullName().Name())),
					node: v,
				}
				switch {
				case fd.IsList(), fd.IsMap():
					if repeatedChild != nil {
						// This is implied by the rule that any repeated fields must be
						// last in the paths.
						// Ref: https://github.com/protocolbuffers/protobuf/blob/6b0ff74ecf63e26c7315f6745de36aff66deb59d/src/google/protobuf/field_mask.proto#L85-L86
						return nil, fmt.Errorf("only one repeated value is allowed per field_mask")
					}
					repeatedChild = &child
					// Don't add to paths until the end
				case fd.Message() != nil:
					child.msg = item.msg.Get(fd).Message()
					fallthrough
				default:
					queue = append(queue, child)
				}
			}
		case len(item.path) > 0:
			// otherwise, it's a leaf node so print its path
			fm.Paths = append(fm.Paths, strings.Join(item.path, "."))
		}
	}

	// Add any repeated fields last, as per
	// https://github.com/protocolbuffers/protobuf/blob/6b0ff74ecf63e26c7315f6745de36aff66deb59d/src/google/protobuf/field_mask.proto#L85-L86
	if repeatedChild != nil {
		fm.Paths = append(fm.Paths, strings.Join(repeatedChild.path, "."))
	}

	return fm, nil
}

// fieldMaskPathItem stores a in-progress deconstruction of a path for a fieldmask
type fieldMaskPathItem struct {
	// the list of prior fields leading up to node
	path []string

	// a generic decoded json object the current item to inspect for further path extraction
	node interface{}

	// parent message
	msg protoreflect.Message
}
