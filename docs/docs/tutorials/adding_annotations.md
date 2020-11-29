---
layout: default
title: Adding the grpc-gateway annotations to an existing protobuf file
parent: Tutorials
nav_order: 5
---

# Adding the grpc-gateway annotations to an existing protobuf file

Now that we've got a working Go gRPC server, we need to add the grpc-gateway annotations.

The annotations define how gRPC services map to the JSON request and response. When using protocol buffers, each RPC must define the HTTP method and path using the `google.api.http` annotation.

So you will need to add `import "google/api/http.proto";` to the gRPC proto file.

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

HttpRule is typically specified as an `google.api.http` annotation on the gRPC method. Each mapping specifies a URL path template and an HTTP method.

Also, See [a_bit_of_everything.proto](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/proto/examplepb/a_bit_of_everything.proto) for examples of more annotations you can add to customize gateway behavior and generated OpenAPI output.

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

Usage examples can be found on this [Usage](https://github.com/grpc-ecosystem/grpc-gateway#usage)

### In addition to the main.go

```go
package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	helloworldpb "github.com/iamrajiv/helloworld/proto/helloworld"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(dialAddr string) error {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	gwmux := runtime.NewServeMux()
	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "11000"
	}
	gatewayAddr := "0.0.0.0:" + port
	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "") {
				gwmux.ServeHTTP(w, r)
				return
			}
		}),
	}
	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}


```

For more refer to this boilerplate repository [grpc-gateway-boilerplate
](https://github.com/johanbrandhorst/grpc-gateway-boilerplate)

## Testing the gRPC-Gateway

Then we use curl to send HTTP requests:

```sh
$ curl -X POST -k http://localhost:8080/v1/example/echo -d '{"name": " Hello"}'
```

```
{"message":"Hello  World"}
```

The process is as follows:

`curl` sends a request to the gateway with the post, gateway as proxy forwards the request to greeter_server through grpc, greeter_server returns the result through grpc, the gateway receives the result, and JSON returns to the front end.

In this way, the transformation process from HTTP JSON to internal grpc is completed through gRPC-Gateway.

[Next](learn_more.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
