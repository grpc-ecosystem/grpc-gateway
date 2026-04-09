package genopenapi

import (
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Proto path indices into the FileDescriptorProto for source code info.
// These mirror the field numbers in descriptor.proto and are how protoc
// records line ranges per declaration.
var (
	messageProtoPath  = protoFieldNumber[descriptorpb.FileDescriptorProto]("MessageType")
	enumProtoPath     = protoFieldNumber[descriptorpb.FileDescriptorProto]("EnumType")
	serviceProtoPath  = protoFieldNumber[descriptorpb.FileDescriptorProto]("Service")
	methodProtoPath   = protoFieldNumber[descriptorpb.ServiceDescriptorProto]("Method")
	nestedProtoPath   = protoFieldNumber[descriptorpb.DescriptorProto]("NestedType")
	fieldProtoPath    = protoFieldNumber[descriptorpb.DescriptorProto]("Field")
	enumTypeProtoPath = protoFieldNumber[descriptorpb.DescriptorProto]("EnumType")
)

// protoFieldNumber pulls the field number out of a descriptorpb struct's
// protobuf tag. We resolve this once at startup so a future descriptor.proto
// rename is caught immediately rather than silently producing wrong comment
// lookups.
func protoFieldNumber[T any](fieldName string) int32 {
	t := reflect.TypeFor[T]()
	field, ok := t.FieldByName(fieldName)
	if !ok {
		panic("genopenapi: missing field " + fieldName + " on " + t.Name())
	}
	// Tag format: "bytes,4,rep,name=message_type,..."; the second comma-
	// separated component is always the field number.
	parts := strings.Split(field.Tag.Get("protobuf"), ",")
	if len(parts) < 2 {
		panic("genopenapi: malformed protobuf tag on " + t.Name() + "." + fieldName)
	}
	n, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		panic("genopenapi: bad field number on " + t.Name() + "." + fieldName + ": " + err.Error())
	}
	return int32(n)
}

// internalCommentPattern matches AIP-192 internal comment markers (-- ... --).
var internalCommentPattern = regexp.MustCompile(`(?s)\(--.*?--\)`)

// extractComments returns the leading comment for the given proto path within
// a file. Returns an empty string if no comment is present or the file has no
// source-code info.
func extractComments(file *descriptor.File, path []int32) string {
	if file.SourceCodeInfo == nil {
		return ""
	}
	for _, loc := range file.SourceCodeInfo.Location {
		if !slices.Equal(loc.Path, path) {
			continue
		}
		var comment string
		if loc.LeadingComments != nil {
			comment = strings.TrimRight(*loc.LeadingComments, "\n")
			// Strip a single leading space from continuation lines. This
			// handles the common protoc output shape but not tabs, deeper
			// indentation, or doxygen-style "\n *" prefixes; comments using
			// those will render with ragged margins.
			comment = strings.ReplaceAll(comment, "\n ", "\n")
		}
		comment = strings.TrimSpace(comment)
		comment = internalCommentPattern.ReplaceAllString(comment, "")
		comment = strings.TrimSpace(comment)
		return comment
	}
	return ""
}

// splitSummaryDescription splits a comment block into a one-line summary
// (the first paragraph) and a longer description (the remainder), following
// the convention used by Google AIP and the v2 generator.
func splitSummaryDescription(comment string) (summary, description string) {
	if comment == "" {
		return "", ""
	}
	parts := strings.SplitN(comment, "\n\n", 2)
	summary = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		description = strings.TrimSpace(parts[1])
	}
	return summary, description
}

// serviceComments returns the leading comment on a service declaration.
func serviceComments(svc *descriptor.Service) string {
	idx := slices.Index(svc.File.Services, svc)
	if idx < 0 {
		return ""
	}
	return extractComments(svc.File, []int32{serviceProtoPath, int32(idx)})
}

// methodComments returns the leading comment on a method declaration.
func methodComments(m *descriptor.Method) string {
	svcIdx := slices.Index(m.Service.File.Services, m.Service)
	if svcIdx < 0 {
		return ""
	}
	mIdx := slices.Index(m.Service.Methods, m)
	if mIdx < 0 {
		return ""
	}
	return extractComments(m.Service.File, []int32{
		serviceProtoPath, int32(svcIdx),
		methodProtoPath, int32(mIdx),
	})
}

// messageComments returns the leading comment on a message declaration,
// walking outer messages so nested types resolve correctly.
func messageComments(reg *descriptor.Registry, msg *descriptor.Message) string {
	path := messagePath(reg, msg)
	if path == nil {
		return ""
	}
	return extractComments(msg.File, path)
}

// fieldComments returns the leading comment on a field declaration.
func fieldComments(reg *descriptor.Registry, field *descriptor.Field) string {
	parentPath := messagePath(reg, field.Message)
	if parentPath == nil {
		return ""
	}
	idx := slices.Index(field.Message.Fields, field)
	if idx < 0 {
		return ""
	}
	return extractComments(field.Message.File, append(parentPath, fieldProtoPath, int32(idx)))
}

// enumComments returns the leading comment on an enum declaration.
func enumComments(reg *descriptor.Registry, enum *descriptor.Enum) string {
	path := enumPath(reg, enum)
	if path == nil {
		return ""
	}
	return extractComments(enum.File, path)
}

// messagePath builds the source-code-info path to a (possibly nested) message.
func messagePath(reg *descriptor.Registry, msg *descriptor.Message) []int32 {
	path := []int32{messageProtoPath}
	if len(msg.Outers) > 0 {
		outers := outerPath(reg, msg.File, msg.Outers)
		if outers == nil {
			return nil
		}
		path = append(path, outers...)
		path = append(path, nestedProtoPath)
	}
	path = append(path, int32(msg.Index))
	return path
}

// enumPath builds the source-code-info path to a (possibly nested) enum.
func enumPath(reg *descriptor.Registry, enum *descriptor.Enum) []int32 {
	if len(enum.Outers) == 0 {
		return []int32{enumProtoPath, int32(enum.Index)}
	}
	path := []int32{messageProtoPath}
	outers := outerPath(reg, enum.File, enum.Outers)
	if outers == nil {
		return nil
	}
	path = append(path, outers...)
	path = append(path, enumTypeProtoPath, int32(enum.Index))
	return path
}

// outerPath walks a chain of outer message names and returns their indices,
// separated by the nested-type field number, suitable for appending to a
// source-code-info path.
func outerPath(reg *descriptor.Registry, file *descriptor.File, outers []string) []int32 {
	var result []int32
	pkg := ""
	if file.Package != nil {
		pkg = file.GetPackage()
	}
	for i := range outers {
		fqn := strings.Join(outers[:i+1], ".")
		msg, err := reg.LookupMsg(pkg, fqn)
		if err != nil {
			return nil
		}
		if i > 0 {
			result = append(result, nestedProtoPath)
		}
		result = append(result, int32(msg.Index))
	}
	return result
}

