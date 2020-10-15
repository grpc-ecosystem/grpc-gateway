package gengateway

import (
	"errors"
	"fmt"
	"go/format"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	gen "github.com/grpc-ecosystem/grpc-gateway/v2/internal/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type generator struct {
	reg                *descriptor.Registry
	baseImports        []descriptor.GoPackage
	useRequestContext  bool
	registerFuncSuffix string
	pathType           pathType
	allowPatchFeature  bool
	standalone         bool
}

// New returns a new generator which generates grpc gateway files and the registry to be used for loading files.
func New(useRequestContext bool, registerFuncSuffix, pathTypeString string, allowPatchFeature, standalone bool) (gen.Generator, *descriptor.Registry) {
	reg := descriptor.NewRegistry()

	var imports []descriptor.GoPackage
	for _, pkgpath := range []string{
		"context",
		"io",
		"net/http",
		"github.com/grpc-ecosystem/grpc-gateway/v2/runtime",
		"github.com/grpc-ecosystem/grpc-gateway/v2/utilities",
		"google.golang.org/protobuf/proto",
		"google.golang.org/grpc",
		"google.golang.org/grpc/codes",
		"google.golang.org/grpc/grpclog",
		"google.golang.org/grpc/metadata",
		"google.golang.org/grpc/status",
	} {
		pkg := descriptor.GoPackage{
			Path: pkgpath,
			Name: path.Base(pkgpath),
		}
		if err := reg.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
			for i := 0; ; i++ {
				alias := fmt.Sprintf("%s_%d", pkg.Name, i)
				if err := reg.ReserveGoPackageAlias(alias, pkg.Path); err != nil {
					continue
				}
				pkg.Alias = alias
				break
			}
		}
		imports = append(imports, pkg)
	}

	var pathType pathType
	switch pathTypeString {
	case "", "import":
		// paths=import is default
	case "source_relative":
		pathType = pathTypeSourceRelative
	default:
		glog.Fatalf(`Unknown path type %q: want "import" or "source_relative".`, pathTypeString)
	}

	g := &generator{
		reg:                reg,
		baseImports:        imports,
		useRequestContext:  useRequestContext,
		registerFuncSuffix: registerFuncSuffix,
		pathType:           pathType,
		allowPatchFeature:  allowPatchFeature,
		standalone:         standalone,
	}

	return g, reg
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
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
		formatted, err := format.Source([]byte(code))
		if err != nil {
			glog.Errorf("%v: %s", err, code)
			return nil, err
		}

		name := g.getFilePath(file)
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		filename := fmt.Sprintf("%s.pb.gw.go", base)
		files = append(files, &descriptor.ResponseFile{
			GoPkg: file.GoPkg,
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(filename),
				Content: proto.String(string(formatted)),
			},
		})
	}
	return files, nil
}

func (g *generator) getFilePath(file *descriptor.File) string {
	name := file.GetName()
	if g.pathType == pathTypeImport && file.GoPkg.Path != "" {
		return fmt.Sprintf("%s/%s", file.GoPkg.Path, filepath.Base(name))
	}

	return name
}

func (g *generator) generate(file *descriptor.File) (string, error) {
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.baseImports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}

	if g.standalone {
		imports = append(imports, file.GoPkg)
	}

	for _, svc := range file.Services {
		for _, m := range svc.Methods {
			imports = append(imports, g.addEnumPathParamImports(file, m, pkgSeen)...)
			pkg := m.RequestType.File.GoPkg
			if len(m.Bindings) == 0 ||
				pkg == file.GoPkg || pkgSeen[pkg.Path] {
				continue
			}
			pkgSeen[pkg.Path] = true
			imports = append(imports, pkg)
		}
	}
	params := param{
		File:               file,
		Imports:            imports,
		UseRequestContext:  g.useRequestContext,
		RegisterFuncSuffix: g.registerFuncSuffix,
		AllowPatchFeature:  g.allowPatchFeature,
	}
	if g.reg != nil {
		params.OmitPackageDoc = g.reg.GetOmitPackageDoc()
	}
	return applyTemplate(params, g.reg)
}

// addEnumPathParamImports handles adding import of enum path parameter go packages
func (g *generator) addEnumPathParamImports(file *descriptor.File, m *descriptor.Method, pkgSeen map[string]bool) []descriptor.GoPackage {
	var imports []descriptor.GoPackage
	for _, b := range m.Bindings {
		for _, p := range b.PathParams {
			e, err := g.reg.LookupEnum("", p.Target.GetTypeName())
			if err != nil {
				continue
			}
			pkg := e.File.GoPkg
			if pkg == file.GoPkg || pkgSeen[pkg.Path] {
				continue
			}
			pkgSeen[pkg.Path] = true
			imports = append(imports, pkg)
		}
	}
	return imports
}
