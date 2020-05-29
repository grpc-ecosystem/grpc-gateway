# This is a Makefile which maintains files automatically generated but to be
# shipped together with other files.
# You don't have to rebuild these targets by yourself unless you develop
# grpc-gateway itself.

GO_PLUGIN=bin/protoc-gen-go
GO_PLUGIN_PKG=google.golang.org/protobuf/cmd/protoc-gen-go
GO_GRPC_PLUGIN=bin/protoc-gen-go-grpc
GO_GRPC_PLUGIN_PKG=google.golang.org/grpc/cmd/protoc-gen-go-grpc
OPENAPI_PLUGIN=bin/protoc-gen-openapiv2
OPENAPI_PLUGIN_SRC= ./utilities/doc.go \
			./utilities/pattern.go \
			./utilities/trie.go \
			protoc-gen-openapiv2/internal/genopenapi/generator.go \
			protoc-gen-openapiv2/internal/genopenapi/template.go \
			protoc-gen-openapiv2/main.go
OPENAPI_PLUGIN_PKG=./protoc-gen-openapiv2
GATEWAY_PLUGIN=bin/protoc-gen-grpc-gateway
GATEWAY_PLUGIN_PKG=./protoc-gen-grpc-gateway
GATEWAY_PLUGIN_SRC= ./utilities/doc.go \
			./utilities/pattern.go \
			./utilities/trie.go \
			./internal/descriptor \
			./internal/descriptor/registry.go \
			./internal/descriptor/services.go \
			./internal/descriptor/types.go \
			./internal/descriptor/grpc_api_configuration.go \
			./internal/generator \
			./internal/generator/generator.go \
			protoc-gen-grpc-gateway \
			protoc-gen-grpc-gateway/internal/gengateway \
			protoc-gen-grpc-gateway/internal/gengateway/doc.go \
			protoc-gen-grpc-gateway/internal/gengateway/generator.go \
			protoc-gen-grpc-gateway/internal/gengateway/template.go \
			internal/httprule \
			internal/httprule/compile.go \
			internal/httprule/parse.go \
			internal/httprule/types.go \
			protoc-gen-grpc-gateway/main.go
GATEWAY_PLUGIN_FLAGS?=
OPENAPI_PLUGIN_FLAGS?=

GOOGLEAPIS_DIR=third_party/googleapis
OUTPUT_DIR=_output

OPENAPIV2_PROTO=protoc-gen-openapiv2/options/openapiv2.proto protoc-gen-openapiv2/options/annotations.proto
OPENAPIV2_GO=$(OPENAPIV2_PROTO:.proto=.pb.go)

ADDITIONAL_GW_FLAGS=
ifneq "$(GATEWAY_PLUGIN_FLAGS)" ""
	ADDITIONAL_GW_FLAGS=,$(GATEWAY_PLUGIN_FLAGS)
endif
ADDITIONAL_SWG_FLAGS=
ifneq "$(OPENAPI_PLUGIN_FLAGS)" ""
	ADDITIONAL_SWG_FLAGS=,$(OPENAPI_PLUGIN_FLAGS)
endif
OPENAPI_EXAMPLES=examples/internal/proto/examplepb/echo_service.proto \
	 examples/internal/proto/examplepb/a_bit_of_everything.proto \
	 examples/internal/proto/examplepb/wrappers.proto \
	 examples/internal/proto/examplepb/stream.proto \
	 examples/internal/proto/examplepb/unannotated_echo_service.proto \
	 examples/internal/proto/examplepb/use_go_template.proto \
	 examples/internal/proto/examplepb/response_body_service.proto

EXAMPLES=examples/internal/proto/examplepb/echo_service.proto \
	 examples/internal/proto/examplepb/a_bit_of_everything.proto \
	 examples/internal/proto/examplepb/stream.proto \
	 examples/internal/proto/examplepb/flow_combination.proto \
	 examples/internal/proto/examplepb/non_standard_names.proto \
	 examples/internal/proto/examplepb/wrappers.proto \
	 examples/internal/proto/examplepb/unannotated_echo_service.proto \
	 examples/internal/proto/examplepb/use_go_template.proto \
	 examples/internal/proto/examplepb/response_body_service.proto

STANDALONE_EXAMPLES=examples/internal/proto/examplepb/unannotated_echo_service.proto

HELLOWORLD=examples/internal/helloworld/helloworld.proto

EXAMPLE_SVCSRCS=$(EXAMPLES:.proto=.pb.go)
EXAMPLE_GWSRCS=$(EXAMPLES:.proto=.pb.gw.go)
EXAMPLE_OPENAPISRCS=$(OPENAPI_EXAMPLES:.proto=.swagger.json)
EXAMPLE_DEPS=examples/internal/proto/pathenum/path_enum.proto examples/internal/proto/sub/message.proto examples/internal/proto/sub2/message.proto
EXAMPLE_DEPSRCS=$(EXAMPLE_DEPS:.proto=.pb.go)

HELLOWORLD_SVCSRCS=$(HELLOWORLD:.proto=.pb.go)
HELLOWORLD_GWSRCS=$(HELLOWORLD:.proto=.pb.gw.go)

RUNTIME_TEST_PROTO=runtime/internal/examplepb/example.proto \
	runtime/internal/examplepb/proto2.proto \
	runtime/internal/examplepb/proto3.proto \
	runtime/internal/examplepb/non_standard_names.proto
RUNTIME_TEST_SRCS=$(RUNTIME_TEST_PROTO:.proto=.pb.go)

APICONFIG_PROTO=internal/descriptor/apiconfig/apiconfig.proto
APICONFIG_SRCS=$(APICONFIG_PROTO:.proto=.pb.go)

EXAMPLE_CLIENT_DIR=examples/internal/clients
ECHO_EXAMPLE_SPEC=examples/internal/proto/examplepb/echo_service.swagger.json
ECHO_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/echo/client.go \
		  $(EXAMPLE_CLIENT_DIR)/echo/response.go \
		  $(EXAMPLE_CLIENT_DIR)/echo/configuration.go \
		  $(EXAMPLE_CLIENT_DIR)/echo/api_echo_service.go \
		  $(EXAMPLE_CLIENT_DIR)/echo/model_examplepb_simple_message.go \
		  $(EXAMPLE_CLIENT_DIR)/echo/model_examplepb_embedded.go
ABE_EXAMPLE_SPEC=examples/internal/proto/examplepb/a_bit_of_everything.swagger.json
ABE_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/abe/model_a_bit_of_everything_nested.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/api_a_bit_of_everything_service.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/client.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/api_camel_case_service_name.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/configuration.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/api_echo_rpc.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_examplepb_a_bit_of_everything.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_examplepb_a_bit_of_everything_repeated.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_examplepb_body.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_examplepb_numeric_enum.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_examplepb_update_v2_request.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_message_path_enum_nested_path_enum.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_nested_deep_enum.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_pathenum_path_enum.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/model_protobuf_field_mask.go \
		 $(EXAMPLE_CLIENT_DIR)/abe/response.go
UNANNOTATED_ECHO_EXAMPLE_SPEC=examples/internal/proto/examplepb/unannotated_echo_service.swagger.json
UNANNOTATED_ECHO_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/unannotatedecho/client.go \
		 $(EXAMPLE_CLIENT_DIR)/unannotatedecho/response.go \
		 $(EXAMPLE_CLIENT_DIR)/unannotatedecho/configuration.go \
		 $(EXAMPLE_CLIENT_DIR)/unannotatedecho/model_examplepb_unannotated_simple_message.go \
		 $(EXAMPLE_CLIENT_DIR)/unannotatedecho/api_unannotated_echo_service.go
RESPONSE_BODY_EXAMPLE_SPEC=examples/internal/proto/examplepb/response_body_service.swagger.json
RESPONSE_BODY_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/responsebody/client.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/response.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/configuration.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/model_examplepb_repeated_response_body_out.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/model_examplepb_repeated_response_body_out_response.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/model_examplepb_repeated_response_strings.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/model_examplepb_response_body_out.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/model_examplepb_response_body_out_response.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/model_response_response_type.go \
		 $(EXAMPLE_CLIENT_DIR)/responsebody/api_response_body_service.go

EXAMPLE_CLIENT_SRCS=$(ECHO_EXAMPLE_SRCS) $(ABE_EXAMPLE_SRCS) $(UNANNOTATED_ECHO_EXAMPLE_SRCS) $(RESPONSE_BODY_EXAMPLE_SRCS)
SWAGGER_CODEGEN=swagger-codegen

PROTOC_INC_PATH=$(dir $(shell which protoc))/../include

.SUFFIXES: .go .proto

$(GO_PLUGIN):
	go build -o $(GO_PLUGIN) $(GO_PLUGIN_PKG)

$(GO_GRPC_PLUGIN):
	go build -o $(GO_GRPC_PLUGIN) $(GO_GRPC_PLUGIN_PKG)

$(OPENAPIV2_GO): $(OPENAPIV2_PROTO) $(GO_PLUGIN)
	protoc -I $(PROTOC_INC_PATH) --plugin=$(GO_PLUGIN) -I. --go_out=paths=source_relative:. $(OPENAPIV2_PROTO)

$(GATEWAY_PLUGIN): $(GATEWAY_PLUGIN_SRC)
	go build -o $@ $(GATEWAY_PLUGIN_PKG)

$(OPENAPI_PLUGIN): $(OPENAPI_PLUGIN_SRC) $(OPENAPIV2_GO)
	go build -o $@ $(OPENAPI_PLUGIN_PKG)

$(EXAMPLE_SVCSRCS): $(GO_PLUGIN) $(GO_GRPC_PLUGIN) $(EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GO_PLUGIN) --plugin=$(GO_GRPC_PLUGIN) --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(EXAMPLES)
$(EXAMPLE_DEPSRCS): $(GO_PLUGIN) $(GO_GRPC_PLUGIN) $(EXAMPLE_DEPS)
	mkdir -p $(OUTPUT_DIR)
	protoc -I $(PROTOC_INC_PATH) -I. --plugin=$(GO_PLUGIN) --plugin=$(GO_GRPC_PLUGIN) --go_out=paths=source_relative:$(OUTPUT_DIR) --go-grpc_out=paths=source_relative:$(OUTPUT_DIR) $(@:.pb.go=.proto)
	cp $(OUTPUT_DIR)/$@ $@ || cp $(OUTPUT_DIR)/$@ $@

$(RUNTIME_TEST_SRCS): $(GO_PLUGIN) $(GO_GRPC_PLUGIN) $(RUNTIME_TEST_PROTO)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GO_PLUGIN) --plugin=$(GO_GRPC_PLUGIN) --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(RUNTIME_TEST_PROTO)

$(APICONFIG_SRCS): $(GO_PLUGIN) $(APICONFIG_PROTO)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GO_PLUGIN) --go_out=paths=source_relative:. $(APICONFIG_PROTO)

$(EXAMPLE_GWSRCS): ADDITIONAL_GW_FLAGS:=$(ADDITIONAL_GW_FLAGS),grpc_api_configuration=examples/internal/proto/examplepb/unannotated_echo_service.yaml
$(EXAMPLE_GWSRCS): ADDITIONAL_SA_FLAGS:=,standalone=true,grpc_api_configuration=examples/internal/proto/examplepb/standalone_echo_service.yaml
$(EXAMPLE_GWSRCS): $(GATEWAY_PLUGIN) $(EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GATEWAY_PLUGIN) --grpc-gateway_out=logtostderr=true,allow_repeated_fields_in_body=true,paths=source_relative$(ADDITIONAL_SA_FLAGS):. $(STANDALONE_EXAMPLES)
	mv examples/internal/proto/examplepb/unannotated_echo_service.pb.gw.go examples/internal/proto/standalone/
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GATEWAY_PLUGIN) --grpc-gateway_out=logtostderr=true,allow_repeated_fields_in_body=true,paths=source_relative$(ADDITIONAL_GW_FLAGS):. $(EXAMPLES)


$(EXAMPLE_OPENAPISRCS): ADDITIONAL_SWG_FLAGS:=$(ADDITIONAL_SWG_FLAGS),grpc_api_configuration=examples/internal/proto/examplepb/unannotated_echo_service.yaml
$(EXAMPLE_OPENAPISRCS): $(OPENAPI_PLUGIN) $(OPENAPI_EXAMPLES)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(OPENAPI_PLUGIN) --openapiv2_out=logtostderr=true,allow_repeated_fields_in_body=true,use_go_templates=true$(ADDITIONAL_SWG_FLAGS):. $(OPENAPI_EXAMPLES)

$(HELLOWORLD_SVCSRCS): $(GO_PLUGIN) $(GO_GRPC_PLUGIN) $(HELLOWORLD)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GO_PLUGIN) --plugin=$(GO_GRPC_PLUGIN) --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(HELLOWORLD)

$(HELLOWORLD_GWSRCS):
$(HELLOWORLD_GWSRCS): $(GATEWAY_PLUGIN) $(HELLOWORLD)
	protoc -I $(PROTOC_INC_PATH) -I. -I$(GOOGLEAPIS_DIR) --plugin=$(GATEWAY_PLUGIN) --grpc-gateway_out=logtostderr=true,allow_repeated_fields_in_body=true,paths=source_relative$(ADDITIONAL_GW_FLAGS):. $(HELLOWORLD)


$(ECHO_EXAMPLE_SRCS): $(ECHO_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(ECHO_EXAMPLE_SPEC) \
		-l go -o examples/internal/clients/echo --additional-properties packageName=echo
	@rm -f $(EXAMPLE_CLIENT_DIR)/echo/README.md \
		$(EXAMPLE_CLIENT_DIR)/echo/git_push.sh
$(ABE_EXAMPLE_SRCS): $(ABE_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(ABE_EXAMPLE_SPEC) \
		-l go -o examples/internal/clients/abe --additional-properties packageName=abe
	@rm -f $(EXAMPLE_CLIENT_DIR)/abe/README.md \
		$(EXAMPLE_CLIENT_DIR)/abe/git_push.sh
$(UNANNOTATED_ECHO_EXAMPLE_SRCS): $(UNANNOTATED_ECHO_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(UNANNOTATED_ECHO_EXAMPLE_SPEC) \
		-l go -o examples/internal/clients/unannotatedecho --additional-properties packageName=unannotatedecho
	@rm -f $(EXAMPLE_CLIENT_DIR)/unannotatedecho/README.md \
		$(EXAMPLE_CLIENT_DIR)/unannotatedecho/git_push.sh
$(RESPONSE_BODY_EXAMPLE_SRCS): $(RESPONSE_BODY_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(RESPONSE_BODY_EXAMPLE_SPEC) \
		-l go -o examples/internal/clients/responsebody --additional-properties packageName=responsebody
	@rm -f $(EXAMPLE_CLIENT_DIR)/responsebody/README.md \
		$(EXAMPLE_CLIENT_DIR)/responsebody/git_push.sh

examples: $(EXAMPLE_DEPSRCS) $(EXAMPLE_SVCSRCS) $(EXAMPLE_GWSRCS) $(EXAMPLE_OPENAPISRCS) $(EXAMPLE_CLIENT_SRCS) $(HELLOWORLD_SVCSRCS) $(HELLOWORLD_GWSRCS)
testproto: $(RUNTIME_TEST_SRCS) $(APICONFIG_SRCS)
test: examples testproto
	go test -short -race ./...
	go test -race ./examples/internal/integration -args -network=unix -endpoint=test.sock
changelog:
	docker run --rm \
		--interactive \
		--tty \
		-e "CHANGELOG_GITHUB_TOKEN=${CHANGELOG_GITHUB_TOKEN}" \
		-v "$(PWD):/usr/local/src/your-app" \
		ferrarimarco/github-changelog-generator:1.14.3 \
				-u grpc-ecosystem \
				-p grpc-gateway \
				--author \
				--compare-link \
				--github-site=https://github.com \
				--unreleased-label "**Next release**" \
				--release-branch=v2 \
				--future-release=v2.0.0-beta.2

clean:
	rm -f $(GATEWAY_PLUGIN) $(OPENAPI_PLUGIN)
distclean: clean
	rm -f $(GO_PLUGIN)
realclean: distclean
	rm -f $(EXAMPLE_SVCSRCS) $(EXAMPLE_DEPSRCS)
	rm -f $(EXAMPLE_GWSRCS)
	rm -f $(EXAMPLE_OPENAPISRCS)
	rm -f $(EXAMPLE_CLIENT_SRCS)
	rm -f $(HELLOWORLD_SVCSRCS)
	rm -f $(HELLOWORLD_GWSRCS)
	rm -f $(OPENAPIV2_GO)
	rm -f $(RUNTIME_TEST_SRCS)

.PHONY: generate examples test lint clean distclean realclean
