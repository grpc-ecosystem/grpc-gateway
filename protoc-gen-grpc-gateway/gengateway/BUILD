package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "generator.go",
        "template.go",
    ],
    deps = [
        "//protoc-gen-grpc-gateway/descriptor:go_default_library",
        "//protoc-gen-grpc-gateway/generator:go_default_library",
        "//utilities:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:protoc-gen-go/plugin",
    ],
)

go_test(
    name = "template_test",
    srcs = [
        "template_test.go",
    ],
    library = ":go_default_library",
    deps = [
        "//protoc-gen-grpc-gateway/descriptor:go_default_library",
        "//protoc-gen-grpc-gateway/generator:go_default_library",
        "//protoc-gen-grpc-gateway/httprule:go_default_library",
        "//utilities:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:protoc-gen-go/descriptor",
        "@com_github_golang_protobuf//:protoc-gen-go/plugin",
    ],
)
