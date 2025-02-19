"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "573345c2039889a4001b9933a7ebde8dcaf910c47787993aecccebc3117a4425",
        strip_prefix = "buildtools-8.0.3",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v8.0.3.tar.gz"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
