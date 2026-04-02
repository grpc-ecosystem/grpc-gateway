workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_features",
    sha256 = "ccf85bbf0613d12bf6df2c8470ecec544a6fe8ceab684e970e8ed4dde4cb24ec",
    strip_prefix = "bazel_features-1.44.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.44.0/bazel_features-v1.44.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

http_archive(
    name = "rules_python",
    sha256 = "098ba13578e796c00c853a2161f382647f32eb9a77099e1c88bc5299333d0d6e",
    strip_prefix = "rules_python-1.9.0",
    url = "https://github.com/bazelbuild/rules_python/releases/download/1.9.0/rules_python-1.9.0.tar.gz",
)

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

http_archive(
    name = "com_google_googletest",
    sha256 = "40d4ec942217dcc84a9ebe2a68584ada7d4a33a8ee958755763278ea1c5e18ff",
    strip_prefix = "googletest-1.17.0",
    urls = ["https://github.com/google/googletest/archive/v1.17.0.zip"],
)

# Define before rules_proto, otherwise we receive the version of com_google_protobuf from there
http_archive(
    name = "com_google_protobuf",
    sha256 = "a83103b7ed3afaeedee9a212c8f65825444f58144f5e075b73c83f2b4ff27b62",
    strip_prefix = "protobuf-34.1",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v34.1.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "4a5cd7834a53c1fc62252ec2f1f905b326c24fec7a2e28ba9d62bbaf87a9a26f",
    strip_prefix = "googleapis-8d52a0bd5332bec365647fd792102a25a665e9fe",
    urls = [
        "https://github.com/googleapis/googleapis/archive/8d52a0bd5332bec365647fd792102a25a665e9fe.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "3b5b49006181f5f8ff626ef8ddceaa95e9bb8ad294f7b5d7b11ea9f7ddaf8c59",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.9.0/bazel-skylib-1.9.0.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.9.0/bazel-skylib-1.9.0.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "rules_proto",
    sha256 = "14a225870ab4e91869652cfd69ef2028277fc1dc4910d65d353b62d6e0ae21f4",
    strip_prefix = "rules_proto-7.1.0",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/7.1.0.tar.gz",
    ],
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

load("@rules_proto//proto:toolchains.bzl", "rules_proto_toolchains")

rules_proto_toolchains()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "86d3dc8f59d253524f933aaf2f3c05896cb0b605fc35b460c0b4b039996124c6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.26.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "6549aff70998217292406776024d6da91b4e764c679d180ea072c557c70dacf2",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.48.0/bazel-gazelle-v0.48.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.48.0/bazel-gazelle-v0.48.0.tar.gz",
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
    name = "rules_shell",
    sha256 = "3709d1745ba4be4ef054449647b62e424267066eca887bb00dd29242cb8463a0",
    strip_prefix = "rules_shell-0.7.1",
    url = "https://github.com/bazelbuild/rules_shell/releases/download/v0.7.1/rules_shell-v0.7.1.tar.gz",
)

load("@rules_shell//shell:repositories.bzl", "rules_shell_dependencies", "rules_shell_toolchains")

rules_shell_dependencies()

rules_shell_toolchains()

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "f3b800e9f6ca60bdef3709440f393348f7c18a29f30814288a7326285c80aab9",
    strip_prefix = "buildtools-8.5.1",
    urls = ["https://github.com/bazelbuild/buildtools/archive/v8.5.1.tar.gz"],
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
