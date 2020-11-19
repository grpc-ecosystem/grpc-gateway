---
layout: default
title: Creating a Basic Arithmetic gRPC Server
parent: Tutorials
nav_order: 3
---

## Creating a Basic Arithmetic gRPC Server

First look at the `main.go` file of Arithmetic server that will serve on a given port and listening at the endpoint we defined in proto files.

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

Now rest fo `main.go` file are:

```go
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
	pbExample.RegisterArithmeticServiceServer(s, &Server{})

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///" + addr)
	log.Fatalln(err)
}
```
