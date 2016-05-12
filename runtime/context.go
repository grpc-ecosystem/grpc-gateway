package runtime

import (
	"net/http"
	"strings"

	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const metadataHeaderPrefix = "Grpc-Metadata-"
const metadataTrailerPrefix = "Grpc-Trailer-"
const xForwardedFor = "X-Forwarded-For"
const xForwardedHost = "X-Forwarded-Host"

/*
AnnotateContext adds context information such as metadata from the request.

If there are no metadata headers in the request, then the context returned
will be the same context.
*/
func AnnotateContext(ctx context.Context, req *http.Request) context.Context {
	var pairs []string
	for key, vals := range req.Header {
		for _, val := range vals {
			if key == "Authorization" {
				pairs = append(pairs, "authorization", val)
				continue
			}
			if strings.HasPrefix(key, metadataHeaderPrefix) {
				pairs = append(pairs, key[len(metadataHeaderPrefix):], val)
			}
		}
	}
	if host := req.Header.Get(xForwardedHost); host != "" {
		pairs = append(pairs, strings.ToLower(xForwardedHost), host)
	} else if req.Host != "" {
		pairs = append(pairs, strings.ToLower(xForwardedHost), req.Host)
	}
	remoteIp, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		if req.Header.Get(xForwardedFor) == "" {
			pairs = append(pairs, strings.ToLower(xForwardedFor), remoteIp)
		} else {
			pairs = append(pairs, strings.ToLower(xForwardedFor), req.Header.Get(xForwardedFor)+", "+remoteIp)
		}
	}

	if len(pairs) == 0 {
		return ctx
	}
	return metadata.NewContext(ctx, metadata.Pairs(pairs...))
}

// ServerMetadata consists of metadata sent from gRPC server.
type ServerMetadata struct {
	HeaderMD  metadata.MD
	TrailerMD metadata.MD
}

type serverMetadataKey struct{}

// NewServerMetadataContext creates a new context with ServerMetadata
func NewServerMetadataContext(ctx context.Context, md ServerMetadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// ServerMetadataFromContext returns the ServerMetadata in ctx
func ServerMetadataFromContext(ctx context.Context) (md ServerMetadata, ok bool) {
	md, ok = ctx.Value(serverMetadataKey{}).(ServerMetadata)
	return
}
