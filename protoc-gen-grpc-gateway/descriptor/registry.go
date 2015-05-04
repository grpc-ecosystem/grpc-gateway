package descriptor

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Registry is a registry of information extracted from plugin.CodeGeneratorRequest.
type Registry struct {
	// msgs is a mapping from fully-qualified message name to descriptor
	msgs map[string]*Message

	// files is a mapping from file path to descriptor
	files map[string]*File

	// prefix is a prefix to be inserted to golang pacakge paths generated from proto package names.
	prefix string

	// pkgMap is a user-specified mapping from file path to proto package.
	pkgMap map[string]string

	// pkgAliases is a mapping from package aliases to package paths in go which are already taken.
	pkgAliases map[string]string
}

// NewRegistry returns a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		msgs:       make(map[string]*Message),
		files:      make(map[string]*File),
		pkgMap:     make(map[string]string),
		pkgAliases: make(map[string]string),
	}
}

// Load loads definitions of services, methods, messages and fields from "req".
func (r *Registry) Load(req *plugin.CodeGeneratorRequest) error {
	for _, file := range req.GetProtoFile() {
		r.loadFile(file)
	}
	for _, target := range req.FileToGenerate {
		if err := r.loadServices(target); err != nil {
			return err
		}
	}
	return nil
}

// loadFile loads messages and fiels from "file".
// It does not loads services and methods in "file".  You need to call
// loadServices after loadFiles is called for all files to load services and methods.
func (r *Registry) loadFile(file *descriptor.FileDescriptorProto) {
	pkg := GoPackage{
		Path: r.goPackagePath(file),
		Name: defaultGoPackageName(file),
	}
	if err := r.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
		for i := 0; ; i++ {
			alias := fmt.Sprintf("%s_%d", pkg.Name, i)
			if err := r.ReserveGoPackageAlias(alias, pkg.Path); err == nil {
				pkg.Alias = alias
				break
			}
		}
	}
	f := &File{
		FileDescriptorProto: file,
		GoPkg:               pkg,
	}

	r.files[file.GetName()] = f
	r.registerMsg(f, nil, file.GetMessageType())
}

func (r *Registry) registerMsg(file *File, outerPath []string, msgs []*descriptor.DescriptorProto) {
	for _, md := range msgs {
		m := &Message{
			File:            file,
			Outers:          outerPath,
			DescriptorProto: md,
		}
		for _, fd := range md.GetField() {
			m.Fields = append(m.Fields, &Field{
				Message:              m,
				FieldDescriptorProto: fd,
			})
		}
		file.Messages = append(file.Messages, m)
		r.msgs[m.FQMN()] = m
		glog.V(1).Infof("register name: %s", m.FQMN())

		var outers []string
		outers = append(outers, outerPath...)
		outers = append(outers, m.GetName())
		r.registerMsg(file, outers, m.GetNestedType())
	}
}

// LookupMsg looks up a message type by "name".
// It tries to resolve "name" from "location" if "name" is a relative message name.
func (r *Registry) LookupMsg(location, name string) (*Message, error) {
	glog.V(1).Infof("lookup %s from %s", name, location)
	if strings.HasPrefix(name, ".") {
		m, ok := r.msgs[name]
		if !ok {
			return nil, fmt.Errorf("no message found: %s", name)
		}
		return m, nil
	}

	if !strings.HasPrefix(location, ".") {
		location = fmt.Sprintf(".%s", location)
	}
	components := strings.Split(location, ".")
	for len(components) > 0 {
		fqmn := strings.Join(append(components, name), ".")
		if m, ok := r.msgs[fqmn]; ok {
			return m, nil
		}
		components = components[:len(components)-1]
	}
	return nil, fmt.Errorf("no message found: %s", name)
}

// LookupFile looks up a file by name.
func (r *Registry) LookupFile(name string) (*File, error) {
	f, ok := r.files[name]
	if !ok {
		return nil, fmt.Errorf("no such file given: %s", name)
	}
	return f, nil
}

// AddPkgMap adds a mapping from a .proto file to proto package name.
func (r *Registry) AddPkgMap(file, protoPkg string) {
	r.pkgMap[file] = protoPkg
}

// SetPrefix registeres the perfix to be added to go package paths generated from proto package names.
func (r *Registry) SetPrefix(prefix string) {
	r.prefix = prefix
}

// ReserveGoPackageAlias reserves the unique alias of go package.
// If succeeded, the alias will be never used for other packages in generated go files.
// If failed, the alias is already taken by another package, so you need to use another
// alias for the package in your go files.
func (r *Registry) ReserveGoPackageAlias(alias, pkgpath string) error {
	if taken, ok := r.pkgAliases[alias]; ok {
		if taken == pkgpath {
			return nil
		}
		return fmt.Errorf("package name %s is already taken. Use another alias", alias)
	}
	r.pkgAliases[alias] = pkgpath
	return nil
}

// goPackagePath returns the go package path which go files generated from "f" should have.
// It respects the mapping registered by AddPkgMap if exists. Or it generates a path from
// the file name of "f" if otherwise.
func (r *Registry) goPackagePath(f *descriptor.FileDescriptorProto) string {
	name := f.GetName()
	if pkg, ok := r.pkgMap[name]; ok {
		return path.Join(r.prefix, pkg)
	}

	ext := filepath.Ext(name)
	if ext == ".protodevel" || ext == ".proto" {
		name = strings.TrimSuffix(name, ext)
	}
	return path.Join(r.prefix, fmt.Sprintf("%s.pb", name))
}

// defaultGoPackageName returns the default go package name to be used for go files generated from "f".
// You might need to use an unique alias for the package when you import it.  Use ReserveGoPackageAlias to get a unique alias.
func defaultGoPackageName(f *descriptor.FileDescriptorProto) string {
	if f.Options != nil && f.Options.GoPackage != nil {
		return f.Options.GetGoPackage()
	}

	if f.Package == nil {
		base := filepath.Base(f.GetName())
		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)
	}
	return strings.Replace(f.GetPackage(), ".", "_", -1)
}
