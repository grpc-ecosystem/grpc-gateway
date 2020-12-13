<div align="center">
<h1>gRPC-Gateway</h1>
<p>
gRPC to JSON proxy generator following the gRPC HTTP spec
</p>
<a href="https://circleci.com/gh/grpc-ecosystem/grpc-gateway"><img src="https://img.shields.io/circleci/build/github/grpc-ecosystem/grpc-gateway?color=379c9c&logo=circleci&logoColor=ffffff&style=flat-square"/></a>
<a href="https://codecov.io/gh/grpc-ecosystem/grpc-gateway"><img src="https://img.shields.io/codecov/c/github/grpc-ecosystem/grpc-gateway?color=379c9c&logo=codecov&logoColor=ffffff&style=flat-square"/></a>
<a href="https://app.slack.com/client/T029RQSE6/CBATURP1D"><img src="https://img.shields.io/badge/slack-grpc--gateway-379c9c?logo=slack&logoColor=ffffff&style=flat-square"/></a>
<a href="https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt"><img src="https://img.shields.io/github/license/grpc-ecosystem/grpc-gateway?color=379c9c&style=flat-square"/></a>
<a href="https://github.com/grpc-ecosystem/grpc-gateway/releases"><img src="https://img.shields.io/github/v/release/grpc-ecosystem/grpc-gateway?color=379c9c&logoColor=ffffff&style=flat-square"/></a>
<a href="https://github.com/grpc-ecosystem/grpc-gateway/stargazers"><img src="https://img.shields.io/github/stars/grpc-ecosystem/grpc-gateway?color=379c9c&style=flat-square"/></a>
</div>

## About

The gRPC-Gateway is a plugin of the Google protocol buffers compiler
[protoc](https://github.com/protocolbuffers/protobuf).
It reads protobuf service definitions and generates a reverse-proxy server which
translates a RESTful HTTP API into gRPC. This server is generated according to the
[`google.api.http`](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L46)
annotations in your service definitions.

This helps you provide your APIs in both gRPC and RESTful style at the same time.

<div align="center">
<img src="docs/assets/images/architecture_introduction_diagram.svg" />
</div>

## Testimonials

> We use the gRPC-Gateway to serve millions of API requests per day,
> and have been since 2018 and through all of that,
> we have never had any issues with it.
>
> _- William Mill, [Ad Hoc](http://adhocteam.us/)_

## Background

gRPC is great -- it generates API clients and server stubs in many programming
languages, it is fast, easy-to-use, bandwidth-efficient and its design is
combat-proven by Google. However, you might still want to provide a traditional
RESTful JSON API as well. Reasons can range from maintaining
backward-compatibility, supporting languages or clients that are not well supported by
gRPC, to simply maintaining the aesthetics and tooling involved with a RESTful
JSON architecture.

This project aims to provide that HTTP+JSON interface to your gRPC service.
A small amount of configuration in your service to attach HTTP semantics is all
that's needed to generate a reverse-proxy with this library.

## Installation

The gRPC-Gateway requires a local installation of the Google protocol buffers
compiler `protoc` v3.0.0 or above. Please install this via your local package
manager or by downloading one of the releases from the official repository:

https://github.com/protocolbuffers/protobuf/releases

The following instructions assume you are using
[Go Modules](https://github.com/golang/go/wiki/Modules) for dependency
management. Use a
[tool dependency](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
to track the versions of the following executable packages:

```go
// +build tools

package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

Run `go mod tidy` to resolve the versions. Install by running

```sh
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

This will place four binaries in your `$GOBIN`;

- `protoc-gen-grpc-gateway`
- `protoc-gen-openapiv2`
- `protoc-gen-go`
- `protoc-gen-go-grpc`

Make sure that your `$GOBIN` is in your `$PATH`.

## Usage

1. Define your [gRPC](https://grpc.io/docs/) service using protocol buffers

   `your_service.proto`:

   ```protobuf
    syntax = "proto3";
    package your.service.v1;
    option go_package = "github.com/yourorg/yourprotos/gen/go/your/service/v1";
    message StringMessage {
      string value = 1;
    }

    service YourService {
      rpc Echo(StringMessage) returns (StringMessage) {}
    }
   ```

2. Generate gRPC stubs

   This step generates the gRPC stubs that you can use to implement the service and consume from clients:

   Here's an example of what a `protoc` command might look like to generate Go stubs:

   ```sh
   protoc -I . \
      --go_out ./gen/go/ --go_opt paths=source_relative \
      --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative \
      your/service/v1/your_service.proto
   ```

3. Implement your service in gRPC as usual

   1. (Optional) Generate gRPC stub in the [other programming languages](https://grpc.io/docs/).

   For example, the following generates gRPC code for Ruby based on `your/service/v1/your_service.proto`:

   ```sh
   protoc -I . --ruby_out ./gen/ruby your/service/v1/your_service.proto

   protoc -I . --grpc-ruby_out ./gen/ruby your/service/v1/your_service.proto
   ```

   2. Add the googleapis-common-protos gem (or your language equivalent) as a dependency to your project.
   3. Implement your gRPC service stubs

4. Generate reverse-proxy using `protoc-gen-grpc-gateway`

   At this point, you have 3 options:

   - no further modifications, use the default mapping to HTTP semantics (method, path, etc.)
     - this will work on any `.proto` file, but will not allow setting HTTP paths, request parameters or similar
   - additional `.proto` modifications to use a custom mapping
     - relies on parameters in the `.proto` file to set custom HTTP mappings
   - no `.proto` modifications, but use an external configuration file
     - relies on an external configuration file to set custom HTTP mappings
     - mostly useful when the source proto file isn't under your control

   1. Using the default mapping

   This requires no additional modification to the `.proto` file but does require enabling a specific option when executing the plugin.
   The `generate_unbound_methods` should be enabled.

   Here's what a `protoc` execution might look like with this option enabled:

   ```sh
      protoc -I . --grpc-gateway_out ./gen/go \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        your/service/v1/your_service.proto
   ```

   2. With custom annotations

   Add a [`google.api.http`](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L46)
   annotation to your .proto file

   `your_service.proto`:

   ```diff
    syntax = "proto3";
    package your.service.v1;
    option go_package = "github.com/yourorg/yourprotos/gen/go/your/service/v1";
   +
   +import "google/api/annotations.proto";
   +
    message StringMessage {
      string value = 1;
    }

    service YourService {
   -  rpc Echo(StringMessage) returns (StringMessage) {}
   +  rpc Echo(StringMessage) returns (StringMessage) {
   +    option (google.api.http) = {
   +      post: "/v1/example/echo"
   +      body: "*"
   +    };
   +  }
    }
   ```

   > You will need to provide the required third party protobuf files to the `protoc` compiler.
   > They are included in this repo under the `third_party/googleapis` folder, and we recommend copying
   > them into your `protoc` generation file structure. If you've structured your proto files according
   > to something like [the Buf style guide](https://buf.build/docs/style-guide#files-and-packages),
   > you could copy the files into a top-level `./google` folder.

   See [a_bit_of_everything.proto](examples/internal/proto/examplepb/a_bit_of_everything.proto)
   for examples of more annotations you can add to customize gateway behavior
   and generated OpenAPI output.

   Here's what a `protoc` execution might look like:

   ```sh
      protoc -I . --grpc-gateway_out ./gen/go \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        your/service/v1/your_service.proto
   ```

   3. External configuration
      If you do not want to (or cannot) modify the proto file for use with gRPC-Gateway you can
      alternatively use an external
      [gRPC Service Configuration](https://cloud.google.com/endpoints/docs/grpc/grpc-service-config) file.
      [Check our documentation](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/grpc_api_configuration/)
      for more information.

   Here's what a `protoc` execution might look like with this option enabled:

   ```sh
      protoc -I . --grpc-gateway_out ./gen/go \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt grpc_api_configuration=path/to/config.yaml \
        your/service/v1/your_service.proto
   ```

5. Write an entrypoint for the HTTP reverse-proxy server

   ```go
   package main

   import (
     "context"
     "flag"
     "net/http"

     "github.com/golang/glog"
     "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
     "google.golang.org/grpc"

     gw "github.com/yourorg/yourrepo/proto/gen/go/your/service/v1/your_service"  // Update
   )

   var (
     // command-line options:
     // gRPC server endpoint
     grpcServerEndpoint = flag.String("grpc-server-endpoint",  "localhost:9090", "gRPC server endpoint")
   )

   func run() error {
     ctx := context.Background()
     ctx, cancel := context.WithCancel(ctx)
     defer cancel()

     // Register gRPC server endpoint
     // Note: Make sure the gRPC server is running properly and accessible
     mux := runtime.NewServeMux()
     opts := []grpc.DialOption{grpc.WithInsecure()}
     err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
     if err != nil {
       return err
     }

     // Start HTTP server (and proxy calls to gRPC server endpoint)
     return http.ListenAndServe(":8081", mux)
   }

   func main() {
     flag.Parse()
     defer glog.Flush()

     if err := run(); err != nil {
       glog.Fatal(err)
     }
   }
   ```

6. (Optional) Generate OpenAPI definitions using `protoc-gen-openapiv2`

   ```sh
   protoc -I . --openapiv2_out ./gen/openapiv2 --openapiv2_opt logtostderr=true your/service/v1/your_service.proto
   ```

   Note that this plugin also supports generating OpenAPI definitions for unannotated methods; use the `generate_unbound_methods` option to enable this.

## Video intro

This GopherCon UK 2019 presentation from our maintainer [@JohanBrandhorst](https://github.com/johanbrandhorst) provides a good intro to using the gRPC-Gateway. It uses the following boilerplate repo as a base: https://github.com/johanbrandhorst/grpc-gateway-boilerplate.

<div align="center">
<a href="https://www.youtube.com/watch?v=Pq1paKC-fXk">
<img src="https://img.youtube.com/vi/Pq1paKC-fXk/0.jpg" />
</a>
</div>

## Parameters and flags

During code generation with `protoc`, flags to gRPC-Gateway tools must be passed
through `protoc` using one of 2 patterns:

- as part of the `--<tool_suffix>_out` `protoc` parameter: `--<tool_suffix>_out=<flags>:<path>`

```sh
--grpc-gateway_out=logtostderr=true,repeated_path_param_separator=ssv:.
--openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:.
```

- using additional `--<tool_suffix>_opt` parameters: `--<tool_suffix>_opt=<flag>[,<flag>]*`

```sh
--grpc-gateway_opt logtostderr=true,repeated_path_param_separator=ssv
# or separately
--grpc-gateway_opt logtostderr=true --grpc-gateway_opt repeated_path_param_separator=ssv

--openapiv2_opt logtostderr=true,repeated_path_param_separator=ssv
# or separately
--openapiv2_opt logtostderr=true --openapiv2_opt repeated_path_param_separator=ssv
```

`protoc-gen-grpc-gateway` supports custom mapping from Protobuf `import` to
Golang import paths. They are compatible with
[the parameters with the same names in `protoc-gen-go`](https://github.com/golang/protobuf#parameters).

In addition, we also support the `request_context` parameter in order to use the
`http.Request`'s Context (only for Go 1.7 and above). This parameter can be
useful to pass the request-scoped context between the gateway and the gRPC service.

`protoc-gen-grpc-gateway` also supports some more command line flags to control
logging. You can give these flags together with the parameters above. Run
`protoc-gen-grpc-gateway --help` for more details about the flags.

Similarly, `protoc-gen-openapiv2` supports command-line flags to control OpenAPI
output (for example, `json_names_for_fields` to output JSON names for fields
instead of protobuf names). Run `protoc-gen-openapiv2 --help` for more flag
details. Further OpenAPI customization is possible by annotating your `.proto`
files with options from
[openapiv2.proto](protoc-gen-openapiv2/options/openapiv2.proto) - see
[a_bit_of_everything.proto](examples/internal/proto/examplepb/a_bit_of_everything.proto)
for examples.

## More examples

More examples are available under the `examples` directory.

- `proto/examplepb/echo_service.proto`, `proto/examplepb/a_bit_of_everything.proto`, `proto/examplepb/unannotated_echo_service.proto`: service definition
  - `proto/examplepb/echo_service.pb.go`, `proto/examplepb/a_bit_of_everything.pb.go`, `proto/examplepb/unannotated_echo_service.pb.go`: [generated] stub of the service
  - `proto/examplepb/echo_service.pb.gw.go`, `proto/examplepb/a_bit_of_everything.pb.gw.go`, `proto/examplepb/uannotated_echo_service.pb.gw.go`: [generated] reverse proxy for the service
  - `proto/examplepb/unannotated_echo_service.yaml`: gRPC API Configuration for `unannotated_echo_service.proto`
- `server/main.go`: service implementation
- `main.go`: entrypoint of the generated reverse proxy

To use the same port for custom HTTP handlers (e.g. serving `swagger.json`),
gRPC-Gateway, and a gRPC server, see
[this example by CoreOS](https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go)
(and its accompanying [blog post](https://coreos.com/blog/grpc-protobufs-swagger.html)).

## Features

### Supported

- Generating JSON API handlers.
- Method parameters in the request body.
- Method parameters in the request path.
- Method parameters in the query string.
- Enum fields in the path parameter (including repeated enum fields).
- Mapping streaming APIs to newline-delimited JSON streams.
- Mapping HTTP headers with `Grpc-Metadata-` prefix to gRPC metadata (prefixed with `grpcgateway-`)
- Optionally emitting API definitions for
  [OpenAPI (Swagger) v2](https://swagger.io/docs/specification/2-0/basic-structure/).
- Setting [gRPC timeouts](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests)
  through inbound HTTP `Grpc-Timeout` header.
- Partial support for [gRPC API Configuration](https://cloud.google.com/endpoints/docs/grpc/grpc-service-config)
  files as an alternative to annotation.
- Automatically translating PATCH requests into Field Mask gRPC requests. See
  [the docs](https://grpc-ecosystem.github.io/grpc-gateway/docs/patch.html)
  for more information.

### No plan to support

But patches are welcome.

- Method parameters in HTTP headers.
- Handling trailer metadata.
- Encoding request/response body in XML.
- True bi-directional streaming.

## Mapping gRPC to HTTP

- [How gRPC error codes map to HTTP status codes in the response](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15).
- HTTP request source IP is added as `X-Forwarded-For` gRPC request header.
- HTTP request host is added as `X-Forwarded-Host` gRPC request header.
- HTTP `Authorization` header is added as `authorization` gRPC request header.
- Remaining Permanent HTTP header keys (as specified by the IANA
  [here](http://www.iana.org/assignments/message-headers/message-headers.xhtml)
  are prefixed with `grpcgateway-` and added with their values to gRPC request
  header.
- HTTP headers that start with 'Grpc-Metadata-' are mapped to gRPC metadata
  (prefixed with `grpcgateway-`).
- While configurable, the default {un,}marshaling uses
  [jsonpb](https://pkg.go.dev/github.com/golang/protobuf/jsonpb) with
  `OrigName: true`.

## Contribution

See [CONTRIBUTING.md](http://github.com/grpc-ecosystem/grpc-gateway/blob/master/CONTRIBUTING.md).

## License

gRPC-Gateway is licensed under the BSD 3-Clause License.
See [LICENSE.txt](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt) for more details.
