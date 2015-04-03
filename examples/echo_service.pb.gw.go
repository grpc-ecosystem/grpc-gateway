package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func handle_EchoService_Echo(ctx context.Context, c EchoServiceClient, req *http.Request) (proto.Message, error) {
	protoReq := new(SimpleMessage)

	if err := json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}

	return c.Echo(ctx, protoReq)
}

var (
	endpointEchoService string
)

func init() {

	flag.StringVar(&endpointEchoService, "echo_service_endpoint", "", "endpoint host:port of EchoService")

}

type handler struct {
	mux   http.ServeMux
	conns map[string]*grpc.ClientConn
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) Close() error {
	var err error
	for svc, conn := range h.conns {
		cerr := conn.Close()
		if err == nil {
			err = cerr
		}
		if cerr != nil {
			glog.Errorf("Failed to close gRPC connection to %s: %v", svc, err)
		}
	}
	return err
}

func NewHandler(ctx context.Context) (http.Handler, error) {
	h := &handler{
		conns: make(map[string]*grpc.ClientConn),
	}
	var err error
	defer func() {
		if err != nil {
			h.Close()
		}
	}()

	err = func() error {
		conn, err := grpc.Dial(endpointEchoService)
		if err != nil {
			return err
		}
		h.conns["EchoService"] = conn
		client := NewEchoServiceClient(conn)

		h.mux.HandleFunc("/v1/example/echo", func(w http.ResponseWriter, req *http.Request) {
			resp, err := handle_EchoService_Echo(ctx, client, req)
			if err != nil {
				glog.Errorf("RPC error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf, err := proto.Marshal(resp)
			if err != nil {
				glog.Errorf("Marshal error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err = w.Write(buf); err != nil {
				glog.Errorf("Failed to write response: %v", err)
			}
		})

		return nil
	}()
	if err != nil {
		return nil, err
	}

	return h, nil
}
