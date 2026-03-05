package genopenapiv3

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Proto path indices for various descriptor types.
var (
	messageProtoPath = protoPathIndex(reflect.TypeOf((*descriptorpb.FileDescriptorProto)(nil)), "MessageType")
	nestedProtoPath  = protoPathIndex(reflect.TypeOf((*descriptorpb.DescriptorProto)(nil)), "NestedType")
	packageProtoPath = protoPathIndex(reflect.TypeOf((*descriptorpb.FileDescriptorProto)(nil)), "Package")
	serviceProtoPath = protoPathIndex(reflect.TypeOf((*descriptorpb.FileDescriptorProto)(nil)), "Service")
	methodProtoPath  = protoPathIndex(reflect.TypeOf((*descriptorpb.ServiceDescriptorProto)(nil)), "Method")
	enumProtoPath    = protoPathIndex(reflect.TypeOf((*descriptorpb.FileDescriptorProto)(nil)), "EnumType")
)

// protoPathIndex returns the proto path index for a field in a protobuf message type.
// This is used to build source code info paths for comment extraction.
func protoPathIndex(target reflect.Type, fieldName string) int32 {
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	field, ok := target.FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("field %q not found in %s", fieldName, target.Name()))
	}

	// Parse the protobuf tag to get the field number
	protoTag := field.Tag.Get("protobuf")
	if protoTag == "" {
		panic(fmt.Sprintf("field %q in %s has no protobuf tag", fieldName, target.Name()))
	}

	// Tag format: "bytes,4,rep,name=message_type,json=messageType,proto3"
	// The field number is always the second element in the tag
	parts := strings.Split(protoTag, ",")
	if len(parts) >= 2 {
		var fieldNum int32
		if _, err := fmt.Sscanf(parts[1], "%d", &fieldNum); err == nil {
			return fieldNum
		}
	}

	panic(fmt.Sprintf("could not find field number for %q in %s", fieldName, target.Name()))
}

// serviceComments extracts comments for a service.
func serviceComments(reg *descriptor.Registry, svc *descriptor.Service) string {
	if reg.GetIgnoreComments() {
		return ""
	}

	file := svc.File
	if file.SourceCodeInfo == nil {
		return ""
	}

	// Find service index
	svcIndex := int32(-1)
	for i, s := range file.Services {
		if s == svc {
			svcIndex = int32(i)
			break
		}
	}
	if svcIndex < 0 {
		return ""
	}

	path := []int32{serviceProtoPath, svcIndex}
	return extractComments(file, path, reg.GetRemoveInternalComments())
}

// methodComments extracts comments for a method.
func methodComments(reg *descriptor.Registry, method *descriptor.Method) string {
	if reg.GetIgnoreComments() {
		return ""
	}

	svc := method.Service
	file := svc.File
	if file.SourceCodeInfo == nil {
		return ""
	}

	// Find service index
	svcIndex := int32(-1)
	for i, s := range file.Services {
		if s == svc {
			svcIndex = int32(i)
			break
		}
	}
	if svcIndex < 0 {
		return ""
	}

	// Find method index
	methodIndex := int32(-1)
	for i, m := range svc.Methods {
		if m == method {
			methodIndex = int32(i)
			break
		}
	}
	if methodIndex < 0 {
		return ""
	}

	path := []int32{serviceProtoPath, svcIndex, methodProtoPath, methodIndex}
	return extractComments(file, path, reg.GetRemoveInternalComments())
}

// messageComments extracts comments for a message.
func messageComments(reg *descriptor.Registry, msg *descriptor.Message) string {
	if reg.GetIgnoreComments() {
		return ""
	}

	file := msg.File
	if file.SourceCodeInfo == nil {
		return ""
	}

	// Build path through nested messages
	var path []int32
	path = append(path, messageProtoPath)

	// Handle nested messages
	if len(msg.Outers) > 0 {
		// Find the outermost message first
		outerPath := buildMessagePath(reg, file, msg.Outers)
		if outerPath == nil {
			return ""
		}
		path = append(path[:1], outerPath...)
		path = append(path, nestedProtoPath)
	}

	path = append(path, int32(msg.Index))

	return extractComments(file, path, reg.GetRemoveInternalComments())
}

// buildMessagePath builds the path to a nested message.
func buildMessagePath(reg *descriptor.Registry, file *descriptor.File, outers []string) []int32 {
	if len(outers) == 0 {
		return nil
	}

	var path []int32
	location := ""
	if file.Package != nil {
		location = file.GetPackage()
	}

	for i := range outers {
		outerFQMN := strings.Join(outers[:i+1], ".")
		msg, err := reg.LookupMsg(location, outerFQMN)
		if err != nil {
			return nil
		}
		if i > 0 {
			path = append(path, nestedProtoPath)
		}
		path = append(path, int32(msg.Index))
	}

	return path
}

// fieldComments extracts comments for a field.
func fieldComments(reg *descriptor.Registry, field *descriptor.Field) string {
	if reg.GetIgnoreComments() {
		return ""
	}

	msg := field.Message
	file := msg.File
	if file.SourceCodeInfo == nil {
		return ""
	}

	// Build path to the field
	fieldProtoPath := protoPathIndex(reflect.TypeOf((*descriptorpb.DescriptorProto)(nil)), "Field")

	var path []int32
	path = append(path, messageProtoPath)

	// Handle nested messages
	if len(msg.Outers) > 0 {
		outerPath := buildMessagePath(reg, file, msg.Outers)
		if outerPath == nil {
			return ""
		}
		path = append(path[:1], outerPath...)
		path = append(path, nestedProtoPath)
	}

	path = append(path, int32(msg.Index))

	// Find field index
	fieldIndex := int32(-1)
	for i, f := range msg.Fields {
		if f == field {
			fieldIndex = int32(i)
			break
		}
	}
	if fieldIndex < 0 {
		return ""
	}

	path = append(path, fieldProtoPath, fieldIndex)

	return extractComments(file, path, reg.GetRemoveInternalComments())
}

// enumComments extracts comments for an enum.
func enumComments(reg *descriptor.Registry, enum *descriptor.Enum) string {
	if reg.GetIgnoreComments() {
		return ""
	}

	file := enum.File
	if file.SourceCodeInfo == nil {
		return ""
	}

	var path []int32

	// Handle nested enums (inside messages)
	if len(enum.Outers) > 0 {
		path = append(path, messageProtoPath)
		outerPath := buildMessagePath(reg, file, enum.Outers)
		if outerPath == nil {
			return ""
		}
		path = append(path, outerPath...)
		enumTypeProtoPath := protoPathIndex(reflect.TypeOf((*descriptorpb.DescriptorProto)(nil)), "EnumType")
		path = append(path, enumTypeProtoPath)
	} else {
		path = append(path, enumProtoPath)
	}

	path = append(path, int32(enum.Index))

	return extractComments(file, path, reg.GetRemoveInternalComments())
}

// extractComments extracts leading and trailing comments at the given path.
func extractComments(file *descriptor.File, path []int32, removeInternal bool) string {
	if file.SourceCodeInfo == nil {
		fmt.Fprintln(os.Stderr, file.GetName(), "descriptor.File should not contain nil SourceCodeInfo")
		return ""
	}

	for _, loc := range file.SourceCodeInfo.Location {
		if !pathMatches(loc.Path, path) {
			continue
		}

		comments := ""
		if loc.LeadingComments != nil {
			comments = strings.TrimRight(*loc.LeadingComments, "\n")
			comments = strings.TrimSpace(comments)
			// Fix "// " being interpreted as "//"
			comments = strings.ReplaceAll(comments, "\n ", "\n")
		}
		if loc.TrailingComments != nil {
			trailing := strings.TrimSpace(*loc.TrailingComments)
			if comments == "" {
				comments = trailing
			} else {
				comments += "\n\n" + trailing
			}
		}

		if removeInternal {
			comments = removeInternalComments(comments)
		}

		return comments
	}

	return ""
}

// pathMatches checks if two proto source code info paths match.
func pathMatches(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// internalCommentPattern matches internal comments: (-- ... --)
// Uses (?s) flag to make . match newlines for multi-line internal comments
var internalCommentPattern = regexp.MustCompile(`(?s)\(--.*?--\)`)

// bufLintIgnorePattern matches buf lint ignore directives: buf:lint:ignore RULE_NAME
var bufLintIgnorePattern = regexp.MustCompile(`(?m)^\s*buf:lint:ignore\s+\S+\s*$`)

// removeInternalComments removes internal comments per Google AIP-192.
// These are marked with (-- ... --) and should not appear in public docs.
// Also removes buf:lint:ignore directives which are tooling-specific.
func removeInternalComments(comment string) string {
	// Remove AIP-192 internal comments
	result := internalCommentPattern.ReplaceAllString(comment, "")
	// Remove buf lint ignore directives
	result = bufLintIgnorePattern.ReplaceAllString(result, "")
	// Clean up extra blank lines
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")
	return strings.TrimSpace(result)
}

// splitSummaryDescription splits comments into summary (first paragraph) and description (rest).
// Following Google AIP style guide conventions.
func splitSummaryDescription(comment string) (summary, description string) {
	if comment == "" {
		return "", ""
	}

	// Split on double newline (paragraph break)
	const paragraphDeliminator = "\n\n"
	parts := strings.SplitN(comment, paragraphDeliminator, 2)

	summary = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		description = strings.TrimSpace(parts[1])
	}

	return summary, description
}
