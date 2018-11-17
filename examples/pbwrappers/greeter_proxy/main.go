package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "../helloworld"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	grpcServiceHost := ":8090"
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, grpcServiceHost, opts)
	if err != nil {
		return err
	}

	proxyHostAddr := ":8091"
	log.Printf("Proxy from: %s to: %s", proxyHostAddr, grpcServiceHost)
	return http.ListenAndServe(proxyHostAddr, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
