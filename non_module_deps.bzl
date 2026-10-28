"""Module extension for non-module dependencies."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(
        # buildifier: disable=unused-variable
        mctx):
    # TODO(bazelbuild/buildtools#1204): Remove when available as module.
    http_archive(
        name = "com_github_bazelbuild_buildtools",
        sha256 = "5af90068e9e8812abdc455ad04cc474470bccc8bbabc8b571e3f58a66799680f",
        strip_prefix = "buildtools-10.2.0",
        urls = ["https://github.com/bazelbuild/buildtools/archive/v10.2.0.tar.gz"],
    )
    http_archive(
        name = "com_github_google_safeopen",
        patch_cmds = [
            "find . -name BUILD.bazel -type f -print0 | xargs -0 sed -i 's|@go_sys//|@org_golang_x_sys//|g'",
        ],
        sha256 = "0f10a4457d47fef681d60b88b2f771b0248ca37ee22b448ce6c46920faf9d1cd",
        strip_prefix = "safeopen-43626d6f468560fc611ad54f1f17f3a455dbe92d",
        urls = ["https://github.com/google/safeopen/archive/43626d6f468560fc611ad54f1f17f3a455dbe92d.zip"],
    )

non_module_deps = module_extension(
    implementation = _non_module_deps_impl,
)
