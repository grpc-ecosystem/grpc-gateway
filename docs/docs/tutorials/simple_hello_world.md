---
layout: default
title: Creating a simple hello world with gRPC-Gateway
parent: Tutorials
nav_order: 4
---

### Creating a simple hello world with gRPC-Gateway

### gRPC-Gateway

We all know that gRPC is not a tool for everything. There are cases where we still want to provide a traditional RESTful JSON API. The reasons can range from maintaining backwards-compatibility to supporting programming languages or clients not well supported by gRPC. But coding another API for REST is quite a time consuming and tedious.

So is there any way to code just once, but can provide APIs in both gRPC and REST at the same time?

The answer is Yes.

### Usage

One way to achieve that is to use gRPC gateway. gRPC gateway is a plugin of the protocol buffer compiler. It reads the protobuf service definitions and generates a proxy server, which translates a RESTful HTTP call into gRPC request.

All we need to do is a small amount of configuration
in the service. Before start coding, we have to install some
tools.

Use `go get -u` to download the following packages:

```sh
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Define and Generate proto files

Before we create a gRPC service, we should create a proto file to define what we need, here we create a file named `hello.proto` to show.

```proto
syntax = "proto3";

option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";

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
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
helloworld/helloworld.proto
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

The gateway code is as follows:

```go
package main

import (
    "flag"
    "net/http"

    "github.com/golang/glog"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "golang.org/x/net/context"
    "google.golang.org/grpc"

    gw "grpc-helloworld-gateway/helloworld"
)

var (
    echoEndpoint = flag.String("echo_endpoint", "localhost:50051", "endpoint of YourService")
)

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
    if err != nil {
        return err
    }

    return http.ListenAndServe(":8080", mux)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if err := run(); err != nil {
        glog.Fatal(err)
    }
}
```

First echoEndpoint stores the server information that needs to be connected, and then registers and binds the new server with RegisterGreeter Handler FromEndpoint in `gw.go`. Then the lower level will connect to the remote server address provided by echoEndpoint, so the gateway establishes the connection as the client and remote server, and then starts the new server with http, gateway. It serves as a server to provide HTTP services to the outside world.

This is the end of the code. Let's test it.

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

### More

#### google.api.http

GRPC transcoding is a conversion function between the gRPC method and one or more HTTP REST endpoints. This allows developers to create a single API service that supports both the gRPC API and the REST API. Many systems, including the API Google, Cloud Endpoints, gRPC Gateway, and the Envoy proxy server support this feature and use it for large-scale production services.

The grcp-gateway the server is created according to the `google.api.http` annotations in your service definitions.

HttpRule defines the gRPC / REST mapping scheme. The mapping defines how different parts of a gRPC request message are mapped to the URL path, URL request parameters, and HTTP request body. It also controls how the gRPC response message is displayed in the HTTP response body. HttpRule is usually specified as a `google.api.http` annotation in the gRPC method.

Each mapping defines a URL path template and an HTTP method. A path template can refer to one or more fields in a gRPC request message if each field is a non-repeating field with a primitive type. The path template controls how the request message fields are mapped to the URL path.

```proto
import "google/api/annotations.proto";
import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
message StatusResponse {
 google.protobuf.Timestamp current_time = 1;
}
service MyService {
 rpc Status(google.protobuf.Empty)
  returns (StatusResponse) {
   option (google.api.http) = {
    get: "/status"
   };
  }
 }
}
```

You will need to provide the necessary third-party `protobuf` files to the `protoc` compiler. They have included in the `grpc-gateway` repository in the `[third_party/googleapis](https://github.com/grpc-ecosystem/grpc-gateway/tree/master/third_party/googleapis)` folder, and we recommend copying them to the project file structure.

You will need to provide the necessary third-party `protobuf` files to the `protoc` compiler. They are included in the `[grpc-gateway` repository](https://github.com/grpc-ecosystem/grpc-gateway/) in the `[third_party/googleapis](https://github.com/grpc-ecosystem/grpc-gateway/tree/master/third_party/googleapis)` folder and we recommend copying them to the project file structure.
