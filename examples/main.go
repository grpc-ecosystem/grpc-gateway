package main

import (
	"flag"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/examples/examplepb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:9090", "endpoint of EchoService")
	abeEndpoint  = flag.String("more_endpoint", "localhost:9090", "endpoint of ABitOfEverythingService")
	flowEndpoint = flag.String("flow_endpoint", "localhost:9090", "endpoint of FlowCombination")

	swaggerDir = flag.String("swagger_dir", "examples/examplepb", "path to the directory which contains swagger definitions")
)

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := examplepb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = examplepb.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, *abeEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = examplepb.RegisterABitOfEverythingServiceHandlerFromEndpoint(ctx, mux, *abeEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = examplepb.RegisterFlowCombinationHandlerFromEndpoint(ctx, mux, *flowEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		glog.Errorf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(*swaggerDir, p)
	http.ServeFile(w, r, p)
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}

// Run starts a HTTP server and blocks forever if successful.
func Run(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", serveSwagger)

	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	return http.ListenAndServe(address, allowCORS(mux))
}

func main() {
	flag.Parse()
	defer glog.Flush()
	if err := Run(":8080", runtime.WithProtoErrorHandler(HTTPError)); err != nil {
		glog.Fatal(err)
	}
}

func HTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	type errorBody struct {
		Error string `protobuf:"bytes,1,name=error" json:"error"`
		Code  int32  `protobuf:"varint,2,name=code" json:"code"`
	}

	body := &errorBody{
		Error: s.Message(),
		Code:  200,
	}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Printf("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Printf("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Printf("Failed to extract ServerMetadata from context")
	}

	runtime.HandleForwardResponseServerMetadata(w, mux, md)
	runtime.HandleForwardResponseTrailerHeader(w, md)
	st := runtime.HTTPStatusFromCode(0)
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Printf("Failed to write response: %v", err)
	}

	runtime.HandleForwardResponseTrailer(w, md)
}
