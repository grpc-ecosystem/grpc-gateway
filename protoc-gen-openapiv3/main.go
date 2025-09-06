package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/codegenerator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
)

var (
	file                           = flag.String("file", "-", "where to load data from")
	oneOfStrategy                  = flag.String("oneof_strategy", "oneOf", "how to handle oneofs")
	outputFormat                   = flag.String("output_format", string(genopenapiv3.FormatJSON), fmt.Sprintf("output content format. Allowed values are: `%s`, `%s`", genopenapiv3.FormatJSON, genopenapiv3.FormatYAML))
	visibilityRestrictionSelectors = utilities.StringArrayFlag(flag.CommandLine, "visibility_restriction_selectors", "list of `google.api.VisibilityRule` visibility labels to include in the generated output when a visibility annotation is defined. Repeat this option to supply multiple values. Elements without visibility annotations are unaffected by this setting.")
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {

	flag.Parse()

	var openFile *os.File
	var err error

	if *file == "-" {
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
		err := parseReqParam(req.GetParameter(), flag.CommandLine)
		if err != nil {
			emitError(err)
			return
		}
	}

	reg.SetOneOfStrategy(*oneOfStrategy)
	reg.SetVisibilityRestrictionSelectors(*visibilityRestrictionSelectors)

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
			return fmt.Errorf("no value exists for flag: %s", flagName)
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
