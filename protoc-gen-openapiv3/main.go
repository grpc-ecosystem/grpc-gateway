// Command protoc-gen-openapiv3 implements an OpenAPI 3.1.0 generator for
// proto files annotated with google.api.http rules.
//
// Status: alpha. The emitted JSON shape is not yet stable — encodings for
// oneofs, wrapper types, enums, and path-template expansion may change
// between minor releases while the mapping rules settle in response to
// real-world feedback. For a production-stable OpenAPI pipeline today,
// use protoc-gen-openapiv2.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/codegenerator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	visibilityRestrictionSelectors = utilities.StringArrayFlag(flag.CommandLine, "visibility_restriction_selectors", "list of `google.api.VisibilityRule` visibility labels to include in the generated output when a visibility annotation is defined. Repeat this option to supply multiple values. Elements without visibility annotations are unaffected by this setting.")
	disableDefaultErrors           = flag.Bool("disable_default_errors", false, "if set, disables generation of default errors. This is useful if you have defined custom error handling")
	allowMerge                     = flag.Bool("allow_merge", false, "if set, generates a single OpenAPI file out of multiple protos")
	mergeFileName                  = flag.String("merge_file_name", "apidocs", "target OpenAPI file name prefix after merge")
)

func main() {
	if err := run(); err != nil {
		emitError(err)
		os.Exit(1)
	}
}

func run() error {
	req, err := codegenerator.ParseRequest(os.Stdin)
	if err != nil {
		return err
	}

	if param := req.GetParameter(); param != "" {
		if err := parseReqParam(param, flag.CommandLine); err != nil {
			return err
		}
	}

	reg := descriptor.NewRegistry()
	reg.SetVisibilityRestrictionSelectors(*visibilityRestrictionSelectors)
	reg.SetDisableDefaultErrors(*disableDefaultErrors)
	reg.SetAllowMerge(*allowMerge)
	reg.SetMergeFileName(*mergeFileName)
	if err := reg.Load(req); err != nil {
		return err
	}

	var targets []*descriptor.File
	for _, name := range req.FileToGenerate {
		f, err := reg.LookupFile(name)
		if err != nil {
			return err
		}
		targets = append(targets, f)
	}

	out, err := genopenapi.Generate(reg, targets)
	if err != nil {
		return err
	}

	emitFiles(out)
	return nil
}

func emitFiles(files []*pluginpb.CodeGeneratorResponse_File) {
	resp := &pluginpb.CodeGeneratorResponse{File: files}
	codegenerator.SetSupportedFeaturesOnCodeGeneratorResponse(resp)
	emitResp(resp)
}

func emitError(err error) {
	// Echo to stderr in addition to the proto response
	grpclog.Infoln(err)
	emitResp(&pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitResp(resp *pluginpb.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// parseReqParam parses the protoc plugin parameter string (comma-separated
// key=value pairs) and sets the corresponding flags. When a flag name is
// present without a value (e.g. "disable_default_errors"), it is treated as
// "true", so that boolean flags can be enabled by name alone.
func parseReqParam(param string, f *flag.FlagSet) error {
	for p := range strings.SplitSeq(param, ",") {
		if p == "" {
			continue
		}
		flagName, val, ok := strings.Cut(p, "=")
		if !ok {
			if err := f.Set(flagName, "true"); err != nil {
				return fmt.Errorf("cannot set flag %s: %w", flagName, err)
			}
			continue
		}
		if err := f.Set(flagName, val); err != nil {
			return fmt.Errorf("cannot set flag %s: %w", flagName, err)
		}
	}
	return nil
}
