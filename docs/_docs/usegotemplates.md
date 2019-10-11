| Title                                  | Category      |
| -------------------------------------- | ------------- |
| Use go templates in protofile comments | Documentation |

# Use go templates in protofile comments

Use [Go templates](https://golang.org/pkg/text/template/ "Package template") in your protofile comments to allow more advanced documentation such as:  
* Documentation about fields in the proto objects.  
* Import the content of external files (such as [Markdown](https://en.wikipedia.org/wiki/Markdown "Markdown Github")). 

## How to use it

By default this function is turned off, so if you want to use it you have to set the `use_go_templates` flag to true inside of the `swagger_out` flag.
```bash
--swagger_out=use_go_templates=true:.
```

### Example script

Example of a bash script with the `use_go_templates` flag set to true:

```bash
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  --go_out=plugins=grpc:. \
  --grpc-gateway_out=logtostderr=true:. \
  --swagger_out=logtostderr=true,use_go_templates=true:. \
  *.proto 
```

### Example proto file

Example of a protofile with Go templates. This proto file imports documentation from another file, `tables.md`:
```protobuf
service LoginService {
    // Login
    // 
    // {{.MethodDescriptorProto.Name}} is a call with the method(s) {{$first := true}}{{range .Bindings}}{{if $first}}{{$first = false}}{{else}}, {{end}}{{.HTTPMethod}}{{end}} within the "{{.Service.Name}}" service.
    // It takes in "{{.RequestType.Name}}" and returns a "{{.ResponseType.Name}}".
    //
    // {{import "tables.md"}}
    rpc Login (LoginRequest) returns (LoginReply) {
        option (google.api.http) = {
            post: "/v1/example/login"
            body: "*"
        };
    }
}

message LoginRequest {
    // The entered username 
    string username = 1;
    // The entered password
    string password = 2;
}

message LoginReply {
    // Whether you have access or not
    bool access = 1;
}
```

The content of `tables.md`:

```markdown
## {{.RequestType.Name}}
| Field ID    | Name      | Type                                                       | Description                  |
| ----------- | --------- | ---------------------------------------------------------  | ---------------------------- | {{range .RequestType.Fields}}
| {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}  
 
## {{.ResponseType.Name}}
| Field ID    | Name      | Type                                                       | Description                  |
| ----------- | --------- | ---------------------------------------------------------- | ---------------------------- | {{range .ResponseType.Fields}}
| {{.Number}} | {{.Name}} | {{if eq .Label.String "LABEL_REPEATED"}}[]{{end}}{{.Type}} | {{fieldcomments .Message .}} | {{end}}  
```

## Swagger output

### SwaggerUI

This is how the swagger file would be rendered in [SwaggerUI](https://swagger.io/tools/swagger-ui/ "SwaggerUI site")

![Screenshot swaggerfile in SwaggerUI](../_imgs/gotemplates/swaggerui.png "SwaggerUI")

### Postman

This is how the swagger file would be rendered in [Postman](https://www.getpostman.com/ "Postman site")

![Screenshot swaggerfile in Postman](../_imgs/gotemplates/postman.png "Postman")

For a more detailed example of a protofile that has Go templates enabled, [click here](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/proto/examplepb/use_go_template.proto "Example protofile with Go template").
