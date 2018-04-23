# This is a Makefile which maintains files automatically generated but to be
# shipped together with other files.
# You don't have to rebuild these targets by yourself unless you develop
# grpc-gateway itself.

PKG=github.com/grpc-ecosystem/grpc-gateway
GO_PLUGIN=bin/protoc-gen-go
GO_PROTOBUF_REPO=github.com/golang/protobuf
GO_PLUGIN_PKG=$(GO_PROTOBUF_REPO)/protoc-gen-go
GO_PTYPES_ANY_PKG=$(GO_PROTOBUF_REPO)/ptypes/any
SWAGGER_PLUGIN=bin/protoc-gen-swagger
SWAGGER_PLUGIN_SRC= utilities/doc.go \
		    utilities/pattern.go \
		    utilities/trie.go \
		    protoc-gen-swagger/genswagger/generator.go \
		    protoc-gen-swagger/genswagger/template.go \
		    protoc-gen-swagger/main.go
SWAGGER_PLUGIN_PKG=$(PKG)/protoc-gen-swagger
GATEWAY_PLUGIN=bin/protoc-gen-grpc-gateway
GATEWAY_PLUGIN_PKG=$(PKG)/protoc-gen-grpc-gateway
GATEWAY_PLUGIN_SRC= utilities/doc.go \
		    utilities/pattern.go \
		    utilities/trie.go \
		    protoc-gen-grpc-gateway \
		    protoc-gen-grpc-gateway/descriptor \
		    protoc-gen-grpc-gateway/descriptor/registry.go \
		    protoc-gen-grpc-gateway/descriptor/services.go \
		    protoc-gen-grpc-gateway/descriptor/types.go \
		    protoc-gen-grpc-gateway/generator \
		    protoc-gen-grpc-gateway/generator/generator.go \
		    protoc-gen-grpc-gateway/gengateway \
		    protoc-gen-grpc-gateway/gengateway/doc.go \
		    protoc-gen-grpc-gateway/gengateway/generator.go \
		    protoc-gen-grpc-gateway/gengateway/template.go \
		    protoc-gen-grpc-gateway/httprule \
		    protoc-gen-grpc-gateway/httprule/compile.go \
		    protoc-gen-grpc-gateway/httprule/parse.go \
		    protoc-gen-grpc-gateway/httprule/types.go \
		    protoc-gen-grpc-gateway/main.go
GATEWAY_PLUGIN_FLAGS?=
SWAGGER_PLUGIN_FLAGS?=

GOOGLEAPIS_DIR=third_party/googleapis
OUTPUT_DIR=_output

RUNTIME_PROTO=runtime/internal/stream_chunk.proto
RUNTIME_GO=$(RUNTIME_PROTO:.proto=.pb.go)

OPENAPIV2_PROTO=protoc-gen-swagger/options/openapiv2.proto protoc-gen-swagger/options/annotations.proto
OPENAPIV2_GO=$(OPENAPIV2_PROTO:.proto=.pb.go)

PKGMAP=Mgoogle/protobuf/descriptor.proto=$(GO_PLUGIN_PKG)/descriptor,Mexamples/sub/message.proto=$(PKG)/examples/sub
ADDITIONAL_GW_FLAGS=
ifneq "$(GATEWAY_PLUGIN_FLAGS)" ""
	ADDITIONAL_GW_FLAGS=,$(GATEWAY_PLUGIN_FLAGS)
endif
ADDITIONAL_SWG_FLAGS=
ifneq "$(SWAGGER_PLUGIN_FLAGS)" ""
	ADDITIONAL_SWG_FLAGS=,$(SWAGGER_PLUGIN_FLAGS)
endif
SWAGGER_EXAMPLES=examples/examplepb/echo_service.proto \
	 examples/examplepb/a_bit_of_everything.proto \
	 examples/examplepb/wrappers.proto \
	 examples/examplepb/unannotated_echo_service.proto
EXAMPLES=examples/examplepb/echo_service.proto \
	 examples/examplepb/a_bit_of_everything.proto \
	 examples/examplepb/stream.proto \
	 examples/examplepb/flow_combination.proto \
	 examples/examplepb/wrappers.proto \
	 examples/examplepb/unannotated_echo_service.proto
EXAMPLE_SVCSRCS=$(EXAMPLES:.proto=.pb.go)
EXAMPLE_GWSRCS=$(EXAMPLES:.proto=.pb.gw.go)
EXAMPLE_SWAGGERSRCS=$(EXAMPLES:.proto=.swagger.json)
EXAMPLE_DEPS=examples/sub/message.proto examples/sub2/message.proto
EXAMPLE_DEPSRCS=$(EXAMPLE_DEPS:.proto=.pb.go)

EXAMPLE_CLIENT_DIR=examples/clients
ECHO_EXAMPLE_SPEC=examples/examplepb/echo_service.swagger.json
ECHO_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/echo/EchoServiceApi.go \
		  $(EXAMPLE_CLIENT_DIR)/echo/ExamplepbSimpleMessage.go
ABE_EXAMPLE_SPEC=examples/examplepb/a_bit_of_everything.swagger.json
ABE_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/abe/ABitOfEverythingServiceApi.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/ABitOfEverythingNested.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/ExamplepbABitOfEverything.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/ExamplepbNumericEnum.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/ExamplepbIdMessage.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/NestedDeepEnum.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/ProtobufEmpty.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/Sub2IdMessage.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/SubStringMessage.go
UNANNOTATED_ECHO_EXAMPLE_SPEC=examples/examplepb/unannotated_echo_service.swagger.json
UNANNOTATED_ECHO_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/unannotatedecho/UnannotatedEchoServiceApi.go \
		 $(EXAMPLE_CLIENT_DIR)/unannotatedecho/ExamplepbUnannotatedSimpleMessage.go
EXAMPLE_CLIENT_SRCS=$(ECHO_EXAMPLE_SRCS) $(ABE_EXAMPLE_SRCS) $(UNANNOTATED_ECHO_EXAMPLE_SRCS)
SWAGGER_CODEGEN=swagger-codegen

PROTOC_INC_PATH=$(dir $(shell which protoc))/../include

generate: $(RUNTIME_GO)

.SUFFIXES: .go .proto

$(GO_PLUGIN):
	go get $(GO_PLUGIN_PKG)
	go build -o $@ $(GO_PLUGIN_PKG)

$(RUNTIME_GO): $(RUNTIME_PROTO) $(GO_PLUGIN)
	protoc -I $(PROTOC_INC_PATH) --plugin=$(GO_PLUGIN) -I $(GOPATH)/src/$(GO_PTYPES_ANY_PKG) -I. --go_out=$(PKGMAP):. $(RUNTIME_PROTO)

$(OPENAPIV2_GO): $(OPENAPIV2_PROTO) $(GO_PLUGIN)
	protoc -I $(PROTOC_INC_PATH) --plugin=$(GO_PLUGIN) -I. --go_out=$(PKGMAP):$(GOPATH)/src $(OPENAPIV2_PROTO)

$(GATEWAY_PLUGIN): $(RUNTIME_GO) $(GATEWAY_PLUGIN_SRC)
	go build -o $@ $(GATEWAY_PLUGIN_PKG)

$(SWAGGER_PLUGIN): $(SWAGGER_PLUGIN_SRC) $(OPENAPIV2_GO)
	go build -o $@ $(SWAGGER_PLUGIN_PKG)

$(EXAMPLE_SVCSRCS): $(GO_PLUGIN) $(EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GO_PLUGIN) --go_out=$(PKGMAP),plugins=grpc:. $(EXAMPLES)
$(EXAMPLE_DEPSRCS): $(GO_PLUGIN) $(EXAMPLE_DEPS)
	mkdir -p $(OUTPUT_DIR)
	protoc -I $(PROTOC_INC_PATH) -I. --plugin=$(GO_PLUGIN) --go_out=$(PKGMAP),plugins=grpc:$(OUTPUT_DIR) $(@:.pb.go=.proto)
	cp $(OUTPUT_DIR)/$(PKG)/$@ $@ || cp $(OUTPUT_DIR)/$@ $@

$(EXAMPLE_GWSRCS): ADDITIONAL_GW_FLAGS:=$(ADDITIONAL_GW_FLAGS),grpc_api_configuration=examples/examplepb/unannotated_echo_service.yaml
$(EXAMPLE_GWSRCS): $(GATEWAY_PLUGIN) $(EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GATEWAY_PLUGIN) --grpc-gateway_out=logtostderr=true,$(PKGMAP)$(ADDITIONAL_GW_FLAGS):. $(EXAMPLES)

$(EXAMPLE_SWAGGERSRCS): ADDITIONAL_SWG_FLAGS:=$(ADDITIONAL_SWG_FLAGS),grpc_api_configuration=examples/examplepb/unannotated_echo_service.yaml
$(EXAMPLE_SWAGGERSRCS): $(SWAGGER_PLUGIN) $(SWAGGER_EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(SWAGGER_PLUGIN) --swagger_out=logtostderr=true,$(PKGMAP)$(ADDITIONAL_SWG_FLAGS):. $(SWAGGER_EXAMPLES)

$(ECHO_EXAMPLE_SRCS): $(ECHO_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(ECHO_EXAMPLE_SPEC) \
	    -l go -o examples/clients/echo --additional-properties packageName=echo
	@rm -f $(EXAMPLE_CLIENT_DIR)/echo/README.md \
		$(EXAMPLE_CLIENT_DIR)/echo/git_push.sh \
		$(EXAMPLE_CLIENT_DIR)/echo/.travis.yml
$(ABE_EXAMPLE_SRCS): $(ABE_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(ABE_EXAMPLE_SPEC) \
	    -l go -o examples/clients/abe --additional-properties packageName=abe
	@rm -f $(EXAMPLE_CLIENT_DIR)/abe/README.md \
		$(EXAMPLE_CLIENT_DIR)/abe/git_push.sh \
		$(EXAMPLE_CLIENT_DIR)/abe/.travis.yml
$(UNANNOTATED_ECHO_EXAMPLE_SRCS): $(UNANNOTATED_ECHO_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(UNANNOTATED_ECHO_EXAMPLE_SPEC) \
	    -l go -o examples/clients/unannotatedecho --additional-properties packageName=unannotatedecho
	@rm -f $(EXAMPLE_CLIENT_DIR)/unannotatedecho/README.md \
		$(EXAMPLE_CLIENT_DIR)/unannotatedecho/git_push.sh \
		$(EXAMPLE_CLIENT_DIR)/unannotatedecho/.travis.yml

examples: $(EXAMPLE_SVCSRCS) $(EXAMPLE_GWSRCS) $(EXAMPLE_DEPSRCS) $(EXAMPLE_SWAGGERSRCS) $(EXAMPLE_CLIENT_SRCS)
test: examples
	go test -race $(PKG)/...
	go test -race $(PKG)/examples -args -network=unix -endpoint=test.sock

lint:
	golint --set_exit_status $(PKG)/runtime
	golint --set_exit_status $(PKG)/utilities/...
	golint --set_exit_status $(PKG)/protoc-gen-grpc-gateway/...
	golint --set_exit_status $(PKG)/protoc-gen-swagger/...
	go vet $(PKG)/runtime || true
	go vet $(PKG)/utilities/...
	go vet $(PKG)/protoc-gen-grpc-gateway/...
	go vet $(PKG)/protoc-gen-swagger/...

clean distclean:
	rm -f $(GATEWAY_PLUGIN)
realclean: distclean
	rm -f $(EXAMPLE_SVCSRCS) $(EXAMPLE_DEPSRCS)
	rm -f $(EXAMPLE_GWSRCS)
	rm -f $(EXAMPLE_SWAGGERSRCS)
	rm -f $(GO_PLUGIN)
	rm -f $(SWAGGER_PLUGIN)
	rm -f $(EXAMPLE_CLIENT_SRCS)
	rm -f $(OPENAPIV2_GO)

.PHONY: generate examples test lint clean distclean realclean
