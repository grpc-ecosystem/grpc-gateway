package(
    default_visibility = ["//:__subpackages__"],
    default_testonly = 1,
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "a_bit_of_everything.pb.go",
        "a_bit_of_everything.pb.gw.go",
        "echo_service.pb.go",
        "echo_service.pb.gw.go",
        "flow_combination.pb.go",
        "flow_combination.pb.gw.go",
    ],
    deps = [
        "//examples/sub2:go_default_library",
        "//examples/sub:go_default_library",
        "//runtime:go_default_library",
        "//third_party/googleapis:google/api",
        "//utilities:go_default_library",
        "@com_github_golang_protobuf//:proto",
        "@org_golang_google_grpc//:codes",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//:grpclog",
        "@org_golang_x_net//:context",
    ],
)
