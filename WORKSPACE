workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_features",
    sha256 = "5450bfb2c8b4bc961c75368838f86156f563cc9adef1be7d504fc5619d54daab",
    strip_prefix = "bazel_features-1.51.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.51.0/bazel_features-v1.51.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

http_archive(
    name = "rules_python",
    sha256 = "25a00d55e4376ef6d9f28938b68038faca17c95580c9afe6953e7dde013d00e4",
    strip_prefix = "rules_python-2.3.2",
    url = "https://github.com/bazelbuild/rules_python/releases/download/2.3.2/rules_python-2.3.2.tar.gz",
)

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

http_archive(
    name = "com_google_googletest",
    sha256 = "63b9c77751a5b8f492486005f67533fdc58682b67476fcb91650be5958d5195a",
    strip_prefix = "googletest-1.18.0",
    urls = ["https://github.com/google/googletest/archive/v1.18.0.zip"],
)

# Define before rules_proto, otherwise we receive the version of com_google_protobuf from there
http_archive(
    name = "com_google_protobuf",
    sha256 = "d120933b88a2857c1b5fb6dafe387cca1b96762c2263cb53d945dcae34e88a3f",
    strip_prefix = "protobuf-36.0",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v36.0.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "a2451deaf575865e27071f2f95b49d7a7716714f02fa7130819991a82bc2c0f0",
    strip_prefix = "googleapis-240b58fe7058f6fff77cba02f51d600170b6c421",
    urls = [
        "https://github.com/googleapis/googleapis/archive/1f38da6aa7661cf22e17247ad33d2d566a2c356e.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "37cdfbc6faefea94f7b37760a305c98c08981116c2bc9e821e3b423221fad8c8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.9.2/bazel-skylib-1.9.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.9.2/bazel-skylib-1.9.2.tar.gz",
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
    sha256 = "c3e253237109ab2e2a8d3cb075688b98a6e6fce43d849849648e8e3a84f20d6f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.63.0/rules_go-v0.63.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.63.0/rules_go-v0.63.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.26.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "fec7a24acb027359eb7d2b31d116b17c81d2fab35842f6bfb2dd5945e1071346",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.53.0/bazel-gazelle-v0.53.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.53.0/bazel-gazelle-v0.53.0.tar.gz",
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
    sha256 = "1df78b782b8df6d34ee05bccea1f996f2b5d9cf961656905704d41bfb4c0ea7f",
    strip_prefix = "rules_shell-0.9.0",
    url = "https://github.com/bazelbuild/rules_shell/releases/download/v0.9.0/rules_shell-v0.9.0.tar.gz",
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
