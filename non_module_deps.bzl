"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "91727456f1338f511442c50a8d827ae245552642d63de2bc832e6d27632ec300",
        strip_prefix = "buildtools-8.0.1",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v8.0.1.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
