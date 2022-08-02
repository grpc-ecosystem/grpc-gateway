FROM golang:1.19.0 as builder

ARG RELEASE_VERSION

# Buf plugins must be built for linux/amd64
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${RELEASE_VERSION}

FROM scratch

ARG RELEASE_VERSION
ARG GO_PROTOBUF_RELEASE_VERSION
ARG GO_GRPC_RELEASE_VERSION

# Runtime dependencies
LABEL "build.buf.plugins.runtime_library_versions.0.name"="github.com/grpc-ecosystem/grpc-gateway/v2"
LABEL "build.buf.plugins.runtime_library_versions.0.version"="${RELEASE_VERSION}"
LABEL "build.buf.plugins.runtime_library_versions.1.name"="google.golang.org/protobuf"
LABEL "build.buf.plugins.runtime_library_versions.1.version"="${GO_PROTOBUF_RELEASE_VERSION}"
LABEL "build.buf.plugins.runtime_library_versions.2.name"="google.golang.org/grpc"
LABEL "build.buf.plugins.runtime_library_versions.2.version"="${GO_GRPC_RELEASE_VERSION}"

COPY --from=builder /go/bin/protoc-gen-grpc-gateway /usr/local/bin/protoc-gen-grpc-gateway

ENTRYPOINT ["/usr/local/bin/protoc-gen-grpc-gateway"]
