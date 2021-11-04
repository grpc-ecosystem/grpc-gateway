FROM golang:1.17.3 as builder

ARG RELEASE_VERSION

# Buf plugins must be built for linux/amd64
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${RELEASE_VERSION}

FROM scratch

COPY --from=builder /go/bin/protoc-gen-openapiv2 /usr/local/bin/protoc-gen-openapiv2

ENTRYPOINT ["/usr/local/bin/protoc-gen-openapiv2"]
