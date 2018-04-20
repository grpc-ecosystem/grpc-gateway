package main


import (
	"flag"
	"net/http"
	"strconv"
	"strings"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "github.com/tronprotocol/grpc-gateway/api"

)

var (
	port = flag.Int("port",50051, "port of your tron grpc service" )
	host = flag.String("host", "localhost", "host of your tron grpc service")
)

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

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	echoEndpoint := *host + ":"  + strconv.Itoa(*port)
	opts := []grpc.DialOption{grpc.WithInsecure()}


	err := gw.RegisterWalletHandlerFromEndpoint(ctx, mux, echoEndpoint, opts)
	if err != nil {
		return err
	}

	//fmt.Printf("connecting %s", echoEndpoint)

	return http.ListenAndServe(":8086", allowCORS(mux))
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
