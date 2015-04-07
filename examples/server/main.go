package main

import (
	"errors"
	"flag"
	"fmt"
	"net"

	examples "github.com/gengo/grpc-gateway/examples"
	sub "github.com/gengo/grpc-gateway/examples/sub"
	"github.com/golang/glog"
	"github.com/rogpeppe/fastuuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Implements of EchoServiceServer

type echoServer struct{}

func newEchoServer() examples.EchoServiceServer {
	return new(echoServer)
}

func (s *echoServer) Echo(ctx context.Context, msg *examples.SimpleMessage) (*examples.SimpleMessage, error) {
	glog.Info(msg)
	return msg, nil
}

func (s *echoServer) EchoBody(ctx context.Context, msg *examples.SimpleMessage) (*examples.SimpleMessage, error) {
	glog.Info(msg)
	return msg, nil
}

// Implements of ABitOfEverythingServiceServer

var uuidgen = fastuuid.MustNewGenerator()

type _ABitOfEverythingServer struct {
	m map[string]*examples.ABitOfEverything
}

func newABitOfEverythingServer() examples.ABitOfEverythingServiceServer {
	return &_ABitOfEverythingServer{
		m: make(map[string]*examples.ABitOfEverything),
	}
}

func (s *_ABitOfEverythingServer) Create(ctx context.Context, msg *examples.ABitOfEverything) (*examples.ABitOfEverything, error) {
	glog.Info(msg)
	var uuid string
	for {
		uuid = fmt.Sprintf("%x", uuidgen.Next())
		if _, ok := s.m[uuid]; !ok {
			break
		}
	}
	s.m[uuid] = msg
	s.m[uuid].Uuid = uuid
	return s.m[uuid], nil
}

func (s *_ABitOfEverythingServer) CreateBody(ctx context.Context, msg *examples.ABitOfEverything) (*examples.ABitOfEverything, error) {
	return s.Create(ctx, msg)
}

func (s *_ABitOfEverythingServer) Lookup(ctx context.Context, msg *examples.IdMessage) (*examples.ABitOfEverything, error) {
	glog.Info(msg)
	if a, ok := s.m[msg.Uuid]; ok {
		return a, nil
	}
	return nil, errors.New("not found")
}

func (s *_ABitOfEverythingServer) Update(ctx context.Context, msg *examples.ABitOfEverything) (*examples.EmptyMessage, error) {
	glog.Info(msg)
	if _, ok := s.m[msg.Uuid]; ok {
		s.m[msg.Uuid] = msg
	} else {
		return nil, errors.New("not found")
	}
	return new(examples.EmptyMessage), nil
}

func (s *_ABitOfEverythingServer) Delete(ctx context.Context, msg *examples.IdMessage) (*examples.EmptyMessage, error) {
	glog.Info(msg)
	if _, ok := s.m[msg.Uuid]; ok {
		delete(s.m, msg.Uuid)
	} else {
		return nil, errors.New("not found")
	}
	return new(examples.EmptyMessage), nil
}

func (s *_ABitOfEverythingServer) Echo(ctx context.Context, msg *sub.StringMessage) (*sub.StringMessage, error) {
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
	examples.RegisterEchoServiceServer(s, newEchoServer())
	examples.RegisterABitOfEverythingServiceServer(s, newABitOfEverythingServer())
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
