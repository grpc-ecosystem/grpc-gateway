package genopenapiv3

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// getAnnotation is a generic function that retrieves OpenAPI v3 annotations from proto options.
// It reduces code duplication by handling the common pattern of:
// 1. Checking if options is nil
// 2. Getting the extension
// 3. Type asserting to the expected type
func getAnnotation[T any](opts proto.Message, ext protoreflect.ExtensionType) *T {
	if opts == nil {
		return nil
	}
	extVal := proto.GetExtension(opts, ext)
	result, ok := extVal.(*T)
	if !ok || result == nil {
		return nil
	}
	return result
}

// getFileAnnotation retrieves the OpenAPI v3 file-level annotation if present.
func getFileAnnotation(file *descriptor.File) *options.OpenAPI {
	if file.Options == nil {
		return nil
	}
	return getAnnotation[options.OpenAPI](file.Options, options.E_Openapiv3Document)
}

// getMethodAnnotation retrieves the OpenAPI v3 method-level annotation if present.
func getMethodAnnotation(method *descriptor.Method) *options.Operation {
	if method.Options == nil {
		return nil
	}
	return getAnnotation[options.Operation](method.Options, options.E_Openapiv3Operation)
}

// getMessageAnnotation retrieves the OpenAPI v3 message-level annotation if present.
func getMessageAnnotation(msg *descriptor.Message) *options.Schema {
	if msg.Options == nil {
		return nil
	}
	return getAnnotation[options.Schema](msg.Options, options.E_Openapiv3Schema)
}

// getFieldAnnotation retrieves the OpenAPI v3 field-level annotation if present.
func getFieldAnnotation(field *descriptor.Field) *options.Schema {
	if field.Options == nil {
		return nil
	}
	return getAnnotation[options.Schema](field.Options, options.E_Openapiv3Field)
}

// getServiceAnnotation retrieves the OpenAPI v3 service-level annotation if present.
func getServiceAnnotation(svc *descriptor.Service) *options.Tag {
	if svc.Options == nil {
		return nil
	}
	return getAnnotation[options.Tag](svc.Options, options.E_Openapiv3Tag)
}

// getEnumAnnotation retrieves the OpenAPI v3 enum-level annotation if present.
func getEnumAnnotation(enum *descriptor.Enum) *options.EnumSchema {
	if enum.Options == nil {
		return nil
	}
	return getAnnotation[options.EnumSchema](enum.Options, options.E_Openapiv3Enum)
}
