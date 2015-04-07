package main

import (
	"flag"
	"fmt"
	"net"
	"sync"

	examples "github.com/gengo/grpc-gateway/examples"
	sub "github.com/gengo/grpc-gateway/examples/sub"
	"github.com/golang/glog"
	"github.com/rogpeppe/fastuuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	v map[string]*examples.ABitOfEverything
	m sync.Mutex
}

func newABitOfEverythingServer() examples.ABitOfEverythingServiceServer {
	return &_ABitOfEverythingServer{
		v: make(map[string]*examples.ABitOfEverything),
	}
}

func (s *_ABitOfEverythingServer) Create(ctx context.Context, msg *examples.ABitOfEverything) (*examples.ABitOfEverything, error) {
	s.m.Lock()
	defer s.m.Unlock()

	glog.Info(msg)
	var uuid string
	for {
		uuid = fmt.Sprintf("%x", uuidgen.Next())
		if _, ok := s.v[uuid]; !ok {
			break
		}
	}
	s.v[uuid] = msg
	s.v[uuid].Uuid = uuid
	return s.v[uuid], nil
}

func (s *_ABitOfEverythingServer) CreateBody(ctx context.Context, msg *examples.ABitOfEverything) (*examples.ABitOfEverything, error) {
	return s.Create(ctx, msg)
}

func (s *_ABitOfEverythingServer) Lookup(ctx context.Context, msg *examples.IdMessage) (*examples.ABitOfEverything, error) {
	s.m.Lock()
	defer s.m.Unlock()

	glog.Info(msg)
	if a, ok := s.v[msg.Uuid]; ok {
		return a, nil
	}
	return nil, grpc.Errorf(codes.NotFound, "not found")
}

func (s *_ABitOfEverythingServer) Update(ctx context.Context, msg *examples.ABitOfEverything) (*examples.EmptyMessage, error) {
	s.m.Lock()
	defer s.m.Unlock()

	glog.Info(msg)
	if _, ok := s.v[msg.Uuid]; ok {
		s.v[msg.Uuid] = msg
	} else {
		return nil, grpc.Errorf(codes.NotFound, "not found")
	}
	return new(examples.EmptyMessage), nil
}

func (s *_ABitOfEverythingServer) Delete(ctx context.Context, msg *examples.IdMessage) (*examples.EmptyMessage, error) {
	s.m.Lock()
	defer s.m.Unlock()

	glog.Info(msg)
	if _, ok := s.v[msg.Uuid]; ok {
		delete(s.v, msg.Uuid)
	} else {
		return nil, grpc.Errorf(codes.NotFound, "not found")
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
