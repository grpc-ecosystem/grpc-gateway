package gateway

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/examples/examplepb"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type optSet struct {
	mux  []gwruntime.ServeMuxOption
	dial []grpc.DialOption

	echoEndpoint, abeEndpoint, flowEndpoint string
}

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(ctx context.Context, opts optSet) (http.Handler, error) {
	mux := gwruntime.NewServeMux(opts.mux...)

	err := examplepb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, opts.echoEndpoint, opts.dial)
	if err != nil {
		return nil, err
	}
	err = examplepb.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, opts.abeEndpoint, opts.dial)
	if err != nil {
		return nil, err
	}
	err = examplepb.RegisterABitOfEverythingServiceHandlerFromEndpoint(ctx, mux, opts.abeEndpoint, opts.dial)
	if err != nil {
		return nil, err
	}
	err = examplepb.RegisterFlowCombinationHandlerFromEndpoint(ctx, mux, opts.flowEndpoint, opts.dial)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

// NewTCPGateway returns a new gateway server which connect to the gRPC service with TCP.
// "addr" must be a valid TCP address with a port number.
func NewTCPGateway(ctx context.Context, addr string, opts ...gwruntime.ServeMuxOption) (http.Handler, error) {
	return newGateway(ctx, optSet{
		mux:          opts,
		dial:         []grpc.DialOption{grpc.WithInsecure()},
		echoEndpoint: addr,
		abeEndpoint:  addr,
		flowEndpoint: addr,
	})
}

// NewUnixGatway returns a new gateway server which connect to the gRPC service with a unix domain socket.
// "addr" must be a valid path to the socket.
func NewUnixGateway(ctx context.Context, addr string, opts ...gwruntime.ServeMuxOption) (http.Handler, error) {
	return newGateway(ctx, optSet{
		mux: opts,
		dial: []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout("unix", addr, timeout)
			}),
		},
		echoEndpoint: addr,
		abeEndpoint:  addr,
		flowEndpoint: addr,
	})
}
