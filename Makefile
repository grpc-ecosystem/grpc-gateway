PKG=github.com/gengo/grpc-gateway
GO_PLUGIN=bin/protoc-gen-go
GO_PLUGIN_PKG=github.com/golang/protobuf/protoc-gen-go
GATEWAY_PLUGIN=bin/protoc-gen-grpc-gateway
GATEWAY_PLUGIN_PKG=$(PKG)/protoc-gen-grpc-gateway
GATEWAY_PLUGIN_SRC=protoc-gen-grpc-gateway/main.go \
		   protoc-gen-grpc-gateway/generator.go
OPTIONS_GO=options/options.pb.go
OPTIONS_PROTO=options/options.proto
PKGMAP=Mgoogle/protobuf/descriptor.proto=$(GO_PLUGIN_PKG)/descriptor
EXAMPLES=examples/echo_service.proto
EXAMPLE_SVCSRCS=$(EXAMPLES:.proto=.pb.go)
EXAMPLE_GWSRCS=$(EXAMPLES:.proto=.pb.gw.go)
PROTOC_INC_PATH=$(dir $(shell which protoc))/../include

generate: $(OPTIONS_GO)

.SUFFIXES: .go .proto

$(GO_PLUGIN): 
	go get $(GO_PLUGIN_PKG)
	go build -o $@ $(GO_PLUGIN_PKG)

$(OPTIONS_GO): $(OPTIONS_PROTO) $(GO_PLUGIN)
	protoc -I $(PROTOC_INC_PATH)  -I. --plugin=$(GO_PLUGIN) --go_out=$(PKGMAP):. $(OPTIONS_PROTO)

$(GATEWAY_PLUGIN): $(OPTIONS_GO) $(GATEWAY_PLUGIN_SRC)
	go build -o $@ $(GATEWAY_PLUGIN_PKG)

$(EXAMPLE_SVCSRCS): $(GO_PLUGIN) $(EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. --plugin=$(GO_PLUGIN) --go_out=$(PKGMAP),plugins=grpc:. $(EXAMPLES)
$(EXAMPLE_GWSRCS): $(GATEWAY_PLUGIN) $(EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. --plugin=$(GATEWAY_PLUGIN) --grpc-gateway_out=logtostderr=true:. $(EXAMPLES)

test: $(EXAMPLE_SVCSRCS) $(EXAMPLE_GWSRCS)
	go test $(PKG)/...

realclean:
	rm -f $(OPTIONS_GO)
	rm -f $(EXAMPLE_SVCSRCS)
	rm -f $(EXAMPLE_GWSRCS)
