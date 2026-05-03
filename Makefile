# This is a Makefile which maintains files automatically generated but to be
# shipped together with other files.
# You don't have to rebuild these targets by yourself unless you develop
# gRPC-Gateway itself.

EXAMPLE_CLIENT_DIR=examples/internal/clients
ECHO_EXAMPLE_SPEC=examples/internal/proto/examplepb/echo_service.swagger.json
ABE_EXAMPLE_SPEC=examples/internal/proto/examplepb/a_bit_of_everything.swagger.json
UNANNOTATED_ECHO_EXAMPLE_SPEC=examples/internal/proto/examplepb/unannotated_echo_service.swagger.json
RESPONSE_BODY_EXAMPLE_SPEC=examples/internal/proto/examplepb/response_body_service.swagger.json
GENERATE_UNBOUND_METHODS_EXAMPLE_SPEC=examples/internal/proto/examplepb/generate_unbound_methods.swagger.json

# go-swagger generates a per-spec tree:
#   <client-dir>/client/                 facade (package "client")
#   <client-dir>/client/<service-name>/  one subpackage per service tag
#   <client-dir>/models/                 shared model types (package "models")
# Targets are PHONY because the file set isn't fully predictable from the
# spec without parsing it.
#
# The jq pass coerces non-body parameters with no `type`/`schema` (map<K,V>
# request fields) and non-body parameters typed as `object` (oneof-of-empty
# fields) to `string`. They can't round-trip through a query string in a
# meaningful way, but go-swagger needs *some* primitive to chew on. The
# accompanying default-value rewrite turns string defaults on numeric/integer
# parameters into the corresponding number — protoc-gen-openapiv2 emits e.g.
# `"default": "0"` on integer params, which go-swagger refuses.
define swagger_client
	rm -rf $(1)/client $(1)/models
	mkdir -p $(1)
	jq 'walk( \
		if type == "object" and has("in") and (.in != "body") then \
			(if (has("type") | not) and (has("schema") | not) then . + {type: "string"} \
			 elif .type == "object" then . + {type: "string"} \
			 else . end) \
			| (if (.type == "number" or .type == "integer") and (.default | type) == "string" \
			   then .default |= tonumber else . end) \
		else . end \
	)' $(2) > $(1)/.spec.json
	swagger generate client -q -f $(1)/.spec.json -t $(1) \
		--client-package=client --model-package=models --skip-validation
	rm -f $(1)/.spec.json
endef

echo-client:
	$(call swagger_client,$(EXAMPLE_CLIENT_DIR)/echo,$(ECHO_EXAMPLE_SPEC))
abe-client:
	$(call swagger_client,$(EXAMPLE_CLIENT_DIR)/abe,$(ABE_EXAMPLE_SPEC))
unannotatedecho-client:
	$(call swagger_client,$(EXAMPLE_CLIENT_DIR)/unannotatedecho,$(UNANNOTATED_ECHO_EXAMPLE_SPEC))
responsebody-client:
	$(call swagger_client,$(EXAMPLE_CLIENT_DIR)/responsebody,$(RESPONSE_BODY_EXAMPLE_SPEC))
generateunboundmethods-client:
	$(call swagger_client,$(EXAMPLE_CLIENT_DIR)/generateunboundmethods,$(GENERATE_UNBOUND_METHODS_EXAMPLE_SPEC))

swagger-clients: echo-client abe-client unannotatedecho-client responsebody-client generateunboundmethods-client

OAPI_CODEGEN_VERSION := v2.6.0
GO_SWAGGER_VERSION := v0.33.2

install:
	go install github.com/bufbuild/buf/cmd/buf@v1.45.0
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	go install github.com/go-swagger/go-swagger/cmd/swagger@$(GO_SWAGGER_VERSION)
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

generate: proto swagger-clients openapiv3-clients

test: proto
	go test -short -race ./...
	go test -race ./examples/internal/integration -args -network=unix -endpoint=test.sock

clean:
	find . -type f -name '*.pb.go' -delete
	find . -type f -name '*.swagger.json' -delete
	find . -type f -name '*.pb.gw.go' -delete
	rm -rf $(EXAMPLE_CLIENT_DIR)/echo/client $(EXAMPLE_CLIENT_DIR)/echo/models
	rm -rf $(EXAMPLE_CLIENT_DIR)/abe/client $(EXAMPLE_CLIENT_DIR)/abe/models
	rm -rf $(EXAMPLE_CLIENT_DIR)/unannotatedecho/client $(EXAMPLE_CLIENT_DIR)/unannotatedecho/models
	rm -rf $(EXAMPLE_CLIENT_DIR)/responsebody/client $(EXAMPLE_CLIENT_DIR)/responsebody/models
	rm -rf $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/client $(EXAMPLE_CLIENT_DIR)/generateunboundmethods/models

.PHONY: generate test clean proto install openapiv3-clients \
	swagger-clients echo-client abe-client unannotatedecho-client \
	responsebody-client generateunboundmethods-client
