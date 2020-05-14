// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: examples/internal/proto/examplepb/response_body_service.proto

/*
Package examplepb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package examplepb

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/utilities"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

func request_ResponseBodyService_GetResponseBody_0(ctx context.Context, marshaler runtime.Marshaler, client ResponseBodyServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := client.GetResponseBody(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ResponseBodyService_GetResponseBody_0(ctx context.Context, marshaler runtime.Marshaler, server ResponseBodyServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := server.GetResponseBody(ctx, &protoReq)
	return msg, metadata, err

}

func request_ResponseBodyService_ListResponseBodies_0(ctx context.Context, marshaler runtime.Marshaler, client ResponseBodyServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := client.ListResponseBodies(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ResponseBodyService_ListResponseBodies_0(ctx context.Context, marshaler runtime.Marshaler, server ResponseBodyServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := server.ListResponseBodies(ctx, &protoReq)
	return msg, metadata, err

}

func request_ResponseBodyService_ListResponseStrings_0(ctx context.Context, marshaler runtime.Marshaler, client ResponseBodyServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := client.ListResponseStrings(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ResponseBodyService_ListResponseStrings_0(ctx context.Context, marshaler runtime.Marshaler, server ResponseBodyServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	msg, err := server.ListResponseStrings(ctx, &protoReq)
	return msg, metadata, err

}

func request_ResponseBodyService_GetResponseBodyStream_0(ctx context.Context, marshaler runtime.Marshaler, client ResponseBodyServiceClient, req *http.Request, pathParams map[string]string) (ResponseBodyService_GetResponseBodyStreamClient, runtime.ServerMetadata, error) {
	var protoReq ResponseBodyIn
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["data"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "data")
	}

	protoReq.Data, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "data", err)
	}

	stream, err := client.GetResponseBodyStream(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil

}

// RegisterResponseBodyServiceHandlerServer registers the http handlers for service ResponseBodyService to "mux".
// UnaryRPC     :call ResponseBodyServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
func RegisterResponseBodyServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ResponseBodyServiceServer) error {

	mux.Handle("GET", pattern_ResponseBodyService_GetResponseBody_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/GetResponseBody")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ResponseBodyService_GetResponseBody_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_GetResponseBody_0(ctx, mux, outboundMarshaler, w, req, response_ResponseBodyService_GetResponseBody_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ResponseBodyService_ListResponseBodies_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/ListResponseBodies")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ResponseBodyService_ListResponseBodies_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_ListResponseBodies_0(ctx, mux, outboundMarshaler, w, req, response_ResponseBodyService_ListResponseBodies_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ResponseBodyService_ListResponseStrings_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/ListResponseStrings")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ResponseBodyService_ListResponseStrings_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_ListResponseStrings_0(ctx, mux, outboundMarshaler, w, req, response_ResponseBodyService_ListResponseStrings_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ResponseBodyService_GetResponseBodyStream_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	})

	return nil
}

// RegisterResponseBodyServiceHandlerFromEndpoint is same as RegisterResponseBodyServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterResponseBodyServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterResponseBodyServiceHandler(ctx, mux, conn)
}

// RegisterResponseBodyServiceHandler registers the http handlers for service ResponseBodyService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterResponseBodyServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterResponseBodyServiceHandlerClient(ctx, mux, NewResponseBodyServiceClient(conn))
}

// RegisterResponseBodyServiceHandlerClient registers the http handlers for service ResponseBodyService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ResponseBodyServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ResponseBodyServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ResponseBodyServiceClient" to call the correct interceptors.
func RegisterResponseBodyServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ResponseBodyServiceClient) error {

	mux.Handle("GET", pattern_ResponseBodyService_GetResponseBody_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/GetResponseBody")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ResponseBodyService_GetResponseBody_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_GetResponseBody_0(ctx, mux, outboundMarshaler, w, req, response_ResponseBodyService_GetResponseBody_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ResponseBodyService_ListResponseBodies_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/ListResponseBodies")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ResponseBodyService_ListResponseBodies_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_ListResponseBodies_0(ctx, mux, outboundMarshaler, w, req, response_ResponseBodyService_ListResponseBodies_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ResponseBodyService_ListResponseStrings_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/ListResponseStrings")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ResponseBodyService_ListResponseStrings_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_ListResponseStrings_0(ctx, mux, outboundMarshaler, w, req, response_ResponseBodyService_ListResponseStrings_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ResponseBodyService_GetResponseBodyStream_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/grpc.gateway.examples.internal.proto.examplepb.ResponseBodyService/GetResponseBodyStream")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ResponseBodyService_GetResponseBodyStream_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ResponseBodyService_GetResponseBodyStream_0(ctx, mux, outboundMarshaler, w, req, func() (proto.Message, error) {
			res, err := resp.Recv()
			return response_ResponseBodyService_GetResponseBodyStream_0{res}, err
		}, mux.GetForwardResponseOptions()...)

	})

	return nil
}

type response_ResponseBodyService_GetResponseBody_0 struct {
	proto.Message
}

func (m response_ResponseBodyService_GetResponseBody_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*ResponseBodyOut)
	return response.Response
}

type response_ResponseBodyService_ListResponseBodies_0 struct {
	proto.Message
}

func (m response_ResponseBodyService_ListResponseBodies_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*RepeatedResponseBodyOut)
	return response.Response
}

type response_ResponseBodyService_ListResponseStrings_0 struct {
	proto.Message
}

func (m response_ResponseBodyService_ListResponseStrings_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*RepeatedResponseStrings)
	return response.Values
}

type response_ResponseBodyService_GetResponseBodyStream_0 struct {
	proto.Message
}

func (m response_ResponseBodyService_GetResponseBodyStream_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*ResponseBodyOut)
	return response.Response
}

var (
	pattern_ResponseBodyService_GetResponseBody_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"responsebody", "data"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ResponseBodyService_ListResponseBodies_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"responsebodies", "data"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ResponseBodyService_ListResponseStrings_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"responsestrings", "data"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ResponseBodyService_GetResponseBodyStream_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2}, []string{"responsebody", "stream", "data"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_ResponseBodyService_GetResponseBody_0 = runtime.ForwardResponseMessage

	forward_ResponseBodyService_ListResponseBodies_0 = runtime.ForwardResponseMessage

	forward_ResponseBodyService_ListResponseStrings_0 = runtime.ForwardResponseMessage

	forward_ResponseBodyService_GetResponseBodyStream_0 = runtime.ForwardResponseStream
)
