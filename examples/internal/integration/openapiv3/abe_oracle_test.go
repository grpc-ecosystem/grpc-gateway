// abe_oracle_test.go is a Tier 2 oracle test for protoc-gen-openapiv3: it
// stands up a real grpc-gateway in front of the existing ABitOfEverything
// gRPC server impl and drives it through a Go client generated (via
// oapi-codegen) from the OpenAPI 3.1 spec our generator emits. A successful
// round-trip proves three things at once:
//
//  1. the generated spec is valid enough for oapi-codegen to consume,
//  2. the operations, parameters and schemas in the spec match what the
//     grpc-gateway actually accepts and returns on the wire,
//  3. the generated client and the gateway agree about JSON shapes.
//
// The generated client under examples/internal/clients/abev3 is regenerated
// by `make openapiv3-clients`; this test is the consumer.
//
// Tier 1 (abe_spec_test.go) checks structural shape of the checked-in spec
// without any codegen; this file is the opt-in complement that actually
// exercises the gateway for a small representative slice of operations.
package openapiv3_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/clients/abev3"
	examples "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// startABEGateway wires an ABitOfEverything gRPC server on a random
// loopback port and a grpc-gateway HTTP proxy in front of it, returning
// the HTTP base URL. It piggy-backs on the existing server impl under
// examples/internal/server so the oracle isn't testing a fake.
//
// The listener is opened once and handed directly to server.ServeGRPC.
func startABEGateway(t *testing.T) string {
	t.Helper()

	grpcLis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen grpc: %v", err)
	}
	grpcAddr := grpcLis.Addr().String()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	serveErr := make(chan error, 1)
	go func() { serveErr <- server.ServeGRPC(ctx, grpcLis) }()
	t.Cleanup(func() {
		cancel()
		// server.ServeGRPC closes the listener via grpc.Server.Serve's
		// GracefulStop path; waiting on the channel here guarantees the
		// goroutine has finished before the test exits so we can surface
		// unexpected errors.
		if err := <-serveErr; err != nil && ctx.Err() == nil {
			t.Errorf("gRPC server exited: %v", err)
		}
	})

	mux := runtime.NewServeMux()
	if err := examples.RegisterABitOfEverythingServiceHandlerFromEndpoint(
		ctx, mux, grpcAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		t.Fatalf("register gateway: %v", err)
	}

	httpSrv := httptest.NewServer(mux)
	t.Cleanup(httpSrv.Close)
	return httpSrv.URL
}

// newABEClient configures an abev3.ClientWithResponses pointed at baseURL.
func newABEClient(t *testing.T, baseURL string) *abev3.ClientWithResponses {
	t.Helper()
	client, err := abev3.NewClientWithResponses(baseURL, abev3.WithHTTPClient(&http.Client{}))
	if err != nil {
		t.Fatalf("construct client: %v", err)
	}
	return client
}

// TestABEOracle_Echo exercises a simple GET with a {value} path param. The
// server impl echoes the string back, so the oracle is asserting that the
// spec's path template and the gateway agree on where `value` lands.
func TestABEOracle_Echo(t *testing.T) {
	baseURL := startABEGateway(t)
	client := newABEClient(t, baseURL)

	resp, err := client.ABitOfEverythingServiceEchoWithResponse(context.Background(), "hello-oracle")
	if err != nil {
		t.Fatalf("Echo: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("Echo status: got %d, want 200 (body=%s)", resp.StatusCode(), string(resp.Body))
	}
	if resp.JSON200 == nil || resp.JSON200.Value == nil || *resp.JSON200.Value != "hello-oracle" {
		t.Errorf("Echo response: got %+v, want value=hello-oracle", resp.JSON200)
	}
}

// TestABEOracle_CheckStatus exercises an RPC that returns a message wrapping
// google.rpc.Status — this verifies that the inlined Status component schema
// is compatible with the type the server actually writes on the wire.
func TestABEOracle_CheckStatus(t *testing.T) {
	baseURL := startABEGateway(t)
	client := newABEClient(t, baseURL)

	resp, err := client.ABitOfEverythingServiceCheckStatusWithResponse(context.Background())
	if err != nil {
		t.Fatalf("CheckStatus: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("CheckStatus status: got %d (body=%s)", resp.StatusCode(), string(resp.Body))
	}
	if resp.JSON200 == nil {
		t.Fatal("CheckStatus: nil JSON200")
	}
}

// TestABEOracle_CreateLookupDelete is the full CRUD smoke path: POST a new
// ABE with body="*", GET it back by {uuid}, and DELETE it (which returns
// google.protobuf.Empty → 200 with `{}`, matching runtime behavior).
// This covers body="*" inline schema synthesis, path params, and the
// WKT Empty schema path.
func TestABEOracle_CreateLookupDelete(t *testing.T) {
	baseURL := startABEGateway(t)
	client := newABEClient(t, baseURL)
	ctx := context.Background()

	wantString := "oracle-string"
	wantInt64 := "123456789"
	createResp, err := client.ABitOfEverythingServiceCreateBodyWithResponse(ctx, abev3.ABitOfEverythingServiceCreateBodyJSONRequestBody{
		StringValue: &wantString,
		Int64Value:  &wantInt64,
	})
	if err != nil {
		t.Fatalf("CreateBody: %v", err)
	}
	if createResp.StatusCode() != http.StatusOK {
		t.Fatalf("CreateBody status: got %d (body=%s)", createResp.StatusCode(), string(createResp.Body))
	}
	if createResp.JSON200 == nil || createResp.JSON200.Uuid == nil || *createResp.JSON200.Uuid == "" {
		t.Fatalf("CreateBody: missing uuid in response: %+v", createResp.JSON200)
	}
	uuid := *createResp.JSON200.Uuid

	lookupResp, err := client.ABitOfEverythingServiceLookupWithResponse(ctx, uuid)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if lookupResp.StatusCode() != http.StatusOK {
		t.Fatalf("Lookup status: got %d (body=%s)", lookupResp.StatusCode(), string(lookupResp.Body))
	}
	if lookupResp.JSON200 == nil || lookupResp.JSON200.StringValue == nil || *lookupResp.JSON200.StringValue != wantString {
		t.Errorf("Lookup: stringValue = %+v, want %q", lookupResp.JSON200, wantString)
	}

	deleteResp, err := client.ABitOfEverythingServiceDeleteWithResponse(ctx, uuid)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	// Delete returns google.protobuf.Empty. The runtime writes `{}` with a
	// 200 regardless of the proto response type, so the spec documents
	// 200 and the oracle enforces that.
	if deleteResp.StatusCode() != http.StatusOK {
		t.Errorf("Delete status: got %d, want 200 (body=%s)", deleteResp.StatusCode(), string(deleteResp.Body))
	}
}

// TestABEOracle_DeepPathEcho exercises a POST with a dotted nested path
// parameter {single_nested.name} and body="*". This is the interesting
// combination the v2 generator has historically fumbled: the path component
// and the body both touch the same request type.
func TestABEOracle_DeepPathEcho(t *testing.T) {
	baseURL := startABEGateway(t)
	client := newABEClient(t, baseURL)

	wantName := "oracle-deep"
	amount := 42
	resp, err := client.ABitOfEverythingServiceDeepPathEchoWithResponse(
		context.Background(),
		wantName,
		abev3.ABitOfEverythingServiceDeepPathEchoJSONRequestBody{
			SingleNested: &abev3.GrpcGatewayExamplesInternalProtoExamplepbABitOfEverythingNested{
				Amount: func() *int64 { v := int64(amount); return &v }(),
			},
		},
	)
	if err != nil {
		t.Fatalf("DeepPathEcho: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("DeepPathEcho status: got %d (body=%s)", resp.StatusCode(), string(resp.Body))
	}
	if resp.JSON200 == nil || resp.JSON200.SingleNested == nil {
		t.Fatalf("DeepPathEcho: missing singleNested, got %+v", resp.JSON200)
	}
	// The path param overrides what's in the body — the gateway sets
	// single_nested.name from the URL, so the echoed response reflects
	// the path value.
	if resp.JSON200.SingleNested.Name == nil || *resp.JSON200.SingleNested.Name != wantName {
		t.Errorf("DeepPathEcho: singleNested.name = %+v, want %q", resp.JSON200.SingleNested.Name, wantName)
	}
}
