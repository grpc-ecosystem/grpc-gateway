# **Issue Address**: [Link Here](https://github.com/grpc-ecosystem/grpc-gateway/issues/5241)



# A complete example in the issue.

First enter command`touch your_service.proto`.

I will provide an example from the issue where unbounded messages are not used in the Swagger display.

```protobuf
syntax = "proto3";
package your.service.v1;
option go_package = "testTool/v1";

import "protoc-gen-openapiv2/options/annotations.proto";
message GenericResponse {
  repeated string resources = 1;
  repeated string errors = 2;
}

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  responses: {
    key: "400"
    value: {
      description: "Returned when the request is malformed."
      schema: {
        json_schema: {ref: "#/your.service.v1/GenericResponse"}
      }
    }
  }
};
```

Next, we will use commands to generate the Swagger file and the `pb.go` file.

```sh
 protoc -I. -I./grpc-gateway -I./googleapis/google \
    --go_out=protos --go_opt=paths=source_relative \
    --go-grpc_out=protos --go-grpc_opt=paths=source_relative \
    --openapiv2_out=protos --openapiv2_opt=allow_merge=true,merge_file_name=your_service_swagger \
    your_service.proto
```

> ```bash
> git clone https://github.com/googleapis/googleapis.git
> ```
>
> You should clone it before command enter, and replace ***/googleapis/goog*le** with your file address.

All files will be generated in the `protos` folder under the current directory. So make sure that this folder exists in the current directory. If it doesn't, please enter the command `mkdir protos` first.

We will see this:

```sh
├── protos
│   ├── your_service.pb.go
│   └── your_service_swagger.swagger.json
```

In your_service_swagger.swagger.json,we can't search anything about GenericResponse

# A right example of rendering `GenericResponse`.

I will make some change to render GenericResponse

```protobuf
syntax = "proto3";
package example.service.v1;
option go_package = "test/v1";

import "protoc-gen-openapiv2/options/annotations.proto";
import "google/api/annotations.proto";
message GenericResponse {
  repeated string resources = 1;
  repeated string errors = 2;
}

message Something{
  string id = 1;
}

service ThisService {
  rpc CreateBody(Something) returns (Something){
    option (google.api.http) = {
      post: "/v1/example/a_bit_of_everything"
      body: "*"
    };
  };
};

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  responses: {
    key: "400"
    value: {
      description: "Returned when the request is malformed."
      schema: {
        json_schema: {ref: ".example.service.v1.GenericResponse"}
      }
    }
  }
};
```

In this example, `ref` has been changed into`.example.service.v1.GenericResponse`

And I've written an interface named `ThisService`, and specified the request method and request body for it.

So we can reproduce `your_service_swagger.swagger.json`

This time,we will see this in`your_service_swagger.swagger.json`

```json
"400": {
  "description": "Returned when the request is malformed.",
  "schema": {
    "$ref": "#/definitions/v1GenericResponse"
  }
},
```

Now we can see `GenericResponse` in swagger



