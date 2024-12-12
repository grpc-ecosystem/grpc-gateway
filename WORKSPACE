workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_features",
    sha256 = "9c7f9a4c997cbf0eb93572b68b0ff0bb9822004633101a16d89499e814190323",
    strip_prefix = "bazel_features-1.22.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.22.0/bazel_features-v1.22.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

http_archive(
    name = "rules_python",
    sha256 = "4f7e2aa1eb9aa722d96498f5ef514f426c1f55161c3c9ae628c857a7128ceb07",
    strip_prefix = "rules_python-1.0.0",
    url = "https://github.com/bazelbuild/rules_python/releases/download/1.0.0/rules_python-1.0.0.tar.gz",
)

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

http_archive(
    name = "com_google_googletest",
    sha256 = "f179ec217f9b3b3f3c6e8b02d3e7eda997b49e4ce26d6b235c9053bec9c0bf9f",
    strip_prefix = "googletest-1.15.2",
    urls = ["https://github.com/google/googletest/archive/v1.15.2.zip"],
)

# Define before rules_proto, otherwise we receive the version of com_google_protobuf from there
http_archive(
    name = "com_google_protobuf",
    sha256 = "3d32940e975c4ad9b8ba69640e78f5527075bae33ca2890275bf26b853c0962c",
    strip_prefix = "protobuf-29.1",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v29.1.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "bf51bf7278711ab666f23c211910054b81ff81e4a9096b7c6c2b422241aa6e7b",
    strip_prefix = "googleapis-f17e0aa970506b47d0552574ea7f57a440aa71f2",
    urls = [
        "https://github.com/googleapis/googleapis/archive/f17e0aa970506b47d0552574ea7f57a440aa71f2.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "bc283cdfcd526a52c3201279cda4bc298652efa898b10b4db0837dc51652756f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.7.1/bazel-skylib-1.7.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.7.1/bazel-skylib-1.7.1.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "rules_proto",
    sha256 = "0e5c64a2599a6e26c6a03d6162242d231ecc0de219534c38cb4402171def21e8",
    strip_prefix = "rules_proto-7.0.2",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/7.0.2.tar.gz",
    ],
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

load("@rules_proto//proto:toolchains.bzl", "rules_proto_toolchains")

rules_proto_toolchains()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f4a9314518ca6acfa16cc4ab43b0b8ce1e4ea64b81c38d8a3772883f153346b8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.50.1/rules_go-v0.50.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.50.1/rules_go-v0.50.1.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.22.7")

http_archive(
    name = "bazel_gazelle",
    sha256 = "a80893292ae1d78eaeedd50d1cab98f242a17e3d5741b1b9fb58b5fd9d2d57bc",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.40.0/bazel-gazelle-v0.40.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.40.0/bazel-gazelle-v0.40.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

# Use gazelle to declare Go dependencies in Bazel.
# gazelle:repository_macro repositories.bzl%go_repositories

load("//:repositories.bzl", "go_repositories")

go_repositories()

# This must be invoked after our explicit dependencies
# See https://github.com/bazelbuild/bazel-gazelle/issues/1115.
gazelle_dependencies()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "051951c10ff8addeb4f10be3b0cf474b304b2ccd675f2cc7683cdd9010320ca9",
    strip_prefix = "buildtools-7.3.1",
    urls = ["https://github.com/bazelbuild/buildtools/archive/v7.3.1.tar.gz"],
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
