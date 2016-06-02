# Copied a part of protobuf.bzl from github.com/google/protobuf.
# Locally modified.

def _GetPath(ctx, path):
  if ctx.label.workspace_root:
    return ctx.label.workspace_root + '/' + path
  else:
    return path

def _GenDir(ctx):
  if not ctx.attr.includes:
    return ctx.label.workspace_root
  if not ctx.attr.includes[0]:
    return _GetPath(ctx, ctx.label.package)
  if not ctx.label.package:
    return _GetPath(ctx, ctx.attr.includes[0])
  return _GetPath(ctx, ctx.label.package + '/' + ctx.attr.includes[0])

def _proto_gen_impl(ctx):
  """General implementation for generating protos"""
  srcs = ctx.files.srcs
  deps = []
  deps += ctx.files.srcs
  gen_dir = _GenDir(ctx)
  if gen_dir:
    import_flags = ["-I" + gen_dir, "-I" + ctx.var["GENDIR"] + "/" + gen_dir]
  else:
    import_flags = ["-I."]

  for dep in ctx.attr.deps:
    import_flags += dep.proto.import_flags
    deps += dep.proto.deps

  args = []
  if ctx.attr.gen_cc:
    args += ["--cpp_out=" + ctx.var["GENDIR"] + "/" + gen_dir]
  if ctx.attr.gen_py:
    args += ["--python_out=" + ctx.var["GENDIR"] + "/" + gen_dir]

  if ctx.executable.plugin:
    plugin = ctx.executable.plugin
    lang = ctx.attr.plugin_language
    if not lang and plugin.basename.startswith('protoc-gen-'):
      lang = plugin.basename[len('protoc-gen-'):]
    if not lang:
      fail("cannot infer the target language of plugin", "plugin_language")

    outdir = ctx.var["GENDIR"] + "/" + gen_dir
    if ctx.attr.plugin_options:
      outdir = ",".join(ctx.attr.plugin_options) + ":" + outdir
    args += ["--plugin=protoc-gen-%s=%s" % (lang, plugin.path)]
    args += ["--%s_out=%s" % (lang, outdir)]

  if args:
    ctx.action(
        inputs=srcs + deps,
        outputs=ctx.outputs.outs,
        arguments=args + import_flags + [s.path for s in srcs],
        executable=ctx.executable.protoc,
        mnemonic="ProtoCompile",
    )

  return struct(
      proto=struct(
          srcs=srcs,
          import_flags=import_flags,
          deps=deps,
      ),
  )

proto_gen = rule(
    attrs = {
        "srcs": attr.label_list(allow_files = True),
        "deps": attr.label_list(providers = ["proto"]),
        "includes": attr.string_list(),
        "protoc": attr.label(
            cfg = HOST_CFG,
            executable = True,
            single_file = True,
            mandatory = True,
        ),
        "plugin": attr.label(
            cfg = HOST_CFG,
            allow_files = True,
            executable = True,
        ),
        "plugin_language": attr.string(),
        "plugin_options": attr.string_list(),
        "gen_cc": attr.bool(),
        "gen_py": attr.bool(),
        "outs": attr.output_list(),
    },
    output_to_genfiles = True,
    implementation = _proto_gen_impl,
)
