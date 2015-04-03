package main

import (
	"flag"
	"net"

	examples "github.com/gengo/grpc-gateway/examples"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type echoServer struct{}

func (s echoServer) Echo(ctx context.Context, msg *examples.SimpleMessage) (*examples.SimpleMessage, error) {
	glog.Info(msg)
	return msg, nil
}

func (s echoServer) EchoBody(ctx context.Context, msg *examples.SimpleMessage) (*examples.SimpleMessage, error) {
	glog.Info(msg)
	return msg, nil
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	examples.RegisterEchoServiceServer(s, echoServer{})
	s.Serve(l)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
