"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "f3b800e9f6ca60bdef3709440f393348f7c18a29f30814288a7326285c80aab9",
        strip_prefix = "buildtools-8.5.1",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v8.5.1.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
