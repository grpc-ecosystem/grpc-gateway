package genswagger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

type generator struct {
	reg *descriptor.Registry
}

// New returns a new generator which generates grpc gateway files.
func New(reg *descriptor.Registry) *generator {
	return &generator{reg: reg}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		glog.V(1).Infof("Processing %s", file.GetName())
		code, err := g.generate(file)
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
			Content: proto.String(string(formatted.Bytes())),
		})
		glog.V(1).Infof("Will emit %s", output)
	}
	return files, nil
}

func (g *generator) generate(file *descriptor.File) (string, error) {
	pkgSeen := make(map[string]bool)
	for _, svc := range file.Services {
		for _, m := range svc.Methods {
			pkg := m.RequestType.File.GoPkg
			if pkg == file.GoPkg {
				continue
			}
			if pkgSeen[pkg.Path] {
				continue
			}
			pkgSeen[pkg.Path] = true
		}
	}
	return applyTemplate(param{File: file, reg: g.reg})
}
