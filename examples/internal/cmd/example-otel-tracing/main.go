/*
Command example-otel-tracing is a self-contained, runnable example that shows how
to propagate OpenTelemetry traces through grpc-gateway.

In a single process it wires up:

  - an OpenTelemetry TracerProvider that exports spans to stdout, plus a W3C
    TraceContext propagator;
  - an in-process gRPC backend (the helloworld Greeter) instrumented with the
    otelgrpc server stats handler;
  - a gRPC client dialed with the otelgrpc client stats handler;
  - runtime.NewServeMux wrapped with otelhttp.NewHandler.

A single incoming HTTP request therefore produces three connected spans that
share one trace ID, exported to stdout:

	HTTP server span (otelhttp) -> gRPC client span (otelgrpc) -> gRPC server span (otelgrpc)

Because otelhttp extracts the incoming trace context and otelgrpc injects it into
the outgoing gRPC call, no custom tracingWrapper middleware is required.

Run it and issue a request:

	go run ./examples/internal/cmd/example-otel-tracing
	curl http://localhost:8081/say/OpenTelemetry
*/
package main

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

const (
	httpAddr = ":8081"
	grpcAddr = "localhost:9091"
)

// greeterServer implements the helloworld Greeter service.
type greeterServer struct {
	helloworld.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(_ context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + req.GetName()}, nil
}

// initTracer configures a TracerProvider that exports spans to stdout and installs
// a W3C TraceContext (plus Baggage) propagator. The returned function flushes and
// shuts the provider down.
func initTracer() (func(context.Context) error, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// WithSyncer flushes every span as soon as it ends so the connected spans are
	// printed to stdout immediately after each request.
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// runGRPCServer starts the instrumented gRPC backend and blocks serving it.
func runGRPCServer(lis net.Listener) error {
	// The server stats handler starts a gRPC server span as a child of the
	// propagated client span.
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	helloworld.RegisterGreeterServer(s, &greeterServer{})
	return s.Serve(lis)
}

func main() {
	ctx := context.Background()

	shutdown, err := initTracer()
	if err != nil {
		grpclog.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			grpclog.Errorf("failed to shut down tracer: %v", err)
		}
	}()

	// Start the in-process gRPC backend.
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		grpclog.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}
	go func() {
		if err := runGRPCServer(lis); err != nil {
			grpclog.Fatalf("failed to serve gRPC backend: %v", err)
		}
	}()

	// Dial the backend with the otelgrpc client stats handler so the gateway's
	// outgoing gRPC call carries the trace context and emits a client span.
	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		grpclog.Fatalf("failed to dial gRPC backend: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			grpclog.Errorf("failed to close gRPC connection: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	if err := helloworld.RegisterGreeterHandler(ctx, mux, conn); err != nil {
		grpclog.Fatalf("failed to register gateway handler: %v", err)
	}

	// otelhttp.NewHandler extracts the incoming W3C trace context from the request
	// headers and starts the HTTP server span. No custom tracingWrapper is needed.
	handler := otelhttp.NewHandler(mux, "grpc-gateway")

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: handler,
	}

	grpclog.Infof("gRPC backend listening on %s", grpcAddr)
	grpclog.Infof("gateway listening on %s", httpAddr)
	grpclog.Infof("try: curl http://localhost%s/say/OpenTelemetry", httpAddr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		grpclog.Fatalf("failed to serve gateway: %v", err)
	}
}
