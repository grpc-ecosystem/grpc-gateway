---
layout: default
title: Creating a basic protos with HTTP annotations
parent: Tutorials
nav_order: 5
---

## Creating a basic protos with HTTP annotations

### Annotating your gRPC .proto file with HTTP bindings and routes

The annotations define how gRPC services map to the JSON request and response. When using protocol buffers, each RPC must define the HTTP method and path using the `google.api.http` annotation.

So you will need to add `import "google/api/http.proto";` to the gRPC proto file.

Now lets add the HTTP annotations to our proto file.

```proto
syntax = "proto3";

package helloworld;

import "google/api/annotations.proto";

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {
        option (google.api.http) = {
        post: "/v1/example/echo"
        body: "*"
    };
  }
}
```
