---
layout: default
title: Learn More
parent: Tutorials
nav_order: 6
---

## Learn More

#### google.api.http

GRPC transcoding is a conversion function between the gRPC method and one or more HTTP REST endpoints. This allows developers to create a single API service that supports both the gRPC API and the REST API. Many systems, including the API Google, Cloud Endpoints, gRPC Gateway, and the Envoy proxy server support this feature and use it for large-scale production services.

The grcp-gateway the server is created according to the `google.api.http` annotations in your service definitions.

HttpRule defines the gRPC / REST mapping scheme. The mapping defines how different parts of a gRPC request message are mapped to the URL path, URL request parameters, and HTTP request body. It also controls how the gRPC response message is displayed in the HTTP response body. HttpRule is usually specified as a `google.api.http` annotation in the gRPC method.

Each mapping defines a URL path template and an HTTP method. A path template can refer to one or more fields in a gRPC request message if each field is a non-repeating field with a primitive type. The path template controls how the request message fields are mapped to the URL path.

Read more about HTTP and gRPC Transcoding on https://google.aip.dev/127.
