package main

import (
	"flag"
	"net/http"

	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:9090", "endpoint of EchoService")
	abeEndpoint  = flag.String("more_endpoint", "localhost:9090", "endpoint of ABitOfEverythingService")
)

func Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := RegisterEchoServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint)
	if err != nil {
		return err
	}
	err = RegisterABitOfEverythingServiceHandlerFromEndpoint(ctx, mux, *abeEndpoint)
	if err != nil {
		return err
	}

	http.ListenAndServe(":8080", mux)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := Run(); err != nil {
		glog.Fatal(err)
	}
}
