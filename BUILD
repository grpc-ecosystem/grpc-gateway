load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:exclude third_party

gazelle(
    name = "gazelle",
    mode = "diff",
    prefix = "github.com/grpc-ecosystem/grpc-gateway",
)

package_group(
    name = "generators",
    packages = [
        "//protoc-gen-grpc-gateway/...",
        "//protoc-gen-swagger/...",
    ],
)
