# One way to run the example

```bash
# Handle dependencies (this repository uses Go modules)
$ go -C examples mod download
```

Run the gRPC server and the gateway in two terminals, from the root of the repository:

```bash
$ go -C examples run ./internal/cmd/example-grpc-server
```

```bash
$ go -C examples run ./internal/cmd/example-gateway-server -openapi_dir internal/proto/examplepb
```

Then you can use curl or a browser to test:

```bash
# List all apis
$ curl http://localhost:8080/openapiv2/echo_service.swagger.json

# Visit the apis
$ curl -XPOST http://localhost:8080/v1/example/echo/foo
{"id":"foo"}

$ curl  http://localhost:8080/v1/example/echo/foo/123
{"id":"foo","num":"123"}

```

So you have visited the apis by HTTP successfully. You can also try other apis.

The integration tests in [integration](./integration) start both servers themselves and exercise
the gateway with clients generated from the OpenAPI definitions (see `client_test.go`):

```bash
$ go -C examples test ./internal/integration
```
