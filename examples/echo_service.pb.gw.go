package main

import (
	"encoding/json"

	"fmt"

	"net/http"

	"google.golang.org/grpc"

	"github.com/gengo/grpc-gateway/convert"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

func handle_EchoService_Echo(ctx context.Context, c web.C, client EchoServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(SimpleMessage)

	var val string
	var ok bool

	val, ok = c.URLParams["id"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "id")
	}
	protoReq.Id, err = convert.String(val)
	if err != nil {
		return nil, err
	}

	return client.Echo(ctx, protoReq)
}

func handle_EchoService_EchoBody(ctx context.Context, c web.C, client EchoServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(SimpleMessage)

	if err = json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}

	return client.EchoBody(ctx, protoReq)
}

func RegisterEchoServiceHandlerFromEndpoint(ctx context.Context, mux *web.Mux, endpoint string) (err error) {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				glog.Error("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				glog.Error("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterEchoServiceHandler(ctx, mux, conn)
}

func RegisterEchoServiceHandler(ctx context.Context, mux *web.Mux, conn *grpc.ClientConn) error {
	client := NewEchoServiceClient(conn)

	mux.Post("/v1/example/echo/:id", func(c web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_EchoService_Echo(ctx, c, client, req)
		if err != nil {
			glog.Errorf("RPC error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf, err := json.Marshal(resp)
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

	mux.Post("/v1/example/echo_body", func(c web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_EchoService_EchoBody(ctx, c, client, req)
		if err != nil {
			glog.Errorf("RPC error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf, err := json.Marshal(resp)
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
}
