# Experimental implementation of Bazel rules to generate files with grpc-gateway.
# The APIs are not stable. They are subject to change without notice.

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//third_party/protobuf:protobuf.bzl", "proto_gen")

def go_proto_library(name,
                     srcs,
                     deps=[],
                     includes=[],
                     protoc="@com_github_google_protobuf//:protoc",
                     go_plugin="@com_github_golang_protobuf//:go_plugin",
                     genopts=[],
                     pkgmap={},
                     go_deps=[],
                     go_extra_srcs=[],
                     go_extra_library=None,
                     **kwargs):
  """Experimental support of protoc-gen-go

  Compiles Protocol Buffers definitions into Go source codes.

  This rule is experimental and a subject to change without notice.

  Args:
    name: A unique name for this rule
    srcs: Source files to be compiled into Go
    includes: A list of include directories to be passed to Protocol Buffers
      compiler.
    protoc: A label of protoc
    go_plugin: A label of protoc-gen-go
    genopts: options to be passed to go_plugin
    pkgmap: custom mapping from protocol buffers import path to go importpath.
    go_deps: Extra dependencies to be passed to go_library
    go_extra_srcs: Extra Go source files to be compiled together with the
       generated files
    go_extra_library: Extra library to be passed to go_library
  """
  genfiles = []
  for s in srcs:
    if not s.endswith('.proto'):
      fail("non proto source file %s" % s, "srcs")
    out = s[:-len('.proto')] + ".pb.go"
    genfiles += [out]

  opts = ["M%s=%s" % (proto, pkgmap[proto]) for proto in pkgmap] + genopts
  proto_gen(
      name = name + "_genproto",
      srcs = srcs,
      deps = deps,
      includes = includes,
      protoc = protoc,
      plugin = go_plugin,
      plugin_language = "go",
      plugin_options = opts,
      outs = genfiles,
      visibility = ["//visibility:private"],
  )
  go_library(
      name = name,
      srcs = genfiles + go_extra_srcs,
      deps = set(go_deps + ["@com_github_golang_protobuf//:proto"]),
      library = go_extra_library,
      **kwargs
  )

def grpc_gateway_proto_library(
    name,
    srcs,
    deps=[],
    includes=[],
    protoc="@com_github_google_protobuf//:protoc",
    gw_plugin="//protoc-gen-grpc-gateway",
    genopts=[],
    pkgmap={},
    go_deps=[],
    go_extra_srcs=[],
    go_extra_library=None,
    **kwargs):
  """Experimental support of protoc-gen-grpc-gateway

  Compiles Protocol Buffers definitions into reverse proxy implementation
  with grpc-gateway.

  This rule is experimental and a subject to change without notice.

  Args:
    name: A unique name for this rule
    srcs: Source files to be compiled into Go
    includes: A list of include directories to be passed to Protocol Buffers
      compiler.
    protoc: A label of protoc
    gw_plugin: A label of protoc-gen-grpc-gateway
    genopts: options to be passed to gw_plugin
    pkgmap: custom mapping from protocol buffers import path to go importpath.
    go_deps: Extra dependencies to be passed to go_library
    go_extra_srcs: Extra Go source files to be compiled together with the
       generated files
    go_extra_library: Extra library to be passed to go_library
  """
  genfiles = []
  for s in srcs:
    if not s.endswith('.proto'):
      fail("non proto source file %s" % s, "srcs")
    out = s[:-len('.proto')] + ".pb.gw.go"
    genfiles += [out]

  opts = ["M%s=%s" % (proto, pkgmap[proto]) for proto in pkgmap]
  opts += genopts + ["logtostderr=true"]
  proto_gen(
      name = name + "_genproto",
      srcs = srcs,
      deps = deps,
      includes = includes,
      protoc = protoc,
      plugin = gw_plugin,
      plugin_language = "grpc-gateway",
      plugin_options = opts,
      outs = genfiles,
      visibility = ["//visibility:private"],
  )
  deps = set(go_deps) + [
      "//runtime:go_default_library",
      "//utilities:go_default_library",
      "@com_github_golang_protobuf//:proto",
      "@org_golang_x_net//:context",
      "@org_golang_google_grpc//:codes",
      "@org_golang_google_grpc//:go_default_library",
      "@org_golang_google_grpc//:grpclog",
  ]
  go_library(
      name = name,
      srcs = genfiles + go_extra_srcs,
      deps = deps,
      library = go_extra_library,
      **kwargs
  )

def grpc_swagger_proto_library(name,
                               srcs,
                               deps=[],
                               includes=[],
                               protoc="@com_github_google_protobuf//:protoc",
                               sw_plugin="//protoc-gen-swagger",
                               genopts=[]):
  """Experimental support of protoc-gen-swagger

  Compiles Protocol Buffers definitions into Swagger definitions.
  This rule is experimental and a subject to change without notice.

  Args:
    name: A unique name for this rule
    srcs: Source files to be compiled into Go
    includes: A list of include directories to be passed to Protocol Buffers
      compiler.
    protoc: A label of protoc
    sw_plugin: A label of protoc-gen-swagger
    genopts: options to be passed to sw_plugin
  """
  genfiles = []
  for s in srcs:
    if not s.endswith('.proto'):
      fail("non proto source file %s" % s, "srcs")
    out = s[:-len('.proto')] + ".swagger.json"
    genfiles += [out]

  proto_gen(
      name = name,
      srcs = srcs,
      deps = deps,
      includes = includes,
      protoc = protoc,
      plugin = sw_plugin,
      plugin_language = "swagger",
      plugin_options = genopts + ["logtostderr=true"],
      outs = genfiles,
      visibility = ["//visibility:private"],
  )
