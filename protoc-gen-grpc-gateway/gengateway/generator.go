package gengateway

import (
	"errors"
	"fmt"
	"go/format"
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

func New(reg *descriptor.Registry) *generator {
	if err := reg.ReserveGoPackageAlias("json", "encoding/json"); err != nil {
		glog.Fatal("package name json collides to another")
	}
	return &generator{reg: reg}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		glog.V(1).Infof("Processing %s", file.GetName())
		code, err := applyTemplate(file)
		if err == errNoTargetService {
			glog.V(1).Infof("%s: %v", file.GetName(), err)
			continue
		}
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(code))
		if err != nil {
			glog.Errorf("%v: %s", err, code)
			return nil, err
		}
		name := file.GetName()
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		output := fmt.Sprintf("%s.pb.gw.go", base)
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(string(formatted)),
		})
		glog.Infof("Will emit %s", output)
	}
	return files, nil
}
