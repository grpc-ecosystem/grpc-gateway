package(default_visibility = ["//visibility:public"])
licenses(["notice"])  # 3-Clause BSD

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library")

go_prefix("google.golang.org/grpc")

go_library(
    name = "go_default_library",
    srcs = glob(
        include = ["*.go"],
        exclude = ["*_test.go"],
    ),
    deps = [
        ":codes",
        ":credentials",
        ":grpclog",
        ":internal",
        ":metadata",
        ":naming",
        ":transport",
        "@com_github_golang_protobuf//:proto",
        "@org_golang_x_net//:context",
        "@org_golang_x_net//:http2",
        "@org_golang_x_net//:trace",
    ],
)

go_library(
    name = "metadata",
    srcs = ["metadata/metadata.go"],
    deps = [
        "@org_golang_x_net//:context",
    ],
)

go_library(
    name = "codes",
    srcs = glob(
        include = ["codes/*.go"],
        exclude = ["codes/*_test.go"],
    ),
    deps = [
        "@org_golang_x_net//:context",
    ],
)

go_library(
    name = "transport",
    srcs = glob(
        include = ["transport/*.go"],
        exclude = ["transport/*_test.go"],
    ),
    deps = [
        ":codes",
        ":credentials",
        ":grpclog",
        ":metadata",
        ":peer",
        "@org_golang_x_net//:context",
        "@org_golang_x_net//:http2",
        "@org_golang_x_net//:http2/hpack",
        "@org_golang_x_net//:trace",
    ],
)

go_library(
    name = "credentials",
    srcs = glob(
        include = ["credentials/*.go"],
        exclude = ["credentials/*_test.go"],
    ),
    deps = [
        "@org_golang_x_net//:context",
    ],
)

go_library(
    name = "peer",
    srcs = glob(
        include = ["peer/*.go"],
        exclude = ["peer/*_test.go"],
    ),
    deps = [
        ":credentials",
        "@org_golang_x_net//:context",
    ],
)

go_library(
    name = "grpclog",
    srcs = glob(
        include = ["grpclog/*.go"],
        exclude = ["grpclog/*_test.go"],
    ),
)

go_library(
    name = "naming",
    srcs = glob(
        include = ["naming/*.go"],
        exclude = ["naming/*_test.go"],
    ),
)

go_library(
    name = "internal",
    srcs = glob(
        include = ["internal/*.go"],
        exclude = ["internal/*_test.go"],
    ),
    visibility = ["//visibility:private"],
)
