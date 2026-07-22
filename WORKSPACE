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
    sha256 = "e11d2e1efce1589e5bdfa93986712c74fc7467a0f093143d489d2ef5ebb1ed2a",
    strip_prefix = "rules_python-2.2.0",
    url = "https://github.com/bazelbuild/rules_python/releases/download/2.2.0/rules_python-2.2.0.tar.gz",
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
    sha256 = "22775f9376938295efa2d59a59bde4cd075a42df5a9b4d27aa9b99fa6a413bd2",
    strip_prefix = "protobuf-35.1",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v35.1.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "e4228c8984fd02afc7491a3dd125c40a7584f12866165516585a9ba66b810b85",
    strip_prefix = "googleapis-a15a5b7adf695d1b6b17359ff64bfb4a83d0750b",
    urls = [
        "https://github.com/googleapis/googleapis/archive/a15a5b7adf695d1b6b17359ff64bfb4a83d0750b.zip",
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
    sha256 = "763f4a3f6b03469fdb00a77a333dd0b5546d3ee1fa29db373128c08fee73e0e8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.61.1/rules_go-v0.61.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.61.1/rules_go-v0.61.1.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.26.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "564508c5329de6fbffeade560cb5c02a5b9b86f47eb13856f0f00b8548d8e5d6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.52.0/bazel-gazelle-v0.52.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.52.0/bazel-gazelle-v0.52.0.tar.gz",
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
    sha256 = "20721f63908879c083f94869e618ea8d4ff5edb91ff9a72a2ebee357fdbc352d",
    strip_prefix = "rules_shell-0.8.0",
    url = "https://github.com/bazelbuild/rules_shell/releases/download/v0.8.0/rules_shell-v0.8.0.tar.gz",
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
