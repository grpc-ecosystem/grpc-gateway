load("@bazel_gazelle//:def.bzl", "gazelle")
load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")

buildifier(
    name = "buildifier",
)

buildifier(
    name = "buildifier_check",
    mode = "check",
)

# gazelle:exclude third_party
# gazelle:exclude _output
# gazelle:prefix github.com/grpc-ecosystem/grpc-gateway/v2
# gazelle:go_proto_compilers //:go_apiv2
# gazelle:go_grpc_compilers //:go_apiv2, //:go_grpc

gazelle(name = "gazelle")

package_group(
    name = "generators",
    packages = [
        "//protoc-gen-grpc-gateway/...",
        "//protoc-gen-openapiv2/...",
    ],
)

go_proto_compiler(
    name = "go_apiv2",
    options = [
        "paths=source_relative",
    ],
    plugin = "@org_golang_google_protobuf//cmd/protoc-gen-go",
    suffix = ".pb.go",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:any_go_proto",
        "@io_bazel_rules_go//proto/wkt:api_go_proto",
        "@io_bazel_rules_go//proto/wkt:compiler_plugin_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
        "@io_bazel_rules_go//proto/wkt:field_mask_go_proto",
        "@io_bazel_rules_go//proto/wkt:source_context_go_proto",
        "@io_bazel_rules_go//proto/wkt:struct_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@io_bazel_rules_go//proto/wkt:type_go_proto",
        "@io_bazel_rules_go//proto/wkt:wrappers_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//runtime/protoimpl:go_default_library",
    ],
)

go_proto_compiler(
    name = "go_grpc",
    options = [
        "paths=source_relative",
        "require_unimplemented_servers=false",
    ],
    plugin = "@org_golang_google_grpc_cmd_protoc_gen_go_grpc//:protoc-gen-go-grpc",
    suffix = "_grpc.pb.go",
    visibility = ["//visibility:public"],
    deps = [
        "@io_bazel_rules_go//proto/wkt:any_go_proto",
        "@io_bazel_rules_go//proto/wkt:api_go_proto",
        "@io_bazel_rules_go//proto/wkt:compiler_plugin_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
        "@io_bazel_rules_go//proto/wkt:field_mask_go_proto",
        "@io_bazel_rules_go//proto/wkt:source_context_go_proto",
        "@io_bazel_rules_go//proto/wkt:struct_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@io_bazel_rules_go//proto/wkt:type_go_proto",
        "@io_bazel_rules_go//proto/wkt:wrappers_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
