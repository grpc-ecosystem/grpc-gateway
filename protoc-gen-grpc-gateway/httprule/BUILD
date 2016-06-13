package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "compile.go",
        "parse.go",
        "types.go",
    ],
    deps = [
        "//utilities:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "compile_test",
    srcs = ["compile_test.go"],
    library = ":go_default_library",
    deps = [
        "//utilities:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "parse_test",
    srcs = ["parse_test.go"],
    library = ":go_default_library",
    deps = [
        "//utilities:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "types_test",
    srcs = ["types_test.go"],
    library = ":go_default_library",
    deps = [
        "//utilities:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)
