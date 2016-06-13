package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["generator.go"],
    deps = [
        "//protoc-gen-grpc-gateway/descriptor:go_default_library",
        "@com_github_golang_protobuf//:protoc-gen-go/descriptor",
        "@com_github_golang_protobuf//:protoc-gen-go/plugin",
    ],
)
