---
layout: default
title: Creating a basic protos with HTTP annotations
parent: Tutorials
nav_order: 3
---

## Creating a basic protos with HTTP annotations

### Annotating your gRPC .proto file with HTTP bindings and routes

The annotations define how gRPC services map to the JSON request and response. When using protocol buffers, each RPC must define the HTTP method and path using the `google.api.http` annotation.

So you will need to add `import "google/api/hello_world.proto";` to the gRPC proto file.

To make a Basic Arithmetic service first we have to define our services in proto files.

```proto
syntax="proto3";

package example;

import "google/api/hello_world.proto";
import "protoc-gen-openapiv2/options/hello_world.proto";

// Defines the import path that should be used to import the generated package,
// and the package name.
option go_package = "github.com/iamrajiv/Basic-Arithmetic-gRPC-Server/proto;example";

// These annotations are used when generating the OpenAPI file.
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    version: "1.0";
  };
  external_docs: {
    url: "https://github.com/iamrajiv/Basic-Arithmetic-gRPC-Server";
    description: "Basic-Arithmetic-gRPC-Server";
  }
  schemes: HTTPS;
};

// This will be our request
message Request {
// First integer in the body of the request
  int64 a = 1;
// Second integer in the body of the request
  int64 b = 2;
}

// This will be our response
message Response {
// A integer in the body of the response
  int64 result = 1;
}

// Here is the overall service where we define all our endpoints
service ArithmeticService {
// Here is our Add method to use as an service
  rpc Add(Request) returns (Response){
    option (google.api.http) = {
      post: "/api/v1/arithmetic/add"
      body: "*"
    };
  };
  // Here is our Divide method to use as an service
  rpc Divide(Request) returns (Response){
    option (google.api.http) = {
      post: "/api/v1/arithmetic/div"
      body: "*"
    };
  };
  // Here is our Multiply method to use as an service
  rpc Multiply(Request) returns (Response){
    option (google.api.http) = {
      post: "/api/v1/arithmetic/mul"
      body: "*"
    };
  };
// Here is our Subtract method to use as an service
  rpc Subtract(Request) returns (Response){
    option (google.api.http) = {
      post: "/api/v1/arithmetic/sub"
      body: "*"
    };
  };
}
```
