package genopenapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	anypb "github.com/golang/protobuf/ptypes/any"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	openapi_options "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"

	//nolint:staticcheck // Known issue, will be replaced when possible
	legacydescriptor "github.com/golang/protobuf/descriptor"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

type OutputFormat string

const (
	OutputFormatJSON OutputFormat = "json"
	OutputFormatYAML OutputFormat = "yaml"
)

type Config struct {
	OutputFormat OutputFormat
}

type generator struct {
	cfg Config
	reg *descriptor.Registry
}

type wrapper struct {
	fileName string
	openapi  *Openapi
}

type param struct {
	*descriptor.File
	reg *descriptor.Registry
}

type GeneratorOptions struct {
	Registry       *descriptor.Registry
	RecursiveDepth int
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry, cfg Config) gen.Generator {
	return &generator{
		reg: reg,
		cfg: cfg,
	}
}

// Merge a lot of OpenAPI file (wrapper) to single one OpenAPI file
func mergeTargetFile(targets []*wrapper, mergeFileName string) *wrapper {
	var mergedTarget *wrapper
	for _, f := range targets {
		if mergedTarget == nil {
			mergedTarget = &wrapper{
				fileName: mergeFileName,
				openapi:  f.openapi,
			}
		} else {
			for k, v := range f.openapi.Components.Schemas {
				mergedTarget.openapi.Components.Schemas[k] = v
			}
			for k, v := range f.openapi.Paths {
				mergedTarget.openapi.Paths[k] = v
			}
			for k, v := range f.openapi.Tags {
				mergedTarget.openapi.Tags[k] = v
			}
			mergedTarget.openapi.Security = append(mergedTarget.openapi.Security, f.openapi.Security...)
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
func (so Openapi) MarshalJSON() ([]byte, error) {
	type alias Openapi
	return extensionMarshalJSON(alias(so), so.Extensions)
}

func (so Info) MarshalJSON() ([]byte, error) {
	type alias Info
	return extensionMarshalJSON(alias(so), so.Extensions)
}

func (so SecurityScheme) MarshalJSON() ([]byte, error) {
	type alias SecurityScheme
	return extensionMarshalJSON(alias(so), so.Extensions)
}

func (so Operation) MarshalJSON() ([]byte, error) {
	type alias Operation
	return extensionMarshalJSON(alias(so), so.Extensions)
}

func (so Response) MarshalJSON() ([]byte, error) {
	type alias Response
	return extensionMarshalJSON(alias(so), so.Extensions)
}

func extensionMarshalJSON(so interface{}, extensions map[string]interface{}) ([]byte, error) {
	// TODO(anjmao): May need to sort extensions by key for predictable output.
	// To append arbitrary keys to the struct we'll render into json,
	// we're creating another struct that embeds the original one, and
	// its extra fields:
	//
	// The struct will look like
	// struct {
	//   *openapiCore
	//   XGrpcGatewayFoo json.RawMessage `json:"x-grpc-gateway-foo"`
	//   XGrpcGatewayBar json.RawMessage `json:"x-grpc-gateway-bar"`
	// }
	// and thus render into what we want -- the JSON of openapiCore with the
	// extensions appended.
	fields := []reflect.StructField{
		{ // embedded
			Name:      "Embedded",
			Type:      reflect.TypeOf(so),
			Anonymous: true,
		},
	}
	for key, value := range extensions {
		fields = append(fields, reflect.StructField{
			Name: fieldName(key),
			Type: reflect.TypeOf(value),
			Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", key)),
		})
	}

	t := reflect.StructOf(fields)
	s := reflect.New(t).Elem()
	s.Field(0).Set(reflect.ValueOf(so))
	for key, value := range extensions {
		s.FieldByName(fieldName(key)).Set(reflect.ValueOf(value))
	}
	return json.Marshal(s.Interface())
}

// encodeOpenAPI converts OpenAPI file obj to pluginpb.CodeGeneratorResponse_File in JSON format.
func encodeOpenAPI(file *wrapper, format OutputFormat) (*descriptor.ResponseFile, error) {
	switch format {
	case OutputFormatJSON:
		return encodeOpenAPIJSON(file)
	case OutputFormatYAML:
		return encodeOpenAPIJYAML(file)
	default:
		return nil, fmt.Errorf("unsupported output format %q", format)
	}
}

// encodeOpenAPIJSON converts OpenAPI file obj to pluginpb.CodeGeneratorResponse_File in JSON format.
func encodeOpenAPIJSON(file *wrapper) (*descriptor.ResponseFile, error) {
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")
	if err := enc.Encode(*file.openapi); err != nil {
		return nil, err
	}
	name := file.fileName
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.openapi.json", base)
	return &descriptor.ResponseFile{
		CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(formatted.String()),
		},
	}, nil
}

// encodeOpenAPIYAML converts OpenAPI file obj to pluginpb.CodeGeneratorResponse_File in YAML format.
func encodeOpenAPIJYAML(file *wrapper) (*descriptor.ResponseFile, error) {
	formatted, err := yaml.Marshal(file.openapi)
	if err != nil {
		return nil, err
	}
	name := file.fileName
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.openapi.yaml", base)
	return &descriptor.ResponseFile{
		CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(string(formatted)),
		},
	}, nil
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	if g.reg.IsAllowMerge() {
		var mergedTarget *descriptor.File
		// try to find proto leader
		for _, f := range targets {
			if proto.HasExtension(f.Options, openapi_options.E_Openapiv3Document) {
				mergedTarget = f
				break
			}
		}
		// merge protos to leader
		for _, f := range targets {
			if mergedTarget == nil {
				mergedTarget = f
			} else if mergedTarget != f {
				mergedTarget.Enums = append(mergedTarget.Enums, f.Enums...)
				mergedTarget.Messages = append(mergedTarget.Messages, f.Messages...)
				mergedTarget.Services = append(mergedTarget.Services, f.Services...)
			}
		}

		targets = nil
		targets = append(targets, mergedTarget)
	}

	var openapis []*wrapper
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
		openapis = append(openapis, &wrapper{
			fileName: file.GetName(),
			openapi:  swagger,
		})
	}

	if g.reg.IsAllowMerge() {
		targetOpenAPI := mergeTargetFile(openapis, g.reg.GetMergeFileName())
		f, err := encodeOpenAPI(targetOpenAPI, g.cfg.OutputFormat)
		if err != nil {
			return nil, fmt.Errorf("failed to encode OpenAPI for %s: %s", g.reg.GetMergeFileName(), err)
		}
		files = append(files, f)
		glog.V(1).Infof("New OpenAPI file will emit")
	} else {
		for _, file := range openapis {
			f, err := encodeOpenAPI(file, g.cfg.OutputFormat)
			if err != nil {
				return nil, fmt.Errorf("failed to encode OpenAPI for %s: %s", file.fileName, err)
			}
			files = append(files, f)
			glog.V(1).Infof("New OpenAPI file will emit")
		}
	}
	return files, nil
}

// AddErrorDefs Adds google.rpc.Status and google.protobuf.Any
// to registry (used for error-related API responses)
func AddErrorDefs(reg *descriptor.Registry) error {
	// load internal protos
	any, _ := legacydescriptor.MessageDescriptorProto(&anypb.Any{})
	any.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	status, _ := legacydescriptor.MessageDescriptorProto(&statuspb.Status{})
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
