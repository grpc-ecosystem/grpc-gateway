package(default_visibility = ["//visibility:public"])
licenses(["notice"])  # 3-Clause BSD

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library", "go_binary")

go_prefix("github.com/golang/protobuf")

go_library(
    name = "proto",
    srcs = glob(
        include = ["proto/*.go"],
        exclude = [
            "proto/*_test.go",
            "proto/pointer_reflect.go",
        ],
    ),
)

go_binary(
    name = "go_plugin",
    srcs = [
      "protoc-gen-go/main.go",
      "protoc-gen-go/link_grpc.go",
    ],
    deps = [
        ":proto",
        ":protoc-gen-go/generator",
        ":protoc-gen-go/grpc",
    ],
)

go_library(
    name = "protoc-gen-go/grpc",
    srcs = ["protoc-gen-go/grpc/grpc.go"],
    deps = [
        "protoc-gen-go/descriptor",
        "protoc-gen-go/generator",
    ],
)

go_library(
    name = "protoc-gen-go/generator",
    srcs = ["protoc-gen-go/generator/generator.go"],
    deps = [
        "protoc-gen-go/descriptor",
        "protoc-gen-go/plugin",
        "proto",
    ],
)

go_library(
    name = "protoc-gen-go/descriptor",
    srcs = ["protoc-gen-go/descriptor/descriptor.pb.go"],
    deps = [":proto"]
)

go_library(
    name = "protoc-gen-go/plugin",
    srcs = ["protoc-gen-go/plugin/plugin.pb.go"],
    deps = [
        ":proto",
        ":protoc-gen-go/descriptor",
    ]
)

go_library(
    name = "jsonpb",
    srcs = glob(
        include = ["jsonpb/*.go"],
        exclude = ["jsonpb/*_test.go"],
    ),
    deps = [":proto"],
)

go_library(
    name = "ptypes",
    srcs = glob(
        include = ["ptypes/go"],
        exclude = ["ptypes/_test.go"],
    ),
    deps = [
        ":ptypes/any",
        ":ptypes/duration",
        ":ptypes/empty",
        ":ptypes/struct",
        ":ptypes/timestamp",
        ":ptypes/wrappers",
    ],
)

go_library(
    name = "ptypes/any",
    srcs = glob(
        include = ["ptypes/any/*.go"],
        exclude = ["ptypes/any/*_test.go"],
    ),
    deps = [":proto"],
)

go_library(
    name = "ptypes/duration",
    srcs = glob(
        include = ["ptypes/duration/*.go"],
        exclude = ["ptypes/duration/*_test.go"],
    ),
    deps = [":proto"],
)

go_library(
    name = "ptypes/empty",
    srcs = glob(
        include = ["ptypes/empty/*.go"],
        exclude = ["ptypes/empty/*_test.go"],
    ),
    deps = [":proto"],
)

go_library(
    name = "ptypes/struct",
    srcs = glob(
        include = ["ptypes/struct/*.go"],
        exclude = ["ptypes/struct/*_test.go"],
    ),
    deps = [":proto"],
)

go_library(
    name = "ptypes/timestamp",
    srcs = glob(
        include = ["ptypes/timestamp/*.go"],
        exclude = ["ptypes/timestamp/*_test.go"],
    ),
    deps = [":proto"],
)

go_library(
    name = "ptypes/wrappers",
    srcs = glob(
        include = ["ptypes/wrappers/*.go"],
        exclude = ["ptypes/wrappers/*_test.go"],
    ),
    deps = [":proto"],
)
