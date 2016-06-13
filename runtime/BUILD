package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "context.go",
        "convert.go",
        "doc.go",
        "errors.go",
        "handler.go",
        "marshal_json.go",
        "marshal_jsonpb.go",
        "marshaler.go",
        "marshaler_registry.go",
        "mux.go",
        "pattern.go",
        "proto2_convert.go",
        "query.go",
        "stream_chunk.pb.go",
    ],
    deps = [
        "//utilities:go_default_library",
        "@com_github_golang_protobuf//:jsonpb",
        "@com_github_golang_protobuf//:proto",
        "@org_golang_google_grpc//:codes",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//:grpclog",
        "@org_golang_google_grpc//:metadata",
        "@org_golang_x_net//:context",
    ],
)

go_test(
    name = "marshaler_registry_test",
    srcs = ["marshaler_registry_test.go"],
    deps = [
        ":go_default_library",
    ],
)

go_test(
    name = "marshal_json_test",
    srcs = ["marshal_json_test.go"],
    deps = [
        ":go_default_library",
        "//examples/examplepb:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:ptypes/empty",
        "@com_github_golang_protobuf//:ptypes/struct",
        "@com_github_golang_protobuf//:ptypes/timestamp",
        "@com_github_golang_protobuf//:ptypes/wrappers",
    ],
)

go_test(
    name = "marshal_jsonpb_test",
    srcs = ["marshal_jsonpb_test.go"],
    deps = [
        ":go_default_library",
        "//examples/examplepb:go_default_library",
        "@com_github_golang_protobuf//:jsonpb",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:ptypes/duration",
        "@com_github_golang_protobuf//:ptypes/empty",
        "@com_github_golang_protobuf//:ptypes/struct",
        "@com_github_golang_protobuf//:ptypes/timestamp",
        "@com_github_golang_protobuf//:ptypes/wrappers",
    ],
)

go_test(
    name = "context_test",
    srcs = ["context_test.go"],
    deps = [
        ":go_default_library",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:metadata",
    ],
)

go_test(
    name = "errors_test",
    srcs = ["errors_test.go"],
    deps = [
        ":go_default_library",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//:codes",
    ],
)

go_test(
    name = "mux_test",
    srcs = ["mux_test.go"],
    deps = [
        ":go_default_library",
        "//utilities:go_default_library",
    ],
)

go_test(
    name = "pattern_test",
    srcs = ["pattern_test.go"],
    library = ":go_default_library",
    deps = [
        "//utilities:go_default_library",
        "@com_github_golang_protobuf//:jsonpb",
        "@com_github_golang_protobuf//:proto",
        "@org_golang_google_grpc//:codes",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//:grpclog",
        "@org_golang_google_grpc//:metadata",
        "@org_golang_x_net//:context",
    ],
)

go_test(
    name = "query_test",
    srcs = ["query_test.go"],
    deps = [
        ":go_default_library",
        "//utilities:go_default_library",
        "@com_github_golang_protobuf//:proto",
    ],
)
