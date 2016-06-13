package main

import (
	"flag"
	"net"

	examples "github.com/gengo/grpc-gateway/examples/examplepb"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

func Run() error {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	examples.RegisterEchoServiceServer(s, newEchoServer())
	examples.RegisterFlowCombinationServer(s, newFlowCombinationServer())

	abe := newABitOfEverythingServer()
	examples.RegisterABitOfEverythingServiceServer(s, abe)
	examples.RegisterStreamServiceServer(s, abe)

	s.Serve(l)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := Run(); err != nil {
		glog.Fatal(err)
	}
}
