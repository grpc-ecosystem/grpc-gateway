workspace(name = "grpc_ecosystem_grpc_gateway")

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.12.1/rules_go-0.12.1.tar.gz",
    sha256 = "8b68d0630d63d95dacc0016c3bb4b76154fe34fca93efd65d1c366de3fcb4294",
)

http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.10.1/bazel-gazelle-0.10.1.tar.gz",
    sha256 = "d03625db67e9fb0905bbd206fa97e32ae9da894fe234a493e7517fd25faec914",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repository")
load("//:repositories.bzl", "repositories")
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

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

repositories()

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
