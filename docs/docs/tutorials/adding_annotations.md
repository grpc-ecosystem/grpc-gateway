---
layout: default
title: Adding the grpc-gateway annotations to an existing protobuf file
parent: Tutorials
nav_order: 6
---

## Adding the grpc-gateway annotations to an existing protobuf file

Start the greeter_server service first, and then start the gateway. Then gateway connects to greeter_server and establishes http monitoring.

Then we use curl to send http requests:

```sh
curl -X POST -k http://localhost:8080/v1/example/echo -d '{"name": " world"}
```

```
{"message":"Hello  world"}
```

The process is as follows:

curl sends a request to the gateway with the post, gateway as proxy forwards the request to greeter_server through grpc, greeter_server returns the result through grpc, the gateway receives the result, and json returns to the front end.

In this way, the transformation process from http json to internal grpc is completed through grpc-gateway.
