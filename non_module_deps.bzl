"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "051951c10ff8addeb4f10be3b0cf474b304b2ccd675f2cc7683cdd9010320ca9",
        strip_prefix = "buildtools-7.3.1",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v7.3.1.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
