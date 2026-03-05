package genopenapiv3

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/proto"
)

// getFileAnnotation retrieves the OpenAPI v3 file-level annotation if present.
func getFileAnnotation(file *descriptor.File) *options.OpenAPI {
	if file.Options == nil {
		return nil
	}
	ext := proto.GetExtension(file.Options, options.E_Openapiv3Openapi)
	opts, ok := ext.(*options.OpenAPI)
	if !ok || opts == nil {
		return nil
	}
	return opts
}

// getMethodAnnotation retrieves the OpenAPI v3 method-level annotation if present.
func getMethodAnnotation(method *descriptor.Method) *options.Operation {
	if method.Options == nil {
		return nil
	}
	ext := proto.GetExtension(method.Options, options.E_Openapiv3Operation)
	opts, ok := ext.(*options.Operation)
	if !ok || opts == nil {
		return nil
	}
	return opts
}

// getMessageAnnotation retrieves the OpenAPI v3 message-level annotation if present.
func getMessageAnnotation(msg *descriptor.Message) *options.Schema {
	if msg.Options == nil {
		return nil
	}
	ext := proto.GetExtension(msg.Options, options.E_Openapiv3Schema)
	opts, ok := ext.(*options.Schema)
	if !ok || opts == nil {
		return nil
	}
	return opts
}

// getFieldAnnotation retrieves the OpenAPI v3 field-level annotation if present.
func getFieldAnnotation(field *descriptor.Field) *options.Schema {
	if field.Options == nil {
		return nil
	}
	ext := proto.GetExtension(field.Options, options.E_Openapiv3Field)
	opts, ok := ext.(*options.Schema)
	if !ok || opts == nil {
		return nil
	}
	return opts
}

// getServiceAnnotation retrieves the OpenAPI v3 service-level annotation if present.
func getServiceAnnotation(svc *descriptor.Service) *options.Tag {
	if svc.Options == nil {
		return nil
	}
	ext := proto.GetExtension(svc.Options, options.E_Openapiv3Tag)
	opts, ok := ext.(*options.Tag)
	if !ok || opts == nil {
		return nil
	}
	return opts
}

// getEnumAnnotation retrieves the OpenAPI v3 enum-level annotation if present.
func getEnumAnnotation(enum *descriptor.Enum) *options.EnumSchema {
	if enum.Options == nil {
		return nil
	}
	ext := proto.GetExtension(enum.Options, options.E_Openapiv3Enum)
	opts, ok := ext.(*options.EnumSchema)
	if !ok || opts == nil {
		return nil
	}
	return opts
}
