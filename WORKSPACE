workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_features",
    sha256 = "966c211ec42c4deb2af4c6dd6948408100b752f61753c97055bdac9bfb5cc0c7",
    strip_prefix = "bazel_features-1.41.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.41.0/bazel_features-v1.41.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

http_archive(
    name = "rules_python",
    sha256 = "94a2b4c5d9c45323a9737f8de8f841923bb628cace1e8e51fec5525ed9ccfb2d",
    strip_prefix = "rules_python-1.8.3",
    url = "https://github.com/bazelbuild/rules_python/releases/download/1.8.3/rules_python-1.8.3.tar.gz",
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
    sha256 = "440848dffa209beb8a04e41cc352762e44f8e91342b2a43aab6af9b30713c2f6",
    strip_prefix = "protobuf-33.5",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v33.5.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "aecc0465d7b6d17dc86efd1817a13cbed2b51ab8fe441bee10de30367367deac",
    strip_prefix = "googleapis-85b3f459f91cd08669c27da47d07bb57f6d9aa4d",
    urls = [
        "https://github.com/googleapis/googleapis/archive/85b3f459f91cd08669c27da47d07bb57f6d9aa4d.zip",
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
    sha256 = "68af54cb97fbdee5e5e8fe8d210d15a518f9d62abfd71620c3eaff3b26a5ff86",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.59.0/rules_go-v0.59.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.59.0/rules_go-v0.59.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.24.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "675114d8b433d0a9f54d81171833be96ebc4113115664b791e6f204d58e93446",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.47.0/bazel-gazelle-v0.47.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.47.0/bazel-gazelle-v0.47.0.tar.gz",
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
    sha256 = "e6b87c89bd0b27039e3af2c5da01147452f240f75d505f5b6880874f31036307",
    strip_prefix = "rules_shell-0.6.1",
    url = "https://github.com/bazelbuild/rules_shell/releases/download/v0.6.1/rules_shell-v0.6.1.tar.gz",
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
