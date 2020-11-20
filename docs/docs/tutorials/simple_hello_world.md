---
layout: default
title: Creating a simple hello world with gRPC-Gateway
parent: Tutorials
nav_order: 5
---

### Creating a simple hello world with gRPC-Gateway

### gRPC-Gateway

We all know that gRPC is not a tool for everything. There are cases where we still want to provide a traditional RESTful JSON API. The reasons can range from maintaining backwards-compatibility to supporting programming languages or clients not well supported by gRPC. But coding another API for REST is quite a time consuming and tedious.

So is there any way to code just once, but can provide APIs in both gRPC and REST at the same time?

The answer is Yes.

### Define and Generate proto files

Before we create a gRPC service, we should create a proto file to define what we need, here we create a file named `hello.proto` to show.

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

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

Then we should generate gRPC stub via protoc

```sh
protoc -I . \
   --go_out ./gen/go/ --go_opt paths=source_relative \
   --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative \
   your/service/v1/hello_world.proto
```

The `helloworld.pb.go` file will be generated. So `helloworld.pb.go` is required by the server service.

Next we need to use protoc to generate the go files needed by the gateway.

```sh
protoc -I/usr/local/include -I. \
-I$GOPATH/src  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--swagger_out=logtostderr=true:. \
helloworld/helloworld.proto
```

The `helloworld.pb.gw.go` file will be generated. This file is the protocol file used by gateway for protocol conversion between grpc and http.

After the protocol file is processed, the gateway code is needed.

### Implementation

For `main.go` file of hello world we can refer to the boilerplate repository [grpc-gateway-boilerplate
Template](https://github.com/johanbrandhorst/grpc-gateway-boilerplate).
