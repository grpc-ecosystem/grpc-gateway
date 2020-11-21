---
layout: default
title: Adding the grpc-gateway annotations to an existing protobuf file
parent: Tutorials
nav_order: 6
---

## Adding the grpc-gateway annotations to an existing protobuf file

Now that we've got a working Go gRPC server, we need to add the grpc-gateway annotations:

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

### Generating the grpc-gateway stubs

Now that we've got the grpc-gateway annotations added to the proto file, we need to use the grpc-gateway generator to generate the stubs.

Before we can do that, we need to copy some dependencies into our protofile structure. Copy the `third_party/googleapis` folder from the grpc-gateway repository to your local protofile structure. It should look like this afterwards:

''' (backticks)
proto/
helloworld/
hello_world.proto
google/
api/
http.proto
annotations.proto
''' (backticks)

#### Using buf

[Using buf](generating_stubs/using_buf.md)

#### Using protoc

[Using protoc](generating_stubs/using_protoc.md)

Start the greeter_server service first, and then start the gateway. Then gateway connects to greeter_server and establishes http monitoring.

Then we use curl to send http requests:

```sh
curl -X POST -k http://localhost:8080/v1/example/echo -d '{"name": " world"}
```

```
{"message":"Hello  world"}
```

The process is as follows:

curl sends a request to the gateway with the post, gateway as proxy forwards the request to greeter_server through grpc, greeter_server returns the result through grpc, the gateway receives the result, and json returns to the front end.

In this way, the transformation process from http json to internal grpc is completed through grpc-gateway.
