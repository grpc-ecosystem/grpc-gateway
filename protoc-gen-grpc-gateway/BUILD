package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "protoc-gen-grpc-gateway",
    srcs = ["main.go"],
    deps = [
        "//protoc-gen-grpc-gateway/descriptor:go_default_library",
        "//protoc-gen-grpc-gateway/gengateway:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:protoc-gen-go/plugin",
    ],
)
