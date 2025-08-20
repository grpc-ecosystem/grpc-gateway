workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_features",
    sha256 = "f08776c1430e8f35209a054f828d0985019879e554c4eb32a093e1f49bc9231e",
    strip_prefix = "bazel_features-1.34.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.34.0/bazel_features-v1.34.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

http_archive(
    name = "rules_python",
    sha256 = "0a1cefefb4a7b550fb0b43f54df67d6da95b7ba352637669e46c987f69986f6a",
    strip_prefix = "rules_python-1.5.3",
    url = "https://github.com/bazelbuild/rules_python/releases/download/1.5.3/rules_python-1.5.3.tar.gz",
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
    sha256 = "3ad017543e502ffaa9cd1f4bd4fe96cf117ce7175970f191705fa0518aff80cd",
    strip_prefix = "protobuf-32.0",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v32.0.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "8b9798b03011ab58bce025c43b7530782b720eb6bbde9bbdb43c6495d7e223a1",
    strip_prefix = "googleapis-3b2a2ae91db23a9c879b2b725d6a5de6bd64a800",
    urls = [
        "https://github.com/googleapis/googleapis/archive/3b2a2ae91db23a9c879b2b725d6a5de6bd64a800.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "51b5105a760b353773f904d2bbc5e664d0987fbaf22265164de65d43e910d8ac",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.8.1/bazel-skylib-1.8.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.8.1/bazel-skylib-1.8.1.tar.gz",
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
    sha256 = "89d2050410602142c9acafd01c95baf48b65f8dd16f4771d37c89f82f5e147f2",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.56.1/rules_go-v0.56.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.56.1/rules_go-v0.56.1.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.24.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "e467b801046b6598c657309b45d2426dc03513777bd1092af2c62eebf990aca5",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.45.0/bazel-gazelle-v0.45.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.45.0/bazel-gazelle-v0.45.0.tar.gz",
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
    sha256 = "fce2a7a974aa70e9367068122e19c39a6a27a5aca30698bcf9030beb529612b6",
    strip_prefix = "rules_shell-0.6.0",
    url = "https://github.com/bazelbuild/rules_shell/releases/download/v0.6.0/rules_shell-v0.6.0.tar.gz",
)

load("@rules_shell//shell:repositories.bzl", "rules_shell_dependencies", "rules_shell_toolchains")

rules_shell_dependencies()

rules_shell_toolchains()

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "53119397bbce1cd7e4c590e117dcda343c2086199de62932106c80733526c261",
    strip_prefix = "buildtools-8.2.1",
    urls = ["https://github.com/bazelbuild/buildtools/archive/v8.2.1.tar.gz"],
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
