package(
    default_visibility = ["//:__subpackages__"],
    default_testonly = 1,
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["message.pb.go"],
    deps = [
        "@com_github_golang_protobuf//:proto",
    ],
)


