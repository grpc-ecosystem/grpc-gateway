package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	field_mask "google.golang.org/protobuf/types/known/fieldmaskpb"
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
		if errors.Is(err, io.EOF) {
			return fm, nil
		}
		return nil, err
	}

	queue := []fieldMaskPathItem{{node: root, msg: msg.ProtoReflect()}}
	for len(queue) > 0 {
		// dequeue an item
		item := queue[0]
		queue = queue[1:]

		m, ok := item.node.(map[string]interface{})
		switch {
		case ok && len(m) > 0:
			// if the item is an object, then enqueue all of its children
			for k, v := range m {
				if item.msg == nil {
					return nil, errors.New("JSON structure did not match request type")
				}

				fd := getFieldByName(item.msg.Descriptor().Fields(), k)
				if fd == nil {
					return nil, fmt.Errorf("could not find field %q in %q", k, item.msg.Descriptor().FullName())
				}

				if isDynamicProtoMessage(fd.Message()) {
					for _, p := range buildPathsBlindly(string(fd.FullName().Name()), v) {
						newPath := p
						if item.path != "" {
							newPath = item.path + "." + newPath
						}
						queue = append(queue, fieldMaskPathItem{path: newPath})
					}
					continue
				}

				if isProtobufAnyMessage(fd.Message()) && !fd.IsList() {
					vm, _ := v.(map[string]interface{})
					if _, hasTypeField := vm["@type"]; !hasTypeField {
						return nil, fmt.Errorf("could not find field @type in %q in message %q", k, item.msg.Descriptor().FullName())
					}
					newPath := string(fd.FullName().Name())
					if item.path != "" {
						newPath = item.path + "." + newPath
					}
					queue = append(queue, fieldMaskPathItem{path: newPath})
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
					// As per: https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/field_mask.proto#L85-L86
					// Do not recurse into repeated fields. The repeated field goes on the end of the path and we stop.
					fm.Paths = append(fm.Paths, child.path)
				case fd.Message() != nil:
					child.msg = item.msg.Get(fd).Message()
					fallthrough
				default:
					queue = append(queue, child)
				}
			}
		case ok && len(m) == 0:
			fallthrough
		case len(item.path) > 0:
			// otherwise, it's a leaf node so print its path
			fm.Paths = append(fm.Paths, item.path)
		}
	}

	// Sort for deterministic output in the presence
	// of repeated fields.
	sort.Strings(fm.Paths)

	return fm, nil
}

func isProtobufAnyMessage(md protoreflect.MessageDescriptor) bool {
	return md != nil && (md.FullName() == "google.protobuf.Any")
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
	if !ok || len(m) == 0 {
		return []string{name}
	}

	// A google.protobuf.Struct/Value field carries arbitrary client JSON, so
	// its nesting depth is attacker controlled. Keep each queue entry pointing
	// at its parent and its own segment rather than carrying the full path
	// string, so a path is concatenated once at the leaf (O(depth)) instead of
	// rebuilding the growing prefix at every level (O(depth^2)). A deeply
	// nested object otherwise makes this allocate quadratically in the depth
	// while producing only a handful of paths.
	type pathNode struct {
		parent  *pathNode
		segment string
	}
	type queueItem struct {
		path *pathNode
		node map[string]interface{}
	}
	// joinPath reconstructs a dotted path by walking from a leaf back to the
	// root, so each path is materialized exactly once.
	joinPath := func(leaf *pathNode) string {
		var segments []string
		n := 0
		for p := leaf; p != nil; p = p.parent {
			segments = append(segments, p.segment)
			if p != leaf {
				n++ // separating dot
			}
			n += len(p.segment)
		}
		var b strings.Builder
		b.Grow(n)
		for i := len(segments) - 1; i >= 0; i-- {
			if i != len(segments)-1 {
				b.WriteByte('.')
			}
			b.WriteString(segments[i])
		}
		return b.String()
	}

	var paths []string
	queue := []queueItem{{path: &pathNode{segment: name}, node: m}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for k, v := range cur.node {
			child := &pathNode{parent: cur.path, segment: k}
			if mi, ok := v.(map[string]interface{}); ok && len(mi) > 0 {
				queue = append(queue, queueItem{path: child, node: mi})
			} else {
				// No populated object remains, so there are no more levels to descend.
				paths = append(paths, joinPath(child))
			}
		}
	}
	return paths
}

// fieldMaskPathItem stores an in-progress deconstruction of a path for a fieldmask
type fieldMaskPathItem struct {
	// the list of prior fields leading up to node connected by dots
	path string

	// a generic decoded json object the current item to inspect for further path extraction
	node interface{}

	// parent message
	msg protoreflect.Message
}
