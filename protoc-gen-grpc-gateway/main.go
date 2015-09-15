package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gengo/grpc-gateway/log"
	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway/gengateway"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var (
	importPrefix = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
)

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	log.Infoln("Parsing code generator request")
	input, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("Failed to read code generator request: %v", err)
		return nil, err
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(input, req); err != nil {
		log.Errorf("Failed to unmarshal code generator request: %v", err)
		return nil, err
	}
	log.Infoln("Parsed code generator request")
	return req, nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	reg := descriptor.NewRegistry()

	log.Infoln("Processing code generator request")
	req, err := parseReq(os.Stdin)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	if req.Parameter != nil {
		for _, p := range strings.Split(req.GetParameter(), ",") {
			spec := strings.SplitN(p, "=", 2)
			if len(spec) == 1 {
				if err := flag.CommandLine.Set(spec[0], ""); err != nil {
					log.Errorf("Cannot set flag %s", p)
					os.Exit(1)
				}
				continue
			}
			name, value := spec[0], spec[1]
			if strings.HasPrefix(name, "M") {
				reg.AddPkgMap(name[1:], value)
				continue
			}
			if err := flag.CommandLine.Set(name, value); err != nil {
				log.Errorf("Cannot set flag %s", p)
				os.Exit(1)
			}
		}
	}

	g := gengateway.New(reg)

	reg.SetPrefix(*importPrefix)
	if err := reg.Load(req); err != nil {
		emitError(err)
		return
	}

	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		targets = append(targets, f)
	}

	out, err := g.Generate(targets)
	log.Infoln("Processed code generator request")
	if err != nil {
		emitError(err)
		return
	}
	emitFiles(out)
}

func emitFiles(out []*plugin.CodeGeneratorResponse_File) {
	emitResp(&plugin.CodeGeneratorResponse{File: out})
}

func emitError(err error) {
	emitResp(&plugin.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitResp(resp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
