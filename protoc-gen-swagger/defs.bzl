def _collect_includes(srcs):
    includes = ["."]
    for src in srcs:
        include = src.dirname
        if include and not include in includes:
            includes += [include]

    return includes

def _run_proto_gen_swagger(direct_proto_srcs, transitive_proto_srcs, actions, protoc, protoc_gen_swagger):
    swagger_files = []
    for proto in direct_proto_srcs:
      swagger_file = actions.declare_file(
          "%s.swagger.json" % proto.basename[:-len(".proto")],
          sibling = proto,
      )

      args = actions.args()
      args.add("--plugin=%s" % protoc_gen_swagger.path)
      args.add("--swagger_out=logtostderr=true:%s" % swagger_file.dirname)
      args.add("-Iexternal/com_google_protobuf/src")
      args.add("-Iexternal/com_github_googleapis_googleapis")
      args.add(["-I%s" % include for include in _collect_includes(direct_proto_srcs + transitive_proto_srcs)])
      args.add(proto.basename)

      actions.run(
          executable = protoc,
          inputs = direct_proto_srcs + transitive_proto_srcs + [protoc_gen_swagger],
          outputs = [swagger_file],
          arguments = [args],
      )

      swagger_files.append(swagger_file)

    return swagger_files

def _proto_gen_swagger_impl(ctx):
    proto = ctx.attr.proto.proto

    return struct(
        files=depset(
            _run_proto_gen_swagger(
                direct_proto_srcs = proto.direct_sources,
                transitive_proto_srcs = ctx.files._well_known_protos + proto.transitive_sources.to_list(),
                actions = ctx.actions,
                protoc = ctx.executable._protoc,
                protoc_gen_swagger = ctx.executable._protoc_gen_swagger,
            )
        )
    )

protoc_gen_swagger = rule(
    attrs = {
        "proto": attr.label(
            allow_rules = ["proto_library"],
            mandatory = True,
            providers = ['proto'],
        ),
        "_protoc": attr.label(
            default = "@com_google_protobuf//:protoc",
            executable = True,
            cfg = "host",
        ),
        "_well_known_protos": attr.label(
            default = "@com_google_protobuf//:well_known_protos",
            allow_files = True,
        ),
        "_protoc_gen_swagger": attr.label(
            default = Label("@grpc_ecosystem_grpc_gateway//protoc-gen-swagger:protoc-gen-swagger"),
            executable = True,
            cfg = "host",
        ),
    },
    implementation = _proto_gen_swagger_impl,
)
