package main

import (
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapi"
)

func main() {
	req, err := readRequest()
	if err != nil {
		emitError(err)
		return
	}

	g := genopenapi.New()
	resp, err := g.Generate(req)
	if err != nil {
		emitError(err)
		return
	}

	emitResponse(resp)
}

func readRequest() (*pluginpb.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(data, req); err != nil {
		return nil, err
	}

	return req, nil
}

func emitResponse(resp *pluginpb.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		panic(err)
	}
}

func emitError(err error) {
	emitResponse(&pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())})
}
