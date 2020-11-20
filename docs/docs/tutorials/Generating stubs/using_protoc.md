---
layout: default
title: Generating stubs using protoc
parent: Generating stubs
grand_parent: Tutorials
nav_order: 2
---

## Generating stubs using protoc

Generate gRPC stubs:

This step generates the gRPC stubs that you can use to implement the service and consume from clients.

Here's an example of what a `protoc` command might look like to generate Go stubs:

```sh
protoc -I . \
   --go_out ./gen/go/ --go_opt paths=source_relative \
   --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative \
   your/service/v1/your_service.proto
```
