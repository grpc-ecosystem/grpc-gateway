workspace(name = "grpc_ecosystem_grpc_gateway")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# master as of 2018-11-09. This should be updated when a new release is tagged.
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "c23db3b50b8822e153bc5accfea75baeecedb481a162391c4f3b9aec451e34b4",
    strip_prefix = "rules_go-109c520465fcb418f2c4be967f3744d959ad66d3",
    urls = [
        "https://github.com/bazelbuild/rules_go/archive/109c520465fcb418f2c4be967f3744d959ad66d3.tar.gz",
    ],
)

# master as of 2018-11-09. This should be updated when a new release is tagged.
http_archive(
    name = "bazel_gazelle",
    sha256 = "a3af4a61d7b2f2c5386761f94a21f474871a32f0e94b13f08824248c4df25229",
    strip_prefix = "bazel-gazelle-7b1e3c6eb5447c6647955fc93c012635f274f0f0",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/archive/7b1e3c6eb5447c6647955fc93c012635f274f0f0.tar.gz",
    ],
)

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "e4c83a7a5d0712e2cea2077112a5eb6bb1af75a84e34c8c9b77330e322966b8b",
    strip_prefix = "buildtools-e90e7cc6ef3e6d08d4ca8a982935c3eed638e058",
    url = "https://github.com/bazelbuild/buildtools/archive/e90e7cc6ef3e6d08d4ca8a982935c3eed638e058.tar.gz",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

# Also define in Gopkg.toml
go_repository(
    name = "org_golang_google_genproto",
    commit = "383e8b2c3b9e36c4076b235b32537292176bae20",
    importpath = "google.golang.org/genproto",
)

# Also define in Gopkg.toml
go_repository(
    name = "com_github_rogpeppe_fastuuid",
    commit = "6724a57986aff9bff1a1770e9347036def7c89f6",
    importpath = "github.com/rogpeppe/fastuuid",
)

# Also define in Gopkg.toml
go_repository(
    name = "com_github_go_resty_resty",
    commit = "f8815663de1e64d57cdd4ee9e2b2fa96977a030e",
    importpath = "github.com/go-resty/resty",
)

# Also define in Gopkg.toml
go_repository(
    name = "com_github_ghodss_yaml",
    commit = "0ca9ea5df5451ffdf184b4428c902747c2c11cd7",
    importpath = "github.com/ghodss/yaml",
)

# Also define in Gopkg.toml
go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "eb3733d160e74a9c7e442f435eb3bea458e1d19f",
    importpath = "gopkg.in/yaml.v2",
)

go_rules_dependencies()

go_register_toolchains()

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
