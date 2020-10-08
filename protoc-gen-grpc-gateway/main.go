// Command protoc-gen-grpc-gateway is a plugin for Google protocol buffer
// compiler to generate a reverse-proxy, which converts incoming RESTful
// HTTP/1 requests gRPC invocation.
// You rarely need to run this program directly. Instead, put this program
// into your $PATH with a name "protoc-gen-grpc-gateway" and run
//   protoc --grpc-gateway_out=output_directory path/to/input.proto
//
// See README.md for more details.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway/internal/gengateway"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	importPrefix               = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	importPath                 = flag.String("import_path", "", "used as the package if no input files declare go_package. If it contains slashes, everything up to the rightmost slash is ignored.")
	registerFuncSuffix         = flag.String("register_func_suffix", "Handler", "used to construct names of generated Register*<Suffix> methods.")
	useRequestContext          = flag.Bool("request_context", true, "determine whether to use http.Request's context or not")
	allowDeleteBody            = flag.Bool("allow_delete_body", false, "unless set, HTTP DELETE methods may not have a body")
	grpcAPIConfiguration       = flag.String("grpc_api_configuration", "", "path to gRPC API Configuration in YAML format")
	pathType                   = flag.String("paths", "", "specifies how the paths of generated files are structured")
	modulePath                 = flag.String("module", "", "specifies a module prefix that will be stripped from the go package to determine the output directory")
	allowRepeatedFieldsInBody  = flag.Bool("allow_repeated_fields_in_body", false, "allows to use repeated field in `body` and `response_body` field of `google.api.http` annotation option")
	repeatedPathParamSeparator = flag.String("repeated_path_param_separator", "csv", "configures how repeated fields should be split. Allowed values are `csv`, `pipes`, `ssv` and `tsv`.")
	allowPatchFeature          = flag.Bool("allow_patch_feature", true, "determines whether to use PATCH feature involving update masks (using google.protobuf.FieldMask).")
	omitPackageDoc             = flag.Bool("omit_package_doc", false, "if true, no package comment will be included in the generated code")
	standalone                 = flag.Bool("standalone", false, "generates a standalone gateway package, which imports the target service package")
	versionFlag                = flag.Bool("version", false, "print the current version")
	warnOnUnboundMethods       = flag.Bool("warn_on_unbound_methods", false, "emit a warning message if an RPC method has no HttpRule annotation")
	generateUnboundMethods     = flag.Bool("generate_unbound_methods", false, "generate proxy methods even for RPC methods that have no HttpRule annotation")
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if *versionFlag {
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	reg := descriptor.NewRegistry()

	protogen.Options{
		// FIXME: ParamFunc is not enough at this point because it does not receive all params.
		//        Some are swallowed by protogen like "paths".
		//        This problem will go away when the code generation is completely rewritten
		//        to support protogen.Plugin.
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		// FIXME: still needed to parse request parameter and apply flags manually, see the comment above.
		parseFlags(reg, plugin.Request.GetParameter())
		if err := applyFlags(reg); err != nil {
			return err
		}

		glog.V(1).Infof("Parsing code generator request")

		if err := reg.Load(plugin.Request); err != nil {
			return err
		}

		unboundHTTPRules := reg.UnboundExternalHTTPRules()
		if len(unboundHTTPRules) != 0 {
			return fmt.Errorf("HTTP rules without a matching selector: %s", strings.Join(unboundHTTPRules, ", "))
		}

		var targets []*descriptor.File
		for _, target := range plugin.Request.FileToGenerate {
			f, err := reg.LookupFile(target)
			if err != nil {
				return err
			}
			targets = append(targets, f)
		}

		g := gengateway.New(reg, *useRequestContext, *registerFuncSuffix, *pathType, *modulePath, *allowPatchFeature, *standalone)
		files, err := g.Generate(targets)
		for _, f := range files {
			glog.V(1).Infof("NewGeneratedFile %q in %s", f.GetName(), f.GoPkg)
			genFile := plugin.NewGeneratedFile(f.GetName(), protogen.GoImportPath(f.GoPkg.Path))
			if _, err := genFile.Write([]byte(f.GetContent())); err != nil {
				return err
			}
		}

		glog.V(1).Info("Processed code generator request")

		return err
	})
}

func parseFlags(reg *descriptor.Registry, parameter string) {
	for _, p := range strings.Split(parameter, ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 1 {
			if err := flag.CommandLine.Set(spec[0], ""); err != nil {
				glog.Fatalf("Cannot set flag %s", p)
			}
			continue
		}

		name, value := spec[0], spec[1]

		if strings.HasPrefix(name, "M") {
			reg.AddPkgMap(name[1:], value)
			continue
		}
		if err := flag.CommandLine.Set(name, value); err != nil {
			glog.Fatalf("Cannot set flag %s", p)
		}
	}
}

func applyFlags(reg *descriptor.Registry) error {
	if *grpcAPIConfiguration != "" {
		if err := reg.LoadGrpcAPIServiceFromYAML(*grpcAPIConfiguration); err != nil {
			return err
		}
	}
	if *warnOnUnboundMethods && *generateUnboundMethods {
		glog.Warningf("Option warn_on_unbound_methods has no effect when generate_unbound_methods is used.")
	}
	reg.SetStandalone(*standalone)
	reg.SetPrefix(*importPrefix)
	reg.SetImportPath(*importPath)
	reg.SetAllowDeleteBody(*allowDeleteBody)
	reg.SetAllowRepeatedFieldsInBody(*allowRepeatedFieldsInBody)
	reg.SetOmitPackageDoc(*omitPackageDoc)
	reg.SetWarnOnUnboundMethods(*warnOnUnboundMethods)
	reg.SetGenerateUnboundMethods(*generateUnboundMethods)
	return reg.SetRepeatedPathParamSeparator(*repeatedPathParamSeparator)
}
