package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "registry.go",
        "services.go",
        "types.go",
    ],
    deps = [
        "//protoc-gen-grpc-gateway/httprule:go_default_library",
        "//third_party/googleapis:google/api",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:protoc-gen-go/descriptor",
        "@com_github_golang_protobuf//:protoc-gen-go/generator",
        "@com_github_golang_protobuf//:protoc-gen-go/plugin",
    ],
)

go_test(
    name = "library_test",
    srcs = [
        "registry_test.go",
        "services_test.go",
        "types_test.go",
    ],
    library = ":go_default_library",
    deps = [
        "//protoc-gen-grpc-gateway/httprule:go_default_library",
        "//third_party/googleapis:google/api",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:protoc-gen-go/descriptor",
        "@com_github_golang_protobuf//:protoc-gen-go/generator",
        "@com_github_golang_protobuf//:protoc-gen-go/plugin",
    ],
)
