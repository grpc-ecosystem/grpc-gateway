"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "e6de6eb19a368efe1f56549c6afe9f25dbcee818161865ee703081307581ef4b",
        strip_prefix = "buildtools-8.5.1",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v8.5.1.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
