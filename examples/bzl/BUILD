load(
    "//:grpc-gateway.bzl",
    "go_proto_library",
    "grpc_gateway_proto_library",
    "grpc_swagger_proto_library"
)
load("@com_github_google_protobuf//:protobuf.bzl", "cc_proto_library")

PKGMAP = {
    "google/protobuf/descriptor.proto": "github.com/golang/protobuf/protoc-gen-go/descriptor",
    "google/api/annotations.proto": "github.com/gengo/grpc-gateway/examples/bzl/google/api",
    "examples/sub/message.proto": "github.com/gengo/grpc-gateway/examples/bzl/examples/sub",
}

# Workaround of the conflict issue of generated files in bazel.
# c.f. https://github.com/bazelbuild/rules_go/issues/16#issuecomment-216233626
#
# TODO(yugui) Remove this workaround and move the example rules below into individual
# packages once the issue is resolved.
genrule(
    name = "collect_sources",
    cmd = " && ".join([
        "rm -rf $(@D)/google $(@D)/examples",
        "mkdir -p $(@D)/google/api",
        "mkdir -p $(@D)/examples/examplepb",
        "mkdir -p $(@D)/examples/sub",
        "mkdir -p $(@D)/examples/sub2",
        "basedir=`pwd`",
        "ln -s $$basedir/examples/examplepb/echo_service.proto $(@D)/examples/examplepb/",
        "ln -s $$basedir/examples/examplepb/a_bit_of_everything.proto $(@D)/examples/examplepb/",
        "ln -s $$basedir/examples/examplepb/flow_combination.proto $(@D)/examples/examplepb/",
        "ln -s $$basedir/examples/sub/message.proto $(@D)/examples/sub/",
        "ln -s $$basedir/examples/sub2/message.proto $(@D)/examples/sub2/",
        "ln -s $$basedir/third_party/googleapis/google/api/http.proto $(@D)/google/api/",
        "ln -s $$basedir/third_party/googleapis/google/api/annotations.proto $(@D)/google/api/",
    ]),
    srcs = [
        "//examples/examplepb:echo_service.proto",
        "//examples/examplepb:a_bit_of_everything.proto",
        "//examples/examplepb:flow_combination.proto",
        "//examples/sub:message.proto",
        "//examples/sub2:message.proto",
        "//third_party/googleapis:google/api/http.proto",
        "//third_party/googleapis:google/api/annotations.proto",
    ],
    outs = [
        "examples/examplepb/echo_service.proto",
        "examples/examplepb/a_bit_of_everything.proto",
        "examples/examplepb/flow_combination.proto",
        "examples/sub/message.proto",
        "examples/sub2/message.proto",
        "google/api/http.proto",
        "google/api/annotations.proto",
    ],
)

grpc_swagger_proto_library(
    name = "examples/examplepb_swagger",
    srcs = [
        "examples/examplepb/echo_service.proto",
        "examples/examplepb/a_bit_of_everything.proto",
        "examples/examplepb/flow_combination.proto",
    ],
    deps = [
        ":google/api_genproto",
    ],
    includes = [
        ".",
        "examples/bzl",
    ],
)

grpc_gateway_proto_library(
    name = "examples/examplepb",
    srcs = [
        "examples/examplepb/echo_service.proto",
        "examples/examplepb/a_bit_of_everything.proto",
        "examples/examplepb/flow_combination.proto",
    ],
    deps = [
        ":google/api_genproto",
    ],
    pkgmap = PKGMAP,
    includes = [
        ".",
        "examples/bzl",
    ],
    go_deps = [
        ":google/api",
        ":examples/sub",
        # It is hard to support examples/sub2 in generic go_proto_library
        # because we cannot predict its output file name without reading
        # content of the source file.
        # c.f. https://github.com/golang/protobuf/issues/139
        #
        # TODO(yugui) Use :examples/sub2 instead once we get a good solution
        # for the unpredictability.
        "//examples/sub2:go_default_library",
    ],
    go_extra_library = ":examples/examplepb_proto",
)

go_proto_library(
    name = "examples/examplepb_proto",
    srcs = [
        "examples/examplepb/echo_service.proto",
        "examples/examplepb/a_bit_of_everything.proto",
        "examples/examplepb/flow_combination.proto",
    ],
    deps = [
        ":google/api_genproto",
    ],
    genopts = ["plugins=grpc"],
    pkgmap = PKGMAP,
    includes = [
        ".",
        "examples/bzl",
    ],
    go_deps = [
        ":google/api",
        ":examples/sub",
        # Same as above.
        # TODO(yugui) Use :examples/sub2 instead.
        "//examples/sub2:go_default_library",
        "@org_golang_x_net//:context",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_proto_library(
    name = "examples/sub",
    srcs = [
        "examples/sub/message.proto",
    ],
    pkgmap = PKGMAP,
    includes = [
        ".",
        "examples/bzl",
    ],
)

go_proto_library(
    name = "google/api",
    srcs = [
        "google/api/http.proto",
        "google/api/annotations.proto",
    ],
    deps = [
        "@com_github_google_protobuf//:cc_wkt_protos_genproto",
    ],
    pkgmap = PKGMAP,
    includes = [
        ".",
        "examples/bzl",
    ],
    go_deps = [
        "@com_github_golang_protobuf//:protoc-gen-go/descriptor",
    ],
)
