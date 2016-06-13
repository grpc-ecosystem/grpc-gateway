package(default_visibility = ["//visibility:public"])
licenses(["notice"])  # Apache 2.0

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library")

go_prefix("github.com/golang/glog")

go_library(
    name = "go_default_library",
    srcs = [
        "glog.go",
        "glog_file.go",
    ],
)
