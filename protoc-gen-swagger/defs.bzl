def _collect_includes(srcs):
    includes = ["."]
    for src in srcs:
        include = src.dirname
        if include and not include in includes:
            includes += [include]

    return includes

def _run_proto_gen_swagger(proto_service, transitive_proto_srcs, actions, protoc, protoc_gen_swagger):
    swagger_file = actions.declare_file(
        "%s.swagger.json" % proto_service.basename[:-len(".proto")],
        sibling = proto_service,
    )

    args = actions.args()
    args.add("--plugin=%s" % protoc_gen_swagger.path)
    args.add("--swagger_out=logtostderr=true:%s" % swagger_file.dirname)
    args.add("-Iexternal/com_google_protobuf/src")
    args.add("-Iexternal/com_github_googleapis_googleapis")
    args.add(["-I%s" % include for include in _collect_includes(transitive_proto_srcs)])
    args.add(proto_service.basename)

    actions.run(
        executable = protoc,
        inputs = transitive_proto_srcs + [protoc_gen_swagger],
        outputs = [swagger_file],
        arguments = [args],
    )

    return swagger_file

def _proto_gen_swagger_impl(ctx):
    transitive_proto_srcs = depset([ctx.file.proto_service, ctx.file._annotations] + ctx.files._well_known_protos)
    for dep in ctx.attr.deps:
        transitive_proto_srcs = depset(transitive=[transitive_proto_srcs, dep.proto.transitive_sources])

    return struct(
        files=depset([
            _run_proto_gen_swagger(
                ctx.file.proto_service,
                transitive_proto_srcs.to_list(),
                ctx.actions,
                ctx.executable._protoc,
                ctx.executable._protoc_gen_swagger,
            )
        ])
    )

protoc_gen_swagger = rule(
    attrs = {
        "proto_service": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "deps": attr.label_list(
            allow_rules = ["proto_library"],
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
        "_annotations": attr.label(
            default = Label("@com_github_googleapis_googleapis//google/api:annotations.proto"),
            allow_single_file = True,
        ),
    },
    implementation = _proto_gen_swagger_impl,
)
