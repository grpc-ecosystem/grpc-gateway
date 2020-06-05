---
category: documentation
---

# How do I use this?

The grpc gateway v2 allows you inject custom router in grpc gateway.

It may become a good alternative to `gorilla/mux` or other golang router matcher.

## Example

```go
package main

import (
	"context"
	"net/http"
	
	pb "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	ctx := context.TODO()
	mux := runtime.NewServeMux()
	// Register proto's GET /say/{name} url
	_ = pb.RegisterGreeterHandlerServer(ctx, mux, &GreeterServer{})
	// Register custom router urls GET /hello/{name} url
	_ = mux.HandlePath("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("hello " + pathParams["name"]))
	})
	panic(http.ListenAndServe(":8080", mux))
}


// GreeterServer is the server API for Greeter service.
type GreeterServer struct {
	
}

// SayHello implement to say hello
func (h *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
```

> Note: You may need change the `pb "github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/helloworld"` to your code