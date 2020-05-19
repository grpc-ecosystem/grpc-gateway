package genswagger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/golang/glog"
	legacydescriptor "github.com/golang/protobuf/descriptor"
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	pluginpb "github.com/golang/protobuf/protoc-gen-go/plugin"
	anypb "github.com/golang/protobuf/ptypes/any"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	swagger_options "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-swagger/options"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

type generator struct {
	reg *descriptor.Registry
}

type wrapper struct {
	fileName string
	swagger  *swaggerObject
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry) gen.Generator {
	return &generator{reg: reg}
}

// Merge a lot of swagger file (wrapper) to single one swagger file
func mergeTargetFile(targets []*wrapper, mergeFileName string) *wrapper {
	var mergedTarget *wrapper
	for _, f := range targets {
		if mergedTarget == nil {
			mergedTarget = &wrapper{
				fileName: mergeFileName,
				swagger:  f.swagger,
			}
		} else {
			for k, v := range f.swagger.Definitions {
				mergedTarget.swagger.Definitions[k] = v
			}
			for k, v := range f.swagger.Paths {
				mergedTarget.swagger.Paths[k] = v
			}
			for k, v := range f.swagger.SecurityDefinitions {
				mergedTarget.swagger.SecurityDefinitions[k] = v
			}
			mergedTarget.swagger.Security = append(mergedTarget.swagger.Security, f.swagger.Security...)
		}
	}
	return mergedTarget
}

// Q: What's up with the alias types here?
// A: We don't want to completely override how these structs are marshaled into
//    JSON, we only want to add fields (see below, extensionMarshalJSON).
//    An infinite recursion would happen if we'd call json.Marshal on the struct
//    that has swaggerObject as an embedded field. To avoid that, we'll create
//    type aliases, and those don't have the custom MarshalJSON methods defined
//    on them. See http://choly.ca/post/go-json-marshalling/ (or, if it ever
//    goes away, use
//    https://web.archive.org/web/20190806073003/http://choly.ca/post/go-json-marshalling/.
func (so swaggerObject) MarshalJSON() ([]byte, error) {
	type alias swaggerObject
	return extensionMarshalJSON(alias(so), so.extensions)
}

func (so swaggerInfoObject) MarshalJSON() ([]byte, error) {
	type alias swaggerInfoObject
	return extensionMarshalJSON(alias(so), so.extensions)
}

func (so swaggerSecuritySchemeObject) MarshalJSON() ([]byte, error) {
	type alias swaggerSecuritySchemeObject
	return extensionMarshalJSON(alias(so), so.extensions)
}

func (so swaggerOperationObject) MarshalJSON() ([]byte, error) {
	type alias swaggerOperationObject
	return extensionMarshalJSON(alias(so), so.extensions)
}

func (so swaggerResponseObject) MarshalJSON() ([]byte, error) {
	type alias swaggerResponseObject
	return extensionMarshalJSON(alias(so), so.extensions)
}

func extensionMarshalJSON(so interface{}, extensions []extension) ([]byte, error) {
	// To append arbitrary keys to the struct we'll render into json,
	// we're creating another struct that embeds the original one, and
	// its extra fields:
	//
	// The struct will look like
	// struct {
	//   *swaggerCore
	//   XGrpcGatewayFoo json.RawMessage `json:"x-grpc-gateway-foo"`
	//   XGrpcGatewayBar json.RawMessage `json:"x-grpc-gateway-bar"`
	// }
	// and thus render into what we want -- the JSON of swaggerCore with the
	// extensions appended.
	fields := []reflect.StructField{
		{ // embedded
			Name:      "Embedded",
			Type:      reflect.TypeOf(so),
			Anonymous: true,
		},
	}
	for _, ext := range extensions {
		fields = append(fields, reflect.StructField{
			Name: fieldName(ext.key),
			Type: reflect.TypeOf(ext.value),
			Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", ext.key)),
		})
	}

	t := reflect.StructOf(fields)
	s := reflect.New(t).Elem()
	s.Field(0).Set(reflect.ValueOf(so))
	for _, ext := range extensions {
		s.FieldByName(fieldName(ext.key)).Set(reflect.ValueOf(ext.value))
	}
	return json.Marshal(s.Interface())
}

// encodeSwagger converts swagger file obj to pluginpb.CodeGeneratorResponse_File
func encodeSwagger(file *wrapper) (*pluginpb.CodeGeneratorResponse_File, error) {
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")
	if err := enc.Encode(*file.swagger); err != nil {
		return nil, err
	}
	name := file.fileName
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.swagger.json", base)
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(output),
		Content: proto.String(formatted.String()),
	}, nil
}

func (g *generator) Generate(targets []*descriptor.File) ([]*pluginpb.CodeGeneratorResponse_File, error) {
	var files []*pluginpb.CodeGeneratorResponse_File
	if g.reg.IsAllowMerge() {
		var mergedTarget *descriptor.File
		// try to find proto leader
		for _, f := range targets {
			if proto.HasExtension(f.Options, swagger_options.E_Openapiv2Swagger) {
				mergedTarget = f
				break
			}
		}
		// merge protos to leader
		for _, f := range targets {
			if mergedTarget == nil {
				mergedTarget = f
			} else {
				mergedTarget.Enums = append(mergedTarget.Enums, f.Enums...)
				mergedTarget.Messages = append(mergedTarget.Messages, f.Messages...)
				mergedTarget.Services = append(mergedTarget.Services, f.Services...)
			}
		}

		targets = nil
		targets = append(targets, mergedTarget)
	}

	var swaggers []*wrapper
	for _, file := range targets {
		glog.V(1).Infof("Processing %s", file.GetName())
		swagger, err := applyTemplate(param{File: file, reg: g.reg})
		if err == errNoTargetService {
			glog.V(1).Infof("%s: %v", file.GetName(), err)
			continue
		}
		if err != nil {
			return nil, err
		}
		swaggers = append(swaggers, &wrapper{
			fileName: file.GetName(),
			swagger:  swagger,
		})
	}

	if g.reg.IsAllowMerge() {
		targetSwagger := mergeTargetFile(swaggers, g.reg.GetMergeFileName())
		f, err := encodeSwagger(targetSwagger)
		if err != nil {
			return nil, fmt.Errorf("failed to encode swagger for %s: %s", g.reg.GetMergeFileName(), err)
		}
		files = append(files, f)
		glog.V(1).Infof("New swagger file will emit")
	} else {
		for _, file := range swaggers {
			f, err := encodeSwagger(file)
			if err != nil {
				return nil, fmt.Errorf("failed to encode swagger for %s: %s", file.fileName, err)
			}
			files = append(files, f)
			glog.V(1).Infof("New swagger file will emit")
		}
	}
	return files, nil
}

// AddErrorDefs Adds google.rpc.Status and google.protobuf.Any
// to registry (used for error-related API responses)
func AddErrorDefs(reg *descriptor.Registry) error {
	// load internal protos
	any, _ := legacydescriptor.ForMessage(&anypb.Any{})
	any.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	status, _ := legacydescriptor.ForMessage(&statuspb.Status{})
	status.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	// TODO(johanbrandhorst): Use new conversion later when possible
	// any := protodesc.ToFileDescriptorProto((&anypb.Any{}).ProtoReflect().Descriptor().ParentFile())
	// status := protodesc.ToFileDescriptorProto((&statuspb.Status{}).ProtoReflect().Descriptor().ParentFile())
	return reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			any,
			status,
		},
	})
}
