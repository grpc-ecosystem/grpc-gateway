package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/codegenerator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	file                           = flag.String("file", "-", "where to load data from")
	proto3OptionalNullable         = flag.Bool("proto3_optional_nullable", false, "whether Proto3 Optional fields should be marked as x-nullable")
	openAPIConfiguration           = flag.String("openapi_configuration", "", "path to file which describes the OpenAPI Configuration in YAML format")
	recursiveDepth                 = flag.Int("recursive-depth", 1000, "maximum recursion count allowed for a field type")
	versionFlag                    = flag.Bool("gizmocore.testversion", false, "print the current version")
	omitEnumDefaultValue           = flag.Bool("omit_enum_default_value", false, "if set, omit default enum value")
	visibilityRestrictionSelectors = utilities.StringArrayFlag(flag.CommandLine, "visibility_restriction_selectors", "list of `google.api.VisibilityRule` visibility labels to include in the generated output when a visibility annotation is defined. Repeat this option to supply multiple values. Elements without visibility annotations are unaffected by this setting.")
	disableServiceTags             = flag.Bool("disable_service_tags", false, "if set, disables generation of service tags. This is useful if you do not want to expose the names of your backend grpc services.")
	preserveRPCOrder               = flag.Bool("preserve_rpc_order", false, "if true, will ensure the order of paths emitted in openapi swagger files mirror the order of RPC methods found in proto files. If false, emitted paths will be ordered alphabetically.")
	enableRpcDeprecation           = flag.Bool("enable_rpc_deprecation", false, "whether to process grpc method's deprecated option.")
	outputFormat                   = flag.String("output_format", string(genopenapiv3.FormatJSON), fmt.Sprintf("output content format. Allowed values are: `%s`, `%s`", genopenapiv3.FormatJSON, genopenapiv3.FormatYAML))
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {

	flag.Parse()

	if *versionFlag {
		if commit == "unknown" {
			buildInfo, ok := debug.ReadBuildInfo()
			if ok {
				version = buildInfo.Main.Version
				for _, setting := range buildInfo.Settings {
					if setting.Key == "vcs.revision" {
						commit = setting.Value
					}
					if setting.Key == "vcs.time" {
						date = setting.Value
					}
				}
			}
		}
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	var openFile *os.File
	var err error

	if  *file == "-" {
		openFile = os.Stdin
	} else {
		openFile, err = os.Open(*file)
		if err != nil {
			grpclog.Fatal("could not open file", err)
			os.Exit(1)
		}

		defer openFile.Close()
	}

	reg := descriptor.NewRegistry()
	if grpclog.V(1) {
		grpclog.Info("Processing code generator request")
	}

	req, err := codegenerator.ParseRequest(openFile)
	if err != nil {
		grpclog.Fatal(err)
	}

	
	if req.Parameter != nil {
		parseReqParam(req.GetParameter(), flag.CommandLine)
	}

	// Set the naming strategy either directly from the flag, or via the value of the legacy fqn_for_openapi_name
	// flag.
	reg.SetProto3OptionalNullable(*proto3OptionalNullable)
	reg.SetRecursiveDepth(*recursiveDepth)
	reg.SetOmitEnumDefaultValue(*omitEnumDefaultValue)
	reg.SetVisibilityRestrictionSelectors(*visibilityRestrictionSelectors)
	reg.SetDisableServiceTags(*disableServiceTags)
	reg.SetPreserveRPCOrder(*preserveRPCOrder)
	reg.SetEnableRpcDeprecation(*enableRpcDeprecation)

	if err := reg.Load(req); err != nil {
		emitError(err)
		return
	}

	generator := genopenapiv3.NewGenerator(reg, genopenapiv3.Format(*outputFormat))


	targets := make([]*descriptor.File, 0, len(req.FileToGenerate))
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			grpclog.Fatal(err)
		}
		targets = append(targets, f)
	}

	out, err := generator.Generate(targets)
	if grpclog.V(1) {
		grpclog.Info("Processed code generator request")
	}
	if err != nil {
		emitError(err)
		return
	}
	emitFiles(out)
}

// parseReqParam parses a CodeGeneratorRequest parameter and adds the
// extracted values to the given FlagSet and pkgMap. Returns a non-nil
// error if setting a flag failed.
func parseReqParam(param string, f *flag.FlagSet) error {
	if param == "" {
		return nil
	}
	for _, p := range strings.Split(param, ",") {
		flagName, value, valueExists := strings.Cut(p, "=")
		if !valueExists {
			switch flagName {
			case "allow_delete_body":
				if err := f.Set(flagName, "true"); err != nil {
					return fmt.Errorf("cannot set flag %s: %w", p, err)
				}
				continue
			case "allow_merge":
				if err := f.Set(flagName, "true"); err != nil {
					return fmt.Errorf("cannot set flag %s: %w", p, err)
				}
				continue
			case "allow_repeated_fields_in_body":
				if err := f.Set(flagName, "true"); err != nil {
					return fmt.Errorf("cannot set flag %s: %w", p, err)
				}
				continue
			case "include_package_in_tags":
				if err := f.Set(flagName, "true"); err != nil {
					return fmt.Errorf("cannot set flag %s: %w", p, err)
				}
				continue
			}
			if err := f.Set(flagName, ""); err != nil {
				return fmt.Errorf("cannot set flag %s: %w", p, err)
			}
			continue
		}

		if strings.HasPrefix(flagName, "M") {
			continue
		}

		if err := f.Set(flagName, value); err != nil {
			return fmt.Errorf("cannot set flag %s: %w", p, err)
		}
	}
	return nil
}

func emitError(err error) {
	emitResp(&pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitFiles(out []*descriptor.ResponseFile) {
	files := make([]*pluginpb.CodeGeneratorResponse_File, len(out))
	for idx, item := range out {
		files[idx] = item.CodeGeneratorResponse_File
	}
	resp := &pluginpb.CodeGeneratorResponse{File: files}
	codegenerator.SetSupportedFeaturesOnCodeGeneratorResponse(resp)
	emitResp(resp)
}

func emitResp(resp *pluginpb.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		grpclog.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		grpclog.Fatal(err)
	}
}
