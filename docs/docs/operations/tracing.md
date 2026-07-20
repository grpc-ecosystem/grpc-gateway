---
layout: default
title: Tracing
nav_order: 2
parent: Operations
---

# Tracing

If you are starting a new project, jump to the [OpenTelemetry](#opentelemetry)
section below. OpenTelemetry is the current standard for tracing in the Go
ecosystem and, when combined with `otelhttp` and `otelgrpc`, propagates spans
through the gateway with **no custom middleware**. The
[OpenCensus](#with-opencensusio-and-aws-x-ray) and
[OpenTracing](#opentracing-support-legacy) sections below are kept for existing
integrations and are considered legacy.

## With [OpenCensus.io](https://opencensus.io/) and [AWS X-ray](https://aws.amazon.com/xray/)

> **Legacy:** OpenCensus is deprecated in favor of OpenTelemetry. See the
> [OpenTelemetry](#opentelemetry) section for the current approach.

### Adding tracing using AWS-Xray as the exporter

This example uses the AWS-Xray exporter with a global trace setting. Note that AWS X-ray exporter does not handle any metrics only tracing.

1. Add the following imports

```go
xray "contrib.go.opencensus.io/exporter/aws"
"go.opencensus.io/plugin/ocgrpc"
"go.opencensus.io/plugin/ochttp"
"go.opencensus.io/trace"
```

2. Register the AWS X-ray exporter for the GRPC server

```go
xrayExporter, err := xray.NewExporter(
    xray.WithVersion("latest"),
    // Add your AWS region.
    xray.WithRegion("ap-southeast-1"),
)
if err != nil {
    // Handle any error.
}
// Do not forget to call Flush() before the application terminates.
defer xrayExporter.Flush()

// Register the trace exporter.
trace.RegisterExporter(xrayExporter)
```

3. Add a global tracing configuration

```go
// Always trace in this example.
// In production this can be set to a trace.ProbabilitySampler.
trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
```

4. Add `ocgrpc.ClientHandler` for tracing the gRPC client calls

```go
conn, err := grpc.NewClient(
    // Other options goes here.
    // Add ocgrpc.ClientHandler for tracing the grpc client calls.
    grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
)
```

5. Wrap the gateway mux with the OpenCensus HTTP handler

```go
gwmux := runtime.NewServeMux()

openCensusHandler := &ochttp.Handler{
		Handler: gwmux,
}

gwServer := &http.Server{
    Addr: "0.0.0.0:10000",
    Handler: openCensusHandler,
    }),
}
```

### Without a global configuration

In this example we have added the [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) and we do not wish to trace any health checks.

1. Follow step `1`, `2` and `4` from the previous section.

2. Since we are not using a global configuration we can decide what paths we want to trace.

```go
gwmux := runtime.NewServeMux()

openCensusHandler := &ochttp.Handler{
    Handler: gwmux,
    GetStartOptions: func(r *http.Request) trace.StartOptions {
        startOptions := trace.StartOptions{}
        if strings.HasPrefix(r.URL.Path, "/api") {
            // This example will always trace anything starting with /api.
            startOptions.Sampler = trace.AlwaysSample()
        }
        return startOptions
    },
}
```

4. No global configuration means we have to use the [per span sampler](https://opencensus.io/tracing/sampling/#per-span-sampler).

#### A method we want to trace

```go
func (s *service) Name(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // Here we add the span ourselves.
    ctx, span := trace.StartSpan(ctx, "name.to.use.in.trace", trace.
    // Select a sampler that fits your implementation.
    WithSampler(trace.AlwaysSample()))
    defer span.End()
    /// Other stuff goes here.
}
```

#### A method we do not wish to trace

```go
func (s *service) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
    // Note no span here.
    return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}
```

## OpenTracing Support (legacy)

> **Legacy:** OpenTracing is deprecated in favor of OpenTelemetry. The custom
> `tracingWrapper` middleware below is **only** needed for OpenTracing. With
> OpenTelemetry's `otelhttp` and `otelgrpc` no such wrapper is required — see the
> [OpenTelemetry](#opentelemetry) section.

If your project uses [OpenTracing](https://github.com/opentracing/opentracing-go) and you'd like spans to propagate through the gateway, you can add some middleware which parses the incoming HTTP headers to create a new span correctly.

```go
import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || err == opentracing.ErrSpanContextNotFound {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}

// Then just wrap the mux returned by runtime.NewServeMux() like this
if err := http.ListenAndServe(":8080", tracingWrapper(mux)); err != nil {
	log.Fatalf("failed to start gateway server on 8080: %v", err)
}
```

Finally, don't forget to add a tracing interceptor when registering
the services. E.g.

```go
import (
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
)

opts := []grpc.DialOption{
	grpc.WithUnaryInterceptor(
		grpc_opentracing.UnaryClientInterceptor(
			grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
		),
	),
}
if err := pb.RegisterMyServiceHandlerFromEndpoint(ctx, mux, serviceEndpoint, opts); err != nil {
	log.Fatalf("could not register HTTP service: %v", err)
}
```

## OpenTelemetry

[OpenTelemetry](https://opentelemetry.io/) is the recommended way to trace
requests through the gateway. **You do not need a custom `tracingWrapper` (as in
the legacy OpenTracing section above).** The official instrumentation libraries
already do the work:

- [`otelhttp`](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp)
  wraps the gateway `http.Handler`, extracts the incoming
  [W3C Trace Context](https://www.w3.org/TR/trace-context/) from the request
  headers, and starts the **HTTP server span**.
- [`otelgrpc`](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)
  stats handlers start the **gRPC client span** for the gateway's outgoing call
  (injecting the trace context into the gRPC metadata) and the **gRPC server
  span** on the backend.

An incoming HTTP request therefore produces three connected parent/child spans
that share a single trace ID:

```
HTTP server span (otelhttp)  ->  gRPC client span (otelgrpc)  ->  gRPC server span (otelgrpc)
```

### Runnable example

A complete, self-contained example lives at
[`examples/internal/cmd/example-otel-tracing`](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/examples/internal/cmd/example-otel-tracing).
In a single process it wires up a stdout span exporter, an instrumented
in-process gRPC backend, and the gateway. Run it with:

```shell
go run ./examples/internal/cmd/example-otel-tracing
# then, in another terminal:
curl http://localhost:8081/say/OpenTelemetry
```

### Wiring

Install a `TracerProvider` and a propagator once at startup:

```go
import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
if err != nil {
	// Handle any error.
}
tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
defer tp.Shutdown(context.Background())

otel.SetTracerProvider(tp)
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
	propagation.TraceContext{},
	propagation.Baggage{},
))
```

Instrument the gRPC backend with the server stats handler:

```go
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
```

Dial the backend with the client stats handler so the gateway's outgoing call is
traced and the trace context is propagated:

```go
conn, err := grpc.NewClient(
	"localhost:9091",
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
)
```

Finally, wrap the gateway mux with `otelhttp` — this replaces the custom
`tracingWrapper` entirely:

```go
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

mux := runtime.NewServeMux()
if err := pb.RegisterGreeterHandler(ctx, mux, conn); err != nil {
	// Handle any error.
}

handler := otelhttp.NewHandler(mux, "grpc-gateway")
if err := http.ListenAndServe(":8081", handler); err != nil {
	log.Fatalf("failed to start gateway server: %v", err)
}
```

### Sample output

Issuing `curl http://localhost:8081/say/OpenTelemetry` prints three spans that
all share the same `TraceID`, confirming the trace propagates through the
gateway. Trimmed for brevity:

```jsonc
// 1. HTTP server span (otelhttp) — the root span.
{
	"Name": "GET",
	"SpanContext": { "TraceID": "6952...4aa5", "SpanID": "3b9c...fd0e" },
	"Parent":      { "TraceID": "0000...0000", "SpanID": "0000...0000" },
	"SpanKind": 2, // server
	"InstrumentationScope": { "Name": ".../net/http/otelhttp" }
}
// 2. gRPC client span (otelgrpc) — child of the HTTP span.
{
	"Name": "grpc.gateway.examples.internal.helloworld.Greeter/SayHello",
	"SpanContext": { "TraceID": "6952...4aa5", "SpanID": "1224...75ca" },
	"Parent":      { "TraceID": "6952...4aa5", "SpanID": "3b9c...fd0e" },
	"SpanKind": 3, // client
	"InstrumentationScope": { "Name": ".../grpc/otelgrpc" }
}
// 3. gRPC server span (otelgrpc) — child of the client span, context received over the wire.
{
	"Name": "grpc.gateway.examples.internal.helloworld.Greeter/SayHello",
	"SpanContext": { "TraceID": "6952...4aa5", "SpanID": "950f...d85e" },
	"Parent":      { "TraceID": "6952...4aa5", "SpanID": "1224...75ca", "Remote": true },
	"SpanKind": 2, // server
	"InstrumentationScope": { "Name": ".../grpc/otelgrpc" }
}
```
