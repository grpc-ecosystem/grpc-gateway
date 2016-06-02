workspace(name = "com_github_gengo_grpc_gateway")

#git_repository(
#    name = "io_bazel_rules_go",
#    remote = "https://github.com/bazelbuild/rules_go.git",
#    commit = "bb7e87b63d33c24770cb21974bcc8ae00f8f9d07",
#)
local_repository(
    name = "io_bazel_rules_go",
    path = "/Users/yugui/dev/bazel/rules_go",
)
load("@io_bazel_rules_go//go:def.bzl", "go_repositories")

go_repositories()

git_repository(
    name = "com_google_developers_protocol_buffers",
    remote = "https://github.com/google/protobuf.git",
    tag = "v3.0.0-beta-3",
)

new_git_repository(
    name = "com_github_golang_protobuf",
    remote = "https://github.com/golang/protobuf.git",
    commit = "7cc19b78d562895b13596ddce7aafb59dd789318",
    build_file = "goprotobuf.BUILD",
)

new_git_repository(
    name = "org_golang_x_net",
    remote = "https://github.com/golang/net.git",
    commit = "2a35e686583654a1b89ca79c4ac78cb3d6529ca3",
    build_file = "gonet.BUILD",
)

new_git_repository(
    name = "org_golang_google_grpc",
    remote = "https://github.com/grpc/grpc-go.git",
    commit = "8213ee577a465c1f314d85748fb29e4eeed59baf",
    build_file = "grpc-go.BUILD",
)

new_git_repository(
    name = "com_github_golang_glog",
    remote = "https://github.com/golang/glog.git",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    build_file = "glog.BUILD",
)

new_git_repository(
    name = "com_github_rogpeppe_fastuuid",
    remote = "https://github.com/rogpeppe/fastuuid.git",
    commit = "6724a57986aff9bff1a1770e9347036def7c89f6",
    build_file = "fastuuid.BUILD",
)
