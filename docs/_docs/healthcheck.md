---
category: documentation
---

# Health check

## With the [GRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

To use the GRPC health checking protocol you must add the two health checking methods, `Watch` and `Check`.

#### Registering the health server

1. Add `google.golang.org/grpc/health/grpc_health_v1` to your imports
2. Register the health server with `grpc_health_v1.RegisterHealthServer(grpcServer, yourService)`

#### Adding the health check methods

1. Check method

```
func (s *serviceServer) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}
```

2. Watch method 

```
func (s *serviceServer) Watch(in *health.HealthCheckRequest, _ health.Health_WatchServer) error {
    // Example of how to register both methods but only implement the Check method.
	return status.Error(codes.Unimplemented, "unimplemented")
}
```

3. You can test the functionality with [GRPC health probe](https://github.com/grpc-ecosystem/grpc-health-probe)
