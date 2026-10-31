load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        importpath = "golang.org/x/xerrors",
        sum = "h1:GGU+dLjvlC3qDwqYgL6UgRmHXhOOgns0bZu2Ty5mm6U=",
        version = "v0.0.0-20220411194840-2f41105eb62f",
    )
    go_repository(
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:+jPbCzuYxYrcoOvCecYUFDj+XJQRmzIhJmlbpps21Zk=",
        version = "v0.8.2",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:pxW6RcqyfI9/kWtOwnv/G+AzdKuy2ZrqINhenH4HyNs=",
        version = "v1.4.1-0.20240526193622-a339e1f7089c",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:aBangftG7EVZoUb69Os8IaYg++6uMOdKK83QtkkvJik=",
        version = "v0.0.0-20260202195803-dba9d589def2",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:hbG2kr4RuFj222B6+7T83thSPqLjwBIfQawTkC++2HA=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_envoy",
        importpath = "github.com/envoyproxy/go-control-plane/envoy",
        sum = "h1:SWe0x6yfglnxuvOiYgTTnNq7QD/yvqthOh1RBI8Bj8w=",
        version = "v1.39.1-0.20260819172001-e6e3fd93e4be",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_ratelimit",
        importpath = "github.com/envoyproxy/go-control-plane/ratelimit",
        sum = "h1:/G9QYbddjL25KvtKTv3an9lx6VBE2cnb8wp1vEGNYGI=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:MVQghNeW+LZcmXe7SY1V36Z+WFMDjpqGAGacLe2T0ds=",
        version = "v1.3.3",
    )
    go_repository(
        name = "com_github_go_jose_go_jose_v4",
        importpath = "github.com/go-jose/go-jose/v4",
        sum = "h1:moDMcTHmvE6Groj34emNPLs/qtYXRVcd6S7NHbHz3kA=",
        version = "v4.1.4",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:tG4xh9yMsRCAiodLVTxyrkzSZ9+o0L1Kg/+cPVcbP/8=",
        version = "v1.4.4",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:DrW6hGnjIhtvhOIiAKT6Psh/Kd/ldepEa81DKeiRJ5I=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=",
        version = "v1.5.4",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_detectors_gcp",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp",
        sum = "h1:bN1gA3of5bXtbnLsRPrwfmbbe7A5UWFlcTHseujLnpc=",
        version = "v1.35.0",
    )
    go_repository(
        name = "com_github_planetscale_vtprotobuf",
        importpath = "github.com/planetscale/vtprotobuf",
        sum = "h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=",
        version = "v0.6.1-0.20240319094008-0393e58bdf10",
    )
    go_repository(
        name = "com_github_spiffe_go_spiffe_v2",
        importpath = "github.com/spiffe/go-spiffe/v2",
        sum = "h1:eXZMLsu+3MLEPJyGJkolqtVrteZfQdUpOWj6LTiDl/E=",
        version = "v2.8.1",
    )
    go_repository(
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
        version = "v1.4.13",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:pDUj4QMoPejqq20dK0Pg2N4yG9zIkYGdBtwLoEkH9Zs=",
        version = "v0.9.0",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:A2jO8jwOugrrovveCWfj0KEZOfqiLgAcwjpHPhzIGw0=",
        version = "v0.25.3",
    )
    go_repository(
        name = "in_yaml_go_yaml_v3",
        importpath = "go.yaml.in/yaml/v3",
        sum = "h1:zcDfvjXEqTQmpjniUBLc91xYJMB2c4QOFxEF/LCEckM=",
        version = "v3.0.6",
    )
    go_repository(
        name = "io_opentelemetry_go_auto_sdk",
        importpath = "go.opentelemetry.io/auto/sdk",
        sum = "h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=",
        version = "v1.2.1",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_detectors_gcp",
        importpath = "go.opentelemetry.io/contrib/detectors/gcp",
        sum = "h1:9jR0ZPRok9ryaOQ2Wx8rg5F7Aon59mxrqbVI60/vlBk=",
        version = "v1.45.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:pdrWmLHofpubmArBv1LgFSv1Z0Ie/ppdZzu+kUN5EeU=",
        version = "v1.45.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:7Eg1uH7CJ5cXv9is6tnBe1FI6rj1nwUdbFypRm3br/M=",
        version = "v1.45.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:4VVSMgQ83dUgW2aoX5f6JgLvHwIvzcuLnF9lUdCSpCw=",
        version = "v1.45.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:oVFszMfyj1Am6s24Vtc7wBb8BKLcwepJjNEYILuiE3o=",
        version = "v1.45.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:l/mP6Uv7oNO7/TblbhpbgMidxhq1uO/rPsikOyVhxag=",
        version = "v1.45.0",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:e789sLZjeTdc02ZCQGTSzHOu0fR73+dwPQVKUp67Z88=",
        version = "v0.0.0-20261026165753-2ffab32bce7e",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:6WNWHZk91dp4gZ25642OALkuG1UVNZbwENJD19DnB8k=",
        version = "v0.0.0-20261019180522-ce908d29299b",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:BPBCfuVRQ1qvBhKDOq9uvQllMvO59rKsYhXpp7Dd8vw=",
        version = "v1.85.0",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:pJOKDDOyeXErUroCihFAd5LQuwXBSpVnKGrj5o/fwxc=",
        version = "v1.36.12",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:COkYLr4k7nDI4r6QtMML2AfzeIswWHAJthLS6K/M/54=",
        version = "v0.58.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:HgI6W1OolKK6tcV74is+rl9SMT0njsHlp0TbQiGXI/A=",
        version = "v0.0.0-20261014143641-f0251eb4aa3e",
    )
    go_repository(
        name = "org_golang_x_exp_typeparams",
        importpath = "golang.org/x/exp/typeparams",
        sum = "h1:1P7xPZEwZMoBoz0Yze5Nx2/4pxj6nw9ZqHWXqP0iRgQ=",
        version = "v0.0.0-20231108232855-2478ac86f678",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:ICeL3t2L30ljSPHHxDnR4fqVcBt/zR2zJAb9juavmQg=",
        version = "v0.42.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:18ELalfM5+XNRwzM2nYqXsZQ9ByS1uhj2O/U7+5jqRY=",
        version = "v0.62.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=",
        version = "v0.36.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:iqiDrGoDgyuvAKrhLQthUEzU/R378N78gXbBjxed7oA=",
        version = "v0.24.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:XbzkgYJHdqh/8m2Uu0W/dQv8nktxx4BFHp1M0gROTXA=",
        version = "v0.49.0",
    )
    go_repository(
        name = "org_golang_x_telemetry",
        importpath = "golang.org/x/telemetry",
        sum = "h1:wnveMBKHwx8/nIb8WYT5P9vtefdXOvIWxzycCkc5n6Q=",
        version = "v0.0.0-20261009151841-97e8ea77aa2a",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:cpyOwQkr2Nzxe2p6KnrTG7Ei2wtzjsuSjhpXXepZ/BM=",
        version = "v0.47.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:1QivrlhwAsnMOcqgOdl9my7s9OyZl8O5/UGUTJVXoWc=",
        version = "v0.43.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:ryxg4U63oaUYk7uOLTF37n43f4qgtGvqV4Ioit1BaEs=",
        version = "v0.52.0",
    )
    go_repository(
        name = "org_golang_x_tools_go_expect",
        importpath = "golang.org/x/tools/go/expect",
        sum = "h1:jpBZDwmgPhXsKZC6WhL20P4b/wmnpsEAGHaNy0n/rJM=",
        version = "v0.1.1-deprecated",
    )
    go_repository(
        name = "org_golang_x_tools_go_packages_packagestest",
        importpath = "golang.org/x/tools/go/packages/packagestest",
        sum = "h1:1h2MnaIAIXISqTFKdENegdpAgUXz6NrPEsbIeWaBRvM=",
        version = "v0.1.1-deprecated",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:VbpOemQlsSMrYmn7T2OUvQ4dqxQXU+ouZFQsZOx50z4=",
        version = "v0.17.0",
    )
