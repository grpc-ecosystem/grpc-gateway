"""Generated an open-api v3 spec for a grpc api spec.

Reads the api spec in protobuf format and generate an open-api spec.
Optionally applies settings from the grpc-service configuration.
"""

load("@rules_proto//proto:defs.bzl", "ProtoInfo")

# TODO(yannic): Replace with |proto_common.direct_source_infos| when
# https://github.com/bazelbuild/rules_proto/pull/22 lands.
def _direct_source_infos(proto_info, provided_sources = []):
    """Returns sequence of `ProtoFileInfo` for `proto_info`'s direct sources.

    Files that are both in `proto_info`'s direct sources and in
    `provided_sources` are skipped. This is useful, e.g., for well-known
    protos that are already provided by the Protobuf runtime.

    Args:
      proto_info: An instance of `ProtoInfo`.
      provided_sources: Optional. A sequence of files to ignore.
          Usually, these files are already provided by the
          Protocol Buffer runtime (e.g. Well-Known protos).

    Returns: A sequence of `ProtoFileInfo` containing information about
        `proto_info`'s direct sources.
    """

    source_root = proto_info.proto_source_root
    if "." == source_root:
        return [struct(file = src, import_path = src.path) for src in proto_info.check_deps_sources.to_list()]

    offset = len(source_root) + 1  # + '/'.

    infos = []
    for src in proto_info.check_deps_sources.to_list():
        # TODO(yannic): Remove this hack when we drop support for Bazel < 1.0.
        local_offset = offset
        if src.root.path and not source_root.startswith(src.root.path):
            # Before Bazel 1.0, `proto_source_root` wasn't guaranteed to be a
            # prefix of `src.path`. This could happened, e.g., if `file` was
            # generated (https://github.com/bazelbuild/bazel/issues/9215).
            local_offset += len(src.root.path) + 1  # + '/'.
        infos.append(struct(file = src, import_path = src.path[local_offset:]))

    return infos

def _run_proto_gen_openapiv3(
        actions,
        proto_info,
        target_name,
        transitive_proto_srcs,
        protoc,
        protoc_gen_openapiv2,
        single_output,
        allow_delete_body,
        grpc_api_configuration,
        json_names_for_fields,
        repeated_path_param_separator,
        include_package_in_tags,
        fqn_for_openapi_name,
        openapi_naming_strategy,
        use_go_templates,
        go_template_args,
        ignore_comments,
        remove_internal_comments,
        disable_default_errors,
        disable_service_tags,
        enums_as_ints,
        omit_enum_default_value,
        output_format,
        simple_operation_ids,
        proto3_optional_nullable,
        openapi_configuration,
        generate_unbound_methods,
        visibility_restriction_selectors,
        use_allof_for_refs,
        disable_default_responses,
        enable_rpc_deprecation,
        expand_slashed_path_patterns,
        preserve_rpc_order,
        generate_x_go_type):
    args = actions.args()

    args.add("--plugin", "protoc-gen-openapiv3=%s" % protoc_gen_openapiv2.path)

    extra_inputs = []

    if output_format:
        args.add("--openapiv3_opt", "output_format=%s" % output_format)

    for visibility_restriction_selector in visibility_restriction_selectors:
        args.add("--openapiv3_opt", "visibility_restriction_selectors=%s" % visibility_restriction_selector)

    proto_file_infos = _direct_source_infos(proto_info)

    # TODO(yannic): Use |proto_info.transitive_descriptor_sets| when
    # https://github.com/bazelbuild/bazel/issues/9337 is fixed.
    args.add_all(proto_info.transitive_proto_path, format_each = "--proto_path=%s")

openapi_file = actions.declare_file("%s.swagger.json" % target_name)
args.add("--openapiv3_out", openapi_file.dirname)

openapi_files = []
for proto_file_info in proto_file_infos:
    file_name = "%s.openapiv3.json" % proto_file_info.import_path[:-len(".proto")]
openapi_file = actions.declare_file(
    "_virtual_imports/%s/%s" % (target_name, file_name),
)

file_args = actions.args()

offset = len(file_name) + 1  # + '/'.
file_args.add("--openapiv3_out", openapi_file.path[:-offset])

file_args.add(proto_file_info.import_path)

actions.run(
    executable = protoc,
    tools = [protoc_gen_openapiv3],
    inputs = depset(
        direct = extra_inputs,
        transitive = [transitive_proto_srcs],
    ),
    outputs = [openapi_file],
    arguments = [args, file_args],
)
openapi_files.append(openapi_file)

return openapi_files

def _proto_gen_openapiv3_impl(ctx):
    proto = ctx.attr.proto[ProtoInfo]
    return [
        DefaultInfo(
            files = depset(
                _run_proto_gen_openapiv3(
                    actions = ctx.actions,
                    proto_info = proto,
                    target_name = ctx.attr.name,
                    transitive_proto_srcs = depset(
                        direct = ctx.files._well_known_protos,
                        transitive = [proto.transitive_sources],
                    ),
                    protoc = ctx.executable._protoc,
                    protoc_gen_openapiv2 = ctx.executable._protoc_gen_openapi,
                    output_format = ctx.attr.output_format,
                    visibility_restriction_selectors = ctx.attr.visibility_restriction_selectors,
                ),
            ),
        ),
    ]

protoc_gen_openapiv3 = rule(
    attrs = {
        "proto": attr.label(
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "output_format": attr.string(
            default = "json",
            mandatory = False,
            values = ["json", "yaml"],
            doc = "output content format. Allowed values are: `json`, `yaml`",
        ),
        "visibility_restriction_selectors": attr.string_list(
            mandatory = False,
            doc = "list of `google.api.VisibilityRule` visibility labels to include" +
                  " in the generated output when a visibility annotation is defined." +
                  " Repeat this option to supply multiple values. Elements without" +
                  " visibility annotations are unaffected by this setting.",
        ),
        "_protoc": attr.label(
            default = "@com_google_protobuf//:protoc",
            executable = True,
            cfg = "exec",
        ),
        "_well_known_protos": attr.label(
            default = "@com_google_protobuf//:well_known_type_protos",
            allow_files = True,
        ),
        "_protoc_gen_openapi": attr.label(
            default = Label("//protoc-gen-openapiv3:protoc-gen-openapiv3"),
            executable = True,
            cfg = "exec",
        ),
    },
    implementation = _proto_gen_openapiv3_impl,
)
