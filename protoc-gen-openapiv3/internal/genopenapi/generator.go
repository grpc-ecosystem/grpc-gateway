package genopenapi

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	openapioptions "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/pluginpb"
)

var errNoTargetService = errors.New("no target service defined in the file")

type wrapperv3 struct {
	fileName      string
	openapiv3Spec *OpenAPIV3Document
}

type generator struct {
	reg    *descriptor.Registry
	format Format
}

type GeneratorOptions struct {
	Registry       *descriptor.Registry
	RecursiveDepth int
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry, format Format) gen.Generator {
	return &generator{
		reg:    reg,
		format: format,
	}
}

// Merge a lot of OpenAPI file (wrapper) to single one OpenAPI file.
// The OpenAPI info object corresponds to the one in the first file.
func mergeTargetFile(targets []*wrapperv3, mergeFileName string) *wrapperv3 {
	var mergedTarget *wrapperv3
	for _, f := range targets {
		if mergedTarget == nil {
			mergedTarget = &wrapperv3{
				fileName:      mergeFileName,
				openapiv3Spec: f.openapiv3Spec,
			}
		} else {
			// Merge Components
			if mergedTarget.openapiv3Spec.Components == nil {
				mergedTarget.openapiv3Spec.Components = &OpenAPIV3Components{}
			}
			if f.openapiv3Spec.Components != nil {
				for k, v := range f.openapiv3Spec.Components.Schemas {
					mergedTarget.openapiv3Spec.Components.Schemas[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.Responses {
					mergedTarget.openapiv3Spec.Components.Responses[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.Parameters {
					mergedTarget.openapiv3Spec.Components.Parameters[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.RequestBodies {
					mergedTarget.openapiv3Spec.Components.RequestBodies[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.Headers {
					mergedTarget.openapiv3Spec.Components.Headers[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.SecuritySchemes {
					mergedTarget.openapiv3Spec.Components.SecuritySchemes[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.Links {
					mergedTarget.openapiv3Spec.Components.Links[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.Callbacks {
					mergedTarget.openapiv3Spec.Components.Callbacks[k] = v
				}
				for k, v := range f.openapiv3Spec.Components.Examples {
					mergedTarget.openapiv3Spec.Components.Examples[k] = v
				}
			}
			// Merge Paths
			for k, v := range f.openapiv3Spec.Paths {
				mergedTarget.openapiv3Spec.Paths[k] = v
			}
			// Merge Security
			mergedTarget.openapiv3Spec.Security = append(mergedTarget.openapiv3Spec.Security, f.openapiv3Spec.Security...)
			// Merge Tags
			mergedTarget.openapiv3Spec.Tags = append(mergedTarget.openapiv3Spec.Tags, f.openapiv3Spec.Tags...)
		}
	}
	return mergedTarget
}

// encodeOpenAPI converts OpenAPI file obj to pluginpb.CodeGeneratorResponse_File
func encodeOpenAPI(file *wrapperv3, format Format) (*descriptor.ResponseFile, error) {
	var contentBuf bytes.Buffer
	enc, err := format.NewEncoder(&contentBuf)
	if err != nil {
		return nil, err
	}

	if err := enc.Encode(*file.openapiv3Spec); err != nil {
		return nil, err
	}

	name := file.fileName
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.swagger."+string(format), base)
	return &descriptor.ResponseFile{
		CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(contentBuf.String()),
		},
	}, nil
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	if g.reg.IsAllowMerge() {
		var mergedTarget *descriptor.File
		// try to find proto leader
		for _, f := range targets {
			if proto.HasExtension(f.Options, openapioptions.E_Openapiv3Swagger) {
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

	var openapis []*wrapperv3
	for _, file := range targets {
		if grpclog.V(1) {
			grpclog.Infof("Processing %s", file.GetName())
		}
		swagger, err := applyTemplateV3(param{File: file, reg: g.reg})
		if errors.Is(err, errNoTargetService) {
			if grpclog.V(1) {
				grpclog.Infof("%s: %v", file.GetName(), err)
			}
			continue
		}
		if err != nil {
			return nil, err
		}
		openapis = append(openapis, &wrapperv3{
			fileName:      file.GetName(),
			openapiv3Spec: &swagger,
		})
	}

	if g.reg.IsAllowMerge() {
		targetOpenAPI := mergeTargetFile(openapis, g.reg.GetMergeFileName())
		if !g.reg.IsPreserveRPCOrder() {
			targetOpenAPI.openapiv3Spec.sortPathsAlphabetically()
		}
		f, err := encodeOpenAPI(targetOpenAPI, g.format)
		if err != nil {
			return nil, fmt.Errorf("failed to encode OpenAPI for %s: %w", g.reg.GetMergeFileName(), err)
		}
		files = append(files, f)
		if grpclog.V(1) {
			grpclog.Infof("New OpenAPI file will emit")
		}
	} else {
		for _, file := range openapis {
			if !g.reg.IsPreserveRPCOrder() {
				file.openapiv3Spec.sortPathsAlphabetically()
			}
			f, err := encodeOpenAPI(file, g.format)
			if err != nil {
				return nil, fmt.Errorf("failed to encode OpenAPI for %s: %w", file.fileName, err)
			}
			files = append(files, f)
			if grpclog.V(1) {
				grpclog.Infof("New OpenAPI file will emit")
			}
		}
	}
	return files, nil
}

func (so *OpenAPIV3Document) sortPathsAlphabetically() {
	sorted := make(OpenAPIV3Paths)
	var keys []string
	for k := range so.Paths {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sorted[k] = so.Paths[k]
	}
	so.Paths = sorted
}

// AddErrorDefs Adds google.rpc.Status and google.protobuf.Any
// to registry (used for error-related API responses)
func AddErrorDefs(reg *descriptor.Registry) error {
	// load internal protos
	any := protodesc.ToFileDescriptorProto((&anypb.Any{}).ProtoReflect().Descriptor().ParentFile())
	any.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	status := protodesc.ToFileDescriptorProto((&statuspb.Status{}).ProtoReflect().Descriptor().ParentFile())
	status.SourceCodeInfo = new(descriptorpb.SourceCodeInfo)
	return reg.Load(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			any,
			status,
		},
	})
}
