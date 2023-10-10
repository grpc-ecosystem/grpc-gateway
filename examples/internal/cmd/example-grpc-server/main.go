/*
Command example-grpc-server is an example grpc server
to be called by example-gateway-server.
*/
package main

import (
	"context"
	"flag"

	"github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/server"
	"google.golang.org/grpc/grpclog"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC service")
	network = flag.String("network", "tcp", "a valid network type which is consistent to -addr")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	if err := server.Run(ctx, *network, *addr); err != nil {
		grpclog.Fatal(err)
	}
}
