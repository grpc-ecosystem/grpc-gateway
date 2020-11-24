---
layout: default
title: Creating main.go
parent: Tutorials
nav_order: 4
---

## Creating main.go

```go
package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	// Static files

	"github.com/hello/hello_world/gateway"
	"github.com/hello/hello_world/insecure"
	hello "github.com/hello/hello_world/proto"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {

	log.Println("request: ", in.Name)
	return &hello.HelloReply{Message: in.Name + " World"}, nil
}

func main() {

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)
	hello.RegisterGreeterServer(s, &server{})

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///" + addr)
	log.Fatalln(err)
}
```

### Read More

For more refer to this boilerplate repository [grpc-gateway-boilerplate Template](https://github.com/johanbrandhorst/grpc-gateway-boilerplate).
