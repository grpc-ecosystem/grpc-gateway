---
layout: default
title: Creating a simple hello world with gRPC
parent: Tutorials
nav_order: 2
---

## Creating a simple hello world with gRPC

To understand the gRPC-Gateway we are going to first make a hello world gRPC service.

### Defining your gRPC service using protocol buffers

Before we create a gRPC service, we should create a proto file to define what we need, here we create a file named `hello_world.proto` in the directory `proto/helloworld/hello_world.proto`.

```proto
syntax = "proto3";

package helloworld;

// The greeting service definition
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```
