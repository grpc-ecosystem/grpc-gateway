"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "fa0b905032d49a621679e7318875736e451895a1417d992fbbebd27f82b83c38",
        strip_prefix = "buildtools-10.1.0",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v10.1.0.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
