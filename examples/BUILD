package(
    default_visibility = ["//visibility:private"],
    default_testonly = 1,
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    deps = [
        "//examples/examplepb:go_default_library",
        "//runtime:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "examples",
    library = ":go_default_library",
    deps = [
        "//examples/examplepb:go_default_library",
        "//runtime:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_test(
    name = "integration_test",
    srcs = ["integration_test.go"],
    library = ":go_default_library",
    deps = [
        "//examples/examplepb:go_default_library",
        "//examples/sub:go_default_library",
        "//examples/server:go_default_library",
        "//runtime:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//:jsonpb",
        "@com_github_golang_protobuf//:proto",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:codes",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
