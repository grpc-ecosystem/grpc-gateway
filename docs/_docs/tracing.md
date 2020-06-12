---
category: documentation
---

# Tracing

## With [OpenCensus.io](https://opencenus.io) and [AWS X-ray](https://aws.amazon.com/xray/)

#### Adding tracing using AWS-Xray as the exporter

This example uses the AWS-Xray exporter with a global trace setting. Note that AWS X-ray exporter does not handle any metrics only tracing. 

1. Add the following imports

```
xray "contrib.go.opencensus.io/exporter/aws"
"go.opencensus.io/plugin/ocgrpc"
"go.opencensus.io/plugin/ochttp"
"go.opencensus.io/trace"
```

2. Register the AWS X-ray exporter for the GRPC server

```
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

```
// Always trace in this example.
// In production this can be set to a trace.ProbabilitySampler.
trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
```

4. Add `ocgrpc.ClientHandler` for tracing the grpc client calls

```
// Example using DialContext
conn, err := grpc.DialContext(
    // Other options goes here.
    // Add ocgrpc.ClientHandler for tracing the grpc client calls.
    grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
)
```

5. Wrap the gateway mux with the OpenCensus HTTP handler
```
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

#### Without a global configuration

In this example we have added the [GRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) and we do not wish to trace any health checks.

1. Follow step `1`, `2` and `4` from the previous section

2. Since we are not using a global configuration we can decide what paths we want to trace

```
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

4. No global configuration means we have to use the [per span sampler](https://opencensus.io/tracing/sampling/#per-span-sampler)


##### A method we __want__ to trace
```
func (s *service) Name(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // Here we add the span ourselves.
    ctx, span := trace.StartSpan(ctx, "name.to.use.in.trace", trace.
    // Select a sampler that fits your implementation.
    WithSampler(trace.AlwaysSample()))
    defer span.End()
    /// Other stuff goes here.
}
```

##### A method we __do not__ wish to trace
```
func (s *service) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
    // Note no span here.
    return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}
```
