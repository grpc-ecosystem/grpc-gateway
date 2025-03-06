# Using arbitrary messages in response description

Assuming a protobuf file of the structure

```protobuf
syntax = "proto3";

package example.service.v1;

import "protoc-gen-openapiv2/options/annotations.proto";

service GenericService {
  rpc GenericRPC(GenericRPCRequest) returns (GenericRPCResponse);
}

message GenericRPCRequest {
  string id = 1;
}

message GenericRPCResponse {
  string result = 1;
}
```

If you want your OpenAPI document to include a custom response for all RPCs defined in this protobuf file, you can add the following:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  responses: {
    key: "400"
    value: {
      description: "Returned when the request is malformed."
      schema: {
        json_schema: {ref: ".example.service.v1.GenericResponse"} // Must match the fully qualified name of the message
      }
    }
  }
};

message GenericResponse {
  repeated string resources = 1;
  repeated string errors = 2;
}
```

When generating, you will see the following response included in your OpenAPI document:

```json
"400": {
  "description": "Returned when the request is malformed.",
  "schema": {
    "$ref": "#/definitions/v1GenericResponse"
  }
},
```

The annotation can also be specified per-rpc.