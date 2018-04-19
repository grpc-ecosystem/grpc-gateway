package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/examples/server"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		glog.Fatal(err)
	}
}
