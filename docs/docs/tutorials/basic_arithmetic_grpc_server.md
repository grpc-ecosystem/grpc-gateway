---
layout: default
title: Creating a Basic Arithmetic gRPC Server
parent: Tutorials
nav_order: 4
---

## Creating a Basic Arithmetic gRPC Server

Before we make `main.go` file first we have to make `server.go` file in a separate folder.

At the top of the file, we define the handlers for our gRPC methods.

Handlers are:

```go
// This is the struct that we will implement all the handlers on
type Server struct {
    a      *pbExample.Request
    b      *pbExample.Request
    Result *pbExample.Response
}

// This is where you implement the handler for the Add method
func (s *Server) Add(ctx context.Context, request *pbExample.Request) (*pbExample.Response, error) {
    a, b := request.GetA(), request.GetB()

    result := a + b

    return &pbExample.Response{Result: result}, nil
}

// This is where you implement the handler for the Divide method
func (s *Server) Divide(ctx context.Context, request *pbExample.Request) (*pbExample.Response, error) {
    a, b := request.GetA(), request.GetB()

    result := a / b

    return &pbExample.Response{Result: result}, nil
}

// This is where you implement the handler for the Multiply method
func (s *Server) Multiply(ctx context.Context, request *pbExample.Request) (*pbExample.Response, error) {
    a, b := request.GetA(), request.GetB()

    result := a * b

    return &pbExample.Response{Result: result}, nil
}

// This is where you implement the handler for the Subtract method
func (s *Server) Subtract(ctx context.Context, request *pbExample.Request) (*pbExample.Response, error) {
    a, b := request.GetA(), request.GetB()

    result := a - b

    return &pbExample.Response{Result: result}, nil
}
```

For `main.go` file of Arithmetic server that will serve on a given port and listening at the endpoint, we defined in proto files we can refer to the boilerplate repository [grpc-gateway-boilerplate
Template](https://github.com/johanbrandhorst/grpc-gateway-boilerplate).
