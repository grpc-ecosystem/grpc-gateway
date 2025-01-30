"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "0063f317e135481783f3dc14c82bc15e0bf873c5e9aeece63b4f94d151aeb09f",
        strip_prefix = "buildtools-8.0.2",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v8.0.2.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
