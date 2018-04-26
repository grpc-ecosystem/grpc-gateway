def swagger_codegen(name, spec, package_name, outputs, out_dir=".", **kwarg):
  out_dir_spec = "$(@D)"
  if out_dir != ".":
    out_dir_spec += "/" + out_dir
    for out in outputs:
      if not out.startswith(out_dir + "/"):
        fail("%s must reside in the out_dir %s" % (out, out_dir), "outputs")

  native.genrule(
      name = name,
      srcs = [spec],
      outs = outputs,
      cmd = " ".join([
          "$(location //examples/clients:swagger-codegen-cli)",
          "generate",
          "-i", "$<",
          "-l", "go",
          "-o", out_dir_spec,
          "--additional-properties",
          "packageName=%s" % package_name,
      ]),
      tools = [
          "//examples/clients:swagger-codegen-cli",
      ],
      **kwarg
  )

def gen_file_integrity_test(name, generated_files, generated_dir, golden_dir, **kwargs):
  golden_prefix = _prefix(golden_dir)
  generated_prefix = _prefix(generated_dir)

  args = []
  data = []
  for gen in generated_files:
    if not gen.startswith(generated_prefix):
      fail("%s must reside in dir %s" % (gen, generated_dir))
    else:
      golden = golden_prefix + gen[len(generated_prefix):]
      data.extend([gen, golden])
      args.extend(["$(location %s)" % f for f in [gen, golden]])

  native.sh_test(
      name = name,
      srcs = ["//examples/clients:diff_test.sh"],
      args = args,
      data = data,
      **kwargs
  )

def _prefix(dir):
  if dir == ".":
    return ""
  return dir + "/"
