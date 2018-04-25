---
category: documentation
---

# Examples

Examples are available under `examples` directory.
* `examplepb/echo_service.proto`, `examplepb/a_bit_of_everything.proto`, `examplepb/unannotated_echo_service.proto`: service definition
  * `examplepb/echo_service.pb.go`, `examplepb/a_bit_of_everything.pb.go`, `examplepb/unannotated_echo_service.pb.go`: [generated] stub of the service
  * `examplepb/echo_service.pb.gw.go`, `examplepb/a_bit_of_everything.pb.gw.go`, `examplepb/uannotated_echo_service.pb.gw.go`: [generated] reverse proxy for the service
  * `examplepb/unannotated_echo_service.yaml`: gRPC API Configuration for ```unannotated_echo_service.proto```
* `server/main.go`: service implementation
* `main.go`: entrypoint of the generated reverse proxy

To use the same port for custom HTTP handlers (e.g. serving `swagger.json`), gRPC-gateway, and a gRPC server, see [this code example by CoreOS](https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go) (and its accompanying [blog post](https://coreos.com/blog/gRPC-protobufs-swagger.html))


