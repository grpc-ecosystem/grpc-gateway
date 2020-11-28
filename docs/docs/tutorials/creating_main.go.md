---
layout: default
title: Creating main.go
parent: Tutorials
nav_order: 4
---

# Creating main.go

Before creating `main.go` file we are assuming that the user has created a `go.mod` with the name `github.com/myuser/myrepo`. Then import here is using the path to the generated files in `proto/helloworld` relative to the root of the repository.

```go
package main

import (
	"context"
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

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
	log.Info("Serving gRPC on 0.0.0.0:8080")
	log.Fatal(s.Serve(lis))
}
```

## Read More

For more refer to gRPC docs [https://grpc.io/docs/languages/go/](https://grpc.io/docs/languages/go/).

[Next](adding_annotations.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
