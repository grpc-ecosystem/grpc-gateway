---
layout: default
title: Adding the grpc-gateway annotations to an existing protobuf file
parent: Tutorials
nav_order: 5
---

# Adding the grpc-gateway annotations to an existing protobuf file

Now that we've got a working Go gRPC server, we need to add the grpc-gateway annotations.

The annotations define how gRPC services map to the JSON request and response. When using protocol buffers, each RPC must define the HTTP method and path using the `google.api.http` annotation.

So we will need to add the `google/api/http.proto` import to the proto file. We also need to add the HTTP->gRPC mapping we want. In this case, we're mapping `POST /v1/example/echo` to our `SayHello` rpc.

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

// The request message containing the user's name
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

See [a_bit_of_everything.proto](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/proto/examplepb/a_bit_of_everything.proto) for examples of more annotations you can add to customize gateway behavior.

## Generating the grpc-gateway stubs

Now that we've got the grpc-gateway annotations added to the proto file, we need to use the grpc-gateway generator to generate the stubs.

Before we can do that, we need to copy some dependencies into our protofile structure. Copy the `third_party/googleapis` folder from the grpc-gateway repository to your local protofile structure. It should look like this afterwards:

```
proto
├── google
│   └── api
│       ├── annotations.proto
│       └── http.proto
└── helloworld
    └── hello_world.proto
```

### Using buf

We'll need to add the grpc-gateway generator to the generation configuration:

```yml
version: v1beta1
plugins:
  - name: go
    out: proto
    opt: paths=source_relative
  - name: go-grpc
    out: proto
    opt: paths=source_relative,require_unimplemented_servers=false
  - name: grpc-gateway
    out: proto
    opt: paths=source_relative
```

And that's it! Now if you run:

```sh
$ buf generate
```

It should produce a `*.gw.pb.go` file.

### Using protoc

Now we need to add the grpc-gateway generator to the protoc invocation:

```sh
$ protoc -I ./proto \
  --go_out ./proto --go_opt paths=source_relative \
  --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
  ./proto/helloworld/hello_world.proto
```

This should generate a `*.gw.pb.go` file.

We also need to add and serve the gRPC-gateway mux in our `main.go` file.

```go
package main
import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"log"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	helloworldpb "github.com/myuser/myrepo/proto/helloworld"
)
type server struct{}
func NewServer() *server {
	return &server{}
}
func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	return &helloworldpb.HelloReply{Message: in.Name + " World"}, nil
}
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer()
	helloworldpb.RegisterGreeterServer(s, &server{})
	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr: ":8090",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
```

For more examples, please refer to [our boilerplate repository](https://github.com/johanbrandhorst/grpc-gateway-boilerplate).

## Testing the gRPC-Gateway

Now we can start the server:

```sh
$ go run main.go
```

Then we use cURL to send HTTP requests:

```sh
$ curl -X POST -k http://localhost:8090/v1/example/echo -d '{"name": " Hello"}'
```

```
{"message":"Hello  World"}
```

The process is as follows:

Hopefully, that gives a bit of understanding of how to use the gRPC-Gateway.

[Next](learn_more.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
