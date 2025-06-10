"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "53119397bbce1cd7e4c590e117dcda343c2086199de62932106c80733526c261",
        strip_prefix = "buildtools-8.2.1",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v8.2.1.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
