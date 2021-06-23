---
layout: default
title: AnnotatedContext
nav_order: 4
parent: Operations
---

# AnnotatedContext

## Get HTTP Path pattern
1. Define HTTP path in proto option like below, following template ```google.api.http```.

```proto
syntax = "proto3";
option go_package = "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/proto/examplepb";
package grpc.gateway.examples.internal.proto.examplepb;

import "google/api/annotations.proto";

service LoginService {
  rpc Login (LoginRequest) returns (LoginReply) {
    option (google.api.http) = {
        post: "/v1/example/login"
        body: "*"
    };
  }
}

message LoginRequest {}

message LoginReply {}
```

2. Get values from annotated context e.g. in WithMetadata function.
You can pass data to backend by adding them to metadata or push them to metrics server.
```go
mux := runtime.NewServeMux(
	runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
		md := make(map[string]string)
		if method, ok := runtime.RPCMethod(ctx); ok {
			md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
		}
		if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
			md["pattern"] = pattern // /v1/example/login
		}
		return metadata.New(md)
	}),
)
```