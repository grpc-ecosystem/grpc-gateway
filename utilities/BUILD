package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "pattern.go",
        "trie.go",
    ],
    deps = [
        "@org_golang_google_grpc//:metadata",
        "@org_golang_x_net//:context",
    ],
)

go_test(
    name = "trie_test",
    srcs = ["trie_test.go"],
    deps = [
        ":go_default_library",
    ],
)
