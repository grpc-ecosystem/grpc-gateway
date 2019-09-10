proto_library(
    visibility = ["//visibility:public"],
    name = "code_proto",
    srcs = ["code.proto"],
    strip_import_prefix="/third_party/googleapis",
)

proto_library(
    visibility = ["//visibility:public"],
    name = "error_details_proto",
    srcs = ["error_details.proto"],
    strip_import_prefix="/third_party/googleapis",
    deps = [
      "@com_google_protobuf//:duration_proto",
    ],
)

proto_library(
    visibility = ["//visibility:public"],
    name = "status_proto",
    srcs = ["status.proto"],
    strip_import_prefix="/third_party/googleapis",
    deps = [
      "@com_google_protobuf//:any_proto",
    ],
)
