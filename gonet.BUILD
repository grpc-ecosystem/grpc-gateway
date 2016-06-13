package(default_visibility = ["//visibility:public"])
licenses(["notice"])  # 3-Clause BSD

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library")

go_prefix("golang.org/x/net")

go_library(
    name = "context",
    srcs = [
        "context/context.go",
        "context/pre_go17.go",
    ]
)

go_library(
    name = "trace",
    srcs = glob(
        include = ["trace/*.go"],
        exclude = ["trace/*_test.go"],
    ),
    deps = [
        ":context",
        ":internal/timeseries",
    ],
)

go_library(
    name = "http2",
    srcs = glob(
        include = ["http2/*.go"],
        exclude = [
            "http2/*_test.go",
            "http2/go17.go",
            "http2/not_go16.go",
        ],
    ),
    deps = [
        ":http2/hpack",
    ],
)

go_library(
    name = "http2/hpack",
    srcs = glob(
        include = ["http2/hpack/*.go"],
        exclude = ["http2/hpack/*_test.go"],
    ),
    deps = [
    ],
)

go_library(
    name = "internal/timeseries",
    srcs = glob(
        include = ["internal/timeseries/*.go"],
        exclude = ["internal/timeseries/*_test.go"],
    ),
    visibility = ["//visibility:private"],
)
