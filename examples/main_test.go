package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/glog"
	server "github.com/grpc-ecosystem/grpc-gateway/examples/server"
	"golang.org/x/net/context"
)

func runServers(ctx context.Context) <-chan error {
	ch := make(chan error, 2)
	go func() {
		if err := server.Run(ctx, *network, *endpoint); err != nil {
			ch <- fmt.Errorf("cannot run grpc service: %v", err)
		}
	}()
	go func() {
		if err := Run(ctx, ":8080"); err != nil {
			ch <- fmt.Errorf("cannot run gateway service: %v", err)
		}
	}()
	return ch
}

func TestMain(m *testing.M) {
	flag.Parse()
	defer glog.Flush()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := runServers(ctx)

	ch := make(chan int, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- m.Run()
	}()

	select {
	case err := <-errCh:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	case status := <-ch:
		cancel()
		os.Exit(status)
	}
}
