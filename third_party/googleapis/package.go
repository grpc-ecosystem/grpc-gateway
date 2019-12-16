package googleapis
// This is an empty package that allows vendoring the folder with `go mod`.
// For example, a `tools.go` file may look like the followings:
//
// package your_package
//
// import (
//     _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"
// )
//
// And you vendor this package by running `go mod vendor`. To generate grpc gateway code,
// you can run with the followings:
//
// protoc -I. \
//   -I./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
//   --plugin=protoc-gen-grpc=grpc_ruby_plugin \
//   --grpc-gateway_out=logtostderr=true:. \
//   path/to/your_service.proto
