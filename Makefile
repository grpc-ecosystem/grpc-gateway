GO_PLUGIN=bin/protoc-gen-go
GO_PLUGIN_PKG=github.com/golang/protobuf/protoc-gen-go
OPTIONS_GO=options/options.pb.go
OPTIONS_PROTO=options/options.proto
PKGMAP=Mgoogle/protobuf/descriptor.proto=$(GO_PLUGIN_PKG)/descriptor

generate: $(OPTIONS_GO)

.SUFFIXES: .go .proto

$(GO_PLUGIN): 
	go get $(GO_PLUGIN_PKG)
	go build -o $@ $(GO_PLUGIN_PKG)

$(OPTIONS_GO): $(OPTIONS_PROTO) $(GO_PLUGIN)
	protoc -I $(dir $(shell which protoc))/../include -I. --plugin=$(GO_PLUGIN) --go_out=$(PKGMAP):. $(OPTIONS_PROTO)

realcelan:
	rm -f $(OPTIONS_GO)
