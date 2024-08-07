# Custom Query Parameter Parsing in gRPC-Gateway
The gRPC-Gateway allows you to customize the way query parameters are parsed and mapped to your gRPC request messages. This can be particularly useful when dealing with nested fields or complex query parameter mappings that are not directly supported by the default behavior.

## Example: Custom Query Parameter Parser for Nested Fields
Suppose you have a gRPC service definition where you want to map query parameters to a nested message field. Here is how you can achieve this using a custom query parameter parser.

## Protobuf Definition:
```protobuf
syntax = "proto3";

import "google/api/annotations.proto";

package your.service.v1;

message PageOptions {
  int32 limit = 1;
  int32 page = 2;
}

message ListStuffRequest {
  string stuff_uuid = 1;
  PageOptions pagination = 2;
}

message ListStuffResponse {
  repeated string stuff = 1;
}

service MyService {
  rpc ListStuff(ListStuffRequest) returns (ListStuffResponse) {
    option (google.api.http) = {
      get: "/path/to/{stuff_uuid}/stuff"
    };
  }
}
```

## The Problem
The ListStuffRequest message contains a nested PageOptions message that represents pagination options for the list operation. By default, creating a ServerMux in gRPC-Gateway will be able to parse query parameters for the `stuff_uuid` and other basic fields (that are not your own custom message). Given that `PageOptions` is a custom message, and therefore will lead to a nested RequestMessage (the root message), it can only parse query parameters for the `PageOptions` message if they are in the format `url?pagination.limit=x`. If your endpoint constraints require query parameters to be in the format `url?limit=x`, without specifying the full qualified path in the url, you will need to implement a custom parser to map these parameters correctly into the nested PageOptions message.

## Custom Query Parameter Parser

Create a custom [QueryParameterParser](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/query.go#L30) to handle the nested PageOptions message.
```go
package customparser

import (
	"net/url"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"

	your_service_v1 "path/to/your/service/v1"
)

// CustomQueryParameterParser parses query parameters into the appropriate gRPC message fields.
type CustomQueryParameterParser struct{}

// Parse parses query parameters and populates the appropriate fields in the gRPC request message.
func (p *CustomQueryParameterParser) Parse(target proto.Message, values url.Values, filter *utilities.DoubleArray) error {
	switch req := target.(type) {
	// Different messages/requests can have different parsers, of course
	case *your_service_v1.ListStuffRequest:
		return populateListStuffParams(values, req)
	}
	
	return (runtime.DefaultQueryParser{}).Parse(target, values, filter)
}

// populateListStuffParams populates the ListStuffRequest with query parameters.
func populateListStuffParams(values url.Values, r *your_service_v1.ListStuffRequest) error {
	pageOptions := &your_service_v1.PageOptions{}
	
	if limit := values.Get("limit"); limit != "" {
		if parsedLimit, err := strconv.Atoi(limit); err == nil {
			pageOptions.Limit = int32(parsedLimit)
		}
	}
	if page := values.Get("page"); page != "" {
		if parsedPage, err := strconv.Atoi(page); err == nil {
			pageOptions.Page = int32(parsedPage)
		}
	}

	r.Pagination = pageOptions
	return nil
}
```
## Integrate Your Custom Query Parameter Parser in gRPC-Gateway Setup

All you need to do now is update the gRPC-Gateway setup, particularly wherever you are [defining the mux server](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/mux.go#L293), to [use the custom query parameter parser](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/mux.go#L110).
```go
package main

import (
	"context"
	"net/http"
	
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"your_module/path/customparser"
)

// create a new ServeMux with custom parser and other runtime options as needed
func createGRPCGatewayMux() *runtime.ServeMux {
	// whatever custom code you may need before you create the mux...
	
	return runtime.NewServeMux(
		// Custom query parameter parser
		runtime.SetQueryParameterParser(&customparser.CustomQueryParameterParser{}),
		
		// other runtime options you may need...
	)
}
```