package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

var (
	echoEndpoint = flag.String("--echo_endpoint", "localhost:9090", "endpoint of EchoService")
	abeEndpoint  = flag.String("--more_endpoint", "localhost:9090", "endpoint of ABitOfEverythingService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := web.New()
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

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
