---
layout: default
title: Creating a simple hello world with gRPC
parent: Tutorials
nav_order: 2
---

## Creating a simple hello world with gRPC

To understand gRPC-Gateway we are going to make hello world gRPC service which uses gRPC-Gateway.

### Creating a basic protos with HTTP annotations

The annotations define how gRPC services map to the JSON request and response. When using protocol buffers, each RPC must define the HTTP method and path using the `google.api.http` annotation.

So you will need to add `import "google/api/http.proto";` to the gRPC proto file.

Now, let's add the HTTP annotations to our proto file.

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

### Defining your gRPC service using protocol buffers

Before we create a gRPC service, we should create a proto file to define what we need, here we create a file named `hello_world.proto` in the directory `proto/helloworld/hello_world.proto`.

```proto
syntax = "proto3";

package helloworld;

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```
