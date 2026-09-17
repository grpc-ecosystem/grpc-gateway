"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "f495fe290cf2a009e80b20d5623c7756890a8a45d83fc93d44a31ab92060c752",
        strip_prefix = "buildtools-10.0.0",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v10.0.0.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
