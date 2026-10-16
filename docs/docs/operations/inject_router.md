---
layout: default
title: Adding custom routes to the mux
nav_order: 0
parent: Operations
---

# Adding custom routes to the mux

The gRPC-Gateway allows you to add custom routes to the serve mux, for example, if you want to support a use case that isn't supported by the gRPC-Gateway, like file uploads.

By default, a deep wildcard in a custom route can match zero path segments, so `/v3/files/{path=**}` also matches `/v3/files`.
To require a remaining path component after `/v3/files`, configure the mux with `runtime.WithStrictDeepWildcardMatching()`:

```go
mux := runtime.NewServeMux(runtime.WithStrictDeepWildcardMatching())
```

The option applies to both generated routes and routes registered with `HandlePath`.
A trailing slash still supplies an empty component, so `/v3/files/` continues to match.
The `strict_deep_wildcard` generator flag remains available when strict matching should be encoded only in generated patterns.

## Example

```go
package main

import (
	"context"
	"net/http"

	pb "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	ctx := context.TODO()
	mux := runtime.NewServeMux()
	// Register generated routes to mux
	err := pb.RegisterGreeterHandlerServer(ctx, mux, &GreeterServer{})
	if err != nil {
		panic(err)
	}
	// Register custom route for  GET /hello/{name}
	err = mux.HandlePath("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("hello " + pathParams["name"]))
	})
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", mux)
}

// GreeterServer is the server API for Greeter service.
type GreeterServer struct {

}

// SayHello implement to say hello
func (h *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
```
