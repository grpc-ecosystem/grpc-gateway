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
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	importPrefix               = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	file                       = flag.String("file", "-", "where to load data from")
	allowDeleteBody            = flag.Bool("allow_delete_body", false, "unless set, HTTP DELETE methods may not have a body")
	grpcAPIConfiguration       = flag.String("grpc_api_configuration", "", "path to file which describes the gRPC API Configuration in YAML format")
	allowMerge                 = flag.Bool("allow_merge", false, "if set, generation one OpenAPI file out of multiple protos")
	mergeFileName              = flag.String("merge_file_name", "apidocs", "target OpenAPI file name prefix after merge")
	useJSONNamesForFields      = flag.Bool("json_names_for_fields", true, "if disabled, the original proto name will be used for generating OpenAPI definitions")
	versionFlag                = flag.Bool("version", false, "print the current version")
	includePackageInTags       = flag.Bool("include_package_in_tags", false, "if unset, the gRPC service name is added to the `Tags` field of each operation")
	openAPINamingStrategy      = flag.String("openapi_naming_strategy", "legacy", "use the given OpenAPI naming strategy. Allowed values are `legacy`, `fqn`, `simple`, `package`")
	ignoreComments             = flag.Bool("ignore_comments", false, "if set, all protofile comments are excluded from output")
	removeInternalComments     = flag.Bool("remove_internal_comments", true, "if set, removes all substrings in comments that start with `(--` and end with `--)`")
	disableDefaultErrors       = flag.Bool("disable_default_errors", false, "if set, disables generation of default errors")
	enumsAsInts                = flag.Bool("enums_as_ints", false, "whether to render enum values as integers, as opposed to string values")
	simpleOperationIDs         = flag.Bool("simple_operation_ids", false, "whether to remove the service prefix in the operationID generation")
	generateUnboundMethods     = flag.Bool("generate_unbound_methods", false, "generate OpenAPI metadata even for RPC methods that have no HttpRule annotation")
	recursiveDepth             = flag.Int("recursive-depth", 1000, "maximum recursion count allowed for a field type")
	omitEnumDefaultValue       = flag.Bool("omit_enum_default_value", false, "if set, omit default enum value")
	outputFormat               = flag.String("output_format", string(genopenapiv3.FormatJSON), fmt.Sprintf("output content format. Allowed values are: `%s`, `%s`", genopenapiv3.FormatJSON, genopenapiv3.FormatYAML))
	disableServiceTags         = flag.Bool("disable_service_tags", false, "if set, disables generation of service tags")
	disableDefaultResponses    = flag.Bool("disable_default_responses", false, "if set, disables generation of default responses")
	preserveRPCOrder           = flag.Bool("preserve_rpc_order", false, "if true, will ensure the order of paths emitted in openapi files mirror the order of RPC methods found in proto files")
	enableRpcDeprecation       = flag.Bool("enable_rpc_deprecation", false, "whether to process grpc method's deprecated option")
	enableFieldDeprecation     = flag.Bool("enable_field_deprecation", false, "whether to process proto field's deprecated option")
	expandSlashedPathPatterns  = flag.Bool("expand_slashed_path_patterns", false, "if set, expands path parameters with URI sub-paths into the URI")
	proto3OptionalNullable     = flag.Bool("proto3_optional_nullable", false, "whether Proto3 Optional fields should be marked as nullable")
	useProto3FieldSemantics    = flag.Bool("use_proto3_field_semantics", true, "if set, uses proto3 field semantics for the OpenAPI schema. This means that non-optional fields are required by default")
	openAPIVersion             = flag.String("openapi_version", "3.0.3", "OpenAPI version to use (3.0.3, 3.1.0)")

	_ = flag.Bool("logtostderr", false, "Legacy glog compatibility. This flag is a no-op, you can safely remove it")
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

	reg := descriptor.NewRegistry()
	if grpclog.V(1) {
		grpclog.Info("Processing code generator request")
	}
	f := os.Stdin
	if *file != "-" {
		var err error
		f, err = os.Open(*file)
		if err != nil {
			grpclog.Fatal(err)
		}
	}
	if grpclog.V(1) {
		grpclog.Info("Parsing code generator request")
	}
	req, err := codegenerator.ParseRequest(f)
	if err != nil {
		grpclog.Fatal(err)
	}
	if grpclog.V(1) {
		grpclog.Info("Parsed code generator request")
	}
	pkgMap := make(map[string]string)
	if req.Parameter != nil {
		if err := parseReqParam(req.GetParameter(), flag.CommandLine, pkgMap); err != nil {
			grpclog.Fatalf("Error parsing flags: %v", err)
		}
	}

	reg.SetPrefix(*importPrefix)
	reg.SetAllowDeleteBody(*allowDeleteBody)
	reg.SetAllowMerge(*allowMerge)
	reg.SetMergeFileName(*mergeFileName)
	reg.SetUseJSONNamesForFields(*useJSONNamesForFields)
	reg.SetIncludePackageInTags(*includePackageInTags)
	reg.SetOpenAPINamingStrategy(*openAPINamingStrategy)
	reg.SetIgnoreComments(*ignoreComments)
	reg.SetRemoveInternalComments(*removeInternalComments)
	reg.SetEnumsAsInts(*enumsAsInts)
	reg.SetDisableDefaultErrors(*disableDefaultErrors)
	reg.SetSimpleOperationIDs(*simpleOperationIDs)
	reg.SetGenerateUnboundMethods(*generateUnboundMethods)
	reg.SetRecursiveDepth(*recursiveDepth)
	reg.SetOmitEnumDefaultValue(*omitEnumDefaultValue)
	reg.SetDisableServiceTags(*disableServiceTags)
	reg.SetDisableDefaultResponses(*disableDefaultResponses)
	reg.SetPreserveRPCOrder(*preserveRPCOrder)
	reg.SetEnableRpcDeprecation(*enableRpcDeprecation)
	reg.SetEnableFieldDeprecation(*enableFieldDeprecation)
	reg.SetExpandSlashedPathPatterns(*expandSlashedPathPatterns)
	reg.SetProto3OptionalNullable(*proto3OptionalNullable)
	reg.SetUseProto3FieldSemantics(*useProto3FieldSemantics)

	for k, v := range pkgMap {
		reg.AddPkgMap(k, v)
	}

	if *grpcAPIConfiguration != "" {
		if err := reg.LoadGrpcAPIServiceFromYAML(*grpcAPIConfiguration); err != nil {
			emitError(err)
			return
		}
	}

	format := genopenapiv3.Format(*outputFormat)
	if err := format.Validate(); err != nil {
		emitError(err)
		return
	}

	// Validate naming strategy
	if *openAPINamingStrategy != "" {
		if genopenapiv3.LookupNamingStrategy(*openAPINamingStrategy) == nil {
			emitError(fmt.Errorf("invalid openapi_naming_strategy %q: allowed values are fqn, legacy, simple, package", *openAPINamingStrategy))
			return
		}
	}

	g := genopenapiv3.New(reg, format, *openAPIVersion)

	if err := genopenapiv3.AddErrorDefs(reg); err != nil {
		emitError(err)
		return
	}

	if err := reg.Load(req); err != nil {
		emitError(err)
		return
	}

	targets := make([]*descriptor.File, 0, len(req.FileToGenerate))
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			grpclog.Fatal(err)
		}
		targets = append(targets, f)
	}

	out, err := g.Generate(targets)
	if grpclog.V(1) {
		grpclog.Info("Processed code generator request")
	}
	if err != nil {
		emitError(err)
		return
	}
	emitFiles(out)
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

func emitError(err error) {
	emitResp(&pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())})
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

// parseReqParam parses a CodeGeneratorRequest parameter and adds the
// extracted values to the given FlagSet and pkgMap.
func parseReqParam(param string, f *flag.FlagSet, pkgMap map[string]string) error {
	if param == "" {
		return nil
	}
	for _, p := range strings.Split(param, ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 1 {
			switch spec[0] {
			case "allow_delete_body", "allow_merge", "include_package_in_tags":
				if err := f.Set(spec[0], "true"); err != nil {
					return fmt.Errorf("cannot set flag %s: %w", p, err)
				}
				continue
			}
			if err := f.Set(spec[0], ""); err != nil {
				return fmt.Errorf("cannot set flag %s: %w", p, err)
			}
			continue
		}
		name, value := spec[0], spec[1]
		if strings.HasPrefix(name, "M") {
			pkgMap[name[1:]] = value
			continue
		}
		if err := f.Set(name, value); err != nil {
			return fmt.Errorf("cannot set flag %s: %w", p, err)
		}
	}
	return nil
}
