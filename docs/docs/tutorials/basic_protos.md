---
layout: default
title: Creating a basic protos with HTTP annotations
parent: Tutorials
nav_order: 2
---

## Creating a basic protos with HTTP annotations

### Annotating your gRPC .proto file with HTTP bindings and routes

The annotations define how gRPC services map to the JSON request and response. When using protocol buffers, each RPC must define the HTTP method and path using the `google.api.http` annotation.

So you will need to add `import "google/api/annotations.proto";` to the gRPC proto file.

To make a Basic Arithmetic service first we have to define our services in proto files.

```proto
syntax="proto3";

package example;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

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

### HTTP method and path

- The first key (`post` in this service) corresponds to the HTTP method. RPCs **may** use `get`, `post`, `patch`, or `delete`.
  - RPCs **must** use the prescribed HTTP verb for each standard method, as discussed in [AIP-131](https://google.aip.dev/131), [AIP-132](https://google.aip.dev/132), [AIP-133](https://google.aip.dev/133), [AIP-134](https://google.aip.dev/134), and [AIP-135](https://google.aip.dev/135)
  - RPCs **should** use the prescribed HTTP verb for custom methods, as discussed in [AIP-136](https://google.aip.dev/136).
  - RPCs **should not** use `put` or `custom`.
- The corresponding value represents the URI.
  - URIs **must** use the `{foo=bar/*}` syntax to represent a variable that should be populated in the request proto. When extracting a [resource name](https://google.aip.dev/122), the variable **must** include the entire resource name, not just the ID component.
  - URIs **must** use the `` character to represent ID components, which matches all URI-safe characters except for `/`. URIs **may** use `*` as the final segment of a URI if matching `/` is required.
- The `body` key defines which field in the request will be sent as the HTTP body. If the body is ``, then this indicates that the request object itself is the HTTP body. The request body is encoded as JSON as defined by protocol buffers' canonical [JSON encoding](https://developers.google.com/protocol-buffers/docs/proto3#json).
  - RPCs **must not** define a `body` at all for RPCs that use the `GET` or `DELETE` HTTP verbs.
  - RPCs **must** use the prescribed `body` for Create ([AIP-133](https://google.aip.dev/133)) and Update ([AIP-134](https://google.aip.dev/134)) requests.
  - RPCs **should** use the prescribed `body` for custom methods ([AIP-136](https://google.aip.dev/136)).
  - Fields **should not** use the `json_name` annotation to alter the field name in JSON, unless doing so for backwards-compatibility reasons.

You’ll see some atypical fields in the .proto that are leveraged by grpc-gateway. One of the most important of these fields is the option (google.api.http) where we define what HTTP URL(s) will be used to handle our request.

We also specify POST as the HTTP method we will accept.

Finally, if you have a request body that you expect (typical for POST requests and others), you must use the body field. If you don’t, then the request won’t be passed along to the handler.

Read more about HTTP and gRPC Transcoding on https://google.aip.dev/127.
