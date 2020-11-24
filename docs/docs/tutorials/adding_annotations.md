---
layout: default
title: Adding the grpc-gateway annotations to an existing protobuf file
parent: Tutorials
nav_order: 5
---

## Adding the grpc-gateway annotations to an existing protobuf file

Now that we've got a working Go gRPC server, we need to add the grpc-gateway annotations:

```proto
syntax = "proto3";

package helloworld;

import "google/api/annotations.proto";

// Here is the overall greeting service definition where we define all our endpoints
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

Also, See [a_bit_of_everything.proto](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/proto/examplepb/a_bit_of_everything.proto) for examples of more annotations you can add to customize gateway behavior and generated OpenAPI output.

### Generating the grpc-gateway stubs

Now that we've got the grpc-gateway annotations added to the proto file, we need to use the grpc-gateway generator to generate the stubs.

Before we can do that, we need to copy some dependencies into our protofile structure. Copy the `third_party/googleapis` folder from the grpc-gateway repository to your local protofile structure. It should look like this afterwards:

```
proto/
helloworld/
hello_world.proto
google/
api/
http.proto
annotations.proto
```

#### Using buf

We'll need to add the grpc-gateway generator to the generation configuration:

```yml
plugins:
  - name: go
    out: proto
    opt: paths=source_relative
  - name: go-grpc
    out: proto
    opt: paths=source_relative
  - name: grpc-gateway
    out: proto
    opt: paths=source_relative
```

And that's it! Now if you run:

```sh
buf generate
```

It should produce a `*.gw.pb.go` file.

#### Using protoc

Now we need to add the grpc-gateway generator to the protoc invocation:

```
protoc -I ./proto \
 ... other plugins ...
--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative
./proto/helloworld/hello_world.proto
```

```sh
protoc -I ./proto \
   --go_out ./proto --go_opt paths=source_relative \
   --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
   ./proto/helloworld/hello_world.proto
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative
   ./proto/helloworld/hello_world.proto
```

This should generate a `*.gw.pb.go` file.

### Testing the grpc-gateway

Then we use curl to send http requests:

```sh
curl -X POST -k http://localhost:8080/v1/example/echo -d '{"name": " Hello"}'
```

```
{"message":"Hello  World"}
```

The process is as follows:

`curl` sends a request to the gateway with the post, gateway as proxy forwards the request to greeter_server through grpc, greeter_server returns the result through grpc, the gateway receives the result, and json returns to the front end.

In this way, the transformation process from http json to internal grpc is completed through grpc-gateway.
