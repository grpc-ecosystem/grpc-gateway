---
layout: default
title: Learn More
parent: Tutorials
nav_order: 6
---

# Learn More

## How it works

After we use cURL to send HTTP requests `curl` sends a request to the gateway with the post, gateway as proxy forwards the request to `GreeterServer` through gRPC, `GreeterServer` returns the result through gRPC, the gateway receives the result, and JSON returns to the front end.

In this way, the transformation process from HTTP JSON to internal gRPC is completed through gRPC-Gateway.

## google.api.http

Read more about `google.api.http` on [https://github.com/googleapis/googleapis/blob/master/google/api/http.proto](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)

## HTTP and gRPC Transcoding

Read more about HTTP and gRPC Transcoding on [https://google.aip.dev/127](https://google.aip.dev/127)
