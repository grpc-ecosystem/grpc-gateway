package(
    default_visibility = ["//examples:__pkg__"],
    default_testonly = 1,
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "a_bit_of_everything.go",
        "echo.go",
        "flow_combination.go",
        "main.go",
    ],
    deps = [
        "//examples/examplepb:go_default_library",
        "//examples/sub:go_default_library",
        "//examples/sub2:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_rogpeppe_fastuuid//:go_default_library",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:codes",
        "@org_golang_google_grpc//:metadata",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
