---
layout: default
title: Introduction about gRPC-Gateway
parent: Tutorials
nav_order: 1
---

## Introduction

We all know that gRPC is not a tool for everything. There are cases where we still want to provide a traditional RESTful JSON API. The reasons can range from maintaining backwards-compatibility to supporting programming languages or clients not well supported by gRPC. But coding another API for REST is quite a time consuming and tedious.

So is there any way to code just once, but can provide APIs in both gRPC and REST at the same time?

The answer is Yes.

The grpc-gateway is a plugin of the Google protocol buffers compiler [protoc](https://github.com/protocolbuffers/protobuf). It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC. This server is generated according to the [`google.api.http`](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L46) annotations in your service definitions.

This helps you provide your APIs in both gRPC and RESTful style at the same time.

![architecture introduction diagram](https://docs.google.com/drawings/d/12hp4CPqrNPFhattL_cIoJptFvlAqm5wLQ0ggqI5mkCg/pub?w=749&h=370)

### Prerequisites

Before we start coding, we have to install some tools and need to do a small amount of configuration in the service.

We will be using a Go gRPC server in the examples, so please install Go first from [https://golang.org/dl/](https://golang.org/dl/).

After installing Go, use `go get` to download the following packages:

```sh
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

This installs the protoc generator plugins we need to generate the stubs. Make sure to add `$GOPATH/bin` to your `$PATH` so that executables installed via `go get` are available on your `$PATH`.
