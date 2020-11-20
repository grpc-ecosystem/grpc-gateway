---
layout: default
title: Prerequisites
parent: Tutorials
nav_order: 1
---

## Prerequisites

One way to achieve that is to use gRPC gateway. gRPC gateway is a plugin of the protocol buffer compiler. It reads the protobuf service definitions and generates a proxy server, which translates a RESTful HTTP call into gRPC request.

All we need to do is a small amount of configuration
in the service. Before start coding, we have to install some
tools.

We will be using a Go gRPC server in the examples, so please install Go first from [https://golang.org/dl/](https://golang.org/dl/).

After installing Go, use `go get` to download the following packages:

```sh
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```

This installs the protoc generator plugins we need to generate the stubs. Make sure to add `$GOPATH/bin` to your `$PATH` so that executables installed via `go get` are available on your `$PATH`.
