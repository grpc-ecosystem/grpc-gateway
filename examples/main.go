package main

import (
	"flag"
	"net/http"

	"github.com/gengo/grpc-gateway/examples/examplepb"
	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:9090", "endpoint of EchoService")
	abeEndpoint  = flag.String("more_endpoint", "localhost:9090", "endpoint of ABitOfEverythingService")
	flowEndpoint = flag.String("flow_endpoint", "localhost:9090", "endpoint of FlowCombination")
)

func Run(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := examplepb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, dialOpts)
	if err != nil {
		return err
	}
	err = examplepb.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, *abeEndpoint, dialOpts)
	if err != nil {
		return err
	}
	err = examplepb.RegisterABitOfEverythingServiceHandlerFromEndpoint(ctx, mux, *abeEndpoint, dialOpts)
	if err != nil {
		return err
	}
	err = examplepb.RegisterFlowCombinationHandlerFromEndpoint(ctx, mux, *flowEndpoint, dialOpts)
	if err != nil {
		return err
	}

	http.ListenAndServe(address, mux)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := Run(":8080"); err != nil {
		glog.Fatal(err)
	}
}
