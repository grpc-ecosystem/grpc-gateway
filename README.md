# grpc-gateway

grpc-gateway is a plugin of [protoc](http://github.com/google/protobuf).
It reads service [gRPC](http://github.com/grpc/grpc) service definition,
and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.

It helps you to provide your APIs in both gRPC and RESTful style at the same time.

## Installation
First you need to install ProtocolBuffers 3.0 or later.

```sh
mkdir tmp
cd tmp
git clone https://github.com/google/protobuf
./autogen.sh
./configure
make
make check
sudo make install
```

Then, `go get` as usual.

```sh
go get github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/golang/protobuf/protoc-gen-go
go install github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway
```
 
## Usage
Make sure that your `$GOPATH/bin` is in your `$PATH`.

### Generate gRPC stub

```sh
protoc -I. -I$GOPATH --go_out=plugins=grpc:. path/to/your_service.proto
```

It will generate a stub file `path/to/your_service.pb.go`.
Now you can implement your service on top of the stub.

### Generate reverse proxy

```sh
protoc --grpc-gateway_out=logtostderr=true:. path/to/your_service.proto
```

It will generate a reverse proxy `path/to/your_service.pb.gw.go`.
Now you need to write an entrypoint of the proxy server.


## Example
* `examples/echo_service.proto`: service definition
  * `examples/echo_service.pb.go`: [generated] stub of the service
  * `examples/echo_service.pb.gw.go`: [generated] reverse proxy for the service
* `examples/server/main.go`: service implementation
* `examples/main.go`: entrypoint of the generated reverse proxy
