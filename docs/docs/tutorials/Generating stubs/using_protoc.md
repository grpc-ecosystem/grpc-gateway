---
layout: default
title: Generating stubs using protoc
parent: Generating stubs
grand_parent: Tutorials
nav_order: 2
---

## Generating stubs using protoc

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

   This step generates the gRPC stubs that you can use to implement the service and consume from clients.

   Here's an example of what a `protoc` command might look like to generate Go stubs:

   ```sh
   protoc -I . \
      --go_out ./gen/go/ --go_opt paths=source_relative \
      --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative \
      your/service/v1/your_service.proto
   ```
