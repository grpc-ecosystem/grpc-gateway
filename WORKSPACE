workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_features",
    sha256 = "091d8b1e1f0bf1f7bd688b95007687e862cc489f8d9bc21c14be5fd032a8362f",
    strip_prefix = "bazel_features-1.26.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.26.0/bazel_features-v1.26.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

http_archive(
    name = "rules_python",
    sha256 = "2ef40fdcd797e07f0b6abda446d1d84e2d9570d234fddf8fcd2aa262da852d1c",
    strip_prefix = "rules_python-1.2.0",
    url = "https://github.com/bazelbuild/rules_python/releases/download/1.2.0/rules_python-1.2.0.tar.gz",
)

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

http_archive(
    name = "com_google_googletest",
    sha256 = "a9607c9215866bd425a725610c5e0f739eeb50887a57903df48891446ce6fb3c",
    strip_prefix = "googletest-1.16.0",
    urls = ["https://github.com/google/googletest/archive/v1.16.0.zip"],
)

# Define before rules_proto, otherwise we receive the version of com_google_protobuf from there
http_archive(
    name = "com_google_protobuf",
    sha256 = "1451b03faec83aed17cdc71671d1bbdfd72e54086b827f5f6fd02bf7a4041b68",
    strip_prefix = "protobuf-30.1",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v30.1.tar.gz"],
)

http_archive(
    name = "googleapis",
    sha256 = "bc06e2de7dac791e5aea1a9e5a42ccc7fdf670c80f7aca690c39bfafa3dd7df6",
    strip_prefix = "googleapis-803b23488149d2d00725edc52d8a6498b6ab0fec",
    urls = [
        "https://github.com/googleapis/googleapis/archive/803b23488149d2d00725edc52d8a6498b6ab0fec.zip",
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
    sha256 = "b78f77458e77162f45b4564d6b20b6f92f56431ed59eaaab09e7819d1d850313",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.53.0/rules_go-v0.53.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.53.0/rules_go-v0.53.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.24.0")

http_archive(
    name = "bazel_gazelle",
    sha256 = "5d80e62a70314f39cc764c1c3eaa800c5936c9f1ea91625006227ce4d20cd086",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.42.0/bazel-gazelle-v0.42.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.42.0/bazel-gazelle-v0.42.0.tar.gz",
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
    sha256 = "3e114424a5c7e4fd43e0133cc6ecdfe54e45ae8affa14fadd839f29901424043",
    strip_prefix = "rules_shell-0.4.0",
    url = "https://github.com/bazelbuild/rules_shell/releases/download/v0.4.0/rules_shell-v0.4.0.tar.gz",
)

load("@rules_shell//shell:repositories.bzl", "rules_shell_dependencies", "rules_shell_toolchains")

rules_shell_dependencies()

rules_shell_toolchains()

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "573345c2039889a4001b9933a7ebde8dcaf910c47787993aecccebc3117a4425",
    strip_prefix = "buildtools-8.0.3",
    urls = ["https://github.com/bazelbuild/buildtools/archive/v8.0.3.tar.gz"],
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
