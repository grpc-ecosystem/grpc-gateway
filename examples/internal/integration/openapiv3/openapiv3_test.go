// Package openapiv3 contains an end-to-end oracle test for
// protoc-gen-openapiv3: it stands up a real grpc-gateway HTTP proxy in front
// of a tiny in-process Greeter gRPC server, then drives it with a Go client
// that was produced by oapi-codegen from the spec our generator emits. A
// successful round-trip proves three things at once:
//
//  1. the generated OpenAPI 3.1 spec is valid enough for oapi-codegen to
//     consume,
//  2. the operations, parameters and schemas in the spec match what the
//     grpc-gateway actually accepts and returns on the wire,
//  3. the generated client and the gateway agree about JSON shapes.
//
// The Go client under examples/internal/clients/helloworldv3 is regenerated
// by `make generate`; this test is the consumer.
package openapiv3_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/helloworldv3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type greeterServer struct {
	helloworld.UnimplementedGreeterServer
}

func (greeterServer) SayHello(_ context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if v := req.GetInt32Val(); v != nil {
		return &helloworld.HelloReply{Message: fmt.Sprintf("int32=%d", v.GetValue())}, nil
	}
	return &helloworld.HelloReply{Message: "hello " + req.GetName()}, nil
}

// startGateway wires a Greeter gRPC server on a random loopback port and a
// grpc-gateway HTTP proxy in front of it, returning the HTTP base URL.
func startGateway(t *testing.T) string {
	t.Helper()

	grpcLis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen grpc: %v", err)
	}
	t.Cleanup(func() { _ = grpcLis.Close() })

	grpcSrv := grpc.NewServer()
	helloworld.RegisterGreeterServer(grpcSrv, greeterServer{})
	go func() { _ = grpcSrv.Serve(grpcLis) }()
	t.Cleanup(grpcSrv.Stop)

	mux := runtime.NewServeMux()
	if err := helloworld.RegisterGreeterHandlerFromEndpoint(
		context.Background(), mux, grpcLis.Addr().String(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		t.Fatalf("register gateway: %v", err)
	}

	httpSrv := httptest.NewServer(mux)
	t.Cleanup(httpSrv.Close)
	return httpSrv.URL
}

// newHelloworldClient configures a helloworldv3.ClientWithResponses pointed
// at baseURL.
func newHelloworldClient(t *testing.T, baseURL string) *helloworldv3.ClientWithResponses {
	t.Helper()
	client, err := helloworldv3.NewClientWithResponses(baseURL, helloworldv3.WithHTTPClient(&http.Client{}))
	if err != nil {
		t.Fatalf("construct client: %v", err)
	}
	return client
}

// TestOracle_SayHello exercises the primary {name} binding. The generated
// client's path-param method is the oracle: if the spec says the path is
// /say/{name}, the client synthesizes that URL, and if the gateway disagrees
// the request 404s and the test fails.
func TestOracle_SayHello(t *testing.T) {
	baseURL := startGateway(t)
	client := newHelloworldClient(t, baseURL)

	resp, err := client.GreeterSayHelloWithResponse(context.Background(), "Alice", nil)
	if err != nil {
		t.Fatalf("SayHello: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("SayHello status: got %d, want 200 (body=%s)", resp.StatusCode(), string(resp.Body))
	}
	if resp.JSON200 == nil || resp.JSON200.Message == nil {
		t.Fatalf("SayHello: missing message in response: %+v", resp.JSON200)
	}
	if got, want := *resp.JSON200.Message, "hello Alice"; got != want {
		t.Errorf("message: got %q, want %q", got, want)
	}
}

// TestOracle_SayHello_Int32ValBinding exercises one of helloworld.proto's
// additional_bindings paths (/say/int32val/{int32Val}). The v3 generator
// disambiguates additional bindings by appending `_<idx>` to the operationId;
// oapi-codegen reflects that as a numbered client method. Hitting that
// binding proves the disambiguation keeps the spec actionable.
func TestOracle_SayHello_Int32ValBinding(t *testing.T) {
	baseURL := startGateway(t)
	client := newHelloworldClient(t, baseURL)

	resp, err := client.GreeterSayHello6WithResponse(context.Background(), 42, nil)
	if err != nil {
		t.Fatalf("SayHello6: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("SayHello6 status: got %d, want 200 (body=%s)", resp.StatusCode(), string(resp.Body))
	}
	if resp.JSON200 == nil || resp.JSON200.Message == nil {
		t.Fatalf("SayHello6: missing message in response: %+v", resp.JSON200)
	}
	if got, want := *resp.JSON200.Message, "int32=42"; got != want {
		t.Errorf("message: got %q, want %q", got, want)
	}
}
