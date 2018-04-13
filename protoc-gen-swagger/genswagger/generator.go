package genswagger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/generator"
	swagger_options "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

type generator struct {
	reg *descriptor.Registry
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry) gen.Generator {
	return &generator{reg: reg}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
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

	for _, file := range targets {
		glog.V(1).Infof("Processing %s", file.GetName())
		code, err := applyTemplate(param{File: file, reg: g.reg})
		if err == errNoTargetService {
			glog.V(1).Infof("%s: %v", file.GetName(), err)
			continue
		}
		if err != nil {
			return nil, err
		}

		var formatted bytes.Buffer
		json.Indent(&formatted, []byte(code), "", "  ")

		name := file.GetName()
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		output := fmt.Sprintf("%s.swagger.json", base)
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(formatted.String()),
		})
		glog.V(1).Infof("Will emit %s", output)
	}
	return files, nil
}
