proto_library(
    visibility = ["//visibility:public"],
    name = "annotations_proto",
    srcs = ["annotations.proto"],
    strip_import_prefix="/third_party/googleapis",
    deps = [
      ":http_proto",
      "@com_google_protobuf//:descriptor_proto",
    ],
)

proto_library(
    visibility = ["//visibility:public"],
    name = "http_proto",
    strip_import_prefix="/third_party/googleapis",
    srcs = ["http.proto"],
)

proto_library(
    visibility = ["//visibility:public"],
    name = "httpbody_proto",
    strip_import_prefix="/third_party/googleapis",
    srcs = ["httpbody.proto"],
    deps = [
      "@com_google_protobuf//:any_proto"
    ],
)
