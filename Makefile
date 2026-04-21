# This is a Makefile which maintains files automatically generated but to be
# shipped together with other files.
# You don't have to rebuild these targets by yourself unless you develop
# gRPC-Gateway itself.

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
GENERATE_UNBOUND_METHODS_EXAMPLE_SPEC=examples/internal/proto/examplepb/generate_unbound_methods.swagger.json
GENERATE_UNBOUND_METHODS_EXAMPLE_SRCS=$(EXAMPLE_CLIENT_DIR)/generateunboundmethods/client.go \
		 $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/response.go \
		 $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/configuration.go \
		 $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/model_examplepb_generate_unbound_methods_simple_message.go \
		 $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/api_generate_unbound_methods.go

EXAMPLE_CLIENT_SRCS=$(ECHO_EXAMPLE_SRCS) $(ABE_EXAMPLE_SRCS) $(UNANNOTATED_ECHO_EXAMPLE_SRCS) $(RESPONSE_BODY_EXAMPLE_SRCS) $(GENERATE_UNBOUND_METHODS_EXAMPLE_SRCS)
SWAGGER_CODEGEN=swagger-codegen

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
$(GENERATE_UNBOUND_METHODS_EXAMPLE_SRCS): $(GENERATE_UNBOUND_METHODS_EXAMPLE_SPEC)
	$(SWAGGER_CODEGEN) generate -i $(GENERATE_UNBOUND_METHODS_EXAMPLE_SPEC) \
	    -l go -o examples/internal/clients/generateunboundmethods --additional-properties packageName=generateunboundmethods
	@rm -f $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/README.md \
		$(EXAMPLE_CLIENT_DIR)/generateunboundmethods/git_push.sh

OAPI_CODEGEN_VERSION := v2.6.0

install:
	go install github.com/bufbuild/buf/cmd/buf@v1.45.0
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	go install \
		./protoc-gen-openapiv2 \
		./protoc-gen-openapiv3 \
		./protoc-gen-grpc-gateway

proto:
	# These generation steps are run in order so that later steps can
	# overwrite files produced by previous steps, if necessary.
	buf generate
	# Remove generated gateway in runtime tests, causes import cycle
	rm ./runtime/internal/examplepb/non_standard_names.pb.gw.go
	# Remove generated_input.proto files, bazel genrule relies on these
	# *not* being generated (to avoid conflicts).
	rm ./examples/internal/proto/examplepb/generated_input.pb.go
	rm ./examples/internal/proto/examplepb/generated_input_grpc.pb.go
	rm ./examples/internal/proto/examplepb/generated_input.pb.gw.go
	# Remove swagger files for openapiv2 definitions, they're unused
	rm ./protoc-gen-openapiv2/options/annotations.swagger.json
	rm ./protoc-gen-openapiv2/options/openapiv2.swagger.json
	buf generate \
		--template ./examples/internal/proto/examplepb/openapi_merge.buf.gen.yaml \
		--path ./examples/internal/proto/examplepb/openapi_merge_a.proto \
		--path ./examples/internal/proto/examplepb/openapi_merge_b.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/standalone_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/unannotated_echo_service.proto
	mv examples/internal/proto/examplepb/unannotated_echo_service.pb.gw.go examples/internal/proto/standalone/
	buf generate \
		--template ./examples/internal/proto/examplepb/unannotated_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/unannotated_echo_service.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/generate_unbound_methods.buf.gen.yaml \
		--path examples/internal/proto/examplepb/generate_unbound_methods.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/use_go_template.buf.gen.yaml \
		--path examples/internal/proto/examplepb/use_go_template.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/ignore_comment.buf.gen.yaml \
		--path examples/internal/proto/examplepb/ignore_comment.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/remove_internal_comment.buf.gen.yaml \
		--path examples/internal/proto/examplepb/remove_internal_comment.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/visibility_rule_preview_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/visibility_rule_echo_service.proto
	mv examples/internal/proto/examplepb/visibility_rule_echo_service.swagger.json examples/internal/proto/examplepb/visibility_rule_preview_echo_service.swagger.json
	buf generate \
		--template ./examples/internal/proto/examplepb/visibility_rule_internal_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/visibility_rule_echo_service.proto
	mv examples/internal/proto/examplepb/visibility_rule_echo_service.swagger.json examples/internal/proto/examplepb/visibility_rule_internal_echo_service.swagger.json
	buf generate \
		--template ./examples/internal/proto/examplepb/visibility_rule_none_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/visibility_rule_echo_service.proto
	mv examples/internal/proto/examplepb/visibility_rule_echo_service.swagger.json examples/internal/proto/examplepb/visibility_rule_none_echo_service.swagger.json
	buf generate \
		--template ./examples/internal/proto/examplepb/visibility_rule_preview_and_internal_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/visibility_rule_echo_service.proto
	mv examples/internal/proto/examplepb/visibility_rule_echo_service.swagger.json examples/internal/proto/examplepb/visibility_rule_preview_and_internal_echo_service.swagger.json
	buf generate \
		--template ./examples/internal/proto/examplepb/visibility_rule_enums_as_ints_echo_service.buf.gen.yaml \
		--path examples/internal/proto/examplepb/visibility_rule_echo_service.proto
	mv examples/internal/proto/examplepb/visibility_rule_echo_service.swagger.json examples/internal/proto/examplepb/visibility_rule_enums_as_ints_echo_service.swagger.json
	buf generate \
		--template examples/internal/proto/examplepb/enum_with_single_value.buf.gen.yaml \
		--path examples/internal/proto/examplepb/enum_with_single_value.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/proto3_field_semantics.buf.gen.yaml \
		--path examples/internal/proto/examplepb/proto3_field_semantics.proto
	buf generate \
		--template ./protoc-gen-openapiv2/options/buf.gen.yaml \
		--path ./protoc-gen-openapiv2/options/annotations.proto \
		--path ./protoc-gen-openapiv2/options/openapiv2.proto
	buf generate \
		--template ./examples/internal/proto/examplepb/opaque.buf.gen.yaml \
		--path examples/internal/proto/examplepb/opaque.proto
	buf generate \
		--template ./protoc-gen-openapiv3/buf.gen.yaml \
		--path ./examples/internal/helloworld/helloworld.proto \
		--path ./examples/internal/proto/examplepb/a_bit_of_everything.proto

# openapiv3-clients regenerates the Go clients consumed by the end-to-end
# oracle tests under examples/internal/integration/openapiv3. Both clients
# are produced by oapi-codegen from the specs emitted by `make proto`.
#
# The helloworldv3 oracle test (openapiv3_test.go) imports the generated
# client, stands up a real grpc-gateway in front of an in-process Greeter
# gRPC server, and round-trips a request through it. abe_oracle_test.go
# does the same for the ABE server impl under examples/internal/server. If
# the specs drift from what the gateway actually accepts, these tests fail.
#
# A third test in the same package (abe_spec_test.go) loads the checked-in
# ABE spec and asserts structural facts without any codegen — that one runs
# in normal CI and is the bulk of the coverage.
HELLOWORLD_V3_SPEC := examples/internal/helloworld/helloworld.openapi.json
HELLOWORLD_V3_CLIENT_DIR := examples/internal/clients/helloworldv3
ABE_V3_SPEC := examples/internal/proto/examplepb/a_bit_of_everything.openapi.json
ABE_V3_CLIENT_DIR := examples/internal/clients/abev3

openapiv3-clients: proto
	@rm -rf $(HELLOWORLD_V3_CLIENT_DIR)
	@mkdir -p $(HELLOWORLD_V3_CLIENT_DIR)
	oapi-codegen -package helloworldv3 -generate types,client \
		-o $(HELLOWORLD_V3_CLIENT_DIR)/helloworldv3.go $(HELLOWORLD_V3_SPEC)
	@rm -rf $(ABE_V3_CLIENT_DIR)
	@mkdir -p $(ABE_V3_CLIENT_DIR)
	oapi-codegen -package abev3 -generate types,client \
		-o $(ABE_V3_CLIENT_DIR)/abev3.go $(ABE_V3_SPEC)

generate: proto $(ECHO_EXAMPLE_SRCS) $(ABE_EXAMPLE_SRCS) $(UNANNOTATED_ECHO_EXAMPLE_SRCS) $(RESPONSE_BODY_EXAMPLE_SRCS) $(GENERATE_UNBOUND_METHODS_EXAMPLE_SRCS) openapiv3-clients

test: proto
	go test -short -race ./...
	go test -race ./examples/internal/integration -args -network=unix -endpoint=test.sock

clean:
	find . -type f -name '*.pb.go' -delete
	find . -type f -name '*.swagger.json' -delete
	find . -type f -name '*.pb.gw.go' -delete
	rm -f $(EXAMPLE_CLIENT_SRCS)

.PHONY: generate test clean proto install openapiv3-clients
