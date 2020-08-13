package runtime

import (
	"encoding/json"
	"fmt"
	"io"

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

				if isDynamicProtoMessage(fd.Message()) {
					for _, p := range buildPathsBlindly(k, v) {
						newPath := p
						if item.path != "" {
							newPath = item.path + "." + newPath
						}
						queue = append(queue, fieldMaskPathItem{path: newPath})
					}
					continue
				}

				child := fieldMaskPathItem{
					node: v,
				}
				if item.path == "" {
					child.path = string(fd.FullName().Name())
				} else {
					child.path = item.path + "." + string(fd.FullName().Name())
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
			fm.Paths = append(fm.Paths, item.path)
		}
	}

	// Add any repeated fields last, as per
	// https://github.com/protocolbuffers/protobuf/blob/6b0ff74ecf63e26c7315f6745de36aff66deb59d/src/google/protobuf/field_mask.proto#L85-L86
	if repeatedChild != nil {
		fm.Paths = append(fm.Paths, repeatedChild.path)
	}

	return fm, nil
}

func isDynamicProtoMessage(md protoreflect.MessageDescriptor) bool {
	return md != nil && (md.FullName() == "google.protobuf.Struct" || md.FullName() == "google.protobuf.Value")
}

// buildPathsBlindly does not attempt to match proto field names to the
// json value keys.  Instead it relies completely on the structure of
// the unmarshalled json contained within in.
// Returns a slice containing all subpaths with the root at the
// passed in name and json value.
func buildPathsBlindly(name string, in interface{}) []string {
	m, ok := in.(map[string]interface{})
	if !ok {
		return []string{name}
	}

	var paths []string
	queue := []fieldMaskPathItem{{path: name, node: m}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		m, ok := cur.node.(map[string]interface{})
		if !ok {
			// This should never happen since we should always check that we only add
			// nodes of type map[string]interface{} to the queue.
			continue
		}
		for k, v := range m {
			if mi, ok := v.(map[string]interface{}); ok {
				queue = append(queue, fieldMaskPathItem{path: cur.path + "." + k, node: mi})
			} else {
				// This is not a struct, so there are no more levels to descend.
				curPath := cur.path + "." + k
				paths = append(paths, curPath)
			}
		}
	}
	return paths
}

// fieldMaskPathItem stores a in-progress deconstruction of a path for a fieldmask
type fieldMaskPathItem struct {
	// the list of prior fields leading up to node connected by dots
	path string

	// a generic decoded json object the current item to inspect for further path extraction
	node interface{}

	// parent message
	msg protoreflect.Message
}
