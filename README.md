# Requirements

#### Go 1.6 or higher

You can download the binary release from [releases](https://golang.org/dl/) or compile the source release.

Make sure that your `$GOPATH/bin` is in your `$PATH`.



#### (Optional)ProtocolBuffers 3.0.0-beta-3 or later

In project, we have provided  gRPC stub and reverse-proxy so that this step is optional.

If you want to enerate gRPC stub and reverse-proxy by yourself, you need to install ProtocolBuffers 3.0.0-beta-3 or later. 

You can download the binary release from [releases](https://github.com/google/protobuf/releases) or compile the source release. 

```sh
cd /tmp
wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
mkdir protobuf
tar -xzvf protoc-3.5.1-linux-x86_64.zip -C ./protobuf
exprot PATH=$PATH:./protobuf/bin
```

> note: you can chose to  download other os realeas on https://github.com/google/protobuf/releases instead of linux realease



# Usage

1. Make sure your Tron grpc serivice has been started on `localhost:50051`  , **you can visit [Tron wiki: quick start](http://wiki.tron.network/en/latest/quick_start.html) for starting Tron service.**
2. Get the source code and change word dir

```
# download project
go get -u github.com/tronprotocol/grpc-gateway

# change to project dir
cd $GOPATH/src/github.com/tronprotocol/grpc-gateway
```

3. (Optional) Generate gRPC stub and reverse-proxy. Make sure you have install protoc

```
./gen_proto.sh
```
4. run proxy-server
```
go run tron_http/main.go
or
go run tron_http/main.go -port 50051 -host localhost
or
go run tron_http/main.go -port 50051 -host 10.0.8.214
```
5. Test API of tron http

```
curl -X POST -k http://localhost:8086/wallet/listaccount
```

If you get account-list jsonn data, congratulations





# grpc-gateway

[![Build Status](https://travis-ci.org/grpc-ecosystem/grpc-gateway.svg?branch=master)](https://travis-ci.org/grpc-ecosystem/grpc-gateway)

grpc-gateway is a plugin of [protoc](http://github.com/google/protobuf).
It reads [gRPC](http://github.com/grpc/grpc-common) service definition,
and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.
This server is generated according to [custom options](https://cloud.google.com/service-management/reference/rpc/google.api#http) in your gRPC definition.

It helps you to provide your APIs in both gRPC and RESTful style at the same time.

![architecture introduction diagram](https://docs.google.com/drawings/d/12hp4CPqrNPFhattL_cIoJptFvlAqm5wLQ0ggqI5mkCg/pub?w=749&h=370)

## Background

gRPC is great -- it generates API clients and server stubs in many programming languages, it is fast, easy-to-use, bandwidth-efficient and its design is combat-proven by Google.
However, you might still want to provide a traditional RESTful API as well. Reasons can range from maintaining backwards-compatibility, supporting languages or clients not well supported by gRPC to simply maintaining the aesthetics and tooling involved with a RESTful architecture.

This project aims to provide that HTTP+JSON interface to your gRPC service. A small amount of configuration in your service to attach HTTP semantics is all that's needed to generate a reverse-proxy with this library.

## Parameters and flags
`protoc-gen-grpc-gateway` supports custom mapping from Protobuf `import` to Golang import path.
They are compatible to [the parameters with same names in `protoc-gen-go`](https://github.com/golang/protobuf#parameters).

In addition we also support the `request_context` parameter in order to use the `http.Request`'s Context (only for Go 1.7 and above).
This parameter can be useful to pass request scoped context between the gateway and the gRPC service.

`protoc-gen-grpc-gateway` also supports some more command line flags to control logging. You can give these flags together with parameters above. Run `protoc-gen-grpc-gateway --help` for more details about the flags.



## Features
### Supported
* Generating JSON API handlers
* Method parameters in request body
* Method parameters in request path
* Method parameters in query string
* Enum fields in path parameter (including repeated enum fields).
* Mapping streaming APIs to newline-delimited JSON streams
* Mapping HTTP headers with `Grpc-Metadata-` prefix to gRPC metadata (prefixed with `grpcgateway-`)
* Optionally emitting API definition for [Swagger](http://swagger.io).
* Setting [gRPC timeouts](http://www.grpc.io/docs/guides/wire.html) through inbound HTTP `Grpc-Timeout` header.



## Mapping gRPC to HTTP

* [How gRPC error codes map to HTTP status codes in the response](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15)
* HTTP request source IP is added as `X-Forwarded-For` gRPC request header
* HTTP request host is added as `X-Forwarded-Host` gRPC request header
* HTTP `Authorization` header is added as `authorization` gRPC request header 
* Remaining Permanent HTTP header keys (as specified by the IANA [here](http://www.iana.org/assignments/message-headers/message-headers.xhtml) are prefixed with `grpcgateway-` and added with their values to gRPC request header
* HTTP headers that start with 'Grpc-Metadata-' are mapped to gRPC metadata (prefixed with `grpcgateway-`)
* While configurable, the default {un,}marshaling uses [jsonpb](https://godoc.org/github.com/golang/protobuf/jsonpb) with `OrigName: true`.

## Contribution
See [CONTRIBUTING.md](http://github.com/grpc-ecosystem/grpc-gateway/blob/master/CONTRIBUTING.md).

## License
grpc-gateway is licensed under the BSD 3-Clause License.
See [LICENSE.txt](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt) for more details.
