package main


import (
	"flag"
	"net/http"
	"strconv"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "github.com/tronprotocol/grpc-gateway/api"

	"fmt"
)

var (
	port = flag.Int("port",50051, "port of your tron grpc service" )
	host = flag.String("host", "localhost", "host of your tron grpc service")
	listen = flag.Int("listen", 8086, "the port that http server listen")

)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	grpcEndpoint := *host + ":"  + strconv.Itoa(*port)
	opts := []grpc.DialOption{grpc.WithInsecure()}


	fmt.Printf("grpc server: %s\n", grpcEndpoint)
	fmt.Printf("http port: %d\n", *listen)

	err := gw.RegisterWalletHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}



	return http.ListenAndServe(":" + strconv.Itoa(*listen), mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
