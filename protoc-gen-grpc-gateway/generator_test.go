package main

import (
	"testing"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func runGenerateRaw(t *testing.T, input string, expected *plugin.CodeGeneratorResponse) {
	msgTbl = make(map[string]*descriptor.DescriptorProto)
	var req plugin.CodeGeneratorRequest
	if err := proto.UnmarshalText(input, &req); err != nil {
		t.Fatalf("proto.Unmarshal(%q, &req) failed with %v; want success", input, err)
	}

	resp := generate(&req)
	if !proto.Equal(resp, expected) {
		t.Errorf("generate(%s) = %s; want %s", input, proto.MarshalTextString(resp), proto.MarshalTextString(expected))
	}
}

func runGenerate(t *testing.T, input string, outputs map[string]string) {
	var expected plugin.CodeGeneratorResponse
	for fname, content := range outputs {
		expected.File = append(expected.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fname),
			Content: proto.String(content),
		})
	}
	runGenerateRaw(t, input, &expected)
}

func TestGenerateEmtpy(t *testing.T) {
	var expected plugin.CodeGeneratorResponse
	runGenerateRaw(t, "", &expected)
}

func TestGenerator(t *testing.T) {
	runGenerate(t, `
file_to_generate: "example.proto"
proto_file <
  name: "example.proto"
  package: "example"
  message_type <
    name: "SimpleMessage"
    field <
      name: "id"
      number: 1
      label: LABEL_REQUIRED,
      type: TYPE_STRING
    >
  >
  service <
    name: "EchoService"
    method <
      name: "Echo"
      input_type: ".example.SimpleMessage"
      output_type: ".example.SimpleMessage"
      options <
        [gengo.grpc.gateway.ApiMethodOptions.api_options] <
          path: "/v1/example/echo"
          method: "POST"
        >
      >
    >
  >
>`, map[string]string{
		"example.pb.gw.go": `package example

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func handle_EchoService_Echo(ctx context.Context, c EchoServiceClient, req *http.Request) (proto.Message, error) {
	protoReq := new(SimpleMessage)

	if err := json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}

	return c.Echo(ctx, protoReq)
}

var (
	endpointEchoService string
)

func init() {

	flag.StringVar(&endpointEchoService, "echo_service_endpoint", "", "endpoint host:port of EchoService")

}

type handler struct {
	mux   http.ServeMux
	conns map[string]*grpc.ClientConn
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) Close() error {
	var err error
	for svc, conn := range h.conns {
		cerr := conn.Close()
		if err == nil {
			err = cerr
		}
		if cerr != nil {
			glog.Errorf("Failed to close gRPC connection to %s: %v", svc, err)
		}
	}
	return err
}

func NewHandler(ctx context.Context) (http.Handler, error) {
	h := &handler{
		conns: make(map[string]*grpc.ClientConn),
	}
	var err error
	defer func() {
		if err != nil {
			h.Close()
		}
	}()

	err = func() error {
		conn, err := grpc.Dial(endpointEchoService)
		if err != nil {
			return err
		}
		h.conns["EchoService"] = conn
		client := NewEchoServiceClient(conn)

		h.mux.HandleFunc("/v1/example/echo", func(w http.ResponseWriter, req *http.Request) {
			resp, err := handle_EchoService_Echo(ctx, client, req)
			if err != nil {
				glog.Errorf("RPC error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf, err := proto.Marshal(resp)
			if err != nil {
				glog.Errorf("Marshal error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err = w.Write(buf); err != nil {
				glog.Errorf("Failed to write response: %v", err)
			}
		})

		return nil
	}()
	if err != nil {
		return nil, err
	}

	return h, nil
}
`,
	})
}
