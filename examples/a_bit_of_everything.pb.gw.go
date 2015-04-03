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

func handle_ABitOfEverythingService_Create(ctx context.Context, c *web.C, client ABitOfEverythingServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(ABitOfEverything)

	if err = json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}

	var val string
	var ok bool

	val, ok = c.URLParams["float_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "float_value")
	}
	protoReq.FloatValue, err = convert.Float32(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["double_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "double_value")
	}
	protoReq.DoubleValue, err = convert.Float64(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["int64_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "int64_value")
	}
	protoReq.Int64Value, err = convert.Int64(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["uint64_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "uint64_value")
	}
	protoReq.Uint64Value, err = convert.Uint64(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["int32_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "int32_value")
	}
	protoReq.Int32Value, err = convert.Int32(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["fixed64_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "fixed64_value")
	}
	protoReq.Fixed64Value, err = convert.Uint64(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["fixed32_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "fixed32_value")
	}
	protoReq.Fixed32Value, err = convert.Uint32(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["bool_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "bool_value")
	}
	protoReq.BoolValue, err = convert.Bool(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["string_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "string_value")
	}
	protoReq.StringValue, err = convert.String(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["uint32_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "uint32_value")
	}
	protoReq.Uint32Value, err = convert.Uint32(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["sfixed32_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "sfixed32_value")
	}
	protoReq.Sfixed32Value, err = convert.Int32(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["sfixed64_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "sfixed64_value")
	}
	protoReq.Sfixed64Value, err = convert.Int64(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["sint32_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "sint32_value")
	}
	protoReq.Sint32Value, err = convert.Int32(val)
	if err != nil {
		return nil, err
	}

	val, ok = c.URLParams["sint64_value"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "sint64_value")
	}
	protoReq.Sint64Value, err = convert.Int64(val)
	if err != nil {
		return nil, err
	}

	return client.Create(ctx, protoReq)
}

func handle_ABitOfEverythingService_CreateBody(ctx context.Context, c *web.C, client ABitOfEverythingServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(ABitOfEverything)

	if err = json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}

	return client.CreateBody(ctx, protoReq)
}

func handle_ABitOfEverythingService_Lookup(ctx context.Context, c *web.C, client ABitOfEverythingServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(IdMessage)

	var val string
	var ok bool

	val, ok = c.URLParams["uuid"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "uuid")
	}
	protoReq.Uuid, err = convert.String(val)
	if err != nil {
		return nil, err
	}

	return client.Lookup(ctx, protoReq)
}

func handle_ABitOfEverythingService_Update(ctx context.Context, c *web.C, client ABitOfEverythingServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(ABitOfEverything)

	if err = json.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, err
	}

	var val string
	var ok bool

	val, ok = c.URLParams["uuid"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "uuid")
	}
	protoReq.Uuid, err = convert.String(val)
	if err != nil {
		return nil, err
	}

	return client.Update(ctx, protoReq)
}

func handle_ABitOfEverythingService_Delete(ctx context.Context, c *web.C, client ABitOfEverythingServiceClient, req *http.Request) (msg proto.Message, err error) {
	protoReq := new(IdMessage)

	var val string
	var ok bool

	val, ok = c.URLParams["uuid"]
	if !ok {
		return nil, fmt.Errorf("missing parameter %s", "uuid")
	}
	protoReq.Uuid, err = convert.String(val)
	if err != nil {
		return nil, err
	}

	return client.Delete(ctx, protoReq)
}

func RegisterABitOfEverythingServiceHandlerFromEndpoint(ctx context.Context, mux *web.Mux, endpoint string) (err error) {
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

	return RegisterABitOfEverythingServiceHandler(ctx, mux, conn)
}

func RegisterABitOfEverythingServiceHandler(ctx context.Context, mux *web.Mux, conn *grpc.ClientConn) error {
	client := NewABitOfEverythingServiceClient(conn)

	mux.Post("/v1/example/a_bit_of_everything/:float_value/:double_value/:int64_value/separator/:uint64_value/:int32_value/:fixed64_value/:fixed32_value/:bool_value/:string_value/:uint32_value/:sfixed32_value/:sfixed64_value/:sint32_value/:sint64_value", func(c *web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_ABitOfEverythingService_Create(ctx, c, client, req)
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

	mux.Post("/v1/example/a_bit_of_everything", func(c *web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_ABitOfEverythingService_CreateBody(ctx, c, client, req)
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

	mux.Get("/v1/example/a_bit_of_everything/:uuid", func(c *web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_ABitOfEverythingService_Lookup(ctx, c, client, req)
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

	mux.Put("/v1/example/a_bit_of_everything/:uuid", func(c *web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_ABitOfEverythingService_Update(ctx, c, client, req)
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

	mux.Delete("/v1/example/a_bit_of_everything/:uuid", func(c *web.C, w http.ResponseWriter, req *http.Request) {
		resp, err := handle_ABitOfEverythingService_Delete(ctx, c, client, req)
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
}
