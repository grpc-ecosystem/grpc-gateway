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

// GeneratorConfig groups all configuration flags for the OpenAPI v3 generator.
// This improves organization, testability, and documentation of available options.
type GeneratorConfig struct {
	// Input/Output
	ImportPrefix string
	File         string
	OutputFormat string

	// HTTP Options
	AllowDeleteBody bool

	// Merge Options
	AllowMerge    bool
	MergeFileName string

	// Naming Options
	UseJSONNamesForFields     bool
	IncludePackageInTags      bool
	OpenAPINamingStrategy     string
	SimpleOperationIDs        bool
	ExpandSlashedPathPatterns bool

	// Comment Options
	IgnoreComments         bool
	RemoveInternalComments bool

	// Schema Options
	EnumsAsInts             bool
	OmitEnumDefaultValue    bool
	Proto3OptionalNullable  bool
	UseProto3FieldSemantics bool
	RecursiveDepth          int

	// Generation Control
	GenerateUnboundMethods  bool
	DisableDefaultErrors    bool
	DisableServiceTags      bool
	DisableDefaultResponses bool
	PreserveRPCOrder        bool
	EnableRPCDeprecation    bool
	EnableFieldDeprecation  bool

	// Visibility
	VisibilityRestrictionSelectors []string

	// Version
	OpenAPIVersion string

	// External Config
	GrpcAPIConfiguration string
}

// flags holds the parsed command line flags
var cfg = &GeneratorConfig{}

// Command line flags
var (
	importPrefix                   = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	file                           = flag.String("file", "-", "where to load data from")
	allowDeleteBody                = flag.Bool("allow_delete_body", false, "unless set, HTTP DELETE methods may not have a body")
	grpcAPIConfiguration           = flag.String("grpc_api_configuration", "", "path to file which describes the gRPC API Configuration in YAML format")
	allowMerge                     = flag.Bool("allow_merge", false, "if set, generation one OpenAPI file out of multiple protos")
	mergeFileName                  = flag.String("merge_file_name", "apidocs", "target OpenAPI file name prefix after merge")
	useJSONNamesForFields          = flag.Bool("json_names_for_fields", true, "if disabled, the original proto name will be used for generating OpenAPI definitions")
	versionFlag                    = flag.Bool("version", false, "print the current version")
	includePackageInTags           = flag.Bool("include_package_in_tags", false, "if unset, the gRPC service name is added to the `Tags` field of each operation")
	openAPINamingStrategy          = flag.String("openapi_naming_strategy", "legacy", "use the given OpenAPI naming strategy. Allowed values are `legacy`, `fqn`, `simple`, `package`")
	ignoreComments                 = flag.Bool("ignore_comments", false, "if set, all protofile comments are excluded from output")
	removeInternalComments         = flag.Bool("remove_internal_comments", true, "if set, removes all substrings in comments that start with `(--` and end with `--)`")
	disableDefaultErrors           = flag.Bool("disable_default_errors", false, "if set, disables generation of default errors")
	enumsAsInts                    = flag.Bool("enums_as_ints", false, "whether to render enum values as integers, as opposed to string values")
	simpleOperationIDs             = flag.Bool("simple_operation_ids", false, "whether to remove the service prefix in the operationID generation")
	generateUnboundMethods         = flag.Bool("generate_unbound_methods", false, "generate OpenAPI metadata even for RPC methods that have no HttpRule annotation")
	recursiveDepth                 = flag.Int("recursive-depth", 1000, "maximum recursion count allowed for a field type")
	omitEnumDefaultValue           = flag.Bool("omit_enum_default_value", false, "if set, omit default enum value")
	outputFormat                   = flag.String("output_format", string(genopenapiv3.FormatJSON), fmt.Sprintf("output content format. Allowed values are: `%s`, `%s`", genopenapiv3.FormatJSON, genopenapiv3.FormatYAML))
	visibilityRestrictionSelectors = utilities.StringArrayFlag(flag.CommandLine, "visibility_restriction_selectors", "list of `google.api.VisibilityRule` visibility labels to include in the generated output when a visibility annotation is defined. Repeat this option to supply multiple values. Elements without visibility annotations are unaffected by this setting.")
	disableServiceTags             = flag.Bool("disable_service_tags", false, "if set, disables generation of service tags")
	disableDefaultResponses        = flag.Bool("disable_default_responses", false, "if set, disables generation of default responses")
	preserveRPCOrder               = flag.Bool("preserve_rpc_order", false, "if true, will ensure the order of paths emitted in openapi files mirror the order of RPC methods found in proto files")
	enableRpcDeprecation           = flag.Bool("enable_rpc_deprecation", false, "whether to process grpc method's deprecated option")
	enableFieldDeprecation         = flag.Bool("enable_field_deprecation", false, "whether to process proto field's deprecated option")
	expandSlashedPathPatterns      = flag.Bool("expand_slashed_path_patterns", false, "if set, expands path parameters with URI sub-paths into the URI")
	proto3OptionalNullable         = flag.Bool("proto3_optional_nullable", false, "whether Proto3 Optional fields should be marked as nullable")
	useProto3FieldSemantics        = flag.Bool("use_proto3_field_semantics", false, "if set, uses proto3 field semantics for the OpenAPI schema. This means that non-optional fields are required by default")
	openAPIVersion                 = flag.String("openapi_version", "3.1.0", "OpenAPI version to use (currently only 3.1.0 is supported)")

	_ = flag.Bool("logtostderr", false, "Legacy glog compatibility. This flag is a no-op, you can safely remove it")
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// populateConfig copies flag values to the config struct.
func populateConfig() {
	cfg.ImportPrefix = *importPrefix
	cfg.File = *file
	cfg.OutputFormat = *outputFormat
	cfg.AllowDeleteBody = *allowDeleteBody
	cfg.AllowMerge = *allowMerge
	cfg.MergeFileName = *mergeFileName
	cfg.UseJSONNamesForFields = *useJSONNamesForFields
	cfg.IncludePackageInTags = *includePackageInTags
	cfg.OpenAPINamingStrategy = *openAPINamingStrategy
	cfg.SimpleOperationIDs = *simpleOperationIDs
	cfg.ExpandSlashedPathPatterns = *expandSlashedPathPatterns
	cfg.IgnoreComments = *ignoreComments
	cfg.RemoveInternalComments = *removeInternalComments
	cfg.EnumsAsInts = *enumsAsInts
	cfg.OmitEnumDefaultValue = *omitEnumDefaultValue
	cfg.Proto3OptionalNullable = *proto3OptionalNullable
	cfg.UseProto3FieldSemantics = *useProto3FieldSemantics
	cfg.RecursiveDepth = *recursiveDepth
	cfg.GenerateUnboundMethods = *generateUnboundMethods
	cfg.DisableDefaultErrors = *disableDefaultErrors
	cfg.DisableServiceTags = *disableServiceTags
	cfg.DisableDefaultResponses = *disableDefaultResponses
	cfg.PreserveRPCOrder = *preserveRPCOrder
	cfg.EnableRPCDeprecation = *enableRpcDeprecation
	cfg.EnableFieldDeprecation = *enableFieldDeprecation
	cfg.VisibilityRestrictionSelectors = *visibilityRestrictionSelectors
	cfg.OpenAPIVersion = *openAPIVersion
	cfg.GrpcAPIConfiguration = *grpcAPIConfiguration
}

// configureRegistry applies the configuration to a registry.
func configureRegistry(reg *descriptor.Registry) {
	reg.SetPrefix(cfg.ImportPrefix)
	reg.SetAllowDeleteBody(cfg.AllowDeleteBody)
	reg.SetAllowMerge(cfg.AllowMerge)
	reg.SetMergeFileName(cfg.MergeFileName)
	reg.SetUseJSONNamesForFields(cfg.UseJSONNamesForFields)
	reg.SetIncludePackageInTags(cfg.IncludePackageInTags)
	reg.SetOpenAPINamingStrategy(cfg.OpenAPINamingStrategy)
	reg.SetIgnoreComments(cfg.IgnoreComments)
	reg.SetRemoveInternalComments(cfg.RemoveInternalComments)
	reg.SetEnumsAsInts(cfg.EnumsAsInts)
	reg.SetDisableDefaultErrors(cfg.DisableDefaultErrors)
	reg.SetSimpleOperationIDs(cfg.SimpleOperationIDs)
	reg.SetGenerateUnboundMethods(cfg.GenerateUnboundMethods)
	reg.SetRecursiveDepth(cfg.RecursiveDepth)
	reg.SetOmitEnumDefaultValue(cfg.OmitEnumDefaultValue)
	reg.SetVisibilityRestrictionSelectors(cfg.VisibilityRestrictionSelectors)
	reg.SetDisableServiceTags(cfg.DisableServiceTags)
	reg.SetDisableDefaultResponses(cfg.DisableDefaultResponses)
	reg.SetPreserveRPCOrder(cfg.PreserveRPCOrder)
	reg.SetEnableRpcDeprecation(cfg.EnableRPCDeprecation)
	reg.SetEnableFieldDeprecation(cfg.EnableFieldDeprecation)
	reg.SetExpandSlashedPathPatterns(cfg.ExpandSlashedPathPatterns)
	reg.SetProto3OptionalNullable(cfg.Proto3OptionalNullable)
	reg.SetUseProto3FieldSemantics(cfg.UseProto3FieldSemantics)
}

func main() {
	flag.Parse()

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}

	populateConfig()

	if err := run(); err != nil {
		emitError(err)
	}
}

// printVersion prints version information from build info or goreleaser.
func printVersion() {
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
}

// run contains the main generation logic.
func run() error {
	reg := descriptor.NewRegistry()

	if grpclog.V(1) {
		grpclog.Info("Processing code generator request")
	}

	// Read input
	f := os.Stdin
	if cfg.File != "-" {
		var err error
		f, err = os.Open(cfg.File)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
	}

	if grpclog.V(1) {
		grpclog.Info("Parsing code generator request")
	}

	req, err := codegenerator.ParseRequest(f)
	if err != nil {
		return err
	}

	if grpclog.V(1) {
		grpclog.Info("Parsed code generator request")
	}

	// Parse request parameters
	pkgMap := make(map[string]string)
	if req.Parameter != nil {
		if err := parseReqParam(req.GetParameter(), flag.CommandLine, pkgMap); err != nil {
			return fmt.Errorf("error parsing flags: %w", err)
		}
	}

	// Re-populate config after potential parameter override
	populateConfig()

	// Configure registry
	configureRegistry(reg)

	for k, v := range pkgMap {
		reg.AddPkgMap(k, v)
	}

	// Load gRPC API configuration if specified
	if cfg.GrpcAPIConfiguration != "" {
		if err := reg.LoadGrpcAPIServiceFromYAML(cfg.GrpcAPIConfiguration); err != nil {
			return err
		}
	}

	// Validate format
	format := genopenapiv3.Format(cfg.OutputFormat)
	if err := format.Validate(); err != nil {
		return err
	}

	// Validate naming strategy
	if cfg.OpenAPINamingStrategy != "" {
		if genopenapiv3.LookupNamingStrategy(cfg.OpenAPINamingStrategy) == nil {
			return fmt.Errorf("invalid openapi_naming_strategy %q: allowed values are fqn, legacy, simple, package", cfg.OpenAPINamingStrategy)
		}
	}

	// Create generator
	g, err := genopenapiv3.New(reg, format, cfg.OpenAPIVersion)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	// Add error definitions
	if err := genopenapiv3.AddErrorDefs(reg); err != nil {
		return err
	}

	// Load proto files
	if err := reg.Load(req); err != nil {
		return err
	}

	// Build target files
	targets := make([]*descriptor.File, 0, len(req.FileToGenerate))
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			grpclog.Fatal(err)
		}
		targets = append(targets, f)
	}

	// Generate output
	out, err := g.Generate(targets)
	if grpclog.V(1) {
		grpclog.Info("Processed code generator request")
	}
	if err != nil {
		return err
	}

	emitFiles(out)
	return nil
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
	for p := range strings.SplitSeq(param, ",") {
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
