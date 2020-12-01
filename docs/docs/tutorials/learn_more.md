---
layout: default
title: Learn More
parent: Tutorials
nav_order: 6
---

# Learn More

## How it works

When the HTTP request arrives at the gRPC-gateway, it parses the JSON data into a protobuf message. It then makes a normal Go gRPC client request using the parsed protobuf message. The Go gRPC client encodes the protobuf structure into the protobuf binary format and sends it to the gRPC server. The gRPC Server handles the request and returns the response in the protobuf binary format. The Go gRPC client parses it into a protobuf message and returns it to the gRPC-gateway, which encodes the protobuf message to JSON and returns it to the original client.

## google.api.http

Read more about `google.api.http` on [https://github.com/googleapis/googleapis/blob/master/google/api/http.proto](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)

## HTTP and gRPC Transcoding

Read more about HTTP and gRPC Transcoding on [https://google.aip.dev/127](https://google.aip.dev/127)
