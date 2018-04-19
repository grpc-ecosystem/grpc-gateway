package server

import (
	"net"

	"github.com/golang/glog"
	examples "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Run starts the example gRPC service.
func Run(ctx context.Context) error {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			glog.Errorf("Failed to close tcp :9090: %v", err)
		}
	}()

	s := grpc.NewServer()
	examples.RegisterEchoServiceServer(s, newEchoServer())
	examples.RegisterFlowCombinationServer(s, newFlowCombinationServer())

	abe := newABitOfEverythingServer()
	examples.RegisterABitOfEverythingServiceServer(s, abe)
	examples.RegisterStreamServiceServer(s, abe)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()
	return s.Serve(l)
}
