package(default_visibility = ["//visibility:public"])
licenses(["notice"])  # Apache 2.0

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "google/api",
    srcs = [
        "google/api/http.pb.go",
        "google/api/annotations.pb.go",
    ],
    deps = [
        "@com_github_golang_protobuf//:proto",
        "@com_github_golang_protobuf//:protoc-gen-go/descriptor",
    ],
)

